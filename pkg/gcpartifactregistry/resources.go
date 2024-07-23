package gcpartifactregistry

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gcpartifactregistry/model"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input     *model.GcpArtifactRegistryStackInput
	GcpLabels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	gcpProvider, err := pulumigoogleprovider.Get(ctx, s.Input.GcpCredential)
	if err != nil {
		return errors.Wrap(err, "failed to setup gcp provider")
	}

	addedReaderServiceAccount, addedWriterServiceAccount, err := s.serviceAccounts(ctx, gcpProvider)
	if err != nil {
		return errors.Wrap(err, "failed to add service accounts")
	}

	if err := s.dockerRepo(ctx, gcpProvider, addedReaderServiceAccount, addedWriterServiceAccount); err != nil {
		return errors.Wrap(err, "failed to add docker repo")
	}

	if err := s.mavenRepo(ctx, gcpProvider, addedReaderServiceAccount, addedWriterServiceAccount); err != nil {
		return errors.Wrap(err, "failed to add docker repo")
	}

	if err := s.npmRepo(ctx, gcpProvider, addedReaderServiceAccount, addedWriterServiceAccount); err != nil {
		return errors.Wrap(err, "failed to add docker repo")
	}

	if err := s.pythonRepo(ctx, gcpProvider, addedReaderServiceAccount, addedWriterServiceAccount); err != nil {
		return errors.Wrap(err, "failed to add python repo")
	}

	return nil
}

func (s *ResourceStack) getGcpArtifactRegistryId() string {
	return s.Input.ApiResource.Metadata.Id
}
