package gcpartifactregistry

import (
	"fmt"
	"github.com/pkg/errors"
	commonsgcpiamsa "github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/serviceaccount"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gcpartifactregistry/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/gcp/pulumigoogleprovider"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	readOnlyServiceAccountNameSuffix  = "ro"
	readWriteServiceAccountNameSuffix = "rw"
)

func (s *ResourceStack) serviceAccounts(ctx *pulumi.Context, gcpProvider *pulumigcp.Provider) (addedReaderServiceAccount,
	WriterServiceAccount *serviceaccount.Account, err error) {
	readerServiceAccountFullName := getReaderServiceAccountName(s.getGcpArtifactRegistryId())

	addedReaderServiceAccount, addedReaderServiceAccountKey, err := addServiceAccount(ctx, gcpProvider,
		s.Input.ApiResource, readerServiceAccountFullName)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to add %s svc acct",
			readerServiceAccountFullName)
	}

	writerServiceAccountFullName := getWriterServiceAccountName(s.getGcpArtifactRegistryId())
	addedWriterServiceAccount, addedWriterServiceAccountKey, err := addServiceAccount(ctx, gcpProvider,
		s.Input.ApiResource, writerServiceAccountFullName)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to add %s svc acct",
			writerServiceAccountFullName)
	}

	ctx.Export(GetReaderServiceAccountEmailOutputName(s.getGcpArtifactRegistryId()), addedReaderServiceAccount.Email)
	ctx.Export(GetReaderServiceAccountKeyOutputName(s.getGcpArtifactRegistryId()), addedReaderServiceAccountKey.PrivateKey)
	ctx.Export(GetWriterServiceAccountEmailOutputName(s.getGcpArtifactRegistryId()), addedWriterServiceAccount.Email)
	ctx.Export(GetWriterServiceAccountKeyOutputName(s.getGcpArtifactRegistryId()), addedWriterServiceAccountKey.PrivateKey)

	return addedReaderServiceAccount, addedWriterServiceAccount, nil
}

func addServiceAccount(ctx *pulumi.Context, gcpProvider *pulumigcp.Provider,
	gcpArtifactRegistry *model.GcpArtifactRegistry, serviceAccountName string) (*serviceaccount.Account,
	*serviceaccount.Key, error) {
	addedServiceAccount, err := serviceaccount.NewAccount(ctx, serviceAccountName, &serviceaccount.AccountArgs{
		Project:     pulumi.String(gcpArtifactRegistry.Spec.ProjectId),
		AccountId:   pulumi.String(serviceAccountName),
		DisplayName: pulumi.String(serviceAccountName),
	}, pulumi.Provider(gcpProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed add new %s svc acct", serviceAccountName)
	}

	addedServiceAccountKey, err := serviceaccount.NewKey(ctx, serviceAccountName, &serviceaccount.KeyArgs{
		ServiceAccountId: addedServiceAccount.Name,
		PublicKeyType:    pulumi.String(commonsgcpiamsa.KeyTypeX509PemFile),
	}, pulumi.Parent(addedServiceAccount))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create key for %s svc acct",
			serviceAccountName)
	}

	return addedServiceAccount, addedServiceAccountKey, nil
}

func getReaderServiceAccountName(gcpArtifactRegistryId string) string {
	return fmt.Sprintf("%s-%s", gcpArtifactRegistryId, readOnlyServiceAccountNameSuffix)
}

func getWriterServiceAccountName(gcpArtifactRegistryId string) string {
	return fmt.Sprintf("%s-%s", gcpArtifactRegistryId, readWriteServiceAccountNameSuffix)
}

func GetReaderServiceAccountEmailOutputName(gcpArtifactRegistryId string) string {
	return pulumigoogleprovider.PulumiOutputName(serviceaccount.Account{},
		getReaderServiceAccountName(gcpArtifactRegistryId), englishword.EnglishWord_email.String())
}

func GetReaderServiceAccountKeyOutputName(gcpArtifactRegistryId string) string {
	return pulumigoogleprovider.PulumiOutputName(serviceaccount.Key{},
		getReaderServiceAccountName(gcpArtifactRegistryId), englishword.EnglishWord_key.String())
}

func GetWriterServiceAccountEmailOutputName(gcpArtifactRegistryId string) string {
	return pulumigoogleprovider.PulumiOutputName(serviceaccount.Account{},
		getWriterServiceAccountName(gcpArtifactRegistryId), englishword.EnglishWord_email.String())
}

func GetWriterServiceAccountKeyOutputName(gcpArtifactRegistryId string) string {
	return pulumigoogleprovider.PulumiOutputName(serviceaccount.Key{},
		getWriterServiceAccountName(gcpArtifactRegistryId), englishword.EnglishWord_key.String())
}
