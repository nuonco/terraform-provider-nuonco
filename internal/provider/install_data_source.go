package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &InstallDataSource{}

func NewInstallDataSource() datasource.DataSource {
	return &InstallDataSource{}
}

// InstallDataSource defines the data source implementation.
type InstallDataSource struct {
	baseDataSource
}

// InstallDataSourceModel describes the data source data model.
type InstallDataSourceModel struct {
	Name types.String `tfsdk:"name"`
	Id   types.String `tfsdk:"id"`
}

func (d *InstallDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_install"
}

func (d *InstallDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides information about a Nuon install. This data source can be useful if you need to begin programmatically managing an install created in the UI.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The human-readable name of the install.",
				Optional:    true,
			},
			"id": schema.StringAttribute{
				Description: "The unique ID of the install.",
				Optional:    true,
			},
		},
	}
}

func (d *InstallDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InstallDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	installResp, err := d.restClient.GetInstall(ctx, data.Id.ValueString())
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "get install")
		return
	}
	data.Name = types.StringValue(installResp.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
