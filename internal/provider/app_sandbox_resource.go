package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nuonco/nuon-go"
	"github.com/nuonco/nuon-go/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AppSandboxResource{}
var _ resource.ResourceWithImportState = &AppSandboxResource{}

func NewAppSandboxResource() resource.Resource {
	return &AppSandboxResource{}
}

// AppSandboxResource defines the resource implementation.
type AppSandboxResource struct {
	baseResource
}

// AppSandboxResourceModel describes the resource data model.
type AppSandboxResourceModel struct {
	ID    types.String `tfsdk:"id"`
	AppID types.String `tfsdk:"app_id"`

	// one of the following sources must be set for the app sandbox
	BuiltinSandboxReleaseID types.String   `tfsdk:"builtin_sandbox_release_id"`
	PublicRepo              *PublicRepo    `tfsdk:"public_repo"`
	ConnectedRepo           *ConnectedRepo `tfsdk:"connected_repo"`

	Inputs           []SandboxInput `tfsdk:"input"`
	TerraformVersion types.String   `tfsdk:"terraform_version"`
}

type SandboxInput struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (r *AppSandboxResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_sandbox"
}

func (r *AppSandboxResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Sandbox configuration for an app.",
		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				Description:   "The application ID.",
				Optional:      false,
				Required:      true,
				PlanModifiers: []planmodifier.String{},
			},
			"id": schema.StringAttribute{
				Description:   "The sandbox config id",
				Computed:      true,
				PlanModifiers: []planmodifier.String{},
			},
			"builtin_sandbox_release_id": schema.StringAttribute{
				Description:   "release ID for a built in sandbox to use",
				Optional:      true,
				Required:      false,
				PlanModifiers: []planmodifier.String{},
			},
			"terraform_version": schema.StringAttribute{
				Description:   "terraform version to use with sandbox",
				Optional:      false,
				Required:      true,
				PlanModifiers: []planmodifier.String{},
			},
			"public_repo":    publicRepoAttribute(),
			"connected_repo": connectedRepoAttribute(),
		},
		Blocks: map[string]schema.Block{
			"input": schema.SetNestedBlock{
				Description: "default sandbox inputs that will be used on each install. Can use Nuon interpolation language.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The input name to be used, which will be used as a terraform variable input to the sandbox.",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The static value, or interpolate value to set.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *AppSandboxResource) getConfigRequest(data *AppSandboxResourceModel) (*models.ServiceCreateAppSandboxConfigRequest, error) {
	if data.ConnectedRepo == nil && data.PublicRepo == nil && data.BuiltinSandboxReleaseID.ValueString() == "" {
		return nil, fmt.Errorf("must set one of connected_repo, public_repo or builtin_sandbox_release_id")
	}

	cfgReq := &models.ServiceCreateAppSandboxConfigRequest{
		SandboxInputs: make(map[string]string),
	}

	// configure source
	if data.PublicRepo != nil {
		cfgReq.PublicGitVcsConfig = &models.ServicePublicGitVCSSandboxConfigRequest{
			Branch:    data.PublicRepo.Branch.ValueStringPointer(),
			Directory: data.PublicRepo.Directory.ValueStringPointer(),
			Repo:      data.PublicRepo.Repo.ValueStringPointer(),
		}
	}
	if data.ConnectedRepo != nil {
		cfgReq.ConnectedGithubVcsConfig = &models.ServiceConnectedGithubVCSSandboxConfigRequest{
			Branch:    data.ConnectedRepo.Branch.ValueString(),
			Directory: data.ConnectedRepo.Directory.ValueStringPointer(),
			Repo:      data.ConnectedRepo.Repo.ValueStringPointer(),
		}
	}
	if data.BuiltinSandboxReleaseID.ValueString() != "" {
		cfgReq.SandboxReleaseID = data.BuiltinSandboxReleaseID.ValueString()
	}

	// configure inputs
	for _, input := range data.Inputs {
		cfgReq.SandboxInputs[input.Name.ValueString()] = input.Value.ValueString()
	}
	cfgReq.TerraformVersion = toPtr(data.TerraformVersion.ValueString())

	return cfgReq, nil
}

func (r *AppSandboxResource) writeStateData(data *AppSandboxResourceModel, resp *models.AppAppSandboxConfig) {
	data.ID = types.StringValue(resp.ID)
	if resp.ConnectedGithubVcsConfig != nil {
		connected := resp.ConnectedGithubVcsConfig
		data.ConnectedRepo = &ConnectedRepo{
			Branch:    types.StringValue(connected.Branch),
			Directory: types.StringValue(connected.Directory),
			Repo:      types.StringValue(connected.Repo),
		}
	}

	if resp.PublicGitVcsConfig != nil {
		public := resp.PublicGitVcsConfig
		data.PublicRepo = &PublicRepo{
			Branch:    types.StringValue(public.Branch),
			Directory: types.StringValue(public.Directory),
			Repo:      types.StringValue(public.Repo),
		}
	}
	if resp.SandboxReleaseID != "" {
		data.BuiltinSandboxReleaseID = types.StringValue(resp.SandboxReleaseID)
	}

	inputs := []SandboxInput{}
	for key, val := range resp.SandboxInputs {
		inputs = append(inputs, SandboxInput{
			Name:  types.StringValue(key),
			Value: types.StringValue(val),
		})
	}
	data.Inputs = inputs
	data.TerraformVersion = types.StringValue(resp.TerraformVersion)
}

func (r *AppSandboxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// get terraform model
	var data *AppSandboxResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create app
	tflog.Trace(ctx, "creating app sandbox")
	cfgReq, err := r.getConfigRequest(data)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create app sandbox")
		return
	}

	appResp, err := r.restClient.CreateAppSandboxConfig(ctx, data.AppID.ValueString(), cfgReq)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create app sandbox")
		return
	}

	r.writeStateData(data, appResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppSandboxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AppSandboxResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading app sandbox")
	appResp, err := r.restClient.GetAppSandboxLatestConfig(ctx, data.AppID.ValueString())
	if nuon.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "read app sandbox")
		return
	}

	r.writeStateData(data, appResp)
	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully read app sandbox")
}

func (r *AppSandboxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// get terraform model
	var data *AppSandboxResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updating app installer")

	// update app
	cfgReq, err := r.getConfigRequest(data)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create app sandbox")
		return
	}

	cfgResp, err := r.restClient.CreateAppSandboxConfig(ctx, data.AppID.ValueString(), cfgReq)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "update app sandbox")
		return
	}
	if nuon.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "read app sandbox")
		return
	}

	r.writeStateData(data, cfgResp)
	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully updated app sandbox")
}

func (r *AppSandboxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *AppSandboxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_id"), req, resp)
}
