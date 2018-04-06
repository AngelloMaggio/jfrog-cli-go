package commands

import (
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/utils/config"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-cli/utils/ioutils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/utils/errorutils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/utils/log"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func Config(details, defaultDetails *config.BintrayDetails, interactive bool) (*config.BintrayDetails, error) {
	if details == nil {
		details = new(config.BintrayDetails)
	}
	if interactive {
		if defaultDetails == nil {
			var err error
			defaultDetails, err = config.ReadBintrayConf()
			if err != nil {
				return nil, err
			}
		}
		if details.User == "" {
			ioutils.ScanFromConsole("User", &details.User, defaultDetails.User)
		}
		if details.Key == "" {
			print("Key: ")
			byteKey, err := terminal.ReadPassword(int(syscall.Stdin))
			err = errorutils.CheckError(err)
			if err != nil {
				return nil, err
			}
			details.Key = string(byteKey)
			if details.Key == "" {
				details.Key = defaultDetails.Key
			}
		}
		if details.DefPackageLicense == "" {
			ioutils.ScanFromConsole("\nDefault package licenses",
				&details.DefPackageLicense, defaultDetails.DefPackageLicense)
		}
	}
	err := config.SaveBintrayConf(details)
	return details, err
}

func ShowConfig() error {
	details, err := config.ReadBintrayConf()
	if err != nil {
		return err
	}
	if details.User != "" {
		log.Output("User: " + details.User)
	}
	if details.Key != "" {
		log.Output("Key: ***")
	}
	if details.DefPackageLicense != "" {
		log.Output("Default package license: " + details.DefPackageLicense)
	}
	return nil
}

func ClearConfig() {
	config.SaveBintrayConf(new(config.BintrayDetails))
}

func GetConfig() (*config.BintrayDetails, error) {
	return config.ReadBintrayConf()
}
