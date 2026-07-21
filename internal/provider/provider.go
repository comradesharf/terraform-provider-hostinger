// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"os"

	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &hostingerProvider{}
)

// hostingerProvider defines the provider implementation.
type hostingerProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// hostingerProviderModel describes the provider data model.
type hostingerProviderModel struct {
	APIToken types.String `tfsdk:"api_token"`
	Host     types.String `tfsdk:"host"`
}

func (p *hostingerProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hostinger"
	resp.Version = p.version
}

func (p *hostingerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API token used to authenticate with the Hostinger API.",
				Sensitive:           true,
				Optional:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The Hostinger API host URL. Defaults to the production API server.",
				Optional:            true,
			},
		},
	}
}

func (p *hostingerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data hostingerProviderModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.APIToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown API Token",
			"The provider cannot create the client as there is an unknown configuration value for the API token. "+
				"Set the value of the API token in the configuration or use the environment variable HOSTINGER_API_TOKEN.",
		)
	}

	if data.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Host",
			"The provider cannot create the client as there is an unknown configuration value for the host. "+
				"Set the value of the host in the configuration or use the default production API server.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiToken := os.Getenv("HOSTINGER_API_TOKEN")
	host := os.Getenv("HOSTINGER_HOST")

	if !data.APIToken.IsNull() {
		apiToken = data.APIToken.ValueString()
		ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "hostinger_api_token")
		tflog.Info(ctx, "Using API token from configuration")
	}

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
		tflog.Info(ctx, "Using host from configuration")
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing API Token",
			"The provider cannot create the client as there is a missing or empty value for the API token. "+
				"Set the value of the API token in the configuration or use the environment variable HOSTINGER_API_TOKEN.",
		)
	}

	if host == "" {
		host = client.ServerUrlProductionAPIServer
		ctx = tflog.SetField(ctx, "hostinger_host", host)
		tflog.Info(ctx, "Using default production API server")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c, err := client.NewClientWithResponses(
		host,
		client.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+apiToken)
			return nil
		}),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create API client",
			"An unexpected error occurred while creating the Hostinger API client.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *hostingerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *hostingerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDataSourceBillingCatalogs,
		NewDataSourceReachContacts,
		NewDataSourceReachSegments,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &hostingerProvider{
			version: version,
		}
	}
}
