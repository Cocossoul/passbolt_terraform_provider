// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/Cocossoul/passbolt_terraform_provider/internal/model"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/passbolt/go-passbolt/api"
	"github.com/passbolt/go-passbolt/helper"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PasswordAndDescriptionDataSource{}

func NewPasswordAndDescriptionDataSource() datasource.DataSource {
	return &PasswordAndDescriptionDataSource{}
}

// PasswordAndDescriptionDataSource defines the data source implementation.
type PasswordAndDescriptionDataSource struct {
	client *api.Client
}

// PasswordAndDescriptionDataSourceModel describes the data source data model.
func (d *PasswordAndDescriptionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_and_description"
}

func (d *PasswordAndDescriptionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "password-and-description type resource data source",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Resource name",
				Required:            true,
			},
			"folder_parent_id": schema.StringAttribute{
				MarkdownDescription: "Resource folder ID.",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Resource username value.",
				Computed:            true,
			},
			"uri": schema.StringAttribute{
				MarkdownDescription: "Resource URI value.",
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Resource password value",
				Computed:            true,
				Sensitive:           true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Resource description.",
				Sensitive:           true,
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Resource ID",
			},
		},
	}
}

func (d *PasswordAndDescriptionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *PasswordAndDescriptionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.PasswordAndDescriptionResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	folderParentID, name, username, uri, password, description, err := helper.GetResource(
		ctx,
		d.client,
		data.Id.String(),
	)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("API Error when reading password-and-description resource with id %s", data.Id.String()), err.Error())
		return
	}

	data.FolderParentId = types.StringValue(folderParentID)
	data.Name = types.StringValue(name)
	data.Username = types.StringValue(username)
	data.URI = types.StringValue(uri)
	data.Password = types.StringValue(password)
	data.Description = types.StringValue(description)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
