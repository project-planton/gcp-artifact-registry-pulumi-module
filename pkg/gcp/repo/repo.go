package repo

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/repo/docker"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/repo/maven"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/repo/npm"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/repo/python"
	"github.com/plantoncloud-inc/artifact-store-pulumi-blueprint/pkg/gcp/serviceaccount"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	GcpProvider             *pulumigcp.Provider
	ArtifactStoreId         string
	GcpProjectId            string
	GcpRegion               string
	ServiceAccountResources *serviceaccount.AddedResources
	IsExternal              bool
	Labels                  map[string]string
}

func Resources(ctx *pulumi.Context, input *Input) error {
	err := docker.Resources(ctx, &docker.Input{
		GcpProvider:          input.GcpProvider,
		ArtifactStoreId:      input.ArtifactStoreId,
		GcpProjectId:         input.GcpProjectId,
		GcpRegion:            input.GcpRegion,
		ServiceAccountOutput: input.ServiceAccountResources,
		IsExternal:           input.IsExternal,
		Labels:               input.Labels,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add docker repo")
	}
	err = maven.Resources(ctx, &maven.Input{
		GcpProvider:          input.GcpProvider,
		ArtifactStoreId:      input.ArtifactStoreId,
		GcpProjectId:         input.GcpProjectId,
		GcpRegion:            input.GcpRegion,
		ServiceAccountOutput: input.ServiceAccountResources,
		IsExternal:           input.IsExternal,
		Labels:               input.Labels,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add maven repo")
	}
	err = npm.Resources(ctx, &npm.Input{
		GcpProvider:          input.GcpProvider,
		ArtifactStoreId:      input.ArtifactStoreId,
		GcpProjectId:         input.GcpProjectId,
		GcpRegion:            input.GcpRegion,
		ServiceAccountOutput: input.ServiceAccountResources,
		IsExternal:           input.IsExternal,
		Labels:               input.Labels,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add npm repo")
	}
	err = python.Resources(ctx, &python.Input{
		GcpProvider:          input.GcpProvider,
		ArtifactStoreId:      input.ArtifactStoreId,
		GcpProjectId:         input.GcpProjectId,
		GcpRegion:            input.GcpRegion,
		ServiceAccountOutput: input.ServiceAccountResources,
		IsExternal:           input.IsExternal,
		Labels:               input.Labels,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add python repo")
	}
	return nil
}
