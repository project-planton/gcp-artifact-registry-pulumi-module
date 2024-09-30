package pkg

import (
	gcpartifactregistryv1 "buf.build/gen/go/plantoncloud/project-planton/protocolbuffers/go/project/planton/apis/provider/gcp/gcpartifactregistry/v1"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strconv"
)

type Locals struct {
	GcpArtifactRegistry *gcpartifactregistryv1.GcpArtifactRegistry
	GcpLabels           map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpartifactregistryv1.GcpArtifactRegistryStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.GcpArtifactRegistry = stackInput.Target

	locals.GcpLabels = map[string]string{
		gcplabelkeys.Resource:     strconv.FormatBool(true),
		gcplabelkeys.Organization: locals.GcpArtifactRegistry.Spec.EnvironmentInfo.OrgId,
		gcplabelkeys.ResourceKind: "gcp_artifact_registry",
		gcplabelkeys.ResourceId:   locals.GcpArtifactRegistry.Metadata.Id,
	}

	return locals
}
