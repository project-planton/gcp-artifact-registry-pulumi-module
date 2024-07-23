package gcpartifactregistry

import (
	"fmt"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gcpartifactregistry/enums/gcpartifactregistryrepotype"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/gcp/pulumigoogleprovider"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (s *ResourceStack) npmRepo(ctx *pulumi.Context, gcpProvider *pulumigcp.Provider,
	readerServiceAccount *serviceaccount.Account, writerServiceAccount *serviceaccount.Account) error {

	gcpArtifactRegistry := s.Input.ApiResource

	repoName := GetNpmRepoName(gcpArtifactRegistry.Metadata.Id)

	addedNpmRepo, err := artifactregistry.NewRepository(ctx, repoName,
		&artifactregistry.RepositoryArgs{
			Project:      pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
			Location:     pulumi.String(gcpArtifactRegistry.Spec.Region),
			RepositoryId: pulumi.String(repoName),
			Format:       pulumi.String(gcpartifactregistryrepotype.GcpArtifactRegistryRepoType_NPM.String()),
			Labels:       pulumi.ToStringMap(s.GcpLabels),
		}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create %s repo", repoName)
	}

	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-reader",
		repoName), &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
		Location:   pulumi.String(gcpArtifactRegistry.Spec.Region),
		Repository: addedNpmRepo.RepositoryId,
		Role:       pulumi.String(standard.Artifactregistry_reader),
		Member:     pulumi.Sprintf("serviceAccounts:%s", readerServiceAccount.Email),
	}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role to svc acct on %s repo",
			standard.Artifactregistry_reader, repoName)
	}

	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-writer",
		repoName), &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
		Location:   pulumi.String(gcpArtifactRegistry.Spec.Region),
		Repository: addedNpmRepo.RepositoryId,
		Role:       pulumi.String(standard.Artifactregistry_writer),
		Member:     pulumi.Sprintf("serviceAccounts:%s", writerServiceAccount.Email),
	}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role svc acct on %s repo",
			standard.Artifactregistry_writer, repoName)
	}

	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-admin",
		repoName), &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
		Location:   pulumi.String(gcpArtifactRegistry.Spec.Region),
		Repository: addedNpmRepo.RepositoryId,
		Role:       pulumi.String(standard.Artifactregistry_repoAdmin),
		Member:     pulumi.Sprintf("serviceAccounts:%s", writerServiceAccount.Email),
	}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role svc acct on %s repo",
			standard.Artifactregistry_repoAdmin, repoName)
	}

	ctx.Export(GetNpmRepoNameOutputName(repoName), addedNpmRepo.RepositoryId)

	return nil
}

func GetNpmRepoNameOutputName(repoName string) string {
	return pulumigoogleprovider.PulumiOutputName(artifactregistry.Repository{}, repoName,
		gcpartifactregistryrepotype.GcpArtifactRegistryRepoType_NPM.String(), englishword.EnglishWord_name.String())
}

func GetNpmRepoName(gcpArtifactRegistryId string) string {
	return fmt.Sprintf("%s-npm", gcpArtifactRegistryId)
}
