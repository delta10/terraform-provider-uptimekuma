package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure UptimeKumaProvider satisfies various provider interfaces.
var _ provider.Provider = &UptimeKumaProvider{}

// UptimeKumaProvider defines the provider implementation.
type UptimeKumaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// UptimeKumaProviderModel describes the provider data model.
type UptimeKumaProviderModel struct {
	URL      types.String `tfsdk:"url"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *UptimeKumaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "uptimekuma"
	resp.Version = p.version
}

func (p *UptimeKumaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the Uptime Kuma instance",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for Uptime Kuma authentication",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for Uptime Kuma authentication",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *UptimeKumaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data UptimeKumaProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available:
	// data.URL, data.Username, data.Password

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if data.URL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown Uptime Kuma URL",
			"The provider cannot create the Uptime Kuma API client as there is an unknown configuration value for the Uptime Kuma URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the UPTIMEKUMA_URL environment variable.",
		)
	}

	if data.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Uptime Kuma Username",
			"The provider cannot create the Uptime Kuma API client as there is an unknown configuration value for the Uptime Kuma username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the UPTIMEKUMA_USERNAME environment variable.",
		)
	}

	if data.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Uptime Kuma Password",
			"The provider cannot create the Uptime Kuma API client as there is an unknown configuration value for the Uptime Kuma password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the UPTIMEKUMA_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	url := os.Getenv("UPTIMEKUMA_URL")
	username := os.Getenv("UPTIMEKUMA_USERNAME")
	password := os.Getenv("UPTIMEKUMA_PASSWORD")

	if !data.URL.IsNull() {
		url = data.URL.ValueString()
	}

	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Missing Uptime Kuma URL",
			"The provider requires a URL for the Uptime Kuma instance. "+
				"Set the url value in the configuration or use the UPTIMEKUMA_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Uptime Kuma Username",
			"The provider requires a username for Uptime Kuma authentication. "+
				"Set the username value in the configuration or use the UPTIMEKUMA_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Uptime Kuma Password",
			"The provider requires a password for Uptime Kuma authentication. "+
				"Set the password value in the configuration or use the UPTIMEKUMA_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "uptimekuma_url", url)
	ctx = tflog.SetField(ctx, "uptimekuma_username", username)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "uptimekuma_password")

	tflog.Debug(ctx, "Creating Uptime Kuma client")

	// Create a new Uptime Kuma client using the configuration values
	client, err := NewClient(url, username, password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Uptime Kuma API Client",
			"An unexpected error occurred when creating the Uptime Kuma API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Uptime Kuma Client Error: "+err.Error(),
		)
		return
	}

	// Make the Uptime Kuma client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Uptime Kuma client", map[string]any{"success": true})
}

func (p *UptimeKumaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewMonitorResource,
		NewNotificationResource,
	}
}

func (p *UptimeKumaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMonitorDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &UptimeKumaProvider{
			version: version,
		}
	}
}
