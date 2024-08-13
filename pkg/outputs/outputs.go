package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gcpartifactregistry"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	DockerRepoHostname        = "docker-repo-hostname"
	DockerRepoName            = "docker-repo-name"
	DockerRepoUrl             = "docker-repo-url"
	MavenRepoName             = "maven-repo-name"
	MavenRepoUrl              = "maven-repo-url"
	NpmRepoName               = "npm-repo-name"
	PythonRepoName            = "python-repo-name"
	ReaderServiceAccountEmail = "reader-service-account-email"
	ReaderServiceAccountKey   = "reader-service-account-key"
	WriterServiceAccountEmail = "writer-service-account-email"
	WriterServiceAccountKey   = "writer-service-account-key"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *gcpartifactregistry.GcpArtifactRegistryStackInput) *gcpartifactregistry.GcpArtifactRegistryStackOutputs {
	return &gcpartifactregistry.GcpArtifactRegistryStackOutputs{
		ReaderServiceAccountEmail:     autoapistackoutput.GetVal(pulumiOutputs, ReaderServiceAccountEmail),
		ReaderServiceAccountKeyBase64: autoapistackoutput.GetVal(pulumiOutputs, ReaderServiceAccountKey),
		WriterServiceAccountEmail:     autoapistackoutput.GetVal(pulumiOutputs, WriterServiceAccountEmail),
		WriterServiceAccountKeyBase64: autoapistackoutput.GetVal(pulumiOutputs, WriterServiceAccountKey),
		DockerRepoName:                autoapistackoutput.GetVal(pulumiOutputs, DockerRepoName),
		DockerRepoHostname:            autoapistackoutput.GetVal(pulumiOutputs, DockerRepoHostname),
		DockerRepoUrl:                 autoapistackoutput.GetVal(pulumiOutputs, DockerRepoUrl),
		MavenRepoName:                 autoapistackoutput.GetVal(pulumiOutputs, MavenRepoName),
		MavenRepoUrl:                  autoapistackoutput.GetVal(pulumiOutputs, MavenRepoUrl),
		NpmRepoName:                   autoapistackoutput.GetVal(pulumiOutputs, NpmRepoName),
		PythonRepoName:                autoapistackoutput.GetVal(pulumiOutputs, PythonRepoName),
	}
}
