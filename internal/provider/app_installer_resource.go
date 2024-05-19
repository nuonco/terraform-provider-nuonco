package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	PostInstallMarkdown types.String `tfsdk:"post_install_markdown"`

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
			"post_install_markdown": schema.StringAttribute{
				MarkdownDescription: "Markdown that will be shown to users after a successful install. Supports interpolation.",
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
}

func (r *AppInstallerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *AppInstallerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *AppInstallerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *AppInstallerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	// resource.ImportStatePassthroughID(ctx, path.Root("org_id"), req, resp)
}
