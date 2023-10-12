package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func envVarSharedBlock() schema.SetNestedBlock {
	return schema.SetNestedBlock{
		Description: "Environment variables to export into the env when running the image.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "The variable name to export to the env (e.g. API_TOKEN or PORT.)",
					Required:    true,
				},
				"value": schema.StringAttribute{
					Description: "The variable value to export to the env. Can be any valid env var value, or interpolated from Nuon.",
					Required:    true,
				},
			},
		},
	}
}
