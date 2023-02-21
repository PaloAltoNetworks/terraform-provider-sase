package provider

import (
	"context"
    "fmt"

    sdk "github.com/paloaltonetworks/sase-go"
    sdkapi "github.com/paloaltonetworks/sase-go/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	_ "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &SaseProvider{}
)

// SaseProvider is the provider implementation.
type SaseProvider struct{
    version string
}

// SaseProviderModel maps provider schema data to a Go type.
type SaseProviderModel struct {
    Host types.String `tfsdk:"host"`
    ClientId types.String `tfsdk:"client_id"`
    ClientSecret types.String `tfsdk:"client_secret"`
    Scope types.String `tfsdk:"scope"`
    AuthFile types.String `tfsdk:"auth_file"`
    Logging types.String `tfsdk:"logging"`
}

// Metadata returns the provider type name.
func (p *SaseProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sase"
    resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *SaseProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with HashiCups.",
		Attributes: map[string]schema.Attribute{
            "host": schema.StringAttribute{
                Description: ProviderParamDescription(
                    "The hostname.",
                    "api.sase.paloaltonetworks.com",
                    "SASE_HOST",
                    "host",
                ),
                Optional: true,
            },
            "client_id": schema.StringAttribute{
                Description: ProviderParamDescription(
                    "The CLient ID for the connection.",
                    "",
                    "SASE_CLIENT_ID",
                    "client_id",
                ),
                Required: true,
            },
            "client_secret": schema.StringAttribute{
                Description: ProviderParamDescription(
                    "The client secret for the connection.",
                    "",
                    "SASE_CLIENT_SECRET",
                    "client_secret",
                ),
                Required: true,
                Sensitive: true,
            },
            "scope": schema.StringAttribute{
                Description: ProviderParamDescription(
                    "The client scope.",
                    "",
                    "SASE_SCOPE",
                    "scope",
                ),
                Required: true,
            },
            "logging": schema.StringAttribute{
                Description: ProviderParamDescription(
                    "The logging level of the provider and underlying communication.",
                    sdkapi.LogQuiet,
                    "SASE_LOGGING",
                    "logging",
                ),
                Optional: true,
            },
            "auth_file": schema.StringAttribute{
                Description: ProviderParamDescription(
                    "The file path to the JSON file with auth creds for SASE.",
                    "",
                    "",
                    "",
                ),
                Optional: true,
            },
		},
	}
}

// Configure prepares a HashiCups API client for data sources and resources.
func (p *SaseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring SASE client")

	var config SaseProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

    con := &sdk.Client{
        Host: config.Host.ValueString(),
        ClientId: config.ClientId.ValueString(),
        ClientSecret: config.ClientSecret.ValueString(),
        Scope: config.Scope.ValueString(),
        AuthFile: config.AuthFile.ValueString(),
        Logging: config.Logging.ValueString(),
        Logger: tflog.Debug,
        CheckEnvironment: true,
        Agent: fmt.Sprintf("Terraform/%s Provider/%s", req.TerraformVersion, p.version),
    }

    if err := con.Setup(); err != nil {
        resp.Diagnostics.AddError("Provider parameter value error", err.Error())
        return
    }

    con.HttpClient.Transport = sdkapi.NewTransport(con.HttpClient.Transport, con)

    if err := con.RefreshJwt(ctx); err != nil {
        resp.Diagnostics.AddError("Authentication error", err.Error())
        return
    }

	resp.DataSourceData = con
	resp.ResourceData = con

	tflog.Info(ctx, "Configured SASE client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *SaseProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
        NewAddressObjectListDataSource,
        //NewAddressObjectDataSource,
		//NewCoffeesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *SaseProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
        NewAddressObjectResource,
		//NewOrderResource,
	}
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SaseProvider{
			version: version,
		}
	}
}
