package npm

import (
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/artifactory/utils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/utils/config"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/utils/errorutils"
	"github.com/mattn/go-shellwords"
	"io"
	"io/ioutil"
)

// This method runs "npm config list --json" command and returns the json object that contains the current configurations of npm
// Fore more info see https://docs.npmjs.com/cli/config
func GetConfigList(npmFlags, executablePath string) ([]byte, error) {
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()
	splitFlags, err := shellwords.Parse(npmFlags)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	configListCmdConfig := createConfigListCmdConfig(executablePath, splitFlags, pipeWriter)
	var npmError error
	go func() {
		npmError = utils.RunCmd(configListCmdConfig)
	}()

	data, err := ioutil.ReadAll(pipeReader)
	if err != nil {
		return nil, errorutils.CheckError(err)
	}

	if npmError != nil {
		return nil, errorutils.CheckError(npmError)
	}
	return data, nil
}

func createConfigListCmdConfig(executablePath string, splitFlags []string, pipeWriter *io.PipeWriter) *NpmConfig {
	return &NpmConfig{
		Npm:          executablePath,
		Command:      []string{"c", "ls"},
		CommandFlags: append(splitFlags, "-json=true"),
		StrWriter:    pipeWriter,
		ErrWriter:    nil,
	}
}

type CliConfiguration struct {
	BuildName   string
	BuildNumber string
	NpmArgs     string
	ArtDetails  *config.ArtifactoryDetails
}
