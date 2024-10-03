package pkg

import (
	gcpartifactregistryv1 "buf.build/gen/go/plantoncloud/project-planton/protocolbuffers/go/project/planton/provider/gcp/gcpartifactregistry/v1"
	"encoding/base64"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *gcpartifactregistryv1.GcpArtifactRegistryStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	gcpProviderArgs := &gcp.ProviderArgs{}

	//if stack-input contains gcp-credentials, provider will be created with those credentials
	if stackInput.GcpCredential != nil {
		serviceAccountKey, err := base64.StdEncoding.DecodeString(stackInput.GcpCredential.ServiceAccountKeyBase64)
		if err != nil {
			return errors.Wrap(err, "failed to decode base64 encoded"+
				" google service account credential")
		}
		gcpProviderArgs = &gcp.ProviderArgs{Credentials: pulumi.String(serviceAccountKey)}
	}

	googleProvider, err := gcp.NewProvider(ctx, "google", gcpProviderArgs)
	if err != nil {
		return errors.Wrap(err, "failed to create google provider")
	}

	//create reader and writer service accounts and their keys
	addedReaderServiceAccount, addedWriterServiceAccount, err := serviceAccounts(ctx, locals, googleProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create service accounts")
	}

	//create docker repository and grant reader and writer roles for the service accounts on the repo
	if err := dockerRepo(ctx, locals, googleProvider, addedReaderServiceAccount, addedWriterServiceAccount); err != nil {
		return errors.Wrap(err, "failed to create docker repo")
	}

	//create maven repository and grant reader and writer roles for the service accounts on the repo
	if err := mavenRepo(ctx, locals, googleProvider, addedReaderServiceAccount, addedWriterServiceAccount); err != nil {
		return errors.Wrap(err, "failed to create maven repo")
	}

	//create npm repository and grant reader and writer roles for the service accounts on the repo
	if err := npmRepo(ctx, locals, googleProvider, addedReaderServiceAccount, addedWriterServiceAccount); err != nil {
		return errors.Wrap(err, "failed to create npm repo")
	}

	//create python repository and grant reader and writer roles for the service accounts on the repo
	if err := pythonRepo(ctx, locals, googleProvider, addedReaderServiceAccount, addedWriterServiceAccount); err != nil {
		return errors.Wrap(err, "failed to create python repo")
	}

	return nil
}
