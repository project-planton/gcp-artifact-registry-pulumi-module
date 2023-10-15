package gcp

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/repo"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/serviceaccount"
	pulumigcpprovider "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/automation/provider/google"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/develop/artifactstore/stack/gcp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input     *gcp.ArtifactStoreGcpStackInput
	GcpLabels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	gcpProvider, err := pulumigcpprovider.Get(ctx, s.Input.CredentialsInput.Google)
	if err != nil {
		return errors.Wrap(err, "failed to setup gcp provider")
	}
	serviceAccountResources, err := serviceaccount.Resources(ctx, &serviceaccount.Input{
		GcpProvider:     gcpProvider,
		GcpProjectId:    s.Input.ResourceInput.ArtifactStore.Spec.GcpArtifactRegistry.ProjectId,
		ArtifactStoreId: s.Input.ResourceInput.ArtifactStore.Metadata.Id,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add service accounts")
	}
	if err := repo.Resources(ctx, &repo.Input{
		GcpProvider:             gcpProvider,
		ArtifactStoreId:         s.Input.ResourceInput.ArtifactStore.Metadata.Id,
		GcpProjectId:            s.Input.ResourceInput.ArtifactStore.Spec.GcpArtifactRegistry.ProjectId,
		GcpRegion:               s.Input.ResourceInput.ArtifactStore.Spec.GcpArtifactRegistry.Region,
		ServiceAccountResources: serviceAccountResources,
		IsExternal:              s.Input.ResourceInput.ArtifactStore.Spec.GcpArtifactRegistry.IsExternal,
		Labels:                  s.GcpLabels,
	}); err != nil {
		return errors.Wrap(err, "failed to add repo resources")
	}
	return nil
}