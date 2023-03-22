package provider

import (
	"context"
	"fmt"

	sdk "github.com/paloaltonetworks/sase-go"
	sdkapi "github.com/paloaltonetworks/sase-go/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the provider implementation interface is sound.
var (
	_ provider.Provider = &SaseProvider{}
)

// SaseProvider is the provider implementation.
type SaseProvider struct {
	version string
}

// SaseProviderModel maps provider schema data to a Go type.
type SaseProviderModel struct {
	Host         types.String `tfsdk:"host"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Scope        types.String `tfsdk:"scope"`
	Logging      types.String `tfsdk:"logging"`
	AuthFile     types.String `tfsdk:"auth_file"`
}

// Metadata returns the provider type name.
func (p *SaseProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sase"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *SaseProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider to interact with Palo Alto Networks SASE API.",
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
					"The client ID for the connection.",
					"",
					"SASE_CLIENT_ID",
					"client_id",
				),
				Optional: true,
			},
			"client_secret": schema.StringAttribute{
				Description: ProviderParamDescription(
					"The client secret for the connection.",
					"",
					"SASE_CLIENT_SECRET",
					"client_secret",
				),
				Optional:  true,
				Sensitive: true,
			},
			"scope": schema.StringAttribute{
				Description: ProviderParamDescription(
					"The client scope.",
					"",
					"SASE_SCOPE",
					"scope",
				),
				Optional: true,
			},
			"logging": schema.StringAttribute{
				Description: ProviderParamDescription(
					"The logging level of the provider and the underlying communication.",
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

// Configure prepares the provider.
func (p *SaseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring the provider client")

	var config SaseProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	con := &sdk.Client{
		Host:             config.Host.ValueString(),
		ClientId:         config.ClientId.ValueString(),
		ClientSecret:     config.ClientSecret.ValueString(),
		Scope:            config.Scope.ValueString(),
		AuthFile:         config.AuthFile.ValueString(),
		Logging:          config.Logging.ValueString(),
		Logger:           tflog.Debug,
		CheckEnvironment: true,
		Agent:            fmt.Sprintf("Terraform/%s Provider/%s", req.TerraformVersion, p.version),
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

	tflog.Info(ctx, "Configured client", map[string]any{"success": true})
}

// DataSources defines the data sources for this provider.
func (p *SaseProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Section: netsec
		NewAntiSpywareProfilesDataSource,
		NewAntiSpywareProfilesListDataSource,
		NewAntiSpywareSignaturesListDataSource,
		NewAppOverrideRulesDataSource,
		NewAppOverrideRulesListDataSource,
		NewAuthenticationPortalsListDataSource,
		NewAuthenticationProfilesDataSource,
		NewAuthenticationProfilesListDataSource,
		NewAuthenticationRulesListDataSource,
		NewAuthenticationSequencesDataSource,
		NewAuthenticationSequencesListDataSource,
		NewAuthenticationSettingsListDataSource,
		NewAutoTagActionsListDataSource,
		NewBandwidthAllocationsListDataSource,
		NewBgpRoutingListDataSource,
		NewCertificateProfilesDataSource,
		NewCertificateProfilesListDataSource,
		NewCertificatesImportListDataSource,
		NewDecryptionExclusionsDataSource,
		NewDecryptionProfilesDataSource,
		NewDecryptionProfilesListDataSource,
		NewDecryptionRulesDataSource,
		NewDecryptionRulesListDataSource,
		NewDnsSecurityProfilesDataSource,
		NewDnsSecurityProfilesListDataSource,
		NewFileBlockingProfilesDataSource,
		NewFileBlockingProfilesListDataSource,
		NewHttpHeaderProfilesDataSource,
		NewHttpHeaderProfilesListDataSource,
		NewIkeCryptoProfilesDataSource,
		NewIkeCryptoProfilesListDataSource,
		NewIkeGatewaysDataSource,
		NewIkeGatewaysListDataSource,
		NewIpsecCryptoProfilesDataSource,
		NewIpsecCryptoProfilesListDataSource,
		NewIpsecTunnelsDataSource,
		NewIpsecTunnelsListDataSource,
		NewJobsDataSource,
		NewJobsListDataSource,
		NewKerberosServerProfilesDataSource,
		NewKerberosServerProfilesListDataSource,
		NewLdapServerProfilesDataSource,
		NewLdapServerProfilesListDataSource,
		NewLoadConfigDataSource,
		NewLocalUsersDataSource,
		NewLocalUsersListDataSource,
		NewMfaServersDataSource,
		NewMobileAgentInfrastructureSettingsListDataSource,
		NewMobileAgentLocationsListDataSource,
		NewObjectsAddressGroupsDataSource,
		NewObjectsAddressGroupsListDataSource,
		NewObjectsAddressesDataSource,
		NewObjectsAddressesListDataSource,
		NewObjectsApplicationFiltersDataSource,
		NewObjectsApplicationFiltersListDataSource,
		NewObjectsApplicationGroupsDataSource,
		NewObjectsApplicationGroupsListDataSource,
		NewObjectsApplicationsDataSource,
		NewObjectsApplicationsListDataSource,
		NewObjectsDynamicUserGroupsDataSource,
		NewObjectsDynamicUserGroupsListDataSource,
		NewObjectsExternalDynamicListsDataSource,
		NewObjectsExternalDynamicListsListDataSource,
		NewObjectsHipObjectsDataSource,
		NewObjectsHipObjectsListDataSource,
		NewObjectsHipProfilesDataSource,
		NewObjectsHipProfilesListDataSource,
		NewObjectsRegionsDataSource,
		NewObjectsRegionsListDataSource,
		NewObjectsSchedulesDataSource,
		NewObjectsSchedulesListDataSource,
		NewObjectsServiceGroupsDataSource,
		NewObjectsServiceGroupsListDataSource,
		NewObjectsServicesDataSource,
		NewObjectsServicesListDataSource,
		NewObjectsTagsDataSource,
		NewObjectsTagsListDataSource,
		NewOcspResponderDataSource,
		NewOcspResponderListDataSource,
		NewProfileGroupsDataSource,
		NewProfileGroupsListDataSource,
		NewQosPolicyRulesDataSource,
		NewQosPolicyRulesListDataSource,
		NewQosProfilesDataSource,
		NewQosProfilesListDataSource,
		NewRadiusServerProfilesDataSource,
		NewRadiusServerProfilesListDataSource,
		NewRemoteNetworksDataSource,
		NewRemoteNetworksListDataSource,
		NewSamlServerProfilesDataSource,
		NewSamlServerProfilesListDataSource,
		NewScepProfilesDataSource,
		NewScepProfilesListDataSource,
		NewSecurityRulesDataSource,
		NewSecurityRulesListDataSource,
		NewServiceConnectionGroupsListDataSource,
		NewServiceConnectionsListDataSource,
		NewSharedInfrastructureSettingsListDataSource,
		NewTacacsServerProfilesDataSource,
		NewTacacsServerProfilesListDataSource,
		NewTlsServiceProfilesDataSource,
		NewTlsServiceProfilesListDataSource,
		NewTrafficSteeringRulesListDataSource,
		NewTrustedCertificateAuthoritiesListDataSource,
		NewUrlAccessProfilesDataSource,
		NewUrlAccessProfilesListDataSource,
		NewUrlCategoriesListDataSource,
		NewUrlFilteringCategoriesListDataSource,
		NewVulnerabilityProtectionProfilesDataSource,
		NewVulnerabilityProtectionProfilesListDataSource,
		NewVulnerabilityProtectionSignaturesDataSource,
		NewVulnerabilityProtectionSignaturesListDataSource,
		NewWildfireAntiVirusProfilesDataSource,
		NewWildfireAntiVirusProfilesListDataSource,
	}
}

// Resources defines the data sources for this provider.
func (p *SaseProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Section: netsec
		NewAntiSpywareProfilesResource,
		NewAppOverrideRulesResource,
		NewAuthenticationProfilesResource,
		NewAuthenticationSequencesResource,
		NewCertificateProfilesResource,
		NewDecryptionExclusionsResource,
		NewDecryptionProfilesResource,
		NewDecryptionRulesResource,
		NewDnsSecurityProfilesResource,
		NewFileBlockingProfilesResource,
		NewHttpHeaderProfilesResource,
		NewIkeCryptoProfilesResource,
		NewIkeGatewaysResource,
		NewIpsecCryptoProfilesResource,
		NewIpsecTunnelsResource,
		NewKerberosServerProfilesResource,
		NewLdapServerProfilesResource,
		NewLocalUsersResource,
		NewMfaServersResource,
		NewObjectsAddressGroupsResource,
		NewObjectsAddressesResource,
		NewObjectsApplicationFiltersResource,
		NewObjectsApplicationsResource,
		NewObjectsDynamicUserGroupsResource,
		NewObjectsExternalDynamicListsResource,
		NewObjectsHipObjectsResource,
		NewObjectsHipProfilesResource,
		NewObjectsRegionsResource,
		NewObjectsSchedulesResource,
		NewObjectsServiceGroupsResource,
		NewObjectsServicesResource,
		NewObjectsTagsResource,
		NewOcspResponderResource,
		NewProfileGroupsResource,
		NewQosPolicyRulesResource,
		NewQosProfilesResource,
		NewRadiusServerProfilesResource,
		NewRemoteNetworksResource,
		NewSamlServerProfilesResource,
		NewScepProfilesResource,
		NewSecurityRulesResource,
		NewTacacsServerProfilesResource,
		NewTlsServiceProfilesResource,
		NewUrlAccessProfilesResource,
		NewVulnerabilityProtectionProfilesResource,
		NewVulnerabilityProtectionSignaturesResource,
		NewWildfireAntiVirusProfilesResource,
	}
}

// New is a helper function to get the provider implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SaseProvider{
			version: version,
		}
	}
}
