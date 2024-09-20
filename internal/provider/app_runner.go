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
var _ resource.Resource = &AppRunnerResource{}
var _ resource.ResourceWithImportState = &AppRunnerResource{}

func NewAppRunnerResource() resource.Resource {
	return &AppRunnerResource{}
}

// AppRunnerResource defines the resource implementation.
type AppRunnerResource struct {
	baseResource
}

// AppRunnerResourceModel describes the resource data model.
type AppRunnerResourceModel struct {
	ID    types.String `tfsdk:"id"`
	AppID types.String `tfsdk:"app_id"`

	EnvVar     []EnvVar     `tfsdk:"env_var"`
	RunnerType types.String `tfsdk:"runner_type"`
}

func (r *AppRunnerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_runner"
}

func (r *AppRunnerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "runner configuration for an app.",
		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				Description: "The application ID.",
				Optional:    false,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description:   "The runner config id",
				Computed:      true,
				PlanModifiers: []planmodifier.String{},
			},
			"runner_type": schema.StringAttribute{
				Description:   "runner type",
				Optional:      false,
				Required:      true,
				PlanModifiers: []planmodifier.String{},
			},
		},
		Blocks: map[string]schema.Block{
			"env_var": envVarSharedBlock(),
		},
	}
}

func (r *AppRunnerResource) getConfigRequest(data *AppRunnerResourceModel) (*models.ServiceCreateAppRunnerConfigRequest, error) {
	cfgReq := &models.ServiceCreateAppRunnerConfigRequest{
		EnvVars: make(map[string]string),
	}

	switch data.RunnerType.ValueString() {
	case "aws-eks":
		cfgReq.Type = models.NewAppAppRunnerType(models.AppAppRunnerTypeAwsDashEks)
	case "aws-ecs":
		cfgReq.Type = models.NewAppAppRunnerType(models.AppAppRunnerTypeAwsDashEcs)
	case "azure-aks":
		cfgReq.Type = models.NewAppAppRunnerType(models.AppAppRunnerTypeAzureDashAks)
	case "azure-acs":
		cfgReq.Type = models.NewAppAppRunnerType(models.AppAppRunnerTypeAzureDashAcs)
	default:
		return nil, fmt.Errorf("invalid runner-type, must be one of (aws-eks, aws-ecs, azure-aks, azure-acs)")
	}

	// configure inputs
	for _, input := range data.EnvVar {
		cfgReq.EnvVars[input.Name.ValueString()] = input.Value.ValueString()
	}

	return cfgReq, nil
}

func (r *AppRunnerResource) writeStateData(data *AppRunnerResourceModel, resp *models.AppAppRunnerConfig) {
	data.ID = types.StringValue(resp.ID)

	envVars := []EnvVar{}
	for key, val := range resp.EnvVars {
		envVars = append(envVars, EnvVar{
			Name:  types.StringValue(key),
			Value: types.StringValue(val),
		})
	}
	data.EnvVar = envVars
	data.RunnerType = types.StringValue(string(resp.AppRunnerType))
}

func (r *AppRunnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// get terraform model
	var data *AppRunnerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create app
	tflog.Trace(ctx, "creating app runner")
	cfgReq, err := r.getConfigRequest(data)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create app runner")
		return
	}

	appResp, err := r.restClient.CreateAppRunnerConfig(ctx, data.AppID.ValueString(), cfgReq)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create app runner")
		return
	}

	r.writeStateData(data, appResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppRunnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AppRunnerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading app runner")
	appResp, err := r.restClient.GetAppRunnerLatestConfig(ctx, data.AppID.ValueString())
	if nuon.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "read app runner")
		return
	}

	r.writeStateData(data, appResp)
	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully read app runner")
}

func (r *AppRunnerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// get terraform model
	var data *AppRunnerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updating app installer")

	// update app
	cfgReq, err := r.getConfigRequest(data)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create app runner")
		return
	}

	cfgResp, err := r.restClient.CreateAppRunnerConfig(ctx, data.AppID.ValueString(), cfgReq)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "update app runner")
		return
	}
	if nuon.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "read app runner")
		return
	}

	r.writeStateData(data, cfgResp)

	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully updated app runner")
}

func (r *AppRunnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *AppRunnerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_id"), req, resp)
}
