package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MonitorDataSource{}

func NewMonitorDataSource() datasource.DataSource {
	return &MonitorDataSource{}
}

// MonitorDataSource defines the data source implementation.
type MonitorDataSource struct {
	client *Client
}

// MonitorDataSourceModel describes the data source data model.
type MonitorDataSourceModel struct {
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
	Active                types.Bool   `tfsdk:"active"`
	IgnoreTLS             types.Bool   `tfsdk:"ignore_tls"`
	HTTPMethod            types.String `tfsdk:"http_method"`
	Body                  types.String `tfsdk:"body"`
	BasicAuthUser         types.String `tfsdk:"basic_auth_user"`
}

func (d *MonitorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (d *MonitorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Uptime Kuma monitor data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Monitor identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Monitor name",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Monitor type (http, tcp, ping, etc.)",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL to monitor (for HTTP monitors)",
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname to monitor (for TCP/Ping monitors)",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port to monitor (for TCP monitors)",
				Computed:            true,
			},
			"interval": schema.Int64Attribute{
				MarkdownDescription: "Check interval in seconds",
				Computed:            true,
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "Timeout in seconds",
				Computed:            true,
			},
			"retry_interval": schema.Int64Attribute{
				MarkdownDescription: "Retry interval in seconds",
				Computed:            true,
			},
			"resend_interval": schema.Int64Attribute{
				MarkdownDescription: "Resend interval in seconds",
				Computed:            true,
			},
			"max_retries": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of retries",
				Computed:            true,
			},
			"upside_down": schema.BoolAttribute{
				MarkdownDescription: "Upside down mode (monitor expects failure)",
				Computed:            true,
			},
			"max_redirects": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of redirects to follow",
				Computed:            true,
			},
			"accepted_status_codes": schema.ListAttribute{
				MarkdownDescription: "List of accepted HTTP status codes",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"follow_redirect": schema.BoolAttribute{
				MarkdownDescription: "Follow redirects",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "List of tags",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the monitor is active",
				Computed:            true,
			},
			"ignore_tls": schema.BoolAttribute{
				MarkdownDescription: "Ignore TLS certificate errors",
				Computed:            true,
			},
			"http_method": schema.StringAttribute{
				MarkdownDescription: "HTTP method (GET, POST, etc.)",
				Computed:            true,
			},
			"body": schema.StringAttribute{
				MarkdownDescription: "Request body for POST/PUT requests",
				Computed:            true,
			},
			"basic_auth_user": schema.StringAttribute{
				MarkdownDescription: "Basic authentication username",
				Computed:            true,
			},
		},
	}
}

func (d *MonitorDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *MonitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MonitorDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert string ID to int
	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor ID: %s", err))
		return
	}

	// Get monitor from API
	monitor, err := d.client.GetMonitor(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor, got error: %s", err))
		return
	}

	// Map response body to model
	data.Name = types.StringValue(monitor.Name)
	data.Type = types.StringValue(monitor.Type)
	data.URL = types.StringValue(monitor.URL)
	data.Hostname = types.StringValue(monitor.Hostname)
	data.Port = types.Int64Value(int64(monitor.Port))
	data.Interval = types.Int64Value(int64(monitor.Interval))
	data.Timeout = types.Int64Value(int64(monitor.Timeout))
	data.RetryInterval = types.Int64Value(int64(monitor.RetryInterval))
	data.ResendInterval = types.Int64Value(int64(monitor.ResendInterval))
	data.MaxRetries = types.Int64Value(int64(monitor.MaxRetries))
	data.UpsideDown = types.BoolValue(monitor.UpsideDown)
	data.MaxRedirects = types.Int64Value(int64(monitor.MaxRedirects))
	data.FollowRedirect = types.BoolValue(monitor.FollowRedirect)
	data.Active = types.BoolValue(monitor.Active)
	data.IgnoreTLS = types.BoolValue(monitor.IgnoreTLS)
	data.HTTPMethod = types.StringValue(monitor.HTTPMethod)
	data.Body = types.StringValue(monitor.Body)
	data.BasicAuthUser = types.StringValue(monitor.BasicAuthUser)

	// Convert lists
	if len(monitor.AcceptedStatusCodes) > 0 {
		statusCodesList, diags := types.ListValueFrom(ctx, types.StringType, monitor.AcceptedStatusCodes)
		resp.Diagnostics.Append(diags...)
		data.AcceptedStatusCodes = statusCodesList
	}

	if len(monitor.Tags) > 0 {
		tagsList, diags := types.ListValueFrom(ctx, types.StringType, monitor.Tags)
		resp.Diagnostics.Append(diags...)
		data.Tags = tagsList
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
