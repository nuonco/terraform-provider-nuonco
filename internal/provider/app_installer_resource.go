package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nuonco/nuon-go"
	"github.com/nuonco/nuon-go/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AppInstallerResource{}
var _ resource.ResourceWithImportState = &AppInstallerResource{}

func NewAppInstallerResource() resource.Resource {
	return &AppInstallerResource{}
}

// AppInstallerResource defines the resource implementation.
type AppInstallerResource struct {
	baseResource
}

// AppInstallerResourceModel describes the resource data model.
type AppInstallerResourceModel struct {
	Id types.String `tfsdk:"id"`

	AppID types.String `tfsdk:"app_id"`

	Slug types.String `tfsdk:"slug"`

	// metadata
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	DocumentationURL types.String `tfsdk:"documentation_url"`
	HomepageURL      types.String `tfsdk:"homepage_url"`
	CommunityURL     types.String `tfsdk:"community_url"`
	GithubURL        types.String `tfsdk:"github_url"`
	LogoURL          types.String `tfsdk:"logo_url"`
	DemoURL          types.String `tfsdk:"demo_url"`
}

func (r *AppInstallerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_installer"
}

func (r *AppInstallerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A public installer page for a nuon app.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "App name to render on install page.",
				Optional:            false,
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Short description of app.",
				Optional:            false,
				Required:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL slug to access app from.",
				Optional:            false,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"documentation_url": schema.StringAttribute{
				MarkdownDescription: "Documentation url",
				Optional:            false,
				Required:            true,
			},
			"homepage_url": schema.StringAttribute{
				MarkdownDescription: "Homepage url",
				Optional:            false,
				Required:            true,
			},
			"community_url": schema.StringAttribute{
				MarkdownDescription: "Community url to a slack or discord, etc.",
				Optional:            false,
				Required:            true,
			},
			"github_url": schema.StringAttribute{
				MarkdownDescription: "GitHub url, link to application.",
				Optional:            false,
				Required:            true,
			},
			"logo_url": schema.StringAttribute{
				MarkdownDescription: "Logo image to display on page.",
				Optional:            false,
				Required:            true,
			},
			"demo_url": schema.StringAttribute{
				MarkdownDescription: "Demo url to show",
				Optional:            true,
				Required:            false,
			},
			"app_id": schema.StringAttribute{
				Description: "The application ID.",
				Optional:    false,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique ID of the app.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *AppInstallerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// get terraform model
	var data *AppInstallerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create app
	tflog.Trace(ctx, "creating app installer")
	appResp, err := r.restClient.CreateAppInstaller(ctx, &models.ServiceCreateAppInstallerRequest{
		AppID:       data.AppID.ValueStringPointer(),
		Name:        data.Name.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Slug:        data.Slug.ValueStringPointer(),
		Links: &models.ServiceCreateAppInstallerRequestLinks{
			Community:     data.CommunityURL.ValueStringPointer(),
			Documentation: data.DocumentationURL.ValueStringPointer(),
			Homepage:      data.HomepageURL.ValueStringPointer(),
			Github:        data.GithubURL.ValueStringPointer(),
			Logo:          data.LogoURL.ValueStringPointer(),
			Demo:          data.DemoURL.ValueString(),
		},
	})
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create app installer")
		return
	}

	// update the state with the returned values
	data.Id = types.StringValue(appResp.ID)
	data.Slug = types.StringValue(appResp.Slug)
	data.Name = types.StringValue(appResp.AppInstallerMetadata.Name)
	data.Description = types.StringValue(appResp.AppInstallerMetadata.Description)
	data.CommunityURL = types.StringValue(appResp.AppInstallerMetadata.CommunityURL)
	data.GithubURL = types.StringValue(appResp.AppInstallerMetadata.GithubURL)
	data.DocumentationURL = types.StringValue(appResp.AppInstallerMetadata.DocumentationURL)
	data.HomepageURL = types.StringValue(appResp.AppInstallerMetadata.HomepageURL)
	data.LogoURL = types.StringValue(appResp.AppInstallerMetadata.LogoURL)

	if appResp.AppInstallerMetadata.DemoURL != "" {
		data.DemoURL = types.StringValue(appResp.AppInstallerMetadata.DemoURL)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppInstallerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AppInstallerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading app installer")

	appResp, err := r.restClient.GetAppInstaller(ctx, data.Id.ValueString())
	if nuon.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "read app installer")
		return
	}

	// update the state with the returned values
	data.Id = types.StringValue(appResp.ID)
	data.Slug = types.StringValue(appResp.Slug)
	data.Name = types.StringValue(appResp.AppInstallerMetadata.Name)
	data.Description = types.StringValue(appResp.AppInstallerMetadata.Description)
	data.CommunityURL = types.StringValue(appResp.AppInstallerMetadata.CommunityURL)
	data.GithubURL = types.StringValue(appResp.AppInstallerMetadata.GithubURL)
	data.DocumentationURL = types.StringValue(appResp.AppInstallerMetadata.DocumentationURL)
	data.HomepageURL = types.StringValue(appResp.AppInstallerMetadata.HomepageURL)
	data.LogoURL = types.StringValue(appResp.AppInstallerMetadata.LogoURL)
	if appResp.AppInstallerMetadata.DemoURL != "" {
		data.DemoURL = types.StringValue(appResp.AppInstallerMetadata.DemoURL)
	}

	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully read app installer")
}

func (r *AppInstallerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// get terraform model
	var data *AppInstallerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updating app installer")

	// update app
	_, err := r.restClient.UpdateAppInstaller(ctx, data.Id.ValueString(), &models.ServiceUpdateAppInstallerRequest{
		Name:        data.Name.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Links: &models.ServiceUpdateAppInstallerRequestLinks{
			Community:     data.CommunityURL.ValueStringPointer(),
			Documentation: data.DocumentationURL.ValueStringPointer(),
			Homepage:      data.HomepageURL.ValueStringPointer(),
			Github:        data.GithubURL.ValueStringPointer(),
			Logo:          data.LogoURL.ValueStringPointer(),
			Demo:          data.DemoURL.ValueStringPointer(),
		},
	})
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "update app installer")
		return
	}

	appResp, err := r.restClient.GetAppInstaller(ctx, data.Id.ValueString())
	if nuon.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "read app installer")
		return
	}

	// update the state with the returned values
	data.Id = types.StringValue(appResp.ID)
	data.Slug = types.StringValue(appResp.Slug)
	data.Name = types.StringValue(appResp.AppInstallerMetadata.Name)
	data.Description = types.StringValue(appResp.AppInstallerMetadata.Description)
	data.CommunityURL = types.StringValue(appResp.AppInstallerMetadata.CommunityURL)
	data.GithubURL = types.StringValue(appResp.AppInstallerMetadata.GithubURL)
	data.DocumentationURL = types.StringValue(appResp.AppInstallerMetadata.DocumentationURL)
	data.HomepageURL = types.StringValue(appResp.AppInstallerMetadata.HomepageURL)
	data.LogoURL = types.StringValue(appResp.AppInstallerMetadata.LogoURL)
	if appResp.AppInstallerMetadata.DemoURL != "" {
		data.DemoURL = types.StringValue(appResp.AppInstallerMetadata.DemoURL)
	}

	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully updated app")
}

func (r *AppInstallerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AppInstallerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "deleting app installer")

	deleted, err := r.restClient.DeleteAppInstaller(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete app installer",
			fmt.Sprintf("Please make sure your app_id (%s) is correct, and that the auth token has permissions for this org. ", data.Id.String()),
		)
		return
	}
	if !deleted {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "delete app installer")
		return
	}
}

func (r *AppInstallerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	// resource.ImportStatePassthroughID(ctx, path.Root("org_id"), req, resp)
}
