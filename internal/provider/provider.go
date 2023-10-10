// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nuonco/nuon-go"
	"github.com/nuonco/terraform-provider-nuon/internal/config"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &Provider{}

// Provider defines the provider implementation.
type Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ProviderModel describes the provider data model.
type ProviderModel struct {
	APIAuthToken types.String `tfsdk:"api_token"`
	OrgID        types.String `tfsdk:"org_id"`
}

type ProviderData struct {
	OrgID      string
	RestClient nuon.Client
}

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nuon"
	resp.Version = p.version
}

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Terraform provider for managing apps on the Nuon platform.",
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Description: "A valid API token to access the api.",
				Optional:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "Your Nuon organization ID.",
				Optional:    true,
			},
		},
	}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read sdk config from config file, env vars, then terraform
	cfg, err := config.NewConfig("")
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "initialize nuon")
		return
	}

	apiToken := cfg.APIToken
	if val := data.APIAuthToken.ValueString(); val != "" {
		apiToken = val
	}

	orgID := cfg.OrgID
	if val := data.OrgID.ValueString(); val != "" {
		orgID = val
	}

	// initialize sdk
	restClient, err := nuon.New(
		validator.New(),
		nuon.WithAuthToken(apiToken),
		nuon.WithOrgID(orgID),
		nuon.WithURL(cfg.APIURL),
	)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "initialize nuon")
		return
	}

	resp.DataSourceData = &ProviderData{
		RestClient: restClient,
	}
	resp.ResourceData = &ProviderData{
		RestClient: restClient,
	}
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAppResource,
		NewAppInstallerResource,
		NewInstallResource,
		NewContainerImageComponentResource,
		NewDockerBuildComponentResource,
		NewHelmChartComponentResource,
		NewTerraformModuleComponentResource,
	}
}

func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAppDataSource,
		NewConnectedRepoDataSource,
		NewInstallDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{
			version: version,
		}
	}
}
