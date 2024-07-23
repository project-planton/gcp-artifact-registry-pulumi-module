package gcpartifactregistry

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gcpartifactregistry/enums/gcpartifactregistryrepotype"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/gcp/pulumigoogleprovider"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (s *ResourceStack) pythonRepo(ctx *pulumi.Context, gcpProvider *pulumigcp.Provider,
	readerServiceAccount *serviceaccount.Account, writerServiceAccount *serviceaccount.Account) error {

	gcpArtifactRegistry := s.Input.ApiResource

	repoName := GetPythonRepoName(gcpArtifactRegistry.Metadata.Id)

	addedPythonRepo, err := artifactregistry.NewRepository(ctx, fmt.Sprintf("%s-python", repoName),
		&artifactregistry.RepositoryArgs{
			Project:      pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
			Location:     pulumi.String(gcpArtifactRegistry.Spec.Region),
			RepositoryId: pulumi.String(repoName),
			Format:       pulumi.String(gcpartifactregistryrepotype.GcpArtifactRegistryRepoType_PYTHON.String()),
			Labels:       pulumi.ToStringMap(s.GcpLabels),
		}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create %s repo", repoName)
	}

	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-reader", repoName),
		&artifactregistry.RepositoryIamMemberArgs{
			Project:    pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
			Location:   pulumi.String(gcpArtifactRegistry.Spec.Region),
			Repository: addedPythonRepo.RepositoryId,
			Role:       pulumi.String(standard.Artifactregistry_reader),
			Member: pulumi.Sprintf("serviceAccounts:%s",
				readerServiceAccount.Email),
		}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role to svc acct on %s repo",
			standard.Artifactregistry_reader, repoName)
	}

	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-writer", repoName),
		&artifactregistry.RepositoryIamMemberArgs{
			Project:    pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
			Location:   pulumi.String(gcpArtifactRegistry.Spec.Region),
			Repository: addedPythonRepo.RepositoryId,
			Role:       pulumi.String(standard.Artifactregistry_writer),
			Member: pulumi.Sprintf("serviceAccounts:%s",
				writerServiceAccount.Email),
		}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role svc acct on %s repo",
			standard.Artifactregistry_writer, repoName)
	}

	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-admin", repoName),
		&artifactregistry.RepositoryIamMemberArgs{
			Project:    pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
			Location:   pulumi.String(gcpArtifactRegistry.Spec.Region),
			Repository: addedPythonRepo.RepositoryId,
			Role:       pulumi.String(standard.Artifactregistry_repoAdmin),
			Member: pulumi.Sprintf("serviceAccounts:%s",
				writerServiceAccount.Email),
		}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role svc acct on %s repo",
			standard.Artifactregistry_repoAdmin, repoName)
	}

	ctx.Export(GetPythonRepoNameOutputName(repoName), addedPythonRepo.RepositoryId)

	return nil
}

func GetPythonRepoNameOutputName(gcpArtifactRegistryId string) string {
	return pulumigoogleprovider.PulumiOutputName(artifactregistry.Repository{}, gcpArtifactRegistryId,
		gcpartifactregistryrepotype.GcpArtifactRegistryRepoType_PYTHON.String(), englishword.EnglishWord_name.String())
}

func GetPythonRepoName(gcpArtifactRegistryId string) string {
	return fmt.Sprintf("%s-python", gcpArtifactRegistryId)
}
