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

// Devinプロバイダーの型を定義
type DevinProvider struct {
	// プロバイダーのバージョン
	version string
}

// プロバイダーの設定構造体
type DevinProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

// New関数は新しいDevinプロバイダーインスタンスを返します
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DevinProvider{
			version: version,
		}
	}
}

// Metadata関数はプロバイダーのメタデータを返します
func (p *DevinProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "devin"
	resp.Version = p.version
}

// Schema関数はプロバイダーのスキーマを定義します
func (p *DevinProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "Devin APIのAPI Key。環境変数 DEVIN_API_KEY で設定することもできます。",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure関数はプロバイダーの設定を行います
func (p *DevinProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Devinプロバイダーの設定を開始します")

	var config DevinProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// API Keyの設定
	apiKey := os.Getenv("DEVIN_API_KEY")
	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"API Keyが未設定です",
			"Devin APIのAPI Keyを設定してください。Terraform設定または環境変数 DEVIN_API_KEY で設定できます。",
		)
		return
	}

	// クライアントの作成
	client := NewClient(apiKey)

	resp.ResourceData = client
	resp.DataSourceData = client

	tflog.Info(ctx, "Devinプロバイダーの設定が完了しました")
}

// Resources関数はプロバイダーが提供するリソースを定義します
func (p *DevinProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewKnowledgeResource,
	}
}

// DataSources関数はプロバイダーが提供するデータソースを定義します
func (p *DevinProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewKnowledgeDataSource,
	}
}
