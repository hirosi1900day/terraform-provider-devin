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

// ScheduleResource defines the type for schedule resources
type ScheduleResource struct {
	client *DevinClient
}

// ScheduleResourceModel represents the schema structure for the Terraform resource
type ScheduleResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Prompt     types.String `tfsdk:"prompt"`
	Cron       types.String `tfsdk:"cron"`
	PlaybookID types.String `tfsdk:"playbook_id"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

// NewScheduleResource creates an instance of the schedule resource
func NewScheduleResource() resource.Resource {
	return &ScheduleResource{}
}

// Metadata sets the resource metadata
func (r *ScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schedule"
}

// Schema defines the resource schema
func (r *ScheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages schedule resources in Devin API v3",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The scheduled_session_id of the schedule resource",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the schedule",
				Required:    true,
			},
			"prompt": schema.StringAttribute{
				Description: "The prompt for the scheduled task",
				Required:    true,
			},
			"cron": schema.StringAttribute{
				Description: "Cron expression for the schedule (maps to API 'frequency' field)",
				Required:    true,
			},
			"playbook_id": schema.StringAttribute{
				Description: "The ID of the playbook to use",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the schedule is enabled (read-only)",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Creation timestamp (ISO 8601)",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update timestamp (ISO 8601)",
				Computed:    true,
			},
		},
	}
}

// Configure configures the resource
func (r *ScheduleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a schedule resource
func (r *ScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ScheduleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting schedule resource creation")

	reqBody := CreateScheduleRequest{
		Name:      plan.Name.ValueString(),
		Prompt:    plan.Prompt.ValueString(),
		Frequency: plan.Cron.ValueString(),
	}

	if !plan.PlaybookID.IsNull() && !plan.PlaybookID.IsUnknown() {
		reqBody.PlaybookID = plan.PlaybookID.ValueString()
	}

	schedule, err := r.client.CreateSchedule(reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create schedule",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	mapScheduleToModel(schedule, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Schedule resource creation completed", map[string]interface{}{
		"id": schedule.ScheduledSessionID,
	})
}

// Read reads a schedule resource
func (r *ScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ScheduleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Retrieving schedule resource information", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	schedule, err := r.client.GetSchedule(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to retrieve schedule",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	mapScheduleToModel(schedule, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Schedule resource information retrieval completed")
}

// Update updates a schedule resource (PATCH)
func (r *ScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ScheduleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ScheduleResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting schedule resource update", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	reqBody := UpdateScheduleRequest{}

	name := plan.Name.ValueString()
	reqBody.Name = &name

	prompt := plan.Prompt.ValueString()
	reqBody.Prompt = &prompt

	cron := plan.Cron.ValueString()
	reqBody.Frequency = &cron

	if !plan.PlaybookID.IsNull() && !plan.PlaybookID.IsUnknown() {
		v := plan.PlaybookID.ValueString()
		reqBody.PlaybookID = &v
	}

	schedule, err := r.client.UpdateSchedule(state.ID.ValueString(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update schedule",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	mapScheduleToModel(schedule, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Schedule resource update completed")
}

// Delete deletes a schedule resource
func (r *ScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ScheduleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting schedule resource deletion", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := r.client.DeleteSchedule(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete schedule",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	tflog.Info(ctx, "Schedule resource deletion completed")
}

// ImportState imports a schedule resource
func (r *ScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapScheduleToModel maps a Schedule API response to the Terraform model
func mapScheduleToModel(schedule *Schedule, model *ScheduleResourceModel) {
	model.ID = types.StringValue(schedule.ScheduledSessionID)
	model.Name = types.StringValue(schedule.Name)
	model.Prompt = types.StringValue(schedule.Prompt)
	model.Enabled = types.BoolValue(schedule.Enabled)
	model.CreatedAt = types.StringValue(schedule.CreatedAt)
	model.UpdatedAt = types.StringValue(schedule.UpdatedAt)

	if schedule.Frequency != nil {
		model.Cron = types.StringValue(*schedule.Frequency)
	} else {
		model.Cron = types.StringValue("")
	}

	if schedule.Playbook != nil && schedule.Playbook.PlaybookID != "" {
		model.PlaybookID = types.StringValue(schedule.Playbook.PlaybookID)
	} else {
		model.PlaybookID = types.StringNull()
	}
}
