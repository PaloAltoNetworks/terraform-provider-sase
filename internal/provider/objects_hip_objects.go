package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	dJpBrWV "github.com/paloaltonetworks/sase-go/netsec/schema/objects/hip/objects"
	yCYVNEN "github.com/paloaltonetworks/sase-go/netsec/service/v1/hipobjects"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data source (listing).
var (
	_ datasource.DataSource              = &objectsHipObjectsListDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsHipObjectsListDataSource{}
)

func NewObjectsHipObjectsListDataSource() datasource.DataSource {
	return &objectsHipObjectsListDataSource{}
}

type objectsHipObjectsListDataSource struct {
	client *sase.Client
}

type objectsHipObjectsListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Name   types.String `tfsdk:"name"`
	Folder types.String `tfsdk:"folder"`

	// Output.
	Data []objectsHipObjectsListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type objectsHipObjectsListDsModelConfig struct {
	AntiMalware        *objectsHipObjectsListDsModelAntiMalwareObject        `tfsdk:"anti_malware"`
	Certificate        *objectsHipObjectsListDsModelCertificateObject        `tfsdk:"certificate"`
	CustomChecks       *objectsHipObjectsListDsModelCustomChecksObject       `tfsdk:"custom_checks"`
	DataLossPrevention *objectsHipObjectsListDsModelDataLossPreventionObject `tfsdk:"data_loss_prevention"`
	Description        types.String                                          `tfsdk:"description"`
	DiskBackup         *objectsHipObjectsListDsModelDiskBackupObject         `tfsdk:"disk_backup"`
	DiskEncryption     *objectsHipObjectsListDsModelDiskEncryptionObject     `tfsdk:"disk_encryption"`
	Firewall           *objectsHipObjectsListDsModelFirewallObject           `tfsdk:"firewall"`
	HostInfo           *objectsHipObjectsListDsModelHostInfoObject           `tfsdk:"host_info"`
	ObjectId           types.String                                          `tfsdk:"object_id"`
	MobileDevice       *objectsHipObjectsListDsModelMobileDeviceObject       `tfsdk:"mobile_device"`
	Name               types.String                                          `tfsdk:"name"`
	NetworkInfo        *objectsHipObjectsListDsModelNetworkInfoObject        `tfsdk:"network_info"`
	PatchManagement    *objectsHipObjectsListDsModelPatchManagementObject    `tfsdk:"patch_management"`
}

type objectsHipObjectsListDsModelAntiMalwareObject struct {
	Criteria      *objectsHipObjectsListDsModelCriteriaObject `tfsdk:"criteria"`
	ExcludeVendor types.Bool                                  `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsListDsModelVendorObject  `tfsdk:"vendor"`
}

type objectsHipObjectsListDsModelCriteriaObject struct {
	IsInstalled        types.Bool                                        `tfsdk:"is_installed"`
	LastScanTime       *objectsHipObjectsListDsModelLastScanTimeObject   `tfsdk:"last_scan_time"`
	ProductVersion     *objectsHipObjectsListDsModelProductVersionObject `tfsdk:"product_version"`
	RealTimeProtection types.String                                      `tfsdk:"real_time_protection"`
	VirdefVersion      *objectsHipObjectsListDsModelVirdefVersionObject  `tfsdk:"virdef_version"`
}

type objectsHipObjectsListDsModelLastScanTimeObject struct {
	NotAvailable types.Bool                                   `tfsdk:"not_available"`
	NotWithin    *objectsHipObjectsListDsModelNotWithinObject `tfsdk:"not_within"`
	Within       *objectsHipObjectsListDsModelWithinObject    `tfsdk:"within"`
}

type objectsHipObjectsListDsModelNotWithinObject struct {
	Days  types.Int64 `tfsdk:"days"`
	Hours types.Int64 `tfsdk:"hours"`
}

type objectsHipObjectsListDsModelWithinObject struct {
	Days  types.Int64 `tfsdk:"days"`
	Hours types.Int64 `tfsdk:"hours"`
}

type objectsHipObjectsListDsModelProductVersionObject struct {
	Contains     types.String                                  `tfsdk:"contains"`
	GreaterEqual types.String                                  `tfsdk:"greater_equal"`
	GreaterThan  types.String                                  `tfsdk:"greater_than"`
	Is           types.String                                  `tfsdk:"is"`
	IsNot        types.String                                  `tfsdk:"is_not"`
	LessEqual    types.String                                  `tfsdk:"less_equal"`
	LessThan     types.String                                  `tfsdk:"less_than"`
	NotWithin    *objectsHipObjectsListDsModelNotWithinObject1 `tfsdk:"not_within"`
	Within       *objectsHipObjectsListDsModelWithinObject1    `tfsdk:"within"`
}

type objectsHipObjectsListDsModelNotWithinObject1 struct {
	Versions types.Int64 `tfsdk:"versions"`
}

type objectsHipObjectsListDsModelWithinObject1 struct {
	Versions types.Int64 `tfsdk:"versions"`
}

type objectsHipObjectsListDsModelVirdefVersionObject struct {
	NotWithin *objectsHipObjectsListDsModelNotWithinObject2 `tfsdk:"not_within"`
	Within    *objectsHipObjectsListDsModelWithinObject2    `tfsdk:"within"`
}

type objectsHipObjectsListDsModelNotWithinObject2 struct {
	Days     types.Int64 `tfsdk:"days"`
	Versions types.Int64 `tfsdk:"versions"`
}

type objectsHipObjectsListDsModelWithinObject2 struct {
	Days     types.Int64 `tfsdk:"days"`
	Versions types.Int64 `tfsdk:"versions"`
}

type objectsHipObjectsListDsModelVendorObject struct {
	Name    types.String   `tfsdk:"name"`
	Product []types.String `tfsdk:"product"`
}

type objectsHipObjectsListDsModelCertificateObject struct {
	Criteria *objectsHipObjectsListDsModelCriteriaObject1 `tfsdk:"criteria"`
}

type objectsHipObjectsListDsModelCriteriaObject1 struct {
	CertificateAttributes []objectsHipObjectsListDsModelCertificateAttributesObject `tfsdk:"certificate_attributes"`
	CertificateProfile    types.String                                              `tfsdk:"certificate_profile"`
}

type objectsHipObjectsListDsModelCertificateAttributesObject struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type objectsHipObjectsListDsModelCustomChecksObject struct {
	Criteria objectsHipObjectsListDsModelCriteriaObject2 `tfsdk:"criteria"`
}

type objectsHipObjectsListDsModelCriteriaObject2 struct {
	Plist       []objectsHipObjectsListDsModelPlistObject       `tfsdk:"plist"`
	ProcessList []objectsHipObjectsListDsModelProcessListObject `tfsdk:"process_list"`
	RegistryKey []objectsHipObjectsListDsModelRegistryKeyObject `tfsdk:"registry_key"`
}

type objectsHipObjectsListDsModelPlistObject struct {
	Key    []objectsHipObjectsListDsModelKeyObject `tfsdk:"key"`
	Name   types.String                            `tfsdk:"name"`
	Negate types.Bool                              `tfsdk:"negate"`
}

type objectsHipObjectsListDsModelKeyObject struct {
	Name   types.String `tfsdk:"name"`
	Negate types.Bool   `tfsdk:"negate"`
	Value  types.String `tfsdk:"value"`
}

type objectsHipObjectsListDsModelProcessListObject struct {
	Name    types.String `tfsdk:"name"`
	Running types.Bool   `tfsdk:"running"`
}

type objectsHipObjectsListDsModelRegistryKeyObject struct {
	DefaultValueData types.String                                      `tfsdk:"default_value_data"`
	Name             types.String                                      `tfsdk:"name"`
	Negate           types.Bool                                        `tfsdk:"negate"`
	RegistryValue    []objectsHipObjectsListDsModelRegistryValueObject `tfsdk:"registry_value"`
}

type objectsHipObjectsListDsModelRegistryValueObject struct {
	Name      types.String `tfsdk:"name"`
	Negate    types.Bool   `tfsdk:"negate"`
	ValueData types.String `tfsdk:"value_data"`
}

type objectsHipObjectsListDsModelDataLossPreventionObject struct {
	Criteria      *objectsHipObjectsListDsModelCriteriaObject3 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                                   `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsListDsModelVendorObject1  `tfsdk:"vendor"`
}

type objectsHipObjectsListDsModelCriteriaObject3 struct {
	IsEnabled   types.String `tfsdk:"is_enabled"`
	IsInstalled types.Bool   `tfsdk:"is_installed"`
}

type objectsHipObjectsListDsModelVendorObject1 struct {
	Name    types.String   `tfsdk:"name"`
	Product []types.String `tfsdk:"product"`
}

type objectsHipObjectsListDsModelDiskBackupObject struct {
	Criteria      *objectsHipObjectsListDsModelCriteriaObject4 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                                   `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsListDsModelVendorObject   `tfsdk:"vendor"`
}

type objectsHipObjectsListDsModelCriteriaObject4 struct {
	IsInstalled    types.Bool                                        `tfsdk:"is_installed"`
	LastBackupTime *objectsHipObjectsListDsModelLastBackupTimeObject `tfsdk:"last_backup_time"`
}

type objectsHipObjectsListDsModelLastBackupTimeObject struct {
	NotAvailable types.Bool                                   `tfsdk:"not_available"`
	NotWithin    *objectsHipObjectsListDsModelNotWithinObject `tfsdk:"not_within"`
	Within       *objectsHipObjectsListDsModelWithinObject    `tfsdk:"within"`
}

type objectsHipObjectsListDsModelDiskEncryptionObject struct {
	Criteria      *objectsHipObjectsListDsModelCriteriaObject5 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                                   `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsListDsModelVendorObject   `tfsdk:"vendor"`
}

type objectsHipObjectsListDsModelCriteriaObject5 struct {
	EncryptedLocations []objectsHipObjectsListDsModelEncryptedLocationsObject `tfsdk:"encrypted_locations"`
	IsInstalled        types.Bool                                             `tfsdk:"is_installed"`
}

type objectsHipObjectsListDsModelEncryptedLocationsObject struct {
	EncryptionState *objectsHipObjectsListDsModelEncryptionStateObject `tfsdk:"encryption_state"`
	Name            types.String                                       `tfsdk:"name"`
}

type objectsHipObjectsListDsModelEncryptionStateObject struct {
	Is    types.String `tfsdk:"is"`
	IsNot types.String `tfsdk:"is_not"`
}

type objectsHipObjectsListDsModelFirewallObject struct {
	Criteria      *objectsHipObjectsListDsModelCriteriaObject3 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                                   `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsListDsModelVendorObject   `tfsdk:"vendor"`
}

type objectsHipObjectsListDsModelHostInfoObject struct {
	Criteria objectsHipObjectsListDsModelCriteriaObject6 `tfsdk:"criteria"`
}

type objectsHipObjectsListDsModelCriteriaObject6 struct {
	ClientVersion *objectsHipObjectsListDsModelClientVersionObject `tfsdk:"client_version"`
	Domain        *objectsHipObjectsListDsModelDomainObject        `tfsdk:"domain"`
	HostId        *objectsHipObjectsListDsModelHostIdObject        `tfsdk:"host_id"`
	HostName      *objectsHipObjectsListDsModelHostNameObject      `tfsdk:"host_name"`
	Managed       types.Bool                                       `tfsdk:"managed"`
	Os            *objectsHipObjectsListDsModelOsObject            `tfsdk:"os"`
	SerialNumber  *objectsHipObjectsListDsModelSerialNumberObject  `tfsdk:"serial_number"`
}

type objectsHipObjectsListDsModelClientVersionObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsListDsModelDomainObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsListDsModelHostIdObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsListDsModelHostNameObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsListDsModelOsObject struct {
	Contains *objectsHipObjectsListDsModelContainsObject `tfsdk:"contains"`
}

type objectsHipObjectsListDsModelContainsObject struct {
	Apple     types.String `tfsdk:"apple"`
	Google    types.String `tfsdk:"google"`
	Linux     types.String `tfsdk:"linux"`
	Microsoft types.String `tfsdk:"microsoft"`
	Other     types.String `tfsdk:"other"`
}

type objectsHipObjectsListDsModelSerialNumberObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsListDsModelMobileDeviceObject struct {
	Criteria *objectsHipObjectsListDsModelCriteriaObject7 `tfsdk:"criteria"`
}

type objectsHipObjectsListDsModelCriteriaObject7 struct {
	Applications    *objectsHipObjectsListDsModelApplicationsObject    `tfsdk:"applications"`
	DiskEncrypted   types.Bool                                         `tfsdk:"disk_encrypted"`
	Imei            *objectsHipObjectsListDsModelImeiObject            `tfsdk:"imei"`
	Jailbroken      types.Bool                                         `tfsdk:"jailbroken"`
	LastCheckinTime *objectsHipObjectsListDsModelLastCheckinTimeObject `tfsdk:"last_checkin_time"`
	Model           *objectsHipObjectsListDsModelModelObject           `tfsdk:"model"`
	PasscodeSet     types.Bool                                         `tfsdk:"passcode_set"`
	PhoneNumber     *objectsHipObjectsListDsModelPhoneNumberObject     `tfsdk:"phone_number"`
	Tag             *objectsHipObjectsListDsModelTagObject             `tfsdk:"tag"`
}

type objectsHipObjectsListDsModelApplicationsObject struct {
	HasMalware      *objectsHipObjectsListDsModelHasMalwareObject `tfsdk:"has_malware"`
	HasUnmanagedApp types.Bool                                    `tfsdk:"has_unmanaged_app"`
	Includes        []objectsHipObjectsListDsModelIncludesObject  `tfsdk:"includes"`
}

type objectsHipObjectsListDsModelHasMalwareObject struct {
	No  types.Bool                             `tfsdk:"no"`
	Yes *objectsHipObjectsListDsModelYesObject `tfsdk:"yes"`
}

type objectsHipObjectsListDsModelYesObject struct {
	Excludes []objectsHipObjectsListDsModelExcludesObject `tfsdk:"excludes"`
}

type objectsHipObjectsListDsModelExcludesObject struct {
	Hash    types.String `tfsdk:"hash"`
	Name    types.String `tfsdk:"name"`
	Package types.String `tfsdk:"package"`
}

type objectsHipObjectsListDsModelIncludesObject struct {
	Hash    types.String `tfsdk:"hash"`
	Name    types.String `tfsdk:"name"`
	Package types.String `tfsdk:"package"`
}

type objectsHipObjectsListDsModelImeiObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsListDsModelLastCheckinTimeObject struct {
	NotWithin *objectsHipObjectsListDsModelNotWithinObject3 `tfsdk:"not_within"`
	Within    *objectsHipObjectsListDsModelWithinObject3    `tfsdk:"within"`
}

type objectsHipObjectsListDsModelNotWithinObject3 struct {
	Days types.Int64 `tfsdk:"days"`
}

type objectsHipObjectsListDsModelWithinObject3 struct {
	Days types.Int64 `tfsdk:"days"`
}

type objectsHipObjectsListDsModelModelObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsListDsModelPhoneNumberObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsListDsModelTagObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsListDsModelNetworkInfoObject struct {
	Criteria *objectsHipObjectsListDsModelCriteriaObject8 `tfsdk:"criteria"`
}

type objectsHipObjectsListDsModelCriteriaObject8 struct {
	Network *objectsHipObjectsListDsModelNetworkObject `tfsdk:"network"`
}

type objectsHipObjectsListDsModelNetworkObject struct {
	Is    *objectsHipObjectsListDsModelIsObject    `tfsdk:"is"`
	IsNot *objectsHipObjectsListDsModelIsNotObject `tfsdk:"is_not"`
}

type objectsHipObjectsListDsModelIsObject struct {
	Mobile  *objectsHipObjectsListDsModelMobileObject `tfsdk:"mobile"`
	Unknown types.Bool                                `tfsdk:"unknown"`
	Wifi    *objectsHipObjectsListDsModelWifiObject   `tfsdk:"wifi"`
}

type objectsHipObjectsListDsModelMobileObject struct {
	Carrier types.String `tfsdk:"carrier"`
}

type objectsHipObjectsListDsModelWifiObject struct {
	Ssid types.String `tfsdk:"ssid"`
}

type objectsHipObjectsListDsModelIsNotObject struct {
	Ethernet types.Bool                                `tfsdk:"ethernet"`
	Mobile   *objectsHipObjectsListDsModelMobileObject `tfsdk:"mobile"`
	Unknown  types.Bool                                `tfsdk:"unknown"`
	Wifi     *objectsHipObjectsListDsModelWifiObject   `tfsdk:"wifi"`
}

type objectsHipObjectsListDsModelPatchManagementObject struct {
	Criteria      *objectsHipObjectsListDsModelCriteriaObject9 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                                   `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsListDsModelVendorObject1  `tfsdk:"vendor"`
}

type objectsHipObjectsListDsModelCriteriaObject9 struct {
	IsEnabled      types.String                                      `tfsdk:"is_enabled"`
	IsInstalled    types.Bool                                        `tfsdk:"is_installed"`
	MissingPatches *objectsHipObjectsListDsModelMissingPatchesObject `tfsdk:"missing_patches"`
}

type objectsHipObjectsListDsModelMissingPatchesObject struct {
	Check    types.String                                `tfsdk:"check"`
	Patches  []types.String                              `tfsdk:"patches"`
	Severity *objectsHipObjectsListDsModelSeverityObject `tfsdk:"severity"`
}

type objectsHipObjectsListDsModelSeverityObject struct {
	GreaterEqual types.Int64 `tfsdk:"greater_equal"`
	GreaterThan  types.Int64 `tfsdk:"greater_than"`
	Is           types.Int64 `tfsdk:"is"`
	IsNot        types.Int64 `tfsdk:"is_not"`
	LessEqual    types.Int64 `tfsdk:"less_equal"`
	LessThan     types.Int64 `tfsdk:"less_than"`
}

// Metadata returns the data source type name.
func (d *objectsHipObjectsListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_hip_objects_list"
}

// Schema defines the schema for this listing data source.
func (d *objectsHipObjectsListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves a listing of config items.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description: "The object ID.",
				Computed:    true,
			},

			// Input.
			"limit": dsschema.Int64Attribute{
				Description: "The max count in result entry (count per page)",
				Optional:    true,
				Computed:    true,
			},
			"offset": dsschema.Int64Attribute{
				Description: "The offset of the result entry",
				Optional:    true,
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "The name of the entry",
				Optional:    true,
			},
			"folder": dsschema.StringAttribute{
				Description: "The folder of the entry",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"anti_malware": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"criteria": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"is_installed": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"last_scan_time": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"not_available": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"not_within": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"days": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
														"hours": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"within": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"days": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
														"hours": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
											},
										},
										"product_version": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"contains": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"greater_equal": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"greater_than": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is_not": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"less_equal": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"less_than": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"not_within": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"versions": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"within": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"versions": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
											},
										},
										"real_time_protection": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"virdef_version": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"not_within": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"days": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
														"versions": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"within": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"days": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
														"versions": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
											},
										},
									},
								},
								"exclude_vendor": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"vendor": dsschema.ListNestedAttribute{
									Description: "",
									Computed:    true,
									NestedObject: dsschema.NestedAttributeObject{
										Attributes: map[string]dsschema.Attribute{
											"name": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"product": dsschema.ListAttribute{
												Description: "",
												Computed:    true,
												ElementType: types.StringType,
											},
										},
									},
								},
							},
						},
						"certificate": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"criteria": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"certificate_attributes": dsschema.ListNestedAttribute{
											Description: "",
											Computed:    true,
											NestedObject: dsschema.NestedAttributeObject{
												Attributes: map[string]dsschema.Attribute{
													"name": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"value": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
												},
											},
										},
										"certificate_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
							},
						},
						"custom_checks": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"criteria": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"plist": dsschema.ListNestedAttribute{
											Description: "",
											Computed:    true,
											NestedObject: dsschema.NestedAttributeObject{
												Attributes: map[string]dsschema.Attribute{
													"key": dsschema.ListNestedAttribute{
														Description: "",
														Computed:    true,
														NestedObject: dsschema.NestedAttributeObject{
															Attributes: map[string]dsschema.Attribute{
																"name": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
																"negate": dsschema.BoolAttribute{
																	Description: "",
																	Computed:    true,
																},
																"value": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
															},
														},
													},
													"name": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"negate": dsschema.BoolAttribute{
														Description: "",
														Computed:    true,
													},
												},
											},
										},
										"process_list": dsschema.ListNestedAttribute{
											Description: "",
											Computed:    true,
											NestedObject: dsschema.NestedAttributeObject{
												Attributes: map[string]dsschema.Attribute{
													"name": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"running": dsschema.BoolAttribute{
														Description: "",
														Computed:    true,
													},
												},
											},
										},
										"registry_key": dsschema.ListNestedAttribute{
											Description: "",
											Computed:    true,
											NestedObject: dsschema.NestedAttributeObject{
												Attributes: map[string]dsschema.Attribute{
													"default_value_data": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"name": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"negate": dsschema.BoolAttribute{
														Description: "",
														Computed:    true,
													},
													"registry_value": dsschema.ListNestedAttribute{
														Description: "",
														Computed:    true,
														NestedObject: dsschema.NestedAttributeObject{
															Attributes: map[string]dsschema.Attribute{
																"name": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
																"negate": dsschema.BoolAttribute{
																	Description: "",
																	Computed:    true,
																},
																"value_data": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"data_loss_prevention": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"criteria": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"is_enabled": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"is_installed": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"exclude_vendor": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"vendor": dsschema.ListNestedAttribute{
									Description: "",
									Computed:    true,
									NestedObject: dsschema.NestedAttributeObject{
										Attributes: map[string]dsschema.Attribute{
											"name": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"product": dsschema.ListAttribute{
												Description: "",
												Computed:    true,
												ElementType: types.StringType,
											},
										},
									},
								},
							},
						},
						"description": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"disk_backup": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"criteria": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"is_installed": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"last_backup_time": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"not_available": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"not_within": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"days": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
														"hours": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"within": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"days": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
														"hours": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
											},
										},
									},
								},
								"exclude_vendor": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"vendor": dsschema.ListNestedAttribute{
									Description: "",
									Computed:    true,
									NestedObject: dsschema.NestedAttributeObject{
										Attributes: map[string]dsschema.Attribute{
											"name": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"product": dsschema.ListAttribute{
												Description: "",
												Computed:    true,
												ElementType: types.StringType,
											},
										},
									},
								},
							},
						},
						"disk_encryption": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"criteria": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"encrypted_locations": dsschema.ListNestedAttribute{
											Description: "",
											Computed:    true,
											NestedObject: dsschema.NestedAttributeObject{
												Attributes: map[string]dsschema.Attribute{
													"encryption_state": dsschema.SingleNestedAttribute{
														Description: "",
														Computed:    true,
														Attributes: map[string]dsschema.Attribute{
															"is": dsschema.StringAttribute{
																Description: "",
																Computed:    true,
															},
															"is_not": dsschema.StringAttribute{
																Description: "",
																Computed:    true,
															},
														},
													},
													"name": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
												},
											},
										},
										"is_installed": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"exclude_vendor": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"vendor": dsschema.ListNestedAttribute{
									Description: "",
									Computed:    true,
									NestedObject: dsschema.NestedAttributeObject{
										Attributes: map[string]dsschema.Attribute{
											"name": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"product": dsschema.ListAttribute{
												Description: "",
												Computed:    true,
												ElementType: types.StringType,
											},
										},
									},
								},
							},
						},
						"firewall": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"criteria": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"is_enabled": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"is_installed": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"exclude_vendor": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"vendor": dsschema.ListNestedAttribute{
									Description: "",
									Computed:    true,
									NestedObject: dsschema.NestedAttributeObject{
										Attributes: map[string]dsschema.Attribute{
											"name": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"product": dsschema.ListAttribute{
												Description: "",
												Computed:    true,
												ElementType: types.StringType,
											},
										},
									},
								},
							},
						},
						"host_info": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"criteria": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"client_version": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"contains": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is_not": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"domain": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"contains": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is_not": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"host_id": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"contains": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is_not": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"host_name": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"contains": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is_not": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"managed": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"os": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"contains": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"apple": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"google": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"linux": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"microsoft": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
														"other": dsschema.StringAttribute{
															Description: "",
															Computed:    true,
														},
													},
												},
											},
										},
										"serial_number": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"contains": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is_not": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"mobile_device": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"criteria": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"applications": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"has_malware": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"no": dsschema.BoolAttribute{
															Description: "",
															Computed:    true,
														},
														"yes": dsschema.SingleNestedAttribute{
															Description: "",
															Computed:    true,
															Attributes: map[string]dsschema.Attribute{
																"excludes": dsschema.ListNestedAttribute{
																	Description: "",
																	Computed:    true,
																	NestedObject: dsschema.NestedAttributeObject{
																		Attributes: map[string]dsschema.Attribute{
																			"hash": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																			"name": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																			"package": dsschema.StringAttribute{
																				Description: "",
																				Computed:    true,
																			},
																		},
																	},
																},
															},
														},
													},
												},
												"has_unmanaged_app": dsschema.BoolAttribute{
													Description: "",
													Computed:    true,
												},
												"includes": dsschema.ListNestedAttribute{
													Description: "",
													Computed:    true,
													NestedObject: dsschema.NestedAttributeObject{
														Attributes: map[string]dsschema.Attribute{
															"hash": dsschema.StringAttribute{
																Description: "",
																Computed:    true,
															},
															"name": dsschema.StringAttribute{
																Description: "",
																Computed:    true,
															},
															"package": dsschema.StringAttribute{
																Description: "",
																Computed:    true,
															},
														},
													},
												},
											},
										},
										"disk_encrypted": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"imei": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"contains": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is_not": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"jailbroken": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"last_checkin_time": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"not_within": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"days": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
												"within": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"days": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
											},
										},
										"model": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"contains": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is_not": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"passcode_set": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"phone_number": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"contains": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is_not": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"tag": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"contains": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is_not": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"network_info": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"criteria": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"network": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"is": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"mobile": dsschema.SingleNestedAttribute{
															Description: "",
															Computed:    true,
															Attributes: map[string]dsschema.Attribute{
																"carrier": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
															},
														},
														"unknown": dsschema.BoolAttribute{
															Description: "",
															Computed:    true,
														},
														"wifi": dsschema.SingleNestedAttribute{
															Description: "",
															Computed:    true,
															Attributes: map[string]dsschema.Attribute{
																"ssid": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
															},
														},
													},
												},
												"is_not": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"ethernet": dsschema.BoolAttribute{
															Description: "",
															Computed:    true,
														},
														"mobile": dsschema.SingleNestedAttribute{
															Description: "",
															Computed:    true,
															Attributes: map[string]dsschema.Attribute{
																"carrier": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
															},
														},
														"unknown": dsschema.BoolAttribute{
															Description: "",
															Computed:    true,
														},
														"wifi": dsschema.SingleNestedAttribute{
															Description: "",
															Computed:    true,
															Attributes: map[string]dsschema.Attribute{
																"ssid": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"patch_management": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"criteria": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"is_enabled": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"is_installed": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"missing_patches": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"check": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"patches": dsschema.ListAttribute{
													Description: "",
													Computed:    true,
													ElementType: types.StringType,
												},
												"severity": dsschema.SingleNestedAttribute{
													Description: "",
													Computed:    true,
													Attributes: map[string]dsschema.Attribute{
														"greater_equal": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
														"greater_than": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
														"is": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
														"is_not": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
														"less_equal": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
														"less_than": dsschema.Int64Attribute{
															Description: "",
															Computed:    true,
														},
													},
												},
											},
										},
									},
								},
								"exclude_vendor": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"vendor": dsschema.ListNestedAttribute{
									Description: "",
									Computed:    true,
									NestedObject: dsschema.NestedAttributeObject{
										Attributes: map[string]dsschema.Attribute{
											"name": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"product": dsschema.ListAttribute{
												Description: "",
												Computed:    true,
												ElementType: types.StringType,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"total": dsschema.Int64Attribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *objectsHipObjectsListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsHipObjectsListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsHipObjectsListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_objects_hip_objects_list",
		"limit":            state.Limit.ValueInt64(),
		"has_limit":        !state.Limit.IsNull(),
		"offset":           state.Offset.ValueInt64(),
		"has_offset":       !state.Offset.IsNull(),
		"name":             state.Name.ValueString(),
		"has_name":         !state.Name.IsNull(),
		"folder":           state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := yCYVNEN.NewClient(d.client)
	input := yCYVNEN.ListInput{
		Folder: state.Folder.ValueString(),
	}
	if !state.Limit.IsNull() {
		input.Limit = api.Int(state.Limit.ValueInt64())
	}
	if !state.Offset.IsNull() {
		input.Offset = api.Int(state.Offset.ValueInt64())
	}
	if !state.Name.IsNull() {
		input.Name = api.String(state.Name.ValueString())
	}

	// Perform the operation.
	ans, err := svc.List(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error getting listing", err.Error())
		return
	}

	// Store the answer to state.
	state.Id = types.StringValue(strings.Join([]string{strconv.FormatInt(*input.Limit, 10), strconv.FormatInt(*input.Offset, 10), *input.Name, input.Folder}, IdSeparator))
	var var0 []objectsHipObjectsListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]objectsHipObjectsListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 objectsHipObjectsListDsModelConfig
			var var3 *objectsHipObjectsListDsModelAntiMalwareObject
			if var1.AntiMalware != nil {
				var3 = &objectsHipObjectsListDsModelAntiMalwareObject{}
				var var4 *objectsHipObjectsListDsModelCriteriaObject
				if var1.AntiMalware.Criteria != nil {
					var4 = &objectsHipObjectsListDsModelCriteriaObject{}
					var var5 *objectsHipObjectsListDsModelLastScanTimeObject
					if var1.AntiMalware.Criteria.LastScanTime != nil {
						var5 = &objectsHipObjectsListDsModelLastScanTimeObject{}
						var var6 *objectsHipObjectsListDsModelNotWithinObject
						if var1.AntiMalware.Criteria.LastScanTime.NotWithin != nil {
							var6 = &objectsHipObjectsListDsModelNotWithinObject{}
							var6.Days = types.Int64Value(var1.AntiMalware.Criteria.LastScanTime.NotWithin.Days)
							var6.Hours = types.Int64Value(var1.AntiMalware.Criteria.LastScanTime.NotWithin.Hours)
						}
						var var7 *objectsHipObjectsListDsModelWithinObject
						if var1.AntiMalware.Criteria.LastScanTime.Within != nil {
							var7 = &objectsHipObjectsListDsModelWithinObject{}
							var7.Days = types.Int64Value(var1.AntiMalware.Criteria.LastScanTime.Within.Days)
							var7.Hours = types.Int64Value(var1.AntiMalware.Criteria.LastScanTime.Within.Hours)
						}
						if var1.AntiMalware.Criteria.LastScanTime.NotAvailable != nil {
							var5.NotAvailable = types.BoolValue(true)
						}
						var5.NotWithin = var6
						var5.Within = var7
					}
					var var8 *objectsHipObjectsListDsModelProductVersionObject
					if var1.AntiMalware.Criteria.ProductVersion != nil {
						var8 = &objectsHipObjectsListDsModelProductVersionObject{}
						var var9 *objectsHipObjectsListDsModelNotWithinObject1
						if var1.AntiMalware.Criteria.ProductVersion.NotWithin != nil {
							var9 = &objectsHipObjectsListDsModelNotWithinObject1{}
							var9.Versions = types.Int64Value(var1.AntiMalware.Criteria.ProductVersion.NotWithin.Versions)
						}
						var var10 *objectsHipObjectsListDsModelWithinObject1
						if var1.AntiMalware.Criteria.ProductVersion.Within != nil {
							var10 = &objectsHipObjectsListDsModelWithinObject1{}
							var10.Versions = types.Int64Value(var1.AntiMalware.Criteria.ProductVersion.Within.Versions)
						}
						var8.Contains = types.StringValue(var1.AntiMalware.Criteria.ProductVersion.Contains)
						var8.GreaterEqual = types.StringValue(var1.AntiMalware.Criteria.ProductVersion.GreaterEqual)
						var8.GreaterThan = types.StringValue(var1.AntiMalware.Criteria.ProductVersion.GreaterThan)
						var8.Is = types.StringValue(var1.AntiMalware.Criteria.ProductVersion.Is)
						var8.IsNot = types.StringValue(var1.AntiMalware.Criteria.ProductVersion.IsNot)
						var8.LessEqual = types.StringValue(var1.AntiMalware.Criteria.ProductVersion.LessEqual)
						var8.LessThan = types.StringValue(var1.AntiMalware.Criteria.ProductVersion.LessThan)
						var8.NotWithin = var9
						var8.Within = var10
					}
					var var11 *objectsHipObjectsListDsModelVirdefVersionObject
					if var1.AntiMalware.Criteria.VirdefVersion != nil {
						var11 = &objectsHipObjectsListDsModelVirdefVersionObject{}
						var var12 *objectsHipObjectsListDsModelNotWithinObject2
						if var1.AntiMalware.Criteria.VirdefVersion.NotWithin != nil {
							var12 = &objectsHipObjectsListDsModelNotWithinObject2{}
							var12.Days = types.Int64Value(var1.AntiMalware.Criteria.VirdefVersion.NotWithin.Days)
							var12.Versions = types.Int64Value(var1.AntiMalware.Criteria.VirdefVersion.NotWithin.Versions)
						}
						var var13 *objectsHipObjectsListDsModelWithinObject2
						if var1.AntiMalware.Criteria.VirdefVersion.Within != nil {
							var13 = &objectsHipObjectsListDsModelWithinObject2{}
							var13.Days = types.Int64Value(var1.AntiMalware.Criteria.VirdefVersion.Within.Days)
							var13.Versions = types.Int64Value(var1.AntiMalware.Criteria.VirdefVersion.Within.Versions)
						}
						var11.NotWithin = var12
						var11.Within = var13
					}
					var4.IsInstalled = types.BoolValue(var1.AntiMalware.Criteria.IsInstalled)
					var4.LastScanTime = var5
					var4.ProductVersion = var8
					var4.RealTimeProtection = types.StringValue(var1.AntiMalware.Criteria.RealTimeProtection)
					var4.VirdefVersion = var11
				}
				var var14 []objectsHipObjectsListDsModelVendorObject
				if len(var1.AntiMalware.Vendor) != 0 {
					var14 = make([]objectsHipObjectsListDsModelVendorObject, 0, len(var1.AntiMalware.Vendor))
					for var15Index := range var1.AntiMalware.Vendor {
						var15 := var1.AntiMalware.Vendor[var15Index]
						var var16 objectsHipObjectsListDsModelVendorObject
						var16.Name = types.StringValue(var15.Name)
						var16.Product = EncodeStringSlice(var15.Product)
						var14 = append(var14, var16)
					}
				}
				var3.Criteria = var4
				var3.ExcludeVendor = types.BoolValue(var1.AntiMalware.ExcludeVendor)
				var3.Vendor = var14
			}
			var var17 *objectsHipObjectsListDsModelCertificateObject
			if var1.Certificate != nil {
				var17 = &objectsHipObjectsListDsModelCertificateObject{}
				var var18 *objectsHipObjectsListDsModelCriteriaObject1
				if var1.Certificate.Criteria != nil {
					var18 = &objectsHipObjectsListDsModelCriteriaObject1{}
					var var19 []objectsHipObjectsListDsModelCertificateAttributesObject
					if len(var1.Certificate.Criteria.CertificateAttributes) != 0 {
						var19 = make([]objectsHipObjectsListDsModelCertificateAttributesObject, 0, len(var1.Certificate.Criteria.CertificateAttributes))
						for var20Index := range var1.Certificate.Criteria.CertificateAttributes {
							var20 := var1.Certificate.Criteria.CertificateAttributes[var20Index]
							var var21 objectsHipObjectsListDsModelCertificateAttributesObject
							var21.Name = types.StringValue(var20.Name)
							var21.Value = types.StringValue(var20.Value)
							var19 = append(var19, var21)
						}
					}
					var18.CertificateAttributes = var19
					var18.CertificateProfile = types.StringValue(var1.Certificate.Criteria.CertificateProfile)
				}
				var17.Criteria = var18
			}
			var var22 *objectsHipObjectsListDsModelCustomChecksObject
			if var1.CustomChecks != nil {
				var22 = &objectsHipObjectsListDsModelCustomChecksObject{}
				var var23 objectsHipObjectsListDsModelCriteriaObject2
				var var24 []objectsHipObjectsListDsModelPlistObject
				if len(var1.CustomChecks.Criteria.Plist) != 0 {
					var24 = make([]objectsHipObjectsListDsModelPlistObject, 0, len(var1.CustomChecks.Criteria.Plist))
					for var25Index := range var1.CustomChecks.Criteria.Plist {
						var25 := var1.CustomChecks.Criteria.Plist[var25Index]
						var var26 objectsHipObjectsListDsModelPlistObject
						var var27 []objectsHipObjectsListDsModelKeyObject
						if len(var25.Key) != 0 {
							var27 = make([]objectsHipObjectsListDsModelKeyObject, 0, len(var25.Key))
							for var28Index := range var25.Key {
								var28 := var25.Key[var28Index]
								var var29 objectsHipObjectsListDsModelKeyObject
								var29.Name = types.StringValue(var28.Name)
								var29.Negate = types.BoolValue(var28.Negate)
								var29.Value = types.StringValue(var28.Value)
								var27 = append(var27, var29)
							}
						}
						var26.Key = var27
						var26.Name = types.StringValue(var25.Name)
						var26.Negate = types.BoolValue(var25.Negate)
						var24 = append(var24, var26)
					}
				}
				var var30 []objectsHipObjectsListDsModelProcessListObject
				if len(var1.CustomChecks.Criteria.ProcessList) != 0 {
					var30 = make([]objectsHipObjectsListDsModelProcessListObject, 0, len(var1.CustomChecks.Criteria.ProcessList))
					for var31Index := range var1.CustomChecks.Criteria.ProcessList {
						var31 := var1.CustomChecks.Criteria.ProcessList[var31Index]
						var var32 objectsHipObjectsListDsModelProcessListObject
						var32.Name = types.StringValue(var31.Name)
						var32.Running = types.BoolValue(var31.Running)
						var30 = append(var30, var32)
					}
				}
				var var33 []objectsHipObjectsListDsModelRegistryKeyObject
				if len(var1.CustomChecks.Criteria.RegistryKey) != 0 {
					var33 = make([]objectsHipObjectsListDsModelRegistryKeyObject, 0, len(var1.CustomChecks.Criteria.RegistryKey))
					for var34Index := range var1.CustomChecks.Criteria.RegistryKey {
						var34 := var1.CustomChecks.Criteria.RegistryKey[var34Index]
						var var35 objectsHipObjectsListDsModelRegistryKeyObject
						var var36 []objectsHipObjectsListDsModelRegistryValueObject
						if len(var34.RegistryValue) != 0 {
							var36 = make([]objectsHipObjectsListDsModelRegistryValueObject, 0, len(var34.RegistryValue))
							for var37Index := range var34.RegistryValue {
								var37 := var34.RegistryValue[var37Index]
								var var38 objectsHipObjectsListDsModelRegistryValueObject
								var38.Name = types.StringValue(var37.Name)
								var38.Negate = types.BoolValue(var37.Negate)
								var38.ValueData = types.StringValue(var37.ValueData)
								var36 = append(var36, var38)
							}
						}
						var35.DefaultValueData = types.StringValue(var34.DefaultValueData)
						var35.Name = types.StringValue(var34.Name)
						var35.Negate = types.BoolValue(var34.Negate)
						var35.RegistryValue = var36
						var33 = append(var33, var35)
					}
				}
				var23.Plist = var24
				var23.ProcessList = var30
				var23.RegistryKey = var33
				var22.Criteria = var23
			}
			var var39 *objectsHipObjectsListDsModelDataLossPreventionObject
			if var1.DataLossPrevention != nil {
				var39 = &objectsHipObjectsListDsModelDataLossPreventionObject{}
				var var40 *objectsHipObjectsListDsModelCriteriaObject3
				if var1.DataLossPrevention.Criteria != nil {
					var40 = &objectsHipObjectsListDsModelCriteriaObject3{}
					var40.IsEnabled = types.StringValue(var1.DataLossPrevention.Criteria.IsEnabled)
					var40.IsInstalled = types.BoolValue(var1.DataLossPrevention.Criteria.IsInstalled)
				}
				var var41 []objectsHipObjectsListDsModelVendorObject1
				if len(var1.DataLossPrevention.Vendor) != 0 {
					var41 = make([]objectsHipObjectsListDsModelVendorObject1, 0, len(var1.DataLossPrevention.Vendor))
					for var42Index := range var1.DataLossPrevention.Vendor {
						var42 := var1.DataLossPrevention.Vendor[var42Index]
						var var43 objectsHipObjectsListDsModelVendorObject1
						var43.Name = types.StringValue(var42.Name)
						var43.Product = EncodeStringSlice(var42.Product)
						var41 = append(var41, var43)
					}
				}
				var39.Criteria = var40
				var39.ExcludeVendor = types.BoolValue(var1.DataLossPrevention.ExcludeVendor)
				var39.Vendor = var41
			}
			var var44 *objectsHipObjectsListDsModelDiskBackupObject
			if var1.DiskBackup != nil {
				var44 = &objectsHipObjectsListDsModelDiskBackupObject{}
				var var45 *objectsHipObjectsListDsModelCriteriaObject4
				if var1.DiskBackup.Criteria != nil {
					var45 = &objectsHipObjectsListDsModelCriteriaObject4{}
					var var46 *objectsHipObjectsListDsModelLastBackupTimeObject
					if var1.DiskBackup.Criteria.LastBackupTime != nil {
						var46 = &objectsHipObjectsListDsModelLastBackupTimeObject{}
						var var47 *objectsHipObjectsListDsModelNotWithinObject
						if var1.DiskBackup.Criteria.LastBackupTime.NotWithin != nil {
							var47 = &objectsHipObjectsListDsModelNotWithinObject{}
							var47.Days = types.Int64Value(var1.DiskBackup.Criteria.LastBackupTime.NotWithin.Days)
							var47.Hours = types.Int64Value(var1.DiskBackup.Criteria.LastBackupTime.NotWithin.Hours)
						}
						var var48 *objectsHipObjectsListDsModelWithinObject
						if var1.DiskBackup.Criteria.LastBackupTime.Within != nil {
							var48 = &objectsHipObjectsListDsModelWithinObject{}
							var48.Days = types.Int64Value(var1.DiskBackup.Criteria.LastBackupTime.Within.Days)
							var48.Hours = types.Int64Value(var1.DiskBackup.Criteria.LastBackupTime.Within.Hours)
						}
						if var1.DiskBackup.Criteria.LastBackupTime.NotAvailable != nil {
							var46.NotAvailable = types.BoolValue(true)
						}
						var46.NotWithin = var47
						var46.Within = var48
					}
					var45.IsInstalled = types.BoolValue(var1.DiskBackup.Criteria.IsInstalled)
					var45.LastBackupTime = var46
				}
				var var49 []objectsHipObjectsListDsModelVendorObject
				if len(var1.DiskBackup.Vendor) != 0 {
					var49 = make([]objectsHipObjectsListDsModelVendorObject, 0, len(var1.DiskBackup.Vendor))
					for var50Index := range var1.DiskBackup.Vendor {
						var50 := var1.DiskBackup.Vendor[var50Index]
						var var51 objectsHipObjectsListDsModelVendorObject
						var51.Name = types.StringValue(var50.Name)
						var51.Product = EncodeStringSlice(var50.Product)
						var49 = append(var49, var51)
					}
				}
				var44.Criteria = var45
				var44.ExcludeVendor = types.BoolValue(var1.DiskBackup.ExcludeVendor)
				var44.Vendor = var49
			}
			var var52 *objectsHipObjectsListDsModelDiskEncryptionObject
			if var1.DiskEncryption != nil {
				var52 = &objectsHipObjectsListDsModelDiskEncryptionObject{}
				var var53 *objectsHipObjectsListDsModelCriteriaObject5
				if var1.DiskEncryption.Criteria != nil {
					var53 = &objectsHipObjectsListDsModelCriteriaObject5{}
					var var54 []objectsHipObjectsListDsModelEncryptedLocationsObject
					if len(var1.DiskEncryption.Criteria.EncryptedLocations) != 0 {
						var54 = make([]objectsHipObjectsListDsModelEncryptedLocationsObject, 0, len(var1.DiskEncryption.Criteria.EncryptedLocations))
						for var55Index := range var1.DiskEncryption.Criteria.EncryptedLocations {
							var55 := var1.DiskEncryption.Criteria.EncryptedLocations[var55Index]
							var var56 objectsHipObjectsListDsModelEncryptedLocationsObject
							var var57 *objectsHipObjectsListDsModelEncryptionStateObject
							if var55.EncryptionState != nil {
								var57 = &objectsHipObjectsListDsModelEncryptionStateObject{}
								var57.Is = types.StringValue(var55.EncryptionState.Is)
								var57.IsNot = types.StringValue(var55.EncryptionState.IsNot)
							}
							var56.EncryptionState = var57
							var56.Name = types.StringValue(var55.Name)
							var54 = append(var54, var56)
						}
					}
					var53.EncryptedLocations = var54
					var53.IsInstalled = types.BoolValue(var1.DiskEncryption.Criteria.IsInstalled)
				}
				var var58 []objectsHipObjectsListDsModelVendorObject
				if len(var1.DiskEncryption.Vendor) != 0 {
					var58 = make([]objectsHipObjectsListDsModelVendorObject, 0, len(var1.DiskEncryption.Vendor))
					for var59Index := range var1.DiskEncryption.Vendor {
						var59 := var1.DiskEncryption.Vendor[var59Index]
						var var60 objectsHipObjectsListDsModelVendorObject
						var60.Name = types.StringValue(var59.Name)
						var60.Product = EncodeStringSlice(var59.Product)
						var58 = append(var58, var60)
					}
				}
				var52.Criteria = var53
				var52.ExcludeVendor = types.BoolValue(var1.DiskEncryption.ExcludeVendor)
				var52.Vendor = var58
			}
			var var61 *objectsHipObjectsListDsModelFirewallObject
			if var1.Firewall != nil {
				var61 = &objectsHipObjectsListDsModelFirewallObject{}
				var var62 *objectsHipObjectsListDsModelCriteriaObject3
				if var1.Firewall.Criteria != nil {
					var62 = &objectsHipObjectsListDsModelCriteriaObject3{}
					var62.IsEnabled = types.StringValue(var1.Firewall.Criteria.IsEnabled)
					var62.IsInstalled = types.BoolValue(var1.Firewall.Criteria.IsInstalled)
				}
				var var63 []objectsHipObjectsListDsModelVendorObject
				if len(var1.Firewall.Vendor) != 0 {
					var63 = make([]objectsHipObjectsListDsModelVendorObject, 0, len(var1.Firewall.Vendor))
					for var64Index := range var1.Firewall.Vendor {
						var64 := var1.Firewall.Vendor[var64Index]
						var var65 objectsHipObjectsListDsModelVendorObject
						var65.Name = types.StringValue(var64.Name)
						var65.Product = EncodeStringSlice(var64.Product)
						var63 = append(var63, var65)
					}
				}
				var61.Criteria = var62
				var61.ExcludeVendor = types.BoolValue(var1.Firewall.ExcludeVendor)
				var61.Vendor = var63
			}
			var var66 *objectsHipObjectsListDsModelHostInfoObject
			if var1.HostInfo != nil {
				var66 = &objectsHipObjectsListDsModelHostInfoObject{}
				var var67 objectsHipObjectsListDsModelCriteriaObject6
				var var68 *objectsHipObjectsListDsModelClientVersionObject
				if var1.HostInfo.Criteria.ClientVersion != nil {
					var68 = &objectsHipObjectsListDsModelClientVersionObject{}
					var68.Contains = types.StringValue(var1.HostInfo.Criteria.ClientVersion.Contains)
					var68.Is = types.StringValue(var1.HostInfo.Criteria.ClientVersion.Is)
					var68.IsNot = types.StringValue(var1.HostInfo.Criteria.ClientVersion.IsNot)
				}
				var var69 *objectsHipObjectsListDsModelDomainObject
				if var1.HostInfo.Criteria.Domain != nil {
					var69 = &objectsHipObjectsListDsModelDomainObject{}
					var69.Contains = types.StringValue(var1.HostInfo.Criteria.Domain.Contains)
					var69.Is = types.StringValue(var1.HostInfo.Criteria.Domain.Is)
					var69.IsNot = types.StringValue(var1.HostInfo.Criteria.Domain.IsNot)
				}
				var var70 *objectsHipObjectsListDsModelHostIdObject
				if var1.HostInfo.Criteria.HostId != nil {
					var70 = &objectsHipObjectsListDsModelHostIdObject{}
					var70.Contains = types.StringValue(var1.HostInfo.Criteria.HostId.Contains)
					var70.Is = types.StringValue(var1.HostInfo.Criteria.HostId.Is)
					var70.IsNot = types.StringValue(var1.HostInfo.Criteria.HostId.IsNot)
				}
				var var71 *objectsHipObjectsListDsModelHostNameObject
				if var1.HostInfo.Criteria.HostName != nil {
					var71 = &objectsHipObjectsListDsModelHostNameObject{}
					var71.Contains = types.StringValue(var1.HostInfo.Criteria.HostName.Contains)
					var71.Is = types.StringValue(var1.HostInfo.Criteria.HostName.Is)
					var71.IsNot = types.StringValue(var1.HostInfo.Criteria.HostName.IsNot)
				}
				var var72 *objectsHipObjectsListDsModelOsObject
				if var1.HostInfo.Criteria.Os != nil {
					var72 = &objectsHipObjectsListDsModelOsObject{}
					var var73 *objectsHipObjectsListDsModelContainsObject
					if var1.HostInfo.Criteria.Os.Contains != nil {
						var73 = &objectsHipObjectsListDsModelContainsObject{}
						var73.Apple = types.StringValue(var1.HostInfo.Criteria.Os.Contains.Apple)
						var73.Google = types.StringValue(var1.HostInfo.Criteria.Os.Contains.Google)
						var73.Linux = types.StringValue(var1.HostInfo.Criteria.Os.Contains.Linux)
						var73.Microsoft = types.StringValue(var1.HostInfo.Criteria.Os.Contains.Microsoft)
						var73.Other = types.StringValue(var1.HostInfo.Criteria.Os.Contains.Other)
					}
					var72.Contains = var73
				}
				var var74 *objectsHipObjectsListDsModelSerialNumberObject
				if var1.HostInfo.Criteria.SerialNumber != nil {
					var74 = &objectsHipObjectsListDsModelSerialNumberObject{}
					var74.Contains = types.StringValue(var1.HostInfo.Criteria.SerialNumber.Contains)
					var74.Is = types.StringValue(var1.HostInfo.Criteria.SerialNumber.Is)
					var74.IsNot = types.StringValue(var1.HostInfo.Criteria.SerialNumber.IsNot)
				}
				var67.ClientVersion = var68
				var67.Domain = var69
				var67.HostId = var70
				var67.HostName = var71
				var67.Managed = types.BoolValue(var1.HostInfo.Criteria.Managed)
				var67.Os = var72
				var67.SerialNumber = var74
				var66.Criteria = var67
			}
			var var75 *objectsHipObjectsListDsModelMobileDeviceObject
			if var1.MobileDevice != nil {
				var75 = &objectsHipObjectsListDsModelMobileDeviceObject{}
				var var76 *objectsHipObjectsListDsModelCriteriaObject7
				if var1.MobileDevice.Criteria != nil {
					var76 = &objectsHipObjectsListDsModelCriteriaObject7{}
					var var77 *objectsHipObjectsListDsModelApplicationsObject
					if var1.MobileDevice.Criteria.Applications != nil {
						var77 = &objectsHipObjectsListDsModelApplicationsObject{}
						var var78 *objectsHipObjectsListDsModelHasMalwareObject
						if var1.MobileDevice.Criteria.Applications.HasMalware != nil {
							var78 = &objectsHipObjectsListDsModelHasMalwareObject{}
							var var79 *objectsHipObjectsListDsModelYesObject
							if var1.MobileDevice.Criteria.Applications.HasMalware.Yes != nil {
								var79 = &objectsHipObjectsListDsModelYesObject{}
								var var80 []objectsHipObjectsListDsModelExcludesObject
								if len(var1.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes) != 0 {
									var80 = make([]objectsHipObjectsListDsModelExcludesObject, 0, len(var1.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes))
									for var81Index := range var1.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes {
										var81 := var1.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes[var81Index]
										var var82 objectsHipObjectsListDsModelExcludesObject
										var82.Hash = types.StringValue(var81.Hash)
										var82.Name = types.StringValue(var81.Name)
										var82.Package = types.StringValue(var81.Package)
										var80 = append(var80, var82)
									}
								}
								var79.Excludes = var80
							}
							if var1.MobileDevice.Criteria.Applications.HasMalware.No != nil {
								var78.No = types.BoolValue(true)
							}
							var78.Yes = var79
						}
						var var83 []objectsHipObjectsListDsModelIncludesObject
						if len(var1.MobileDevice.Criteria.Applications.Includes) != 0 {
							var83 = make([]objectsHipObjectsListDsModelIncludesObject, 0, len(var1.MobileDevice.Criteria.Applications.Includes))
							for var84Index := range var1.MobileDevice.Criteria.Applications.Includes {
								var84 := var1.MobileDevice.Criteria.Applications.Includes[var84Index]
								var var85 objectsHipObjectsListDsModelIncludesObject
								var85.Hash = types.StringValue(var84.Hash)
								var85.Name = types.StringValue(var84.Name)
								var85.Package = types.StringValue(var84.Package)
								var83 = append(var83, var85)
							}
						}
						var77.HasMalware = var78
						var77.HasUnmanagedApp = types.BoolValue(var1.MobileDevice.Criteria.Applications.HasUnmanagedApp)
						var77.Includes = var83
					}
					var var86 *objectsHipObjectsListDsModelImeiObject
					if var1.MobileDevice.Criteria.Imei != nil {
						var86 = &objectsHipObjectsListDsModelImeiObject{}
						var86.Contains = types.StringValue(var1.MobileDevice.Criteria.Imei.Contains)
						var86.Is = types.StringValue(var1.MobileDevice.Criteria.Imei.Is)
						var86.IsNot = types.StringValue(var1.MobileDevice.Criteria.Imei.IsNot)
					}
					var var87 *objectsHipObjectsListDsModelLastCheckinTimeObject
					if var1.MobileDevice.Criteria.LastCheckinTime != nil {
						var87 = &objectsHipObjectsListDsModelLastCheckinTimeObject{}
						var var88 *objectsHipObjectsListDsModelNotWithinObject3
						if var1.MobileDevice.Criteria.LastCheckinTime.NotWithin != nil {
							var88 = &objectsHipObjectsListDsModelNotWithinObject3{}
							var88.Days = types.Int64Value(var1.MobileDevice.Criteria.LastCheckinTime.NotWithin.Days)
						}
						var var89 *objectsHipObjectsListDsModelWithinObject3
						if var1.MobileDevice.Criteria.LastCheckinTime.Within != nil {
							var89 = &objectsHipObjectsListDsModelWithinObject3{}
							var89.Days = types.Int64Value(var1.MobileDevice.Criteria.LastCheckinTime.Within.Days)
						}
						var87.NotWithin = var88
						var87.Within = var89
					}
					var var90 *objectsHipObjectsListDsModelModelObject
					if var1.MobileDevice.Criteria.Model != nil {
						var90 = &objectsHipObjectsListDsModelModelObject{}
						var90.Contains = types.StringValue(var1.MobileDevice.Criteria.Model.Contains)
						var90.Is = types.StringValue(var1.MobileDevice.Criteria.Model.Is)
						var90.IsNot = types.StringValue(var1.MobileDevice.Criteria.Model.IsNot)
					}
					var var91 *objectsHipObjectsListDsModelPhoneNumberObject
					if var1.MobileDevice.Criteria.PhoneNumber != nil {
						var91 = &objectsHipObjectsListDsModelPhoneNumberObject{}
						var91.Contains = types.StringValue(var1.MobileDevice.Criteria.PhoneNumber.Contains)
						var91.Is = types.StringValue(var1.MobileDevice.Criteria.PhoneNumber.Is)
						var91.IsNot = types.StringValue(var1.MobileDevice.Criteria.PhoneNumber.IsNot)
					}
					var var92 *objectsHipObjectsListDsModelTagObject
					if var1.MobileDevice.Criteria.Tag != nil {
						var92 = &objectsHipObjectsListDsModelTagObject{}
						var92.Contains = types.StringValue(var1.MobileDevice.Criteria.Tag.Contains)
						var92.Is = types.StringValue(var1.MobileDevice.Criteria.Tag.Is)
						var92.IsNot = types.StringValue(var1.MobileDevice.Criteria.Tag.IsNot)
					}
					var76.Applications = var77
					var76.DiskEncrypted = types.BoolValue(var1.MobileDevice.Criteria.DiskEncrypted)
					var76.Imei = var86
					var76.Jailbroken = types.BoolValue(var1.MobileDevice.Criteria.Jailbroken)
					var76.LastCheckinTime = var87
					var76.Model = var90
					var76.PasscodeSet = types.BoolValue(var1.MobileDevice.Criteria.PasscodeSet)
					var76.PhoneNumber = var91
					var76.Tag = var92
				}
				var75.Criteria = var76
			}
			var var93 *objectsHipObjectsListDsModelNetworkInfoObject
			if var1.NetworkInfo != nil {
				var93 = &objectsHipObjectsListDsModelNetworkInfoObject{}
				var var94 *objectsHipObjectsListDsModelCriteriaObject8
				if var1.NetworkInfo.Criteria != nil {
					var94 = &objectsHipObjectsListDsModelCriteriaObject8{}
					var var95 *objectsHipObjectsListDsModelNetworkObject
					if var1.NetworkInfo.Criteria.Network != nil {
						var95 = &objectsHipObjectsListDsModelNetworkObject{}
						var var96 *objectsHipObjectsListDsModelIsObject
						if var1.NetworkInfo.Criteria.Network.Is != nil {
							var96 = &objectsHipObjectsListDsModelIsObject{}
							var var97 *objectsHipObjectsListDsModelMobileObject
							if var1.NetworkInfo.Criteria.Network.Is.Mobile != nil {
								var97 = &objectsHipObjectsListDsModelMobileObject{}
								var97.Carrier = types.StringValue(var1.NetworkInfo.Criteria.Network.Is.Mobile.Carrier)
							}
							var var98 *objectsHipObjectsListDsModelWifiObject
							if var1.NetworkInfo.Criteria.Network.Is.Wifi != nil {
								var98 = &objectsHipObjectsListDsModelWifiObject{}
								var98.Ssid = types.StringValue(var1.NetworkInfo.Criteria.Network.Is.Wifi.Ssid)
							}
							var96.Mobile = var97
							if var1.NetworkInfo.Criteria.Network.Is.Unknown != nil {
								var96.Unknown = types.BoolValue(true)
							}
							var96.Wifi = var98
						}
						var var99 *objectsHipObjectsListDsModelIsNotObject
						if var1.NetworkInfo.Criteria.Network.IsNot != nil {
							var99 = &objectsHipObjectsListDsModelIsNotObject{}
							var var100 *objectsHipObjectsListDsModelMobileObject
							if var1.NetworkInfo.Criteria.Network.IsNot.Mobile != nil {
								var100 = &objectsHipObjectsListDsModelMobileObject{}
								var100.Carrier = types.StringValue(var1.NetworkInfo.Criteria.Network.IsNot.Mobile.Carrier)
							}
							var var101 *objectsHipObjectsListDsModelWifiObject
							if var1.NetworkInfo.Criteria.Network.IsNot.Wifi != nil {
								var101 = &objectsHipObjectsListDsModelWifiObject{}
								var101.Ssid = types.StringValue(var1.NetworkInfo.Criteria.Network.IsNot.Wifi.Ssid)
							}
							if var1.NetworkInfo.Criteria.Network.IsNot.Ethernet != nil {
								var99.Ethernet = types.BoolValue(true)
							}
							var99.Mobile = var100
							if var1.NetworkInfo.Criteria.Network.IsNot.Unknown != nil {
								var99.Unknown = types.BoolValue(true)
							}
							var99.Wifi = var101
						}
						var95.Is = var96
						var95.IsNot = var99
					}
					var94.Network = var95
				}
				var93.Criteria = var94
			}
			var var102 *objectsHipObjectsListDsModelPatchManagementObject
			if var1.PatchManagement != nil {
				var102 = &objectsHipObjectsListDsModelPatchManagementObject{}
				var var103 *objectsHipObjectsListDsModelCriteriaObject9
				if var1.PatchManagement.Criteria != nil {
					var103 = &objectsHipObjectsListDsModelCriteriaObject9{}
					var var104 *objectsHipObjectsListDsModelMissingPatchesObject
					if var1.PatchManagement.Criteria.MissingPatches != nil {
						var104 = &objectsHipObjectsListDsModelMissingPatchesObject{}
						var var105 *objectsHipObjectsListDsModelSeverityObject
						if var1.PatchManagement.Criteria.MissingPatches.Severity != nil {
							var105 = &objectsHipObjectsListDsModelSeverityObject{}
							var105.GreaterEqual = types.Int64Value(var1.PatchManagement.Criteria.MissingPatches.Severity.GreaterEqual)
							var105.GreaterThan = types.Int64Value(var1.PatchManagement.Criteria.MissingPatches.Severity.GreaterThan)
							var105.Is = types.Int64Value(var1.PatchManagement.Criteria.MissingPatches.Severity.Is)
							var105.IsNot = types.Int64Value(var1.PatchManagement.Criteria.MissingPatches.Severity.IsNot)
							var105.LessEqual = types.Int64Value(var1.PatchManagement.Criteria.MissingPatches.Severity.LessEqual)
							var105.LessThan = types.Int64Value(var1.PatchManagement.Criteria.MissingPatches.Severity.LessThan)
						}
						var104.Check = types.StringValue(var1.PatchManagement.Criteria.MissingPatches.Check)
						var104.Patches = EncodeStringSlice(var1.PatchManagement.Criteria.MissingPatches.Patches)
						var104.Severity = var105
					}
					var103.IsEnabled = types.StringValue(var1.PatchManagement.Criteria.IsEnabled)
					var103.IsInstalled = types.BoolValue(var1.PatchManagement.Criteria.IsInstalled)
					var103.MissingPatches = var104
				}
				var var106 []objectsHipObjectsListDsModelVendorObject1
				if len(var1.PatchManagement.Vendor) != 0 {
					var106 = make([]objectsHipObjectsListDsModelVendorObject1, 0, len(var1.PatchManagement.Vendor))
					for var107Index := range var1.PatchManagement.Vendor {
						var107 := var1.PatchManagement.Vendor[var107Index]
						var var108 objectsHipObjectsListDsModelVendorObject1
						var108.Name = types.StringValue(var107.Name)
						var108.Product = EncodeStringSlice(var107.Product)
						var106 = append(var106, var108)
					}
				}
				var102.Criteria = var103
				var102.ExcludeVendor = types.BoolValue(var1.PatchManagement.ExcludeVendor)
				var102.Vendor = var106
			}
			var2.AntiMalware = var3
			var2.Certificate = var17
			var2.CustomChecks = var22
			var2.DataLossPrevention = var39
			var2.Description = types.StringValue(var1.Description)
			var2.DiskBackup = var44
			var2.DiskEncryption = var52
			var2.Firewall = var61
			var2.HostInfo = var66
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.MobileDevice = var75
			var2.Name = types.StringValue(var1.Name)
			var2.NetworkInfo = var93
			var2.PatchManagement = var102
			var0 = append(var0, var2)
		}
	}
	state.Data = var0
	if !state.Limit.IsNull() {
		state.Limit = types.Int64Value(ans.Limit)
	}
	if !state.Offset.IsNull() {
		state.Offset = types.Int64Value(ans.Offset)
	}
	state.Total = types.Int64Value(ans.Total)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Data source.
var (
	_ datasource.DataSource              = &objectsHipObjectsDataSource{}
	_ datasource.DataSourceWithConfigure = &objectsHipObjectsDataSource{}
)

func NewObjectsHipObjectsDataSource() datasource.DataSource {
	return &objectsHipObjectsDataSource{}
}

type objectsHipObjectsDataSource struct {
	client *sase.Client
}

type objectsHipObjectsDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`

	// Output.
	// Ref: #/components/schemas/objects-hip-objects
	AntiMalware        *objectsHipObjectsDsModelAntiMalwareObject        `tfsdk:"anti_malware"`
	Certificate        *objectsHipObjectsDsModelCertificateObject        `tfsdk:"certificate"`
	CustomChecks       *objectsHipObjectsDsModelCustomChecksObject       `tfsdk:"custom_checks"`
	DataLossPrevention *objectsHipObjectsDsModelDataLossPreventionObject `tfsdk:"data_loss_prevention"`
	Description        types.String                                      `tfsdk:"description"`
	DiskBackup         *objectsHipObjectsDsModelDiskBackupObject         `tfsdk:"disk_backup"`
	DiskEncryption     *objectsHipObjectsDsModelDiskEncryptionObject     `tfsdk:"disk_encryption"`
	Firewall           *objectsHipObjectsDsModelFirewallObject           `tfsdk:"firewall"`
	HostInfo           *objectsHipObjectsDsModelHostInfoObject           `tfsdk:"host_info"`
	// input omit: ObjectId
	MobileDevice    *objectsHipObjectsDsModelMobileDeviceObject    `tfsdk:"mobile_device"`
	Name            types.String                                   `tfsdk:"name"`
	NetworkInfo     *objectsHipObjectsDsModelNetworkInfoObject     `tfsdk:"network_info"`
	PatchManagement *objectsHipObjectsDsModelPatchManagementObject `tfsdk:"patch_management"`
}

type objectsHipObjectsDsModelAntiMalwareObject struct {
	Criteria      *objectsHipObjectsDsModelCriteriaObject `tfsdk:"criteria"`
	ExcludeVendor types.Bool                              `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsDsModelVendorObject  `tfsdk:"vendor"`
}

type objectsHipObjectsDsModelCriteriaObject struct {
	IsInstalled        types.Bool                                    `tfsdk:"is_installed"`
	LastScanTime       *objectsHipObjectsDsModelLastScanTimeObject   `tfsdk:"last_scan_time"`
	ProductVersion     *objectsHipObjectsDsModelProductVersionObject `tfsdk:"product_version"`
	RealTimeProtection types.String                                  `tfsdk:"real_time_protection"`
	VirdefVersion      *objectsHipObjectsDsModelVirdefVersionObject  `tfsdk:"virdef_version"`
}

type objectsHipObjectsDsModelLastScanTimeObject struct {
	NotAvailable types.Bool                               `tfsdk:"not_available"`
	NotWithin    *objectsHipObjectsDsModelNotWithinObject `tfsdk:"not_within"`
	Within       *objectsHipObjectsDsModelWithinObject    `tfsdk:"within"`
}

type objectsHipObjectsDsModelNotWithinObject struct {
	Days  types.Int64 `tfsdk:"days"`
	Hours types.Int64 `tfsdk:"hours"`
}

type objectsHipObjectsDsModelWithinObject struct {
	Days  types.Int64 `tfsdk:"days"`
	Hours types.Int64 `tfsdk:"hours"`
}

type objectsHipObjectsDsModelProductVersionObject struct {
	Contains     types.String                              `tfsdk:"contains"`
	GreaterEqual types.String                              `tfsdk:"greater_equal"`
	GreaterThan  types.String                              `tfsdk:"greater_than"`
	Is           types.String                              `tfsdk:"is"`
	IsNot        types.String                              `tfsdk:"is_not"`
	LessEqual    types.String                              `tfsdk:"less_equal"`
	LessThan     types.String                              `tfsdk:"less_than"`
	NotWithin    *objectsHipObjectsDsModelNotWithinObject1 `tfsdk:"not_within"`
	Within       *objectsHipObjectsDsModelWithinObject1    `tfsdk:"within"`
}

type objectsHipObjectsDsModelNotWithinObject1 struct {
	Versions types.Int64 `tfsdk:"versions"`
}

type objectsHipObjectsDsModelWithinObject1 struct {
	Versions types.Int64 `tfsdk:"versions"`
}

type objectsHipObjectsDsModelVirdefVersionObject struct {
	NotWithin *objectsHipObjectsDsModelNotWithinObject2 `tfsdk:"not_within"`
	Within    *objectsHipObjectsDsModelWithinObject2    `tfsdk:"within"`
}

type objectsHipObjectsDsModelNotWithinObject2 struct {
	Days     types.Int64 `tfsdk:"days"`
	Versions types.Int64 `tfsdk:"versions"`
}

type objectsHipObjectsDsModelWithinObject2 struct {
	Days     types.Int64 `tfsdk:"days"`
	Versions types.Int64 `tfsdk:"versions"`
}

type objectsHipObjectsDsModelVendorObject struct {
	Name    types.String   `tfsdk:"name"`
	Product []types.String `tfsdk:"product"`
}

type objectsHipObjectsDsModelCertificateObject struct {
	Criteria *objectsHipObjectsDsModelCriteriaObject1 `tfsdk:"criteria"`
}

type objectsHipObjectsDsModelCriteriaObject1 struct {
	CertificateAttributes []objectsHipObjectsDsModelCertificateAttributesObject `tfsdk:"certificate_attributes"`
	CertificateProfile    types.String                                          `tfsdk:"certificate_profile"`
}

type objectsHipObjectsDsModelCertificateAttributesObject struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type objectsHipObjectsDsModelCustomChecksObject struct {
	Criteria objectsHipObjectsDsModelCriteriaObject2 `tfsdk:"criteria"`
}

type objectsHipObjectsDsModelCriteriaObject2 struct {
	Plist       []objectsHipObjectsDsModelPlistObject       `tfsdk:"plist"`
	ProcessList []objectsHipObjectsDsModelProcessListObject `tfsdk:"process_list"`
	RegistryKey []objectsHipObjectsDsModelRegistryKeyObject `tfsdk:"registry_key"`
}

type objectsHipObjectsDsModelPlistObject struct {
	Key    []objectsHipObjectsDsModelKeyObject `tfsdk:"key"`
	Name   types.String                        `tfsdk:"name"`
	Negate types.Bool                          `tfsdk:"negate"`
}

type objectsHipObjectsDsModelKeyObject struct {
	Name   types.String `tfsdk:"name"`
	Negate types.Bool   `tfsdk:"negate"`
	Value  types.String `tfsdk:"value"`
}

type objectsHipObjectsDsModelProcessListObject struct {
	Name    types.String `tfsdk:"name"`
	Running types.Bool   `tfsdk:"running"`
}

type objectsHipObjectsDsModelRegistryKeyObject struct {
	DefaultValueData types.String                                  `tfsdk:"default_value_data"`
	Name             types.String                                  `tfsdk:"name"`
	Negate           types.Bool                                    `tfsdk:"negate"`
	RegistryValue    []objectsHipObjectsDsModelRegistryValueObject `tfsdk:"registry_value"`
}

type objectsHipObjectsDsModelRegistryValueObject struct {
	Name      types.String `tfsdk:"name"`
	Negate    types.Bool   `tfsdk:"negate"`
	ValueData types.String `tfsdk:"value_data"`
}

type objectsHipObjectsDsModelDataLossPreventionObject struct {
	Criteria      *objectsHipObjectsDsModelCriteriaObject3 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                               `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsDsModelVendorObject1  `tfsdk:"vendor"`
}

type objectsHipObjectsDsModelCriteriaObject3 struct {
	IsEnabled   types.String `tfsdk:"is_enabled"`
	IsInstalled types.Bool   `tfsdk:"is_installed"`
}

type objectsHipObjectsDsModelVendorObject1 struct {
	Name    types.String   `tfsdk:"name"`
	Product []types.String `tfsdk:"product"`
}

type objectsHipObjectsDsModelDiskBackupObject struct {
	Criteria      *objectsHipObjectsDsModelCriteriaObject4 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                               `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsDsModelVendorObject   `tfsdk:"vendor"`
}

type objectsHipObjectsDsModelCriteriaObject4 struct {
	IsInstalled    types.Bool                                    `tfsdk:"is_installed"`
	LastBackupTime *objectsHipObjectsDsModelLastBackupTimeObject `tfsdk:"last_backup_time"`
}

type objectsHipObjectsDsModelLastBackupTimeObject struct {
	NotAvailable types.Bool                               `tfsdk:"not_available"`
	NotWithin    *objectsHipObjectsDsModelNotWithinObject `tfsdk:"not_within"`
	Within       *objectsHipObjectsDsModelWithinObject    `tfsdk:"within"`
}

type objectsHipObjectsDsModelDiskEncryptionObject struct {
	Criteria      *objectsHipObjectsDsModelCriteriaObject5 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                               `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsDsModelVendorObject   `tfsdk:"vendor"`
}

type objectsHipObjectsDsModelCriteriaObject5 struct {
	EncryptedLocations []objectsHipObjectsDsModelEncryptedLocationsObject `tfsdk:"encrypted_locations"`
	IsInstalled        types.Bool                                         `tfsdk:"is_installed"`
}

type objectsHipObjectsDsModelEncryptedLocationsObject struct {
	EncryptionState *objectsHipObjectsDsModelEncryptionStateObject `tfsdk:"encryption_state"`
	Name            types.String                                   `tfsdk:"name"`
}

type objectsHipObjectsDsModelEncryptionStateObject struct {
	Is    types.String `tfsdk:"is"`
	IsNot types.String `tfsdk:"is_not"`
}

type objectsHipObjectsDsModelFirewallObject struct {
	Criteria      *objectsHipObjectsDsModelCriteriaObject3 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                               `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsDsModelVendorObject   `tfsdk:"vendor"`
}

type objectsHipObjectsDsModelHostInfoObject struct {
	Criteria objectsHipObjectsDsModelCriteriaObject6 `tfsdk:"criteria"`
}

type objectsHipObjectsDsModelCriteriaObject6 struct {
	ClientVersion *objectsHipObjectsDsModelClientVersionObject `tfsdk:"client_version"`
	Domain        *objectsHipObjectsDsModelDomainObject        `tfsdk:"domain"`
	HostId        *objectsHipObjectsDsModelHostIdObject        `tfsdk:"host_id"`
	HostName      *objectsHipObjectsDsModelHostNameObject      `tfsdk:"host_name"`
	Managed       types.Bool                                   `tfsdk:"managed"`
	Os            *objectsHipObjectsDsModelOsObject            `tfsdk:"os"`
	SerialNumber  *objectsHipObjectsDsModelSerialNumberObject  `tfsdk:"serial_number"`
}

type objectsHipObjectsDsModelClientVersionObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsDsModelDomainObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsDsModelHostIdObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsDsModelHostNameObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsDsModelOsObject struct {
	Contains *objectsHipObjectsDsModelContainsObject `tfsdk:"contains"`
}

type objectsHipObjectsDsModelContainsObject struct {
	Apple     types.String `tfsdk:"apple"`
	Google    types.String `tfsdk:"google"`
	Linux     types.String `tfsdk:"linux"`
	Microsoft types.String `tfsdk:"microsoft"`
	Other     types.String `tfsdk:"other"`
}

type objectsHipObjectsDsModelSerialNumberObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsDsModelMobileDeviceObject struct {
	Criteria *objectsHipObjectsDsModelCriteriaObject7 `tfsdk:"criteria"`
}

type objectsHipObjectsDsModelCriteriaObject7 struct {
	Applications    *objectsHipObjectsDsModelApplicationsObject    `tfsdk:"applications"`
	DiskEncrypted   types.Bool                                     `tfsdk:"disk_encrypted"`
	Imei            *objectsHipObjectsDsModelImeiObject            `tfsdk:"imei"`
	Jailbroken      types.Bool                                     `tfsdk:"jailbroken"`
	LastCheckinTime *objectsHipObjectsDsModelLastCheckinTimeObject `tfsdk:"last_checkin_time"`
	Model           *objectsHipObjectsDsModelModelObject           `tfsdk:"model"`
	PasscodeSet     types.Bool                                     `tfsdk:"passcode_set"`
	PhoneNumber     *objectsHipObjectsDsModelPhoneNumberObject     `tfsdk:"phone_number"`
	Tag             *objectsHipObjectsDsModelTagObject             `tfsdk:"tag"`
}

type objectsHipObjectsDsModelApplicationsObject struct {
	HasMalware      *objectsHipObjectsDsModelHasMalwareObject `tfsdk:"has_malware"`
	HasUnmanagedApp types.Bool                                `tfsdk:"has_unmanaged_app"`
	Includes        []objectsHipObjectsDsModelIncludesObject  `tfsdk:"includes"`
}

type objectsHipObjectsDsModelHasMalwareObject struct {
	No  types.Bool                         `tfsdk:"no"`
	Yes *objectsHipObjectsDsModelYesObject `tfsdk:"yes"`
}

type objectsHipObjectsDsModelYesObject struct {
	Excludes []objectsHipObjectsDsModelExcludesObject `tfsdk:"excludes"`
}

type objectsHipObjectsDsModelExcludesObject struct {
	Hash    types.String `tfsdk:"hash"`
	Name    types.String `tfsdk:"name"`
	Package types.String `tfsdk:"package"`
}

type objectsHipObjectsDsModelIncludesObject struct {
	Hash    types.String `tfsdk:"hash"`
	Name    types.String `tfsdk:"name"`
	Package types.String `tfsdk:"package"`
}

type objectsHipObjectsDsModelImeiObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsDsModelLastCheckinTimeObject struct {
	NotWithin *objectsHipObjectsDsModelNotWithinObject3 `tfsdk:"not_within"`
	Within    *objectsHipObjectsDsModelWithinObject3    `tfsdk:"within"`
}

type objectsHipObjectsDsModelNotWithinObject3 struct {
	Days types.Int64 `tfsdk:"days"`
}

type objectsHipObjectsDsModelWithinObject3 struct {
	Days types.Int64 `tfsdk:"days"`
}

type objectsHipObjectsDsModelModelObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsDsModelPhoneNumberObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsDsModelTagObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsDsModelNetworkInfoObject struct {
	Criteria *objectsHipObjectsDsModelCriteriaObject8 `tfsdk:"criteria"`
}

type objectsHipObjectsDsModelCriteriaObject8 struct {
	Network *objectsHipObjectsDsModelNetworkObject `tfsdk:"network"`
}

type objectsHipObjectsDsModelNetworkObject struct {
	Is    *objectsHipObjectsDsModelIsObject    `tfsdk:"is"`
	IsNot *objectsHipObjectsDsModelIsNotObject `tfsdk:"is_not"`
}

type objectsHipObjectsDsModelIsObject struct {
	Mobile  *objectsHipObjectsDsModelMobileObject `tfsdk:"mobile"`
	Unknown types.Bool                            `tfsdk:"unknown"`
	Wifi    *objectsHipObjectsDsModelWifiObject   `tfsdk:"wifi"`
}

type objectsHipObjectsDsModelMobileObject struct {
	Carrier types.String `tfsdk:"carrier"`
}

type objectsHipObjectsDsModelWifiObject struct {
	Ssid types.String `tfsdk:"ssid"`
}

type objectsHipObjectsDsModelIsNotObject struct {
	Ethernet types.Bool                            `tfsdk:"ethernet"`
	Mobile   *objectsHipObjectsDsModelMobileObject `tfsdk:"mobile"`
	Unknown  types.Bool                            `tfsdk:"unknown"`
	Wifi     *objectsHipObjectsDsModelWifiObject   `tfsdk:"wifi"`
}

type objectsHipObjectsDsModelPatchManagementObject struct {
	Criteria      *objectsHipObjectsDsModelCriteriaObject9 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                               `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsDsModelVendorObject1  `tfsdk:"vendor"`
}

type objectsHipObjectsDsModelCriteriaObject9 struct {
	IsEnabled      types.String                                  `tfsdk:"is_enabled"`
	IsInstalled    types.Bool                                    `tfsdk:"is_installed"`
	MissingPatches *objectsHipObjectsDsModelMissingPatchesObject `tfsdk:"missing_patches"`
}

type objectsHipObjectsDsModelMissingPatchesObject struct {
	Check    types.String                            `tfsdk:"check"`
	Patches  []types.String                          `tfsdk:"patches"`
	Severity *objectsHipObjectsDsModelSeverityObject `tfsdk:"severity"`
}

type objectsHipObjectsDsModelSeverityObject struct {
	GreaterEqual types.Int64 `tfsdk:"greater_equal"`
	GreaterThan  types.Int64 `tfsdk:"greater_than"`
	Is           types.Int64 `tfsdk:"is"`
	IsNot        types.Int64 `tfsdk:"is_not"`
	LessEqual    types.Int64 `tfsdk:"less_equal"`
	LessThan     types.Int64 `tfsdk:"less_than"`
}

// Metadata returns the data source type name.
func (d *objectsHipObjectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_hip_objects"
}

// Schema defines the schema for this listing data source.
func (d *objectsHipObjectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Retrieves config for a specific item.",

		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Description: "The object ID.",
				Computed:    true,
			},

			// Input.
			"object_id": dsschema.StringAttribute{
				Description: "The uuid of the resource",
				Required:    true,
			},

			// Output.
			"anti_malware": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"criteria": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"is_installed": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"last_scan_time": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"not_available": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"not_within": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"days": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
											"hours": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"within": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"days": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
											"hours": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
							"product_version": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"contains": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"greater_equal": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"greater_than": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is_not": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"less_equal": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"less_than": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"not_within": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"versions": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"within": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"versions": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
							"real_time_protection": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"virdef_version": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"not_within": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"days": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
											"versions": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"within": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"days": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
											"versions": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
						},
					},
					"exclude_vendor": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"vendor": dsschema.ListNestedAttribute{
						Description: "",
						Computed:    true,
						NestedObject: dsschema.NestedAttributeObject{
							Attributes: map[string]dsschema.Attribute{
								"name": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"product": dsschema.ListAttribute{
									Description: "",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"certificate": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"criteria": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"certificate_attributes": dsschema.ListNestedAttribute{
								Description: "",
								Computed:    true,
								NestedObject: dsschema.NestedAttributeObject{
									Attributes: map[string]dsschema.Attribute{
										"name": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"value": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
							},
							"certificate_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
				},
			},
			"custom_checks": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"criteria": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"plist": dsschema.ListNestedAttribute{
								Description: "",
								Computed:    true,
								NestedObject: dsschema.NestedAttributeObject{
									Attributes: map[string]dsschema.Attribute{
										"key": dsschema.ListNestedAttribute{
											Description: "",
											Computed:    true,
											NestedObject: dsschema.NestedAttributeObject{
												Attributes: map[string]dsschema.Attribute{
													"name": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"negate": dsschema.BoolAttribute{
														Description: "",
														Computed:    true,
													},
													"value": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
												},
											},
										},
										"name": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"negate": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
							},
							"process_list": dsschema.ListNestedAttribute{
								Description: "",
								Computed:    true,
								NestedObject: dsschema.NestedAttributeObject{
									Attributes: map[string]dsschema.Attribute{
										"name": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"running": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
							},
							"registry_key": dsschema.ListNestedAttribute{
								Description: "",
								Computed:    true,
								NestedObject: dsschema.NestedAttributeObject{
									Attributes: map[string]dsschema.Attribute{
										"default_value_data": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"name": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"negate": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"registry_value": dsschema.ListNestedAttribute{
											Description: "",
											Computed:    true,
											NestedObject: dsschema.NestedAttributeObject{
												Attributes: map[string]dsschema.Attribute{
													"name": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
													"negate": dsschema.BoolAttribute{
														Description: "",
														Computed:    true,
													},
													"value_data": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"data_loss_prevention": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"criteria": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"is_enabled": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"is_installed": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"exclude_vendor": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"vendor": dsschema.ListNestedAttribute{
						Description: "",
						Computed:    true,
						NestedObject: dsschema.NestedAttributeObject{
							Attributes: map[string]dsschema.Attribute{
								"name": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"product": dsschema.ListAttribute{
									Description: "",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"description": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"disk_backup": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"criteria": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"is_installed": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"last_backup_time": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"not_available": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"not_within": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"days": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
											"hours": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"within": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"days": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
											"hours": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
						},
					},
					"exclude_vendor": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"vendor": dsschema.ListNestedAttribute{
						Description: "",
						Computed:    true,
						NestedObject: dsschema.NestedAttributeObject{
							Attributes: map[string]dsschema.Attribute{
								"name": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"product": dsschema.ListAttribute{
									Description: "",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"disk_encryption": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"criteria": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"encrypted_locations": dsschema.ListNestedAttribute{
								Description: "",
								Computed:    true,
								NestedObject: dsschema.NestedAttributeObject{
									Attributes: map[string]dsschema.Attribute{
										"encryption_state": dsschema.SingleNestedAttribute{
											Description: "",
											Computed:    true,
											Attributes: map[string]dsschema.Attribute{
												"is": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"is_not": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
										"name": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
							},
							"is_installed": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"exclude_vendor": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"vendor": dsschema.ListNestedAttribute{
						Description: "",
						Computed:    true,
						NestedObject: dsschema.NestedAttributeObject{
							Attributes: map[string]dsschema.Attribute{
								"name": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"product": dsschema.ListAttribute{
									Description: "",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"firewall": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"criteria": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"is_enabled": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"is_installed": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"exclude_vendor": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"vendor": dsschema.ListNestedAttribute{
						Description: "",
						Computed:    true,
						NestedObject: dsschema.NestedAttributeObject{
							Attributes: map[string]dsschema.Attribute{
								"name": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"product": dsschema.ListAttribute{
									Description: "",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"host_info": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"criteria": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"client_version": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"contains": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is_not": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
							"domain": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"contains": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is_not": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
							"host_id": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"contains": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is_not": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
							"host_name": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"contains": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is_not": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
							"managed": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"os": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"contains": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"apple": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"google": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"linux": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"microsoft": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"other": dsschema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
							"serial_number": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"contains": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is_not": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"mobile_device": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"criteria": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"applications": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"has_malware": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"no": dsschema.BoolAttribute{
												Description: "",
												Computed:    true,
											},
											"yes": dsschema.SingleNestedAttribute{
												Description: "",
												Computed:    true,
												Attributes: map[string]dsschema.Attribute{
													"excludes": dsschema.ListNestedAttribute{
														Description: "",
														Computed:    true,
														NestedObject: dsschema.NestedAttributeObject{
															Attributes: map[string]dsschema.Attribute{
																"hash": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
																"name": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
																"package": dsschema.StringAttribute{
																	Description: "",
																	Computed:    true,
																},
															},
														},
													},
												},
											},
										},
									},
									"has_unmanaged_app": dsschema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"includes": dsschema.ListNestedAttribute{
										Description: "",
										Computed:    true,
										NestedObject: dsschema.NestedAttributeObject{
											Attributes: map[string]dsschema.Attribute{
												"hash": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"name": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
												"package": dsschema.StringAttribute{
													Description: "",
													Computed:    true,
												},
											},
										},
									},
								},
							},
							"disk_encrypted": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"imei": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"contains": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is_not": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
							"jailbroken": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"last_checkin_time": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"not_within": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"days": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
									"within": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"days": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
							"model": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"contains": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is_not": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
							"passcode_set": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"phone_number": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"contains": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is_not": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
							"tag": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"contains": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is_not": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"network_info": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"criteria": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"network": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"is": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"mobile": dsschema.SingleNestedAttribute{
												Description: "",
												Computed:    true,
												Attributes: map[string]dsschema.Attribute{
													"carrier": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
												},
											},
											"unknown": dsschema.BoolAttribute{
												Description: "",
												Computed:    true,
											},
											"wifi": dsschema.SingleNestedAttribute{
												Description: "",
												Computed:    true,
												Attributes: map[string]dsschema.Attribute{
													"ssid": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
												},
											},
										},
									},
									"is_not": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"ethernet": dsschema.BoolAttribute{
												Description: "",
												Computed:    true,
											},
											"mobile": dsschema.SingleNestedAttribute{
												Description: "",
												Computed:    true,
												Attributes: map[string]dsschema.Attribute{
													"carrier": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
												},
											},
											"unknown": dsschema.BoolAttribute{
												Description: "",
												Computed:    true,
											},
											"wifi": dsschema.SingleNestedAttribute{
												Description: "",
												Computed:    true,
												Attributes: map[string]dsschema.Attribute{
													"ssid": dsschema.StringAttribute{
														Description: "",
														Computed:    true,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"patch_management": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"criteria": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"is_enabled": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"is_installed": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"missing_patches": dsschema.SingleNestedAttribute{
								Description: "",
								Computed:    true,
								Attributes: map[string]dsschema.Attribute{
									"check": dsschema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"patches": dsschema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"severity": dsschema.SingleNestedAttribute{
										Description: "",
										Computed:    true,
										Attributes: map[string]dsschema.Attribute{
											"greater_equal": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
											"greater_than": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
											"is": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
											"is_not": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
											"less_equal": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
											"less_than": dsschema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
							},
						},
					},
					"exclude_vendor": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"vendor": dsschema.ListNestedAttribute{
						Description: "",
						Computed:    true,
						NestedObject: dsschema.NestedAttributeObject{
							Attributes: map[string]dsschema.Attribute{
								"name": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"product": dsschema.ListAttribute{
									Description: "",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (d *objectsHipObjectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *objectsHipObjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state objectsHipObjectsDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name": "sase_objects_hip_objects",
		"object_id":        state.ObjectId.ValueString(),
	})

	// Prepare to run the command.
	svc := yCYVNEN.NewClient(d.client)
	input := yCYVNEN.ReadInput{
		ObjectId: state.ObjectId.ValueString(),
	}

	// Perform the operation.
	ans, err := svc.Read(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error getting singleton", err.Error())
		return
	}

	// Store the answer to state.
	state.Id = types.StringValue(strings.Join([]string{input.ObjectId}, IdSeparator))
	var var0 *objectsHipObjectsDsModelAntiMalwareObject
	if ans.AntiMalware != nil {
		var0 = &objectsHipObjectsDsModelAntiMalwareObject{}
		var var1 *objectsHipObjectsDsModelCriteriaObject
		if ans.AntiMalware.Criteria != nil {
			var1 = &objectsHipObjectsDsModelCriteriaObject{}
			var var2 *objectsHipObjectsDsModelLastScanTimeObject
			if ans.AntiMalware.Criteria.LastScanTime != nil {
				var2 = &objectsHipObjectsDsModelLastScanTimeObject{}
				var var3 *objectsHipObjectsDsModelNotWithinObject
				if ans.AntiMalware.Criteria.LastScanTime.NotWithin != nil {
					var3 = &objectsHipObjectsDsModelNotWithinObject{}
					var3.Days = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.NotWithin.Days)
					var3.Hours = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.NotWithin.Hours)
				}
				var var4 *objectsHipObjectsDsModelWithinObject
				if ans.AntiMalware.Criteria.LastScanTime.Within != nil {
					var4 = &objectsHipObjectsDsModelWithinObject{}
					var4.Days = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.Within.Days)
					var4.Hours = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.Within.Hours)
				}
				if ans.AntiMalware.Criteria.LastScanTime.NotAvailable != nil {
					var2.NotAvailable = types.BoolValue(true)
				}
				var2.NotWithin = var3
				var2.Within = var4
			}
			var var5 *objectsHipObjectsDsModelProductVersionObject
			if ans.AntiMalware.Criteria.ProductVersion != nil {
				var5 = &objectsHipObjectsDsModelProductVersionObject{}
				var var6 *objectsHipObjectsDsModelNotWithinObject1
				if ans.AntiMalware.Criteria.ProductVersion.NotWithin != nil {
					var6 = &objectsHipObjectsDsModelNotWithinObject1{}
					var6.Versions = types.Int64Value(ans.AntiMalware.Criteria.ProductVersion.NotWithin.Versions)
				}
				var var7 *objectsHipObjectsDsModelWithinObject1
				if ans.AntiMalware.Criteria.ProductVersion.Within != nil {
					var7 = &objectsHipObjectsDsModelWithinObject1{}
					var7.Versions = types.Int64Value(ans.AntiMalware.Criteria.ProductVersion.Within.Versions)
				}
				var5.Contains = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.Contains)
				var5.GreaterEqual = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.GreaterEqual)
				var5.GreaterThan = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.GreaterThan)
				var5.Is = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.Is)
				var5.IsNot = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.IsNot)
				var5.LessEqual = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.LessEqual)
				var5.LessThan = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.LessThan)
				var5.NotWithin = var6
				var5.Within = var7
			}
			var var8 *objectsHipObjectsDsModelVirdefVersionObject
			if ans.AntiMalware.Criteria.VirdefVersion != nil {
				var8 = &objectsHipObjectsDsModelVirdefVersionObject{}
				var var9 *objectsHipObjectsDsModelNotWithinObject2
				if ans.AntiMalware.Criteria.VirdefVersion.NotWithin != nil {
					var9 = &objectsHipObjectsDsModelNotWithinObject2{}
					var9.Days = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.NotWithin.Days)
					var9.Versions = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.NotWithin.Versions)
				}
				var var10 *objectsHipObjectsDsModelWithinObject2
				if ans.AntiMalware.Criteria.VirdefVersion.Within != nil {
					var10 = &objectsHipObjectsDsModelWithinObject2{}
					var10.Days = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.Within.Days)
					var10.Versions = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.Within.Versions)
				}
				var8.NotWithin = var9
				var8.Within = var10
			}
			var1.IsInstalled = types.BoolValue(ans.AntiMalware.Criteria.IsInstalled)
			var1.LastScanTime = var2
			var1.ProductVersion = var5
			var1.RealTimeProtection = types.StringValue(ans.AntiMalware.Criteria.RealTimeProtection)
			var1.VirdefVersion = var8
		}
		var var11 []objectsHipObjectsDsModelVendorObject
		if len(ans.AntiMalware.Vendor) != 0 {
			var11 = make([]objectsHipObjectsDsModelVendorObject, 0, len(ans.AntiMalware.Vendor))
			for var12Index := range ans.AntiMalware.Vendor {
				var12 := ans.AntiMalware.Vendor[var12Index]
				var var13 objectsHipObjectsDsModelVendorObject
				var13.Name = types.StringValue(var12.Name)
				var13.Product = EncodeStringSlice(var12.Product)
				var11 = append(var11, var13)
			}
		}
		var0.Criteria = var1
		var0.ExcludeVendor = types.BoolValue(ans.AntiMalware.ExcludeVendor)
		var0.Vendor = var11
	}
	var var14 *objectsHipObjectsDsModelCertificateObject
	if ans.Certificate != nil {
		var14 = &objectsHipObjectsDsModelCertificateObject{}
		var var15 *objectsHipObjectsDsModelCriteriaObject1
		if ans.Certificate.Criteria != nil {
			var15 = &objectsHipObjectsDsModelCriteriaObject1{}
			var var16 []objectsHipObjectsDsModelCertificateAttributesObject
			if len(ans.Certificate.Criteria.CertificateAttributes) != 0 {
				var16 = make([]objectsHipObjectsDsModelCertificateAttributesObject, 0, len(ans.Certificate.Criteria.CertificateAttributes))
				for var17Index := range ans.Certificate.Criteria.CertificateAttributes {
					var17 := ans.Certificate.Criteria.CertificateAttributes[var17Index]
					var var18 objectsHipObjectsDsModelCertificateAttributesObject
					var18.Name = types.StringValue(var17.Name)
					var18.Value = types.StringValue(var17.Value)
					var16 = append(var16, var18)
				}
			}
			var15.CertificateAttributes = var16
			var15.CertificateProfile = types.StringValue(ans.Certificate.Criteria.CertificateProfile)
		}
		var14.Criteria = var15
	}
	var var19 *objectsHipObjectsDsModelCustomChecksObject
	if ans.CustomChecks != nil {
		var19 = &objectsHipObjectsDsModelCustomChecksObject{}
		var var20 objectsHipObjectsDsModelCriteriaObject2
		var var21 []objectsHipObjectsDsModelPlistObject
		if len(ans.CustomChecks.Criteria.Plist) != 0 {
			var21 = make([]objectsHipObjectsDsModelPlistObject, 0, len(ans.CustomChecks.Criteria.Plist))
			for var22Index := range ans.CustomChecks.Criteria.Plist {
				var22 := ans.CustomChecks.Criteria.Plist[var22Index]
				var var23 objectsHipObjectsDsModelPlistObject
				var var24 []objectsHipObjectsDsModelKeyObject
				if len(var22.Key) != 0 {
					var24 = make([]objectsHipObjectsDsModelKeyObject, 0, len(var22.Key))
					for var25Index := range var22.Key {
						var25 := var22.Key[var25Index]
						var var26 objectsHipObjectsDsModelKeyObject
						var26.Name = types.StringValue(var25.Name)
						var26.Negate = types.BoolValue(var25.Negate)
						var26.Value = types.StringValue(var25.Value)
						var24 = append(var24, var26)
					}
				}
				var23.Key = var24
				var23.Name = types.StringValue(var22.Name)
				var23.Negate = types.BoolValue(var22.Negate)
				var21 = append(var21, var23)
			}
		}
		var var27 []objectsHipObjectsDsModelProcessListObject
		if len(ans.CustomChecks.Criteria.ProcessList) != 0 {
			var27 = make([]objectsHipObjectsDsModelProcessListObject, 0, len(ans.CustomChecks.Criteria.ProcessList))
			for var28Index := range ans.CustomChecks.Criteria.ProcessList {
				var28 := ans.CustomChecks.Criteria.ProcessList[var28Index]
				var var29 objectsHipObjectsDsModelProcessListObject
				var29.Name = types.StringValue(var28.Name)
				var29.Running = types.BoolValue(var28.Running)
				var27 = append(var27, var29)
			}
		}
		var var30 []objectsHipObjectsDsModelRegistryKeyObject
		if len(ans.CustomChecks.Criteria.RegistryKey) != 0 {
			var30 = make([]objectsHipObjectsDsModelRegistryKeyObject, 0, len(ans.CustomChecks.Criteria.RegistryKey))
			for var31Index := range ans.CustomChecks.Criteria.RegistryKey {
				var31 := ans.CustomChecks.Criteria.RegistryKey[var31Index]
				var var32 objectsHipObjectsDsModelRegistryKeyObject
				var var33 []objectsHipObjectsDsModelRegistryValueObject
				if len(var31.RegistryValue) != 0 {
					var33 = make([]objectsHipObjectsDsModelRegistryValueObject, 0, len(var31.RegistryValue))
					for var34Index := range var31.RegistryValue {
						var34 := var31.RegistryValue[var34Index]
						var var35 objectsHipObjectsDsModelRegistryValueObject
						var35.Name = types.StringValue(var34.Name)
						var35.Negate = types.BoolValue(var34.Negate)
						var35.ValueData = types.StringValue(var34.ValueData)
						var33 = append(var33, var35)
					}
				}
				var32.DefaultValueData = types.StringValue(var31.DefaultValueData)
				var32.Name = types.StringValue(var31.Name)
				var32.Negate = types.BoolValue(var31.Negate)
				var32.RegistryValue = var33
				var30 = append(var30, var32)
			}
		}
		var20.Plist = var21
		var20.ProcessList = var27
		var20.RegistryKey = var30
		var19.Criteria = var20
	}
	var var36 *objectsHipObjectsDsModelDataLossPreventionObject
	if ans.DataLossPrevention != nil {
		var36 = &objectsHipObjectsDsModelDataLossPreventionObject{}
		var var37 *objectsHipObjectsDsModelCriteriaObject3
		if ans.DataLossPrevention.Criteria != nil {
			var37 = &objectsHipObjectsDsModelCriteriaObject3{}
			var37.IsEnabled = types.StringValue(ans.DataLossPrevention.Criteria.IsEnabled)
			var37.IsInstalled = types.BoolValue(ans.DataLossPrevention.Criteria.IsInstalled)
		}
		var var38 []objectsHipObjectsDsModelVendorObject1
		if len(ans.DataLossPrevention.Vendor) != 0 {
			var38 = make([]objectsHipObjectsDsModelVendorObject1, 0, len(ans.DataLossPrevention.Vendor))
			for var39Index := range ans.DataLossPrevention.Vendor {
				var39 := ans.DataLossPrevention.Vendor[var39Index]
				var var40 objectsHipObjectsDsModelVendorObject1
				var40.Name = types.StringValue(var39.Name)
				var40.Product = EncodeStringSlice(var39.Product)
				var38 = append(var38, var40)
			}
		}
		var36.Criteria = var37
		var36.ExcludeVendor = types.BoolValue(ans.DataLossPrevention.ExcludeVendor)
		var36.Vendor = var38
	}
	var var41 *objectsHipObjectsDsModelDiskBackupObject
	if ans.DiskBackup != nil {
		var41 = &objectsHipObjectsDsModelDiskBackupObject{}
		var var42 *objectsHipObjectsDsModelCriteriaObject4
		if ans.DiskBackup.Criteria != nil {
			var42 = &objectsHipObjectsDsModelCriteriaObject4{}
			var var43 *objectsHipObjectsDsModelLastBackupTimeObject
			if ans.DiskBackup.Criteria.LastBackupTime != nil {
				var43 = &objectsHipObjectsDsModelLastBackupTimeObject{}
				var var44 *objectsHipObjectsDsModelNotWithinObject
				if ans.DiskBackup.Criteria.LastBackupTime.NotWithin != nil {
					var44 = &objectsHipObjectsDsModelNotWithinObject{}
					var44.Days = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.NotWithin.Days)
					var44.Hours = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.NotWithin.Hours)
				}
				var var45 *objectsHipObjectsDsModelWithinObject
				if ans.DiskBackup.Criteria.LastBackupTime.Within != nil {
					var45 = &objectsHipObjectsDsModelWithinObject{}
					var45.Days = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.Within.Days)
					var45.Hours = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.Within.Hours)
				}
				if ans.DiskBackup.Criteria.LastBackupTime.NotAvailable != nil {
					var43.NotAvailable = types.BoolValue(true)
				}
				var43.NotWithin = var44
				var43.Within = var45
			}
			var42.IsInstalled = types.BoolValue(ans.DiskBackup.Criteria.IsInstalled)
			var42.LastBackupTime = var43
		}
		var var46 []objectsHipObjectsDsModelVendorObject
		if len(ans.DiskBackup.Vendor) != 0 {
			var46 = make([]objectsHipObjectsDsModelVendorObject, 0, len(ans.DiskBackup.Vendor))
			for var47Index := range ans.DiskBackup.Vendor {
				var47 := ans.DiskBackup.Vendor[var47Index]
				var var48 objectsHipObjectsDsModelVendorObject
				var48.Name = types.StringValue(var47.Name)
				var48.Product = EncodeStringSlice(var47.Product)
				var46 = append(var46, var48)
			}
		}
		var41.Criteria = var42
		var41.ExcludeVendor = types.BoolValue(ans.DiskBackup.ExcludeVendor)
		var41.Vendor = var46
	}
	var var49 *objectsHipObjectsDsModelDiskEncryptionObject
	if ans.DiskEncryption != nil {
		var49 = &objectsHipObjectsDsModelDiskEncryptionObject{}
		var var50 *objectsHipObjectsDsModelCriteriaObject5
		if ans.DiskEncryption.Criteria != nil {
			var50 = &objectsHipObjectsDsModelCriteriaObject5{}
			var var51 []objectsHipObjectsDsModelEncryptedLocationsObject
			if len(ans.DiskEncryption.Criteria.EncryptedLocations) != 0 {
				var51 = make([]objectsHipObjectsDsModelEncryptedLocationsObject, 0, len(ans.DiskEncryption.Criteria.EncryptedLocations))
				for var52Index := range ans.DiskEncryption.Criteria.EncryptedLocations {
					var52 := ans.DiskEncryption.Criteria.EncryptedLocations[var52Index]
					var var53 objectsHipObjectsDsModelEncryptedLocationsObject
					var var54 *objectsHipObjectsDsModelEncryptionStateObject
					if var52.EncryptionState != nil {
						var54 = &objectsHipObjectsDsModelEncryptionStateObject{}
						var54.Is = types.StringValue(var52.EncryptionState.Is)
						var54.IsNot = types.StringValue(var52.EncryptionState.IsNot)
					}
					var53.EncryptionState = var54
					var53.Name = types.StringValue(var52.Name)
					var51 = append(var51, var53)
				}
			}
			var50.EncryptedLocations = var51
			var50.IsInstalled = types.BoolValue(ans.DiskEncryption.Criteria.IsInstalled)
		}
		var var55 []objectsHipObjectsDsModelVendorObject
		if len(ans.DiskEncryption.Vendor) != 0 {
			var55 = make([]objectsHipObjectsDsModelVendorObject, 0, len(ans.DiskEncryption.Vendor))
			for var56Index := range ans.DiskEncryption.Vendor {
				var56 := ans.DiskEncryption.Vendor[var56Index]
				var var57 objectsHipObjectsDsModelVendorObject
				var57.Name = types.StringValue(var56.Name)
				var57.Product = EncodeStringSlice(var56.Product)
				var55 = append(var55, var57)
			}
		}
		var49.Criteria = var50
		var49.ExcludeVendor = types.BoolValue(ans.DiskEncryption.ExcludeVendor)
		var49.Vendor = var55
	}
	var var58 *objectsHipObjectsDsModelFirewallObject
	if ans.Firewall != nil {
		var58 = &objectsHipObjectsDsModelFirewallObject{}
		var var59 *objectsHipObjectsDsModelCriteriaObject3
		if ans.Firewall.Criteria != nil {
			var59 = &objectsHipObjectsDsModelCriteriaObject3{}
			var59.IsEnabled = types.StringValue(ans.Firewall.Criteria.IsEnabled)
			var59.IsInstalled = types.BoolValue(ans.Firewall.Criteria.IsInstalled)
		}
		var var60 []objectsHipObjectsDsModelVendorObject
		if len(ans.Firewall.Vendor) != 0 {
			var60 = make([]objectsHipObjectsDsModelVendorObject, 0, len(ans.Firewall.Vendor))
			for var61Index := range ans.Firewall.Vendor {
				var61 := ans.Firewall.Vendor[var61Index]
				var var62 objectsHipObjectsDsModelVendorObject
				var62.Name = types.StringValue(var61.Name)
				var62.Product = EncodeStringSlice(var61.Product)
				var60 = append(var60, var62)
			}
		}
		var58.Criteria = var59
		var58.ExcludeVendor = types.BoolValue(ans.Firewall.ExcludeVendor)
		var58.Vendor = var60
	}
	var var63 *objectsHipObjectsDsModelHostInfoObject
	if ans.HostInfo != nil {
		var63 = &objectsHipObjectsDsModelHostInfoObject{}
		var var64 objectsHipObjectsDsModelCriteriaObject6
		var var65 *objectsHipObjectsDsModelClientVersionObject
		if ans.HostInfo.Criteria.ClientVersion != nil {
			var65 = &objectsHipObjectsDsModelClientVersionObject{}
			var65.Contains = types.StringValue(ans.HostInfo.Criteria.ClientVersion.Contains)
			var65.Is = types.StringValue(ans.HostInfo.Criteria.ClientVersion.Is)
			var65.IsNot = types.StringValue(ans.HostInfo.Criteria.ClientVersion.IsNot)
		}
		var var66 *objectsHipObjectsDsModelDomainObject
		if ans.HostInfo.Criteria.Domain != nil {
			var66 = &objectsHipObjectsDsModelDomainObject{}
			var66.Contains = types.StringValue(ans.HostInfo.Criteria.Domain.Contains)
			var66.Is = types.StringValue(ans.HostInfo.Criteria.Domain.Is)
			var66.IsNot = types.StringValue(ans.HostInfo.Criteria.Domain.IsNot)
		}
		var var67 *objectsHipObjectsDsModelHostIdObject
		if ans.HostInfo.Criteria.HostId != nil {
			var67 = &objectsHipObjectsDsModelHostIdObject{}
			var67.Contains = types.StringValue(ans.HostInfo.Criteria.HostId.Contains)
			var67.Is = types.StringValue(ans.HostInfo.Criteria.HostId.Is)
			var67.IsNot = types.StringValue(ans.HostInfo.Criteria.HostId.IsNot)
		}
		var var68 *objectsHipObjectsDsModelHostNameObject
		if ans.HostInfo.Criteria.HostName != nil {
			var68 = &objectsHipObjectsDsModelHostNameObject{}
			var68.Contains = types.StringValue(ans.HostInfo.Criteria.HostName.Contains)
			var68.Is = types.StringValue(ans.HostInfo.Criteria.HostName.Is)
			var68.IsNot = types.StringValue(ans.HostInfo.Criteria.HostName.IsNot)
		}
		var var69 *objectsHipObjectsDsModelOsObject
		if ans.HostInfo.Criteria.Os != nil {
			var69 = &objectsHipObjectsDsModelOsObject{}
			var var70 *objectsHipObjectsDsModelContainsObject
			if ans.HostInfo.Criteria.Os.Contains != nil {
				var70 = &objectsHipObjectsDsModelContainsObject{}
				var70.Apple = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Apple)
				var70.Google = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Google)
				var70.Linux = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Linux)
				var70.Microsoft = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Microsoft)
				var70.Other = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Other)
			}
			var69.Contains = var70
		}
		var var71 *objectsHipObjectsDsModelSerialNumberObject
		if ans.HostInfo.Criteria.SerialNumber != nil {
			var71 = &objectsHipObjectsDsModelSerialNumberObject{}
			var71.Contains = types.StringValue(ans.HostInfo.Criteria.SerialNumber.Contains)
			var71.Is = types.StringValue(ans.HostInfo.Criteria.SerialNumber.Is)
			var71.IsNot = types.StringValue(ans.HostInfo.Criteria.SerialNumber.IsNot)
		}
		var64.ClientVersion = var65
		var64.Domain = var66
		var64.HostId = var67
		var64.HostName = var68
		var64.Managed = types.BoolValue(ans.HostInfo.Criteria.Managed)
		var64.Os = var69
		var64.SerialNumber = var71
		var63.Criteria = var64
	}
	var var72 *objectsHipObjectsDsModelMobileDeviceObject
	if ans.MobileDevice != nil {
		var72 = &objectsHipObjectsDsModelMobileDeviceObject{}
		var var73 *objectsHipObjectsDsModelCriteriaObject7
		if ans.MobileDevice.Criteria != nil {
			var73 = &objectsHipObjectsDsModelCriteriaObject7{}
			var var74 *objectsHipObjectsDsModelApplicationsObject
			if ans.MobileDevice.Criteria.Applications != nil {
				var74 = &objectsHipObjectsDsModelApplicationsObject{}
				var var75 *objectsHipObjectsDsModelHasMalwareObject
				if ans.MobileDevice.Criteria.Applications.HasMalware != nil {
					var75 = &objectsHipObjectsDsModelHasMalwareObject{}
					var var76 *objectsHipObjectsDsModelYesObject
					if ans.MobileDevice.Criteria.Applications.HasMalware.Yes != nil {
						var76 = &objectsHipObjectsDsModelYesObject{}
						var var77 []objectsHipObjectsDsModelExcludesObject
						if len(ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes) != 0 {
							var77 = make([]objectsHipObjectsDsModelExcludesObject, 0, len(ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes))
							for var78Index := range ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes {
								var78 := ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes[var78Index]
								var var79 objectsHipObjectsDsModelExcludesObject
								var79.Hash = types.StringValue(var78.Hash)
								var79.Name = types.StringValue(var78.Name)
								var79.Package = types.StringValue(var78.Package)
								var77 = append(var77, var79)
							}
						}
						var76.Excludes = var77
					}
					if ans.MobileDevice.Criteria.Applications.HasMalware.No != nil {
						var75.No = types.BoolValue(true)
					}
					var75.Yes = var76
				}
				var var80 []objectsHipObjectsDsModelIncludesObject
				if len(ans.MobileDevice.Criteria.Applications.Includes) != 0 {
					var80 = make([]objectsHipObjectsDsModelIncludesObject, 0, len(ans.MobileDevice.Criteria.Applications.Includes))
					for var81Index := range ans.MobileDevice.Criteria.Applications.Includes {
						var81 := ans.MobileDevice.Criteria.Applications.Includes[var81Index]
						var var82 objectsHipObjectsDsModelIncludesObject
						var82.Hash = types.StringValue(var81.Hash)
						var82.Name = types.StringValue(var81.Name)
						var82.Package = types.StringValue(var81.Package)
						var80 = append(var80, var82)
					}
				}
				var74.HasMalware = var75
				var74.HasUnmanagedApp = types.BoolValue(ans.MobileDevice.Criteria.Applications.HasUnmanagedApp)
				var74.Includes = var80
			}
			var var83 *objectsHipObjectsDsModelImeiObject
			if ans.MobileDevice.Criteria.Imei != nil {
				var83 = &objectsHipObjectsDsModelImeiObject{}
				var83.Contains = types.StringValue(ans.MobileDevice.Criteria.Imei.Contains)
				var83.Is = types.StringValue(ans.MobileDevice.Criteria.Imei.Is)
				var83.IsNot = types.StringValue(ans.MobileDevice.Criteria.Imei.IsNot)
			}
			var var84 *objectsHipObjectsDsModelLastCheckinTimeObject
			if ans.MobileDevice.Criteria.LastCheckinTime != nil {
				var84 = &objectsHipObjectsDsModelLastCheckinTimeObject{}
				var var85 *objectsHipObjectsDsModelNotWithinObject3
				if ans.MobileDevice.Criteria.LastCheckinTime.NotWithin != nil {
					var85 = &objectsHipObjectsDsModelNotWithinObject3{}
					var85.Days = types.Int64Value(ans.MobileDevice.Criteria.LastCheckinTime.NotWithin.Days)
				}
				var var86 *objectsHipObjectsDsModelWithinObject3
				if ans.MobileDevice.Criteria.LastCheckinTime.Within != nil {
					var86 = &objectsHipObjectsDsModelWithinObject3{}
					var86.Days = types.Int64Value(ans.MobileDevice.Criteria.LastCheckinTime.Within.Days)
				}
				var84.NotWithin = var85
				var84.Within = var86
			}
			var var87 *objectsHipObjectsDsModelModelObject
			if ans.MobileDevice.Criteria.Model != nil {
				var87 = &objectsHipObjectsDsModelModelObject{}
				var87.Contains = types.StringValue(ans.MobileDevice.Criteria.Model.Contains)
				var87.Is = types.StringValue(ans.MobileDevice.Criteria.Model.Is)
				var87.IsNot = types.StringValue(ans.MobileDevice.Criteria.Model.IsNot)
			}
			var var88 *objectsHipObjectsDsModelPhoneNumberObject
			if ans.MobileDevice.Criteria.PhoneNumber != nil {
				var88 = &objectsHipObjectsDsModelPhoneNumberObject{}
				var88.Contains = types.StringValue(ans.MobileDevice.Criteria.PhoneNumber.Contains)
				var88.Is = types.StringValue(ans.MobileDevice.Criteria.PhoneNumber.Is)
				var88.IsNot = types.StringValue(ans.MobileDevice.Criteria.PhoneNumber.IsNot)
			}
			var var89 *objectsHipObjectsDsModelTagObject
			if ans.MobileDevice.Criteria.Tag != nil {
				var89 = &objectsHipObjectsDsModelTagObject{}
				var89.Contains = types.StringValue(ans.MobileDevice.Criteria.Tag.Contains)
				var89.Is = types.StringValue(ans.MobileDevice.Criteria.Tag.Is)
				var89.IsNot = types.StringValue(ans.MobileDevice.Criteria.Tag.IsNot)
			}
			var73.Applications = var74
			var73.DiskEncrypted = types.BoolValue(ans.MobileDevice.Criteria.DiskEncrypted)
			var73.Imei = var83
			var73.Jailbroken = types.BoolValue(ans.MobileDevice.Criteria.Jailbroken)
			var73.LastCheckinTime = var84
			var73.Model = var87
			var73.PasscodeSet = types.BoolValue(ans.MobileDevice.Criteria.PasscodeSet)
			var73.PhoneNumber = var88
			var73.Tag = var89
		}
		var72.Criteria = var73
	}
	var var90 *objectsHipObjectsDsModelNetworkInfoObject
	if ans.NetworkInfo != nil {
		var90 = &objectsHipObjectsDsModelNetworkInfoObject{}
		var var91 *objectsHipObjectsDsModelCriteriaObject8
		if ans.NetworkInfo.Criteria != nil {
			var91 = &objectsHipObjectsDsModelCriteriaObject8{}
			var var92 *objectsHipObjectsDsModelNetworkObject
			if ans.NetworkInfo.Criteria.Network != nil {
				var92 = &objectsHipObjectsDsModelNetworkObject{}
				var var93 *objectsHipObjectsDsModelIsObject
				if ans.NetworkInfo.Criteria.Network.Is != nil {
					var93 = &objectsHipObjectsDsModelIsObject{}
					var var94 *objectsHipObjectsDsModelMobileObject
					if ans.NetworkInfo.Criteria.Network.Is.Mobile != nil {
						var94 = &objectsHipObjectsDsModelMobileObject{}
						var94.Carrier = types.StringValue(ans.NetworkInfo.Criteria.Network.Is.Mobile.Carrier)
					}
					var var95 *objectsHipObjectsDsModelWifiObject
					if ans.NetworkInfo.Criteria.Network.Is.Wifi != nil {
						var95 = &objectsHipObjectsDsModelWifiObject{}
						var95.Ssid = types.StringValue(ans.NetworkInfo.Criteria.Network.Is.Wifi.Ssid)
					}
					var93.Mobile = var94
					if ans.NetworkInfo.Criteria.Network.Is.Unknown != nil {
						var93.Unknown = types.BoolValue(true)
					}
					var93.Wifi = var95
				}
				var var96 *objectsHipObjectsDsModelIsNotObject
				if ans.NetworkInfo.Criteria.Network.IsNot != nil {
					var96 = &objectsHipObjectsDsModelIsNotObject{}
					var var97 *objectsHipObjectsDsModelMobileObject
					if ans.NetworkInfo.Criteria.Network.IsNot.Mobile != nil {
						var97 = &objectsHipObjectsDsModelMobileObject{}
						var97.Carrier = types.StringValue(ans.NetworkInfo.Criteria.Network.IsNot.Mobile.Carrier)
					}
					var var98 *objectsHipObjectsDsModelWifiObject
					if ans.NetworkInfo.Criteria.Network.IsNot.Wifi != nil {
						var98 = &objectsHipObjectsDsModelWifiObject{}
						var98.Ssid = types.StringValue(ans.NetworkInfo.Criteria.Network.IsNot.Wifi.Ssid)
					}
					if ans.NetworkInfo.Criteria.Network.IsNot.Ethernet != nil {
						var96.Ethernet = types.BoolValue(true)
					}
					var96.Mobile = var97
					if ans.NetworkInfo.Criteria.Network.IsNot.Unknown != nil {
						var96.Unknown = types.BoolValue(true)
					}
					var96.Wifi = var98
				}
				var92.Is = var93
				var92.IsNot = var96
			}
			var91.Network = var92
		}
		var90.Criteria = var91
	}
	var var99 *objectsHipObjectsDsModelPatchManagementObject
	if ans.PatchManagement != nil {
		var99 = &objectsHipObjectsDsModelPatchManagementObject{}
		var var100 *objectsHipObjectsDsModelCriteriaObject9
		if ans.PatchManagement.Criteria != nil {
			var100 = &objectsHipObjectsDsModelCriteriaObject9{}
			var var101 *objectsHipObjectsDsModelMissingPatchesObject
			if ans.PatchManagement.Criteria.MissingPatches != nil {
				var101 = &objectsHipObjectsDsModelMissingPatchesObject{}
				var var102 *objectsHipObjectsDsModelSeverityObject
				if ans.PatchManagement.Criteria.MissingPatches.Severity != nil {
					var102 = &objectsHipObjectsDsModelSeverityObject{}
					var102.GreaterEqual = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.GreaterEqual)
					var102.GreaterThan = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.GreaterThan)
					var102.Is = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.Is)
					var102.IsNot = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.IsNot)
					var102.LessEqual = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.LessEqual)
					var102.LessThan = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.LessThan)
				}
				var101.Check = types.StringValue(ans.PatchManagement.Criteria.MissingPatches.Check)
				var101.Patches = EncodeStringSlice(ans.PatchManagement.Criteria.MissingPatches.Patches)
				var101.Severity = var102
			}
			var100.IsEnabled = types.StringValue(ans.PatchManagement.Criteria.IsEnabled)
			var100.IsInstalled = types.BoolValue(ans.PatchManagement.Criteria.IsInstalled)
			var100.MissingPatches = var101
		}
		var var103 []objectsHipObjectsDsModelVendorObject1
		if len(ans.PatchManagement.Vendor) != 0 {
			var103 = make([]objectsHipObjectsDsModelVendorObject1, 0, len(ans.PatchManagement.Vendor))
			for var104Index := range ans.PatchManagement.Vendor {
				var104 := ans.PatchManagement.Vendor[var104Index]
				var var105 objectsHipObjectsDsModelVendorObject1
				var105.Name = types.StringValue(var104.Name)
				var105.Product = EncodeStringSlice(var104.Product)
				var103 = append(var103, var105)
			}
		}
		var99.Criteria = var100
		var99.ExcludeVendor = types.BoolValue(ans.PatchManagement.ExcludeVendor)
		var99.Vendor = var103
	}
	state.AntiMalware = var0
	state.Certificate = var14
	state.CustomChecks = var19
	state.DataLossPrevention = var36
	state.Description = types.StringValue(ans.Description)
	state.DiskBackup = var41
	state.DiskEncryption = var49
	state.Firewall = var58
	state.HostInfo = var63
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MobileDevice = var72
	state.Name = types.StringValue(ans.Name)
	state.NetworkInfo = var90
	state.PatchManagement = var99

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &objectsHipObjectsResource{}
	_ resource.ResourceWithConfigure   = &objectsHipObjectsResource{}
	_ resource.ResourceWithImportState = &objectsHipObjectsResource{}
)

func NewObjectsHipObjectsResource() resource.Resource {
	return &objectsHipObjectsResource{}
}

type objectsHipObjectsResource struct {
	client *sase.Client
}

type objectsHipObjectsRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/objects-hip-objects
	AntiMalware        *objectsHipObjectsRsModelAntiMalwareObject        `tfsdk:"anti_malware"`
	Certificate        *objectsHipObjectsRsModelCertificateObject        `tfsdk:"certificate"`
	CustomChecks       *objectsHipObjectsRsModelCustomChecksObject       `tfsdk:"custom_checks"`
	DataLossPrevention *objectsHipObjectsRsModelDataLossPreventionObject `tfsdk:"data_loss_prevention"`
	Description        types.String                                      `tfsdk:"description"`
	DiskBackup         *objectsHipObjectsRsModelDiskBackupObject         `tfsdk:"disk_backup"`
	DiskEncryption     *objectsHipObjectsRsModelDiskEncryptionObject     `tfsdk:"disk_encryption"`
	Firewall           *objectsHipObjectsRsModelFirewallObject           `tfsdk:"firewall"`
	HostInfo           *objectsHipObjectsRsModelHostInfoObject           `tfsdk:"host_info"`
	ObjectId           types.String                                      `tfsdk:"object_id"`
	MobileDevice       *objectsHipObjectsRsModelMobileDeviceObject       `tfsdk:"mobile_device"`
	Name               types.String                                      `tfsdk:"name"`
	NetworkInfo        *objectsHipObjectsRsModelNetworkInfoObject        `tfsdk:"network_info"`
	PatchManagement    *objectsHipObjectsRsModelPatchManagementObject    `tfsdk:"patch_management"`
}

type objectsHipObjectsRsModelAntiMalwareObject struct {
	Criteria      *objectsHipObjectsRsModelCriteriaObject `tfsdk:"criteria"`
	ExcludeVendor types.Bool                              `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsRsModelVendorObject  `tfsdk:"vendor"`
}

type objectsHipObjectsRsModelCriteriaObject struct {
	IsInstalled        types.Bool                                    `tfsdk:"is_installed"`
	LastScanTime       *objectsHipObjectsRsModelLastScanTimeObject   `tfsdk:"last_scan_time"`
	ProductVersion     *objectsHipObjectsRsModelProductVersionObject `tfsdk:"product_version"`
	RealTimeProtection types.String                                  `tfsdk:"real_time_protection"`
	VirdefVersion      *objectsHipObjectsRsModelVirdefVersionObject  `tfsdk:"virdef_version"`
}

type objectsHipObjectsRsModelLastScanTimeObject struct {
	NotAvailable types.Bool                               `tfsdk:"not_available"`
	NotWithin    *objectsHipObjectsRsModelNotWithinObject `tfsdk:"not_within"`
	Within       *objectsHipObjectsRsModelWithinObject    `tfsdk:"within"`
}

type objectsHipObjectsRsModelNotWithinObject struct {
	Days  types.Int64 `tfsdk:"days"`
	Hours types.Int64 `tfsdk:"hours"`
}

type objectsHipObjectsRsModelWithinObject struct {
	Days  types.Int64 `tfsdk:"days"`
	Hours types.Int64 `tfsdk:"hours"`
}

type objectsHipObjectsRsModelProductVersionObject struct {
	Contains     types.String                              `tfsdk:"contains"`
	GreaterEqual types.String                              `tfsdk:"greater_equal"`
	GreaterThan  types.String                              `tfsdk:"greater_than"`
	Is           types.String                              `tfsdk:"is"`
	IsNot        types.String                              `tfsdk:"is_not"`
	LessEqual    types.String                              `tfsdk:"less_equal"`
	LessThan     types.String                              `tfsdk:"less_than"`
	NotWithin    *objectsHipObjectsRsModelNotWithinObject1 `tfsdk:"not_within"`
	Within       *objectsHipObjectsRsModelWithinObject1    `tfsdk:"within"`
}

type objectsHipObjectsRsModelNotWithinObject1 struct {
	Versions types.Int64 `tfsdk:"versions"`
}

type objectsHipObjectsRsModelWithinObject1 struct {
	Versions types.Int64 `tfsdk:"versions"`
}

type objectsHipObjectsRsModelVirdefVersionObject struct {
	NotWithin *objectsHipObjectsRsModelNotWithinObject2 `tfsdk:"not_within"`
	Within    *objectsHipObjectsRsModelWithinObject2    `tfsdk:"within"`
}

type objectsHipObjectsRsModelNotWithinObject2 struct {
	Days     types.Int64 `tfsdk:"days"`
	Versions types.Int64 `tfsdk:"versions"`
}

type objectsHipObjectsRsModelWithinObject2 struct {
	Days     types.Int64 `tfsdk:"days"`
	Versions types.Int64 `tfsdk:"versions"`
}

type objectsHipObjectsRsModelVendorObject struct {
	Name    types.String   `tfsdk:"name"`
	Product []types.String `tfsdk:"product"`
}

type objectsHipObjectsRsModelCertificateObject struct {
	Criteria *objectsHipObjectsRsModelCriteriaObject1 `tfsdk:"criteria"`
}

type objectsHipObjectsRsModelCriteriaObject1 struct {
	CertificateAttributes []objectsHipObjectsRsModelCertificateAttributesObject `tfsdk:"certificate_attributes"`
	CertificateProfile    types.String                                          `tfsdk:"certificate_profile"`
}

type objectsHipObjectsRsModelCertificateAttributesObject struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type objectsHipObjectsRsModelCustomChecksObject struct {
	Criteria objectsHipObjectsRsModelCriteriaObject2 `tfsdk:"criteria"`
}

type objectsHipObjectsRsModelCriteriaObject2 struct {
	Plist       []objectsHipObjectsRsModelPlistObject       `tfsdk:"plist"`
	ProcessList []objectsHipObjectsRsModelProcessListObject `tfsdk:"process_list"`
	RegistryKey []objectsHipObjectsRsModelRegistryKeyObject `tfsdk:"registry_key"`
}

type objectsHipObjectsRsModelPlistObject struct {
	Key    []objectsHipObjectsRsModelKeyObject `tfsdk:"key"`
	Name   types.String                        `tfsdk:"name"`
	Negate types.Bool                          `tfsdk:"negate"`
}

type objectsHipObjectsRsModelKeyObject struct {
	Name   types.String `tfsdk:"name"`
	Negate types.Bool   `tfsdk:"negate"`
	Value  types.String `tfsdk:"value"`
}

type objectsHipObjectsRsModelProcessListObject struct {
	Name    types.String `tfsdk:"name"`
	Running types.Bool   `tfsdk:"running"`
}

type objectsHipObjectsRsModelRegistryKeyObject struct {
	DefaultValueData types.String                                  `tfsdk:"default_value_data"`
	Name             types.String                                  `tfsdk:"name"`
	Negate           types.Bool                                    `tfsdk:"negate"`
	RegistryValue    []objectsHipObjectsRsModelRegistryValueObject `tfsdk:"registry_value"`
}

type objectsHipObjectsRsModelRegistryValueObject struct {
	Name      types.String `tfsdk:"name"`
	Negate    types.Bool   `tfsdk:"negate"`
	ValueData types.String `tfsdk:"value_data"`
}

type objectsHipObjectsRsModelDataLossPreventionObject struct {
	Criteria      *objectsHipObjectsRsModelCriteriaObject3 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                               `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsRsModelVendorObject1  `tfsdk:"vendor"`
}

type objectsHipObjectsRsModelCriteriaObject3 struct {
	IsEnabled   types.String `tfsdk:"is_enabled"`
	IsInstalled types.Bool   `tfsdk:"is_installed"`
}

type objectsHipObjectsRsModelVendorObject1 struct {
	Name    types.String   `tfsdk:"name"`
	Product []types.String `tfsdk:"product"`
}

type objectsHipObjectsRsModelDiskBackupObject struct {
	Criteria      *objectsHipObjectsRsModelCriteriaObject4 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                               `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsRsModelVendorObject   `tfsdk:"vendor"`
}

type objectsHipObjectsRsModelCriteriaObject4 struct {
	IsInstalled    types.Bool                                    `tfsdk:"is_installed"`
	LastBackupTime *objectsHipObjectsRsModelLastBackupTimeObject `tfsdk:"last_backup_time"`
}

type objectsHipObjectsRsModelLastBackupTimeObject struct {
	NotAvailable types.Bool                               `tfsdk:"not_available"`
	NotWithin    *objectsHipObjectsRsModelNotWithinObject `tfsdk:"not_within"`
	Within       *objectsHipObjectsRsModelWithinObject    `tfsdk:"within"`
}

type objectsHipObjectsRsModelDiskEncryptionObject struct {
	Criteria      *objectsHipObjectsRsModelCriteriaObject5 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                               `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsRsModelVendorObject   `tfsdk:"vendor"`
}

type objectsHipObjectsRsModelCriteriaObject5 struct {
	EncryptedLocations []objectsHipObjectsRsModelEncryptedLocationsObject `tfsdk:"encrypted_locations"`
	IsInstalled        types.Bool                                         `tfsdk:"is_installed"`
}

type objectsHipObjectsRsModelEncryptedLocationsObject struct {
	EncryptionState *objectsHipObjectsRsModelEncryptionStateObject `tfsdk:"encryption_state"`
	Name            types.String                                   `tfsdk:"name"`
}

type objectsHipObjectsRsModelEncryptionStateObject struct {
	Is    types.String `tfsdk:"is"`
	IsNot types.String `tfsdk:"is_not"`
}

type objectsHipObjectsRsModelFirewallObject struct {
	Criteria      *objectsHipObjectsRsModelCriteriaObject3 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                               `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsRsModelVendorObject   `tfsdk:"vendor"`
}

type objectsHipObjectsRsModelHostInfoObject struct {
	Criteria objectsHipObjectsRsModelCriteriaObject6 `tfsdk:"criteria"`
}

type objectsHipObjectsRsModelCriteriaObject6 struct {
	ClientVersion *objectsHipObjectsRsModelClientVersionObject `tfsdk:"client_version"`
	Domain        *objectsHipObjectsRsModelDomainObject        `tfsdk:"domain"`
	HostId        *objectsHipObjectsRsModelHostIdObject        `tfsdk:"host_id"`
	HostName      *objectsHipObjectsRsModelHostNameObject      `tfsdk:"host_name"`
	Managed       types.Bool                                   `tfsdk:"managed"`
	Os            *objectsHipObjectsRsModelOsObject            `tfsdk:"os"`
	SerialNumber  *objectsHipObjectsRsModelSerialNumberObject  `tfsdk:"serial_number"`
}

type objectsHipObjectsRsModelClientVersionObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsRsModelDomainObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsRsModelHostIdObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsRsModelHostNameObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsRsModelOsObject struct {
	Contains *objectsHipObjectsRsModelContainsObject `tfsdk:"contains"`
}

type objectsHipObjectsRsModelContainsObject struct {
	Apple     types.String `tfsdk:"apple"`
	Google    types.String `tfsdk:"google"`
	Linux     types.String `tfsdk:"linux"`
	Microsoft types.String `tfsdk:"microsoft"`
	Other     types.String `tfsdk:"other"`
}

type objectsHipObjectsRsModelSerialNumberObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsRsModelMobileDeviceObject struct {
	Criteria *objectsHipObjectsRsModelCriteriaObject7 `tfsdk:"criteria"`
}

type objectsHipObjectsRsModelCriteriaObject7 struct {
	Applications    *objectsHipObjectsRsModelApplicationsObject    `tfsdk:"applications"`
	DiskEncrypted   types.Bool                                     `tfsdk:"disk_encrypted"`
	Imei            *objectsHipObjectsRsModelImeiObject            `tfsdk:"imei"`
	Jailbroken      types.Bool                                     `tfsdk:"jailbroken"`
	LastCheckinTime *objectsHipObjectsRsModelLastCheckinTimeObject `tfsdk:"last_checkin_time"`
	Model           *objectsHipObjectsRsModelModelObject           `tfsdk:"model"`
	PasscodeSet     types.Bool                                     `tfsdk:"passcode_set"`
	PhoneNumber     *objectsHipObjectsRsModelPhoneNumberObject     `tfsdk:"phone_number"`
	Tag             *objectsHipObjectsRsModelTagObject             `tfsdk:"tag"`
}

type objectsHipObjectsRsModelApplicationsObject struct {
	HasMalware      *objectsHipObjectsRsModelHasMalwareObject `tfsdk:"has_malware"`
	HasUnmanagedApp types.Bool                                `tfsdk:"has_unmanaged_app"`
	Includes        []objectsHipObjectsRsModelIncludesObject  `tfsdk:"includes"`
}

type objectsHipObjectsRsModelHasMalwareObject struct {
	No  types.Bool                         `tfsdk:"no"`
	Yes *objectsHipObjectsRsModelYesObject `tfsdk:"yes"`
}

type objectsHipObjectsRsModelYesObject struct {
	Excludes []objectsHipObjectsRsModelExcludesObject `tfsdk:"excludes"`
}

type objectsHipObjectsRsModelExcludesObject struct {
	Hash    types.String `tfsdk:"hash"`
	Name    types.String `tfsdk:"name"`
	Package types.String `tfsdk:"package"`
}

type objectsHipObjectsRsModelIncludesObject struct {
	Hash    types.String `tfsdk:"hash"`
	Name    types.String `tfsdk:"name"`
	Package types.String `tfsdk:"package"`
}

type objectsHipObjectsRsModelImeiObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsRsModelLastCheckinTimeObject struct {
	NotWithin *objectsHipObjectsRsModelNotWithinObject3 `tfsdk:"not_within"`
	Within    *objectsHipObjectsRsModelWithinObject3    `tfsdk:"within"`
}

type objectsHipObjectsRsModelNotWithinObject3 struct {
	Days types.Int64 `tfsdk:"days"`
}

type objectsHipObjectsRsModelWithinObject3 struct {
	Days types.Int64 `tfsdk:"days"`
}

type objectsHipObjectsRsModelModelObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsRsModelPhoneNumberObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsRsModelTagObject struct {
	Contains types.String `tfsdk:"contains"`
	Is       types.String `tfsdk:"is"`
	IsNot    types.String `tfsdk:"is_not"`
}

type objectsHipObjectsRsModelNetworkInfoObject struct {
	Criteria *objectsHipObjectsRsModelCriteriaObject8 `tfsdk:"criteria"`
}

type objectsHipObjectsRsModelCriteriaObject8 struct {
	Network *objectsHipObjectsRsModelNetworkObject `tfsdk:"network"`
}

type objectsHipObjectsRsModelNetworkObject struct {
	Is    *objectsHipObjectsRsModelIsObject    `tfsdk:"is"`
	IsNot *objectsHipObjectsRsModelIsNotObject `tfsdk:"is_not"`
}

type objectsHipObjectsRsModelIsObject struct {
	Mobile  *objectsHipObjectsRsModelMobileObject `tfsdk:"mobile"`
	Unknown types.Bool                            `tfsdk:"unknown"`
	Wifi    *objectsHipObjectsRsModelWifiObject   `tfsdk:"wifi"`
}

type objectsHipObjectsRsModelMobileObject struct {
	Carrier types.String `tfsdk:"carrier"`
}

type objectsHipObjectsRsModelWifiObject struct {
	Ssid types.String `tfsdk:"ssid"`
}

type objectsHipObjectsRsModelIsNotObject struct {
	Ethernet types.Bool                            `tfsdk:"ethernet"`
	Mobile   *objectsHipObjectsRsModelMobileObject `tfsdk:"mobile"`
	Unknown  types.Bool                            `tfsdk:"unknown"`
	Wifi     *objectsHipObjectsRsModelWifiObject   `tfsdk:"wifi"`
}

type objectsHipObjectsRsModelPatchManagementObject struct {
	Criteria      *objectsHipObjectsRsModelCriteriaObject9 `tfsdk:"criteria"`
	ExcludeVendor types.Bool                               `tfsdk:"exclude_vendor"`
	Vendor        []objectsHipObjectsRsModelVendorObject1  `tfsdk:"vendor"`
}

type objectsHipObjectsRsModelCriteriaObject9 struct {
	IsEnabled      types.String                                  `tfsdk:"is_enabled"`
	IsInstalled    types.Bool                                    `tfsdk:"is_installed"`
	MissingPatches *objectsHipObjectsRsModelMissingPatchesObject `tfsdk:"missing_patches"`
}

type objectsHipObjectsRsModelMissingPatchesObject struct {
	Check    types.String                            `tfsdk:"check"`
	Patches  []types.String                          `tfsdk:"patches"`
	Severity *objectsHipObjectsRsModelSeverityObject `tfsdk:"severity"`
}

type objectsHipObjectsRsModelSeverityObject struct {
	GreaterEqual types.Int64 `tfsdk:"greater_equal"`
	GreaterThan  types.Int64 `tfsdk:"greater_than"`
	Is           types.Int64 `tfsdk:"is"`
	IsNot        types.Int64 `tfsdk:"is_not"`
	LessEqual    types.Int64 `tfsdk:"less_equal"`
	LessThan     types.Int64 `tfsdk:"less_than"`
}

// Metadata returns the data source type name.
func (r *objectsHipObjectsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_objects_hip_objects"
}

// Schema defines the schema for this listing data source.
func (r *objectsHipObjectsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rsschema.Schema{
		Description: "Retrieves config for a specific item.",

		Attributes: map[string]rsschema.Attribute{
			"id": rsschema.StringAttribute{
				Description: "The object ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			// Input.
			"folder": rsschema.StringAttribute{
				Description: "The folder of the entry",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			"anti_malware": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"criteria": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"is_installed": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(true),
								},
							},
							"last_scan_time": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"not_available": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.Bool{
											boolvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("not_within"),
												path.MatchRelative().AtParent().AtName("within"),
											),
										},
									},
									"not_within": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"days": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(1),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
											"hours": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(24),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
										},
									},
									"within": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"days": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(1),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
											"hours": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(24),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
										},
									},
								},
							},
							"product_version": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"contains": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("greater_equal"),
												path.MatchRelative().AtParent().AtName("greater_than"),
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
												path.MatchRelative().AtParent().AtName("less_equal"),
												path.MatchRelative().AtParent().AtName("less_than"),
												path.MatchRelative().AtParent().AtName("not_within"),
												path.MatchRelative().AtParent().AtName("within"),
											),
										},
									},
									"greater_equal": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("greater_than"),
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
												path.MatchRelative().AtParent().AtName("less_equal"),
												path.MatchRelative().AtParent().AtName("less_than"),
												path.MatchRelative().AtParent().AtName("not_within"),
												path.MatchRelative().AtParent().AtName("within"),
											),
										},
									},
									"greater_than": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("greater_equal"),
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
												path.MatchRelative().AtParent().AtName("less_equal"),
												path.MatchRelative().AtParent().AtName("less_than"),
												path.MatchRelative().AtParent().AtName("not_within"),
												path.MatchRelative().AtParent().AtName("within"),
											),
										},
									},
									"is": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("greater_equal"),
												path.MatchRelative().AtParent().AtName("greater_than"),
												path.MatchRelative().AtParent().AtName("is_not"),
												path.MatchRelative().AtParent().AtName("less_equal"),
												path.MatchRelative().AtParent().AtName("less_than"),
												path.MatchRelative().AtParent().AtName("not_within"),
												path.MatchRelative().AtParent().AtName("within"),
											),
										},
									},
									"is_not": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("greater_equal"),
												path.MatchRelative().AtParent().AtName("greater_than"),
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("less_equal"),
												path.MatchRelative().AtParent().AtName("less_than"),
												path.MatchRelative().AtParent().AtName("not_within"),
												path.MatchRelative().AtParent().AtName("within"),
											),
										},
									},
									"less_equal": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("greater_equal"),
												path.MatchRelative().AtParent().AtName("greater_than"),
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
												path.MatchRelative().AtParent().AtName("less_than"),
												path.MatchRelative().AtParent().AtName("not_within"),
												path.MatchRelative().AtParent().AtName("within"),
											),
										},
									},
									"less_than": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("greater_equal"),
												path.MatchRelative().AtParent().AtName("greater_than"),
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
												path.MatchRelative().AtParent().AtName("less_equal"),
												path.MatchRelative().AtParent().AtName("not_within"),
												path.MatchRelative().AtParent().AtName("within"),
											),
										},
									},
									"not_within": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"versions": rsschema.Int64Attribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(1),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
										},
									},
									"within": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"versions": rsschema.Int64Attribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(1),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
										},
									},
								},
							},
							"real_time_protection": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.OneOf("no", "yes", "not-available"),
								},
							},
							"virdef_version": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"not_within": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"days": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(1),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
											"versions": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(1),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
										},
									},
									"within": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"days": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(1),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
											"versions": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(1),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
										},
									},
								},
							},
						},
					},
					"exclude_vendor": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"vendor": rsschema.ListNestedAttribute{
						Description: "",
						Optional:    true,
						NestedObject: rsschema.NestedAttributeObject{
							Attributes: map[string]rsschema.Attribute{
								"name": rsschema.StringAttribute{
									Description: "",
									Required:    true,
									PlanModifiers: []planmodifier.String{
										DefaultString(""),
									},
									Validators: []validator.String{
										stringvalidator.LengthAtMost(103),
									},
								},
								"product": rsschema.ListAttribute{
									Description: "",
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"certificate": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"criteria": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"certificate_attributes": rsschema.ListNestedAttribute{
								Description: "",
								Optional:    true,
								NestedObject: rsschema.NestedAttributeObject{
									Attributes: map[string]rsschema.Attribute{
										"name": rsschema.StringAttribute{
											Description: "",
											Required:    true,
											PlanModifiers: []planmodifier.String{
												DefaultString(""),
											},
										},
										"value": rsschema.StringAttribute{
											Description: "",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												DefaultString(""),
											},
											Validators: []validator.String{
												stringvalidator.LengthAtMost(1024),
											},
										},
									},
								},
							},
							"certificate_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
				},
			},
			"custom_checks": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"criteria": rsschema.SingleNestedAttribute{
						Description: "",
						Required:    true,
						Attributes: map[string]rsschema.Attribute{
							"plist": rsschema.ListNestedAttribute{
								Description: "",
								Optional:    true,
								NestedObject: rsschema.NestedAttributeObject{
									Attributes: map[string]rsschema.Attribute{
										"key": rsschema.ListNestedAttribute{
											Description: "",
											Optional:    true,
											NestedObject: rsschema.NestedAttributeObject{
												Attributes: map[string]rsschema.Attribute{
													"name": rsschema.StringAttribute{
														Description: "",
														Required:    true,
														PlanModifiers: []planmodifier.String{
															DefaultString(""),
														},
														Validators: []validator.String{
															stringvalidator.LengthAtMost(1023),
														},
													},
													"negate": rsschema.BoolAttribute{
														Description: "",
														Optional:    true,
														Computed:    true,
														PlanModifiers: []planmodifier.Bool{
															DefaultBool(false),
														},
													},
													"value": rsschema.StringAttribute{
														Description: "",
														Optional:    true,
														Computed:    true,
														PlanModifiers: []planmodifier.String{
															DefaultString(""),
														},
														Validators: []validator.String{
															stringvalidator.LengthAtMost(1024),
														},
													},
												},
											},
										},
										"name": rsschema.StringAttribute{
											Description: "",
											Required:    true,
											PlanModifiers: []planmodifier.String{
												DefaultString(""),
											},
											Validators: []validator.String{
												stringvalidator.LengthAtMost(1023),
											},
										},
										"negate": rsschema.BoolAttribute{
											Description: "",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.Bool{
												DefaultBool(false),
											},
										},
									},
								},
							},
							"process_list": rsschema.ListNestedAttribute{
								Description: "",
								Optional:    true,
								NestedObject: rsschema.NestedAttributeObject{
									Attributes: map[string]rsschema.Attribute{
										"name": rsschema.StringAttribute{
											Description: "",
											Required:    true,
											PlanModifiers: []planmodifier.String{
												DefaultString(""),
											},
											Validators: []validator.String{
												stringvalidator.LengthAtMost(1023),
											},
										},
										"running": rsschema.BoolAttribute{
											Description: "",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.Bool{
												DefaultBool(true),
											},
										},
									},
								},
							},
							"registry_key": rsschema.ListNestedAttribute{
								Description: "",
								Optional:    true,
								NestedObject: rsschema.NestedAttributeObject{
									Attributes: map[string]rsschema.Attribute{
										"default_value_data": rsschema.StringAttribute{
											Description: "",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.String{
												DefaultString(""),
											},
											Validators: []validator.String{
												stringvalidator.LengthAtMost(1024),
											},
										},
										"name": rsschema.StringAttribute{
											Description: "",
											Required:    true,
											PlanModifiers: []planmodifier.String{
												DefaultString(""),
											},
											Validators: []validator.String{
												stringvalidator.LengthAtMost(1023),
											},
										},
										"negate": rsschema.BoolAttribute{
											Description: "",
											Optional:    true,
											Computed:    true,
											PlanModifiers: []planmodifier.Bool{
												DefaultBool(false),
											},
										},
										"registry_value": rsschema.ListNestedAttribute{
											Description: "",
											Optional:    true,
											NestedObject: rsschema.NestedAttributeObject{
												Attributes: map[string]rsschema.Attribute{
													"name": rsschema.StringAttribute{
														Description: "",
														Required:    true,
														PlanModifiers: []planmodifier.String{
															DefaultString(""),
														},
														Validators: []validator.String{
															stringvalidator.LengthAtMost(1023),
														},
													},
													"negate": rsschema.BoolAttribute{
														Description: "",
														Optional:    true,
														Computed:    true,
														PlanModifiers: []planmodifier.Bool{
															DefaultBool(false),
														},
													},
													"value_data": rsschema.StringAttribute{
														Description: "",
														Optional:    true,
														Computed:    true,
														PlanModifiers: []planmodifier.String{
															DefaultString(""),
														},
														Validators: []validator.String{
															stringvalidator.LengthAtMost(1024),
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"data_loss_prevention": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"criteria": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"is_enabled": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.OneOf("no", "yes", "not-available"),
								},
							},
							"is_installed": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(true),
								},
							},
						},
					},
					"exclude_vendor": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"vendor": rsschema.ListNestedAttribute{
						Description: "",
						Optional:    true,
						NestedObject: rsschema.NestedAttributeObject{
							Attributes: map[string]rsschema.Attribute{
								"name": rsschema.StringAttribute{
									Description: "",
									Required:    true,
									PlanModifiers: []planmodifier.String{
										DefaultString(""),
									},
									Validators: []validator.String{
										stringvalidator.LengthAtMost(103),
									},
								},
								"product": rsschema.ListAttribute{
									Description: "",
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"description": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 255),
				},
			},
			"disk_backup": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"criteria": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"is_installed": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(true),
								},
							},
							"last_backup_time": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"not_available": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Validators: []validator.Bool{
											boolvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("not_within"),
												path.MatchRelative().AtParent().AtName("within"),
											),
										},
									},
									"not_within": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"days": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(1),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
											"hours": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(24),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
										},
									},
									"within": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"days": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(1),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
											"hours": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(24),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 65535),
												},
											},
										},
									},
								},
							},
						},
					},
					"exclude_vendor": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"vendor": rsschema.ListNestedAttribute{
						Description: "",
						Optional:    true,
						NestedObject: rsschema.NestedAttributeObject{
							Attributes: map[string]rsschema.Attribute{
								"name": rsschema.StringAttribute{
									Description: "",
									Required:    true,
									PlanModifiers: []planmodifier.String{
										DefaultString(""),
									},
									Validators: []validator.String{
										stringvalidator.LengthAtMost(103),
									},
								},
								"product": rsschema.ListAttribute{
									Description: "",
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"disk_encryption": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"criteria": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"encrypted_locations": rsschema.ListNestedAttribute{
								Description: "",
								Optional:    true,
								NestedObject: rsschema.NestedAttributeObject{
									Attributes: map[string]rsschema.Attribute{
										"encryption_state": rsschema.SingleNestedAttribute{
											Description: "",
											Optional:    true,
											Attributes: map[string]rsschema.Attribute{
												"is": rsschema.StringAttribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.String{
														DefaultString("encrypted"),
													},
													Validators: []validator.String{
														stringvalidator.OneOf("encrypted", "unencrypted", "partial", "unknown"),
													},
												},
												"is_not": rsschema.StringAttribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.String{
														DefaultString("encrypted"),
													},
													Validators: []validator.String{
														stringvalidator.OneOf("encrypted", "unencrypted", "partial", "unknown"),
													},
												},
											},
										},
										"name": rsschema.StringAttribute{
											Description: "",
											Required:    true,
											PlanModifiers: []planmodifier.String{
												DefaultString(""),
											},
											Validators: []validator.String{
												stringvalidator.LengthAtMost(1023),
											},
										},
									},
								},
							},
							"is_installed": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(true),
								},
							},
						},
					},
					"exclude_vendor": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"vendor": rsschema.ListNestedAttribute{
						Description: "",
						Optional:    true,
						NestedObject: rsschema.NestedAttributeObject{
							Attributes: map[string]rsschema.Attribute{
								"name": rsschema.StringAttribute{
									Description: "",
									Required:    true,
									PlanModifiers: []planmodifier.String{
										DefaultString(""),
									},
									Validators: []validator.String{
										stringvalidator.LengthAtMost(103),
									},
								},
								"product": rsschema.ListAttribute{
									Description: "",
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"firewall": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"criteria": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"is_enabled": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.OneOf("no", "yes", "not-available"),
								},
							},
							"is_installed": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(true),
								},
							},
						},
					},
					"exclude_vendor": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"vendor": rsschema.ListNestedAttribute{
						Description: "",
						Optional:    true,
						NestedObject: rsschema.NestedAttributeObject{
							Attributes: map[string]rsschema.Attribute{
								"name": rsschema.StringAttribute{
									Description: "",
									Required:    true,
									PlanModifiers: []planmodifier.String{
										DefaultString(""),
									},
									Validators: []validator.String{
										stringvalidator.LengthAtMost(103),
									},
								},
								"product": rsschema.ListAttribute{
									Description: "",
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"host_info": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"criteria": rsschema.SingleNestedAttribute{
						Description: "",
						Required:    true,
						Attributes: map[string]rsschema.Attribute{
							"client_version": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"contains": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is_not": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is"),
											),
										},
									},
								},
							},
							"domain": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"contains": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is_not": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is"),
											),
										},
									},
								},
							},
							"host_id": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"contains": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is_not": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is"),
											),
										},
									},
								},
							},
							"host_name": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"contains": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is_not": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is"),
											),
										},
									},
								},
							},
							"managed": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"os": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"contains": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"apple": rsschema.StringAttribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("All"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(0, 255),
													stringvalidator.ConflictsWith(
														path.MatchRelative().AtParent().AtName("Google"),
														path.MatchRelative().AtParent().AtName("Linux"),
														path.MatchRelative().AtParent().AtName("Microsoft"),
														path.MatchRelative().AtParent().AtName("Other"),
													),
												},
											},
											"google": rsschema.StringAttribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("All"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(0, 255),
													stringvalidator.ConflictsWith(
														path.MatchRelative().AtParent().AtName("Apple"),
														path.MatchRelative().AtParent().AtName("Linux"),
														path.MatchRelative().AtParent().AtName("Microsoft"),
														path.MatchRelative().AtParent().AtName("Other"),
													),
												},
											},
											"linux": rsschema.StringAttribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("All"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(0, 255),
													stringvalidator.ConflictsWith(
														path.MatchRelative().AtParent().AtName("Apple"),
														path.MatchRelative().AtParent().AtName("Google"),
														path.MatchRelative().AtParent().AtName("Microsoft"),
														path.MatchRelative().AtParent().AtName("Other"),
													),
												},
											},
											"microsoft": rsschema.StringAttribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString("All"),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(0, 255),
													stringvalidator.ConflictsWith(
														path.MatchRelative().AtParent().AtName("Apple"),
														path.MatchRelative().AtParent().AtName("Google"),
														path.MatchRelative().AtParent().AtName("Linux"),
														path.MatchRelative().AtParent().AtName("Other"),
													),
												},
											},
											"other": rsschema.StringAttribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.String{
													DefaultString(""),
												},
												Validators: []validator.String{
													stringvalidator.LengthBetween(0, 255),
													stringvalidator.ConflictsWith(
														path.MatchRelative().AtParent().AtName("Apple"),
														path.MatchRelative().AtParent().AtName("Google"),
														path.MatchRelative().AtParent().AtName("Linux"),
														path.MatchRelative().AtParent().AtName("Microsoft"),
													),
												},
											},
										},
									},
								},
							},
							"serial_number": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"contains": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is_not": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is"),
											),
										},
									},
								},
							},
						},
					},
				},
			},
			"object_id": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mobile_device": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"criteria": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"applications": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"has_malware": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"no": rsschema.BoolAttribute{
												Description: "",
												Optional:    true,
											},
											"yes": rsschema.SingleNestedAttribute{
												Description: "",
												Optional:    true,
												Attributes: map[string]rsschema.Attribute{
													"excludes": rsschema.ListNestedAttribute{
														Description: "",
														Optional:    true,
														NestedObject: rsschema.NestedAttributeObject{
															Attributes: map[string]rsschema.Attribute{
																"hash": rsschema.StringAttribute{
																	Description: "",
																	Optional:    true,
																	Computed:    true,
																	PlanModifiers: []planmodifier.String{
																		DefaultString(""),
																	},
																	Validators: []validator.String{
																		stringvalidator.LengthAtMost(1024),
																	},
																},
																"name": rsschema.StringAttribute{
																	Description: "",
																	Required:    true,
																	PlanModifiers: []planmodifier.String{
																		DefaultString(""),
																	},
																	Validators: []validator.String{
																		stringvalidator.LengthAtMost(31),
																	},
																},
																"package": rsschema.StringAttribute{
																	Description: "",
																	Optional:    true,
																	Computed:    true,
																	PlanModifiers: []planmodifier.String{
																		DefaultString(""),
																	},
																	Validators: []validator.String{
																		stringvalidator.LengthAtMost(1024),
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"has_unmanaged_app": rsschema.BoolAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.Bool{
											DefaultBool(false),
										},
									},
									"includes": rsschema.ListNestedAttribute{
										Description: "",
										Optional:    true,
										NestedObject: rsschema.NestedAttributeObject{
											Attributes: map[string]rsschema.Attribute{
												"hash": rsschema.StringAttribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.String{
														DefaultString(""),
													},
													Validators: []validator.String{
														stringvalidator.LengthAtMost(1024),
													},
												},
												"name": rsschema.StringAttribute{
													Description: "",
													Required:    true,
													PlanModifiers: []planmodifier.String{
														DefaultString(""),
													},
													Validators: []validator.String{
														stringvalidator.LengthAtMost(31),
													},
												},
												"package": rsschema.StringAttribute{
													Description: "",
													Optional:    true,
													Computed:    true,
													PlanModifiers: []planmodifier.String{
														DefaultString(""),
													},
													Validators: []validator.String{
														stringvalidator.LengthAtMost(1024),
													},
												},
											},
										},
									},
								},
							},
							"disk_encrypted": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"imei": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"contains": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is_not": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is"),
											),
										},
									},
								},
							},
							"jailbroken": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"last_checkin_time": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"not_within": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"days": rsschema.Int64Attribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(30),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 365),
												},
											},
										},
									},
									"within": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"days": rsschema.Int64Attribute{
												Description: "",
												Required:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(30),
												},
												Validators: []validator.Int64{
													int64validator.Between(1, 365),
												},
											},
										},
									},
								},
							},
							"model": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"contains": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is_not": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is"),
											),
										},
									},
								},
							},
							"passcode_set": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"phone_number": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"contains": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is_not": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is"),
											),
										},
									},
								},
							},
							"tag": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"contains": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("is"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is_not"),
											),
										},
									},
									"is_not": rsschema.StringAttribute{
										Description: "",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString(""),
										},
										Validators: []validator.String{
											stringvalidator.LengthBetween(0, 255),
											stringvalidator.ConflictsWith(
												path.MatchRelative().AtParent().AtName("contains"),
												path.MatchRelative().AtParent().AtName("is"),
											),
										},
									},
								},
							},
						},
					},
				},
			},
			"name": rsschema.StringAttribute{
				Description: "",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(31),
				},
			},
			"network_info": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"criteria": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"network": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"is": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"mobile": rsschema.SingleNestedAttribute{
												Description: "",
												Optional:    true,
												Attributes: map[string]rsschema.Attribute{
													"carrier": rsschema.StringAttribute{
														Description: "",
														Optional:    true,
														Computed:    true,
														PlanModifiers: []planmodifier.String{
															DefaultString(""),
														},
														Validators: []validator.String{
															stringvalidator.LengthAtMost(1023),
														},
													},
												},
											},
											"unknown": rsschema.BoolAttribute{
												Description: "",
												Optional:    true,
											},
											"wifi": rsschema.SingleNestedAttribute{
												Description: "",
												Optional:    true,
												Attributes: map[string]rsschema.Attribute{
													"ssid": rsschema.StringAttribute{
														Description: "",
														Optional:    true,
														Computed:    true,
														PlanModifiers: []planmodifier.String{
															DefaultString(""),
														},
														Validators: []validator.String{
															stringvalidator.LengthAtMost(1023),
														},
													},
												},
											},
										},
									},
									"is_not": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"ethernet": rsschema.BoolAttribute{
												Description: "",
												Optional:    true,
											},
											"mobile": rsschema.SingleNestedAttribute{
												Description: "",
												Optional:    true,
												Attributes: map[string]rsschema.Attribute{
													"carrier": rsschema.StringAttribute{
														Description: "",
														Optional:    true,
														Computed:    true,
														PlanModifiers: []planmodifier.String{
															DefaultString(""),
														},
														Validators: []validator.String{
															stringvalidator.LengthAtMost(1023),
														},
													},
												},
											},
											"unknown": rsschema.BoolAttribute{
												Description: "",
												Optional:    true,
											},
											"wifi": rsschema.SingleNestedAttribute{
												Description: "",
												Optional:    true,
												Attributes: map[string]rsschema.Attribute{
													"ssid": rsschema.StringAttribute{
														Description: "",
														Optional:    true,
														Computed:    true,
														PlanModifiers: []planmodifier.String{
															DefaultString(""),
														},
														Validators: []validator.String{
															stringvalidator.LengthAtMost(1023),
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"patch_management": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"criteria": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"is_enabled": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.OneOf("no", "yes", "not-available"),
								},
							},
							"is_installed": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(true),
								},
							},
							"missing_patches": rsschema.SingleNestedAttribute{
								Description: "",
								Optional:    true,
								Attributes: map[string]rsschema.Attribute{
									"check": rsschema.StringAttribute{
										Description: "",
										Required:    true,
										PlanModifiers: []planmodifier.String{
											DefaultString("has-any"),
										},
										Validators: []validator.String{
											stringvalidator.OneOf("has-any", "has-none", "has-all"),
										},
									},
									"patches": rsschema.ListAttribute{
										Description: "",
										Optional:    true,
										ElementType: types.StringType,
									},
									"severity": rsschema.SingleNestedAttribute{
										Description: "",
										Optional:    true,
										Attributes: map[string]rsschema.Attribute{
											"greater_equal": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(0),
												},
												Validators: []validator.Int64{
													int64validator.Between(0, 100000),
													int64validator.ConflictsWith(
														path.MatchRelative().AtParent().AtName("greater_than"),
														path.MatchRelative().AtParent().AtName("is"),
														path.MatchRelative().AtParent().AtName("is_not"),
														path.MatchRelative().AtParent().AtName("less_equal"),
														path.MatchRelative().AtParent().AtName("less_than"),
													),
												},
											},
											"greater_than": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(0),
												},
												Validators: []validator.Int64{
													int64validator.Between(0, 100000),
													int64validator.ConflictsWith(
														path.MatchRelative().AtParent().AtName("greater_equal"),
														path.MatchRelative().AtParent().AtName("is"),
														path.MatchRelative().AtParent().AtName("is_not"),
														path.MatchRelative().AtParent().AtName("less_equal"),
														path.MatchRelative().AtParent().AtName("less_than"),
													),
												},
											},
											"is": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(0),
												},
												Validators: []validator.Int64{
													int64validator.Between(0, 100000),
													int64validator.ConflictsWith(
														path.MatchRelative().AtParent().AtName("greater_equal"),
														path.MatchRelative().AtParent().AtName("greater_than"),
														path.MatchRelative().AtParent().AtName("is_not"),
														path.MatchRelative().AtParent().AtName("less_equal"),
														path.MatchRelative().AtParent().AtName("less_than"),
													),
												},
											},
											"is_not": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(0),
												},
												Validators: []validator.Int64{
													int64validator.Between(0, 100000),
													int64validator.ConflictsWith(
														path.MatchRelative().AtParent().AtName("greater_equal"),
														path.MatchRelative().AtParent().AtName("greater_than"),
														path.MatchRelative().AtParent().AtName("is"),
														path.MatchRelative().AtParent().AtName("less_equal"),
														path.MatchRelative().AtParent().AtName("less_than"),
													),
												},
											},
											"less_equal": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(0),
												},
												Validators: []validator.Int64{
													int64validator.Between(0, 100000),
													int64validator.ConflictsWith(
														path.MatchRelative().AtParent().AtName("greater_equal"),
														path.MatchRelative().AtParent().AtName("greater_than"),
														path.MatchRelative().AtParent().AtName("is"),
														path.MatchRelative().AtParent().AtName("is_not"),
														path.MatchRelative().AtParent().AtName("less_than"),
													),
												},
											},
											"less_than": rsschema.Int64Attribute{
												Description: "",
												Optional:    true,
												Computed:    true,
												PlanModifiers: []planmodifier.Int64{
													DefaultInt64(0),
												},
												Validators: []validator.Int64{
													int64validator.Between(0, 100000),
													int64validator.ConflictsWith(
														path.MatchRelative().AtParent().AtName("greater_equal"),
														path.MatchRelative().AtParent().AtName("greater_than"),
														path.MatchRelative().AtParent().AtName("is"),
														path.MatchRelative().AtParent().AtName("is_not"),
														path.MatchRelative().AtParent().AtName("less_equal"),
													),
												},
											},
										},
									},
								},
							},
						},
					},
					"exclude_vendor": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
						},
					},
					"vendor": rsschema.ListNestedAttribute{
						Description: "",
						Optional:    true,
						NestedObject: rsschema.NestedAttributeObject{
							Attributes: map[string]rsschema.Attribute{
								"name": rsschema.StringAttribute{
									Description: "",
									Required:    true,
									PlanModifiers: []planmodifier.String{
										DefaultString(""),
									},
									Validators: []validator.String{
										stringvalidator.LengthAtMost(103),
									},
								},
								"product": rsschema.ListAttribute{
									Description: "",
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *objectsHipObjectsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *objectsHipObjectsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state objectsHipObjectsRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_objects_hip_objects",
		"folder":        state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := yCYVNEN.NewClient(r.client)
	input := yCYVNEN.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 dJpBrWV.Config
	var var1 *dJpBrWV.AntiMalwareObject
	if state.AntiMalware != nil {
		var1 = &dJpBrWV.AntiMalwareObject{}
		var var2 *dJpBrWV.CriteriaObject
		if state.AntiMalware.Criteria != nil {
			var2 = &dJpBrWV.CriteriaObject{}
			var2.IsInstalled = state.AntiMalware.Criteria.IsInstalled.ValueBool()
			var var3 *dJpBrWV.LastScanTimeObject
			if state.AntiMalware.Criteria.LastScanTime != nil {
				var3 = &dJpBrWV.LastScanTimeObject{}
				if state.AntiMalware.Criteria.LastScanTime.NotAvailable.ValueBool() {
					var3.NotAvailable = struct{}{}
				}
				var var4 *dJpBrWV.NotWithinObject
				if state.AntiMalware.Criteria.LastScanTime.NotWithin != nil {
					var4 = &dJpBrWV.NotWithinObject{}
					var4.Days = state.AntiMalware.Criteria.LastScanTime.NotWithin.Days.ValueInt64()
					var4.Hours = state.AntiMalware.Criteria.LastScanTime.NotWithin.Hours.ValueInt64()
				}
				var3.NotWithin = var4
				var var5 *dJpBrWV.WithinObject
				if state.AntiMalware.Criteria.LastScanTime.Within != nil {
					var5 = &dJpBrWV.WithinObject{}
					var5.Days = state.AntiMalware.Criteria.LastScanTime.Within.Days.ValueInt64()
					var5.Hours = state.AntiMalware.Criteria.LastScanTime.Within.Hours.ValueInt64()
				}
				var3.Within = var5
			}
			var2.LastScanTime = var3
			var var6 *dJpBrWV.ProductVersionObject
			if state.AntiMalware.Criteria.ProductVersion != nil {
				var6 = &dJpBrWV.ProductVersionObject{}
				var6.Contains = state.AntiMalware.Criteria.ProductVersion.Contains.ValueString()
				var6.GreaterEqual = state.AntiMalware.Criteria.ProductVersion.GreaterEqual.ValueString()
				var6.GreaterThan = state.AntiMalware.Criteria.ProductVersion.GreaterThan.ValueString()
				var6.Is = state.AntiMalware.Criteria.ProductVersion.Is.ValueString()
				var6.IsNot = state.AntiMalware.Criteria.ProductVersion.IsNot.ValueString()
				var6.LessEqual = state.AntiMalware.Criteria.ProductVersion.LessEqual.ValueString()
				var6.LessThan = state.AntiMalware.Criteria.ProductVersion.LessThan.ValueString()
				var var7 *dJpBrWV.NotWithinObject1
				if state.AntiMalware.Criteria.ProductVersion.NotWithin != nil {
					var7 = &dJpBrWV.NotWithinObject1{}
					var7.Versions = state.AntiMalware.Criteria.ProductVersion.NotWithin.Versions.ValueInt64()
				}
				var6.NotWithin = var7
				var var8 *dJpBrWV.WithinObject1
				if state.AntiMalware.Criteria.ProductVersion.Within != nil {
					var8 = &dJpBrWV.WithinObject1{}
					var8.Versions = state.AntiMalware.Criteria.ProductVersion.Within.Versions.ValueInt64()
				}
				var6.Within = var8
			}
			var2.ProductVersion = var6
			var2.RealTimeProtection = state.AntiMalware.Criteria.RealTimeProtection.ValueString()
			var var9 *dJpBrWV.VirdefVersionObject
			if state.AntiMalware.Criteria.VirdefVersion != nil {
				var9 = &dJpBrWV.VirdefVersionObject{}
				var var10 *dJpBrWV.NotWithinObject2
				if state.AntiMalware.Criteria.VirdefVersion.NotWithin != nil {
					var10 = &dJpBrWV.NotWithinObject2{}
					var10.Days = state.AntiMalware.Criteria.VirdefVersion.NotWithin.Days.ValueInt64()
					var10.Versions = state.AntiMalware.Criteria.VirdefVersion.NotWithin.Versions.ValueInt64()
				}
				var9.NotWithin = var10
				var var11 *dJpBrWV.WithinObject2
				if state.AntiMalware.Criteria.VirdefVersion.Within != nil {
					var11 = &dJpBrWV.WithinObject2{}
					var11.Days = state.AntiMalware.Criteria.VirdefVersion.Within.Days.ValueInt64()
					var11.Versions = state.AntiMalware.Criteria.VirdefVersion.Within.Versions.ValueInt64()
				}
				var9.Within = var11
			}
			var2.VirdefVersion = var9
		}
		var1.Criteria = var2
		var1.ExcludeVendor = state.AntiMalware.ExcludeVendor.ValueBool()
		var var12 []dJpBrWV.VendorObject
		if len(state.AntiMalware.Vendor) != 0 {
			var12 = make([]dJpBrWV.VendorObject, 0, len(state.AntiMalware.Vendor))
			for var13Index := range state.AntiMalware.Vendor {
				var13 := state.AntiMalware.Vendor[var13Index]
				var var14 dJpBrWV.VendorObject
				var14.Name = var13.Name.ValueString()
				var14.Product = DecodeStringSlice(var13.Product)
				var12 = append(var12, var14)
			}
		}
		var1.Vendor = var12
	}
	var0.AntiMalware = var1
	var var15 *dJpBrWV.CertificateObject
	if state.Certificate != nil {
		var15 = &dJpBrWV.CertificateObject{}
		var var16 *dJpBrWV.CriteriaObject1
		if state.Certificate.Criteria != nil {
			var16 = &dJpBrWV.CriteriaObject1{}
			var var17 []dJpBrWV.CertificateAttributesObject
			if len(state.Certificate.Criteria.CertificateAttributes) != 0 {
				var17 = make([]dJpBrWV.CertificateAttributesObject, 0, len(state.Certificate.Criteria.CertificateAttributes))
				for var18Index := range state.Certificate.Criteria.CertificateAttributes {
					var18 := state.Certificate.Criteria.CertificateAttributes[var18Index]
					var var19 dJpBrWV.CertificateAttributesObject
					var19.Name = var18.Name.ValueString()
					var19.Value = var18.Value.ValueString()
					var17 = append(var17, var19)
				}
			}
			var16.CertificateAttributes = var17
			var16.CertificateProfile = state.Certificate.Criteria.CertificateProfile.ValueString()
		}
		var15.Criteria = var16
	}
	var0.Certificate = var15
	var var20 *dJpBrWV.CustomChecksObject
	if state.CustomChecks != nil {
		var20 = &dJpBrWV.CustomChecksObject{}
		var var21 dJpBrWV.CriteriaObject2
		var var22 []dJpBrWV.PlistObject
		if len(state.CustomChecks.Criteria.Plist) != 0 {
			var22 = make([]dJpBrWV.PlistObject, 0, len(state.CustomChecks.Criteria.Plist))
			for var23Index := range state.CustomChecks.Criteria.Plist {
				var23 := state.CustomChecks.Criteria.Plist[var23Index]
				var var24 dJpBrWV.PlistObject
				var var25 []dJpBrWV.KeyObject
				if len(var23.Key) != 0 {
					var25 = make([]dJpBrWV.KeyObject, 0, len(var23.Key))
					for var26Index := range var23.Key {
						var26 := var23.Key[var26Index]
						var var27 dJpBrWV.KeyObject
						var27.Name = var26.Name.ValueString()
						var27.Negate = var26.Negate.ValueBool()
						var27.Value = var26.Value.ValueString()
						var25 = append(var25, var27)
					}
				}
				var24.Key = var25
				var24.Name = var23.Name.ValueString()
				var24.Negate = var23.Negate.ValueBool()
				var22 = append(var22, var24)
			}
		}
		var21.Plist = var22
		var var28 []dJpBrWV.ProcessListObject
		if len(state.CustomChecks.Criteria.ProcessList) != 0 {
			var28 = make([]dJpBrWV.ProcessListObject, 0, len(state.CustomChecks.Criteria.ProcessList))
			for var29Index := range state.CustomChecks.Criteria.ProcessList {
				var29 := state.CustomChecks.Criteria.ProcessList[var29Index]
				var var30 dJpBrWV.ProcessListObject
				var30.Name = var29.Name.ValueString()
				var30.Running = var29.Running.ValueBool()
				var28 = append(var28, var30)
			}
		}
		var21.ProcessList = var28
		var var31 []dJpBrWV.RegistryKeyObject
		if len(state.CustomChecks.Criteria.RegistryKey) != 0 {
			var31 = make([]dJpBrWV.RegistryKeyObject, 0, len(state.CustomChecks.Criteria.RegistryKey))
			for var32Index := range state.CustomChecks.Criteria.RegistryKey {
				var32 := state.CustomChecks.Criteria.RegistryKey[var32Index]
				var var33 dJpBrWV.RegistryKeyObject
				var33.DefaultValueData = var32.DefaultValueData.ValueString()
				var33.Name = var32.Name.ValueString()
				var33.Negate = var32.Negate.ValueBool()
				var var34 []dJpBrWV.RegistryValueObject
				if len(var32.RegistryValue) != 0 {
					var34 = make([]dJpBrWV.RegistryValueObject, 0, len(var32.RegistryValue))
					for var35Index := range var32.RegistryValue {
						var35 := var32.RegistryValue[var35Index]
						var var36 dJpBrWV.RegistryValueObject
						var36.Name = var35.Name.ValueString()
						var36.Negate = var35.Negate.ValueBool()
						var36.ValueData = var35.ValueData.ValueString()
						var34 = append(var34, var36)
					}
				}
				var33.RegistryValue = var34
				var31 = append(var31, var33)
			}
		}
		var21.RegistryKey = var31
		var20.Criteria = var21
	}
	var0.CustomChecks = var20
	var var37 *dJpBrWV.DataLossPreventionObject
	if state.DataLossPrevention != nil {
		var37 = &dJpBrWV.DataLossPreventionObject{}
		var var38 *dJpBrWV.CriteriaObject3
		if state.DataLossPrevention.Criteria != nil {
			var38 = &dJpBrWV.CriteriaObject3{}
			var38.IsEnabled = state.DataLossPrevention.Criteria.IsEnabled.ValueString()
			var38.IsInstalled = state.DataLossPrevention.Criteria.IsInstalled.ValueBool()
		}
		var37.Criteria = var38
		var37.ExcludeVendor = state.DataLossPrevention.ExcludeVendor.ValueBool()
		var var39 []dJpBrWV.VendorObject1
		if len(state.DataLossPrevention.Vendor) != 0 {
			var39 = make([]dJpBrWV.VendorObject1, 0, len(state.DataLossPrevention.Vendor))
			for var40Index := range state.DataLossPrevention.Vendor {
				var40 := state.DataLossPrevention.Vendor[var40Index]
				var var41 dJpBrWV.VendorObject1
				var41.Name = var40.Name.ValueString()
				var41.Product = DecodeStringSlice(var40.Product)
				var39 = append(var39, var41)
			}
		}
		var37.Vendor = var39
	}
	var0.DataLossPrevention = var37
	var0.Description = state.Description.ValueString()
	var var42 *dJpBrWV.DiskBackupObject
	if state.DiskBackup != nil {
		var42 = &dJpBrWV.DiskBackupObject{}
		var var43 *dJpBrWV.CriteriaObject4
		if state.DiskBackup.Criteria != nil {
			var43 = &dJpBrWV.CriteriaObject4{}
			var43.IsInstalled = state.DiskBackup.Criteria.IsInstalled.ValueBool()
			var var44 *dJpBrWV.LastBackupTimeObject
			if state.DiskBackup.Criteria.LastBackupTime != nil {
				var44 = &dJpBrWV.LastBackupTimeObject{}
				if state.DiskBackup.Criteria.LastBackupTime.NotAvailable.ValueBool() {
					var44.NotAvailable = struct{}{}
				}
				var var45 *dJpBrWV.NotWithinObject
				if state.DiskBackup.Criteria.LastBackupTime.NotWithin != nil {
					var45 = &dJpBrWV.NotWithinObject{}
					var45.Days = state.DiskBackup.Criteria.LastBackupTime.NotWithin.Days.ValueInt64()
					var45.Hours = state.DiskBackup.Criteria.LastBackupTime.NotWithin.Hours.ValueInt64()
				}
				var44.NotWithin = var45
				var var46 *dJpBrWV.WithinObject
				if state.DiskBackup.Criteria.LastBackupTime.Within != nil {
					var46 = &dJpBrWV.WithinObject{}
					var46.Days = state.DiskBackup.Criteria.LastBackupTime.Within.Days.ValueInt64()
					var46.Hours = state.DiskBackup.Criteria.LastBackupTime.Within.Hours.ValueInt64()
				}
				var44.Within = var46
			}
			var43.LastBackupTime = var44
		}
		var42.Criteria = var43
		var42.ExcludeVendor = state.DiskBackup.ExcludeVendor.ValueBool()
		var var47 []dJpBrWV.VendorObject
		if len(state.DiskBackup.Vendor) != 0 {
			var47 = make([]dJpBrWV.VendorObject, 0, len(state.DiskBackup.Vendor))
			for var48Index := range state.DiskBackup.Vendor {
				var48 := state.DiskBackup.Vendor[var48Index]
				var var49 dJpBrWV.VendorObject
				var49.Name = var48.Name.ValueString()
				var49.Product = DecodeStringSlice(var48.Product)
				var47 = append(var47, var49)
			}
		}
		var42.Vendor = var47
	}
	var0.DiskBackup = var42
	var var50 *dJpBrWV.DiskEncryptionObject
	if state.DiskEncryption != nil {
		var50 = &dJpBrWV.DiskEncryptionObject{}
		var var51 *dJpBrWV.CriteriaObject5
		if state.DiskEncryption.Criteria != nil {
			var51 = &dJpBrWV.CriteriaObject5{}
			var var52 []dJpBrWV.EncryptedLocationsObject
			if len(state.DiskEncryption.Criteria.EncryptedLocations) != 0 {
				var52 = make([]dJpBrWV.EncryptedLocationsObject, 0, len(state.DiskEncryption.Criteria.EncryptedLocations))
				for var53Index := range state.DiskEncryption.Criteria.EncryptedLocations {
					var53 := state.DiskEncryption.Criteria.EncryptedLocations[var53Index]
					var var54 dJpBrWV.EncryptedLocationsObject
					var var55 *dJpBrWV.EncryptionStateObject
					if var53.EncryptionState != nil {
						var55 = &dJpBrWV.EncryptionStateObject{}
						var55.Is = var53.EncryptionState.Is.ValueString()
						var55.IsNot = var53.EncryptionState.IsNot.ValueString()
					}
					var54.EncryptionState = var55
					var54.Name = var53.Name.ValueString()
					var52 = append(var52, var54)
				}
			}
			var51.EncryptedLocations = var52
			var51.IsInstalled = state.DiskEncryption.Criteria.IsInstalled.ValueBool()
		}
		var50.Criteria = var51
		var50.ExcludeVendor = state.DiskEncryption.ExcludeVendor.ValueBool()
		var var56 []dJpBrWV.VendorObject
		if len(state.DiskEncryption.Vendor) != 0 {
			var56 = make([]dJpBrWV.VendorObject, 0, len(state.DiskEncryption.Vendor))
			for var57Index := range state.DiskEncryption.Vendor {
				var57 := state.DiskEncryption.Vendor[var57Index]
				var var58 dJpBrWV.VendorObject
				var58.Name = var57.Name.ValueString()
				var58.Product = DecodeStringSlice(var57.Product)
				var56 = append(var56, var58)
			}
		}
		var50.Vendor = var56
	}
	var0.DiskEncryption = var50
	var var59 *dJpBrWV.FirewallObject
	if state.Firewall != nil {
		var59 = &dJpBrWV.FirewallObject{}
		var var60 *dJpBrWV.CriteriaObject3
		if state.Firewall.Criteria != nil {
			var60 = &dJpBrWV.CriteriaObject3{}
			var60.IsEnabled = state.Firewall.Criteria.IsEnabled.ValueString()
			var60.IsInstalled = state.Firewall.Criteria.IsInstalled.ValueBool()
		}
		var59.Criteria = var60
		var59.ExcludeVendor = state.Firewall.ExcludeVendor.ValueBool()
		var var61 []dJpBrWV.VendorObject
		if len(state.Firewall.Vendor) != 0 {
			var61 = make([]dJpBrWV.VendorObject, 0, len(state.Firewall.Vendor))
			for var62Index := range state.Firewall.Vendor {
				var62 := state.Firewall.Vendor[var62Index]
				var var63 dJpBrWV.VendorObject
				var63.Name = var62.Name.ValueString()
				var63.Product = DecodeStringSlice(var62.Product)
				var61 = append(var61, var63)
			}
		}
		var59.Vendor = var61
	}
	var0.Firewall = var59
	var var64 *dJpBrWV.HostInfoObject
	if state.HostInfo != nil {
		var64 = &dJpBrWV.HostInfoObject{}
		var var65 dJpBrWV.CriteriaObject6
		var var66 *dJpBrWV.ClientVersionObject
		if state.HostInfo.Criteria.ClientVersion != nil {
			var66 = &dJpBrWV.ClientVersionObject{}
			var66.Contains = state.HostInfo.Criteria.ClientVersion.Contains.ValueString()
			var66.Is = state.HostInfo.Criteria.ClientVersion.Is.ValueString()
			var66.IsNot = state.HostInfo.Criteria.ClientVersion.IsNot.ValueString()
		}
		var65.ClientVersion = var66
		var var67 *dJpBrWV.DomainObject
		if state.HostInfo.Criteria.Domain != nil {
			var67 = &dJpBrWV.DomainObject{}
			var67.Contains = state.HostInfo.Criteria.Domain.Contains.ValueString()
			var67.Is = state.HostInfo.Criteria.Domain.Is.ValueString()
			var67.IsNot = state.HostInfo.Criteria.Domain.IsNot.ValueString()
		}
		var65.Domain = var67
		var var68 *dJpBrWV.HostIdObject
		if state.HostInfo.Criteria.HostId != nil {
			var68 = &dJpBrWV.HostIdObject{}
			var68.Contains = state.HostInfo.Criteria.HostId.Contains.ValueString()
			var68.Is = state.HostInfo.Criteria.HostId.Is.ValueString()
			var68.IsNot = state.HostInfo.Criteria.HostId.IsNot.ValueString()
		}
		var65.HostId = var68
		var var69 *dJpBrWV.HostNameObject
		if state.HostInfo.Criteria.HostName != nil {
			var69 = &dJpBrWV.HostNameObject{}
			var69.Contains = state.HostInfo.Criteria.HostName.Contains.ValueString()
			var69.Is = state.HostInfo.Criteria.HostName.Is.ValueString()
			var69.IsNot = state.HostInfo.Criteria.HostName.IsNot.ValueString()
		}
		var65.HostName = var69
		var65.Managed = state.HostInfo.Criteria.Managed.ValueBool()
		var var70 *dJpBrWV.OsObject
		if state.HostInfo.Criteria.Os != nil {
			var70 = &dJpBrWV.OsObject{}
			var var71 *dJpBrWV.ContainsObject
			if state.HostInfo.Criteria.Os.Contains != nil {
				var71 = &dJpBrWV.ContainsObject{}
				var71.Apple = state.HostInfo.Criteria.Os.Contains.Apple.ValueString()
				var71.Google = state.HostInfo.Criteria.Os.Contains.Google.ValueString()
				var71.Linux = state.HostInfo.Criteria.Os.Contains.Linux.ValueString()
				var71.Microsoft = state.HostInfo.Criteria.Os.Contains.Microsoft.ValueString()
				var71.Other = state.HostInfo.Criteria.Os.Contains.Other.ValueString()
			}
			var70.Contains = var71
		}
		var65.Os = var70
		var var72 *dJpBrWV.SerialNumberObject
		if state.HostInfo.Criteria.SerialNumber != nil {
			var72 = &dJpBrWV.SerialNumberObject{}
			var72.Contains = state.HostInfo.Criteria.SerialNumber.Contains.ValueString()
			var72.Is = state.HostInfo.Criteria.SerialNumber.Is.ValueString()
			var72.IsNot = state.HostInfo.Criteria.SerialNumber.IsNot.ValueString()
		}
		var65.SerialNumber = var72
		var64.Criteria = var65
	}
	var0.HostInfo = var64
	var var73 *dJpBrWV.MobileDeviceObject
	if state.MobileDevice != nil {
		var73 = &dJpBrWV.MobileDeviceObject{}
		var var74 *dJpBrWV.CriteriaObject7
		if state.MobileDevice.Criteria != nil {
			var74 = &dJpBrWV.CriteriaObject7{}
			var var75 *dJpBrWV.ApplicationsObject
			if state.MobileDevice.Criteria.Applications != nil {
				var75 = &dJpBrWV.ApplicationsObject{}
				var var76 *dJpBrWV.HasMalwareObject
				if state.MobileDevice.Criteria.Applications.HasMalware != nil {
					var76 = &dJpBrWV.HasMalwareObject{}
					if state.MobileDevice.Criteria.Applications.HasMalware.No.ValueBool() {
						var76.No = struct{}{}
					}
					var var77 *dJpBrWV.YesObject
					if state.MobileDevice.Criteria.Applications.HasMalware.Yes != nil {
						var77 = &dJpBrWV.YesObject{}
						var var78 []dJpBrWV.ExcludesObject
						if len(state.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes) != 0 {
							var78 = make([]dJpBrWV.ExcludesObject, 0, len(state.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes))
							for var79Index := range state.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes {
								var79 := state.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes[var79Index]
								var var80 dJpBrWV.ExcludesObject
								var80.Hash = var79.Hash.ValueString()
								var80.Name = var79.Name.ValueString()
								var80.Package = var79.Package.ValueString()
								var78 = append(var78, var80)
							}
						}
						var77.Excludes = var78
					}
					var76.Yes = var77
				}
				var75.HasMalware = var76
				var75.HasUnmanagedApp = state.MobileDevice.Criteria.Applications.HasUnmanagedApp.ValueBool()
				var var81 []dJpBrWV.IncludesObject
				if len(state.MobileDevice.Criteria.Applications.Includes) != 0 {
					var81 = make([]dJpBrWV.IncludesObject, 0, len(state.MobileDevice.Criteria.Applications.Includes))
					for var82Index := range state.MobileDevice.Criteria.Applications.Includes {
						var82 := state.MobileDevice.Criteria.Applications.Includes[var82Index]
						var var83 dJpBrWV.IncludesObject
						var83.Hash = var82.Hash.ValueString()
						var83.Name = var82.Name.ValueString()
						var83.Package = var82.Package.ValueString()
						var81 = append(var81, var83)
					}
				}
				var75.Includes = var81
			}
			var74.Applications = var75
			var74.DiskEncrypted = state.MobileDevice.Criteria.DiskEncrypted.ValueBool()
			var var84 *dJpBrWV.ImeiObject
			if state.MobileDevice.Criteria.Imei != nil {
				var84 = &dJpBrWV.ImeiObject{}
				var84.Contains = state.MobileDevice.Criteria.Imei.Contains.ValueString()
				var84.Is = state.MobileDevice.Criteria.Imei.Is.ValueString()
				var84.IsNot = state.MobileDevice.Criteria.Imei.IsNot.ValueString()
			}
			var74.Imei = var84
			var74.Jailbroken = state.MobileDevice.Criteria.Jailbroken.ValueBool()
			var var85 *dJpBrWV.LastCheckinTimeObject
			if state.MobileDevice.Criteria.LastCheckinTime != nil {
				var85 = &dJpBrWV.LastCheckinTimeObject{}
				var var86 *dJpBrWV.NotWithinObject3
				if state.MobileDevice.Criteria.LastCheckinTime.NotWithin != nil {
					var86 = &dJpBrWV.NotWithinObject3{}
					var86.Days = state.MobileDevice.Criteria.LastCheckinTime.NotWithin.Days.ValueInt64()
				}
				var85.NotWithin = var86
				var var87 *dJpBrWV.WithinObject3
				if state.MobileDevice.Criteria.LastCheckinTime.Within != nil {
					var87 = &dJpBrWV.WithinObject3{}
					var87.Days = state.MobileDevice.Criteria.LastCheckinTime.Within.Days.ValueInt64()
				}
				var85.Within = var87
			}
			var74.LastCheckinTime = var85
			var var88 *dJpBrWV.ModelObject
			if state.MobileDevice.Criteria.Model != nil {
				var88 = &dJpBrWV.ModelObject{}
				var88.Contains = state.MobileDevice.Criteria.Model.Contains.ValueString()
				var88.Is = state.MobileDevice.Criteria.Model.Is.ValueString()
				var88.IsNot = state.MobileDevice.Criteria.Model.IsNot.ValueString()
			}
			var74.Model = var88
			var74.PasscodeSet = state.MobileDevice.Criteria.PasscodeSet.ValueBool()
			var var89 *dJpBrWV.PhoneNumberObject
			if state.MobileDevice.Criteria.PhoneNumber != nil {
				var89 = &dJpBrWV.PhoneNumberObject{}
				var89.Contains = state.MobileDevice.Criteria.PhoneNumber.Contains.ValueString()
				var89.Is = state.MobileDevice.Criteria.PhoneNumber.Is.ValueString()
				var89.IsNot = state.MobileDevice.Criteria.PhoneNumber.IsNot.ValueString()
			}
			var74.PhoneNumber = var89
			var var90 *dJpBrWV.TagObject
			if state.MobileDevice.Criteria.Tag != nil {
				var90 = &dJpBrWV.TagObject{}
				var90.Contains = state.MobileDevice.Criteria.Tag.Contains.ValueString()
				var90.Is = state.MobileDevice.Criteria.Tag.Is.ValueString()
				var90.IsNot = state.MobileDevice.Criteria.Tag.IsNot.ValueString()
			}
			var74.Tag = var90
		}
		var73.Criteria = var74
	}
	var0.MobileDevice = var73
	var0.Name = state.Name.ValueString()
	var var91 *dJpBrWV.NetworkInfoObject
	if state.NetworkInfo != nil {
		var91 = &dJpBrWV.NetworkInfoObject{}
		var var92 *dJpBrWV.CriteriaObject8
		if state.NetworkInfo.Criteria != nil {
			var92 = &dJpBrWV.CriteriaObject8{}
			var var93 *dJpBrWV.NetworkObject
			if state.NetworkInfo.Criteria.Network != nil {
				var93 = &dJpBrWV.NetworkObject{}
				var var94 *dJpBrWV.IsObject
				if state.NetworkInfo.Criteria.Network.Is != nil {
					var94 = &dJpBrWV.IsObject{}
					var var95 *dJpBrWV.MobileObject
					if state.NetworkInfo.Criteria.Network.Is.Mobile != nil {
						var95 = &dJpBrWV.MobileObject{}
						var95.Carrier = state.NetworkInfo.Criteria.Network.Is.Mobile.Carrier.ValueString()
					}
					var94.Mobile = var95
					if state.NetworkInfo.Criteria.Network.Is.Unknown.ValueBool() {
						var94.Unknown = struct{}{}
					}
					var var96 *dJpBrWV.WifiObject
					if state.NetworkInfo.Criteria.Network.Is.Wifi != nil {
						var96 = &dJpBrWV.WifiObject{}
						var96.Ssid = state.NetworkInfo.Criteria.Network.Is.Wifi.Ssid.ValueString()
					}
					var94.Wifi = var96
				}
				var93.Is = var94
				var var97 *dJpBrWV.IsNotObject
				if state.NetworkInfo.Criteria.Network.IsNot != nil {
					var97 = &dJpBrWV.IsNotObject{}
					if state.NetworkInfo.Criteria.Network.IsNot.Ethernet.ValueBool() {
						var97.Ethernet = struct{}{}
					}
					var var98 *dJpBrWV.MobileObject
					if state.NetworkInfo.Criteria.Network.IsNot.Mobile != nil {
						var98 = &dJpBrWV.MobileObject{}
						var98.Carrier = state.NetworkInfo.Criteria.Network.IsNot.Mobile.Carrier.ValueString()
					}
					var97.Mobile = var98
					if state.NetworkInfo.Criteria.Network.IsNot.Unknown.ValueBool() {
						var97.Unknown = struct{}{}
					}
					var var99 *dJpBrWV.WifiObject
					if state.NetworkInfo.Criteria.Network.IsNot.Wifi != nil {
						var99 = &dJpBrWV.WifiObject{}
						var99.Ssid = state.NetworkInfo.Criteria.Network.IsNot.Wifi.Ssid.ValueString()
					}
					var97.Wifi = var99
				}
				var93.IsNot = var97
			}
			var92.Network = var93
		}
		var91.Criteria = var92
	}
	var0.NetworkInfo = var91
	var var100 *dJpBrWV.PatchManagementObject
	if state.PatchManagement != nil {
		var100 = &dJpBrWV.PatchManagementObject{}
		var var101 *dJpBrWV.CriteriaObject9
		if state.PatchManagement.Criteria != nil {
			var101 = &dJpBrWV.CriteriaObject9{}
			var101.IsEnabled = state.PatchManagement.Criteria.IsEnabled.ValueString()
			var101.IsInstalled = state.PatchManagement.Criteria.IsInstalled.ValueBool()
			var var102 *dJpBrWV.MissingPatchesObject
			if state.PatchManagement.Criteria.MissingPatches != nil {
				var102 = &dJpBrWV.MissingPatchesObject{}
				var102.Check = state.PatchManagement.Criteria.MissingPatches.Check.ValueString()
				var102.Patches = DecodeStringSlice(state.PatchManagement.Criteria.MissingPatches.Patches)
				var var103 *dJpBrWV.SeverityObject
				if state.PatchManagement.Criteria.MissingPatches.Severity != nil {
					var103 = &dJpBrWV.SeverityObject{}
					var103.GreaterEqual = state.PatchManagement.Criteria.MissingPatches.Severity.GreaterEqual.ValueInt64()
					var103.GreaterThan = state.PatchManagement.Criteria.MissingPatches.Severity.GreaterThan.ValueInt64()
					var103.Is = state.PatchManagement.Criteria.MissingPatches.Severity.Is.ValueInt64()
					var103.IsNot = state.PatchManagement.Criteria.MissingPatches.Severity.IsNot.ValueInt64()
					var103.LessEqual = state.PatchManagement.Criteria.MissingPatches.Severity.LessEqual.ValueInt64()
					var103.LessThan = state.PatchManagement.Criteria.MissingPatches.Severity.LessThan.ValueInt64()
				}
				var102.Severity = var103
			}
			var101.MissingPatches = var102
		}
		var100.Criteria = var101
		var100.ExcludeVendor = state.PatchManagement.ExcludeVendor.ValueBool()
		var var104 []dJpBrWV.VendorObject1
		if len(state.PatchManagement.Vendor) != 0 {
			var104 = make([]dJpBrWV.VendorObject1, 0, len(state.PatchManagement.Vendor))
			for var105Index := range state.PatchManagement.Vendor {
				var105 := state.PatchManagement.Vendor[var105Index]
				var var106 dJpBrWV.VendorObject1
				var106.Name = var105.Name.ValueString()
				var106.Product = DecodeStringSlice(var105.Product)
				var104 = append(var104, var106)
			}
		}
		var100.Vendor = var104
	}
	var0.PatchManagement = var100
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Create(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Store the answer to state.
	state.Id = types.StringValue(strings.Join([]string{input.Folder, ans.ObjectId}, IdSeparator))
	var var107 *objectsHipObjectsRsModelAntiMalwareObject
	if ans.AntiMalware != nil {
		var107 = &objectsHipObjectsRsModelAntiMalwareObject{}
		var var108 *objectsHipObjectsRsModelCriteriaObject
		if ans.AntiMalware.Criteria != nil {
			var108 = &objectsHipObjectsRsModelCriteriaObject{}
			var var109 *objectsHipObjectsRsModelLastScanTimeObject
			if ans.AntiMalware.Criteria.LastScanTime != nil {
				var109 = &objectsHipObjectsRsModelLastScanTimeObject{}
				var var110 *objectsHipObjectsRsModelNotWithinObject
				if ans.AntiMalware.Criteria.LastScanTime.NotWithin != nil {
					var110 = &objectsHipObjectsRsModelNotWithinObject{}
					var110.Days = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.NotWithin.Days)
					var110.Hours = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.NotWithin.Hours)
				}
				var var111 *objectsHipObjectsRsModelWithinObject
				if ans.AntiMalware.Criteria.LastScanTime.Within != nil {
					var111 = &objectsHipObjectsRsModelWithinObject{}
					var111.Days = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.Within.Days)
					var111.Hours = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.Within.Hours)
				}
				if ans.AntiMalware.Criteria.LastScanTime.NotAvailable != nil {
					var109.NotAvailable = types.BoolValue(true)
				}
				var109.NotWithin = var110
				var109.Within = var111
			}
			var var112 *objectsHipObjectsRsModelProductVersionObject
			if ans.AntiMalware.Criteria.ProductVersion != nil {
				var112 = &objectsHipObjectsRsModelProductVersionObject{}
				var var113 *objectsHipObjectsRsModelNotWithinObject1
				if ans.AntiMalware.Criteria.ProductVersion.NotWithin != nil {
					var113 = &objectsHipObjectsRsModelNotWithinObject1{}
					var113.Versions = types.Int64Value(ans.AntiMalware.Criteria.ProductVersion.NotWithin.Versions)
				}
				var var114 *objectsHipObjectsRsModelWithinObject1
				if ans.AntiMalware.Criteria.ProductVersion.Within != nil {
					var114 = &objectsHipObjectsRsModelWithinObject1{}
					var114.Versions = types.Int64Value(ans.AntiMalware.Criteria.ProductVersion.Within.Versions)
				}
				var112.Contains = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.Contains)
				var112.GreaterEqual = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.GreaterEqual)
				var112.GreaterThan = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.GreaterThan)
				var112.Is = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.Is)
				var112.IsNot = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.IsNot)
				var112.LessEqual = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.LessEqual)
				var112.LessThan = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.LessThan)
				var112.NotWithin = var113
				var112.Within = var114
			}
			var var115 *objectsHipObjectsRsModelVirdefVersionObject
			if ans.AntiMalware.Criteria.VirdefVersion != nil {
				var115 = &objectsHipObjectsRsModelVirdefVersionObject{}
				var var116 *objectsHipObjectsRsModelNotWithinObject2
				if ans.AntiMalware.Criteria.VirdefVersion.NotWithin != nil {
					var116 = &objectsHipObjectsRsModelNotWithinObject2{}
					var116.Days = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.NotWithin.Days)
					var116.Versions = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.NotWithin.Versions)
				}
				var var117 *objectsHipObjectsRsModelWithinObject2
				if ans.AntiMalware.Criteria.VirdefVersion.Within != nil {
					var117 = &objectsHipObjectsRsModelWithinObject2{}
					var117.Days = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.Within.Days)
					var117.Versions = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.Within.Versions)
				}
				var115.NotWithin = var116
				var115.Within = var117
			}
			var108.IsInstalled = types.BoolValue(ans.AntiMalware.Criteria.IsInstalled)
			var108.LastScanTime = var109
			var108.ProductVersion = var112
			var108.RealTimeProtection = types.StringValue(ans.AntiMalware.Criteria.RealTimeProtection)
			var108.VirdefVersion = var115
		}
		var var118 []objectsHipObjectsRsModelVendorObject
		if len(ans.AntiMalware.Vendor) != 0 {
			var118 = make([]objectsHipObjectsRsModelVendorObject, 0, len(ans.AntiMalware.Vendor))
			for var119Index := range ans.AntiMalware.Vendor {
				var119 := ans.AntiMalware.Vendor[var119Index]
				var var120 objectsHipObjectsRsModelVendorObject
				var120.Name = types.StringValue(var119.Name)
				var120.Product = EncodeStringSlice(var119.Product)
				var118 = append(var118, var120)
			}
		}
		var107.Criteria = var108
		var107.ExcludeVendor = types.BoolValue(ans.AntiMalware.ExcludeVendor)
		var107.Vendor = var118
	}
	var var121 *objectsHipObjectsRsModelCertificateObject
	if ans.Certificate != nil {
		var121 = &objectsHipObjectsRsModelCertificateObject{}
		var var122 *objectsHipObjectsRsModelCriteriaObject1
		if ans.Certificate.Criteria != nil {
			var122 = &objectsHipObjectsRsModelCriteriaObject1{}
			var var123 []objectsHipObjectsRsModelCertificateAttributesObject
			if len(ans.Certificate.Criteria.CertificateAttributes) != 0 {
				var123 = make([]objectsHipObjectsRsModelCertificateAttributesObject, 0, len(ans.Certificate.Criteria.CertificateAttributes))
				for var124Index := range ans.Certificate.Criteria.CertificateAttributes {
					var124 := ans.Certificate.Criteria.CertificateAttributes[var124Index]
					var var125 objectsHipObjectsRsModelCertificateAttributesObject
					var125.Name = types.StringValue(var124.Name)
					var125.Value = types.StringValue(var124.Value)
					var123 = append(var123, var125)
				}
			}
			var122.CertificateAttributes = var123
			var122.CertificateProfile = types.StringValue(ans.Certificate.Criteria.CertificateProfile)
		}
		var121.Criteria = var122
	}
	var var126 *objectsHipObjectsRsModelCustomChecksObject
	if ans.CustomChecks != nil {
		var126 = &objectsHipObjectsRsModelCustomChecksObject{}
		var var127 objectsHipObjectsRsModelCriteriaObject2
		var var128 []objectsHipObjectsRsModelPlistObject
		if len(ans.CustomChecks.Criteria.Plist) != 0 {
			var128 = make([]objectsHipObjectsRsModelPlistObject, 0, len(ans.CustomChecks.Criteria.Plist))
			for var129Index := range ans.CustomChecks.Criteria.Plist {
				var129 := ans.CustomChecks.Criteria.Plist[var129Index]
				var var130 objectsHipObjectsRsModelPlistObject
				var var131 []objectsHipObjectsRsModelKeyObject
				if len(var129.Key) != 0 {
					var131 = make([]objectsHipObjectsRsModelKeyObject, 0, len(var129.Key))
					for var132Index := range var129.Key {
						var132 := var129.Key[var132Index]
						var var133 objectsHipObjectsRsModelKeyObject
						var133.Name = types.StringValue(var132.Name)
						var133.Negate = types.BoolValue(var132.Negate)
						var133.Value = types.StringValue(var132.Value)
						var131 = append(var131, var133)
					}
				}
				var130.Key = var131
				var130.Name = types.StringValue(var129.Name)
				var130.Negate = types.BoolValue(var129.Negate)
				var128 = append(var128, var130)
			}
		}
		var var134 []objectsHipObjectsRsModelProcessListObject
		if len(ans.CustomChecks.Criteria.ProcessList) != 0 {
			var134 = make([]objectsHipObjectsRsModelProcessListObject, 0, len(ans.CustomChecks.Criteria.ProcessList))
			for var135Index := range ans.CustomChecks.Criteria.ProcessList {
				var135 := ans.CustomChecks.Criteria.ProcessList[var135Index]
				var var136 objectsHipObjectsRsModelProcessListObject
				var136.Name = types.StringValue(var135.Name)
				var136.Running = types.BoolValue(var135.Running)
				var134 = append(var134, var136)
			}
		}
		var var137 []objectsHipObjectsRsModelRegistryKeyObject
		if len(ans.CustomChecks.Criteria.RegistryKey) != 0 {
			var137 = make([]objectsHipObjectsRsModelRegistryKeyObject, 0, len(ans.CustomChecks.Criteria.RegistryKey))
			for var138Index := range ans.CustomChecks.Criteria.RegistryKey {
				var138 := ans.CustomChecks.Criteria.RegistryKey[var138Index]
				var var139 objectsHipObjectsRsModelRegistryKeyObject
				var var140 []objectsHipObjectsRsModelRegistryValueObject
				if len(var138.RegistryValue) != 0 {
					var140 = make([]objectsHipObjectsRsModelRegistryValueObject, 0, len(var138.RegistryValue))
					for var141Index := range var138.RegistryValue {
						var141 := var138.RegistryValue[var141Index]
						var var142 objectsHipObjectsRsModelRegistryValueObject
						var142.Name = types.StringValue(var141.Name)
						var142.Negate = types.BoolValue(var141.Negate)
						var142.ValueData = types.StringValue(var141.ValueData)
						var140 = append(var140, var142)
					}
				}
				var139.DefaultValueData = types.StringValue(var138.DefaultValueData)
				var139.Name = types.StringValue(var138.Name)
				var139.Negate = types.BoolValue(var138.Negate)
				var139.RegistryValue = var140
				var137 = append(var137, var139)
			}
		}
		var127.Plist = var128
		var127.ProcessList = var134
		var127.RegistryKey = var137
		var126.Criteria = var127
	}
	var var143 *objectsHipObjectsRsModelDataLossPreventionObject
	if ans.DataLossPrevention != nil {
		var143 = &objectsHipObjectsRsModelDataLossPreventionObject{}
		var var144 *objectsHipObjectsRsModelCriteriaObject3
		if ans.DataLossPrevention.Criteria != nil {
			var144 = &objectsHipObjectsRsModelCriteriaObject3{}
			var144.IsEnabled = types.StringValue(ans.DataLossPrevention.Criteria.IsEnabled)
			var144.IsInstalled = types.BoolValue(ans.DataLossPrevention.Criteria.IsInstalled)
		}
		var var145 []objectsHipObjectsRsModelVendorObject1
		if len(ans.DataLossPrevention.Vendor) != 0 {
			var145 = make([]objectsHipObjectsRsModelVendorObject1, 0, len(ans.DataLossPrevention.Vendor))
			for var146Index := range ans.DataLossPrevention.Vendor {
				var146 := ans.DataLossPrevention.Vendor[var146Index]
				var var147 objectsHipObjectsRsModelVendorObject1
				var147.Name = types.StringValue(var146.Name)
				var147.Product = EncodeStringSlice(var146.Product)
				var145 = append(var145, var147)
			}
		}
		var143.Criteria = var144
		var143.ExcludeVendor = types.BoolValue(ans.DataLossPrevention.ExcludeVendor)
		var143.Vendor = var145
	}
	var var148 *objectsHipObjectsRsModelDiskBackupObject
	if ans.DiskBackup != nil {
		var148 = &objectsHipObjectsRsModelDiskBackupObject{}
		var var149 *objectsHipObjectsRsModelCriteriaObject4
		if ans.DiskBackup.Criteria != nil {
			var149 = &objectsHipObjectsRsModelCriteriaObject4{}
			var var150 *objectsHipObjectsRsModelLastBackupTimeObject
			if ans.DiskBackup.Criteria.LastBackupTime != nil {
				var150 = &objectsHipObjectsRsModelLastBackupTimeObject{}
				var var151 *objectsHipObjectsRsModelNotWithinObject
				if ans.DiskBackup.Criteria.LastBackupTime.NotWithin != nil {
					var151 = &objectsHipObjectsRsModelNotWithinObject{}
					var151.Days = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.NotWithin.Days)
					var151.Hours = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.NotWithin.Hours)
				}
				var var152 *objectsHipObjectsRsModelWithinObject
				if ans.DiskBackup.Criteria.LastBackupTime.Within != nil {
					var152 = &objectsHipObjectsRsModelWithinObject{}
					var152.Days = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.Within.Days)
					var152.Hours = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.Within.Hours)
				}
				if ans.DiskBackup.Criteria.LastBackupTime.NotAvailable != nil {
					var150.NotAvailable = types.BoolValue(true)
				}
				var150.NotWithin = var151
				var150.Within = var152
			}
			var149.IsInstalled = types.BoolValue(ans.DiskBackup.Criteria.IsInstalled)
			var149.LastBackupTime = var150
		}
		var var153 []objectsHipObjectsRsModelVendorObject
		if len(ans.DiskBackup.Vendor) != 0 {
			var153 = make([]objectsHipObjectsRsModelVendorObject, 0, len(ans.DiskBackup.Vendor))
			for var154Index := range ans.DiskBackup.Vendor {
				var154 := ans.DiskBackup.Vendor[var154Index]
				var var155 objectsHipObjectsRsModelVendorObject
				var155.Name = types.StringValue(var154.Name)
				var155.Product = EncodeStringSlice(var154.Product)
				var153 = append(var153, var155)
			}
		}
		var148.Criteria = var149
		var148.ExcludeVendor = types.BoolValue(ans.DiskBackup.ExcludeVendor)
		var148.Vendor = var153
	}
	var var156 *objectsHipObjectsRsModelDiskEncryptionObject
	if ans.DiskEncryption != nil {
		var156 = &objectsHipObjectsRsModelDiskEncryptionObject{}
		var var157 *objectsHipObjectsRsModelCriteriaObject5
		if ans.DiskEncryption.Criteria != nil {
			var157 = &objectsHipObjectsRsModelCriteriaObject5{}
			var var158 []objectsHipObjectsRsModelEncryptedLocationsObject
			if len(ans.DiskEncryption.Criteria.EncryptedLocations) != 0 {
				var158 = make([]objectsHipObjectsRsModelEncryptedLocationsObject, 0, len(ans.DiskEncryption.Criteria.EncryptedLocations))
				for var159Index := range ans.DiskEncryption.Criteria.EncryptedLocations {
					var159 := ans.DiskEncryption.Criteria.EncryptedLocations[var159Index]
					var var160 objectsHipObjectsRsModelEncryptedLocationsObject
					var var161 *objectsHipObjectsRsModelEncryptionStateObject
					if var159.EncryptionState != nil {
						var161 = &objectsHipObjectsRsModelEncryptionStateObject{}
						var161.Is = types.StringValue(var159.EncryptionState.Is)
						var161.IsNot = types.StringValue(var159.EncryptionState.IsNot)
					}
					var160.EncryptionState = var161
					var160.Name = types.StringValue(var159.Name)
					var158 = append(var158, var160)
				}
			}
			var157.EncryptedLocations = var158
			var157.IsInstalled = types.BoolValue(ans.DiskEncryption.Criteria.IsInstalled)
		}
		var var162 []objectsHipObjectsRsModelVendorObject
		if len(ans.DiskEncryption.Vendor) != 0 {
			var162 = make([]objectsHipObjectsRsModelVendorObject, 0, len(ans.DiskEncryption.Vendor))
			for var163Index := range ans.DiskEncryption.Vendor {
				var163 := ans.DiskEncryption.Vendor[var163Index]
				var var164 objectsHipObjectsRsModelVendorObject
				var164.Name = types.StringValue(var163.Name)
				var164.Product = EncodeStringSlice(var163.Product)
				var162 = append(var162, var164)
			}
		}
		var156.Criteria = var157
		var156.ExcludeVendor = types.BoolValue(ans.DiskEncryption.ExcludeVendor)
		var156.Vendor = var162
	}
	var var165 *objectsHipObjectsRsModelFirewallObject
	if ans.Firewall != nil {
		var165 = &objectsHipObjectsRsModelFirewallObject{}
		var var166 *objectsHipObjectsRsModelCriteriaObject3
		if ans.Firewall.Criteria != nil {
			var166 = &objectsHipObjectsRsModelCriteriaObject3{}
			var166.IsEnabled = types.StringValue(ans.Firewall.Criteria.IsEnabled)
			var166.IsInstalled = types.BoolValue(ans.Firewall.Criteria.IsInstalled)
		}
		var var167 []objectsHipObjectsRsModelVendorObject
		if len(ans.Firewall.Vendor) != 0 {
			var167 = make([]objectsHipObjectsRsModelVendorObject, 0, len(ans.Firewall.Vendor))
			for var168Index := range ans.Firewall.Vendor {
				var168 := ans.Firewall.Vendor[var168Index]
				var var169 objectsHipObjectsRsModelVendorObject
				var169.Name = types.StringValue(var168.Name)
				var169.Product = EncodeStringSlice(var168.Product)
				var167 = append(var167, var169)
			}
		}
		var165.Criteria = var166
		var165.ExcludeVendor = types.BoolValue(ans.Firewall.ExcludeVendor)
		var165.Vendor = var167
	}
	var var170 *objectsHipObjectsRsModelHostInfoObject
	if ans.HostInfo != nil {
		var170 = &objectsHipObjectsRsModelHostInfoObject{}
		var var171 objectsHipObjectsRsModelCriteriaObject6
		var var172 *objectsHipObjectsRsModelClientVersionObject
		if ans.HostInfo.Criteria.ClientVersion != nil {
			var172 = &objectsHipObjectsRsModelClientVersionObject{}
			var172.Contains = types.StringValue(ans.HostInfo.Criteria.ClientVersion.Contains)
			var172.Is = types.StringValue(ans.HostInfo.Criteria.ClientVersion.Is)
			var172.IsNot = types.StringValue(ans.HostInfo.Criteria.ClientVersion.IsNot)
		}
		var var173 *objectsHipObjectsRsModelDomainObject
		if ans.HostInfo.Criteria.Domain != nil {
			var173 = &objectsHipObjectsRsModelDomainObject{}
			var173.Contains = types.StringValue(ans.HostInfo.Criteria.Domain.Contains)
			var173.Is = types.StringValue(ans.HostInfo.Criteria.Domain.Is)
			var173.IsNot = types.StringValue(ans.HostInfo.Criteria.Domain.IsNot)
		}
		var var174 *objectsHipObjectsRsModelHostIdObject
		if ans.HostInfo.Criteria.HostId != nil {
			var174 = &objectsHipObjectsRsModelHostIdObject{}
			var174.Contains = types.StringValue(ans.HostInfo.Criteria.HostId.Contains)
			var174.Is = types.StringValue(ans.HostInfo.Criteria.HostId.Is)
			var174.IsNot = types.StringValue(ans.HostInfo.Criteria.HostId.IsNot)
		}
		var var175 *objectsHipObjectsRsModelHostNameObject
		if ans.HostInfo.Criteria.HostName != nil {
			var175 = &objectsHipObjectsRsModelHostNameObject{}
			var175.Contains = types.StringValue(ans.HostInfo.Criteria.HostName.Contains)
			var175.Is = types.StringValue(ans.HostInfo.Criteria.HostName.Is)
			var175.IsNot = types.StringValue(ans.HostInfo.Criteria.HostName.IsNot)
		}
		var var176 *objectsHipObjectsRsModelOsObject
		if ans.HostInfo.Criteria.Os != nil {
			var176 = &objectsHipObjectsRsModelOsObject{}
			var var177 *objectsHipObjectsRsModelContainsObject
			if ans.HostInfo.Criteria.Os.Contains != nil {
				var177 = &objectsHipObjectsRsModelContainsObject{}
				var177.Apple = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Apple)
				var177.Google = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Google)
				var177.Linux = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Linux)
				var177.Microsoft = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Microsoft)
				var177.Other = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Other)
			}
			var176.Contains = var177
		}
		var var178 *objectsHipObjectsRsModelSerialNumberObject
		if ans.HostInfo.Criteria.SerialNumber != nil {
			var178 = &objectsHipObjectsRsModelSerialNumberObject{}
			var178.Contains = types.StringValue(ans.HostInfo.Criteria.SerialNumber.Contains)
			var178.Is = types.StringValue(ans.HostInfo.Criteria.SerialNumber.Is)
			var178.IsNot = types.StringValue(ans.HostInfo.Criteria.SerialNumber.IsNot)
		}
		var171.ClientVersion = var172
		var171.Domain = var173
		var171.HostId = var174
		var171.HostName = var175
		var171.Managed = types.BoolValue(ans.HostInfo.Criteria.Managed)
		var171.Os = var176
		var171.SerialNumber = var178
		var170.Criteria = var171
	}
	var var179 *objectsHipObjectsRsModelMobileDeviceObject
	if ans.MobileDevice != nil {
		var179 = &objectsHipObjectsRsModelMobileDeviceObject{}
		var var180 *objectsHipObjectsRsModelCriteriaObject7
		if ans.MobileDevice.Criteria != nil {
			var180 = &objectsHipObjectsRsModelCriteriaObject7{}
			var var181 *objectsHipObjectsRsModelApplicationsObject
			if ans.MobileDevice.Criteria.Applications != nil {
				var181 = &objectsHipObjectsRsModelApplicationsObject{}
				var var182 *objectsHipObjectsRsModelHasMalwareObject
				if ans.MobileDevice.Criteria.Applications.HasMalware != nil {
					var182 = &objectsHipObjectsRsModelHasMalwareObject{}
					var var183 *objectsHipObjectsRsModelYesObject
					if ans.MobileDevice.Criteria.Applications.HasMalware.Yes != nil {
						var183 = &objectsHipObjectsRsModelYesObject{}
						var var184 []objectsHipObjectsRsModelExcludesObject
						if len(ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes) != 0 {
							var184 = make([]objectsHipObjectsRsModelExcludesObject, 0, len(ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes))
							for var185Index := range ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes {
								var185 := ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes[var185Index]
								var var186 objectsHipObjectsRsModelExcludesObject
								var186.Hash = types.StringValue(var185.Hash)
								var186.Name = types.StringValue(var185.Name)
								var186.Package = types.StringValue(var185.Package)
								var184 = append(var184, var186)
							}
						}
						var183.Excludes = var184
					}
					if ans.MobileDevice.Criteria.Applications.HasMalware.No != nil {
						var182.No = types.BoolValue(true)
					}
					var182.Yes = var183
				}
				var var187 []objectsHipObjectsRsModelIncludesObject
				if len(ans.MobileDevice.Criteria.Applications.Includes) != 0 {
					var187 = make([]objectsHipObjectsRsModelIncludesObject, 0, len(ans.MobileDevice.Criteria.Applications.Includes))
					for var188Index := range ans.MobileDevice.Criteria.Applications.Includes {
						var188 := ans.MobileDevice.Criteria.Applications.Includes[var188Index]
						var var189 objectsHipObjectsRsModelIncludesObject
						var189.Hash = types.StringValue(var188.Hash)
						var189.Name = types.StringValue(var188.Name)
						var189.Package = types.StringValue(var188.Package)
						var187 = append(var187, var189)
					}
				}
				var181.HasMalware = var182
				var181.HasUnmanagedApp = types.BoolValue(ans.MobileDevice.Criteria.Applications.HasUnmanagedApp)
				var181.Includes = var187
			}
			var var190 *objectsHipObjectsRsModelImeiObject
			if ans.MobileDevice.Criteria.Imei != nil {
				var190 = &objectsHipObjectsRsModelImeiObject{}
				var190.Contains = types.StringValue(ans.MobileDevice.Criteria.Imei.Contains)
				var190.Is = types.StringValue(ans.MobileDevice.Criteria.Imei.Is)
				var190.IsNot = types.StringValue(ans.MobileDevice.Criteria.Imei.IsNot)
			}
			var var191 *objectsHipObjectsRsModelLastCheckinTimeObject
			if ans.MobileDevice.Criteria.LastCheckinTime != nil {
				var191 = &objectsHipObjectsRsModelLastCheckinTimeObject{}
				var var192 *objectsHipObjectsRsModelNotWithinObject3
				if ans.MobileDevice.Criteria.LastCheckinTime.NotWithin != nil {
					var192 = &objectsHipObjectsRsModelNotWithinObject3{}
					var192.Days = types.Int64Value(ans.MobileDevice.Criteria.LastCheckinTime.NotWithin.Days)
				}
				var var193 *objectsHipObjectsRsModelWithinObject3
				if ans.MobileDevice.Criteria.LastCheckinTime.Within != nil {
					var193 = &objectsHipObjectsRsModelWithinObject3{}
					var193.Days = types.Int64Value(ans.MobileDevice.Criteria.LastCheckinTime.Within.Days)
				}
				var191.NotWithin = var192
				var191.Within = var193
			}
			var var194 *objectsHipObjectsRsModelModelObject
			if ans.MobileDevice.Criteria.Model != nil {
				var194 = &objectsHipObjectsRsModelModelObject{}
				var194.Contains = types.StringValue(ans.MobileDevice.Criteria.Model.Contains)
				var194.Is = types.StringValue(ans.MobileDevice.Criteria.Model.Is)
				var194.IsNot = types.StringValue(ans.MobileDevice.Criteria.Model.IsNot)
			}
			var var195 *objectsHipObjectsRsModelPhoneNumberObject
			if ans.MobileDevice.Criteria.PhoneNumber != nil {
				var195 = &objectsHipObjectsRsModelPhoneNumberObject{}
				var195.Contains = types.StringValue(ans.MobileDevice.Criteria.PhoneNumber.Contains)
				var195.Is = types.StringValue(ans.MobileDevice.Criteria.PhoneNumber.Is)
				var195.IsNot = types.StringValue(ans.MobileDevice.Criteria.PhoneNumber.IsNot)
			}
			var var196 *objectsHipObjectsRsModelTagObject
			if ans.MobileDevice.Criteria.Tag != nil {
				var196 = &objectsHipObjectsRsModelTagObject{}
				var196.Contains = types.StringValue(ans.MobileDevice.Criteria.Tag.Contains)
				var196.Is = types.StringValue(ans.MobileDevice.Criteria.Tag.Is)
				var196.IsNot = types.StringValue(ans.MobileDevice.Criteria.Tag.IsNot)
			}
			var180.Applications = var181
			var180.DiskEncrypted = types.BoolValue(ans.MobileDevice.Criteria.DiskEncrypted)
			var180.Imei = var190
			var180.Jailbroken = types.BoolValue(ans.MobileDevice.Criteria.Jailbroken)
			var180.LastCheckinTime = var191
			var180.Model = var194
			var180.PasscodeSet = types.BoolValue(ans.MobileDevice.Criteria.PasscodeSet)
			var180.PhoneNumber = var195
			var180.Tag = var196
		}
		var179.Criteria = var180
	}
	var var197 *objectsHipObjectsRsModelNetworkInfoObject
	if ans.NetworkInfo != nil {
		var197 = &objectsHipObjectsRsModelNetworkInfoObject{}
		var var198 *objectsHipObjectsRsModelCriteriaObject8
		if ans.NetworkInfo.Criteria != nil {
			var198 = &objectsHipObjectsRsModelCriteriaObject8{}
			var var199 *objectsHipObjectsRsModelNetworkObject
			if ans.NetworkInfo.Criteria.Network != nil {
				var199 = &objectsHipObjectsRsModelNetworkObject{}
				var var200 *objectsHipObjectsRsModelIsObject
				if ans.NetworkInfo.Criteria.Network.Is != nil {
					var200 = &objectsHipObjectsRsModelIsObject{}
					var var201 *objectsHipObjectsRsModelMobileObject
					if ans.NetworkInfo.Criteria.Network.Is.Mobile != nil {
						var201 = &objectsHipObjectsRsModelMobileObject{}
						var201.Carrier = types.StringValue(ans.NetworkInfo.Criteria.Network.Is.Mobile.Carrier)
					}
					var var202 *objectsHipObjectsRsModelWifiObject
					if ans.NetworkInfo.Criteria.Network.Is.Wifi != nil {
						var202 = &objectsHipObjectsRsModelWifiObject{}
						var202.Ssid = types.StringValue(ans.NetworkInfo.Criteria.Network.Is.Wifi.Ssid)
					}
					var200.Mobile = var201
					if ans.NetworkInfo.Criteria.Network.Is.Unknown != nil {
						var200.Unknown = types.BoolValue(true)
					}
					var200.Wifi = var202
				}
				var var203 *objectsHipObjectsRsModelIsNotObject
				if ans.NetworkInfo.Criteria.Network.IsNot != nil {
					var203 = &objectsHipObjectsRsModelIsNotObject{}
					var var204 *objectsHipObjectsRsModelMobileObject
					if ans.NetworkInfo.Criteria.Network.IsNot.Mobile != nil {
						var204 = &objectsHipObjectsRsModelMobileObject{}
						var204.Carrier = types.StringValue(ans.NetworkInfo.Criteria.Network.IsNot.Mobile.Carrier)
					}
					var var205 *objectsHipObjectsRsModelWifiObject
					if ans.NetworkInfo.Criteria.Network.IsNot.Wifi != nil {
						var205 = &objectsHipObjectsRsModelWifiObject{}
						var205.Ssid = types.StringValue(ans.NetworkInfo.Criteria.Network.IsNot.Wifi.Ssid)
					}
					if ans.NetworkInfo.Criteria.Network.IsNot.Ethernet != nil {
						var203.Ethernet = types.BoolValue(true)
					}
					var203.Mobile = var204
					if ans.NetworkInfo.Criteria.Network.IsNot.Unknown != nil {
						var203.Unknown = types.BoolValue(true)
					}
					var203.Wifi = var205
				}
				var199.Is = var200
				var199.IsNot = var203
			}
			var198.Network = var199
		}
		var197.Criteria = var198
	}
	var var206 *objectsHipObjectsRsModelPatchManagementObject
	if ans.PatchManagement != nil {
		var206 = &objectsHipObjectsRsModelPatchManagementObject{}
		var var207 *objectsHipObjectsRsModelCriteriaObject9
		if ans.PatchManagement.Criteria != nil {
			var207 = &objectsHipObjectsRsModelCriteriaObject9{}
			var var208 *objectsHipObjectsRsModelMissingPatchesObject
			if ans.PatchManagement.Criteria.MissingPatches != nil {
				var208 = &objectsHipObjectsRsModelMissingPatchesObject{}
				var var209 *objectsHipObjectsRsModelSeverityObject
				if ans.PatchManagement.Criteria.MissingPatches.Severity != nil {
					var209 = &objectsHipObjectsRsModelSeverityObject{}
					var209.GreaterEqual = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.GreaterEqual)
					var209.GreaterThan = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.GreaterThan)
					var209.Is = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.Is)
					var209.IsNot = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.IsNot)
					var209.LessEqual = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.LessEqual)
					var209.LessThan = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.LessThan)
				}
				var208.Check = types.StringValue(ans.PatchManagement.Criteria.MissingPatches.Check)
				var208.Patches = EncodeStringSlice(ans.PatchManagement.Criteria.MissingPatches.Patches)
				var208.Severity = var209
			}
			var207.IsEnabled = types.StringValue(ans.PatchManagement.Criteria.IsEnabled)
			var207.IsInstalled = types.BoolValue(ans.PatchManagement.Criteria.IsInstalled)
			var207.MissingPatches = var208
		}
		var var210 []objectsHipObjectsRsModelVendorObject1
		if len(ans.PatchManagement.Vendor) != 0 {
			var210 = make([]objectsHipObjectsRsModelVendorObject1, 0, len(ans.PatchManagement.Vendor))
			for var211Index := range ans.PatchManagement.Vendor {
				var211 := ans.PatchManagement.Vendor[var211Index]
				var var212 objectsHipObjectsRsModelVendorObject1
				var212.Name = types.StringValue(var211.Name)
				var212.Product = EncodeStringSlice(var211.Product)
				var210 = append(var210, var212)
			}
		}
		var206.Criteria = var207
		var206.ExcludeVendor = types.BoolValue(ans.PatchManagement.ExcludeVendor)
		var206.Vendor = var210
	}
	state.AntiMalware = var107
	state.Certificate = var121
	state.CustomChecks = var126
	state.DataLossPrevention = var143
	state.Description = types.StringValue(ans.Description)
	state.DiskBackup = var148
	state.DiskEncryption = var156
	state.Firewall = var165
	state.HostInfo = var170
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MobileDevice = var179
	state.Name = types.StringValue(ans.Name)
	state.NetworkInfo = var197
	state.PatchManagement = var206

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *objectsHipObjectsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var idType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := idType.ValueString()
	tokens := strings.Split(id, IdSeparator)
	if len(tokens) != 2 {
		resp.Diagnostics.AddError("Error in resource ID format", "Expected 2 tokens")
		return
	}

	var state objectsHipObjectsRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"resource_name": "sase_objects_hip_objects",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	// Prepare to read the config.
	svc := yCYVNEN.NewClient(r.client)
	input := yCYVNEN.ReadInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	ans, err := svc.Read(ctx, input)
	if err != nil {
		if IsObjectNotFound(err) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Error reading config", err.Error())
		}
		return
	}

	// Store the answer to state.
	state.Folder = types.StringValue(tokens[0])
	state.Id = idType
	var var0 *objectsHipObjectsRsModelAntiMalwareObject
	if ans.AntiMalware != nil {
		var0 = &objectsHipObjectsRsModelAntiMalwareObject{}
		var var1 *objectsHipObjectsRsModelCriteriaObject
		if ans.AntiMalware.Criteria != nil {
			var1 = &objectsHipObjectsRsModelCriteriaObject{}
			var var2 *objectsHipObjectsRsModelLastScanTimeObject
			if ans.AntiMalware.Criteria.LastScanTime != nil {
				var2 = &objectsHipObjectsRsModelLastScanTimeObject{}
				var var3 *objectsHipObjectsRsModelNotWithinObject
				if ans.AntiMalware.Criteria.LastScanTime.NotWithin != nil {
					var3 = &objectsHipObjectsRsModelNotWithinObject{}
					var3.Days = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.NotWithin.Days)
					var3.Hours = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.NotWithin.Hours)
				}
				var var4 *objectsHipObjectsRsModelWithinObject
				if ans.AntiMalware.Criteria.LastScanTime.Within != nil {
					var4 = &objectsHipObjectsRsModelWithinObject{}
					var4.Days = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.Within.Days)
					var4.Hours = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.Within.Hours)
				}
				if ans.AntiMalware.Criteria.LastScanTime.NotAvailable != nil {
					var2.NotAvailable = types.BoolValue(true)
				}
				var2.NotWithin = var3
				var2.Within = var4
			}
			var var5 *objectsHipObjectsRsModelProductVersionObject
			if ans.AntiMalware.Criteria.ProductVersion != nil {
				var5 = &objectsHipObjectsRsModelProductVersionObject{}
				var var6 *objectsHipObjectsRsModelNotWithinObject1
				if ans.AntiMalware.Criteria.ProductVersion.NotWithin != nil {
					var6 = &objectsHipObjectsRsModelNotWithinObject1{}
					var6.Versions = types.Int64Value(ans.AntiMalware.Criteria.ProductVersion.NotWithin.Versions)
				}
				var var7 *objectsHipObjectsRsModelWithinObject1
				if ans.AntiMalware.Criteria.ProductVersion.Within != nil {
					var7 = &objectsHipObjectsRsModelWithinObject1{}
					var7.Versions = types.Int64Value(ans.AntiMalware.Criteria.ProductVersion.Within.Versions)
				}
				var5.Contains = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.Contains)
				var5.GreaterEqual = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.GreaterEqual)
				var5.GreaterThan = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.GreaterThan)
				var5.Is = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.Is)
				var5.IsNot = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.IsNot)
				var5.LessEqual = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.LessEqual)
				var5.LessThan = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.LessThan)
				var5.NotWithin = var6
				var5.Within = var7
			}
			var var8 *objectsHipObjectsRsModelVirdefVersionObject
			if ans.AntiMalware.Criteria.VirdefVersion != nil {
				var8 = &objectsHipObjectsRsModelVirdefVersionObject{}
				var var9 *objectsHipObjectsRsModelNotWithinObject2
				if ans.AntiMalware.Criteria.VirdefVersion.NotWithin != nil {
					var9 = &objectsHipObjectsRsModelNotWithinObject2{}
					var9.Days = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.NotWithin.Days)
					var9.Versions = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.NotWithin.Versions)
				}
				var var10 *objectsHipObjectsRsModelWithinObject2
				if ans.AntiMalware.Criteria.VirdefVersion.Within != nil {
					var10 = &objectsHipObjectsRsModelWithinObject2{}
					var10.Days = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.Within.Days)
					var10.Versions = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.Within.Versions)
				}
				var8.NotWithin = var9
				var8.Within = var10
			}
			var1.IsInstalled = types.BoolValue(ans.AntiMalware.Criteria.IsInstalled)
			var1.LastScanTime = var2
			var1.ProductVersion = var5
			var1.RealTimeProtection = types.StringValue(ans.AntiMalware.Criteria.RealTimeProtection)
			var1.VirdefVersion = var8
		}
		var var11 []objectsHipObjectsRsModelVendorObject
		if len(ans.AntiMalware.Vendor) != 0 {
			var11 = make([]objectsHipObjectsRsModelVendorObject, 0, len(ans.AntiMalware.Vendor))
			for var12Index := range ans.AntiMalware.Vendor {
				var12 := ans.AntiMalware.Vendor[var12Index]
				var var13 objectsHipObjectsRsModelVendorObject
				var13.Name = types.StringValue(var12.Name)
				var13.Product = EncodeStringSlice(var12.Product)
				var11 = append(var11, var13)
			}
		}
		var0.Criteria = var1
		var0.ExcludeVendor = types.BoolValue(ans.AntiMalware.ExcludeVendor)
		var0.Vendor = var11
	}
	var var14 *objectsHipObjectsRsModelCertificateObject
	if ans.Certificate != nil {
		var14 = &objectsHipObjectsRsModelCertificateObject{}
		var var15 *objectsHipObjectsRsModelCriteriaObject1
		if ans.Certificate.Criteria != nil {
			var15 = &objectsHipObjectsRsModelCriteriaObject1{}
			var var16 []objectsHipObjectsRsModelCertificateAttributesObject
			if len(ans.Certificate.Criteria.CertificateAttributes) != 0 {
				var16 = make([]objectsHipObjectsRsModelCertificateAttributesObject, 0, len(ans.Certificate.Criteria.CertificateAttributes))
				for var17Index := range ans.Certificate.Criteria.CertificateAttributes {
					var17 := ans.Certificate.Criteria.CertificateAttributes[var17Index]
					var var18 objectsHipObjectsRsModelCertificateAttributesObject
					var18.Name = types.StringValue(var17.Name)
					var18.Value = types.StringValue(var17.Value)
					var16 = append(var16, var18)
				}
			}
			var15.CertificateAttributes = var16
			var15.CertificateProfile = types.StringValue(ans.Certificate.Criteria.CertificateProfile)
		}
		var14.Criteria = var15
	}
	var var19 *objectsHipObjectsRsModelCustomChecksObject
	if ans.CustomChecks != nil {
		var19 = &objectsHipObjectsRsModelCustomChecksObject{}
		var var20 objectsHipObjectsRsModelCriteriaObject2
		var var21 []objectsHipObjectsRsModelPlistObject
		if len(ans.CustomChecks.Criteria.Plist) != 0 {
			var21 = make([]objectsHipObjectsRsModelPlistObject, 0, len(ans.CustomChecks.Criteria.Plist))
			for var22Index := range ans.CustomChecks.Criteria.Plist {
				var22 := ans.CustomChecks.Criteria.Plist[var22Index]
				var var23 objectsHipObjectsRsModelPlistObject
				var var24 []objectsHipObjectsRsModelKeyObject
				if len(var22.Key) != 0 {
					var24 = make([]objectsHipObjectsRsModelKeyObject, 0, len(var22.Key))
					for var25Index := range var22.Key {
						var25 := var22.Key[var25Index]
						var var26 objectsHipObjectsRsModelKeyObject
						var26.Name = types.StringValue(var25.Name)
						var26.Negate = types.BoolValue(var25.Negate)
						var26.Value = types.StringValue(var25.Value)
						var24 = append(var24, var26)
					}
				}
				var23.Key = var24
				var23.Name = types.StringValue(var22.Name)
				var23.Negate = types.BoolValue(var22.Negate)
				var21 = append(var21, var23)
			}
		}
		var var27 []objectsHipObjectsRsModelProcessListObject
		if len(ans.CustomChecks.Criteria.ProcessList) != 0 {
			var27 = make([]objectsHipObjectsRsModelProcessListObject, 0, len(ans.CustomChecks.Criteria.ProcessList))
			for var28Index := range ans.CustomChecks.Criteria.ProcessList {
				var28 := ans.CustomChecks.Criteria.ProcessList[var28Index]
				var var29 objectsHipObjectsRsModelProcessListObject
				var29.Name = types.StringValue(var28.Name)
				var29.Running = types.BoolValue(var28.Running)
				var27 = append(var27, var29)
			}
		}
		var var30 []objectsHipObjectsRsModelRegistryKeyObject
		if len(ans.CustomChecks.Criteria.RegistryKey) != 0 {
			var30 = make([]objectsHipObjectsRsModelRegistryKeyObject, 0, len(ans.CustomChecks.Criteria.RegistryKey))
			for var31Index := range ans.CustomChecks.Criteria.RegistryKey {
				var31 := ans.CustomChecks.Criteria.RegistryKey[var31Index]
				var var32 objectsHipObjectsRsModelRegistryKeyObject
				var var33 []objectsHipObjectsRsModelRegistryValueObject
				if len(var31.RegistryValue) != 0 {
					var33 = make([]objectsHipObjectsRsModelRegistryValueObject, 0, len(var31.RegistryValue))
					for var34Index := range var31.RegistryValue {
						var34 := var31.RegistryValue[var34Index]
						var var35 objectsHipObjectsRsModelRegistryValueObject
						var35.Name = types.StringValue(var34.Name)
						var35.Negate = types.BoolValue(var34.Negate)
						var35.ValueData = types.StringValue(var34.ValueData)
						var33 = append(var33, var35)
					}
				}
				var32.DefaultValueData = types.StringValue(var31.DefaultValueData)
				var32.Name = types.StringValue(var31.Name)
				var32.Negate = types.BoolValue(var31.Negate)
				var32.RegistryValue = var33
				var30 = append(var30, var32)
			}
		}
		var20.Plist = var21
		var20.ProcessList = var27
		var20.RegistryKey = var30
		var19.Criteria = var20
	}
	var var36 *objectsHipObjectsRsModelDataLossPreventionObject
	if ans.DataLossPrevention != nil {
		var36 = &objectsHipObjectsRsModelDataLossPreventionObject{}
		var var37 *objectsHipObjectsRsModelCriteriaObject3
		if ans.DataLossPrevention.Criteria != nil {
			var37 = &objectsHipObjectsRsModelCriteriaObject3{}
			var37.IsEnabled = types.StringValue(ans.DataLossPrevention.Criteria.IsEnabled)
			var37.IsInstalled = types.BoolValue(ans.DataLossPrevention.Criteria.IsInstalled)
		}
		var var38 []objectsHipObjectsRsModelVendorObject1
		if len(ans.DataLossPrevention.Vendor) != 0 {
			var38 = make([]objectsHipObjectsRsModelVendorObject1, 0, len(ans.DataLossPrevention.Vendor))
			for var39Index := range ans.DataLossPrevention.Vendor {
				var39 := ans.DataLossPrevention.Vendor[var39Index]
				var var40 objectsHipObjectsRsModelVendorObject1
				var40.Name = types.StringValue(var39.Name)
				var40.Product = EncodeStringSlice(var39.Product)
				var38 = append(var38, var40)
			}
		}
		var36.Criteria = var37
		var36.ExcludeVendor = types.BoolValue(ans.DataLossPrevention.ExcludeVendor)
		var36.Vendor = var38
	}
	var var41 *objectsHipObjectsRsModelDiskBackupObject
	if ans.DiskBackup != nil {
		var41 = &objectsHipObjectsRsModelDiskBackupObject{}
		var var42 *objectsHipObjectsRsModelCriteriaObject4
		if ans.DiskBackup.Criteria != nil {
			var42 = &objectsHipObjectsRsModelCriteriaObject4{}
			var var43 *objectsHipObjectsRsModelLastBackupTimeObject
			if ans.DiskBackup.Criteria.LastBackupTime != nil {
				var43 = &objectsHipObjectsRsModelLastBackupTimeObject{}
				var var44 *objectsHipObjectsRsModelNotWithinObject
				if ans.DiskBackup.Criteria.LastBackupTime.NotWithin != nil {
					var44 = &objectsHipObjectsRsModelNotWithinObject{}
					var44.Days = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.NotWithin.Days)
					var44.Hours = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.NotWithin.Hours)
				}
				var var45 *objectsHipObjectsRsModelWithinObject
				if ans.DiskBackup.Criteria.LastBackupTime.Within != nil {
					var45 = &objectsHipObjectsRsModelWithinObject{}
					var45.Days = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.Within.Days)
					var45.Hours = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.Within.Hours)
				}
				if ans.DiskBackup.Criteria.LastBackupTime.NotAvailable != nil {
					var43.NotAvailable = types.BoolValue(true)
				}
				var43.NotWithin = var44
				var43.Within = var45
			}
			var42.IsInstalled = types.BoolValue(ans.DiskBackup.Criteria.IsInstalled)
			var42.LastBackupTime = var43
		}
		var var46 []objectsHipObjectsRsModelVendorObject
		if len(ans.DiskBackup.Vendor) != 0 {
			var46 = make([]objectsHipObjectsRsModelVendorObject, 0, len(ans.DiskBackup.Vendor))
			for var47Index := range ans.DiskBackup.Vendor {
				var47 := ans.DiskBackup.Vendor[var47Index]
				var var48 objectsHipObjectsRsModelVendorObject
				var48.Name = types.StringValue(var47.Name)
				var48.Product = EncodeStringSlice(var47.Product)
				var46 = append(var46, var48)
			}
		}
		var41.Criteria = var42
		var41.ExcludeVendor = types.BoolValue(ans.DiskBackup.ExcludeVendor)
		var41.Vendor = var46
	}
	var var49 *objectsHipObjectsRsModelDiskEncryptionObject
	if ans.DiskEncryption != nil {
		var49 = &objectsHipObjectsRsModelDiskEncryptionObject{}
		var var50 *objectsHipObjectsRsModelCriteriaObject5
		if ans.DiskEncryption.Criteria != nil {
			var50 = &objectsHipObjectsRsModelCriteriaObject5{}
			var var51 []objectsHipObjectsRsModelEncryptedLocationsObject
			if len(ans.DiskEncryption.Criteria.EncryptedLocations) != 0 {
				var51 = make([]objectsHipObjectsRsModelEncryptedLocationsObject, 0, len(ans.DiskEncryption.Criteria.EncryptedLocations))
				for var52Index := range ans.DiskEncryption.Criteria.EncryptedLocations {
					var52 := ans.DiskEncryption.Criteria.EncryptedLocations[var52Index]
					var var53 objectsHipObjectsRsModelEncryptedLocationsObject
					var var54 *objectsHipObjectsRsModelEncryptionStateObject
					if var52.EncryptionState != nil {
						var54 = &objectsHipObjectsRsModelEncryptionStateObject{}
						var54.Is = types.StringValue(var52.EncryptionState.Is)
						var54.IsNot = types.StringValue(var52.EncryptionState.IsNot)
					}
					var53.EncryptionState = var54
					var53.Name = types.StringValue(var52.Name)
					var51 = append(var51, var53)
				}
			}
			var50.EncryptedLocations = var51
			var50.IsInstalled = types.BoolValue(ans.DiskEncryption.Criteria.IsInstalled)
		}
		var var55 []objectsHipObjectsRsModelVendorObject
		if len(ans.DiskEncryption.Vendor) != 0 {
			var55 = make([]objectsHipObjectsRsModelVendorObject, 0, len(ans.DiskEncryption.Vendor))
			for var56Index := range ans.DiskEncryption.Vendor {
				var56 := ans.DiskEncryption.Vendor[var56Index]
				var var57 objectsHipObjectsRsModelVendorObject
				var57.Name = types.StringValue(var56.Name)
				var57.Product = EncodeStringSlice(var56.Product)
				var55 = append(var55, var57)
			}
		}
		var49.Criteria = var50
		var49.ExcludeVendor = types.BoolValue(ans.DiskEncryption.ExcludeVendor)
		var49.Vendor = var55
	}
	var var58 *objectsHipObjectsRsModelFirewallObject
	if ans.Firewall != nil {
		var58 = &objectsHipObjectsRsModelFirewallObject{}
		var var59 *objectsHipObjectsRsModelCriteriaObject3
		if ans.Firewall.Criteria != nil {
			var59 = &objectsHipObjectsRsModelCriteriaObject3{}
			var59.IsEnabled = types.StringValue(ans.Firewall.Criteria.IsEnabled)
			var59.IsInstalled = types.BoolValue(ans.Firewall.Criteria.IsInstalled)
		}
		var var60 []objectsHipObjectsRsModelVendorObject
		if len(ans.Firewall.Vendor) != 0 {
			var60 = make([]objectsHipObjectsRsModelVendorObject, 0, len(ans.Firewall.Vendor))
			for var61Index := range ans.Firewall.Vendor {
				var61 := ans.Firewall.Vendor[var61Index]
				var var62 objectsHipObjectsRsModelVendorObject
				var62.Name = types.StringValue(var61.Name)
				var62.Product = EncodeStringSlice(var61.Product)
				var60 = append(var60, var62)
			}
		}
		var58.Criteria = var59
		var58.ExcludeVendor = types.BoolValue(ans.Firewall.ExcludeVendor)
		var58.Vendor = var60
	}
	var var63 *objectsHipObjectsRsModelHostInfoObject
	if ans.HostInfo != nil {
		var63 = &objectsHipObjectsRsModelHostInfoObject{}
		var var64 objectsHipObjectsRsModelCriteriaObject6
		var var65 *objectsHipObjectsRsModelClientVersionObject
		if ans.HostInfo.Criteria.ClientVersion != nil {
			var65 = &objectsHipObjectsRsModelClientVersionObject{}
			var65.Contains = types.StringValue(ans.HostInfo.Criteria.ClientVersion.Contains)
			var65.Is = types.StringValue(ans.HostInfo.Criteria.ClientVersion.Is)
			var65.IsNot = types.StringValue(ans.HostInfo.Criteria.ClientVersion.IsNot)
		}
		var var66 *objectsHipObjectsRsModelDomainObject
		if ans.HostInfo.Criteria.Domain != nil {
			var66 = &objectsHipObjectsRsModelDomainObject{}
			var66.Contains = types.StringValue(ans.HostInfo.Criteria.Domain.Contains)
			var66.Is = types.StringValue(ans.HostInfo.Criteria.Domain.Is)
			var66.IsNot = types.StringValue(ans.HostInfo.Criteria.Domain.IsNot)
		}
		var var67 *objectsHipObjectsRsModelHostIdObject
		if ans.HostInfo.Criteria.HostId != nil {
			var67 = &objectsHipObjectsRsModelHostIdObject{}
			var67.Contains = types.StringValue(ans.HostInfo.Criteria.HostId.Contains)
			var67.Is = types.StringValue(ans.HostInfo.Criteria.HostId.Is)
			var67.IsNot = types.StringValue(ans.HostInfo.Criteria.HostId.IsNot)
		}
		var var68 *objectsHipObjectsRsModelHostNameObject
		if ans.HostInfo.Criteria.HostName != nil {
			var68 = &objectsHipObjectsRsModelHostNameObject{}
			var68.Contains = types.StringValue(ans.HostInfo.Criteria.HostName.Contains)
			var68.Is = types.StringValue(ans.HostInfo.Criteria.HostName.Is)
			var68.IsNot = types.StringValue(ans.HostInfo.Criteria.HostName.IsNot)
		}
		var var69 *objectsHipObjectsRsModelOsObject
		if ans.HostInfo.Criteria.Os != nil {
			var69 = &objectsHipObjectsRsModelOsObject{}
			var var70 *objectsHipObjectsRsModelContainsObject
			if ans.HostInfo.Criteria.Os.Contains != nil {
				var70 = &objectsHipObjectsRsModelContainsObject{}
				var70.Apple = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Apple)
				var70.Google = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Google)
				var70.Linux = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Linux)
				var70.Microsoft = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Microsoft)
				var70.Other = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Other)
			}
			var69.Contains = var70
		}
		var var71 *objectsHipObjectsRsModelSerialNumberObject
		if ans.HostInfo.Criteria.SerialNumber != nil {
			var71 = &objectsHipObjectsRsModelSerialNumberObject{}
			var71.Contains = types.StringValue(ans.HostInfo.Criteria.SerialNumber.Contains)
			var71.Is = types.StringValue(ans.HostInfo.Criteria.SerialNumber.Is)
			var71.IsNot = types.StringValue(ans.HostInfo.Criteria.SerialNumber.IsNot)
		}
		var64.ClientVersion = var65
		var64.Domain = var66
		var64.HostId = var67
		var64.HostName = var68
		var64.Managed = types.BoolValue(ans.HostInfo.Criteria.Managed)
		var64.Os = var69
		var64.SerialNumber = var71
		var63.Criteria = var64
	}
	var var72 *objectsHipObjectsRsModelMobileDeviceObject
	if ans.MobileDevice != nil {
		var72 = &objectsHipObjectsRsModelMobileDeviceObject{}
		var var73 *objectsHipObjectsRsModelCriteriaObject7
		if ans.MobileDevice.Criteria != nil {
			var73 = &objectsHipObjectsRsModelCriteriaObject7{}
			var var74 *objectsHipObjectsRsModelApplicationsObject
			if ans.MobileDevice.Criteria.Applications != nil {
				var74 = &objectsHipObjectsRsModelApplicationsObject{}
				var var75 *objectsHipObjectsRsModelHasMalwareObject
				if ans.MobileDevice.Criteria.Applications.HasMalware != nil {
					var75 = &objectsHipObjectsRsModelHasMalwareObject{}
					var var76 *objectsHipObjectsRsModelYesObject
					if ans.MobileDevice.Criteria.Applications.HasMalware.Yes != nil {
						var76 = &objectsHipObjectsRsModelYesObject{}
						var var77 []objectsHipObjectsRsModelExcludesObject
						if len(ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes) != 0 {
							var77 = make([]objectsHipObjectsRsModelExcludesObject, 0, len(ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes))
							for var78Index := range ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes {
								var78 := ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes[var78Index]
								var var79 objectsHipObjectsRsModelExcludesObject
								var79.Hash = types.StringValue(var78.Hash)
								var79.Name = types.StringValue(var78.Name)
								var79.Package = types.StringValue(var78.Package)
								var77 = append(var77, var79)
							}
						}
						var76.Excludes = var77
					}
					if ans.MobileDevice.Criteria.Applications.HasMalware.No != nil {
						var75.No = types.BoolValue(true)
					}
					var75.Yes = var76
				}
				var var80 []objectsHipObjectsRsModelIncludesObject
				if len(ans.MobileDevice.Criteria.Applications.Includes) != 0 {
					var80 = make([]objectsHipObjectsRsModelIncludesObject, 0, len(ans.MobileDevice.Criteria.Applications.Includes))
					for var81Index := range ans.MobileDevice.Criteria.Applications.Includes {
						var81 := ans.MobileDevice.Criteria.Applications.Includes[var81Index]
						var var82 objectsHipObjectsRsModelIncludesObject
						var82.Hash = types.StringValue(var81.Hash)
						var82.Name = types.StringValue(var81.Name)
						var82.Package = types.StringValue(var81.Package)
						var80 = append(var80, var82)
					}
				}
				var74.HasMalware = var75
				var74.HasUnmanagedApp = types.BoolValue(ans.MobileDevice.Criteria.Applications.HasUnmanagedApp)
				var74.Includes = var80
			}
			var var83 *objectsHipObjectsRsModelImeiObject
			if ans.MobileDevice.Criteria.Imei != nil {
				var83 = &objectsHipObjectsRsModelImeiObject{}
				var83.Contains = types.StringValue(ans.MobileDevice.Criteria.Imei.Contains)
				var83.Is = types.StringValue(ans.MobileDevice.Criteria.Imei.Is)
				var83.IsNot = types.StringValue(ans.MobileDevice.Criteria.Imei.IsNot)
			}
			var var84 *objectsHipObjectsRsModelLastCheckinTimeObject
			if ans.MobileDevice.Criteria.LastCheckinTime != nil {
				var84 = &objectsHipObjectsRsModelLastCheckinTimeObject{}
				var var85 *objectsHipObjectsRsModelNotWithinObject3
				if ans.MobileDevice.Criteria.LastCheckinTime.NotWithin != nil {
					var85 = &objectsHipObjectsRsModelNotWithinObject3{}
					var85.Days = types.Int64Value(ans.MobileDevice.Criteria.LastCheckinTime.NotWithin.Days)
				}
				var var86 *objectsHipObjectsRsModelWithinObject3
				if ans.MobileDevice.Criteria.LastCheckinTime.Within != nil {
					var86 = &objectsHipObjectsRsModelWithinObject3{}
					var86.Days = types.Int64Value(ans.MobileDevice.Criteria.LastCheckinTime.Within.Days)
				}
				var84.NotWithin = var85
				var84.Within = var86
			}
			var var87 *objectsHipObjectsRsModelModelObject
			if ans.MobileDevice.Criteria.Model != nil {
				var87 = &objectsHipObjectsRsModelModelObject{}
				var87.Contains = types.StringValue(ans.MobileDevice.Criteria.Model.Contains)
				var87.Is = types.StringValue(ans.MobileDevice.Criteria.Model.Is)
				var87.IsNot = types.StringValue(ans.MobileDevice.Criteria.Model.IsNot)
			}
			var var88 *objectsHipObjectsRsModelPhoneNumberObject
			if ans.MobileDevice.Criteria.PhoneNumber != nil {
				var88 = &objectsHipObjectsRsModelPhoneNumberObject{}
				var88.Contains = types.StringValue(ans.MobileDevice.Criteria.PhoneNumber.Contains)
				var88.Is = types.StringValue(ans.MobileDevice.Criteria.PhoneNumber.Is)
				var88.IsNot = types.StringValue(ans.MobileDevice.Criteria.PhoneNumber.IsNot)
			}
			var var89 *objectsHipObjectsRsModelTagObject
			if ans.MobileDevice.Criteria.Tag != nil {
				var89 = &objectsHipObjectsRsModelTagObject{}
				var89.Contains = types.StringValue(ans.MobileDevice.Criteria.Tag.Contains)
				var89.Is = types.StringValue(ans.MobileDevice.Criteria.Tag.Is)
				var89.IsNot = types.StringValue(ans.MobileDevice.Criteria.Tag.IsNot)
			}
			var73.Applications = var74
			var73.DiskEncrypted = types.BoolValue(ans.MobileDevice.Criteria.DiskEncrypted)
			var73.Imei = var83
			var73.Jailbroken = types.BoolValue(ans.MobileDevice.Criteria.Jailbroken)
			var73.LastCheckinTime = var84
			var73.Model = var87
			var73.PasscodeSet = types.BoolValue(ans.MobileDevice.Criteria.PasscodeSet)
			var73.PhoneNumber = var88
			var73.Tag = var89
		}
		var72.Criteria = var73
	}
	var var90 *objectsHipObjectsRsModelNetworkInfoObject
	if ans.NetworkInfo != nil {
		var90 = &objectsHipObjectsRsModelNetworkInfoObject{}
		var var91 *objectsHipObjectsRsModelCriteriaObject8
		if ans.NetworkInfo.Criteria != nil {
			var91 = &objectsHipObjectsRsModelCriteriaObject8{}
			var var92 *objectsHipObjectsRsModelNetworkObject
			if ans.NetworkInfo.Criteria.Network != nil {
				var92 = &objectsHipObjectsRsModelNetworkObject{}
				var var93 *objectsHipObjectsRsModelIsObject
				if ans.NetworkInfo.Criteria.Network.Is != nil {
					var93 = &objectsHipObjectsRsModelIsObject{}
					var var94 *objectsHipObjectsRsModelMobileObject
					if ans.NetworkInfo.Criteria.Network.Is.Mobile != nil {
						var94 = &objectsHipObjectsRsModelMobileObject{}
						var94.Carrier = types.StringValue(ans.NetworkInfo.Criteria.Network.Is.Mobile.Carrier)
					}
					var var95 *objectsHipObjectsRsModelWifiObject
					if ans.NetworkInfo.Criteria.Network.Is.Wifi != nil {
						var95 = &objectsHipObjectsRsModelWifiObject{}
						var95.Ssid = types.StringValue(ans.NetworkInfo.Criteria.Network.Is.Wifi.Ssid)
					}
					var93.Mobile = var94
					if ans.NetworkInfo.Criteria.Network.Is.Unknown != nil {
						var93.Unknown = types.BoolValue(true)
					}
					var93.Wifi = var95
				}
				var var96 *objectsHipObjectsRsModelIsNotObject
				if ans.NetworkInfo.Criteria.Network.IsNot != nil {
					var96 = &objectsHipObjectsRsModelIsNotObject{}
					var var97 *objectsHipObjectsRsModelMobileObject
					if ans.NetworkInfo.Criteria.Network.IsNot.Mobile != nil {
						var97 = &objectsHipObjectsRsModelMobileObject{}
						var97.Carrier = types.StringValue(ans.NetworkInfo.Criteria.Network.IsNot.Mobile.Carrier)
					}
					var var98 *objectsHipObjectsRsModelWifiObject
					if ans.NetworkInfo.Criteria.Network.IsNot.Wifi != nil {
						var98 = &objectsHipObjectsRsModelWifiObject{}
						var98.Ssid = types.StringValue(ans.NetworkInfo.Criteria.Network.IsNot.Wifi.Ssid)
					}
					if ans.NetworkInfo.Criteria.Network.IsNot.Ethernet != nil {
						var96.Ethernet = types.BoolValue(true)
					}
					var96.Mobile = var97
					if ans.NetworkInfo.Criteria.Network.IsNot.Unknown != nil {
						var96.Unknown = types.BoolValue(true)
					}
					var96.Wifi = var98
				}
				var92.Is = var93
				var92.IsNot = var96
			}
			var91.Network = var92
		}
		var90.Criteria = var91
	}
	var var99 *objectsHipObjectsRsModelPatchManagementObject
	if ans.PatchManagement != nil {
		var99 = &objectsHipObjectsRsModelPatchManagementObject{}
		var var100 *objectsHipObjectsRsModelCriteriaObject9
		if ans.PatchManagement.Criteria != nil {
			var100 = &objectsHipObjectsRsModelCriteriaObject9{}
			var var101 *objectsHipObjectsRsModelMissingPatchesObject
			if ans.PatchManagement.Criteria.MissingPatches != nil {
				var101 = &objectsHipObjectsRsModelMissingPatchesObject{}
				var var102 *objectsHipObjectsRsModelSeverityObject
				if ans.PatchManagement.Criteria.MissingPatches.Severity != nil {
					var102 = &objectsHipObjectsRsModelSeverityObject{}
					var102.GreaterEqual = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.GreaterEqual)
					var102.GreaterThan = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.GreaterThan)
					var102.Is = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.Is)
					var102.IsNot = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.IsNot)
					var102.LessEqual = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.LessEqual)
					var102.LessThan = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.LessThan)
				}
				var101.Check = types.StringValue(ans.PatchManagement.Criteria.MissingPatches.Check)
				var101.Patches = EncodeStringSlice(ans.PatchManagement.Criteria.MissingPatches.Patches)
				var101.Severity = var102
			}
			var100.IsEnabled = types.StringValue(ans.PatchManagement.Criteria.IsEnabled)
			var100.IsInstalled = types.BoolValue(ans.PatchManagement.Criteria.IsInstalled)
			var100.MissingPatches = var101
		}
		var var103 []objectsHipObjectsRsModelVendorObject1
		if len(ans.PatchManagement.Vendor) != 0 {
			var103 = make([]objectsHipObjectsRsModelVendorObject1, 0, len(ans.PatchManagement.Vendor))
			for var104Index := range ans.PatchManagement.Vendor {
				var104 := ans.PatchManagement.Vendor[var104Index]
				var var105 objectsHipObjectsRsModelVendorObject1
				var105.Name = types.StringValue(var104.Name)
				var105.Product = EncodeStringSlice(var104.Product)
				var103 = append(var103, var105)
			}
		}
		var99.Criteria = var100
		var99.ExcludeVendor = types.BoolValue(ans.PatchManagement.ExcludeVendor)
		var99.Vendor = var103
	}
	state.AntiMalware = var0
	state.Certificate = var14
	state.CustomChecks = var19
	state.DataLossPrevention = var36
	state.Description = types.StringValue(ans.Description)
	state.DiskBackup = var41
	state.DiskEncryption = var49
	state.Firewall = var58
	state.HostInfo = var63
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MobileDevice = var72
	state.Name = types.StringValue(ans.Name)
	state.NetworkInfo = var90
	state.PatchManagement = var99

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *objectsHipObjectsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state objectsHipObjectsRsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource update", map[string]any{
		"resource_name": "sase_objects_hip_objects",
		"object_id":     state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := yCYVNEN.NewClient(r.client)
	input := yCYVNEN.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 dJpBrWV.Config
	var var1 *dJpBrWV.AntiMalwareObject
	if plan.AntiMalware != nil {
		var1 = &dJpBrWV.AntiMalwareObject{}
		var var2 *dJpBrWV.CriteriaObject
		if plan.AntiMalware.Criteria != nil {
			var2 = &dJpBrWV.CriteriaObject{}
			var2.IsInstalled = plan.AntiMalware.Criteria.IsInstalled.ValueBool()
			var var3 *dJpBrWV.LastScanTimeObject
			if plan.AntiMalware.Criteria.LastScanTime != nil {
				var3 = &dJpBrWV.LastScanTimeObject{}
				if plan.AntiMalware.Criteria.LastScanTime.NotAvailable.ValueBool() {
					var3.NotAvailable = struct{}{}
				}
				var var4 *dJpBrWV.NotWithinObject
				if plan.AntiMalware.Criteria.LastScanTime.NotWithin != nil {
					var4 = &dJpBrWV.NotWithinObject{}
					var4.Days = plan.AntiMalware.Criteria.LastScanTime.NotWithin.Days.ValueInt64()
					var4.Hours = plan.AntiMalware.Criteria.LastScanTime.NotWithin.Hours.ValueInt64()
				}
				var3.NotWithin = var4
				var var5 *dJpBrWV.WithinObject
				if plan.AntiMalware.Criteria.LastScanTime.Within != nil {
					var5 = &dJpBrWV.WithinObject{}
					var5.Days = plan.AntiMalware.Criteria.LastScanTime.Within.Days.ValueInt64()
					var5.Hours = plan.AntiMalware.Criteria.LastScanTime.Within.Hours.ValueInt64()
				}
				var3.Within = var5
			}
			var2.LastScanTime = var3
			var var6 *dJpBrWV.ProductVersionObject
			if plan.AntiMalware.Criteria.ProductVersion != nil {
				var6 = &dJpBrWV.ProductVersionObject{}
				var6.Contains = plan.AntiMalware.Criteria.ProductVersion.Contains.ValueString()
				var6.GreaterEqual = plan.AntiMalware.Criteria.ProductVersion.GreaterEqual.ValueString()
				var6.GreaterThan = plan.AntiMalware.Criteria.ProductVersion.GreaterThan.ValueString()
				var6.Is = plan.AntiMalware.Criteria.ProductVersion.Is.ValueString()
				var6.IsNot = plan.AntiMalware.Criteria.ProductVersion.IsNot.ValueString()
				var6.LessEqual = plan.AntiMalware.Criteria.ProductVersion.LessEqual.ValueString()
				var6.LessThan = plan.AntiMalware.Criteria.ProductVersion.LessThan.ValueString()
				var var7 *dJpBrWV.NotWithinObject1
				if plan.AntiMalware.Criteria.ProductVersion.NotWithin != nil {
					var7 = &dJpBrWV.NotWithinObject1{}
					var7.Versions = plan.AntiMalware.Criteria.ProductVersion.NotWithin.Versions.ValueInt64()
				}
				var6.NotWithin = var7
				var var8 *dJpBrWV.WithinObject1
				if plan.AntiMalware.Criteria.ProductVersion.Within != nil {
					var8 = &dJpBrWV.WithinObject1{}
					var8.Versions = plan.AntiMalware.Criteria.ProductVersion.Within.Versions.ValueInt64()
				}
				var6.Within = var8
			}
			var2.ProductVersion = var6
			var2.RealTimeProtection = plan.AntiMalware.Criteria.RealTimeProtection.ValueString()
			var var9 *dJpBrWV.VirdefVersionObject
			if plan.AntiMalware.Criteria.VirdefVersion != nil {
				var9 = &dJpBrWV.VirdefVersionObject{}
				var var10 *dJpBrWV.NotWithinObject2
				if plan.AntiMalware.Criteria.VirdefVersion.NotWithin != nil {
					var10 = &dJpBrWV.NotWithinObject2{}
					var10.Days = plan.AntiMalware.Criteria.VirdefVersion.NotWithin.Days.ValueInt64()
					var10.Versions = plan.AntiMalware.Criteria.VirdefVersion.NotWithin.Versions.ValueInt64()
				}
				var9.NotWithin = var10
				var var11 *dJpBrWV.WithinObject2
				if plan.AntiMalware.Criteria.VirdefVersion.Within != nil {
					var11 = &dJpBrWV.WithinObject2{}
					var11.Days = plan.AntiMalware.Criteria.VirdefVersion.Within.Days.ValueInt64()
					var11.Versions = plan.AntiMalware.Criteria.VirdefVersion.Within.Versions.ValueInt64()
				}
				var9.Within = var11
			}
			var2.VirdefVersion = var9
		}
		var1.Criteria = var2
		var1.ExcludeVendor = plan.AntiMalware.ExcludeVendor.ValueBool()
		var var12 []dJpBrWV.VendorObject
		if len(plan.AntiMalware.Vendor) != 0 {
			var12 = make([]dJpBrWV.VendorObject, 0, len(plan.AntiMalware.Vendor))
			for var13Index := range plan.AntiMalware.Vendor {
				var13 := plan.AntiMalware.Vendor[var13Index]
				var var14 dJpBrWV.VendorObject
				var14.Name = var13.Name.ValueString()
				var14.Product = DecodeStringSlice(var13.Product)
				var12 = append(var12, var14)
			}
		}
		var1.Vendor = var12
	}
	var0.AntiMalware = var1
	var var15 *dJpBrWV.CertificateObject
	if plan.Certificate != nil {
		var15 = &dJpBrWV.CertificateObject{}
		var var16 *dJpBrWV.CriteriaObject1
		if plan.Certificate.Criteria != nil {
			var16 = &dJpBrWV.CriteriaObject1{}
			var var17 []dJpBrWV.CertificateAttributesObject
			if len(plan.Certificate.Criteria.CertificateAttributes) != 0 {
				var17 = make([]dJpBrWV.CertificateAttributesObject, 0, len(plan.Certificate.Criteria.CertificateAttributes))
				for var18Index := range plan.Certificate.Criteria.CertificateAttributes {
					var18 := plan.Certificate.Criteria.CertificateAttributes[var18Index]
					var var19 dJpBrWV.CertificateAttributesObject
					var19.Name = var18.Name.ValueString()
					var19.Value = var18.Value.ValueString()
					var17 = append(var17, var19)
				}
			}
			var16.CertificateAttributes = var17
			var16.CertificateProfile = plan.Certificate.Criteria.CertificateProfile.ValueString()
		}
		var15.Criteria = var16
	}
	var0.Certificate = var15
	var var20 *dJpBrWV.CustomChecksObject
	if plan.CustomChecks != nil {
		var20 = &dJpBrWV.CustomChecksObject{}
		var var21 dJpBrWV.CriteriaObject2
		var var22 []dJpBrWV.PlistObject
		if len(plan.CustomChecks.Criteria.Plist) != 0 {
			var22 = make([]dJpBrWV.PlistObject, 0, len(plan.CustomChecks.Criteria.Plist))
			for var23Index := range plan.CustomChecks.Criteria.Plist {
				var23 := plan.CustomChecks.Criteria.Plist[var23Index]
				var var24 dJpBrWV.PlistObject
				var var25 []dJpBrWV.KeyObject
				if len(var23.Key) != 0 {
					var25 = make([]dJpBrWV.KeyObject, 0, len(var23.Key))
					for var26Index := range var23.Key {
						var26 := var23.Key[var26Index]
						var var27 dJpBrWV.KeyObject
						var27.Name = var26.Name.ValueString()
						var27.Negate = var26.Negate.ValueBool()
						var27.Value = var26.Value.ValueString()
						var25 = append(var25, var27)
					}
				}
				var24.Key = var25
				var24.Name = var23.Name.ValueString()
				var24.Negate = var23.Negate.ValueBool()
				var22 = append(var22, var24)
			}
		}
		var21.Plist = var22
		var var28 []dJpBrWV.ProcessListObject
		if len(plan.CustomChecks.Criteria.ProcessList) != 0 {
			var28 = make([]dJpBrWV.ProcessListObject, 0, len(plan.CustomChecks.Criteria.ProcessList))
			for var29Index := range plan.CustomChecks.Criteria.ProcessList {
				var29 := plan.CustomChecks.Criteria.ProcessList[var29Index]
				var var30 dJpBrWV.ProcessListObject
				var30.Name = var29.Name.ValueString()
				var30.Running = var29.Running.ValueBool()
				var28 = append(var28, var30)
			}
		}
		var21.ProcessList = var28
		var var31 []dJpBrWV.RegistryKeyObject
		if len(plan.CustomChecks.Criteria.RegistryKey) != 0 {
			var31 = make([]dJpBrWV.RegistryKeyObject, 0, len(plan.CustomChecks.Criteria.RegistryKey))
			for var32Index := range plan.CustomChecks.Criteria.RegistryKey {
				var32 := plan.CustomChecks.Criteria.RegistryKey[var32Index]
				var var33 dJpBrWV.RegistryKeyObject
				var33.DefaultValueData = var32.DefaultValueData.ValueString()
				var33.Name = var32.Name.ValueString()
				var33.Negate = var32.Negate.ValueBool()
				var var34 []dJpBrWV.RegistryValueObject
				if len(var32.RegistryValue) != 0 {
					var34 = make([]dJpBrWV.RegistryValueObject, 0, len(var32.RegistryValue))
					for var35Index := range var32.RegistryValue {
						var35 := var32.RegistryValue[var35Index]
						var var36 dJpBrWV.RegistryValueObject
						var36.Name = var35.Name.ValueString()
						var36.Negate = var35.Negate.ValueBool()
						var36.ValueData = var35.ValueData.ValueString()
						var34 = append(var34, var36)
					}
				}
				var33.RegistryValue = var34
				var31 = append(var31, var33)
			}
		}
		var21.RegistryKey = var31
		var20.Criteria = var21
	}
	var0.CustomChecks = var20
	var var37 *dJpBrWV.DataLossPreventionObject
	if plan.DataLossPrevention != nil {
		var37 = &dJpBrWV.DataLossPreventionObject{}
		var var38 *dJpBrWV.CriteriaObject3
		if plan.DataLossPrevention.Criteria != nil {
			var38 = &dJpBrWV.CriteriaObject3{}
			var38.IsEnabled = plan.DataLossPrevention.Criteria.IsEnabled.ValueString()
			var38.IsInstalled = plan.DataLossPrevention.Criteria.IsInstalled.ValueBool()
		}
		var37.Criteria = var38
		var37.ExcludeVendor = plan.DataLossPrevention.ExcludeVendor.ValueBool()
		var var39 []dJpBrWV.VendorObject1
		if len(plan.DataLossPrevention.Vendor) != 0 {
			var39 = make([]dJpBrWV.VendorObject1, 0, len(plan.DataLossPrevention.Vendor))
			for var40Index := range plan.DataLossPrevention.Vendor {
				var40 := plan.DataLossPrevention.Vendor[var40Index]
				var var41 dJpBrWV.VendorObject1
				var41.Name = var40.Name.ValueString()
				var41.Product = DecodeStringSlice(var40.Product)
				var39 = append(var39, var41)
			}
		}
		var37.Vendor = var39
	}
	var0.DataLossPrevention = var37
	var0.Description = plan.Description.ValueString()
	var var42 *dJpBrWV.DiskBackupObject
	if plan.DiskBackup != nil {
		var42 = &dJpBrWV.DiskBackupObject{}
		var var43 *dJpBrWV.CriteriaObject4
		if plan.DiskBackup.Criteria != nil {
			var43 = &dJpBrWV.CriteriaObject4{}
			var43.IsInstalled = plan.DiskBackup.Criteria.IsInstalled.ValueBool()
			var var44 *dJpBrWV.LastBackupTimeObject
			if plan.DiskBackup.Criteria.LastBackupTime != nil {
				var44 = &dJpBrWV.LastBackupTimeObject{}
				if plan.DiskBackup.Criteria.LastBackupTime.NotAvailable.ValueBool() {
					var44.NotAvailable = struct{}{}
				}
				var var45 *dJpBrWV.NotWithinObject
				if plan.DiskBackup.Criteria.LastBackupTime.NotWithin != nil {
					var45 = &dJpBrWV.NotWithinObject{}
					var45.Days = plan.DiskBackup.Criteria.LastBackupTime.NotWithin.Days.ValueInt64()
					var45.Hours = plan.DiskBackup.Criteria.LastBackupTime.NotWithin.Hours.ValueInt64()
				}
				var44.NotWithin = var45
				var var46 *dJpBrWV.WithinObject
				if plan.DiskBackup.Criteria.LastBackupTime.Within != nil {
					var46 = &dJpBrWV.WithinObject{}
					var46.Days = plan.DiskBackup.Criteria.LastBackupTime.Within.Days.ValueInt64()
					var46.Hours = plan.DiskBackup.Criteria.LastBackupTime.Within.Hours.ValueInt64()
				}
				var44.Within = var46
			}
			var43.LastBackupTime = var44
		}
		var42.Criteria = var43
		var42.ExcludeVendor = plan.DiskBackup.ExcludeVendor.ValueBool()
		var var47 []dJpBrWV.VendorObject
		if len(plan.DiskBackup.Vendor) != 0 {
			var47 = make([]dJpBrWV.VendorObject, 0, len(plan.DiskBackup.Vendor))
			for var48Index := range plan.DiskBackup.Vendor {
				var48 := plan.DiskBackup.Vendor[var48Index]
				var var49 dJpBrWV.VendorObject
				var49.Name = var48.Name.ValueString()
				var49.Product = DecodeStringSlice(var48.Product)
				var47 = append(var47, var49)
			}
		}
		var42.Vendor = var47
	}
	var0.DiskBackup = var42
	var var50 *dJpBrWV.DiskEncryptionObject
	if plan.DiskEncryption != nil {
		var50 = &dJpBrWV.DiskEncryptionObject{}
		var var51 *dJpBrWV.CriteriaObject5
		if plan.DiskEncryption.Criteria != nil {
			var51 = &dJpBrWV.CriteriaObject5{}
			var var52 []dJpBrWV.EncryptedLocationsObject
			if len(plan.DiskEncryption.Criteria.EncryptedLocations) != 0 {
				var52 = make([]dJpBrWV.EncryptedLocationsObject, 0, len(plan.DiskEncryption.Criteria.EncryptedLocations))
				for var53Index := range plan.DiskEncryption.Criteria.EncryptedLocations {
					var53 := plan.DiskEncryption.Criteria.EncryptedLocations[var53Index]
					var var54 dJpBrWV.EncryptedLocationsObject
					var var55 *dJpBrWV.EncryptionStateObject
					if var53.EncryptionState != nil {
						var55 = &dJpBrWV.EncryptionStateObject{}
						var55.Is = var53.EncryptionState.Is.ValueString()
						var55.IsNot = var53.EncryptionState.IsNot.ValueString()
					}
					var54.EncryptionState = var55
					var54.Name = var53.Name.ValueString()
					var52 = append(var52, var54)
				}
			}
			var51.EncryptedLocations = var52
			var51.IsInstalled = plan.DiskEncryption.Criteria.IsInstalled.ValueBool()
		}
		var50.Criteria = var51
		var50.ExcludeVendor = plan.DiskEncryption.ExcludeVendor.ValueBool()
		var var56 []dJpBrWV.VendorObject
		if len(plan.DiskEncryption.Vendor) != 0 {
			var56 = make([]dJpBrWV.VendorObject, 0, len(plan.DiskEncryption.Vendor))
			for var57Index := range plan.DiskEncryption.Vendor {
				var57 := plan.DiskEncryption.Vendor[var57Index]
				var var58 dJpBrWV.VendorObject
				var58.Name = var57.Name.ValueString()
				var58.Product = DecodeStringSlice(var57.Product)
				var56 = append(var56, var58)
			}
		}
		var50.Vendor = var56
	}
	var0.DiskEncryption = var50
	var var59 *dJpBrWV.FirewallObject
	if plan.Firewall != nil {
		var59 = &dJpBrWV.FirewallObject{}
		var var60 *dJpBrWV.CriteriaObject3
		if plan.Firewall.Criteria != nil {
			var60 = &dJpBrWV.CriteriaObject3{}
			var60.IsEnabled = plan.Firewall.Criteria.IsEnabled.ValueString()
			var60.IsInstalled = plan.Firewall.Criteria.IsInstalled.ValueBool()
		}
		var59.Criteria = var60
		var59.ExcludeVendor = plan.Firewall.ExcludeVendor.ValueBool()
		var var61 []dJpBrWV.VendorObject
		if len(plan.Firewall.Vendor) != 0 {
			var61 = make([]dJpBrWV.VendorObject, 0, len(plan.Firewall.Vendor))
			for var62Index := range plan.Firewall.Vendor {
				var62 := plan.Firewall.Vendor[var62Index]
				var var63 dJpBrWV.VendorObject
				var63.Name = var62.Name.ValueString()
				var63.Product = DecodeStringSlice(var62.Product)
				var61 = append(var61, var63)
			}
		}
		var59.Vendor = var61
	}
	var0.Firewall = var59
	var var64 *dJpBrWV.HostInfoObject
	if plan.HostInfo != nil {
		var64 = &dJpBrWV.HostInfoObject{}
		var var65 dJpBrWV.CriteriaObject6
		var var66 *dJpBrWV.ClientVersionObject
		if plan.HostInfo.Criteria.ClientVersion != nil {
			var66 = &dJpBrWV.ClientVersionObject{}
			var66.Contains = plan.HostInfo.Criteria.ClientVersion.Contains.ValueString()
			var66.Is = plan.HostInfo.Criteria.ClientVersion.Is.ValueString()
			var66.IsNot = plan.HostInfo.Criteria.ClientVersion.IsNot.ValueString()
		}
		var65.ClientVersion = var66
		var var67 *dJpBrWV.DomainObject
		if plan.HostInfo.Criteria.Domain != nil {
			var67 = &dJpBrWV.DomainObject{}
			var67.Contains = plan.HostInfo.Criteria.Domain.Contains.ValueString()
			var67.Is = plan.HostInfo.Criteria.Domain.Is.ValueString()
			var67.IsNot = plan.HostInfo.Criteria.Domain.IsNot.ValueString()
		}
		var65.Domain = var67
		var var68 *dJpBrWV.HostIdObject
		if plan.HostInfo.Criteria.HostId != nil {
			var68 = &dJpBrWV.HostIdObject{}
			var68.Contains = plan.HostInfo.Criteria.HostId.Contains.ValueString()
			var68.Is = plan.HostInfo.Criteria.HostId.Is.ValueString()
			var68.IsNot = plan.HostInfo.Criteria.HostId.IsNot.ValueString()
		}
		var65.HostId = var68
		var var69 *dJpBrWV.HostNameObject
		if plan.HostInfo.Criteria.HostName != nil {
			var69 = &dJpBrWV.HostNameObject{}
			var69.Contains = plan.HostInfo.Criteria.HostName.Contains.ValueString()
			var69.Is = plan.HostInfo.Criteria.HostName.Is.ValueString()
			var69.IsNot = plan.HostInfo.Criteria.HostName.IsNot.ValueString()
		}
		var65.HostName = var69
		var65.Managed = plan.HostInfo.Criteria.Managed.ValueBool()
		var var70 *dJpBrWV.OsObject
		if plan.HostInfo.Criteria.Os != nil {
			var70 = &dJpBrWV.OsObject{}
			var var71 *dJpBrWV.ContainsObject
			if plan.HostInfo.Criteria.Os.Contains != nil {
				var71 = &dJpBrWV.ContainsObject{}
				var71.Apple = plan.HostInfo.Criteria.Os.Contains.Apple.ValueString()
				var71.Google = plan.HostInfo.Criteria.Os.Contains.Google.ValueString()
				var71.Linux = plan.HostInfo.Criteria.Os.Contains.Linux.ValueString()
				var71.Microsoft = plan.HostInfo.Criteria.Os.Contains.Microsoft.ValueString()
				var71.Other = plan.HostInfo.Criteria.Os.Contains.Other.ValueString()
			}
			var70.Contains = var71
		}
		var65.Os = var70
		var var72 *dJpBrWV.SerialNumberObject
		if plan.HostInfo.Criteria.SerialNumber != nil {
			var72 = &dJpBrWV.SerialNumberObject{}
			var72.Contains = plan.HostInfo.Criteria.SerialNumber.Contains.ValueString()
			var72.Is = plan.HostInfo.Criteria.SerialNumber.Is.ValueString()
			var72.IsNot = plan.HostInfo.Criteria.SerialNumber.IsNot.ValueString()
		}
		var65.SerialNumber = var72
		var64.Criteria = var65
	}
	var0.HostInfo = var64
	var var73 *dJpBrWV.MobileDeviceObject
	if plan.MobileDevice != nil {
		var73 = &dJpBrWV.MobileDeviceObject{}
		var var74 *dJpBrWV.CriteriaObject7
		if plan.MobileDevice.Criteria != nil {
			var74 = &dJpBrWV.CriteriaObject7{}
			var var75 *dJpBrWV.ApplicationsObject
			if plan.MobileDevice.Criteria.Applications != nil {
				var75 = &dJpBrWV.ApplicationsObject{}
				var var76 *dJpBrWV.HasMalwareObject
				if plan.MobileDevice.Criteria.Applications.HasMalware != nil {
					var76 = &dJpBrWV.HasMalwareObject{}
					if plan.MobileDevice.Criteria.Applications.HasMalware.No.ValueBool() {
						var76.No = struct{}{}
					}
					var var77 *dJpBrWV.YesObject
					if plan.MobileDevice.Criteria.Applications.HasMalware.Yes != nil {
						var77 = &dJpBrWV.YesObject{}
						var var78 []dJpBrWV.ExcludesObject
						if len(plan.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes) != 0 {
							var78 = make([]dJpBrWV.ExcludesObject, 0, len(plan.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes))
							for var79Index := range plan.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes {
								var79 := plan.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes[var79Index]
								var var80 dJpBrWV.ExcludesObject
								var80.Hash = var79.Hash.ValueString()
								var80.Name = var79.Name.ValueString()
								var80.Package = var79.Package.ValueString()
								var78 = append(var78, var80)
							}
						}
						var77.Excludes = var78
					}
					var76.Yes = var77
				}
				var75.HasMalware = var76
				var75.HasUnmanagedApp = plan.MobileDevice.Criteria.Applications.HasUnmanagedApp.ValueBool()
				var var81 []dJpBrWV.IncludesObject
				if len(plan.MobileDevice.Criteria.Applications.Includes) != 0 {
					var81 = make([]dJpBrWV.IncludesObject, 0, len(plan.MobileDevice.Criteria.Applications.Includes))
					for var82Index := range plan.MobileDevice.Criteria.Applications.Includes {
						var82 := plan.MobileDevice.Criteria.Applications.Includes[var82Index]
						var var83 dJpBrWV.IncludesObject
						var83.Hash = var82.Hash.ValueString()
						var83.Name = var82.Name.ValueString()
						var83.Package = var82.Package.ValueString()
						var81 = append(var81, var83)
					}
				}
				var75.Includes = var81
			}
			var74.Applications = var75
			var74.DiskEncrypted = plan.MobileDevice.Criteria.DiskEncrypted.ValueBool()
			var var84 *dJpBrWV.ImeiObject
			if plan.MobileDevice.Criteria.Imei != nil {
				var84 = &dJpBrWV.ImeiObject{}
				var84.Contains = plan.MobileDevice.Criteria.Imei.Contains.ValueString()
				var84.Is = plan.MobileDevice.Criteria.Imei.Is.ValueString()
				var84.IsNot = plan.MobileDevice.Criteria.Imei.IsNot.ValueString()
			}
			var74.Imei = var84
			var74.Jailbroken = plan.MobileDevice.Criteria.Jailbroken.ValueBool()
			var var85 *dJpBrWV.LastCheckinTimeObject
			if plan.MobileDevice.Criteria.LastCheckinTime != nil {
				var85 = &dJpBrWV.LastCheckinTimeObject{}
				var var86 *dJpBrWV.NotWithinObject3
				if plan.MobileDevice.Criteria.LastCheckinTime.NotWithin != nil {
					var86 = &dJpBrWV.NotWithinObject3{}
					var86.Days = plan.MobileDevice.Criteria.LastCheckinTime.NotWithin.Days.ValueInt64()
				}
				var85.NotWithin = var86
				var var87 *dJpBrWV.WithinObject3
				if plan.MobileDevice.Criteria.LastCheckinTime.Within != nil {
					var87 = &dJpBrWV.WithinObject3{}
					var87.Days = plan.MobileDevice.Criteria.LastCheckinTime.Within.Days.ValueInt64()
				}
				var85.Within = var87
			}
			var74.LastCheckinTime = var85
			var var88 *dJpBrWV.ModelObject
			if plan.MobileDevice.Criteria.Model != nil {
				var88 = &dJpBrWV.ModelObject{}
				var88.Contains = plan.MobileDevice.Criteria.Model.Contains.ValueString()
				var88.Is = plan.MobileDevice.Criteria.Model.Is.ValueString()
				var88.IsNot = plan.MobileDevice.Criteria.Model.IsNot.ValueString()
			}
			var74.Model = var88
			var74.PasscodeSet = plan.MobileDevice.Criteria.PasscodeSet.ValueBool()
			var var89 *dJpBrWV.PhoneNumberObject
			if plan.MobileDevice.Criteria.PhoneNumber != nil {
				var89 = &dJpBrWV.PhoneNumberObject{}
				var89.Contains = plan.MobileDevice.Criteria.PhoneNumber.Contains.ValueString()
				var89.Is = plan.MobileDevice.Criteria.PhoneNumber.Is.ValueString()
				var89.IsNot = plan.MobileDevice.Criteria.PhoneNumber.IsNot.ValueString()
			}
			var74.PhoneNumber = var89
			var var90 *dJpBrWV.TagObject
			if plan.MobileDevice.Criteria.Tag != nil {
				var90 = &dJpBrWV.TagObject{}
				var90.Contains = plan.MobileDevice.Criteria.Tag.Contains.ValueString()
				var90.Is = plan.MobileDevice.Criteria.Tag.Is.ValueString()
				var90.IsNot = plan.MobileDevice.Criteria.Tag.IsNot.ValueString()
			}
			var74.Tag = var90
		}
		var73.Criteria = var74
	}
	var0.MobileDevice = var73
	var0.Name = plan.Name.ValueString()
	var var91 *dJpBrWV.NetworkInfoObject
	if plan.NetworkInfo != nil {
		var91 = &dJpBrWV.NetworkInfoObject{}
		var var92 *dJpBrWV.CriteriaObject8
		if plan.NetworkInfo.Criteria != nil {
			var92 = &dJpBrWV.CriteriaObject8{}
			var var93 *dJpBrWV.NetworkObject
			if plan.NetworkInfo.Criteria.Network != nil {
				var93 = &dJpBrWV.NetworkObject{}
				var var94 *dJpBrWV.IsObject
				if plan.NetworkInfo.Criteria.Network.Is != nil {
					var94 = &dJpBrWV.IsObject{}
					var var95 *dJpBrWV.MobileObject
					if plan.NetworkInfo.Criteria.Network.Is.Mobile != nil {
						var95 = &dJpBrWV.MobileObject{}
						var95.Carrier = plan.NetworkInfo.Criteria.Network.Is.Mobile.Carrier.ValueString()
					}
					var94.Mobile = var95
					if plan.NetworkInfo.Criteria.Network.Is.Unknown.ValueBool() {
						var94.Unknown = struct{}{}
					}
					var var96 *dJpBrWV.WifiObject
					if plan.NetworkInfo.Criteria.Network.Is.Wifi != nil {
						var96 = &dJpBrWV.WifiObject{}
						var96.Ssid = plan.NetworkInfo.Criteria.Network.Is.Wifi.Ssid.ValueString()
					}
					var94.Wifi = var96
				}
				var93.Is = var94
				var var97 *dJpBrWV.IsNotObject
				if plan.NetworkInfo.Criteria.Network.IsNot != nil {
					var97 = &dJpBrWV.IsNotObject{}
					if plan.NetworkInfo.Criteria.Network.IsNot.Ethernet.ValueBool() {
						var97.Ethernet = struct{}{}
					}
					var var98 *dJpBrWV.MobileObject
					if plan.NetworkInfo.Criteria.Network.IsNot.Mobile != nil {
						var98 = &dJpBrWV.MobileObject{}
						var98.Carrier = plan.NetworkInfo.Criteria.Network.IsNot.Mobile.Carrier.ValueString()
					}
					var97.Mobile = var98
					if plan.NetworkInfo.Criteria.Network.IsNot.Unknown.ValueBool() {
						var97.Unknown = struct{}{}
					}
					var var99 *dJpBrWV.WifiObject
					if plan.NetworkInfo.Criteria.Network.IsNot.Wifi != nil {
						var99 = &dJpBrWV.WifiObject{}
						var99.Ssid = plan.NetworkInfo.Criteria.Network.IsNot.Wifi.Ssid.ValueString()
					}
					var97.Wifi = var99
				}
				var93.IsNot = var97
			}
			var92.Network = var93
		}
		var91.Criteria = var92
	}
	var0.NetworkInfo = var91
	var var100 *dJpBrWV.PatchManagementObject
	if plan.PatchManagement != nil {
		var100 = &dJpBrWV.PatchManagementObject{}
		var var101 *dJpBrWV.CriteriaObject9
		if plan.PatchManagement.Criteria != nil {
			var101 = &dJpBrWV.CriteriaObject9{}
			var101.IsEnabled = plan.PatchManagement.Criteria.IsEnabled.ValueString()
			var101.IsInstalled = plan.PatchManagement.Criteria.IsInstalled.ValueBool()
			var var102 *dJpBrWV.MissingPatchesObject
			if plan.PatchManagement.Criteria.MissingPatches != nil {
				var102 = &dJpBrWV.MissingPatchesObject{}
				var102.Check = plan.PatchManagement.Criteria.MissingPatches.Check.ValueString()
				var102.Patches = DecodeStringSlice(plan.PatchManagement.Criteria.MissingPatches.Patches)
				var var103 *dJpBrWV.SeverityObject
				if plan.PatchManagement.Criteria.MissingPatches.Severity != nil {
					var103 = &dJpBrWV.SeverityObject{}
					var103.GreaterEqual = plan.PatchManagement.Criteria.MissingPatches.Severity.GreaterEqual.ValueInt64()
					var103.GreaterThan = plan.PatchManagement.Criteria.MissingPatches.Severity.GreaterThan.ValueInt64()
					var103.Is = plan.PatchManagement.Criteria.MissingPatches.Severity.Is.ValueInt64()
					var103.IsNot = plan.PatchManagement.Criteria.MissingPatches.Severity.IsNot.ValueInt64()
					var103.LessEqual = plan.PatchManagement.Criteria.MissingPatches.Severity.LessEqual.ValueInt64()
					var103.LessThan = plan.PatchManagement.Criteria.MissingPatches.Severity.LessThan.ValueInt64()
				}
				var102.Severity = var103
			}
			var101.MissingPatches = var102
		}
		var100.Criteria = var101
		var100.ExcludeVendor = plan.PatchManagement.ExcludeVendor.ValueBool()
		var var104 []dJpBrWV.VendorObject1
		if len(plan.PatchManagement.Vendor) != 0 {
			var104 = make([]dJpBrWV.VendorObject1, 0, len(plan.PatchManagement.Vendor))
			for var105Index := range plan.PatchManagement.Vendor {
				var105 := plan.PatchManagement.Vendor[var105Index]
				var var106 dJpBrWV.VendorObject1
				var106.Name = var105.Name.ValueString()
				var106.Product = DecodeStringSlice(var105.Product)
				var104 = append(var104, var106)
			}
		}
		var100.Vendor = var104
	}
	var0.PatchManagement = var100
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var107 *objectsHipObjectsRsModelAntiMalwareObject
	if ans.AntiMalware != nil {
		var107 = &objectsHipObjectsRsModelAntiMalwareObject{}
		var var108 *objectsHipObjectsRsModelCriteriaObject
		if ans.AntiMalware.Criteria != nil {
			var108 = &objectsHipObjectsRsModelCriteriaObject{}
			var var109 *objectsHipObjectsRsModelLastScanTimeObject
			if ans.AntiMalware.Criteria.LastScanTime != nil {
				var109 = &objectsHipObjectsRsModelLastScanTimeObject{}
				var var110 *objectsHipObjectsRsModelNotWithinObject
				if ans.AntiMalware.Criteria.LastScanTime.NotWithin != nil {
					var110 = &objectsHipObjectsRsModelNotWithinObject{}
					var110.Days = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.NotWithin.Days)
					var110.Hours = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.NotWithin.Hours)
				}
				var var111 *objectsHipObjectsRsModelWithinObject
				if ans.AntiMalware.Criteria.LastScanTime.Within != nil {
					var111 = &objectsHipObjectsRsModelWithinObject{}
					var111.Days = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.Within.Days)
					var111.Hours = types.Int64Value(ans.AntiMalware.Criteria.LastScanTime.Within.Hours)
				}
				if ans.AntiMalware.Criteria.LastScanTime.NotAvailable != nil {
					var109.NotAvailable = types.BoolValue(true)
				}
				var109.NotWithin = var110
				var109.Within = var111
			}
			var var112 *objectsHipObjectsRsModelProductVersionObject
			if ans.AntiMalware.Criteria.ProductVersion != nil {
				var112 = &objectsHipObjectsRsModelProductVersionObject{}
				var var113 *objectsHipObjectsRsModelNotWithinObject1
				if ans.AntiMalware.Criteria.ProductVersion.NotWithin != nil {
					var113 = &objectsHipObjectsRsModelNotWithinObject1{}
					var113.Versions = types.Int64Value(ans.AntiMalware.Criteria.ProductVersion.NotWithin.Versions)
				}
				var var114 *objectsHipObjectsRsModelWithinObject1
				if ans.AntiMalware.Criteria.ProductVersion.Within != nil {
					var114 = &objectsHipObjectsRsModelWithinObject1{}
					var114.Versions = types.Int64Value(ans.AntiMalware.Criteria.ProductVersion.Within.Versions)
				}
				var112.Contains = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.Contains)
				var112.GreaterEqual = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.GreaterEqual)
				var112.GreaterThan = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.GreaterThan)
				var112.Is = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.Is)
				var112.IsNot = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.IsNot)
				var112.LessEqual = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.LessEqual)
				var112.LessThan = types.StringValue(ans.AntiMalware.Criteria.ProductVersion.LessThan)
				var112.NotWithin = var113
				var112.Within = var114
			}
			var var115 *objectsHipObjectsRsModelVirdefVersionObject
			if ans.AntiMalware.Criteria.VirdefVersion != nil {
				var115 = &objectsHipObjectsRsModelVirdefVersionObject{}
				var var116 *objectsHipObjectsRsModelNotWithinObject2
				if ans.AntiMalware.Criteria.VirdefVersion.NotWithin != nil {
					var116 = &objectsHipObjectsRsModelNotWithinObject2{}
					var116.Days = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.NotWithin.Days)
					var116.Versions = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.NotWithin.Versions)
				}
				var var117 *objectsHipObjectsRsModelWithinObject2
				if ans.AntiMalware.Criteria.VirdefVersion.Within != nil {
					var117 = &objectsHipObjectsRsModelWithinObject2{}
					var117.Days = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.Within.Days)
					var117.Versions = types.Int64Value(ans.AntiMalware.Criteria.VirdefVersion.Within.Versions)
				}
				var115.NotWithin = var116
				var115.Within = var117
			}
			var108.IsInstalled = types.BoolValue(ans.AntiMalware.Criteria.IsInstalled)
			var108.LastScanTime = var109
			var108.ProductVersion = var112
			var108.RealTimeProtection = types.StringValue(ans.AntiMalware.Criteria.RealTimeProtection)
			var108.VirdefVersion = var115
		}
		var var118 []objectsHipObjectsRsModelVendorObject
		if len(ans.AntiMalware.Vendor) != 0 {
			var118 = make([]objectsHipObjectsRsModelVendorObject, 0, len(ans.AntiMalware.Vendor))
			for var119Index := range ans.AntiMalware.Vendor {
				var119 := ans.AntiMalware.Vendor[var119Index]
				var var120 objectsHipObjectsRsModelVendorObject
				var120.Name = types.StringValue(var119.Name)
				var120.Product = EncodeStringSlice(var119.Product)
				var118 = append(var118, var120)
			}
		}
		var107.Criteria = var108
		var107.ExcludeVendor = types.BoolValue(ans.AntiMalware.ExcludeVendor)
		var107.Vendor = var118
	}
	var var121 *objectsHipObjectsRsModelCertificateObject
	if ans.Certificate != nil {
		var121 = &objectsHipObjectsRsModelCertificateObject{}
		var var122 *objectsHipObjectsRsModelCriteriaObject1
		if ans.Certificate.Criteria != nil {
			var122 = &objectsHipObjectsRsModelCriteriaObject1{}
			var var123 []objectsHipObjectsRsModelCertificateAttributesObject
			if len(ans.Certificate.Criteria.CertificateAttributes) != 0 {
				var123 = make([]objectsHipObjectsRsModelCertificateAttributesObject, 0, len(ans.Certificate.Criteria.CertificateAttributes))
				for var124Index := range ans.Certificate.Criteria.CertificateAttributes {
					var124 := ans.Certificate.Criteria.CertificateAttributes[var124Index]
					var var125 objectsHipObjectsRsModelCertificateAttributesObject
					var125.Name = types.StringValue(var124.Name)
					var125.Value = types.StringValue(var124.Value)
					var123 = append(var123, var125)
				}
			}
			var122.CertificateAttributes = var123
			var122.CertificateProfile = types.StringValue(ans.Certificate.Criteria.CertificateProfile)
		}
		var121.Criteria = var122
	}
	var var126 *objectsHipObjectsRsModelCustomChecksObject
	if ans.CustomChecks != nil {
		var126 = &objectsHipObjectsRsModelCustomChecksObject{}
		var var127 objectsHipObjectsRsModelCriteriaObject2
		var var128 []objectsHipObjectsRsModelPlistObject
		if len(ans.CustomChecks.Criteria.Plist) != 0 {
			var128 = make([]objectsHipObjectsRsModelPlistObject, 0, len(ans.CustomChecks.Criteria.Plist))
			for var129Index := range ans.CustomChecks.Criteria.Plist {
				var129 := ans.CustomChecks.Criteria.Plist[var129Index]
				var var130 objectsHipObjectsRsModelPlistObject
				var var131 []objectsHipObjectsRsModelKeyObject
				if len(var129.Key) != 0 {
					var131 = make([]objectsHipObjectsRsModelKeyObject, 0, len(var129.Key))
					for var132Index := range var129.Key {
						var132 := var129.Key[var132Index]
						var var133 objectsHipObjectsRsModelKeyObject
						var133.Name = types.StringValue(var132.Name)
						var133.Negate = types.BoolValue(var132.Negate)
						var133.Value = types.StringValue(var132.Value)
						var131 = append(var131, var133)
					}
				}
				var130.Key = var131
				var130.Name = types.StringValue(var129.Name)
				var130.Negate = types.BoolValue(var129.Negate)
				var128 = append(var128, var130)
			}
		}
		var var134 []objectsHipObjectsRsModelProcessListObject
		if len(ans.CustomChecks.Criteria.ProcessList) != 0 {
			var134 = make([]objectsHipObjectsRsModelProcessListObject, 0, len(ans.CustomChecks.Criteria.ProcessList))
			for var135Index := range ans.CustomChecks.Criteria.ProcessList {
				var135 := ans.CustomChecks.Criteria.ProcessList[var135Index]
				var var136 objectsHipObjectsRsModelProcessListObject
				var136.Name = types.StringValue(var135.Name)
				var136.Running = types.BoolValue(var135.Running)
				var134 = append(var134, var136)
			}
		}
		var var137 []objectsHipObjectsRsModelRegistryKeyObject
		if len(ans.CustomChecks.Criteria.RegistryKey) != 0 {
			var137 = make([]objectsHipObjectsRsModelRegistryKeyObject, 0, len(ans.CustomChecks.Criteria.RegistryKey))
			for var138Index := range ans.CustomChecks.Criteria.RegistryKey {
				var138 := ans.CustomChecks.Criteria.RegistryKey[var138Index]
				var var139 objectsHipObjectsRsModelRegistryKeyObject
				var var140 []objectsHipObjectsRsModelRegistryValueObject
				if len(var138.RegistryValue) != 0 {
					var140 = make([]objectsHipObjectsRsModelRegistryValueObject, 0, len(var138.RegistryValue))
					for var141Index := range var138.RegistryValue {
						var141 := var138.RegistryValue[var141Index]
						var var142 objectsHipObjectsRsModelRegistryValueObject
						var142.Name = types.StringValue(var141.Name)
						var142.Negate = types.BoolValue(var141.Negate)
						var142.ValueData = types.StringValue(var141.ValueData)
						var140 = append(var140, var142)
					}
				}
				var139.DefaultValueData = types.StringValue(var138.DefaultValueData)
				var139.Name = types.StringValue(var138.Name)
				var139.Negate = types.BoolValue(var138.Negate)
				var139.RegistryValue = var140
				var137 = append(var137, var139)
			}
		}
		var127.Plist = var128
		var127.ProcessList = var134
		var127.RegistryKey = var137
		var126.Criteria = var127
	}
	var var143 *objectsHipObjectsRsModelDataLossPreventionObject
	if ans.DataLossPrevention != nil {
		var143 = &objectsHipObjectsRsModelDataLossPreventionObject{}
		var var144 *objectsHipObjectsRsModelCriteriaObject3
		if ans.DataLossPrevention.Criteria != nil {
			var144 = &objectsHipObjectsRsModelCriteriaObject3{}
			var144.IsEnabled = types.StringValue(ans.DataLossPrevention.Criteria.IsEnabled)
			var144.IsInstalled = types.BoolValue(ans.DataLossPrevention.Criteria.IsInstalled)
		}
		var var145 []objectsHipObjectsRsModelVendorObject1
		if len(ans.DataLossPrevention.Vendor) != 0 {
			var145 = make([]objectsHipObjectsRsModelVendorObject1, 0, len(ans.DataLossPrevention.Vendor))
			for var146Index := range ans.DataLossPrevention.Vendor {
				var146 := ans.DataLossPrevention.Vendor[var146Index]
				var var147 objectsHipObjectsRsModelVendorObject1
				var147.Name = types.StringValue(var146.Name)
				var147.Product = EncodeStringSlice(var146.Product)
				var145 = append(var145, var147)
			}
		}
		var143.Criteria = var144
		var143.ExcludeVendor = types.BoolValue(ans.DataLossPrevention.ExcludeVendor)
		var143.Vendor = var145
	}
	var var148 *objectsHipObjectsRsModelDiskBackupObject
	if ans.DiskBackup != nil {
		var148 = &objectsHipObjectsRsModelDiskBackupObject{}
		var var149 *objectsHipObjectsRsModelCriteriaObject4
		if ans.DiskBackup.Criteria != nil {
			var149 = &objectsHipObjectsRsModelCriteriaObject4{}
			var var150 *objectsHipObjectsRsModelLastBackupTimeObject
			if ans.DiskBackup.Criteria.LastBackupTime != nil {
				var150 = &objectsHipObjectsRsModelLastBackupTimeObject{}
				var var151 *objectsHipObjectsRsModelNotWithinObject
				if ans.DiskBackup.Criteria.LastBackupTime.NotWithin != nil {
					var151 = &objectsHipObjectsRsModelNotWithinObject{}
					var151.Days = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.NotWithin.Days)
					var151.Hours = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.NotWithin.Hours)
				}
				var var152 *objectsHipObjectsRsModelWithinObject
				if ans.DiskBackup.Criteria.LastBackupTime.Within != nil {
					var152 = &objectsHipObjectsRsModelWithinObject{}
					var152.Days = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.Within.Days)
					var152.Hours = types.Int64Value(ans.DiskBackup.Criteria.LastBackupTime.Within.Hours)
				}
				if ans.DiskBackup.Criteria.LastBackupTime.NotAvailable != nil {
					var150.NotAvailable = types.BoolValue(true)
				}
				var150.NotWithin = var151
				var150.Within = var152
			}
			var149.IsInstalled = types.BoolValue(ans.DiskBackup.Criteria.IsInstalled)
			var149.LastBackupTime = var150
		}
		var var153 []objectsHipObjectsRsModelVendorObject
		if len(ans.DiskBackup.Vendor) != 0 {
			var153 = make([]objectsHipObjectsRsModelVendorObject, 0, len(ans.DiskBackup.Vendor))
			for var154Index := range ans.DiskBackup.Vendor {
				var154 := ans.DiskBackup.Vendor[var154Index]
				var var155 objectsHipObjectsRsModelVendorObject
				var155.Name = types.StringValue(var154.Name)
				var155.Product = EncodeStringSlice(var154.Product)
				var153 = append(var153, var155)
			}
		}
		var148.Criteria = var149
		var148.ExcludeVendor = types.BoolValue(ans.DiskBackup.ExcludeVendor)
		var148.Vendor = var153
	}
	var var156 *objectsHipObjectsRsModelDiskEncryptionObject
	if ans.DiskEncryption != nil {
		var156 = &objectsHipObjectsRsModelDiskEncryptionObject{}
		var var157 *objectsHipObjectsRsModelCriteriaObject5
		if ans.DiskEncryption.Criteria != nil {
			var157 = &objectsHipObjectsRsModelCriteriaObject5{}
			var var158 []objectsHipObjectsRsModelEncryptedLocationsObject
			if len(ans.DiskEncryption.Criteria.EncryptedLocations) != 0 {
				var158 = make([]objectsHipObjectsRsModelEncryptedLocationsObject, 0, len(ans.DiskEncryption.Criteria.EncryptedLocations))
				for var159Index := range ans.DiskEncryption.Criteria.EncryptedLocations {
					var159 := ans.DiskEncryption.Criteria.EncryptedLocations[var159Index]
					var var160 objectsHipObjectsRsModelEncryptedLocationsObject
					var var161 *objectsHipObjectsRsModelEncryptionStateObject
					if var159.EncryptionState != nil {
						var161 = &objectsHipObjectsRsModelEncryptionStateObject{}
						var161.Is = types.StringValue(var159.EncryptionState.Is)
						var161.IsNot = types.StringValue(var159.EncryptionState.IsNot)
					}
					var160.EncryptionState = var161
					var160.Name = types.StringValue(var159.Name)
					var158 = append(var158, var160)
				}
			}
			var157.EncryptedLocations = var158
			var157.IsInstalled = types.BoolValue(ans.DiskEncryption.Criteria.IsInstalled)
		}
		var var162 []objectsHipObjectsRsModelVendorObject
		if len(ans.DiskEncryption.Vendor) != 0 {
			var162 = make([]objectsHipObjectsRsModelVendorObject, 0, len(ans.DiskEncryption.Vendor))
			for var163Index := range ans.DiskEncryption.Vendor {
				var163 := ans.DiskEncryption.Vendor[var163Index]
				var var164 objectsHipObjectsRsModelVendorObject
				var164.Name = types.StringValue(var163.Name)
				var164.Product = EncodeStringSlice(var163.Product)
				var162 = append(var162, var164)
			}
		}
		var156.Criteria = var157
		var156.ExcludeVendor = types.BoolValue(ans.DiskEncryption.ExcludeVendor)
		var156.Vendor = var162
	}
	var var165 *objectsHipObjectsRsModelFirewallObject
	if ans.Firewall != nil {
		var165 = &objectsHipObjectsRsModelFirewallObject{}
		var var166 *objectsHipObjectsRsModelCriteriaObject3
		if ans.Firewall.Criteria != nil {
			var166 = &objectsHipObjectsRsModelCriteriaObject3{}
			var166.IsEnabled = types.StringValue(ans.Firewall.Criteria.IsEnabled)
			var166.IsInstalled = types.BoolValue(ans.Firewall.Criteria.IsInstalled)
		}
		var var167 []objectsHipObjectsRsModelVendorObject
		if len(ans.Firewall.Vendor) != 0 {
			var167 = make([]objectsHipObjectsRsModelVendorObject, 0, len(ans.Firewall.Vendor))
			for var168Index := range ans.Firewall.Vendor {
				var168 := ans.Firewall.Vendor[var168Index]
				var var169 objectsHipObjectsRsModelVendorObject
				var169.Name = types.StringValue(var168.Name)
				var169.Product = EncodeStringSlice(var168.Product)
				var167 = append(var167, var169)
			}
		}
		var165.Criteria = var166
		var165.ExcludeVendor = types.BoolValue(ans.Firewall.ExcludeVendor)
		var165.Vendor = var167
	}
	var var170 *objectsHipObjectsRsModelHostInfoObject
	if ans.HostInfo != nil {
		var170 = &objectsHipObjectsRsModelHostInfoObject{}
		var var171 objectsHipObjectsRsModelCriteriaObject6
		var var172 *objectsHipObjectsRsModelClientVersionObject
		if ans.HostInfo.Criteria.ClientVersion != nil {
			var172 = &objectsHipObjectsRsModelClientVersionObject{}
			var172.Contains = types.StringValue(ans.HostInfo.Criteria.ClientVersion.Contains)
			var172.Is = types.StringValue(ans.HostInfo.Criteria.ClientVersion.Is)
			var172.IsNot = types.StringValue(ans.HostInfo.Criteria.ClientVersion.IsNot)
		}
		var var173 *objectsHipObjectsRsModelDomainObject
		if ans.HostInfo.Criteria.Domain != nil {
			var173 = &objectsHipObjectsRsModelDomainObject{}
			var173.Contains = types.StringValue(ans.HostInfo.Criteria.Domain.Contains)
			var173.Is = types.StringValue(ans.HostInfo.Criteria.Domain.Is)
			var173.IsNot = types.StringValue(ans.HostInfo.Criteria.Domain.IsNot)
		}
		var var174 *objectsHipObjectsRsModelHostIdObject
		if ans.HostInfo.Criteria.HostId != nil {
			var174 = &objectsHipObjectsRsModelHostIdObject{}
			var174.Contains = types.StringValue(ans.HostInfo.Criteria.HostId.Contains)
			var174.Is = types.StringValue(ans.HostInfo.Criteria.HostId.Is)
			var174.IsNot = types.StringValue(ans.HostInfo.Criteria.HostId.IsNot)
		}
		var var175 *objectsHipObjectsRsModelHostNameObject
		if ans.HostInfo.Criteria.HostName != nil {
			var175 = &objectsHipObjectsRsModelHostNameObject{}
			var175.Contains = types.StringValue(ans.HostInfo.Criteria.HostName.Contains)
			var175.Is = types.StringValue(ans.HostInfo.Criteria.HostName.Is)
			var175.IsNot = types.StringValue(ans.HostInfo.Criteria.HostName.IsNot)
		}
		var var176 *objectsHipObjectsRsModelOsObject
		if ans.HostInfo.Criteria.Os != nil {
			var176 = &objectsHipObjectsRsModelOsObject{}
			var var177 *objectsHipObjectsRsModelContainsObject
			if ans.HostInfo.Criteria.Os.Contains != nil {
				var177 = &objectsHipObjectsRsModelContainsObject{}
				var177.Apple = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Apple)
				var177.Google = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Google)
				var177.Linux = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Linux)
				var177.Microsoft = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Microsoft)
				var177.Other = types.StringValue(ans.HostInfo.Criteria.Os.Contains.Other)
			}
			var176.Contains = var177
		}
		var var178 *objectsHipObjectsRsModelSerialNumberObject
		if ans.HostInfo.Criteria.SerialNumber != nil {
			var178 = &objectsHipObjectsRsModelSerialNumberObject{}
			var178.Contains = types.StringValue(ans.HostInfo.Criteria.SerialNumber.Contains)
			var178.Is = types.StringValue(ans.HostInfo.Criteria.SerialNumber.Is)
			var178.IsNot = types.StringValue(ans.HostInfo.Criteria.SerialNumber.IsNot)
		}
		var171.ClientVersion = var172
		var171.Domain = var173
		var171.HostId = var174
		var171.HostName = var175
		var171.Managed = types.BoolValue(ans.HostInfo.Criteria.Managed)
		var171.Os = var176
		var171.SerialNumber = var178
		var170.Criteria = var171
	}
	var var179 *objectsHipObjectsRsModelMobileDeviceObject
	if ans.MobileDevice != nil {
		var179 = &objectsHipObjectsRsModelMobileDeviceObject{}
		var var180 *objectsHipObjectsRsModelCriteriaObject7
		if ans.MobileDevice.Criteria != nil {
			var180 = &objectsHipObjectsRsModelCriteriaObject7{}
			var var181 *objectsHipObjectsRsModelApplicationsObject
			if ans.MobileDevice.Criteria.Applications != nil {
				var181 = &objectsHipObjectsRsModelApplicationsObject{}
				var var182 *objectsHipObjectsRsModelHasMalwareObject
				if ans.MobileDevice.Criteria.Applications.HasMalware != nil {
					var182 = &objectsHipObjectsRsModelHasMalwareObject{}
					var var183 *objectsHipObjectsRsModelYesObject
					if ans.MobileDevice.Criteria.Applications.HasMalware.Yes != nil {
						var183 = &objectsHipObjectsRsModelYesObject{}
						var var184 []objectsHipObjectsRsModelExcludesObject
						if len(ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes) != 0 {
							var184 = make([]objectsHipObjectsRsModelExcludesObject, 0, len(ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes))
							for var185Index := range ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes {
								var185 := ans.MobileDevice.Criteria.Applications.HasMalware.Yes.Excludes[var185Index]
								var var186 objectsHipObjectsRsModelExcludesObject
								var186.Hash = types.StringValue(var185.Hash)
								var186.Name = types.StringValue(var185.Name)
								var186.Package = types.StringValue(var185.Package)
								var184 = append(var184, var186)
							}
						}
						var183.Excludes = var184
					}
					if ans.MobileDevice.Criteria.Applications.HasMalware.No != nil {
						var182.No = types.BoolValue(true)
					}
					var182.Yes = var183
				}
				var var187 []objectsHipObjectsRsModelIncludesObject
				if len(ans.MobileDevice.Criteria.Applications.Includes) != 0 {
					var187 = make([]objectsHipObjectsRsModelIncludesObject, 0, len(ans.MobileDevice.Criteria.Applications.Includes))
					for var188Index := range ans.MobileDevice.Criteria.Applications.Includes {
						var188 := ans.MobileDevice.Criteria.Applications.Includes[var188Index]
						var var189 objectsHipObjectsRsModelIncludesObject
						var189.Hash = types.StringValue(var188.Hash)
						var189.Name = types.StringValue(var188.Name)
						var189.Package = types.StringValue(var188.Package)
						var187 = append(var187, var189)
					}
				}
				var181.HasMalware = var182
				var181.HasUnmanagedApp = types.BoolValue(ans.MobileDevice.Criteria.Applications.HasUnmanagedApp)
				var181.Includes = var187
			}
			var var190 *objectsHipObjectsRsModelImeiObject
			if ans.MobileDevice.Criteria.Imei != nil {
				var190 = &objectsHipObjectsRsModelImeiObject{}
				var190.Contains = types.StringValue(ans.MobileDevice.Criteria.Imei.Contains)
				var190.Is = types.StringValue(ans.MobileDevice.Criteria.Imei.Is)
				var190.IsNot = types.StringValue(ans.MobileDevice.Criteria.Imei.IsNot)
			}
			var var191 *objectsHipObjectsRsModelLastCheckinTimeObject
			if ans.MobileDevice.Criteria.LastCheckinTime != nil {
				var191 = &objectsHipObjectsRsModelLastCheckinTimeObject{}
				var var192 *objectsHipObjectsRsModelNotWithinObject3
				if ans.MobileDevice.Criteria.LastCheckinTime.NotWithin != nil {
					var192 = &objectsHipObjectsRsModelNotWithinObject3{}
					var192.Days = types.Int64Value(ans.MobileDevice.Criteria.LastCheckinTime.NotWithin.Days)
				}
				var var193 *objectsHipObjectsRsModelWithinObject3
				if ans.MobileDevice.Criteria.LastCheckinTime.Within != nil {
					var193 = &objectsHipObjectsRsModelWithinObject3{}
					var193.Days = types.Int64Value(ans.MobileDevice.Criteria.LastCheckinTime.Within.Days)
				}
				var191.NotWithin = var192
				var191.Within = var193
			}
			var var194 *objectsHipObjectsRsModelModelObject
			if ans.MobileDevice.Criteria.Model != nil {
				var194 = &objectsHipObjectsRsModelModelObject{}
				var194.Contains = types.StringValue(ans.MobileDevice.Criteria.Model.Contains)
				var194.Is = types.StringValue(ans.MobileDevice.Criteria.Model.Is)
				var194.IsNot = types.StringValue(ans.MobileDevice.Criteria.Model.IsNot)
			}
			var var195 *objectsHipObjectsRsModelPhoneNumberObject
			if ans.MobileDevice.Criteria.PhoneNumber != nil {
				var195 = &objectsHipObjectsRsModelPhoneNumberObject{}
				var195.Contains = types.StringValue(ans.MobileDevice.Criteria.PhoneNumber.Contains)
				var195.Is = types.StringValue(ans.MobileDevice.Criteria.PhoneNumber.Is)
				var195.IsNot = types.StringValue(ans.MobileDevice.Criteria.PhoneNumber.IsNot)
			}
			var var196 *objectsHipObjectsRsModelTagObject
			if ans.MobileDevice.Criteria.Tag != nil {
				var196 = &objectsHipObjectsRsModelTagObject{}
				var196.Contains = types.StringValue(ans.MobileDevice.Criteria.Tag.Contains)
				var196.Is = types.StringValue(ans.MobileDevice.Criteria.Tag.Is)
				var196.IsNot = types.StringValue(ans.MobileDevice.Criteria.Tag.IsNot)
			}
			var180.Applications = var181
			var180.DiskEncrypted = types.BoolValue(ans.MobileDevice.Criteria.DiskEncrypted)
			var180.Imei = var190
			var180.Jailbroken = types.BoolValue(ans.MobileDevice.Criteria.Jailbroken)
			var180.LastCheckinTime = var191
			var180.Model = var194
			var180.PasscodeSet = types.BoolValue(ans.MobileDevice.Criteria.PasscodeSet)
			var180.PhoneNumber = var195
			var180.Tag = var196
		}
		var179.Criteria = var180
	}
	var var197 *objectsHipObjectsRsModelNetworkInfoObject
	if ans.NetworkInfo != nil {
		var197 = &objectsHipObjectsRsModelNetworkInfoObject{}
		var var198 *objectsHipObjectsRsModelCriteriaObject8
		if ans.NetworkInfo.Criteria != nil {
			var198 = &objectsHipObjectsRsModelCriteriaObject8{}
			var var199 *objectsHipObjectsRsModelNetworkObject
			if ans.NetworkInfo.Criteria.Network != nil {
				var199 = &objectsHipObjectsRsModelNetworkObject{}
				var var200 *objectsHipObjectsRsModelIsObject
				if ans.NetworkInfo.Criteria.Network.Is != nil {
					var200 = &objectsHipObjectsRsModelIsObject{}
					var var201 *objectsHipObjectsRsModelMobileObject
					if ans.NetworkInfo.Criteria.Network.Is.Mobile != nil {
						var201 = &objectsHipObjectsRsModelMobileObject{}
						var201.Carrier = types.StringValue(ans.NetworkInfo.Criteria.Network.Is.Mobile.Carrier)
					}
					var var202 *objectsHipObjectsRsModelWifiObject
					if ans.NetworkInfo.Criteria.Network.Is.Wifi != nil {
						var202 = &objectsHipObjectsRsModelWifiObject{}
						var202.Ssid = types.StringValue(ans.NetworkInfo.Criteria.Network.Is.Wifi.Ssid)
					}
					var200.Mobile = var201
					if ans.NetworkInfo.Criteria.Network.Is.Unknown != nil {
						var200.Unknown = types.BoolValue(true)
					}
					var200.Wifi = var202
				}
				var var203 *objectsHipObjectsRsModelIsNotObject
				if ans.NetworkInfo.Criteria.Network.IsNot != nil {
					var203 = &objectsHipObjectsRsModelIsNotObject{}
					var var204 *objectsHipObjectsRsModelMobileObject
					if ans.NetworkInfo.Criteria.Network.IsNot.Mobile != nil {
						var204 = &objectsHipObjectsRsModelMobileObject{}
						var204.Carrier = types.StringValue(ans.NetworkInfo.Criteria.Network.IsNot.Mobile.Carrier)
					}
					var var205 *objectsHipObjectsRsModelWifiObject
					if ans.NetworkInfo.Criteria.Network.IsNot.Wifi != nil {
						var205 = &objectsHipObjectsRsModelWifiObject{}
						var205.Ssid = types.StringValue(ans.NetworkInfo.Criteria.Network.IsNot.Wifi.Ssid)
					}
					if ans.NetworkInfo.Criteria.Network.IsNot.Ethernet != nil {
						var203.Ethernet = types.BoolValue(true)
					}
					var203.Mobile = var204
					if ans.NetworkInfo.Criteria.Network.IsNot.Unknown != nil {
						var203.Unknown = types.BoolValue(true)
					}
					var203.Wifi = var205
				}
				var199.Is = var200
				var199.IsNot = var203
			}
			var198.Network = var199
		}
		var197.Criteria = var198
	}
	var var206 *objectsHipObjectsRsModelPatchManagementObject
	if ans.PatchManagement != nil {
		var206 = &objectsHipObjectsRsModelPatchManagementObject{}
		var var207 *objectsHipObjectsRsModelCriteriaObject9
		if ans.PatchManagement.Criteria != nil {
			var207 = &objectsHipObjectsRsModelCriteriaObject9{}
			var var208 *objectsHipObjectsRsModelMissingPatchesObject
			if ans.PatchManagement.Criteria.MissingPatches != nil {
				var208 = &objectsHipObjectsRsModelMissingPatchesObject{}
				var var209 *objectsHipObjectsRsModelSeverityObject
				if ans.PatchManagement.Criteria.MissingPatches.Severity != nil {
					var209 = &objectsHipObjectsRsModelSeverityObject{}
					var209.GreaterEqual = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.GreaterEqual)
					var209.GreaterThan = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.GreaterThan)
					var209.Is = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.Is)
					var209.IsNot = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.IsNot)
					var209.LessEqual = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.LessEqual)
					var209.LessThan = types.Int64Value(ans.PatchManagement.Criteria.MissingPatches.Severity.LessThan)
				}
				var208.Check = types.StringValue(ans.PatchManagement.Criteria.MissingPatches.Check)
				var208.Patches = EncodeStringSlice(ans.PatchManagement.Criteria.MissingPatches.Patches)
				var208.Severity = var209
			}
			var207.IsEnabled = types.StringValue(ans.PatchManagement.Criteria.IsEnabled)
			var207.IsInstalled = types.BoolValue(ans.PatchManagement.Criteria.IsInstalled)
			var207.MissingPatches = var208
		}
		var var210 []objectsHipObjectsRsModelVendorObject1
		if len(ans.PatchManagement.Vendor) != 0 {
			var210 = make([]objectsHipObjectsRsModelVendorObject1, 0, len(ans.PatchManagement.Vendor))
			for var211Index := range ans.PatchManagement.Vendor {
				var211 := ans.PatchManagement.Vendor[var211Index]
				var var212 objectsHipObjectsRsModelVendorObject1
				var212.Name = types.StringValue(var211.Name)
				var212.Product = EncodeStringSlice(var211.Product)
				var210 = append(var210, var212)
			}
		}
		var206.Criteria = var207
		var206.ExcludeVendor = types.BoolValue(ans.PatchManagement.ExcludeVendor)
		var206.Vendor = var210
	}
	state.AntiMalware = var107
	state.Certificate = var121
	state.CustomChecks = var126
	state.DataLossPrevention = var143
	state.Description = types.StringValue(ans.Description)
	state.DiskBackup = var148
	state.DiskEncryption = var156
	state.Firewall = var165
	state.HostInfo = var170
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.MobileDevice = var179
	state.Name = types.StringValue(ans.Name)
	state.NetworkInfo = var197
	state.PatchManagement = var206

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *objectsHipObjectsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var idType types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &idType)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := idType.ValueString()
	tokens := strings.Split(id, IdSeparator)
	if len(tokens) != 2 {
		resp.Diagnostics.AddError("Error in resource ID format", "Expected 2 tokens")
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource delete", map[string]any{
		"resource_name": "sase_objects_hip_objects",
		"locMap":        map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":        tokens,
	})

	svc := yCYVNEN.NewClient(r.client)
	input := yCYVNEN.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *objectsHipObjectsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
