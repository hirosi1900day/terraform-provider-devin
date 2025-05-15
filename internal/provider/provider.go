package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// DevinProvider represents the Terraform provider for Devin
type DevinProvider struct {
	// provider version
	version string
}

// DevinProviderModel represents the provider configuration structure
type DevinProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

// New returns a new instance of the Devin provider
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DevinProvider{
			version: version,
		}
	}
}

// Metadata returns provider metadata
func (p *DevinProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "devin"
	resp.Version = p.version
}

// Schema defines the provider's schema
func (p *DevinProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "API Key for Devin API. Can also be set via the DEVIN_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure configures the provider
func (p *DevinProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Starting Devin provider configuration")

	var config DevinProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Setting API Key
	apiKey := os.Getenv("DEVIN_API_KEY")
	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"API Key not set",
			"Please set the API Key for Devin API. It can be set in Terraform configuration or via the DEVIN_API_KEY environment variable.",
		)
		return
	}

	// Create client
	client := NewClient(apiKey)

	resp.ResourceData = client
	resp.DataSourceData = client

	tflog.Info(ctx, "Devin provider configuration completed")
}

// Resources defines the resources provided by this provider
func (p *DevinProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewKnowledgeResource,
	}
}

// DataSources defines the data sources provided by this provider
func (p *DevinProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewKnowledgeDataSource,
	}
}
