package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/nuonco/nuon-go/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AppResource{}
var _ resource.ResourceWithImportState = &AppResource{}

func NewAppResource() resource.Resource {
	return &AppResource{}
}

// AppResource defines the resource implementation.
type AppResource struct {
	baseResource
}

// AppResourceModel describes the resource data model.
type AppResourceModel struct {
	Name           types.String          `tfsdk:"name"`
	Id             types.String          `tfsdk:"id"`
	SandboxRelease basetypes.ObjectValue `tfsdk:"sandbox_release"`
}

func (r *AppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (r *AppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Nuon application, required to set up components and installs.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The human readable name of the app.",
				Optional:            false,
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique ID of the app.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sandbox_release": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The sandbox being used for this app's installs.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"version": schema.StringAttribute{
						Computed: true,
					},
					"terraform_version": schema.StringAttribute{
						Computed: true,
					},
					"provision_policy_url": schema.StringAttribute{
						Computed: true,
					},
					"deprovision_policy_url": schema.StringAttribute{
						Computed: true,
					},
					"trust_policy_url": schema.StringAttribute{
						Computed: true,
					},
					"one_click_role_template_url": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

func (r *AppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// get terraform model
	var data *AppResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create app
	tflog.Trace(ctx, "creating app")
	appResp, err := r.restClient.CreateApp(ctx, &models.ServiceCreateAppRequest{
		Name: data.Name.ValueStringPointer(),
	})
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create app")
		return
	}

	// populate terraform model with data from api
	data.Name = types.StringValue(appResp.Name)
	data.Id = types.StringValue(appResp.ID)
	data.SandboxRelease = convertSandboxRelease(*appResp.SandboxRelease)

	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully created app")

	// poll app to completion status
	stateConf := &retry.StateChangeConf{
		Pending: []string{statusQueued, statusProvisioning},
		Target:  []string{statusActive},
		Refresh: func() (interface{}, string, error) {
			tflog.Trace(ctx, "refreshing app status")
			app, err := r.restClient.GetApp(ctx, appResp.ID)
			if err != nil {
				writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "poll status")
				return nil, "unknown", err
			}
			return app.Status, string(app.Status), nil
		},
		Timeout:    time.Minute * 20,
		Delay:      time.Second * 10,
		MinTimeout: 3 * time.Second,
	}
	statusRaw, err := stateConf.WaitForState()
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "get app")
		return
	}
	status, ok := statusRaw.(string)
	if !ok {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, fmt.Errorf("invalid app %s status", status), "create app")
		return
	}
	if status != statusActive {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, fmt.Errorf("status %s", status), "create app")
		return
	}
}

func (r *AppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// get terraform model
	var data *AppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "reading app")

	// get app from api
	appResp, err := r.restClient.GetApp(ctx, data.Id.ValueString())
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "read app")
		return
	}

	// populate terraform model with data from api
	data.Name = types.StringValue(appResp.Name)
	data.Id = types.StringValue(appResp.ID)
	data.SandboxRelease = convertSandboxRelease(*appResp.SandboxRelease)

	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully read app")
}

func (r *AppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// get terraform model
	var data *AppResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "updating app")

	// update app
	_, err := r.restClient.UpdateApp(ctx, data.Id.ValueString(), &models.ServiceUpdateAppRequest{
		Name: data.Name.ValueString(),
	})
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "update app")
		return
	}

	// get app from api
	appResp, err := r.restClient.GetApp(ctx, data.Id.ValueString())
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "update app")
		return
	}

	// populate terraform model with data from api
	data.Name = types.StringValue(appResp.Name)
	data.Id = types.StringValue(appResp.ID)
	data.SandboxRelease = convertSandboxRelease(*appResp.SandboxRelease)

	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully updated app")
}

func (r *AppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "deleting app")

	deleted, err := r.restClient.DeleteApp(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete app",
			fmt.Sprintf("Please make sure your app_id (%s) is correct, and that the auth token has permissions for this org. ", data.Id.String()),
		)
		return
	}
	if !deleted {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "delete app")
		return
	}

	data.Id = types.StringValue(data.Id.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	stateConf := &retry.StateChangeConf{
		Pending: []string{statusDeleteQueued, statusDeprovisioning},
		Target:  []string{""},
		Refresh: func() (interface{}, string, error) {
			tflog.Trace(ctx, "refreshing app status")
			app, err := r.restClient.GetApp(ctx, data.Id.ValueString())
			if err != nil {
				return "", "", nil
			} else {
				return app.Status, app.Status, nil
			}
		},
		Timeout:    time.Minute * 20,
		Delay:      time.Second * 10,
		MinTimeout: 3 * time.Second,
	}
	statusRaw, err := stateConf.WaitForState()
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "get app")
		return
	}

	status, ok := statusRaw.(string)
	if !ok {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, fmt.Errorf("invalid app %s", status), "delete app")
	}
}

func (r *AppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	// resource.ImportStatePassthroughID(ctx, path.Root("org_id"), req, resp)
}
