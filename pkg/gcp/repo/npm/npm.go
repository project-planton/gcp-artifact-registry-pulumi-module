package npm

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/serviceaccount"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	puluminameoutputgcp "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/gcp/output"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/develop/artifactstore/provider/gcpartifactregistry/enums/repotype"
	wordpb "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/commons/english/enums"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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
	repoName := GetRepoName(input.ArtifactStoreId)
	npmRepo, err := artifactregistry.NewRepository(ctx, repoName,
		&artifactregistry.RepositoryArgs{
			Project:      pulumi.String(input.GcpProjectId),
			Location:     pulumi.String(input.GcpRegion),
			RepositoryId: pulumi.String(repoName),
			Format:       pulumi.String(repotype.GcpArtifactRegistryRepoType_NPM.String()),
			Labels:       pulumi.ToStringMap(input.Labels),
		}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create %s repo", repoName)
	}
	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-reader",
		repoName), &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(input.GcpProjectId),
		Location:   pulumi.String(input.GcpRegion),
		Repository: npmRepo.RepositoryId,
		Role:       pulumi.String(standard.Artifactregistry_reader),
		Member:     pulumi.Sprintf("serviceAccount:%s", input.ServiceAccountOutput.ReaderServiceAccount.Email),
	}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role to svc acct on %s repo",
			standard.Artifactregistry_reader, repoName)
	}
	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-writer",
		repoName), &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(input.GcpProjectId),
		Location:   pulumi.String(input.GcpRegion),
		Repository: npmRepo.RepositoryId,
		Role:       pulumi.String(standard.Artifactregistry_writer),
		Member:     pulumi.Sprintf("serviceAccount:%s", input.ServiceAccountOutput.WriterServiceAccount.Email),
	}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role svc acct on %s repo",
			standard.Artifactregistry_writer, repoName)
	}
	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-admin",
		repoName), &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(input.GcpProjectId),
		Location:   pulumi.String(input.GcpRegion),
		Repository: npmRepo.RepositoryId,
		Role:       pulumi.String(standard.Artifactregistry_repoAdmin),
		Member:     pulumi.Sprintf("serviceAccount:%s", input.ServiceAccountOutput.WriterServiceAccount.Email),
	}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role svc acct on %s repo",
			standard.Artifactregistry_repoAdmin, repoName)
	}

	exportOutputs(ctx, repoName, npmRepo)

	return nil
}

func GetNpmRepoNameOutputName(repoName string) string {
	return puluminameoutputgcp.Name(artifactregistry.Repository{}, repoName,
		repotype.GcpArtifactRegistryRepoType_NPM.String(), wordpb.Word_name.String())
}

func exportOutputs(ctx *pulumi.Context, repoName string, npmRepo *artifactregistry.Repository) {
	ctx.Export(GetNpmRepoNameOutputName(repoName), npmRepo.RepositoryId)
}

func GetRepoName(artifactStoreId string) string {
	return fmt.Sprintf("%s-npm", artifactStoreId)
}
