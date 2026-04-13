package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// KnowledgeDataSource defines the type for knowledge data source
type KnowledgeDataSource struct {
	client *DevinClient
}

// KnowledgeDataSourceModel represents the schema structure for the Terraform data source
type KnowledgeDataSourceModel struct {
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

// NewKnowledgeDataSource creates an instance of the knowledge data source
func NewKnowledgeDataSource() datasource.DataSource {
	return &KnowledgeDataSource{}
}

// Metadata sets the data source metadata
func (d *KnowledgeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge"
}

// Schema defines the data source schema
func (d *KnowledgeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves knowledge note information from the Devin API v3",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The note_id of the knowledge resource",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the knowledge resource",
				Computed:    true,
			},
			"body": schema.StringAttribute{
				Description: "The content of the knowledge resource",
				Computed:    true,
			},
			"trigger": schema.StringAttribute{
				Description: "The trigger description for the knowledge resource",
				Computed:    true,
			},
			"folder_id": schema.StringAttribute{
				Description: "The ID of the parent folder",
				Computed:    true,
			},
			"folder_path": schema.StringAttribute{
				Description: "The folder path",
				Computed:    true,
			},
			"pinned_repo": schema.StringAttribute{
				Description: "Pinned repository",
				Computed:    true,
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Whether the knowledge is enabled",
				Computed:    true,
			},
			"macro": schema.StringAttribute{
				Description: "The macro identifier",
				Computed:    true,
			},
			"access_type": schema.StringAttribute{
				Description: "Access type: enterprise or org",
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

// Configure configures the data source
func (d *KnowledgeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*DevinClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected data source configure type",
			fmt.Sprintf("Expected *DevinClient, got: %T.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read reads the knowledge information
func (d *KnowledgeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config KnowledgeDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	noteID := config.ID.ValueString()
	tflog.Info(ctx, "Starting knowledge data retrieval", map[string]interface{}{
		"id": noteID,
	})

	note, err := d.client.GetKnowledgeNote(noteID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to retrieve knowledge",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	// Map response to model
	config.Name = types.StringValue(note.Name)
	config.Body = types.StringValue(note.Body)
	config.Trigger = types.StringValue(note.Trigger)
	config.IsEnabled = types.BoolValue(note.IsEnabled)
	config.FolderPath = types.StringValue(note.FolderPath)
	config.Macro = types.StringValue(note.Macro)
	config.AccessType = types.StringValue(note.AccessType)
	config.CreatedAt = types.Float64Value(note.CreatedAt)
	config.UpdatedAt = types.Float64Value(note.UpdatedAt)

	if note.FolderID != "" {
		config.FolderID = types.StringValue(note.FolderID)
	} else {
		config.FolderID = types.StringValue("")
	}

	if note.PinnedRepo != nil {
		config.PinnedRepo = types.StringValue(*note.PinnedRepo)
	} else {
		config.PinnedRepo = types.StringValue("")
	}

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Knowledge data retrieval completed")
}
