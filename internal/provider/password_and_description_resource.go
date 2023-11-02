// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/Cocossoul/passbolt_terraform_provider/internal/model"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/passbolt/go-passbolt/api"
	"github.com/passbolt/go-passbolt/helper"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PasswordAndDescriptionResource{}
var _ resource.ResourceWithImportState = &PasswordAndDescriptionResource{}

func NewPasswordAndDescriptionResource() resource.Resource {
	return &PasswordAndDescriptionResource{}
}

// PasswordAndDescriptionResource defines the resource implementation.
type PasswordAndDescriptionResource struct {
	client *api.Client
}

func (r *PasswordAndDescriptionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_and_description"
}

func (r *PasswordAndDescriptionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "password-and-description type resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Resource name",
				Required:            true,
			},
			"folder_parent_id": schema.StringAttribute{
				MarkdownDescription: "Resource folder ID. Defaults to root.",
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Resource username value. Defaults to empty string",
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"uri": schema.StringAttribute{
				MarkdownDescription: "Resource URI value. Defaults to empty string",
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Resource password value",
				Required:            true,
				Sensitive:           true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Resource description. Defaults to empty string.",
				Sensitive:           true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *PasswordAndDescriptionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *PasswordAndDescriptionResource) Refresh(ctx context.Context, data model.PasswordAndDescriptionResourceModel) (model.PasswordAndDescriptionResourceModel, error) {
	folderParentID, name, username, uri, password, description, err := helper.GetResource(
		ctx,
		r.client,
		data.Id.String(),
	)

	if err != nil {
		return data, err
	}

	data.FolderParentId = types.StringValue(folderParentID)
	data.Name = types.StringValue(name)
	data.Username = types.StringValue(username)
	data.URI = types.StringValue(uri)
	data.Password = types.StringValue(password)
	data.Description = types.StringValue(description)

	return data, nil
}

func (r *PasswordAndDescriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data model.PasswordAndDescriptionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := helper.CreateResource(
		ctx,
		r.client,
		data.FolderParentId.String(),
		data.Name.String(),
		data.Username.String(),
		data.URI.String(),
		data.Password.String(),
		data.Description.String(),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("API Error when creating %s password-and-description resource", data.Name.String()),
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(id)

	tflog.Trace(ctx, fmt.Sprintf("created %s password-and-description resource with id %s", data.Name, data.Id))

	data, err = r.Refresh(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("API Error when refreshing %s password-and-description resource with id %s just after creation", data.Name.String(), data.Id.String()),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PasswordAndDescriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data model.PasswordAndDescriptionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.Refresh(ctx, data)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("API Error when reading password-and-description resource with id %s", data.Id.String()), err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PasswordAndDescriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data model.PasswordAndDescriptionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := helper.UpdateResource(
		ctx,
		r.client,
		data.Id.String(),
		data.Name.String(),
		data.Username.String(),
		data.URI.String(),
		data.Password.String(),
		data.Description.String(),
	)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("API Error when updating %s password-and-description resource with id %s", data.Name.String(), data.Id.String()), err.Error())
		return
	}

	data, err = r.Refresh(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("API Error when refreshing %s password-and-description resource with id %s just after update", data.Name.String(), data.Id.String()),
			err.Error(),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PasswordAndDescriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data model.PasswordAndDescriptionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteResource(ctx, data.Id.String())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("API Error when deleting %s password-and-description resource with id %s", data.Name.String(), data.Id.String()),
			err.Error(),
		)
	}
}

func (r *PasswordAndDescriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
