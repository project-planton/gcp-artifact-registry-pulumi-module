package serviceaccount

import (
	"fmt"
	"github.com/pkg/errors"
	commonsgcpiamsa "github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/serviceaccount"
	puluminameoutputgcp "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/gcp/output"
	wordpb "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/commons/english/rpc/enums"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	readOnlyServiceAccountNameSuffix  = "ro"
	readWriteServiceAccountNameSuffix = "rw"
)

type Input struct {
	GcpProvider     *pulumigcp.Provider
	GcpProjectId    string
	ArtifactStoreId string
}

type AddedResources struct {
	ReaderServiceAccount *serviceaccount.Account
	WriterServiceAccount *serviceaccount.Account
}

func Resources(ctx *pulumi.Context, input *Input) (*AddedResources, error) {
	readerServiceAccountFullName := getReaderServiceAccountName(input.ArtifactStoreId)
	readerServiceAccount, readerServiceAccountKey, err := addServiceAccount(ctx, input, readerServiceAccountFullName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s svc acct", readerServiceAccountFullName)
	}
	writerServiceAccountFullName := getWriterServiceAccountName(input.ArtifactStoreId)
	writerServiceAccount, writerServiceAccountKey, err := addServiceAccount(ctx, input, writerServiceAccountFullName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s svc acct", writerServiceAccountFullName)
	}

	exportOutputs(ctx, input.ArtifactStoreId, readerServiceAccount, writerServiceAccount, readerServiceAccountKey, writerServiceAccountKey)

	return &AddedResources{
		ReaderServiceAccount: readerServiceAccount,
		WriterServiceAccount: writerServiceAccount,
	}, nil
}

func addServiceAccount(ctx *pulumi.Context, input *Input, serviceAccountName string) (*serviceaccount.Account, *serviceaccount.Key, error) {
	sa, err := serviceaccount.NewAccount(ctx, serviceAccountName, &serviceaccount.AccountArgs{
		Project:     pulumi.String(input.GcpProjectId),
		AccountId:   pulumi.String(serviceAccountName),
		DisplayName: pulumi.String(serviceAccountName),
	}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed add new %s svc acct", serviceAccountName)
	}
	serviceAccountKey, err := serviceaccount.NewKey(ctx, serviceAccountName, &serviceaccount.KeyArgs{
		ServiceAccountId: sa.Name,
		PublicKeyType:    pulumi.String(commonsgcpiamsa.KeyTypeX509PemFile),
	}, pulumi.Parent(sa))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create key for %s svc acct", serviceAccountName)
	}
	return sa, serviceAccountKey, nil
}

func getReaderServiceAccountName(artifactStoreId string) string {
	return fmt.Sprintf("%s-%s", artifactStoreId, readOnlyServiceAccountNameSuffix)
}

func getWriterServiceAccountName(artifactStoreId string) string {
	return fmt.Sprintf("%s-%s", artifactStoreId, readWriteServiceAccountNameSuffix)
}

func GetReaderServiceAccountEmailOutputName(artifactStoreId string) string {
	return puluminameoutputgcp.Name(serviceaccount.Account{}, getReaderServiceAccountName(artifactStoreId), wordpb.Word_email.String())
}

func GetReaderServiceAccountKeyOutputName(artifactStoreId string) string {
	return puluminameoutputgcp.Name(serviceaccount.Key{}, getReaderServiceAccountName(artifactStoreId), wordpb.Word_key.String())
}

func GetWriterServiceAccountEmailOutputName(artifactStoreId string) string {
	return puluminameoutputgcp.Name(serviceaccount.Account{}, getWriterServiceAccountName(artifactStoreId), wordpb.Word_email.String())
}

func GetWriterServiceAccountKeyOutputName(artifactStoreId string) string {
	return puluminameoutputgcp.Name(serviceaccount.Key{}, getWriterServiceAccountName(artifactStoreId), wordpb.Word_key.String())
}

func exportOutputs(ctx *pulumi.Context, artifactStoreId string, readerServiceAccount,
	writerServiceAccount *serviceaccount.Account, readerServiceAccountKey, writerServiceAccountKey *serviceaccount.Key) {
	ctx.Export(GetReaderServiceAccountEmailOutputName(artifactStoreId), readerServiceAccount.Email)
	ctx.Export(GetReaderServiceAccountKeyOutputName(artifactStoreId), readerServiceAccountKey.PrivateKey)
	ctx.Export(GetWriterServiceAccountEmailOutputName(artifactStoreId), writerServiceAccount.Email)
	ctx.Export(GetWriterServiceAccountKeyOutputName(artifactStoreId), writerServiceAccountKey.PrivateKey)
}
