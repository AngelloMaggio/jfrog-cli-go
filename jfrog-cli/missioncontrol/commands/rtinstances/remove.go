package rtinstances

import (
	"errors"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/missioncontrol/utils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/utils/config"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/utils/errorutils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/utils/io/httputils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/utils/log"
)

func Remove(instanceName string, flags *RemoveFlags) error {
	missionControlUrl := flags.MissionControlDetails.Url + "api/v1/instances/" + instanceName
	httpClientDetails := utils.GetMissionControlHttpClientDetails(flags.MissionControlDetails)
	resp, body, err := httputils.SendDelete(missionControlUrl, nil, httpClientDetails)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return errorutils.CheckError(errors.New(resp.Status + ". " + utils.ReadMissionControlHttpMessage(body)))
	}
	log.Debug("Mission Control response: " + resp.Status)
	return nil
}

type RemoveFlags struct {
	MissionControlDetails *config.MissionControlDetails
	Interactive           bool
}
