package pkg

import (
	"fmt"
	"github.com/plantoncloud/gcp-artifact-registry-pulumi-module/pkg/outputs"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gcpartifactregistry/enums/gcpartifactregistryrepotype"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// npmRepo creates npm repository and also grants reader role to the reader service account and writer, admin roles to
// writer service account.
func (s *ResourceStack) npmRepo(ctx *pulumi.Context, gcpProvider *pulumigcp.Provider,
	readerServiceAccount *serviceaccount.Account, writerServiceAccount *serviceaccount.Account) error {
	//create a variable with descriptive name for the api-resource in the input
	gcpArtifactRegistry := s.Input.ApiResource

	//create a name for the npm repo since the name of this repository should be unique with in the gcp project.
	npmRepoName := fmt.Sprintf("%s-npm", gcpArtifactRegistry.Metadata.Id)

	//create npm repository
	createdNpmRepo, err := artifactregistry.NewRepository(ctx,
		npmRepoName,
		&artifactregistry.RepositoryArgs{
			Project:      pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
			Location:     pulumi.String(gcpArtifactRegistry.Spec.Region),
			RepositoryId: pulumi.String(npmRepoName),
			Format:       pulumi.String(gcpartifactregistryrepotype.GcpArtifactRegistryRepoType_NPM.String()),
			Labels:       pulumi.ToStringMap(s.GcpLabels),
		}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to create npm repo")
	}

	//grant "reader" role for the writer service account on the repo
	_, err = artifactregistry.NewRepositoryIamMember(ctx,
		fmt.Sprintf("%s-reader", npmRepoName),
		&artifactregistry.RepositoryIamMemberArgs{
			Project:    pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
			Location:   pulumi.String(gcpArtifactRegistry.Spec.Region),
			Repository: createdNpmRepo.RepositoryId,
			Role:       pulumi.String("roles/artifactregistry.reader"),
			Member:     pulumi.Sprintf("serviceAccount:%s", readerServiceAccount.Email),
		}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to grant reader role on npm repo for reader service account")
	}

	//grant "writer" role for the writer service account on the repo
	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-writer",
		npmRepoName), &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
		Location:   pulumi.String(gcpArtifactRegistry.Spec.Region),
		Repository: createdNpmRepo.RepositoryId,
		Role:       pulumi.String("roles/artifactregistry.writer"),
		Member:     pulumi.Sprintf("serviceAccount:%s", writerServiceAccount.Email),
	}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to grant writer role on npm repo for writer service account")
	}

	//grant "admin" role for writer service account on the repo
	_, err = artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("%s-admin",
		npmRepoName), &artifactregistry.RepositoryIamMemberArgs{
		Project:    pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
		Location:   pulumi.String(gcpArtifactRegistry.Spec.Region),
		Repository: createdNpmRepo.RepositoryId,
		Role:       pulumi.String("roles/artifactregistry.repoAdmin"),
		Member:     pulumi.Sprintf("serviceAccount:%s", writerServiceAccount.Email),
	}, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to grant admin role on npm repo for writer service account")
	}

	//export the name of the npm repository as output
	ctx.Export(outputs.NpmRepoName, createdNpmRepo.RepositoryId)

	return nil
}
