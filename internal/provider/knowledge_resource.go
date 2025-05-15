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
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Body               types.String `tfsdk:"body"`
	TriggerDescription types.String `tfsdk:"trigger_description"`
	ParentFolderID     types.String `tfsdk:"parent_folder_id"`
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
		Description: "Manages knowledge resources in the Devin API",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the knowledge resource",
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
			"trigger_description": schema.StringAttribute{
				Description: "The trigger description for the knowledge resource",
				Required:    true,
			},
			"parent_folder_id": schema.StringAttribute{
				Description: "The ID of the parent folder",
				Optional:    true,
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
			fmt.Sprintf("Expected *DevinClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	// Create knowledge
	knowledge, err := r.client.CreateKnowledge(
		plan.Name.ValueString(),
		plan.Body.ValueString(),
		plan.TriggerDescription.ValueString(),
		plan.ParentFolderID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create knowledge",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	// Update model
	plan.ID = types.StringValue(knowledge.ID)

	// Save state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Knowledge resource creation completed", map[string]interface{}{
		"id": knowledge.ID,
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

	// Get knowledge
	knowledge, err := r.client.GetKnowledge(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to retrieve knowledge",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	// Update model
	state.Name = types.StringValue(knowledge.Name)
	state.Body = types.StringValue(knowledge.Body)
	state.TriggerDescription = types.StringValue(knowledge.TriggerDescription)

	// Update ParentFolderID only if not null
	if knowledge.ParentFolderID != "" {
		state.ParentFolderID = types.StringValue(knowledge.ParentFolderID)
	} else if !state.ParentFolderID.IsNull() {
		// If API returns empty string but current state has a value, update to empty string
		state.ParentFolderID = types.StringValue("")
	}

	// Save state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	// Maintain existing ID
	plan.ID = state.ID

	// Update knowledge
	_, err := r.client.UpdateKnowledge(
		state.ID.ValueString(),
		plan.Name.ValueString(),
		plan.Body.ValueString(),
		plan.TriggerDescription.ValueString(),
		plan.ParentFolderID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update knowledge",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	// Save state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	// Delete knowledge
	err := r.client.DeleteKnowledge(state.ID.ValueString())
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

	// Set import ID as knowledge ID
	// Set knowledge ID to id attribute in state
	diags := resp.State.SetAttribute(ctx, path.Root("id"), req.ID)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Knowledge resource import completed")
}
