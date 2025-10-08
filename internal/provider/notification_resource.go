package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NotificationResource{}
var _ resource.ResourceWithImportState = &NotificationResource{}

func NewNotificationResource() resource.Resource {
	return &NotificationResource{}
}

// NotificationResource defines the resource implementation.
type NotificationResource struct {
	client *Client
}

// NotificationResourceModel describes the resource data model.
type NotificationResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	IsDefault     types.Bool   `tfsdk:"is_default"`
	ApplyExisting types.Bool   `tfsdk:"apply_existing"`
	Active        types.Bool   `tfsdk:"active"`
	Config        types.Map    `tfsdk:"config"`
}

func (r *NotificationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification"
}

func (r *NotificationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Uptime Kuma notification resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Notification identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Notification name",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Notification type (e.g., discord, slack, webhook, smtp, telegram, etc.)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_default": schema.BoolAttribute{
				MarkdownDescription: "Whether this notification is enabled by default for new monitors",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"apply_existing": schema.BoolAttribute{
				MarkdownDescription: "Whether to apply this notification to all existing monitors",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the notification is active",
				Computed:            true,
			},
			"config": schema.MapAttribute{
				MarkdownDescription: "Notification-specific configuration parameters",
				ElementType:         types.StringType,
				Optional:            true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *NotificationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NotificationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check if a notification with this name already exists
	existingNotifications, err := r.client.GetNotifications()
	if err != nil {
		resp.Diagnostics.AddWarning("Warning", fmt.Sprintf("Unable to check for existing notifications: %s", err))
	}

	var existingNotification *Notification
	for i := range existingNotifications {
		if existingNotifications[i].Name == data.Name.ValueString() {
			existingNotification = &existingNotifications[i]
			resp.Diagnostics.AddWarning("Adopting Resource", fmt.Sprintf("Found existing notification with name '%s' and ID %d, adopting it", existingNotification.Name, existingNotification.ID))
			break
		}
	}

	// Convert config map to interface map
	config := make(map[string]interface{})
	if !data.Config.IsNull() && !data.Config.IsUnknown() {
		configMap := make(map[string]string)
		diags := data.Config.ElementsAs(ctx, &configMap, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for k, v := range configMap {
			config[k] = v
		}
	}

	// Create notification struct
	notification := &Notification{
		Name:          data.Name.ValueString(),
		Type:          data.Type.ValueString(),
		IsDefault:     data.IsDefault.ValueBool(),
		ApplyExisting: data.ApplyExisting.ValueBool(),
		Config:        config,
	}

	var createdNotification *Notification
	if existingNotification != nil {
		// Adopt the existing notification and update it
		notification.ID = existingNotification.ID
		createdNotification, err = r.client.UpdateNotification(notification)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update existing notification, got error: %s", err))
			return
		}
	} else {
		// Create notification via API
		createdNotification, err = r.client.CreateNotification(notification)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create notification, got error: %s", err))
			return
		}
	}

	// Update the data model with response values
	data.ID = types.StringValue(strconv.Itoa(createdNotification.ID))
	data.Active = types.BoolValue(createdNotification.Active)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NotificationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse notification ID
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse notification ID: %s", err))
		return
	}

	// Get notification from API
	notification, err := r.client.GetNotification(id)
	if err != nil {
		// If the notification is not found, remove it from state (Terraform will recreate it)
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read notification, got error: %s", err))
		return
	}

	// Update the data model with values from API
	data.Name = types.StringValue(notification.Name)
	data.Type = types.StringValue(notification.Type)
	data.IsDefault = types.BoolValue(notification.IsDefault)
	data.ApplyExisting = types.BoolValue(notification.ApplyExisting)
	data.Active = types.BoolValue(notification.Active)

	// Convert config to map
	if notification.Config != nil && len(notification.Config) > 0 {
		configMap := make(map[string]string)
		for k, v := range notification.Config {
			if str, ok := v.(string); ok {
				configMap[k] = str
			} else {
				configMap[k] = fmt.Sprintf("%v", v)
			}
		}
		configValue, diags := types.MapValueFrom(ctx, types.StringType, configMap)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Config = configValue
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NotificationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse notification ID
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse notification ID: %s", err))
		return
	}

	// Convert config map to interface map
	config := make(map[string]interface{})
	if !data.Config.IsNull() && !data.Config.IsUnknown() {
		configMap := make(map[string]string)
		diags := data.Config.ElementsAs(ctx, &configMap, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for k, v := range configMap {
			config[k] = v
		}
	}

	// Create notification struct
	notification := &Notification{
		ID:            id,
		Name:          data.Name.ValueString(),
		Type:          data.Type.ValueString(),
		IsDefault:     data.IsDefault.ValueBool(),
		ApplyExisting: data.ApplyExisting.ValueBool(),
		Config:        config,
	}

	// Update notification via API
	updatedNotification, err := r.client.UpdateNotification(notification)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update notification, got error: %s", err))
		return
	}

	// Update the data model
	data.Active = types.BoolValue(updatedNotification.Active)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NotificationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse notification ID
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse notification ID: %s", err))
		return
	}

	// Delete notification via API
	err = r.client.DeleteNotification(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete notification, got error: %s", err))
		return
	}
}

func (r *NotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by notification ID
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse notification ID: %s", err))
		return
	}

	// Verify notification exists
	_, err = r.client.GetNotification(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to find notification with ID %d: %s", id, err))
		return
	}

	// Set the ID in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
