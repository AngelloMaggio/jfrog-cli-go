package commands

import (
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/bintray"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/bintray/services/url"
)

func SignVersion(config bintray.Config, params *url.Params) (err error) {
	sm, err := bintray.New(config)
	if err != nil {
		return err
	}
	return sm.SignUrl(params)
}
