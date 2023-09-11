package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnvVar struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func envVarSharedBlock() schema.ListNestedBlock {
	return schema.ListNestedBlock{
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

func helmValueSharedBlock() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Helm values to set when deploying the Helm chart.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "The name of the value to set in the chart (e.g. env.secret or server.container.image.)",
					Required:    true,
				},
				"value": schema.StringAttribute{
					Description: "The value to set in the chart. Can be any valid Helm chart value, or interpolated from Nuon.",
					Required:    true,
				},
			},
		},
	}
}

type TerraformVariable struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func terraformVariableSharedBlock() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description: "Terraform variables to set when applying the Terraform configuration.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "The variable name to write to the terraform.tfvars file (e.g. bucket_name or db_name.)",
					Required:    true,
				},
				"value": schema.StringAttribute{
					Description: "The variable value to write to the terraform.tfvars file. Can be any valid Terraform value, or interpolated from Nuon.",
					Required:    true,
				},
			},
		},
	}
}
