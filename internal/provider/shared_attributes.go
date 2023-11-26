package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PublicRepo struct {
	Repo      types.String `tfsdk:"repo"`
	Directory types.String `tfsdk:"directory"`
	Branch    types.String `tfsdk:"branch"`
}

func publicRepoAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "A publically-accessible git repo.",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"repo": schema.StringAttribute{
				Description: "The https clone url",
				Required:    true,
			},
			"branch": schema.StringAttribute{
				Description: "The default branch to create new builds from.",
				Required:    true,
			},
			"directory": schema.StringAttribute{
				Description: "The directory the component code is in.",
				Optional:    true,
			},
		},
	}
}

type ConnectedRepo struct {
	Repo      types.String `tfsdk:"repo"`
	Branch    types.String `tfsdk:"branch"`
	Directory types.String `tfsdk:"directory"`
}

func connectedRepoAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "A repo accessible via your Nuon connected github account",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"repo": schema.StringAttribute{
				Description: "The https clone url",
				Required:    true,
			},
			"branch": schema.StringAttribute{
				Description: "The default branch to create new builds from.",
				Required:    true,
			},
			"directory": schema.StringAttribute{
				Description: "The directory the component code is in.",
				Optional:    true,
			},
		},
	}
}

type BasicDeploy struct {
	Port            types.Int64  `tfsdk:"port"`
	InstanceCount   types.Int64  `tfsdk:"instance_count"`
	HealthCheckPath types.String `tfsdk:"health_check_path"`
}

func basicDeployAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Create a basic deployment of this image with a listener.",
		Attributes: map[string]schema.Attribute{
			"port": schema.Int64Attribute{
				Description: "The port to listen on.",
				Required:    true,
			},
			"instance_count": schema.Int64Attribute{
				Description: "The number of instances to run.",
				Default:     int64default.StaticInt64(1),
				Optional:    true,
				Computed:    true,
			},
			"health_check_path": schema.StringAttribute{
				Description: "The path to use for health checks.",
				Optional:    true,
				Default:     stringdefault.StaticString("/"),
				Computed:    true,
			},
		},
	}
}
