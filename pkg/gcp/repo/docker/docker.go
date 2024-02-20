package docker

import (
	"fmt"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/artifactstore/enums/gcpartifactregistryrepotype"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/serviceaccount"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	puluminameoutputgcp "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/gcp/output"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	HostnameSuffix = "docker.pkg.dev"
)

type Input struct {
	GcpProvider          *pulumigcp.Provider
	ArtifactStoreId      string
	GcpProjectId         string
	GcpRegion            string
	ServiceAccountOutput *serviceaccount.AddedResources
	IsExternal           bool
	Labels               map[string]string
}

func Resources(ctx *pulumi.Context, input *Input) error {
	addedRepo, err := addRepo(ctx, input)
	if err != nil {
		return errors.Wrapf(err, "failed to add internal repo resources")
	}

	exportOutputs(ctx, GetRepoName(input.ArtifactStoreId), addedRepo)

	return nil
}

func addRepo(ctx *pulumi.Context, input *Input) (*artifactregistry.Repository, error) {
	repoName := GetRepoName(input.ArtifactStoreId)
	r, err := artifactregistry.NewRepository(ctx, repoName,
		&artifactregistry.RepositoryArgs{
			Project:      pulumi.String(input.GcpProjectId),
			Location:     pulumi.String(input.GcpRegion),
			RepositoryId: pulumi.String(repoName),
			Format:       pulumi.String(gcpartifactregistryrepotype.GcpArtifactRegistryRepoType_DOCKER.String()),
			Labels:       pulumi.ToStringMap(input.Labels),
		}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create %s repo", repoName)
	}
	if input.IsExternal {
		_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-reader-%s",
			repoName, iam.AllUsersIdentifier), &artifactregistry.RepositoryIamMemberArgs{
			Project:    pulumi.String(input.GcpProjectId),
			Location:   pulumi.String(input.GcpRegion),
			Repository: r.RepositoryId,
			Role:       pulumi.String(standard.Artifactregistry_reader),
			Member:     pulumi.Sprintf(iam.AllUsersIdentifier),
		}, pulumi.Provider(input.GcpProvider))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to add %s role to %s",
				standard.Artifactregistry_reader, iam.AllUsersIdentifier)
		}
	}
	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-reader",
		repoName), &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(input.GcpProjectId),
		Location:   pulumi.String(input.GcpRegion),
		Repository: r.RepositoryId,
		Role:       pulumi.String(standard.Artifactregistry_reader),
		Member:     pulumi.Sprintf("serviceAccount:%s", input.ServiceAccountOutput.ReaderServiceAccount.Email),
	}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s role to svc acct on %s repo",
			standard.Artifactregistry_reader, repoName)
	}
	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-writer",
		repoName), &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(input.GcpProjectId),
		Location:   pulumi.String(input.GcpRegion),
		Repository: r.RepositoryId,
		Role:       pulumi.String(standard.Artifactregistry_writer),
		Member:     pulumi.Sprintf("serviceAccount:%s", input.ServiceAccountOutput.WriterServiceAccount.Email),
	}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s role svc acct on %s repo",
			standard.Artifactregistry_writer, repoName)
	}
	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-admin",
		repoName), &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(input.GcpProjectId),
		Location:   pulumi.String(input.GcpRegion),
		Repository: r.RepositoryId,
		Role:       pulumi.String(standard.Artifactregistry_repoAdmin),
		Member:     pulumi.Sprintf("serviceAccount:%s", input.ServiceAccountOutput.WriterServiceAccount.Email),
	}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s role svc acct on %s repo",
			standard.Artifactregistry_repoAdmin, repoName)
	}
	return r, nil
}

func GetDockerRepoHostnameOutputName(repoName string) string {
	return puluminameoutputgcp.Name(artifactregistry.Repository{}, repoName,
		gcpartifactregistryrepotype.GcpArtifactRegistryRepoType_DOCKER.String(), englishword.EnglishWord_hostname.String())
}

func GetDockerRepoNameOutputName(repoName string) string {
	return puluminameoutputgcp.Name(artifactregistry.Repository{}, repoName,
		gcpartifactregistryrepotype.GcpArtifactRegistryRepoType_DOCKER.String(), englishword.EnglishWord_name.String())
}

func GetDockerRepoUrlOutputName(repoName string) string {
	return puluminameoutputgcp.Name(artifactregistry.Repository{}, repoName,
		gcpartifactregistryrepotype.GcpArtifactRegistryRepoType_DOCKER.String(), englishword.EnglishWord_url.String())
}

func exportOutputs(ctx *pulumi.Context, repoName string, addedRepo *artifactregistry.Repository) {
	ctx.Export(GetDockerRepoHostnameOutputName(repoName),
		pulumi.Sprintf("%s-%s", addedRepo.Location, HostnameSuffix))
	ctx.Export(GetDockerRepoNameOutputName(repoName), addedRepo.RepositoryId)
	ctx.Export(GetDockerRepoUrlOutputName(repoName), getDockerRepoUrl(addedRepo))
}

// getDockerRepoUrl constructs complete maven repo url using the provided input
// ex: artifactregistry://us-central1-maven.pkg.dev/planton-shared-services-jx/planton-pcs-maven-repo"
func getDockerRepoUrl(repo *artifactregistry.Repository) pulumi.Input {
	return pulumi.Sprintf("%s-%s/%s/%s", repo.Location,
		HostnameSuffix, repo.Project, repo.Name)
}

func GetRepoName(artifactStoreId string) string {
	return fmt.Sprintf("%s-docker", artifactStoreId)
}
