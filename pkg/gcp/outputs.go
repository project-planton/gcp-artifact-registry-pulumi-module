package gcp

import (
	"context"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/stack/job/enums/operationtype"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/repo/docker"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/repo/maven"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/repo/npm"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/repo/python"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/serviceaccount"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	artifactstorestate "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/develop/artifactstore"
	artifactstoregcp "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/develop/artifactstore/stack/gcp"
)

func Outputs(ctx context.Context, input *artifactstoregcp.ArtifactStoreGcpStackInput) (*artifactstoregcp.ArtifactStoreGcpStackOutputs, error) {
	pulumiOrgName, err := org.GetOrgName()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pulumi org name")
	}
	stackOutput, err := backend.StackOutput(pulumiOrgName, input.StackJob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stack output")
	}
	return Get(stackOutput, input), nil
}

func Get(stackOutput map[string]interface{}, input *artifactstoregcp.ArtifactStoreGcpStackInput) *artifactstoregcp.ArtifactStoreGcpStackOutputs {
	if input.StackJob.Spec.OperationType != operationtype.StackJobOperationType_apply || stackOutput == nil {
		return &artifactstoregcp.ArtifactStoreGcpStackOutputs{}
	}
	artifactStoreId := input.ResourceInput.ArtifactStore.Metadata.Id
	dockerRepoName := docker.GetRepoName(artifactStoreId)
	mavenRepoName := maven.GetRepoName(artifactStoreId)
	npmRepoName := npm.GetRepoName(artifactStoreId)
	pythonRepoName := python.GetRepoName(artifactStoreId)
	return &artifactstoregcp.ArtifactStoreGcpStackOutputs{
		GcpArtifactRegistryStatus: &artifactstorestate.ArtifactStoreGcpArtifactRegistryStatus{
			ReaderServiceAccountEmail:     backend.GetVal(stackOutput, serviceaccount.GetReaderServiceAccountEmailOutputName(artifactStoreId)),
			ReaderServiceAccountKeyBase64: backend.GetVal(stackOutput, serviceaccount.GetReaderServiceAccountKeyOutputName(artifactStoreId)),
			WriterServiceAccountEmail:     backend.GetVal(stackOutput, serviceaccount.GetWriterServiceAccountEmailOutputName(artifactStoreId)),
			WriterServiceAccountKeyBase64: backend.GetVal(stackOutput, serviceaccount.GetWriterServiceAccountKeyOutputName(artifactStoreId)),
			DockerRepoName:                backend.GetVal(stackOutput, docker.GetDockerRepoNameOutputName(dockerRepoName)),
			DockerRepoHostname:            backend.GetVal(stackOutput, docker.GetDockerRepoHostnameOutputName(dockerRepoName)),
			DockerRepoUrl:                 backend.GetVal(stackOutput, docker.GetDockerRepoUrlOutputName(dockerRepoName)),
			MavenRepoName:                 backend.GetVal(stackOutput, maven.GetMavenRepoNameOutputName(mavenRepoName)),
			MavenRepoUrl:                  backend.GetVal(stackOutput, maven.GetMavenRepoUrlOutputName(mavenRepoName)),
			NpmRepoName:                   backend.GetVal(stackOutput, npm.GetNpmRepoNameOutputName(npmRepoName)),
			PythonRepoName:                backend.GetVal(stackOutput, python.GetPythonRepoNameOutputName(pythonRepoName)),
		},
	}
}
