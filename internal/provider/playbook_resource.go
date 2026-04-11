package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// PlaybookResource defines the type for playbook resources
type PlaybookResource struct {
	client *DevinClient
}

// PlaybookResourceModel represents the schema structure for the Terraform resource
type PlaybookResourceModel struct {
	ID                types.String  `tfsdk:"id"`
	Title             types.String  `tfsdk:"title"`
	Body              types.String  `tfsdk:"body"`
	Status            types.String  `tfsdk:"status"`
	AccessType        types.String  `tfsdk:"access_type"`
	Macro             types.String  `tfsdk:"macro"`
	OrgID             types.String  `tfsdk:"org_id"`
	CreatedAt         types.Float64 `tfsdk:"created_at"`
	UpdatedAt         types.Float64 `tfsdk:"updated_at"`
	CreatedByUserID   types.String  `tfsdk:"created_by_user_id"`
	CreatedByUserName types.String  `tfsdk:"created_by_user_name"`
	UpdatedByUserID   types.String  `tfsdk:"updated_by_user_id"`
	UpdatedByUserName types.String  `tfsdk:"updated_by_user_name"`
}

// NewPlaybookResource creates an instance of the playbook resource
func NewPlaybookResource() resource.Resource {
	return &PlaybookResource{}
}

// Metadata sets the resource metadata
func (r *PlaybookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_playbook"
}

// Schema defines the resource schema
func (r *PlaybookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages playbook resources in Devin API v3",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The playbook_id of the playbook resource",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "The title of the playbook",
				Required:    true,
			},
			"body": schema.StringAttribute{
				Description: "The body content of the playbook",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the playbook (active/inactive)",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("active"),
			},
			"access_type": schema.StringAttribute{
				Description: "Access type (read-only)",
				Computed:    true,
			},
			"macro": schema.StringAttribute{
				Description: "Macro shortcut identifier (read-only)",
				Computed:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "Organization ID (read-only)",
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
			"created_by_user_id": schema.StringAttribute{
				Description: "User ID of the creator (read-only)",
				Computed:    true,
			},
			"created_by_user_name": schema.StringAttribute{
				Description: "User name of the creator (read-only)",
				Computed:    true,
			},
			"updated_by_user_id": schema.StringAttribute{
				Description: "User ID of the last updater (read-only)",
				Computed:    true,
			},
			"updated_by_user_name": schema.StringAttribute{
				Description: "User name of the last updater (read-only)",
				Computed:    true,
			},
		},
	}
}

// Configure configures the resource
func (r *PlaybookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a playbook resource
func (r *PlaybookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PlaybookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting playbook resource creation")

	reqBody := CreatePlaybookRequest{
		Title: plan.Title.ValueString(),
		Body:  plan.Body.ValueString(),
	}

	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		reqBody.Status = plan.Status.ValueString()
	}

	playbook, err := r.client.CreatePlaybook(reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create playbook",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	mapPlaybookToModel(playbook, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Playbook resource creation completed", map[string]interface{}{
		"id": playbook.PlaybookID,
	})
}

// Read reads a playbook resource
func (r *PlaybookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PlaybookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Retrieving playbook resource information", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	playbook, err := r.client.GetPlaybook(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to retrieve playbook",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	mapPlaybookToModel(playbook, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Playbook resource information retrieval completed")
}

// Update updates a playbook resource
func (r *PlaybookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PlaybookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state PlaybookResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting playbook resource update", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	reqBody := UpdatePlaybookRequest{
		Title: plan.Title.ValueString(),
		Body:  plan.Body.ValueString(),
	}

	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		reqBody.Status = plan.Status.ValueString()
	}

	playbook, err := r.client.UpdatePlaybook(state.ID.ValueString(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update playbook",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	mapPlaybookToModel(playbook, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Playbook resource update completed")
}

// Delete deletes a playbook resource
func (r *PlaybookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PlaybookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting playbook resource deletion", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := r.client.DeletePlaybook(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete playbook",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	tflog.Info(ctx, "Playbook resource deletion completed")
}

// ImportState imports a playbook resource
func (r *PlaybookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapPlaybookToModel maps a Playbook API response to the Terraform model
func mapPlaybookToModel(playbook *Playbook, model *PlaybookResourceModel) {
	model.ID = types.StringValue(playbook.PlaybookID)
	model.Title = types.StringValue(playbook.Title)
	model.Body = types.StringValue(playbook.Body)
	model.Status = types.StringValue(playbook.Status)
	model.AccessType = types.StringValue(playbook.AccessType)
	model.OrgID = types.StringValue(playbook.OrgID)
	model.CreatedAt = types.Float64Value(playbook.CreatedAt)
	model.UpdatedAt = types.Float64Value(playbook.UpdatedAt)
	model.CreatedByUserID = types.StringValue(playbook.CreatedByUserID)
	model.CreatedByUserName = types.StringValue(playbook.CreatedByUserName)
	model.UpdatedByUserID = types.StringValue(playbook.UpdatedByUserID)
	model.UpdatedByUserName = types.StringValue(playbook.UpdatedByUserName)

	if playbook.Macro != nil {
		model.Macro = types.StringValue(*playbook.Macro)
	} else {
		model.Macro = types.StringValue("")
	}
}
