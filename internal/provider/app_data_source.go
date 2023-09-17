package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &AppDataSource{}

func NewAppDataSource() datasource.DataSource {
	return &AppDataSource{}
}

// AppDataSource defines the data source implementation.
type AppDataSource struct {
	baseDataSource
}

// AppDataSourceModel describes the data source data model.
type AppDataSourceModel struct {
	Name	       types.String	     `tfsdk:"name"`
	Id	       types.String	     `tfsdk:"id"`
	SandboxRelease basetypes.ObjectValue `tfsdk:"sandbox_release"`
}

func (d *AppDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (d *AppDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides information about a Nuon app.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The human readable name of the app.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The unique ID of the app.",
				Required:    true,
			},
			"sandbox_release": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The sandbox being used for this app's installs.",
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

func (d *AppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AppDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "fetching app by id")
	appResp, err := d.restClient.GetApp(ctx, data.Id.ValueString())
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "get app")
		return
	}

	data.Name = types.StringValue(appResp.Name)
	data.Id = types.StringValue(appResp.ID)
	data.SandboxRelease = convertSandboxRelease(*appResp.SandboxRelease)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
