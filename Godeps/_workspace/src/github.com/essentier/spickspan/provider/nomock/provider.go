package nomock

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bndr/gopencils"
	"github.com/essentier/nomockutil"
	"github.com/essentier/spickspan/config"
	"github.com/essentier/spickspan/model"
	"github.com/essentier/spickspan/provider"
)

const (
	noReleaseServiceID   string = "noReleaseServiceID"
	containerImagePrefix string = "gcr.io/essentier-nomock/" // IP:5000/nomock/
)

func CreateProvider() model.Provider {
	return &TestingProvider{}
}

type TestingProvider struct {
	nomockApi *gopencils.Resource
	token     string
}

func (p *TestingProvider) Init() error {
	config, err := provider.GetConfig()
	if err != nil {
		return err
	}

	cloudProvider := config.CloudProvider
	token, err := model.LoginToEssentier(cloudProvider.Url, cloudProvider.Username, cloudProvider.Password)
	if err != nil {
		return err
	}

	p.token = token
	p.nomockApi = gopencils.Api(cloudProvider.Url) //  + "/nomockserver"
	return nil
}

func (p *TestingProvider) Detect() bool {
	mode := os.Getenv("SPICKSPAN_MODE")
	return strings.ToLower(mode) == "testing"
}

func (p *TestingProvider) Release(service model.Service) error {
	log.Printf("Releasing service %v", service)
	if service.IP == noReleaseServiceID {
		return nil
	}

	res := p.nomockApi.Res("nomockserver/services")
	res = res.Id(service.Id)
	res.SetHeader("Authorization", "Bearer "+p.token)
	_, err := res.Delete()
	return err
}

func (p *TestingProvider) GetService(serviceName string) (model.Service, error) {
	//When this provider is asked for a service,
	//it will find the service's configuration in the config file
	//and use that configuration to start up the service in the testing cloud.
	service, serviceConfig, err := provider.GetServiceAndConfig(serviceName)
	if err != nil || service.Id != "" {
		return service, err
	}

	newService, err := p.createService(serviceConfig)
	if err != nil {
		return newService, err
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM)
	go func() {
		<-sigchan
		//We can do this only when spickspan is in testing mode.
		p.Release(newService)
	}()

	return newService, nil
}

func (p *TestingProvider) createService(serviceConfig config.Service) (model.Service, error) {
	var newService model.Service
	userId, err := nomockutil.GetSubjectInToken(p.token)
	if err != nil {
		return newService, err
	}

	servicesResource := p.nomockApi.Res("nomockserver/services", &newService)
	if serviceConfig.IsSourceProject() {
		serviceConfig.ContainerImage = containerImagePrefix + userId + "_" + serviceConfig.ServiceName + ":latest"
	}
	log.Printf("service config %#v", serviceConfig)

	servicesResource.SetHeader("Authorization", "Bearer "+p.token)
	_, err = servicesResource.Post(serviceConfig)
	if err != nil {
		log.Printf("Failed to call the service rest api. Error is: %#v. Error string is %v", err, err.Error())
	}
	log.Printf("service is: %#v", newService)
	return newService, nil
}
