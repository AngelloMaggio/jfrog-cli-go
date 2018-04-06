package npm

import (
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/artifactory/utils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/utils/errorutils"
	"github.com/mattn/go-shellwords"
)

func Pack(npmFlags, executablePath string) error {
	splitFlags, err := shellwords.Parse(npmFlags)
	if err != nil {
		return errorutils.CheckError(err)
	}

	configListCmdConfig := createPackCmdConfig(executablePath, splitFlags)
	if err := utils.RunCmd(configListCmdConfig); err != nil {
		return errorutils.CheckError(err)
	}

	return nil
}

func createPackCmdConfig(executablePath string, splitFlags []string) *NpmConfig {
	return &NpmConfig{
		Npm:          executablePath,
		Command:      []string{"pack"},
		CommandFlags: append(splitFlags),
		StrWriter:    nil,
		ErrWriter:    nil,
	}
}
