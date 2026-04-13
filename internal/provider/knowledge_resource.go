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

// KnowledgeResource defines the type for knowledge resources
type KnowledgeResource struct {
	client *DevinClient
}

// KnowledgeResourceModel represents the schema structure for the Terraform resource
type KnowledgeResourceModel struct {
	ID         types.String  `tfsdk:"id"`
	Name       types.String  `tfsdk:"name"`
	Body       types.String  `tfsdk:"body"`
	Trigger    types.String  `tfsdk:"trigger"`
	FolderID   types.String  `tfsdk:"folder_id"`
	FolderPath types.String  `tfsdk:"folder_path"`
	PinnedRepo types.String  `tfsdk:"pinned_repo"`
	IsEnabled  types.Bool    `tfsdk:"is_enabled"`
	Macro      types.String  `tfsdk:"macro"`
	AccessType types.String  `tfsdk:"access_type"`
	CreatedAt  types.Float64 `tfsdk:"created_at"`
	UpdatedAt  types.Float64 `tfsdk:"updated_at"`
}

// NewKnowledgeResource creates an instance of the knowledge resource
func NewKnowledgeResource() resource.Resource {
	return &KnowledgeResource{}
}

// Metadata sets the resource metadata
func (r *KnowledgeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge"
}

// Schema defines the resource schema
func (r *KnowledgeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages knowledge note resources in Devin API v3",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The note_id of the knowledge resource",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the knowledge resource",
				Required:    true,
			},
			"body": schema.StringAttribute{
				Description: "The content of the knowledge resource",
				Required:    true,
			},
			"trigger": schema.StringAttribute{
				Description: "The trigger description for the knowledge resource",
				Required:    true,
			},
			"folder_id": schema.StringAttribute{
				Description: "The ID of the parent folder (read-only, managed via Devin UI)",
				Computed:    true,
			},
			"pinned_repo": schema.StringAttribute{
				Description: "Pinned repository (owner/repo format)",
				Optional:    true,
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the knowledge is enabled (read-only, managed via Devin UI)",
				Computed:    true,
			},
			"folder_path": schema.StringAttribute{
				Description: "The folder path (read-only)",
				Computed:    true,
			},
			"macro": schema.StringAttribute{
				Description: "The macro identifier (read-only)",
				Computed:    true,
			},
			"access_type": schema.StringAttribute{
				Description: "Access type: enterprise or org (read-only)",
				Computed:    true,
			},
			"created_at": schema.Float64Attribute{
				Description: "Creation timestamp (UNIX)",
				Computed:    true,
			},
			"updated_at": schema.Float64Attribute{
				Description: "Last update timestamp (UNIX)",
				Computed:    true,
			},
		},
	}
}

// Configure configures the resource
func (r *KnowledgeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*DevinClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("Expected *DevinClient, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates a knowledge resource
func (r *KnowledgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KnowledgeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting knowledge resource creation")

	reqBody := CreateKnowledgeNoteRequest{
		Name:    plan.Name.ValueString(),
		Body:    plan.Body.ValueString(),
		Trigger: plan.Trigger.ValueString(),
	}

	if !plan.PinnedRepo.IsNull() && !plan.PinnedRepo.IsUnknown() {
		v := plan.PinnedRepo.ValueString()
		reqBody.PinnedRepo = &v
	}

	note, err := r.client.CreateKnowledgeNote(reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create knowledge",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	mapKnowledgeNoteToModel(note, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Knowledge resource creation completed", map[string]interface{}{
		"id": note.NoteID,
	})
}

// Read reads a knowledge resource
func (r *KnowledgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KnowledgeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Retrieving knowledge resource information", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	note, err := r.client.GetKnowledgeNote(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to retrieve knowledge",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	mapKnowledgeNoteToModel(note, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Knowledge resource information retrieval completed")
}

// Update updates a knowledge resource
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

	tflog.Info(ctx, "Starting knowledge resource update", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	reqBody := UpdateKnowledgeNoteRequest{
		Name:    plan.Name.ValueString(),
		Body:    plan.Body.ValueString(),
		Trigger: plan.Trigger.ValueString(),
	}

	if !plan.PinnedRepo.IsNull() && !plan.PinnedRepo.IsUnknown() {
		v := plan.PinnedRepo.ValueString()
		reqBody.PinnedRepo = &v
	}

	note, err := r.client.UpdateKnowledgeNote(state.ID.ValueString(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update knowledge",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	mapKnowledgeNoteToModel(note, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Knowledge resource update completed")
}

// Delete deletes a knowledge resource
func (r *KnowledgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state KnowledgeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting knowledge resource deletion", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := r.client.DeleteKnowledgeNote(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete knowledge",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	tflog.Info(ctx, "Knowledge resource deletion completed")
}

// ImportState imports a knowledge resource
func (r *KnowledgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Starting knowledge resource import", map[string]interface{}{
		"id": req.ID,
	})

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	tflog.Info(ctx, "Knowledge resource import completed")
}

// mapKnowledgeNoteToModel maps a KnowledgeNote API response to the Terraform model
func mapKnowledgeNoteToModel(note *KnowledgeNote, model *KnowledgeResourceModel) {
	model.ID = types.StringValue(note.NoteID)
	model.Name = types.StringValue(note.Name)
	model.Body = types.StringValue(note.Body)
	model.Trigger = types.StringValue(note.Trigger)
	model.IsEnabled = types.BoolValue(note.IsEnabled)
	model.FolderPath = types.StringValue(note.FolderPath)
	model.Macro = types.StringValue(note.Macro)
	model.AccessType = types.StringValue(note.AccessType)
	model.CreatedAt = types.Float64Value(note.CreatedAt)
	model.UpdatedAt = types.Float64Value(note.UpdatedAt)

	if note.FolderID != "" {
		model.FolderID = types.StringValue(note.FolderID)
	} else {
		model.FolderID = types.StringNull()
	}

	if note.PinnedRepo != nil {
		model.PinnedRepo = types.StringValue(*note.PinnedRepo)
	} else {
		model.PinnedRepo = types.StringNull()
	}
}
