package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// FolderDataSource defines the type for folder data source
type FolderDataSource struct {
	client *DevinClient
}

// FolderDataSourceModel represents the schema structure for the Terraform data source
type FolderDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// NewFolderDataSource creates an instance of the folder data source
func NewFolderDataSource() datasource.DataSource {
	return &FolderDataSource{}
}

// Metadata sets the data source metadata
func (d *FolderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_folder"
}

// Schema defines the data source schema
func (d *FolderDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves folder information from the Devin API",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the folder resource",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the folder resource",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the folder resource",
				Computed:    true,
			},
		},
	}
}

// Configure configures the data source
func (d *FolderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read reads the folder information
func (d *FolderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config FolderDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var folder *FolderItem
	var err error

	// Check if we're searching by ID or name
	if !config.ID.IsNull() {
		folderID := config.ID.ValueString()
		tflog.Info(ctx, "Starting folder data retrieval by ID", map[string]interface{}{
			"id": folderID,
		})

		// Get folder by ID
		folder, err = d.client.GetFolderByID(folderID)
	} else if !config.Name.IsNull() {
		folderName := config.Name.ValueString()
		tflog.Info(ctx, "Starting folder data retrieval by name", map[string]interface{}{
			"name": folderName,
		})

		// Get folder by name
		folder, err = d.client.GetFolderByName(folderName)
	} else {
		resp.Diagnostics.AddError(
			"Missing required attributes",
			"Either id or name must be specified",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving folder",
			err.Error(),
		)
		return
	}

	// Map the API response to the Terraform model
	state := FolderDataSourceModel{
		ID:          types.StringValue(folder.ID),
		Name:        types.StringValue(folder.Name),
		Description: types.StringValue(folder.Description),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Successfully retrieved folder data", map[string]interface{}{
		"id":   folder.ID,
		"name": folder.Name,
	})
}
