package gcpartifactregistry

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gcpartifactregistry/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

func OutputMapTransformer(stackOutput auto.OutputMap,
	input *model.GcpArtifactRegistryStackInput) *model.GcpArtifactRegistryStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply ||
		stackOutput == nil {
		return &model.GcpArtifactRegistryStackOutputs{}
	}

	gcpArtifactRegistryId := input.ApiResource.Metadata.Id
	dockerRepoName := GetDockerRepoName(gcpArtifactRegistryId)
	mavenRepoName := GetMavenRepoName(gcpArtifactRegistryId)
	npmRepoName := GetNpmRepoName(gcpArtifactRegistryId)
	pythonRepoName := GetPythonRepoName(gcpArtifactRegistryId)

	return &model.GcpArtifactRegistryStackOutputs{
		ReaderServiceAccountEmail: autoapistackoutput.GetVal(stackOutput,
			GetReaderServiceAccountEmailOutputName(gcpArtifactRegistryId)),
		ReaderServiceAccountKeyBase64: autoapistackoutput.GetVal(stackOutput,
			GetReaderServiceAccountKeyOutputName(gcpArtifactRegistryId)),
		WriterServiceAccountEmail: autoapistackoutput.GetVal(stackOutput,
			GetWriterServiceAccountEmailOutputName(gcpArtifactRegistryId)),
		WriterServiceAccountKeyBase64: autoapistackoutput.GetVal(stackOutput,
			GetWriterServiceAccountKeyOutputName(gcpArtifactRegistryId)),
		DockerRepoName: autoapistackoutput.GetVal(stackOutput,
			GetDockerRepoNameOutputName(dockerRepoName)),
		DockerRepoHostname: autoapistackoutput.GetVal(stackOutput,
			GetDockerRepoHostnameOutputName(dockerRepoName)),
		DockerRepoUrl: autoapistackoutput.GetVal(stackOutput,
			GetDockerRepoUrlOutputName(dockerRepoName)),
		MavenRepoName: autoapistackoutput.GetVal(stackOutput,
			GetMavenRepoNameOutputName(mavenRepoName)),
		MavenRepoUrl: autoapistackoutput.GetVal(stackOutput,
			GetMavenRepoUrlOutputName(mavenRepoName)),
		NpmRepoName: autoapistackoutput.GetVal(stackOutput,
			GetNpmRepoNameOutputName(npmRepoName)),
		PythonRepoName: autoapistackoutput.GetVal(stackOutput,
			GetPythonRepoNameOutputName(pythonRepoName)),
	}
}
