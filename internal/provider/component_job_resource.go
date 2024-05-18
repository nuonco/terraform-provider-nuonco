package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/nuonco/nuon-go"
	"github.com/nuonco/nuon-go/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &JobComponentResource{}
var _ resource.ResourceWithImportState = &JobComponentResource{}

func NewJobComponentResource() resource.Resource {
	return &JobComponentResource{}
}

// JobComponentResource defines the resource implementation.
type JobComponentResource struct {
	baseResource
}

// JobComponentResourceModel describes the resource data model.
type JobComponentResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name         types.String `tfsdk:"name"`
	VarName      types.String `tfsdk:"var_name"`
	Dependencies types.List   `tfsdk:"dependencies"`
	AppID        types.String `tfsdk:"app_id"`
	ImageURL     types.String `tfsdk:"image_url"`
	Tag          types.String `tfsdk:"tag"`
	Cmd          types.List   `tfsdk:"cmd"`
	Args         types.List   `tfsdk:"args"`
	EnvVar       EnvVarSlice  `tfsdk:"env_var"`
}

func (r *JobComponentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job_component"
}

func (r *JobComponentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Release a container as a k8s job.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique ID of the component.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The human-readable name of the component.",
				Optional:    false,
				Required:    true,
			},
			"var_name": schema.StringAttribute{
				Description: "The optional var name to be used when referencing this component.",
				Optional:    true,
				Required:    false,
			},
			"app_id": schema.StringAttribute{
				Description: "The unique ID of the app this component belongs too.",
				Optional:    false,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dependencies": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Component dependencies",
				Optional:    true,
				Required:    false,
			},
			"image_url": schema.StringAttribute{
				Description: "The full image URL or docker hub alias (e.g. kennethreitz/httpbin).",
				Required:    true,
			},
			"tag": schema.StringAttribute{
				Description: "The image tag.",
				Required:    true,
			},
			"cmd": schema.ListAttribute{
				Description: "The command to execute.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"args": schema.ListAttribute{
				Description: "Arguments to pass to the command.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
		Blocks: map[string]schema.Block{
			"env_var": envVarSharedBlock(),
		},
	}
}

func (r *JobComponentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *JobComponentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "creating component")

	dependencies := make([]string, 0)
	resp.Diagnostics.Append(data.Dependencies.ElementsAs(ctx, &dependencies, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	compResp, err := r.restClient.CreateComponent(ctx, data.AppID.ValueString(), &models.ServiceCreateComponentRequest{
		Name:         data.Name.ValueStringPointer(),
		VarName:      data.VarName.ValueString(),
		Dependencies: dependencies,
	})
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create component")
		return
	}
	tflog.Trace(ctx, "got ID -- "+compResp.ID)
	data.ID = types.StringValue(compResp.ID)
	data.VarName = types.StringValue(compResp.VarName)

	configRequest := &models.ServiceCreateJobComponentConfigRequest{
		ImageURL: data.ImageURL.ValueStringPointer(),
		Tag:      data.Tag.ValueStringPointer(),
		Cmd:      listToStringSlice(data.Cmd),
		Args:     listToStringSlice(data.Args),
		EnvVars:  data.EnvVar.ToMap(),
	}
	_, err = r.restClient.CreateJobComponentConfig(ctx, compResp.ID, configRequest)
	if err != nil {
		// attempt to cleanup component, that is in broken state and has no config
		_, cleanupErr := r.restClient.DeleteComponent(ctx, compResp.ID)
		if cleanupErr != nil {
			tflog.Trace(ctx, fmt.Sprintf("unable to cleanup component: %s", cleanupErr))
		}

		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create component config")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully created component")
}

func (r *JobComponentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *JobComponentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	compResp, err := r.restClient.GetComponent(ctx, data.ID.ValueString())
	if nuon.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "get component")
		return
	}
	data.Name = types.StringValue(compResp.Name)
	data.VarName = types.StringValue(compResp.VarName)
	data.AppID = types.StringValue(compResp.AppID)

	configResp, err := r.restClient.GetComponentLatestConfig(ctx, data.ID.ValueString())
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "get component config")
		return
	}
	if configResp.Job == nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, errors.New("did not get job config"), "get component config")
		return
	}
	data.ImageURL = types.StringValue(configResp.Job.ImageURL)
	data.Tag = types.StringValue(configResp.Job.Tag)
	data.Cmd = stringSliceToList(ctx, configResp.Job.Cmd)
	data.Args = stringSliceToList(ctx, configResp.Job.Args)
	data.EnvVar = NewEnvVarSliceFromMap(configResp.Job.EnvVars)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully read component")
}

func (r *JobComponentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *JobComponentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleted, err := r.restClient.DeleteComponent(ctx, data.ID.ValueString())
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "delete component")
		return
	}
	if !deleted {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "delete component")
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{statusActive, statusDeleteQueued, statusDeprovisioning, statusTemporarilyUnavailable},
		Target:  []string{statusNotFound},
		Refresh: func() (interface{}, string, error) {
			tflog.Trace(ctx, "refreshing component status")
			cmp, err := r.restClient.GetComponent(ctx, data.ID.ValueString())
			if err == nil {
				return cmp.Status, cmp.Status, nil
			}
			if nuon.IsNotFound(err) {
				return "", statusNotFound, nil
			}

			logErr(ctx, err, "delete component")
			return statusTemporarilyUnavailable, statusTemporarilyUnavailable, nil
		},
		Timeout:    time.Minute * 20,
		Delay:      time.Second * 10,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "unable to delete component")
		return
	}

	tflog.Trace(ctx, "successfully deleted component")
}

func (r *JobComponentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *JobComponentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updating component "+data.ID.ValueString())

	dependencies := make([]string, 0)
	resp.Diagnostics.Append(data.Dependencies.ElementsAs(ctx, &dependencies, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	compResp, err := r.restClient.UpdateComponent(ctx, data.ID.ValueString(), &models.ServiceUpdateComponentRequest{
		Name:         data.Name.ValueStringPointer(),
		VarName:      data.VarName.ValueString(),
		Dependencies: dependencies,
	})
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "update component")
		return
	}

	data.Name = types.StringValue(compResp.Name)
	data.VarName = types.StringValue(compResp.VarName)

	configRequest := &models.ServiceCreateJobComponentConfigRequest{
		ImageURL: data.ImageURL.ValueStringPointer(),
		Tag:      data.Tag.ValueStringPointer(),
		Cmd:      listToStringSlice(data.Cmd),
		Args:     listToStringSlice(data.Args),
		EnvVars:  data.EnvVar.ToMap(),
	}
	_, err = r.restClient.CreateJobComponentConfig(ctx, compResp.ID, configRequest)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create component config")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully updated component")
}

func (r *JobComponentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
