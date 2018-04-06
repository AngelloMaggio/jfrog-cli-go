package docker

import (
	"encoding/json"
	"errors"
	buildutils "github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/artifactory/utils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/artifactory/utils/buildinfo"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/artifactory"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/artifactory/services"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/artifactory/services/utils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/utils/errorutils"
	"io/ioutil"
	"path"
	"strings"
)

// Docker image build info builder
type Builder interface {
	Build() (*buildinfo.BuildInfo, error)
}

// Create instance of docker build info builder
func BuildInfoBuilder(image Image, targetRepo, buildName, buildNumber string, serviceManager *artifactory.ArtifactoryServicesManager) Builder {
	builder := &buildInfoBuilder{}
	builder.image = image
	builder.targetRepo = targetRepo
	builder.buildName = buildName
	builder.buildNumber = buildNumber
	builder.serviceManager = serviceManager
	return builder
}

type buildInfoBuilder struct {
	image          Image
	targetRepo     string
	buildName      string
	buildNumber    string
	serviceManager *artifactory.ArtifactoryServicesManager

	// internal fields
	imageId      string
	layers       []utils.ResultItem
	artifacts    []buildinfo.Artifacts
	dependencies []buildinfo.Dependencies
}

// Create build info for docker image
func (builder *buildInfoBuilder) Build() (*buildinfo.BuildInfo, error) {
	var err error
	builder.imageId, err = builder.image.Id()
	if err != nil {
		return nil, err
	}

	err = builder.updateArtifactsAndDependencies()
	if err != nil {
		return nil, err
	}

	_, err = builder.setBuildProperties()
	if err != nil {
		return nil, err
	}

	return builder.createBuildInfo()
}

// Search validate and create artifacts and dependencies of docker image.
func (builder *buildInfoBuilder) updateArtifactsAndDependencies() error {
	// Search for all the image layer to get the local path inside Artifactory (supporting virtual repos).
	searchResults, err := searchImageLayers(builder.imageId, path.Join(builder.targetRepo, builder.image.Path(), "*"), builder.serviceManager)
	if err != nil {
		return err
	}

	manifest, artifact, err := getManifest(builder.imageId, searchResults, builder.serviceManager)
	if err != nil {
		return err
	}
	builder.artifacts = append(builder.artifacts, artifact)

	configLayer, artifact, err := getConfigLayer(builder.imageId, searchResults, builder.serviceManager)
	if err != nil {
		return err
	}
	builder.artifacts = append(builder.artifacts, artifact)

	builder.layers = append(builder.layers, searchResults["manifest.json"])
	builder.layers = append(builder.layers, searchResults[digestToLayer(builder.imageId)])
	for i := 0; i < configLayer.getNumberLayers(); i++ {
		layerFileName := digestToLayer(manifest.Layers[i].Digest)
		item, layerExists := searchResults[layerFileName]
		if !layerExists {
			return errorutils.CheckError(errors.New("Could not find layer: " + layerFileName + "in Artifactory"))
		}

		if i < configLayer.getNumberOfDependentLayers() {
			builder.dependencies = append(builder.dependencies, item.ToDependency())
		}

		builder.artifacts = append(builder.artifacts, item.ToArtifact())
		builder.layers = append(builder.layers, item)
	}
	return nil
}

// Set build properties on docker image layers in Artifactory
func (builder *buildInfoBuilder) setBuildProperties() (int, error) {
	props, err := buildutils.CreateBuildProperties(builder.buildName, builder.buildNumber)
	if err != nil {
		return 0, err
	}
	return builder.serviceManager.SetProps(&services.SetPropsParamsImpl{Items: builder.layers, Props: props})
}

// Create docker build info
func (builder *buildInfoBuilder) createBuildInfo() (*buildinfo.BuildInfo, error) {
	imageProperties := map[string]string{}
	imageProperties["docker.image.id"] = builder.imageId
	imageProperties["docker.image.tag"] = builder.image.Tag()

	parentId, err := builder.image.ParentId()
	if err != nil {
		return nil, err
	}
	if parentId != "" {
		imageProperties["docker.image.parent"] = parentId
	}

	buildInfo := &buildinfo.BuildInfo{Modules: []buildinfo.Module{{
		Id:           builder.image.Name(),
		Properties:   imageProperties,
		Artifacts:    builder.artifacts,
		Dependencies: builder.dependencies,
	}}}
	return buildInfo, nil
}

// Download and read the manifest from Artifactory
func getManifest(imageId string, searchResults map[string]utils.ResultItem, serviceManager *artifactory.ArtifactoryServicesManager) (*manifest, buildinfo.Artifacts, error) {
	item := searchResults["manifest.json"]
	ioReaderCloser, err := serviceManager.ReadRemoteFile(item.GetItemRelativePath())
	if err != nil {
		return nil, buildinfo.Artifacts{}, err
	}
	defer ioReaderCloser.Close()
	content, err := ioutil.ReadAll(ioReaderCloser)
	if err != nil {
		return nil, buildinfo.Artifacts{}, err
	}

	var manifest manifest
	err = json.Unmarshal(content, &manifest)
	if errorutils.CheckError(err) != nil {
		return nil, buildinfo.Artifacts{}, err
	}

	// Check that the manifest ID is the right one
	if manifest.Config.Digest != imageId {
		return nil, buildinfo.Artifacts{}, errorutils.CheckError(errors.New("Found incorrect manifest.json file, expecting image ID: " + imageId))
	}

	artifact := buildinfo.Artifacts{Name: "manifest.json", Checksum: &buildinfo.Checksum{Sha1: item.Actual_Sha1, Md5: item.Actual_Md5}}
	return &manifest, artifact, nil
}

// Download and read the config layer from Artifactory
func getConfigLayer(imageId string, searchResults map[string]utils.ResultItem, serviceManager *artifactory.ArtifactoryServicesManager) (*configLayer, buildinfo.Artifacts, error) {
	item := searchResults[digestToLayer(imageId)]
	ioReaderCloser, err := serviceManager.ReadRemoteFile(item.GetItemRelativePath())
	if err != nil {
		return nil, buildinfo.Artifacts{}, err
	}
	defer ioReaderCloser.Close()
	content, err := ioutil.ReadAll(ioReaderCloser)
	if err != nil {
		return nil, buildinfo.Artifacts{}, err
	}

	var configLayer configLayer
	err = json.Unmarshal(content, &configLayer)
	if err != nil {
		return nil, buildinfo.Artifacts{}, err
	}

	artifact := buildinfo.Artifacts{Name: digestToLayer(imageId), Checksum: &buildinfo.Checksum{Sha1: item.Actual_Sha1, Md5: item.Actual_Md5}}
	return &configLayer, artifact, nil
}

// Search for image layers in Artifactory
func searchImageLayers(imageId, imagePathPattern string, serviceManager *artifactory.ArtifactoryServicesManager) (map[string]utils.ResultItem, error) {
	params := &utils.SearchParamsImpl{}
	params.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{}
	params.Pattern = imagePathPattern
	results, err := serviceManager.Search(params)
	if err != nil {
		return nil, err
	}
	resultMap := map[string]utils.ResultItem{}
	for _, v := range results {
		resultMap[v.Name] = v
	}

	// Validate image ID layer exists.
	if _, ok := resultMap[digestToLayer(imageId)]; !ok {
		return nil, errorutils.CheckError(errors.New("Found incorrect docker image, expecting image ID: " + imageId))
	}
	return resultMap, nil
}

// Digest of type sha256:30daa5c11544632449b01f450bebfef6b89644e9e683258ed05797abe7c32a6e to
// sha256__30daa5c11544632449b01f450bebfef6b89644e9e683258ed05797abe7c32a6e
func digestToLayer(digest string) string {
	return strings.Replace(digest, ":", "__", 1)
}

// Get the total number of layers from the config
func (configLayer *configLayer) getNumberLayers() int {
	layersNum := len(configLayer.History)
	for i := len(configLayer.History) - 1; i >= 0; i-- {
		if configLayer.History[i].EmptyLayer {
			layersNum--
		}
	}
	return layersNum
}

// Get the number of dependencies layers from the config
func (configLayer *configLayer) getNumberOfDependentLayers() int {
	layersNum := len(configLayer.History)
	newImageLayers := true
	for i := len(configLayer.History) - 1; i >= 0; i-- {
		if newImageLayers {
			layersNum--
		}

		if !newImageLayers && configLayer.History[i].EmptyLayer {
			layersNum--
		}

		createdBy := configLayer.History[i].CreatedBy
		if strings.Contains(createdBy, "ENTRYPOINT") || strings.Contains(createdBy, "MAINTAINER") {
			newImageLayers = false
		}
	}
	return layersNum
}

// To unmarshal config layer file
type configLayer struct {
	History []history `json:"history,omitempty"`
}

type history struct {
	Created    string `json:"created,omitempty"`
	CreatedBy  string `json:"created_by,omitempty"`
	EmptyLayer bool   `json:"empty_layer,omitempty"`
}

// To unmarshal manifest.json file
type manifest struct {
	Config manifestConfig `json:"config,omitempty"`
	Layers []layer        `json:"layers,omitempty"`
}

type manifestConfig struct {
	Digest string `json:"digest,omitempty"`
}

type layer struct {
	Digest string `json:"digest,omitempty"`
}
