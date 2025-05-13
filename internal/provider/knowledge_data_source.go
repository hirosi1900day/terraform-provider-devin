package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// KnowledgeDataSource はナレッジデータソースの型定義
type KnowledgeDataSource struct {
	client *DevinClient
}

// KnowledgeDataSourceModel はTerraformデータソースのスキーマを表す構造体
type KnowledgeDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Body               types.String `tfsdk:"body"`
	TriggerDescription types.String `tfsdk:"trigger_description"`
	ParentFolderID     types.String `tfsdk:"parent_folder_id"`
	CreatedAt          types.String `tfsdk:"created_at"`
}

// NewKnowledgeDataSource はナレッジデータソースのインスタンスを作成します
func NewKnowledgeDataSource() datasource.DataSource {
	return &KnowledgeDataSource{}
}

// Metadata はデータソースのメタデータを設定します
func (d *KnowledgeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge"
}

// Schema はデータソースのスキーマを定義します
func (d *KnowledgeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Devin APIのナレッジ情報を取得します",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ナレッジのID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "ナレッジの名前",
				Computed:    true,
			},
			"body": schema.StringAttribute{
				Description: "ナレッジの内容",
				Computed:    true,
			},
			"trigger_description": schema.StringAttribute{
				Description: "ナレッジのトリガー説明",
				Computed:    true,
			},
			"parent_folder_id": schema.StringAttribute{
				Description: "親フォルダのID",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "ナレッジの作成日時",
				Computed:    true,
			},
		},
	}
}

// Configure はデータソースの設定を行います
func (d *KnowledgeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*DevinClient)
	if !ok {
		resp.Diagnostics.AddError(
			"予期しないデータソース設定タイプ",
			fmt.Sprintf("予期していないプロバイダデータ型を受け取りました: %T。*DevinClient を期待していました。", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read はナレッジの情報を読み取ります
func (d *KnowledgeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config KnowledgeDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	knowledgeID := config.ID.ValueString()
	tflog.Info(ctx, "ナレッジデータの取得を開始します", map[string]interface{}{
		"id": knowledgeID,
	})

	// ナレッジの取得
	knowledge, err := d.client.GetKnowledge(knowledgeID)
	if err != nil {
		resp.Diagnostics.AddError(
			"ナレッジの取得に失敗しました",
			fmt.Sprintf("Devin APIへのリクエスト中にエラーが発生しました: %s", err),
		)
		return
	}

	// データを設定
	config.Name = types.StringValue(knowledge.Name)
	config.Body = types.StringValue(knowledge.Body)
	config.TriggerDescription = types.StringValue(knowledge.TriggerDescription)
	config.ParentFolderID = types.StringValue(knowledge.ParentFolderID)
	config.CreatedAt = types.StringValue(knowledge.CreatedAt.Format("2006-01-02T15:04:05Z"))

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "ナレッジデータの取得が完了しました")
}
