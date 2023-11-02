// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/passbolt/go-passbolt/api"
)

// Ensure PassboltProvider satisfies various provider interfaces.
var _ provider.Provider = &PassboltProvider{}

// PassboltProvider defines the provider implementation.
type PassboltProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// PassboltProviderModel describes the provider data model.
type PassboltProviderModel struct {
	Endpoint       types.String `tfsdk:"endpoint"`
	UserPassword   types.String `tfsdk:"user_password"`
	UserPrivateKey types.String `tfsdk:"user_private_key"`
}

func (p *PassboltProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "passbolt"
	resp.Version = p.version
}

func (p *PassboltProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Endpoint of the Passbolt server",
				Required:            true,
			},
			"user_password": schema.StringAttribute{
				MarkdownDescription: "The secret user password needed to connect to Passbolt",
				Required:            true,
				Sensitive:           true,
			},
			"user_private_key": schema.StringAttribute{
				MarkdownDescription: "The private GPG key needed to connect to Passbolt",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *PassboltProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PassboltProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := api.NewClient(
		nil,
		fmt.Sprintf("PassboltTerraformProvider/%s", p.version),
		data.Endpoint.String(),
		data.UserPrivateKey.String(),
		data.UserPassword.String(),
	)

	if err != nil {
		resp.Diagnostics.AddError("API Error when creating client", err.Error())
		return
	}

	// Client configuration for data sources and resources
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PassboltProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		SecretResource,
	}
}

func (p *PassboltProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		SecretDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PassboltProvider{
			version: version,
		}
	}
}
