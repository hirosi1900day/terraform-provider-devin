package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SecretResource defines the type for secret resources
type SecretResource struct {
	client *DevinClient
}

// SecretResourceModel represents the schema structure for the Terraform resource
type SecretResourceModel struct {
	ID        types.String  `tfsdk:"id"`
	Name      types.String  `tfsdk:"name"`
	Value     types.String  `tfsdk:"value"`
	CreatedAt types.Float64 `tfsdk:"created_at"`
	UpdatedAt types.Float64 `tfsdk:"updated_at"`
}

// NewSecretResource creates an instance of the secret resource
func NewSecretResource() resource.Resource {
	return &SecretResource{}
}

// Metadata sets the resource metadata
func (r *SecretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

// Schema defines the resource schema
func (r *SecretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages secret resources in Devin API v3. Note: Secrets cannot be updated. Changing name or value will force recreation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The secret_id of the secret resource",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the secret",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Description: "The value of the secret (write-only, not stored in state)",
				Required:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
func (r *SecretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a secret resource
func (r *SecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SecretResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting secret resource creation")

	reqBody := CreateSecretRequest{
		Name:  plan.Name.ValueString(),
		Value: plan.Value.ValueString(),
	}

	secret, err := r.client.CreateSecret(reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create secret",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(secret.SecretID)
	plan.CreatedAt = types.Float64Value(secret.CreatedAt)
	plan.UpdatedAt = types.Float64Value(secret.UpdatedAt)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Secret resource creation completed", map[string]interface{}{
		"id": secret.SecretID,
	})
}

// Read reads a secret resource
func (r *SecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SecretResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Retrieving secret resource information", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	secret, err := r.client.GetSecretByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to retrieve secret",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	state.Name = types.StringValue(secret.Name)
	state.CreatedAt = types.Float64Value(secret.CreatedAt)
	state.UpdatedAt = types.Float64Value(secret.UpdatedAt)
	// value is write-only, keep the existing state value

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Secret resource information retrieval completed")
}

// Update is not supported for secrets (ForceNew handles recreation)
func (r *SecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Secrets cannot be updated. Delete and recreate the secret instead.",
	)
}

// Delete deletes a secret resource
func (r *SecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SecretResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting secret resource deletion", map[string]interface{}{
		"id": state.ID.ValueString(),
	})

	err := r.client.DeleteSecret(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete secret",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	tflog.Info(ctx, "Secret resource deletion completed")
}
