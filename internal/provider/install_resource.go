package provider

import (
	"context"
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
var (
	_ resource.Resource                = &InstallResource{}
	_ resource.ResourceWithImportState = &InstallResource{}
)

func NewInstallResource() resource.Resource {
	return &InstallResource{}
}

// InstallResource defines the resource implementation.
type InstallResource struct {
	baseResource
}

type InstallInput struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type AWSAccount struct {
	Region     types.String `tfsdk:"region"`
	IAMRoleARN types.String `tfsdk:"iam_role_arn"`
}

type AzureAccount struct {
	Location                 types.String `tfsdk:"location"`
	SubscriptionID           types.String `tfsdk:"subscription_id"`
	SubscriptionTenantID     types.String `tfsdk:"subscription_tenant_id"`
	ServicePrincipalAppID    types.String `tfsdk:"service_principal_app_id"`
	ServicePrincipalPassword types.String `tfsdk:"service_principal_password"`
}

// InstallResourceModel describes the resource data model.
type InstallResourceModel struct {
	Name  types.String `tfsdk:"name"`
	AppID types.String `tfsdk:"app_id"`

	AWSAccount   []AWSAccount   `tfsdk:"aws"`
	AzureAccount []AzureAccount `tfsdk:"azure"`
	Inputs       []InstallInput `tfsdk:"input"`

	// computed
	ID types.String `tfsdk:"id"`
}

func (r *InstallResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_install"
}

func (r *InstallResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create an install to release an app to.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The unique ID of the install.",
				Optional:    false,
				Required:    true,
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
				Description: "The unique ID of the install",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"azure": schema.SetNestedBlock{
				Description: "Configuration for an Azure install",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"location": schema.StringAttribute{
							Description: "The Azure location to create the install in.",
							Optional:    false,
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"subscription_id": schema.StringAttribute{
							Description: "The subscription id.",
							Optional:    false,
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"subscription_tenant_id": schema.StringAttribute{
							Description: "The subscription tenant id.",
							Optional:    false,
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"service_principal_app_id": schema.StringAttribute{
							Description: "The service principal app id.",
							Optional:    false,
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"service_principal_password": schema.StringAttribute{
							Description: "The service principal password.",
							Optional:    false,
							Required:    true,
							Sensitive:   true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
			"aws": schema.SetNestedBlock{
				Description: "Configuration for an AWS install",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Description: "The AWS region to create the install in.",
							Optional:    false,
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"iam_role_arn": schema.StringAttribute{
							Description: "The ARN of the AWS IAM role to provision the install with.",
							Optional:    false,
							Required:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
			"input": schema.SetNestedBlock{
				Description: "An input on the install, for configuration",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The input name, which must map to a defined app input",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The static value. Interpolation is not supported here.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *InstallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *InstallResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "creating install")
	createReq := &models.ServiceCreateInstallRequest{
		Name:   data.Name.ValueStringPointer(),
		Inputs: make(map[string]string, 0),
	}
	if len(data.AWSAccount) == 1 {
		createReq.AwsAccount = &models.ServiceCreateInstallRequestAwsAccount{
			Region:     data.AWSAccount[0].Region.ValueString(),
			IamRoleArn: data.AWSAccount[0].IAMRoleARN.ValueStringPointer(),
		}
	}
	if len(data.AzureAccount) == 1 {
		createReq.AzureAccount = &models.ServiceCreateInstallRequestAzureAccount{
			Location:                 data.AzureAccount[0].Location.ValueString(),
			ServicePrincipalAppID:    data.AzureAccount[0].ServicePrincipalAppID.ValueString(),
			ServicePrincipalPassword: data.AzureAccount[0].ServicePrincipalPassword.ValueString(),
			SubscriptionID:           data.AzureAccount[0].SubscriptionID.ValueString(),
			SubscriptionTenantID:     data.AzureAccount[0].SubscriptionTenantID.ValueString(),
		}
	}

	for _, input := range data.Inputs {
		createReq.Inputs[input.Name.ValueString()] = input.Value.ValueString()
	}

	installResp, err := r.restClient.CreateInstall(ctx, data.AppID.ValueString(), createReq)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create install")
		return
	}
	data.ID = types.StringValue(installResp.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Trace(ctx, "successfully created install")

	stateConf := &retry.StateChangeConf{
		Pending: []string{statusQueued, statusProvisioning, statusTemporarilyUnavailable},
		Target:  []string{statusActive},
		Refresh: func() (interface{}, string, error) {
			tflog.Trace(ctx, "refreshing install status")
			install, err := r.restClient.GetInstall(ctx, installResp.ID)
			if err == nil {
				return install.SandboxStatus, install.SandboxStatus, nil
			}

			logErr(ctx, err, "create install")
			return statusTemporarilyUnavailable, statusTemporarilyUnavailable, nil
		},
		Timeout:    time.Minute * 45,
		Delay:      time.Second * 10,
		MinTimeout: 3 * time.Second,
	}
	statusRaw, err := stateConf.WaitForState()
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "get install")
		return
	}

	status, ok := statusRaw.(string)
	if !ok {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, fmt.Errorf("invalid install %s", status), "create install")
		return
	}
	if status != statusActive {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, fmt.Errorf("status %s", status), "create install")
		return
	}
}

func (r *InstallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *InstallResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	installResp, err := r.restClient.GetInstall(ctx, data.ID.ValueString())
	if nuon.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "get install")
		return
	}
	data.Name = types.StringValue(installResp.Name)
	data.AppID = types.StringValue(installResp.AppID)

	if installResp.AwsAccount != nil {
		data.AWSAccount = []AWSAccount{
			{
				IAMRoleARN: types.StringValue(installResp.AwsAccount.IamRoleArn),
				Region:     types.StringValue(installResp.AwsAccount.Region),
			},
		}
	}
	if installResp.AzureAccount != nil {
		data.AzureAccount = []AzureAccount{
			{
				Location:                 types.StringValue(installResp.AzureAccount.Location),
				SubscriptionID:           types.StringValue(installResp.AzureAccount.SubscriptionID),
				SubscriptionTenantID:     types.StringValue(installResp.AzureAccount.SubscriptionTenantID),
				ServicePrincipalAppID:    types.StringValue(installResp.AzureAccount.ServicePrincipalAppID),
				ServicePrincipalPassword: types.StringValue(installResp.AzureAccount.ServicePrincipalPassword),
			},
		}
	}

	inputs, err := r.restClient.GetInstallCurrentInputs(ctx, data.ID.ValueString())
	if err != nil && !nuon.IsNotFound(err) {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "get install inputs")
		return
	}

	// if no inputs are found, it means that no inputs were defined, and this is not an actual error, just empty
	// state.
	data.Inputs = []InstallInput{}
	if nuon.IsNotFound(err) {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}
	for name, value := range inputs.Values {
		data.Inputs = append(data.Inputs, InstallInput{
			Name:  types.StringValue(name),
			Value: types.StringValue(value),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InstallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *InstallResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	installResp, err := r.restClient.UpdateInstall(ctx, data.ID.ValueString(), &models.ServiceUpdateInstallRequest{
		Name: data.Name.ValueString(),
	})
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "update install")
		return
	}
	data.ID = types.StringValue(installResp.ID)
	data.Name = types.StringValue(installResp.Name)

	// TODO: The SDK doesn't return these values.
	// These can't be updated anyway, so it's not a blocker,
	// but it would be ideal to use the API as the source of truth.
	// data.AppID = types.StringValue(installResp.AppID)
	// data.IAMRoleARN = types.StringValue(installResp.AwsAccount.IamRoleArn)
	// data.Region = types.StringValue(installResp.AwsAccount.Region)

	updateReq := &models.ServiceCreateInstallInputsRequest{
		Inputs: make(map[string]string, 0),
	}
	for _, input := range data.Inputs {
		updateReq.Inputs[input.Name.ValueString()] = input.Value.ValueString()
	}

	_, err = r.restClient.CreateInstallInputs(ctx, data.ID.ValueString(), updateReq)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "update install")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InstallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *InstallResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleted, err := r.restClient.DeleteInstall(ctx, data.ID.ValueString())
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "delete install")
		return
	}
	if !deleted {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "delete install")
		return
	}

	data.ID = types.StringValue(data.ID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	stateConf := &retry.StateChangeConf{
		Pending: []string{statusDeleteQueued, statusDeprovisioning, statusTemporarilyUnavailable},
		Target:  []string{statusNotFound},
		Refresh: func() (interface{}, string, error) {
			tflog.Trace(ctx, "refreshing install status")
			install, err := r.restClient.GetInstall(ctx, data.ID.ValueString())
			if err == nil {
				return install.SandboxStatus, install.SandboxStatus, nil
			}

			logErr(ctx, err, "delete install")
			if nuon.IsNotFound(err) {
				return "", statusNotFound, nil
			}

			logErr(ctx, err, "delete install")
			return statusTemporarilyUnavailable, statusTemporarilyUnavailable, nil
		},
		Timeout:    time.Minute * 45,
		Delay:      time.Second * 10,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "delete install")
		return
	}
}

func (r *InstallResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
