package services

import (
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/artifactory/auth"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/artifactory/services/utils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/httpclient"
)

type SearchService struct {
	client     *httpclient.HttpClient
	ArtDetails auth.ArtifactoryDetails
}

func NewSearchService(client *httpclient.HttpClient) *SearchService {
	return &SearchService{client: client}
}

func (s *SearchService) GetArtifactoryDetails() auth.ArtifactoryDetails {
	return s.ArtDetails
}

func (s *SearchService) SetArtifactoryDetails(rt auth.ArtifactoryDetails) {
	s.ArtDetails = rt
}

func (s *SearchService) IsDryRun() bool {
	return false
}

func (s *SearchService) GetJfrogHttpClient() *httpclient.HttpClient {
	return s.client
}

func (s *SearchService) Search(searchParamsImpl utils.SearchParams) ([]utils.ResultItem, error) {
	return utils.SearchBySpecFiles(searchParamsImpl, s)
}
