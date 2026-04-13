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
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Path           types.String `tfsdk:"path"`
	NoteCount      types.Int64  `tfsdk:"note_count"`
	ParentFolderID types.String `tfsdk:"parent_folder_id"`
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
		Description: "Retrieves folder information from the Devin API v3",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The folder_id of the folder resource",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the folder resource",
				Optional:    true,
				Computed:    true,
			},
			"path": schema.StringAttribute{
				Description: "The hierarchical path of the folder",
				Computed:    true,
			},
			"note_count": schema.Int64Attribute{
				Description: "The number of notes in the folder",
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
func (d *FolderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	if !config.ID.IsNull() {
		folderID := config.ID.ValueString()
		tflog.Info(ctx, "Starting folder data retrieval by ID", map[string]interface{}{
			"id": folderID,
		})
		folder, err = d.client.GetFolderByID(folderID)
	} else if !config.Name.IsNull() {
		folderName := config.Name.ValueString()
		tflog.Info(ctx, "Starting folder data retrieval by name", map[string]interface{}{
			"name": folderName,
		})
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

	state := FolderDataSourceModel{
		ID:        types.StringValue(folder.FolderID),
		Name:      types.StringValue(folder.Name),
		Path:      types.StringValue(folder.Path),
		NoteCount: types.Int64Value(int64(folder.NoteCount)),
	}

	if folder.ParentFolderID != "" {
		state.ParentFolderID = types.StringValue(folder.ParentFolderID)
	} else {
		state.ParentFolderID = types.StringValue("")
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Successfully retrieved folder data", map[string]interface{}{
		"id":   folder.FolderID,
		"name": folder.Name,
	})
}
