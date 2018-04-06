package commands

import (
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/artifactory/utils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/utils/config"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/artifactory/services"
)

func BuildPromote(flags *BuildPromotionConfiguration) error {
	servicesManager, err := utils.CreateServiceManager(flags.ArtDetails, flags.DryRun)
	if err != nil {
		return err
	}
	return servicesManager.PromoteBuild(flags.PromotionParamsImpl)
}

type BuildPromotionConfiguration struct {
	*services.PromotionParamsImpl
	ArtDetails *config.ArtifactoryDetails
	DryRun     bool
}
