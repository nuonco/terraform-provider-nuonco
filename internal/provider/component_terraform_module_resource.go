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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/nuonco/nuon-go"
	"github.com/nuonco/nuon-go/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TerraformModuleComponentResource{}
var _ resource.ResourceWithImportState = &TerraformModuleComponentResource{}

func NewTerraformModuleComponentResource() resource.Resource {
	return &TerraformModuleComponentResource{}
}

// TerraformModuleComponentResource defines the resource implementation.
type TerraformModuleComponentResource struct {
	baseResource
}

type TerraformVariable struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// TerraformModuleComponentResourceModel describes the resource data model.
type TerraformModuleComponentResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name             types.String        `tfsdk:"name"`
	VarName          types.String        `tfsdk:"var_name"`
	Dependencies     types.List          `tfsdk:"dependencies"`
	AppID            types.String        `tfsdk:"app_id"`
	TerraformVersion types.String        `tfsdk:"terraform_version"`
	PublicRepo       *PublicRepo         `tfsdk:"public_repo"`
	ConnectedRepo    *ConnectedRepo      `tfsdk:"connected_repo"`
	Var              []TerraformVariable `tfsdk:"var"`
	EnvVar           []EnvVar            `tfsdk:"env_var"`
}

func (r *TerraformModuleComponentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_terraform_module_component"
}

func (r *TerraformModuleComponentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Release a terraform module.",
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
			"terraform_version": schema.StringAttribute{
				Description: "The version of Terraform to use.",
				Optional:    true,
				Default:     stringdefault.StaticString("1.5.3"),
				Computed:    true,
			},
			"public_repo":    publicRepoAttribute(),
			"connected_repo": connectedRepoAttribute(),
		},
		Blocks: map[string]schema.Block{
			"var": schema.SetNestedBlock{
				Description: "Terraform variables to set when applying the Terraform configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The variable name to write to the terraform.tfvars file (e.g. bucket_name or db_name.)",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The variable value to write to the terraform.tfvars file. Can be any valid Terraform value, or interpolated from Nuon.",
							Required:    true,
						},
					},
				},
			},
			"env_var": envVarSharedBlock(),
		},
	}
}

func (r *TerraformModuleComponentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TerraformModuleComponentResourceModel

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
	data.Name = types.StringValue(compResp.Name)
	data.VarName = types.StringValue(compResp.VarName)

	configRequest := &models.ServiceCreateTerraformModuleComponentConfigRequest{
		ConnectedGithubVcsConfig: nil,
		PublicGitVcsConfig:       nil,
		Variables:                map[string]string{},
		EnvVars:                  map[string]string{},
		Version:                  data.TerraformVersion.ValueString(),
	}
	for _, val := range data.Var {
		configRequest.Variables[val.Name.ValueString()] = val.Value.ValueString()
	}
	for _, val := range data.EnvVar {
		configRequest.EnvVars[val.Name.ValueString()] = val.Value.ValueString()
	}

	if data.PublicRepo != nil {
		configRequest.PublicGitVcsConfig = &models.ServicePublicGitVCSConfigRequest{
			Branch:    data.PublicRepo.Branch.ValueStringPointer(),
			Directory: data.PublicRepo.Directory.ValueStringPointer(),
			Repo:      data.PublicRepo.Repo.ValueStringPointer(),
		}
	} else if data.ConnectedRepo != nil {
		configRequest.ConnectedGithubVcsConfig = &models.ServiceConnectedGithubVCSConfigRequest{
			Branch:    data.ConnectedRepo.Branch.ValueString(),
			Directory: data.ConnectedRepo.Directory.ValueStringPointer(),
			Repo:      data.ConnectedRepo.Repo.ValueStringPointer(),
		}
	}

	_, err = r.restClient.CreateTerraformModuleComponentConfig(ctx, compResp.ID, configRequest)
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

func (r *TerraformModuleComponentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// get terraform model
	var data *TerraformModuleComponentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get component from api
	compResp, err := r.restClient.GetComponent(ctx, data.ID.ValueString())
	if nuon.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "get component")
		return
	}

	// get latest config from api
	configResp, err := r.restClient.GetComponentLatestConfig(ctx, data.ID.ValueString())
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "get component config")
		return
	}
	if configResp.TerraformModule == nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, errors.New("did not get terraform config"), "get component config")
		return
	}
	terraformConfig := configResp.TerraformModule

	// populate terraform model with data from api
	data.Name = types.StringValue(compResp.Name)
	data.VarName = types.StringValue(compResp.VarName)
	data.AppID = types.StringValue(compResp.AppID)
	data.TerraformVersion = types.StringValue(terraformConfig.Version)
	if terraformConfig.ConnectedGithubVcsConfig != nil {
		connected := terraformConfig.ConnectedGithubVcsConfig
		data.ConnectedRepo = &ConnectedRepo{
			Branch:    types.StringValue(connected.Branch),
			Directory: types.StringValue(connected.Directory),
			Repo:      types.StringValue(connected.Repo),
		}
	}
	if terraformConfig.PublicGitVcsConfig != nil {
		public := terraformConfig.PublicGitVcsConfig
		data.PublicRepo = &PublicRepo{
			Branch:    types.StringValue(public.Branch),
			Directory: types.StringValue(public.Directory),
			Repo:      types.StringValue(public.Repo),
		}
	}
	apiVars := []TerraformVariable{}
	for key, val := range terraformConfig.Variables {
		apiVars = append(apiVars, TerraformVariable{
			Name:  types.StringValue(key),
			Value: types.StringValue(val),
		})
	}
	data.Var = apiVars

	envVars := []EnvVar{}
	for key, val := range terraformConfig.EnvVars {
		envVars = append(envVars, EnvVar{
			Name:  types.StringValue(key),
			Value: types.StringValue(val),
		})
	}
	data.EnvVar = envVars

	// return populated terraform model
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "successfully read component")
}

func (r *TerraformModuleComponentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TerraformModuleComponentResourceModel

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
}

func (r *TerraformModuleComponentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TerraformModuleComponentResourceModel

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

	configRequest := &models.ServiceCreateTerraformModuleComponentConfigRequest{
		ConnectedGithubVcsConfig: nil,
		PublicGitVcsConfig:       nil,
		Variables:                map[string]string{},
		EnvVars:                  map[string]string{},
		Version:                  data.TerraformVersion.ValueString(),
	}
	for _, value := range data.Var {
		configRequest.Variables[value.Name.ValueString()] = value.Value.ValueString()
	}
	for _, value := range data.EnvVar {
		configRequest.EnvVars[value.Name.ValueString()] = value.Value.ValueString()
	}
	if data.PublicRepo != nil {
		configRequest.PublicGitVcsConfig = &models.ServicePublicGitVCSConfigRequest{
			Branch:    data.PublicRepo.Branch.ValueStringPointer(),
			Directory: data.PublicRepo.Directory.ValueStringPointer(),
			Repo:      data.PublicRepo.Repo.ValueStringPointer(),
		}
	} else {
		configRequest.ConnectedGithubVcsConfig = &models.ServiceConnectedGithubVCSConfigRequest{
			Branch:    data.ConnectedRepo.Branch.ValueString(),
			Directory: data.ConnectedRepo.Directory.ValueStringPointer(),
			Repo:      data.ConnectedRepo.Repo.ValueStringPointer(),
		}
	}
	_, err = r.restClient.CreateTerraformModuleComponentConfig(ctx, compResp.ID, configRequest)
	if err != nil {
		writeDiagnosticsErr(ctx, &resp.Diagnostics, err, "create component config")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TerraformModuleComponentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
