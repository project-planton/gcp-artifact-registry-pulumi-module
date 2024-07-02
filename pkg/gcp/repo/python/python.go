package python

import (
	"fmt"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/pulumigoogleprovider"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/artifactstore/enums/gcpartifactregistryrepotype"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	"github.com/plantoncloud/artifact-store-pulumi-blueprint/pkg/gcp/serviceaccount"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/artifactregistry"
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
	pythonRepo, err := artifactregistry.NewRepository(ctx, fmt.Sprintf("%s-python", repoName),
		&artifactregistry.RepositoryArgs{
			Project:      pulumi.String(input.GcpProjectId),
			Location:     pulumi.String(input.GcpRegion),
			RepositoryId: pulumi.String(repoName),
			Format:       pulumi.String(gcpartifactregistryrepotype.GcpArtifactRegistryRepoType_PYTHON.String()),
			Labels:       pulumi.ToStringMap(input.Labels),
		}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create %s repo", repoName)
	}
	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-reader", repoName),
		&artifactregistry.RepositoryIamMemberArgs{
			Project:    pulumi.String(input.GcpProjectId),
			Location:   pulumi.String(input.GcpRegion),
			Repository: pythonRepo.RepositoryId,
			Role:       pulumi.String(standard.Artifactregistry_reader),
			Member: pulumi.Sprintf("serviceAccount:%s",
				input.ServiceAccountOutput.ReaderServiceAccount.Email),
		}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role to svc acct on %s repo",
			standard.Artifactregistry_reader, repoName)
	}
	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-writer", repoName),
		&artifactregistry.RepositoryIamMemberArgs{
			Project:    pulumi.String(input.GcpProjectId),
			Location:   pulumi.String(input.GcpRegion),
			Repository: pythonRepo.RepositoryId,
			Role:       pulumi.String(standard.Artifactregistry_writer),
			Member: pulumi.Sprintf("serviceAccount:%s",
				input.ServiceAccountOutput.WriterServiceAccount.Email),
		}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role svc acct on %s repo",
			standard.Artifactregistry_writer, repoName)
	}
	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-admin", repoName),
		&artifactregistry.RepositoryIamMemberArgs{
			Project:    pulumi.String(input.GcpProjectId),
			Location:   pulumi.String(input.GcpRegion),
			Repository: pythonRepo.RepositoryId,
			Role:       pulumi.String(standard.Artifactregistry_repoAdmin),
			Member: pulumi.Sprintf("serviceAccount:%s",
				input.ServiceAccountOutput.WriterServiceAccount.Email),
		}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role svc acct on %s repo",
			standard.Artifactregistry_repoAdmin, repoName)
	}

	exportOutputs(ctx, repoName, pythonRepo)
	return nil
}

func GetPythonRepoNameOutputName(artifactStoreId string) string {
	return pulumigoogleprovider.PulumiOutputName(artifactregistry.Repository{}, artifactStoreId,
		gcpartifactregistryrepotype.GcpArtifactRegistryRepoType_PYTHON.String(), englishword.EnglishWord_name.String())
}

func exportOutputs(ctx *pulumi.Context, repoName string, pythonRepo *artifactregistry.Repository) {
	ctx.Export(GetPythonRepoNameOutputName(repoName), pythonRepo.RepositoryId)
}

func GetRepoName(artifactStoreId string) string {
	return fmt.Sprintf("%s-python", artifactStoreId)
}
