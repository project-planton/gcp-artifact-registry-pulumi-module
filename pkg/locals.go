package pkg

import (
	"github.com/plantoncloud/project-planton/apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gcpartifactregistry"
	"github.com/plantoncloud/project-planton/apis/zzgo/cloud/planton/apis/commons/apiresource/enums/apiresourcekind"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strconv"
)

type Locals struct {
	GcpArtifactRegistry *gcpartifactregistry.GcpArtifactRegistry
	GcpLabels           map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpartifactregistry.GcpArtifactRegistryStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.GcpArtifactRegistry = stackInput.Target

	locals.GcpLabels = map[string]string{
		gcplabelkeys.Resource:     strconv.FormatBool(true),
		gcplabelkeys.Organization: locals.GcpArtifactRegistry.Spec.EnvironmentInfo.OrgId,
		gcplabelkeys.ResourceKind: apiresourcekind.ApiResourceKind_gcp_artifact_registry.String(),
		gcplabelkeys.ResourceId:   locals.GcpArtifactRegistry.Metadata.Id,
	}

	return locals
}
