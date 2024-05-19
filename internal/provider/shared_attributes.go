package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
				Description: "The directory the component code is in. Use ./ for root.",
				Required:    true,
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
