package artifactory

import (
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/artifactory/auth/cert"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/artifactory/services"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/artifactory/services/utils"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/httpclient"
	"github.com/AngelloMaggio/jfrog-cli-go/jfrog-client/utils/log"
	"net/http"
	"io"
)

type ArtifactoryServicesManager struct {
	client *httpclient.HttpClient
	config Config
}

func New(config Config) (*ArtifactoryServicesManager, error) {
	var err error
	manager := &ArtifactoryServicesManager{config: config}
	if config.GetCertifactesPath() == "" {
		manager.client = httpclient.NewDefaultHttpClient()
	} else {
		transport, err := cert.GetTransportWithLoadedCert(config.GetCertifactesPath())
		if err != nil {
			return nil, err
		}
		manager.client = httpclient.NewHttpClient(&http.Client{Transport: transport})
	}
	if config.GetLogger() != nil {
		log.SetLogger(config.GetLogger())
	}
	return manager, err
}

func (sm *ArtifactoryServicesManager) DistributeBuild(params services.BuildDistributionParams) error {
	distributionService := services.NewDistributionService(sm.client)
	distributionService.DryRun = sm.config.IsDryRun()
	distributionService.ArtDetails = sm.config.GetArtDetails()
	return distributionService.BuildDistribute(params)
}

func (sm *ArtifactoryServicesManager) PromoteBuild(params services.PromotionParams) error {
	promotionService := services.NewPromotionService(sm.client)
	promotionService.DryRun = sm.config.IsDryRun()
	promotionService.ArtDetails = sm.config.GetArtDetails()
	return promotionService.BuildPromote(params)
}

func (sm *ArtifactoryServicesManager) XrayScanBuild(params services.XrayScanParams) ([]byte, error) {
	xrayScanService := services.NewXrayScanService(sm.client)
	xrayScanService.ArtDetails = sm.config.GetArtDetails()
	return xrayScanService.ScanBuild(params)
}

func (sm *ArtifactoryServicesManager) GetPathsToDelete(params services.DeleteParams) ([]utils.ResultItem, error) {
	deleteService := services.NewDeleteService(sm.client)
	deleteService.DryRun = sm.config.IsDryRun()
	deleteService.ArtDetails = sm.config.GetArtDetails()
	return deleteService.GetPathsToDelete(params)
}

func (sm *ArtifactoryServicesManager) DeleteFiles(resultItems []services.DeleteItem) (int, error) {
	deleteService := services.NewDeleteService(sm.client)
	deleteService.DryRun = sm.config.IsDryRun()
	deleteService.ArtDetails = sm.config.GetArtDetails()
	return deleteService.DeleteFiles(resultItems, deleteService)
}

func (sm *ArtifactoryServicesManager) ReadRemoteFile(readPath string) (io.ReadCloser, error) {
	readFileService := services.NewReadFileService(sm.client)
	readFileService.DryRun = sm.config.IsDryRun()
	readFileService.ArtDetails = sm.config.GetArtDetails()
	return readFileService.ReadRemoteFile(readPath)
}

func (sm *ArtifactoryServicesManager) DownloadFiles(params services.DownloadParams) ([]utils.FileInfo, int, error) {
	downloadService := services.NewDownloadService(sm.client)
	downloadService.DryRun = sm.config.IsDryRun()
	downloadService.ArtDetails = sm.config.GetArtDetails()
	downloadService.Threads = sm.config.GetThreads()
	downloadService.SplitCount = sm.config.GetSplitCount()
	downloadService.MinSplitSize = sm.config.GetMinSplitSize()
	return downloadService.DownloadFiles(params)
}

func (sm *ArtifactoryServicesManager) GetUnreferencedGitLfsFiles(params services.GitLfsCleanParams) ([]utils.ResultItem, error) {
	gitLfsCleanService := services.NewGitLfsCleanService(sm.client)
	gitLfsCleanService.DryRun = sm.config.IsDryRun()
	gitLfsCleanService.ArtDetails = sm.config.GetArtDetails()
	return gitLfsCleanService.GetUnreferencedGitLfsFiles(params)
}

func (sm *ArtifactoryServicesManager) Search(params utils.SearchParams) ([]utils.ResultItem, error) {
	searchService := services.NewSearchService(sm.client)
	searchService.ArtDetails = sm.config.GetArtDetails()
	return searchService.Search(params)
}

func (sm *ArtifactoryServicesManager) Aql(aql string) ([]byte, error) {
	aqlService := services.NewAqlService(sm.client)
	aqlService.ArtDetails = sm.config.GetArtDetails()
	return aqlService.ExecAql(aql)
}

func (sm *ArtifactoryServicesManager) SetProps(params services.SetPropsParams) (int, error) {
	setPropsService := services.NewSetPropsService(sm.client)
	setPropsService.ArtDetails = sm.config.GetArtDetails()
	setPropsService.Threads = sm.config.GetThreads()
	return setPropsService.SetProps(params)
}

func (sm *ArtifactoryServicesManager) UploadFiles(params services.UploadParams) (artifactsFileInfo []utils.FileInfo, totalUploaded, totalFailed int, err error) {
	uploadService := services.NewUploadService(sm.client)
	sm.setCommonServiceConfig(uploadService)
	uploadService.MinChecksumDeploy = sm.config.GetMinChecksumDeploy()
	return uploadService.UploadFiles(params)
}

func (sm *ArtifactoryServicesManager) Copy(params services.MoveCopyParams) (successCount, failedCount int, err error) {
	copyService := services.NewMoveCopyService(sm.client, services.COPY)
	copyService.ArtDetails = sm.config.GetArtDetails()
	return copyService.MoveCopyServiceMoveFilesWrapper(params)
}

func (sm *ArtifactoryServicesManager) Move(params services.MoveCopyParams) (successCount, failedCount int, err error) {
	moveService := services.NewMoveCopyService(sm.client, services.MOVE)
	moveService.ArtDetails = sm.config.GetArtDetails()
	return moveService.MoveCopyServiceMoveFilesWrapper(params)
}

func (sm *ArtifactoryServicesManager) setCommonServiceConfig(commonConfig ArtifactoryServicesSetter) {
	commonConfig.SetThread(sm.config.GetThreads())
	commonConfig.SetArtDetails(sm.config.GetArtDetails())
	commonConfig.SetDryRun(sm.config.IsDryRun())
}
