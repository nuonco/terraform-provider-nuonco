package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nuonco/nuon-go"
	"github.com/nuonco/nuon-go/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &InstallerResource{}
	_ resource.ResourceWithImportState = &InstallerResource{}
)

func NewInstallerResource() resource.Resource {
	return &InstallerResource{}
}

// InstallerResource defines the resource implementation.
type InstallerResource struct {
	baseResource
}

// InstallerResourceModel describes the resource data model.
type InstallerResourceModel struct {
	Id types.String `tfsdk:"id"`

	AppIDs types.Set `tfsdk:"app_ids"`

	// metadata
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	PostInstallMarkdown types.String `tfsdk:"post_install_markdown"`
	FooterMarkdown      types.String `tfsdk:"footer_markdown"`
	CopyrightMarkdown   types.String `tfsdk:"copyright_markdown"`
	DemoURL             types.String `tfsdk:"demo_url"`

	DocumentationURL types.String `tfsdk:"documentation_url"`
	HomepageURL      types.String `tfsdk:"homepage_url"`
	CommunityURL     types.String `tfsdk:"community_url"`
	GithubURL        types.String `tfsdk:"github_url"`
	LogoURL          types.String `tfsdk:"logo_url"`
	FaviconURL       types.String `tfsdk:"favicon_url"`
}

func (r *InstallerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_installer"
}

func (r *InstallerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Nuon installer for one or more apps.",
		Attributes: map[string]schema.Attribute{
			"app_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "App IDs to connect to this installer.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
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
			"post_install_markdown": schema.StringAttribute{
				MarkdownDescription: "Markdown that will be shown to users after a successful install. Supports interpolation.",
				Optional:            true,
				Required:            false,
			},
			"copyright_markdown": schema.StringAttribute{
				MarkdownDescription: "Markdown that rendered in the copyright section.",
				Optional:            true,
				Required:            false,
			},
			"footer_markdown": schema.StringAttribute{
				MarkdownDescription: "Markdown that will be rendered in the footer section.",
				Optional:            true,
				Required:            false,
			},
			"documentation_url": schema.StringAttribute{
				MarkdownDescription: "Documentation url",
				Optional:            false,
				Required:            true,
			},
			"favicon_url": schema.StringAttribute{
				MarkdownDescription: "Favicon url",
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
				Optional:            true,
				Required:            false,
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

func (r *InstallerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// get terraform model
	var data *InstallerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create app
	tflog.Trace(ctx, "creating app installer")
	appIDs := make([]string, 0)
	diags := data.AppIDs.ElementsAs(ctx, &appIDs, false)
	if diags.HasError() {
		return
	}

	appResp, err := r.restClient.CreateInstaller(ctx, &models.ServiceCreateInstallerRequest{
		AppIds: appIDs,
		Name:   data.Name.ValueStringPointer(),
		Metadata: &models.ServiceCreateInstallerRequestMetadata{
			Description: data.Description.ValueStringPointer(),

			CommunityURL:     data.CommunityURL.ValueStringPointer(),
			FaviconURL:       data.FaviconURL.ValueString(),
			DocumentationURL: data.DocumentationURL.ValueStringPointer(),
			HomepageURL:      data.HomepageURL.ValueStringPointer(),
			GithubURL:        data.GithubURL.ValueStringPointer(),
			LogoURL:          data.LogoURL.ValueStringPointer(),

			DemoURL:             data.DemoURL.ValueString(),
			PostInstallMarkdown: data.PostInstallMarkdown.ValueString(),
			FooterMarkdown:      data.FooterMarkdown.ValueString(),
			CopyrightMarkdown:   data.CopyrightMarkdown.ValueString(),
		},
	})
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create app installer")
		return
	}

	// update the state with the returned values
	data.Id = types.StringValue(appResp.ID)
	data.Name = types.StringValue(appResp.Metadata.Name)
	data.Description = types.StringValue(appResp.Metadata.Description)
	data.CommunityURL = types.StringValue(appResp.Metadata.CommunityURL)
	data.GithubURL = types.StringValue(appResp.Metadata.GithubURL)
	data.DocumentationURL = types.StringValue(appResp.Metadata.DocumentationURL)
	data.HomepageURL = types.StringValue(appResp.Metadata.HomepageURL)
	data.LogoURL = types.StringValue(appResp.Metadata.LogoURL)
	data.FaviconURL = types.StringValue(appResp.Metadata.FaviconURL)
	data.PostInstallMarkdown = types.StringValue(appResp.Metadata.PostInstallMarkdown)
	data.FooterMarkdown = types.StringValue(appResp.Metadata.FooterMarkdown)
	data.CopyrightMarkdown = types.StringValue(appResp.Metadata.CopyrightMarkdown)
	data.DemoURL = types.StringValue(appResp.Metadata.DemoURL)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InstallerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *InstallerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading app installer")

	appResp, err := r.restClient.GetInstaller(ctx, data.Id.ValueString())
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
	data.Name = types.StringValue(appResp.Metadata.Name)
	data.Description = types.StringValue(appResp.Metadata.Description)
	data.CommunityURL = types.StringValue(appResp.Metadata.CommunityURL)
	data.GithubURL = types.StringValue(appResp.Metadata.GithubURL)
	data.DocumentationURL = types.StringValue(appResp.Metadata.DocumentationURL)
	data.HomepageURL = types.StringValue(appResp.Metadata.HomepageURL)
	data.LogoURL = types.StringValue(appResp.Metadata.LogoURL)
	data.FaviconURL = types.StringValue(appResp.Metadata.FaviconURL)
	data.PostInstallMarkdown = types.StringValue(appResp.Metadata.PostInstallMarkdown)
	data.FooterMarkdown = types.StringValue(appResp.Metadata.FooterMarkdown)
	data.CopyrightMarkdown = types.StringValue(appResp.Metadata.CopyrightMarkdown)
	data.DemoURL = types.StringValue(appResp.Metadata.DemoURL)

	appIDItems := []attr.Value{}
	for _, app := range appResp.Apps {
		appIDItems = append(appIDItems, types.StringValue(app.ID))
	}
	appIDsResp, diags := types.SetValue(types.StringType, appIDItems)
	if diags.HasError() {
		return
	}
	data.AppIDs = appIDsResp

	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully read app installer")
}

func (r *InstallerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// get terraform model
	var data *InstallerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updating app installer")
	appIDs := make([]string, 0)
	diags := data.AppIDs.ElementsAs(ctx, &appIDs, false)
	if diags.HasError() {
		return
	}

	// update app
	_, err := r.restClient.UpdateInstaller(ctx, data.Id.ValueString(), &models.ServiceUpdateInstallerRequest{
		Name:   data.Name.ValueStringPointer(),
		AppIds: appIDs,
		Metadata: &models.ServiceUpdateInstallerRequestMetadata{
			Description:         data.Description.ValueStringPointer(),
			PostInstallMarkdown: data.PostInstallMarkdown.ValueString(),
			CopyrightMarkdown:   data.CopyrightMarkdown.ValueString(),
			FooterMarkdown:      data.FooterMarkdown.ValueString(),

			CommunityURL:     data.CommunityURL.ValueStringPointer(),
			FaviconURL:       data.FaviconURL.ValueString(),
			DocumentationURL: data.DocumentationURL.ValueStringPointer(),
			HomepageURL:      data.HomepageURL.ValueStringPointer(),
			GithubURL:        data.GithubURL.ValueStringPointer(),
			LogoURL:          data.LogoURL.ValueStringPointer(),
			DemoURL:          data.DemoURL.ValueString(),
		},
	})
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "update app installer")
		return
	}

	appResp, err := r.restClient.GetInstaller(ctx, data.Id.ValueString())
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
	data.Name = types.StringValue(appResp.Metadata.Name)
	data.Description = types.StringValue(appResp.Metadata.Description)
	data.CommunityURL = types.StringValue(appResp.Metadata.CommunityURL)
	data.GithubURL = types.StringValue(appResp.Metadata.GithubURL)
	data.DocumentationURL = types.StringValue(appResp.Metadata.DocumentationURL)
	data.HomepageURL = types.StringValue(appResp.Metadata.HomepageURL)
	data.LogoURL = types.StringValue(appResp.Metadata.LogoURL)
	data.FaviconURL = types.StringValue(appResp.Metadata.FaviconURL)
	data.PostInstallMarkdown = types.StringValue(appResp.Metadata.PostInstallMarkdown)
	data.FooterMarkdown = types.StringValue(appResp.Metadata.FooterMarkdown)
	data.CopyrightMarkdown = types.StringValue(appResp.Metadata.CopyrightMarkdown)
	data.DemoURL = types.StringValue(appResp.Metadata.DemoURL)

	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully updated app")
}

func (r *InstallerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *InstallerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "deleting app installer")

	deleted, err := r.restClient.DeleteInstaller(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete installer",
			fmt.Sprintf("Please make sure your app_id (%s) is correct, and that the auth token has permissions for this org. ", data.Id.String()),
		)
		return
	}
	if !deleted {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "delete app installer")
		return
	}
}

func (r *InstallerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	// resource.ImportStatePassthroughID(ctx, path.Root("org_id"), req, resp)
}
