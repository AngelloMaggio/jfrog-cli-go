package commands

import (
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/artifactory/utils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/utils/config"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/artifactory/services"
)

func BuildDistribute(flags *BuildDistributionConfiguration) error {
	servicesManager, err := utils.CreateServiceManager(flags.ArtDetails, flags.DryRun)
	if err != nil {
		return err
	}
	return servicesManager.DistributeBuild(flags.BuildDistributionParamsImpl)
}

type BuildDistributionConfiguration struct {
	*services.BuildDistributionParamsImpl
	ArtDetails *config.ArtifactoryDetails
	DryRun     bool
}
