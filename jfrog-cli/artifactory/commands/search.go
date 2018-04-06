package commands

import (
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/artifactory/utils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/artifactory/utils/spec"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/utils/config"
	clientutils "github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/artifactory/services/utils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/utils/log"
)

type SearchResult struct {
	Path       string            `json:"path,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

func Search(searchSpec *spec.SpecFiles, artDetails *config.ArtifactoryDetails) ([]SearchResult, error) {
	servicesManager, err := utils.CreateServiceManager(artDetails, false)
	if err != nil {
		return nil, err
	}
	log.Info("Searching artifacts...")
	var resultItems []clientutils.ResultItem
	for i := 0; i < len(searchSpec.Files); i++ {
		params, err := searchSpec.Get(i).ToArtifatorySearchParams()
		if err != nil {
			return nil, err
		}
		currentResultItems, err := servicesManager.Search(&clientutils.SearchParamsImpl{ArtifactoryCommonParams: params})
		if err != nil {
			return nil, err
		}
		resultItems = append(resultItems, currentResultItems...)
	}

	result := aqlResultToSearchResult(resultItems)
	clientutils.LogSearchResults(len(resultItems))
	return result, err
}

func aqlResultToSearchResult(aqlResult []clientutils.ResultItem) (result []SearchResult) {
	result = make([]SearchResult, len(aqlResult))
	for i, v := range aqlResult {
		tempResult := new(SearchResult)
		if v.Path != "." {
			tempResult.Path = v.Repo + "/" + v.Path + "/" + v.Name
		} else {
			tempResult.Path = v.Repo + "/" + v.Name
		}
		tempResult.Properties = make(map[string]string, len(v.Properties))
		for _, prop := range v.Properties {
			tempResult.Properties[prop.Key] = prop.Value
		}
		result[i] = *tempResult
	}
	return
}
