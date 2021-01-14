// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package reporter

import (
	"context"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/managers"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/prometheus"
	"github.com/redhat-marketplace/redhat-marketplace-operator/v2/pkg/utils/reconcileutils"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Injectors from wire.go:

func NewTask(ctx context.Context, reportName ReportName, config2 *Config) (*Task, error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	restMapper, err := managers.NewDynamicRESTMapper(restConfig)
	if err != nil {
		return nil, err
	}
	scheme := provideScheme()
	clientOptions := getClientOptions()
	client, err := managers.ProvideSimpleClient(restConfig, restMapper, scheme, clientOptions)
	if err != nil {
		return nil, err
	}
	logrLogger := _wireLoggerValue
	clientCommandRunner := reconcileutils.NewClientCommand(client, scheme, logrLogger)
	uploaderTarget := config2.UploaderTarget
	uploader, err := ProvideUploader(ctx, clientCommandRunner, logrLogger, uploaderTarget)
	if err != nil {
		return nil, err
	}
	task := &Task{
		ReportName: reportName,
		CC:         clientCommandRunner,
		K8SClient:  client,
		Ctx:        ctx,
		Config:     config2,
		K8SScheme:  scheme,
		Uploader:   uploader,
	}
	return task, nil
}

var (
	_wireLoggerValue = logger
)

func NewReporter(task *Task) (*MarketplaceReporter, error) {
	reporterConfig := task.Config
	client := task.K8SClient
	contextContext := task.Ctx
	scheme := task.K8SScheme
	logrLogger := _wireLogrLoggerValue
	clientCommandRunner := reconcileutils.NewClientCommand(client, scheme, logrLogger)
	reportName := task.ReportName
	meterReport, err := getMarketplaceReport(contextContext, clientCommandRunner, reportName)
	if err != nil {
		return nil, err
	}
	marketplaceConfig, err := getMarketplaceConfig(contextContext, clientCommandRunner)
	if err != nil {
		return nil, err
	}
	service, err := getPrometheusService(contextContext, meterReport, clientCommandRunner)
	if err != nil {
		return nil, err
	}
	prometheusAPISetup := providePrometheusSetup(reporterConfig, meterReport, service)
	prometheusAPI, err := prometheus.NewPrometheusAPIForReporter(prometheusAPISetup)
	if err != nil {
		return nil, err
	}
	marketplaceReporter, err := NewMarketplaceReporter(reporterConfig, client, meterReport, marketplaceConfig, service, prometheusAPI)
	if err != nil {
		return nil, err
	}
	return marketplaceReporter, nil
}

var (
	_wireLogrLoggerValue = logger
)
