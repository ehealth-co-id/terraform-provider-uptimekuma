// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ehealth-co-id/terraform-provider-uptimekuma/internal/client"
)

// Ensure UptimeKumaProvider satisfies various provider interfaces.
var _ provider.Provider = &UptimeKumaProvider{}
var _ provider.ProviderWithFunctions = &UptimeKumaProvider{}
var _ provider.ProviderWithEphemeralResources = &UptimeKumaProvider{}

// UptimeKumaProvider defines the provider implementation.
type UptimeKumaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// UptimeKumaProviderModel describes the provider data model.
type UptimeKumaProviderModel struct {
	BaseURL        types.String `tfsdk:"base_url"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	InsecureHTTPS  types.Bool   `tfsdk:"insecure_https"`
	DefaultHeaders types.Map    `tfsdk:"default_headers"`
}

func (p *UptimeKumaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "uptimekuma"
	resp.Version = p.version
}

func (p *UptimeKumaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Interact with Uptime Kuma",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL of the Uptime Kuma instance (e.g., http://localhost:3001 or https://uptime.example.com)",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for authentication",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for authentication",
				Required:            true,
				Sensitive:           true,
			},
			"insecure_https": schema.BoolAttribute{
				MarkdownDescription: "Skip TLS certificate verification",
				Optional:            true,
			},
			"default_headers": schema.MapAttribute{
				MarkdownDescription: "Default HTTP headers to include with every request (e.g., for Cloudflare).",
				Optional:            true,
				Sensitive:           true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (p *UptimeKumaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data UptimeKumaProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	baseURL := data.BaseURL.ValueString()
	username := data.Username.ValueString()
	password := data.Password.ValueString()
	insecureHTTPS := false
	if !data.InsecureHTTPS.IsNull() {
		insecureHTTPS = data.InsecureHTTPS.ValueBool()
	}

	defaultHeaders := map[string]string{}
	if !data.DefaultHeaders.IsNull() && !data.DefaultHeaders.IsUnknown() {
		for k, v := range data.DefaultHeaders.Elements() {
			if vStr, ok := v.(types.String); ok && !vStr.IsNull() && !vStr.IsUnknown() {
				defaultHeaders[k] = vStr.ValueString()
			}
		}
	}

	config := &client.Config{
		BaseURL:        baseURL,
		Username:       username,
		Password:       password,
		InsecureHTTPS:  insecureHTTPS,
		DefaultHeaders: defaultHeaders,
	}

	// Create client
	apiClient, err := client.New(config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Uptime Kuma API Client",
			"An unexpected error occurred when creating the Uptime Kuma API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Uptime Kuma Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *UptimeKumaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewMonitorResource,
		NewStatusPageResource,
	}
}

func (p *UptimeKumaProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		// No ephemeral resources yet
	}
}

func (p *UptimeKumaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// No data sources yet
	}
}

func (p *UptimeKumaProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// No functions yet
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &UptimeKumaProvider{
			version: version,
		}
	}
}
