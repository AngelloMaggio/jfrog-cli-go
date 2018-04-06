package commands

import (
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/artifactory/utils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/artifactory/utils/spec"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/utils/config"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/artifactory/services"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/utils/log"
)

// Moves the artifacts using the specified move pattern.
func Move(moveSpec *spec.SpecFiles, artDetails *config.ArtifactoryDetails) (successCount, failCount int, err error) {
	servicesManager, err := utils.CreateServiceManager(artDetails, false)
	if err != nil {
		return
	}
	for i := 0; i < len(moveSpec.Files); i++ {
		params, err := moveSpec.Get(i).ToArtifatoryMoveCopyParams()
		if err != nil {
			log.Error(err)
			continue
		}
		flat, err := moveSpec.Get(i).IsFlat(false)
		if err != nil {
			log.Error(err)
			continue
		}
		partialSuccess, partialFailed, err := servicesManager.Move(&services.MoveCopyParamsImpl{ArtifactoryCommonParams: params, Flat: flat})
		successCount += partialSuccess
		failCount += partialFailed
		if err != nil {
			log.Error(err)
			continue
		}
	}
	return
}
