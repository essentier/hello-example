package builder

import (
	"strings"

	"github.com/essentier/servicebuilder/scm"
	"github.com/essentier/spickspan/config"
	"github.com/essentier/spickspan/model"
)

// A service builder knows how to build services.
type ServiceBuilder interface {
	BuildService(serviceName string, project scm.Project) error
}

func CreateDefaultServiceBuilder(cloudProvider config.CloudProvider) (ServiceBuilder, error) {
	token, err := model.LoginToEssentier(cloudProvider.Url, cloudProvider.Username, cloudProvider.Password)
	if err != nil {
		return nil, err
	}

	return serviceBuilder{nativeBuilder: createDefaultNativeBuilder(cloudProvider.Url), token: token}, nil
}

func CreateServiceBuilder(nativeBuilder nativeBuilder, token string) ServiceBuilder {
	return serviceBuilder{nativeBuilder: nativeBuilder, token: token}
}

// This is a low level service builder that pushes code to the native builder and trigger a build there.
type serviceBuilder struct {
	nativeBuilder nativeBuilder
	token         string
}

func (p serviceBuilder) BuildService(serviceName string, project scm.Project) error {
	repoUrl := p.getServiceDepoUrl(serviceName)
	err := project.PushCode(repoUrl)
	if err != nil {
		return err
	}

	return p.nativeBuilder.buildService(serviceName, p.token)
}

func (p serviceBuilder) getServiceDepoUrl(serviceName string) string {
	builderUrl := p.nativeBuilder.url()
	return getEssentierGitRemote(serviceName, builderUrl, p.token)
}

func getEssentierGitRemote(serviceName string, builderUrl string, token string) string {
	remoteUrl := ""
	if strings.HasPrefix(builderUrl, "git") {
		remoteUrl = builderUrl + ":" + serviceName
	} else if strings.HasPrefix(builderUrl, "https://") {
		remoteUrl = "https://" + token + ":@" + builderUrl[8:] + "/" + serviceName
	}
	return remoteUrl
}
