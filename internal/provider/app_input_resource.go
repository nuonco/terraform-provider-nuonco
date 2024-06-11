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
var (
	_ resource.Resource                = &AppInputResource{}
	_ resource.ResourceWithImportState = &AppInputResource{}
)

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

	Inputs []AppInput      `tfsdk:"input"`
	Groups []AppInputGroup `tfsdk:"group"`
}

type AppInput struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	Group       types.String `tfsdk:"group"`
	Required    types.Bool   `tfsdk:"required"`
	Default     types.String `tfsdk:"default"`
	Sensitive   types.Bool   `tfsdk:"sensitive"`
}

type AppInputGroup struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
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
			"group": schema.SetNestedBlock{
				Description: "Input group, which can be used to organize sets of inputs.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The group name to be used.",
							Required:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "Human readable display name.",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of input group.",
							Required:    true,
						},
					},
				},
			},
			"input": schema.SetNestedBlock{
				Description: "Required inputs that each install must provide for this app.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The input name to be used, which will be used to expose this in the interpolation language, using `.nuon.install.inputs.input_name`",
							Required:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "Human readable display name.",
							Required:    true,
						},
						"default": schema.StringAttribute{
							Description: "Default value for input",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of input.",
							Required:    true,
						},
						"group": schema.StringAttribute{
							Description: "Add to a specific group",
							Required:    true,
						},
						"required": schema.BoolAttribute{
							Description: "Mark whether this input is required or not.",
							Optional:    true,
						},
						"sensitive": schema.BoolAttribute{
							Description: "Mark whether the input is required or not",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *AppInputResource) getConfigRequest(data *AppInputResourceModel) (*models.ServiceCreateAppInputConfigRequest, error) {
	cfgReq := &models.ServiceCreateAppInputConfigRequest{
		Inputs: make(map[string]models.ServiceAppInputRequest),
		Groups: make(map[string]models.ServiceAppGroupRequest),
	}

	// configure inputs
	for _, input := range data.Inputs {
		cfgReq.Inputs[input.Name.ValueString()] = models.ServiceAppInputRequest{
			Default:     input.Default.ValueString(),
			DisplayName: toPtr(input.DisplayName.ValueString()),
			Description: toPtr(input.Description.ValueString()),
			Required:    input.Required.ValueBool(),
			Sensitive:   input.Sensitive.ValueBool(),
			Group:       toPtr(input.Group.ValueString()),
		}
	}

	for _, grp := range data.Groups {
		cfgReq.Groups[grp.Name.ValueString()] = models.ServiceAppGroupRequest{
			DisplayName: toPtr(grp.DisplayName.ValueString()),
			Description: toPtr(grp.Description.ValueString()),
		}
	}

	return cfgReq, nil
}

func (r *AppInputResource) writeStateData(data *AppInputResourceModel, resp *models.AppAppInputConfig) {
	data.ID = types.StringValue(resp.ID)

	return
	inputs := []AppInput{}
	for _, inp := range resp.Inputs {

		inpData := AppInput{
			Name:        types.StringValue(inp.Name),
			Description: types.StringValue(inp.Description),
			DisplayName: types.StringValue(inp.DisplayName),
			Default:     types.StringValue(inp.Default),
			Required:    types.BoolValue(inp.Required),
			Sensitive:   types.BoolValue(inp.Sensitive),
		}
		if inp.Group.Name != "" {
			inpData.Group = types.StringValue(inp.Group.Name)
		}
		inputs = append(inputs, inpData)
	}
	data.Inputs = inputs

	groups := []AppInputGroup{}
	for _, grp := range resp.InputGroups {
		groups = append(groups, AppInputGroup{
			Name:        types.StringValue(grp.Name),
			Description: types.StringValue(grp.Description),
			DisplayName: types.StringValue(grp.DisplayName),
		})
	}
	data.Groups = groups
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
