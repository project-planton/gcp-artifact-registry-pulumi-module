package gcpartifactregistry

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gcpartifactregistry/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	DockerRepoHostnameOutputName        = "docker-repo-hostname"
	DockerRepoNameOutputName            = "docker-repo-name"
	DockerRepoUrlOutputName             = "docker-repo-url"
	MavenRepoNameOutputName             = "maven-repo-name"
	MavenRepoUrlOutputName              = "maven-repo-url"
	NpmRepoNameOutputName               = "npm-repo-name"
	PythonRepoNameOutputName            = "python-repo-name"
	ReaderServiceAccountEmailOutputName = "reader-service-account-email"
	ReaderServiceAccountKeyOutputName   = "reader-service-account-key"
	WriterServiceAccountEmailOutputName = "writer-service-account-email"
	WriterServiceAccountKeyOutputName   = "writer-service-account-key"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *model.GcpArtifactRegistryStackInput) *model.GcpArtifactRegistryStackOutputs {
	return &model.GcpArtifactRegistryStackOutputs{
		ReaderServiceAccountEmail:     autoapistackoutput.GetVal(pulumiOutputs, ReaderServiceAccountEmailOutputName),
		ReaderServiceAccountKeyBase64: autoapistackoutput.GetVal(pulumiOutputs, ReaderServiceAccountKeyOutputName),
		WriterServiceAccountEmail:     autoapistackoutput.GetVal(pulumiOutputs, WriterServiceAccountEmailOutputName),
		WriterServiceAccountKeyBase64: autoapistackoutput.GetVal(pulumiOutputs, WriterServiceAccountKeyOutputName),
		DockerRepoName:                autoapistackoutput.GetVal(pulumiOutputs, DockerRepoNameOutputName),
		DockerRepoHostname:            autoapistackoutput.GetVal(pulumiOutputs, DockerRepoHostnameOutputName),
		DockerRepoUrl:                 autoapistackoutput.GetVal(pulumiOutputs, DockerRepoUrlOutputName),
		MavenRepoName:                 autoapistackoutput.GetVal(pulumiOutputs, MavenRepoNameOutputName),
		MavenRepoUrl:                  autoapistackoutput.GetVal(pulumiOutputs, MavenRepoUrlOutputName),
		NpmRepoName:                   autoapistackoutput.GetVal(pulumiOutputs, NpmRepoNameOutputName),
		PythonRepoName:                autoapistackoutput.GetVal(pulumiOutputs, PythonRepoNameOutputName),
	}
}
