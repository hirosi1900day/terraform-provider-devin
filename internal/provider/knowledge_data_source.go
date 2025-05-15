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
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Body               types.String `tfsdk:"body"`
	TriggerDescription types.String `tfsdk:"trigger_description"`
	ParentFolderID     types.String `tfsdk:"parent_folder_id"`
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
		Description: "Retrieves knowledge information from the Devin API",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the knowledge resource",
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
			"trigger_description": schema.StringAttribute{
				Description: "The trigger description for the knowledge resource",
				Computed:    true,
			},
			"parent_folder_id": schema.StringAttribute{
				Description: "The ID of the parent folder",
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
			fmt.Sprintf("Expected *DevinClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	knowledgeID := config.ID.ValueString()
	tflog.Info(ctx, "Starting knowledge data retrieval", map[string]interface{}{
		"id": knowledgeID,
	})

	// Get knowledge
	knowledge, err := d.client.GetKnowledge(knowledgeID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to retrieve knowledge",
			fmt.Sprintf("Error during Devin API request: %s", err),
		)
		return
	}

	// Set data
	config.Name = types.StringValue(knowledge.Name)
	config.Body = types.StringValue(knowledge.Body)
	config.TriggerDescription = types.StringValue(knowledge.TriggerDescription)
	config.ParentFolderID = types.StringValue(knowledge.ParentFolderID)

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Knowledge data retrieval completed")
}
