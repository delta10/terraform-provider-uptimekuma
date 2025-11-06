package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MonitorResource{}
var _ resource.ResourceWithImportState = &MonitorResource{}

func NewMonitorResource() resource.Resource {
	return &MonitorResource{}
}

// MonitorResource defines the resource implementation.
type MonitorResource struct {
	client *Client
}

// MonitorResourceModel describes the resource data model.
type MonitorResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Type                  types.String `tfsdk:"type"`
	URL                   types.String `tfsdk:"url"`
	Hostname              types.String `tfsdk:"hostname"`
	Port                  types.Int64  `tfsdk:"port"`
	Interval              types.Int64  `tfsdk:"interval"`
	Timeout               types.Int64  `tfsdk:"timeout"`
	RetryInterval         types.Int64  `tfsdk:"retry_interval"`
	ResendInterval        types.Int64  `tfsdk:"resend_interval"`
	MaxRetries            types.Int64  `tfsdk:"max_retries"`
	UpsideDown            types.Bool   `tfsdk:"upside_down"`
	MaxRedirects          types.Int64  `tfsdk:"max_redirects"`
	AcceptedStatusCodes   types.List   `tfsdk:"accepted_status_codes"`
	FollowRedirect        types.Bool   `tfsdk:"follow_redirect"`
	Tags                  types.List   `tfsdk:"tags"`
	NotificationIDList    types.List   `tfsdk:"notification_id_list"`
	Active                types.Bool   `tfsdk:"active"`
	IgnoreTLS             types.Bool   `tfsdk:"ignore_tls"`
	HTTPMethod            types.String `tfsdk:"http_method"`
	Body                  types.String `tfsdk:"body"`
	BasicAuthUser         types.String `tfsdk:"basic_auth_user"`
	BasicAuthPass         types.String `tfsdk:"basic_auth_pass"`
}

func (r *MonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *MonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Uptime Kuma monitor resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Monitor identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Monitor name",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Monitor type (http, tcp, ping, etc.)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("http"),
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL to monitor (for HTTP monitors)",
				Optional:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname to monitor (for TCP/Ping monitors)",
				Optional:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port to monitor (for TCP monitors)",
				Optional:            true,
			},
			"interval": schema.Int64Attribute{
				MarkdownDescription: "Check interval in seconds",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(60),
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "Timeout in seconds",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(30),
			},
			"retry_interval": schema.Int64Attribute{
				MarkdownDescription: "Retry interval in seconds",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(60),
			},
			"resend_interval": schema.Int64Attribute{
				MarkdownDescription: "Resend interval in seconds",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"max_retries": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of retries",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(3),
			},
			"upside_down": schema.BoolAttribute{
				MarkdownDescription: "Upside down mode (monitor expects failure)",
				Optional:            true,
			},
			"max_redirects": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of redirects to follow",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(10),
			},
			"accepted_status_codes": schema.ListAttribute{
				MarkdownDescription: "List of accepted HTTP status codes",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"follow_redirect": schema.BoolAttribute{
				MarkdownDescription: "Follow redirects",
				Optional:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "List of tags",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"notification_id_list": schema.ListAttribute{
				MarkdownDescription: "List of notification IDs to associate with this monitor",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the monitor is active",
				Optional:            true,
			},
			"ignore_tls": schema.BoolAttribute{
				MarkdownDescription: "Ignore TLS certificate errors",
				Optional:            true,
			},
			"http_method": schema.StringAttribute{
				MarkdownDescription: "HTTP method (GET, POST, etc.)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("GET"),
			},
			"body": schema.StringAttribute{
				MarkdownDescription: "Request body for POST/PUT requests",
				Optional:            true,
			},
			"basic_auth_user": schema.StringAttribute{
				MarkdownDescription: "Basic authentication username",
				Optional:            true,
			},
			"basic_auth_pass": schema.StringAttribute{
				MarkdownDescription: "Basic authentication password",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *MonitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MonitorResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API model
	monitor := &Monitor{
		Name:           data.Name.ValueString(),
		Type:           data.Type.ValueString(),
		URL:            data.URL.ValueString(),
		Hostname:       data.Hostname.ValueString(),
		Port:           int(data.Port.ValueInt64()),
		Interval:       int(data.Interval.ValueInt64()),
		Timeout:        int(data.Timeout.ValueInt64()),
		RetryInterval:  int(data.RetryInterval.ValueInt64()),
		ResendInterval: int(data.ResendInterval.ValueInt64()),
		MaxRetries:     int(data.MaxRetries.ValueInt64()),
		UpsideDown:     data.UpsideDown.ValueBool(),
		MaxRedirects:   int(data.MaxRedirects.ValueInt64()),
		FollowRedirect: data.FollowRedirect.ValueBool(),
		Active:         data.Active.ValueBool(),
		IgnoreTLS:      data.IgnoreTLS.ValueBool(),
		HTTPMethod:     data.HTTPMethod.ValueString(),
		Body:           data.Body.ValueString(),
		BasicAuthUser:  data.BasicAuthUser.ValueString(),
		BasicAuthPass:  data.BasicAuthPass.ValueString(),
	}

	// Convert lists
	if !data.AcceptedStatusCodes.IsNull() {
		var statusCodes []string
		data.AcceptedStatusCodes.ElementsAs(ctx, &statusCodes, false)
		monitor.AcceptedStatusCodes = statusCodes
	}

	if !data.Tags.IsNull() {
		var tags []string
		data.Tags.ElementsAs(ctx, &tags, false)
		monitor.Tags = tags
	}

	if !data.NotificationIDList.IsNull() {
		var notificationIDs []string
		data.NotificationIDList.ElementsAs(ctx, &notificationIDs, false)
		for _, idStr := range notificationIDs {
			if id, err := strconv.Atoi(idStr); err == nil {
				monitor.NotificationIDList = append(monitor.NotificationIDList, id)
			}
		}
	}

	// Create new monitor
	createdMonitor, err := r.client.CreateMonitor(monitor)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create monitor, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Created new monitor with ID %d", createdMonitor.ID))

	// Update the model with the created monitor ID
	data.ID = types.StringValue(strconv.Itoa(createdMonitor.ID))

	// Don't read back from server - preserve plan values to avoid inconsistent state errors
	// The state should reflect what we sent to the API

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a monitor resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MonitorResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to int
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor ID: %s", err))
		return
	}

	// Refresh monitor data from the API to ensure we have the latest state
	err = r.client.RefreshMonitors()
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Unable to refresh monitors: %s", err))
	}

	// Get monitor from API
	monitor, err := r.client.GetMonitor(id)
	if err != nil {
		// If the monitor is not found, remove it from state (Terraform will recreate it)
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor, got error: %s", err))
		return
	}

	// Update model with current state - only set values that are non-zero/non-empty
	// to preserve null values for optional fields
	data.Name = types.StringValue(monitor.Name)
	data.Type = types.StringValue(monitor.Type)

	if monitor.URL != "" {
		data.URL = types.StringValue(monitor.URL)
	}
	if monitor.Hostname != "" {
		data.Hostname = types.StringValue(monitor.Hostname)
	}
	if monitor.Port != 0 {
		data.Port = types.Int64Value(int64(monitor.Port))
	}

	data.Interval = types.Int64Value(int64(monitor.Interval))
	data.Timeout = types.Int64Value(int64(monitor.Timeout))
	data.RetryInterval = types.Int64Value(int64(monitor.RetryInterval))
	data.ResendInterval = types.Int64Value(int64(monitor.ResendInterval))
	data.MaxRetries = types.Int64Value(int64(monitor.MaxRetries))
	data.MaxRedirects = types.Int64Value(int64(monitor.MaxRedirects))
	data.Active = types.BoolValue(monitor.Active)

	// Only set boolean fields if they were set in the config
	if !data.UpsideDown.IsNull() {
		data.UpsideDown = types.BoolValue(monitor.UpsideDown)
	}
	if !data.FollowRedirect.IsNull() {
		data.FollowRedirect = types.BoolValue(monitor.FollowRedirect)
	}
	if !data.IgnoreTLS.IsNull() {
		data.IgnoreTLS = types.BoolValue(monitor.IgnoreTLS)
	}

	data.HTTPMethod = types.StringValue(monitor.HTTPMethod)

	if monitor.Body != "" {
		data.Body = types.StringValue(monitor.Body)
	}
	if monitor.BasicAuthUser != "" {
		data.BasicAuthUser = types.StringValue(monitor.BasicAuthUser)
	}
	if monitor.BasicAuthPass != "" {
		data.BasicAuthPass = types.StringValue(monitor.BasicAuthPass)
	}

	// Convert accepted status codes to list
	if len(monitor.AcceptedStatusCodes) > 0 {
		listValue, diags := types.ListValueFrom(ctx, types.StringType, monitor.AcceptedStatusCodes)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.AcceptedStatusCodes = listValue
	}

	// Convert notification IDs to string list
	if len(monitor.NotificationIDList) > 0 {
		notificationIDs := make([]string, len(monitor.NotificationIDList))
		for i, id := range monitor.NotificationIDList {
			notificationIDs[i] = strconv.Itoa(id)
		}
		listValue, diags := types.ListValueFrom(ctx, types.StringType, notificationIDs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.NotificationIDList = listValue
	} else {
		// If the plan has an empty list (not null), preserve it as empty list
		if !data.NotificationIDList.IsNull() && !data.NotificationIDList.IsUnknown() {
			emptyList, diags := types.ListValueFrom(ctx, types.StringType, []string{})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			data.NotificationIDList = emptyList
		} else {
			data.NotificationIDList = types.ListNull(types.StringType)
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data MonitorResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to int
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor ID: %s", err))
		return
	}

	// Convert Terraform model to API model
	monitor := &Monitor{
		ID:             id,
		Name:           data.Name.ValueString(),
		Type:           data.Type.ValueString(),
		URL:            data.URL.ValueString(),
		Hostname:       data.Hostname.ValueString(),
		Port:           int(data.Port.ValueInt64()),
		Interval:       int(data.Interval.ValueInt64()),
		Timeout:        int(data.Timeout.ValueInt64()),
		RetryInterval:  int(data.RetryInterval.ValueInt64()),
		ResendInterval: int(data.ResendInterval.ValueInt64()),
		MaxRetries:     int(data.MaxRetries.ValueInt64()),
		UpsideDown:     data.UpsideDown.ValueBool(),
		MaxRedirects:   int(data.MaxRedirects.ValueInt64()),
		FollowRedirect: data.FollowRedirect.ValueBool(),
		Active:         data.Active.ValueBool(),
		IgnoreTLS:      data.IgnoreTLS.ValueBool(),
		HTTPMethod:     data.HTTPMethod.ValueString(),
		Body:           data.Body.ValueString(),
		BasicAuthUser:  data.BasicAuthUser.ValueString(),
		BasicAuthPass:  data.BasicAuthPass.ValueString(),
	}

	// Convert lists
	if !data.AcceptedStatusCodes.IsNull() {
		var statusCodes []string
		data.AcceptedStatusCodes.ElementsAs(ctx, &statusCodes, false)
		monitor.AcceptedStatusCodes = statusCodes
	}

	if !data.Tags.IsNull() {
		var tags []string
		data.Tags.ElementsAs(ctx, &tags, false)
		monitor.Tags = tags
	}

	if !data.NotificationIDList.IsNull() {
		var notificationIDs []string
		data.NotificationIDList.ElementsAs(ctx, &notificationIDs, false)
		for _, idStr := range notificationIDs {
			if id, err := strconv.Atoi(idStr); err == nil {
				monitor.NotificationIDList = append(monitor.NotificationIDList, id)
			}
		}
	}

	// Update monitor
	_, err = r.client.UpdateMonitor(monitor)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update monitor, got error: %s", err))
		return
	}

	// Don't read back from server - preserve plan values to avoid inconsistent state errors
	// The state should reflect what we sent to the API

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated a monitor resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MonitorResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to int
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor ID: %s", err))
		return
	}

	// Delete monitor
	err = r.client.DeleteMonitor(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete monitor, got error: %s", err))
		return
	}

	// Refresh the monitor cache to ensure the deleted monitor is removed
	err = r.client.RefreshMonitors()
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Unable to refresh monitors after delete: %s", err))
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "deleted a monitor resource")
}

func (r *MonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Validate that the ID is a valid integer
	_, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Monitor ID must be a valid integer, got: %s", req.ID))
		return
	}

	// Set the ID in the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
