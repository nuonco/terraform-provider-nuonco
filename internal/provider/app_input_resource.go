package provider

import (
	"context"

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
var _ resource.Resource = &AppInputResource{}
var _ resource.ResourceWithImportState = &AppInputResource{}

func NewAppInputResource() resource.Resource {
	return &AppInputResource{}
}

// AppInputResource defines the resource implementation.
type AppInputResource struct {
	baseResource
}

// AppInputResourceModel describes the resource data model.
type AppInputResourceModel struct {
	ID    types.String `tfsdk:"id"`
	AppID types.String `tfsdk:"app_id"`

	Inputs []AppInput `tfsdk:"input"`
}

type AppInput struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (r *AppInputResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_input"
}

func (r *AppInputResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Input configuration for an app.",
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
				Description:   "The app input config id",
				Computed:      true,
				PlanModifiers: []planmodifier.String{},
			},
		},
		Blocks: map[string]schema.Block{
			"input": schema.SetNestedBlock{
				Description: "Required inputs that each install must provide for this app.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The input name to be used, which will be used to expose this in the interpolation language, using {{.nuon.install.inputs.<name>}}",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value to set.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *AppInputResource) getConfigRequest(data *AppInputResourceModel) (*models.ServiceCreateAppInputConfigRequest, error) {
	cfgReq := &models.ServiceCreateAppInputConfigRequest{
		Inputs: make(map[string]string),
	}

	// configure inputs
	for _, input := range data.Inputs {
		cfgReq.Inputs[input.Name.ValueString()] = input.Value.ValueString()
	}

	return cfgReq, nil
}

func (r *AppInputResource) writeStateData(data *AppInputResourceModel, resp *models.AppAppInputConfig) {
	data.ID = types.StringValue(resp.ID)
	inputs := []AppInput{}
	for key, val := range resp.Inputs {
		inputs = append(inputs, AppInput{
			Name:  types.StringValue(key),
			Value: types.StringValue(val),
		})
	}
	data.Inputs = inputs
}

func (r *AppInputResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// get terraform model
	var data *AppInputResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create app
	tflog.Trace(ctx, "creating app input config")
	cfgReq, err := r.getConfigRequest(data)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create app input")
		return
	}

	appResp, err := r.restClient.CreateAppInputConfig(ctx, data.AppID.ValueString(), cfgReq)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create app input")
		return
	}

	r.writeStateData(data, appResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AppInputResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AppInputResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading app input config")
	appResp, err := r.restClient.GetAppInputLatestConfig(ctx, data.AppID.ValueString())
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
	tflog.Trace(ctx, "successfully read app input config")
}

func (r *AppInputResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// get terraform model
	var data *AppInputResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updating app input config")

	// update app
	cfgReq, err := r.getConfigRequest(data)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create app input config")
		return
	}

	cfgResp, err := r.restClient.CreateAppInputConfig(ctx, data.AppID.ValueString(), cfgReq)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "update app input config")
		return
	}
	if nuon.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "read app input config")
		return
	}

	r.writeStateData(data, cfgResp)
	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully updated app input config")
}

func (r *AppInputResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *AppInputResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_id"), req, resp)
}
