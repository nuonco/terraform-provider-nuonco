package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/nuonco/nuon-go/models"
)

func convertSandboxRelease(sandboxRelease models.AppSandboxRelease) basetypes.ObjectValue {
	obj, _ := basetypes.NewObjectValue(
		map[string]attr.Type{
			"id":                          types.StringType,
			"version":                     types.StringType,
			"provision_policy_url":        types.StringType,
			"deprovision_policy_url":      types.StringType,
			"trust_policy_url":            types.StringType,
			"one_click_role_template_url": types.StringType,
		},
		map[string]attr.Value{
			"id":                          types.StringValue(sandboxRelease.ID),
			"version":                     types.StringValue(sandboxRelease.Version),
			"provision_policy_url":        types.StringValue(sandboxRelease.ProvisionPolicyURL),
			"deprovision_policy_url":      types.StringValue(sandboxRelease.DeprovisionPolicyURL),
			"trust_policy_url":            types.StringValue(sandboxRelease.TrustPolicyURL),
			"one_click_role_template_url": types.StringValue(sandboxRelease.OneClickRoleTemplateURL),
		},
	)
	return obj
}

// convert from a []string{} to a Terraform List
func stringSliceToList(ctx context.Context, stringSlice []string) types.List {
	list, _ := types.ListValueFrom(ctx, types.StringType, stringSlice)
	return list
}

// convert a Terraform List to a []string{}
func listToStringSlice(list types.List) []string {
	stringSlice := []string{}
	elements := list.Elements()
	for _, el := range elements {
		str := el.String()
		str = str[1 : len(str)-1]
		stringSlice = append(stringSlice, str)
	}
	return stringSlice
}
