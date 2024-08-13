package pkg

import (
	"github.com/pkg/errors"
	gcpartifactregistry "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gcpartifactregistry"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input     *gcpartifactregistry.GcpArtifactRegistryStackInput
	GcpLabels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	//create google provider using the credentials from the input
	googleProvider, err := pulumigoogleprovider.Get(ctx, s.Input.GcpCredential)
	if err != nil {
		return errors.Wrap(err, "failed to create google provider")
	}

	//create reader and writer service accounts and their keys
	addedReaderServiceAccount, addedWriterServiceAccount, err := s.serviceAccounts(ctx, googleProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create service accounts")
	}

	//create docker repository and grant reader and writer roles for the service accounts on the repo
	if err := s.dockerRepo(ctx, googleProvider, addedReaderServiceAccount, addedWriterServiceAccount); err != nil {
		return errors.Wrap(err, "failed to create docker repo")
	}

	//create maven repository and grant reader and writer roles for the service accounts on the repo
	if err := s.mavenRepo(ctx, googleProvider, addedReaderServiceAccount, addedWriterServiceAccount); err != nil {
		return errors.Wrap(err, "failed to create maven repo")
	}

	//create npm repository and grant reader and writer roles for the service accounts on the repo
	if err := s.npmRepo(ctx, googleProvider, addedReaderServiceAccount, addedWriterServiceAccount); err != nil {
		return errors.Wrap(err, "failed to create npm repo")
	}

	//create python repository and grant reader and writer roles for the service accounts on the repo
	if err := s.pythonRepo(ctx, googleProvider, addedReaderServiceAccount, addedWriterServiceAccount); err != nil {
		return errors.Wrap(err, "failed to create python repo")
	}

	return nil
}
