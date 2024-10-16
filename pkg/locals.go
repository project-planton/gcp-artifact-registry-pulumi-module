package pkg

import (
	gcpartifactregistryv1 "buf.build/gen/go/project-planton/apis/protocolbuffers/go/project/planton/provider/gcp/gcpartifactregistry/v1"
	"github.com/project-planton/pulumi-module-golang-commons/pkg/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strconv"
)

type Locals struct {
	GcpArtifactRegistry *gcpartifactregistryv1.GcpArtifactRegistry
	GcpLabels           map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpartifactregistryv1.GcpArtifactRegistryStackInput) *Locals {
	locals := &Locals{}

	//if the id is empty, use name as id
	if stackInput.Target.Metadata.Id == "" {
		stackInput.Target.Metadata.Id = stackInput.Target.Metadata.Name
	}

	//assign value for the locals variable to make it available across the project
	locals.GcpArtifactRegistry = stackInput.Target

	locals.GcpLabels = map[string]string{
		gcplabelkeys.Resource:     strconv.FormatBool(true),
		gcplabelkeys.ResourceId:   locals.GcpArtifactRegistry.Metadata.Id,
		gcplabelkeys.ResourceKind: "gcp_artifact_registry",
	}

	if locals.GcpArtifactRegistry.Spec.EnvironmentInfo != nil && locals.GcpArtifactRegistry.Spec.EnvironmentInfo.OrgId != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpArtifactRegistry.Spec.EnvironmentInfo.OrgId
	}

	return locals
}
