package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &BuiltinSandboxDataSource{}

func NewBuiltinSandboxDataSource() datasource.DataSource {
	return &BuiltinSandboxDataSource{}
}

// BuiltinSandboxDataSource defines the data source implementation.
type BuiltinSandboxDataSource struct {
	baseDataSource
}

// BuiltinSandboxDataSourceModel describes the data source data model.
type BuiltinSandboxDataSourceModel struct {
	Name           types.String          `tfsdk:"name"`
	Id             types.String          `tfsdk:"id"`
	SandboxRelease basetypes.ObjectValue `tfsdk:"sandbox_release"`
}

func (d *BuiltinSandboxDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_builtin_sandbox"
}

func (d *BuiltinSandboxDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides information about a built in sandbox",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The sandbox name",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "The unique ID of the app.",
				Computed:    true,
			},
			"sandbox_release": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The latest sandbox release for the built in sandbox",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"version": schema.StringAttribute{
						Computed: true,
					},
					"terraform_version": schema.StringAttribute{
						Computed: true,
					},
					"provision_policy_url": schema.StringAttribute{
						Computed: true,
					},
					"deprovision_policy_url": schema.StringAttribute{
						Computed: true,
					},
					"trust_policy_url": schema.StringAttribute{
						Computed: true,
					},
					"one_click_role_template_url": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

func (d *BuiltinSandboxDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BuiltinSandboxDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "fetching built in sandbox by name")
	sandboxResp, err := d.restClient.GetSandbox(ctx, data.Name.ValueString())
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "get sandbox")
		return
	}

	data.Name = types.StringValue(sandboxResp.Name)
	data.Id = types.StringValue(sandboxResp.ID)
	if len(sandboxResp.Releases) < 1 {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, fmt.Errorf("sandbox did not return any releases."), "invalid sandbox")
		return
	}
	data.SandboxRelease = convertSandboxRelease(*sandboxResp.Releases[0])

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
