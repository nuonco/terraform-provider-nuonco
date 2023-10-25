package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnvVar struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type EnvVarSlice []EnvVar

func NewEnvVarSliceFromMap(stringMap map[string]string) EnvVarSlice {
	blocks := []EnvVar{}
	for key, val := range stringMap {
		blocks = append(blocks, EnvVar{
			Name:  types.StringValue(key),
			Value: types.StringValue(val),
		})
	}
	return blocks
}

func (ev *EnvVarSlice) ToMap() map[string]string {
	stringMap := map[string]string{}
	for _, val := range *ev {
		stringMap[val.Name.ValueString()] = val.Value.ValueString()
	}
	return stringMap
}

func (ev *EnvVarSlice) String() string {
	return fmt.Sprintf("%#v", ev)
}

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
