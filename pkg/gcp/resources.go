package gcp

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/artifact-store-pulumi-blueprint/pkg/gcp/repo"
	"github.com/plantoncloud/artifact-store-pulumi-blueprint/pkg/gcp/serviceaccount"
	code2cloudv1developafsstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/artifactstore/stack/gcp/model"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input     *code2cloudv1developafsstackgcpmodel.ArtifactStoreGcpStackInput
	GcpLabels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	gcpProvider, err := pulumigoogleprovider.Get(ctx, s.Input.CredentialsInput.Google)
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
