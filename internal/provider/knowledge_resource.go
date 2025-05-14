package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// KnowledgeResource はナレッジリソースの型定義
type KnowledgeResource struct {
	client *DevinClient
}

// KnowledgeResourceModel はTerraformリソースのスキーマを表す構造体
type KnowledgeResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Body               types.String `tfsdk:"body"`
	TriggerDescription types.String `tfsdk:"trigger_description"`
	ParentFolderID     types.String `tfsdk:"parent_folder_id"`
}

// NewKnowledgeResource はナレッジリソースのインスタンスを作成します
func NewKnowledgeResource() resource.Resource {
	return &KnowledgeResource{}
}

// Metadata はリソースのメタデータを設定します
func (r *KnowledgeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge"
}

// Schema はリソースのスキーマを定義します
func (r *KnowledgeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Devin APIのナレッジリソースを管理します",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ナレッジのID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "ナレッジの名前",
				Required:    true,
			},
			"body": schema.StringAttribute{
				Description: "ナレッジの内容",
				Required:    true,
			},
			"trigger_description": schema.StringAttribute{
				Description: "ナレッジのトリガー説明",
				Required:    true,
			},
			"parent_folder_id": schema.StringAttribute{
				Description: "親フォルダのID",
				Optional:    true,
			},
		},
	}
}

// Configure はリソースの設定を行います
func (r *KnowledgeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*DevinClient)
	if !ok {
		resp.Diagnostics.AddError(
			"予期しないリソース設定タイプ",
			fmt.Sprintf("予期していないプロバイダデータ型を受け取りました: %T。*DevinClient を期待していました。", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create はナレッジリソースを作成します
func (r *KnowledgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KnowledgeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "ナレッジリソースの作成を開始します")

	// ナレッジの作成
	knowledge, err := r.client.CreateKnowledge(
		plan.Name.ValueString(),
		plan.Body.ValueString(),
		plan.TriggerDescription.ValueString(),
		plan.ParentFolderID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"ナレッジの作成に失敗しました",
			fmt.Sprintf("Devin APIへのリクエスト中にエラーが発生しました: %s", err),
		)
		return
	}

	// モデルを更新
	plan.ID = types.StringValue(knowledge.ID)

	// 状態を保存
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "ナレッジリソースの作成が完了しました", map[string]interface{}{
		"id": knowledge.ID,
	})
}

// Read はナレッジリソースの情報を読み取ります
func (r *KnowledgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KnowledgeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "ナレッジリソースの情報を取得します", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	// ナレッジの取得
	knowledge, err := r.client.GetKnowledge(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"ナレッジの取得に失敗しました",
			fmt.Sprintf("Devin APIへのリクエスト中にエラーが発生しました: %s", err),
		)
		return
	}

	// モデルを更新
	state.Name = types.StringValue(knowledge.Name)
	state.Body = types.StringValue(knowledge.Body)
	state.TriggerDescription = types.StringValue(knowledge.TriggerDescription)

	// ParentFolderIDがnullでない場合のみ値を更新
	if knowledge.ParentFolderID != "" {
		state.ParentFolderID = types.StringValue(knowledge.ParentFolderID)
	} else if !state.ParentFolderID.IsNull() {
		// APIが空文字を返したが、現在の状態には値がある場合は空文字に更新
		state.ParentFolderID = types.StringValue("")
	}

	// 状態を保存
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "ナレッジリソースの情報取得が完了しました")
}

// Update はナレッジリソースを更新します
func (r *KnowledgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan KnowledgeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state KnowledgeResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "ナレッジリソースの更新を開始します", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	// 既存のIDを維持
	plan.ID = state.ID

	// ナレッジの更新
	_, err := r.client.UpdateKnowledge(
		state.ID.ValueString(),
		plan.Name.ValueString(),
		plan.Body.ValueString(),
		plan.TriggerDescription.ValueString(),
		plan.ParentFolderID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"ナレッジの更新に失敗しました",
			fmt.Sprintf("Devin APIへのリクエスト中にエラーが発生しました: %s", err),
		)
		return
	}

	// 状態を保存
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "ナレッジリソースの更新が完了しました")
}

// Delete はナレッジリソースを削除します
func (r *KnowledgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state KnowledgeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "ナレッジリソースの削除を開始します", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	// ナレッジの削除
	err := r.client.DeleteKnowledge(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"ナレッジの削除に失敗しました",
			fmt.Sprintf("Devin APIへのリクエスト中にエラーが発生しました: %s", err),
		)
		return
	}

	tflog.Info(ctx, "ナレッジリソースの削除が完了しました")
}

// ImportState はナレッジリソースをインポートします
func (r *KnowledgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "ナレッジリソースのインポートを開始します", map[string]interface{}{
		"id": req.ID,
	})

	// インポートIDをナレッジIDとして設定
	// ナレッジIDを状態のid属性に設定
	diags := resp.State.SetAttribute(ctx, path.Root("id"), req.ID)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "ナレッジリソースのインポートが完了しました")
}
