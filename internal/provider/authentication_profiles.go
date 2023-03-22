package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/sase-go"
	"github.com/paloaltonetworks/sase-go/api"
	alljvhu "github.com/paloaltonetworks/sase-go/netsec/schema/authentication/profiles"
	cUCsSiw "github.com/paloaltonetworks/sase-go/netsec/service/v1/authenticationprofiles"

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
	_ datasource.DataSource              = &authenticationProfilesListDataSource{}
	_ datasource.DataSourceWithConfigure = &authenticationProfilesListDataSource{}
)

func NewAuthenticationProfilesListDataSource() datasource.DataSource {
	return &authenticationProfilesListDataSource{}
}

type authenticationProfilesListDataSource struct {
	client *sase.Client
}

type authenticationProfilesListDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Limit  types.Int64  `tfsdk:"limit"`
	Offset types.Int64  `tfsdk:"offset"`
	Folder types.String `tfsdk:"folder"`
	Name   types.String `tfsdk:"name"`

	// Output.
	Data []authenticationProfilesListDsModelConfig `tfsdk:"data"`
	// input omit: Limit
	// input omit: Offset
	Total types.Int64 `tfsdk:"total"`
}

type authenticationProfilesListDsModelConfig struct {
	AllowList        []types.String                                          `tfsdk:"allow_list"`
	ObjectId         types.String                                            `tfsdk:"object_id"`
	Lockout          *authenticationProfilesListDsModelLockoutObject         `tfsdk:"lockout"`
	Method           *authenticationProfilesListDsModelMethodObject          `tfsdk:"method"`
	MultiFactorAuth  *authenticationProfilesListDsModelMultiFactorAuthObject `tfsdk:"multi_factor_auth"`
	Name             types.String                                            `tfsdk:"name"`
	SingleSignOn     *authenticationProfilesListDsModelSingleSignOnObject    `tfsdk:"single_sign_on"`
	UserDomain       types.String                                            `tfsdk:"user_domain"`
	UsernameModifier types.String                                            `tfsdk:"username_modifier"`
}

type authenticationProfilesListDsModelLockoutObject struct {
	FailedAttempts types.Int64 `tfsdk:"failed_attempts"`
	LockoutTime    types.Int64 `tfsdk:"lockout_time"`
}

type authenticationProfilesListDsModelMethodObject struct {
	Kerberos      *authenticationProfilesListDsModelKerberosObject `tfsdk:"kerberos"`
	Ldap          *authenticationProfilesListDsModelLdapObject     `tfsdk:"ldap"`
	LocalDatabase types.Bool                                       `tfsdk:"local_database"`
	Radius        *authenticationProfilesListDsModelRadiusObject   `tfsdk:"radius"`
	SamlIdp       *authenticationProfilesListDsModelSamlIdpObject  `tfsdk:"saml_idp"`
	Tacplus       *authenticationProfilesListDsModelTacplusObject  `tfsdk:"tacplus"`
}

type authenticationProfilesListDsModelKerberosObject struct {
	Realm         types.String `tfsdk:"realm"`
	ServerProfile types.String `tfsdk:"server_profile"`
}

type authenticationProfilesListDsModelLdapObject struct {
	LoginAttribute types.String `tfsdk:"login_attribute"`
	PasswdExpDays  types.Int64  `tfsdk:"passwd_exp_days"`
	ServerProfile  types.String `tfsdk:"server_profile"`
}

type authenticationProfilesListDsModelRadiusObject struct {
	Checkgroup    types.Bool   `tfsdk:"checkgroup"`
	ServerProfile types.String `tfsdk:"server_profile"`
}

type authenticationProfilesListDsModelSamlIdpObject struct {
	AttributeNameUsergroup    types.String `tfsdk:"attribute_name_usergroup"`
	AttributeNameUsername     types.String `tfsdk:"attribute_name_username"`
	CertificateProfile        types.String `tfsdk:"certificate_profile"`
	EnableSingleLogout        types.Bool   `tfsdk:"enable_single_logout"`
	RequestSigningCertificate types.String `tfsdk:"request_signing_certificate"`
	ServerProfile             types.String `tfsdk:"server_profile"`
}

type authenticationProfilesListDsModelTacplusObject struct {
	Checkgroup    types.Bool   `tfsdk:"checkgroup"`
	ServerProfile types.String `tfsdk:"server_profile"`
}

type authenticationProfilesListDsModelMultiFactorAuthObject struct {
	Factors   []types.String `tfsdk:"factors"`
	MfaEnable types.Bool     `tfsdk:"mfa_enable"`
}

type authenticationProfilesListDsModelSingleSignOnObject struct {
	KerberosKeytab types.String `tfsdk:"kerberos_keytab"`
	Realm          types.String `tfsdk:"realm"`
}

// Metadata returns the data source type name.
func (d *authenticationProfilesListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_profiles_list"
}

// Schema defines the schema for this listing data source.
func (d *authenticationProfilesListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"folder": dsschema.StringAttribute{
				Description: "The folder of the entry",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},
			"name": dsschema.StringAttribute{
				Description: "The name of the entry",
				Optional:    true,
			},

			// Output.
			"data": dsschema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"allow_list": dsschema.ListAttribute{
							Description: "",
							Computed:    true,
							ElementType: types.StringType,
						},
						"object_id": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"lockout": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"failed_attempts": dsschema.Int64Attribute{
									Description: "",
									Computed:    true,
								},
								"lockout_time": dsschema.Int64Attribute{
									Description: "",
									Computed:    true,
								},
							},
						},
						"method": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"kerberos": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"realm": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"server_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"ldap": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"login_attribute": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"passwd_exp_days": dsschema.Int64Attribute{
											Description: "",
											Computed:    true,
										},
										"server_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"local_database": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
								"radius": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"checkgroup": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"server_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"saml_idp": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"attribute_name_usergroup": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"attribute_name_username": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"certificate_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"enable_single_logout": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"request_signing_certificate": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
										"server_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
								"tacplus": dsschema.SingleNestedAttribute{
									Description: "",
									Computed:    true,
									Attributes: map[string]dsschema.Attribute{
										"checkgroup": dsschema.BoolAttribute{
											Description: "",
											Computed:    true,
										},
										"server_profile": dsschema.StringAttribute{
											Description: "",
											Computed:    true,
										},
									},
								},
							},
						},
						"multi_factor_auth": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"factors": dsschema.ListAttribute{
									Description: "",
									Computed:    true,
									ElementType: types.StringType,
								},
								"mfa_enable": dsschema.BoolAttribute{
									Description: "",
									Computed:    true,
								},
							},
						},
						"name": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"single_sign_on": dsschema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]dsschema.Attribute{
								"kerberos_keytab": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
								"realm": dsschema.StringAttribute{
									Description: "",
									Computed:    true,
								},
							},
						},
						"user_domain": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"username_modifier": dsschema.StringAttribute{
							Description: "",
							Computed:    true,
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
func (d *authenticationProfilesListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *authenticationProfilesListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationProfilesListDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source listing", map[string]any{
		"data_source_name":            "sase_authentication_profiles_list",
		"terraform_provider_function": "Read",
		"limit":                       state.Limit.ValueInt64(),
		"has_limit":                   !state.Limit.IsNull(),
		"offset":                      state.Offset.ValueInt64(),
		"has_offset":                  !state.Offset.IsNull(),
		"folder":                      state.Folder.ValueString(),
		"name":                        state.Name.ValueString(),
		"has_name":                    !state.Name.IsNull(),
	})

	// Prepare to run the command.
	svc := cUCsSiw.NewClient(d.client)
	input := cUCsSiw.ListInput{
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
	var idBuilder strings.Builder
	if input.Limit != nil {
		idBuilder.WriteString(strconv.FormatInt(*input.Limit, 10))
	} else {
		idBuilder.WriteString("0")
	}
	idBuilder.WriteString(IdSeparator)
	if input.Offset != nil {
		idBuilder.WriteString(strconv.FormatInt(*input.Offset, 10))
	} else {
		idBuilder.WriteString("0")
	}
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Folder)
	idBuilder.WriteString(IdSeparator)
	if input.Name != nil {
		idBuilder.WriteString(*input.Name)
	}
	state.Id = types.StringValue(idBuilder.String())
	var var0 []authenticationProfilesListDsModelConfig
	if len(ans.Data) != 0 {
		var0 = make([]authenticationProfilesListDsModelConfig, 0, len(ans.Data))
		for var1Index := range ans.Data {
			var1 := ans.Data[var1Index]
			var var2 authenticationProfilesListDsModelConfig
			var var3 *authenticationProfilesListDsModelLockoutObject
			if var1.Lockout != nil {
				var3 = &authenticationProfilesListDsModelLockoutObject{}
				var3.FailedAttempts = types.Int64Value(var1.Lockout.FailedAttempts)
				var3.LockoutTime = types.Int64Value(var1.Lockout.LockoutTime)
			}
			var var4 *authenticationProfilesListDsModelMethodObject
			if var1.Method != nil {
				var4 = &authenticationProfilesListDsModelMethodObject{}
				var var5 *authenticationProfilesListDsModelKerberosObject
				if var1.Method.Kerberos != nil {
					var5 = &authenticationProfilesListDsModelKerberosObject{}
					var5.Realm = types.StringValue(var1.Method.Kerberos.Realm)
					var5.ServerProfile = types.StringValue(var1.Method.Kerberos.ServerProfile)
				}
				var var6 *authenticationProfilesListDsModelLdapObject
				if var1.Method.Ldap != nil {
					var6 = &authenticationProfilesListDsModelLdapObject{}
					var6.LoginAttribute = types.StringValue(var1.Method.Ldap.LoginAttribute)
					var6.PasswdExpDays = types.Int64Value(var1.Method.Ldap.PasswdExpDays)
					var6.ServerProfile = types.StringValue(var1.Method.Ldap.ServerProfile)
				}
				var var7 *authenticationProfilesListDsModelRadiusObject
				if var1.Method.Radius != nil {
					var7 = &authenticationProfilesListDsModelRadiusObject{}
					var7.Checkgroup = types.BoolValue(var1.Method.Radius.Checkgroup)
					var7.ServerProfile = types.StringValue(var1.Method.Radius.ServerProfile)
				}
				var var8 *authenticationProfilesListDsModelSamlIdpObject
				if var1.Method.SamlIdp != nil {
					var8 = &authenticationProfilesListDsModelSamlIdpObject{}
					var8.AttributeNameUsergroup = types.StringValue(var1.Method.SamlIdp.AttributeNameUsergroup)
					var8.AttributeNameUsername = types.StringValue(var1.Method.SamlIdp.AttributeNameUsername)
					var8.CertificateProfile = types.StringValue(var1.Method.SamlIdp.CertificateProfile)
					var8.EnableSingleLogout = types.BoolValue(var1.Method.SamlIdp.EnableSingleLogout)
					var8.RequestSigningCertificate = types.StringValue(var1.Method.SamlIdp.RequestSigningCertificate)
					var8.ServerProfile = types.StringValue(var1.Method.SamlIdp.ServerProfile)
				}
				var var9 *authenticationProfilesListDsModelTacplusObject
				if var1.Method.Tacplus != nil {
					var9 = &authenticationProfilesListDsModelTacplusObject{}
					var9.Checkgroup = types.BoolValue(var1.Method.Tacplus.Checkgroup)
					var9.ServerProfile = types.StringValue(var1.Method.Tacplus.ServerProfile)
				}
				var4.Kerberos = var5
				var4.Ldap = var6
				if var1.Method.LocalDatabase != nil {
					var4.LocalDatabase = types.BoolValue(true)
				}
				var4.Radius = var7
				var4.SamlIdp = var8
				var4.Tacplus = var9
			}
			var var10 *authenticationProfilesListDsModelMultiFactorAuthObject
			if var1.MultiFactorAuth != nil {
				var10 = &authenticationProfilesListDsModelMultiFactorAuthObject{}
				var10.Factors = EncodeStringSlice(var1.MultiFactorAuth.Factors)
				var10.MfaEnable = types.BoolValue(var1.MultiFactorAuth.MfaEnable)
			}
			var var11 *authenticationProfilesListDsModelSingleSignOnObject
			if var1.SingleSignOn != nil {
				var11 = &authenticationProfilesListDsModelSingleSignOnObject{}
				var11.KerberosKeytab = types.StringValue(var1.SingleSignOn.KerberosKeytab)
				var11.Realm = types.StringValue(var1.SingleSignOn.Realm)
			}
			var2.AllowList = EncodeStringSlice(var1.AllowList)
			var2.ObjectId = types.StringValue(var1.ObjectId)
			var2.Lockout = var3
			var2.Method = var4
			var2.MultiFactorAuth = var10
			var2.Name = types.StringValue(var1.Name)
			var2.SingleSignOn = var11
			var2.UserDomain = types.StringValue(var1.UserDomain)
			var2.UsernameModifier = types.StringValue(var1.UsernameModifier)
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
	_ datasource.DataSource              = &authenticationProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &authenticationProfilesDataSource{}
)

func NewAuthenticationProfilesDataSource() datasource.DataSource {
	return &authenticationProfilesDataSource{}
}

type authenticationProfilesDataSource struct {
	client *sase.Client
}

type authenticationProfilesDsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	ObjectId types.String `tfsdk:"object_id"`
	Folder   types.String `tfsdk:"folder"`

	// Output.
	// Ref: #/components/schemas/authentication-profiles
	AllowList []types.String `tfsdk:"allow_list"`
	// input omit: ObjectId
	Lockout          *authenticationProfilesDsModelLockoutObject         `tfsdk:"lockout"`
	Method           *authenticationProfilesDsModelMethodObject          `tfsdk:"method"`
	MultiFactorAuth  *authenticationProfilesDsModelMultiFactorAuthObject `tfsdk:"multi_factor_auth"`
	Name             types.String                                        `tfsdk:"name"`
	SingleSignOn     *authenticationProfilesDsModelSingleSignOnObject    `tfsdk:"single_sign_on"`
	UserDomain       types.String                                        `tfsdk:"user_domain"`
	UsernameModifier types.String                                        `tfsdk:"username_modifier"`
}

type authenticationProfilesDsModelLockoutObject struct {
	FailedAttempts types.Int64 `tfsdk:"failed_attempts"`
	LockoutTime    types.Int64 `tfsdk:"lockout_time"`
}

type authenticationProfilesDsModelMethodObject struct {
	Kerberos      *authenticationProfilesDsModelKerberosObject `tfsdk:"kerberos"`
	Ldap          *authenticationProfilesDsModelLdapObject     `tfsdk:"ldap"`
	LocalDatabase types.Bool                                   `tfsdk:"local_database"`
	Radius        *authenticationProfilesDsModelRadiusObject   `tfsdk:"radius"`
	SamlIdp       *authenticationProfilesDsModelSamlIdpObject  `tfsdk:"saml_idp"`
	Tacplus       *authenticationProfilesDsModelTacplusObject  `tfsdk:"tacplus"`
}

type authenticationProfilesDsModelKerberosObject struct {
	Realm         types.String `tfsdk:"realm"`
	ServerProfile types.String `tfsdk:"server_profile"`
}

type authenticationProfilesDsModelLdapObject struct {
	LoginAttribute types.String `tfsdk:"login_attribute"`
	PasswdExpDays  types.Int64  `tfsdk:"passwd_exp_days"`
	ServerProfile  types.String `tfsdk:"server_profile"`
}

type authenticationProfilesDsModelRadiusObject struct {
	Checkgroup    types.Bool   `tfsdk:"checkgroup"`
	ServerProfile types.String `tfsdk:"server_profile"`
}

type authenticationProfilesDsModelSamlIdpObject struct {
	AttributeNameUsergroup    types.String `tfsdk:"attribute_name_usergroup"`
	AttributeNameUsername     types.String `tfsdk:"attribute_name_username"`
	CertificateProfile        types.String `tfsdk:"certificate_profile"`
	EnableSingleLogout        types.Bool   `tfsdk:"enable_single_logout"`
	RequestSigningCertificate types.String `tfsdk:"request_signing_certificate"`
	ServerProfile             types.String `tfsdk:"server_profile"`
}

type authenticationProfilesDsModelTacplusObject struct {
	Checkgroup    types.Bool   `tfsdk:"checkgroup"`
	ServerProfile types.String `tfsdk:"server_profile"`
}

type authenticationProfilesDsModelMultiFactorAuthObject struct {
	Factors   []types.String `tfsdk:"factors"`
	MfaEnable types.Bool     `tfsdk:"mfa_enable"`
}

type authenticationProfilesDsModelSingleSignOnObject struct {
	KerberosKeytab types.String `tfsdk:"kerberos_keytab"`
	Realm          types.String `tfsdk:"realm"`
}

// Metadata returns the data source type name.
func (d *authenticationProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_profiles"
}

// Schema defines the schema for this listing data source.
func (d *authenticationProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"folder": dsschema.StringAttribute{
				Description: "The folder of the entry",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Shared", "Mobile Users", "Remote Networks", "Service Connections", "Mobile Users Container", "Mobile Users Explicit Proxy"),
				},
			},

			// Output.
			"allow_list": dsschema.ListAttribute{
				Description: "",
				Computed:    true,
				ElementType: types.StringType,
			},
			"lockout": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"failed_attempts": dsschema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
					"lockout_time": dsschema.Int64Attribute{
						Description: "",
						Computed:    true,
					},
				},
			},
			"method": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"kerberos": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"realm": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"server_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"ldap": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"login_attribute": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"passwd_exp_days": dsschema.Int64Attribute{
								Description: "",
								Computed:    true,
							},
							"server_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"local_database": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
					"radius": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"checkgroup": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"server_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"saml_idp": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"attribute_name_usergroup": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"attribute_name_username": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"certificate_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"enable_single_logout": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"request_signing_certificate": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
							"server_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
					"tacplus": dsschema.SingleNestedAttribute{
						Description: "",
						Computed:    true,
						Attributes: map[string]dsschema.Attribute{
							"checkgroup": dsschema.BoolAttribute{
								Description: "",
								Computed:    true,
							},
							"server_profile": dsschema.StringAttribute{
								Description: "",
								Computed:    true,
							},
						},
					},
				},
			},
			"multi_factor_auth": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"factors": dsschema.ListAttribute{
						Description: "",
						Computed:    true,
						ElementType: types.StringType,
					},
					"mfa_enable": dsschema.BoolAttribute{
						Description: "",
						Computed:    true,
					},
				},
			},
			"name": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"single_sign_on": dsschema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"kerberos_keytab": dsschema.StringAttribute{
						Description: "",
						Computed:    true,
					},
					"realm": dsschema.StringAttribute{
						Description: "",
						Computed:    true,
					},
				},
			},
			"user_domain": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"username_modifier": dsschema.StringAttribute{
				Description: "",
				Computed:    true,
			},
		},
	}
}

// Configure prepares the struct.
func (d *authenticationProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sase.Client)
}

func (d *authenticationProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationProfilesDsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing data source singleton retrieval", map[string]any{
		"terraform_provider_function": "Read",
		"data_source_name":            "sase_authentication_profiles",
		"object_id":                   state.ObjectId.ValueString(),
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to run the command.
	svc := cUCsSiw.NewClient(d.client)
	input := cUCsSiw.ReadInput{
		ObjectId: state.ObjectId.ValueString(),
		Folder:   state.Folder.ValueString(),
	}

	// Perform the operation.
	ans, err := svc.Read(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error getting singleton", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString(input.ObjectId)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(input.Folder)
	state.Id = types.StringValue(idBuilder.String())
	var var0 *authenticationProfilesDsModelLockoutObject
	if ans.Lockout != nil {
		var0 = &authenticationProfilesDsModelLockoutObject{}
		var0.FailedAttempts = types.Int64Value(ans.Lockout.FailedAttempts)
		var0.LockoutTime = types.Int64Value(ans.Lockout.LockoutTime)
	}
	var var1 *authenticationProfilesDsModelMethodObject
	if ans.Method != nil {
		var1 = &authenticationProfilesDsModelMethodObject{}
		var var2 *authenticationProfilesDsModelKerberosObject
		if ans.Method.Kerberos != nil {
			var2 = &authenticationProfilesDsModelKerberosObject{}
			var2.Realm = types.StringValue(ans.Method.Kerberos.Realm)
			var2.ServerProfile = types.StringValue(ans.Method.Kerberos.ServerProfile)
		}
		var var3 *authenticationProfilesDsModelLdapObject
		if ans.Method.Ldap != nil {
			var3 = &authenticationProfilesDsModelLdapObject{}
			var3.LoginAttribute = types.StringValue(ans.Method.Ldap.LoginAttribute)
			var3.PasswdExpDays = types.Int64Value(ans.Method.Ldap.PasswdExpDays)
			var3.ServerProfile = types.StringValue(ans.Method.Ldap.ServerProfile)
		}
		var var4 *authenticationProfilesDsModelRadiusObject
		if ans.Method.Radius != nil {
			var4 = &authenticationProfilesDsModelRadiusObject{}
			var4.Checkgroup = types.BoolValue(ans.Method.Radius.Checkgroup)
			var4.ServerProfile = types.StringValue(ans.Method.Radius.ServerProfile)
		}
		var var5 *authenticationProfilesDsModelSamlIdpObject
		if ans.Method.SamlIdp != nil {
			var5 = &authenticationProfilesDsModelSamlIdpObject{}
			var5.AttributeNameUsergroup = types.StringValue(ans.Method.SamlIdp.AttributeNameUsergroup)
			var5.AttributeNameUsername = types.StringValue(ans.Method.SamlIdp.AttributeNameUsername)
			var5.CertificateProfile = types.StringValue(ans.Method.SamlIdp.CertificateProfile)
			var5.EnableSingleLogout = types.BoolValue(ans.Method.SamlIdp.EnableSingleLogout)
			var5.RequestSigningCertificate = types.StringValue(ans.Method.SamlIdp.RequestSigningCertificate)
			var5.ServerProfile = types.StringValue(ans.Method.SamlIdp.ServerProfile)
		}
		var var6 *authenticationProfilesDsModelTacplusObject
		if ans.Method.Tacplus != nil {
			var6 = &authenticationProfilesDsModelTacplusObject{}
			var6.Checkgroup = types.BoolValue(ans.Method.Tacplus.Checkgroup)
			var6.ServerProfile = types.StringValue(ans.Method.Tacplus.ServerProfile)
		}
		var1.Kerberos = var2
		var1.Ldap = var3
		if ans.Method.LocalDatabase != nil {
			var1.LocalDatabase = types.BoolValue(true)
		}
		var1.Radius = var4
		var1.SamlIdp = var5
		var1.Tacplus = var6
	}
	var var7 *authenticationProfilesDsModelMultiFactorAuthObject
	if ans.MultiFactorAuth != nil {
		var7 = &authenticationProfilesDsModelMultiFactorAuthObject{}
		var7.Factors = EncodeStringSlice(ans.MultiFactorAuth.Factors)
		var7.MfaEnable = types.BoolValue(ans.MultiFactorAuth.MfaEnable)
	}
	var var8 *authenticationProfilesDsModelSingleSignOnObject
	if ans.SingleSignOn != nil {
		var8 = &authenticationProfilesDsModelSingleSignOnObject{}
		var8.KerberosKeytab = types.StringValue(ans.SingleSignOn.KerberosKeytab)
		var8.Realm = types.StringValue(ans.SingleSignOn.Realm)
	}
	state.AllowList = EncodeStringSlice(ans.AllowList)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Lockout = var0
	state.Method = var1
	state.MultiFactorAuth = var7
	state.Name = types.StringValue(ans.Name)
	state.SingleSignOn = var8
	state.UserDomain = types.StringValue(ans.UserDomain)
	state.UsernameModifier = types.StringValue(ans.UsernameModifier)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Resource.
var (
	_ resource.Resource                = &authenticationProfilesResource{}
	_ resource.ResourceWithConfigure   = &authenticationProfilesResource{}
	_ resource.ResourceWithImportState = &authenticationProfilesResource{}
)

func NewAuthenticationProfilesResource() resource.Resource {
	return &authenticationProfilesResource{}
}

type authenticationProfilesResource struct {
	client *sase.Client
}

type authenticationProfilesRsModel struct {
	Id types.String `tfsdk:"id"`

	// Input.
	Folder types.String `tfsdk:"folder"`

	// Request body input.
	// Ref: #/components/schemas/authentication-profiles
	AllowList        []types.String                                      `tfsdk:"allow_list"`
	ObjectId         types.String                                        `tfsdk:"object_id"`
	Lockout          *authenticationProfilesRsModelLockoutObject         `tfsdk:"lockout"`
	Method           *authenticationProfilesRsModelMethodObject          `tfsdk:"method"`
	MultiFactorAuth  *authenticationProfilesRsModelMultiFactorAuthObject `tfsdk:"multi_factor_auth"`
	Name             types.String                                        `tfsdk:"name"`
	SingleSignOn     *authenticationProfilesRsModelSingleSignOnObject    `tfsdk:"single_sign_on"`
	UserDomain       types.String                                        `tfsdk:"user_domain"`
	UsernameModifier types.String                                        `tfsdk:"username_modifier"`
}

type authenticationProfilesRsModelLockoutObject struct {
	FailedAttempts types.Int64 `tfsdk:"failed_attempts"`
	LockoutTime    types.Int64 `tfsdk:"lockout_time"`
}

type authenticationProfilesRsModelMethodObject struct {
	Kerberos      *authenticationProfilesRsModelKerberosObject `tfsdk:"kerberos"`
	Ldap          *authenticationProfilesRsModelLdapObject     `tfsdk:"ldap"`
	LocalDatabase types.Bool                                   `tfsdk:"local_database"`
	Radius        *authenticationProfilesRsModelRadiusObject   `tfsdk:"radius"`
	SamlIdp       *authenticationProfilesRsModelSamlIdpObject  `tfsdk:"saml_idp"`
	Tacplus       *authenticationProfilesRsModelTacplusObject  `tfsdk:"tacplus"`
}

type authenticationProfilesRsModelKerberosObject struct {
	Realm         types.String `tfsdk:"realm"`
	ServerProfile types.String `tfsdk:"server_profile"`
}

type authenticationProfilesRsModelLdapObject struct {
	LoginAttribute types.String `tfsdk:"login_attribute"`
	PasswdExpDays  types.Int64  `tfsdk:"passwd_exp_days"`
	ServerProfile  types.String `tfsdk:"server_profile"`
}

type authenticationProfilesRsModelRadiusObject struct {
	Checkgroup    types.Bool   `tfsdk:"checkgroup"`
	ServerProfile types.String `tfsdk:"server_profile"`
}

type authenticationProfilesRsModelSamlIdpObject struct {
	AttributeNameUsergroup    types.String `tfsdk:"attribute_name_usergroup"`
	AttributeNameUsername     types.String `tfsdk:"attribute_name_username"`
	CertificateProfile        types.String `tfsdk:"certificate_profile"`
	EnableSingleLogout        types.Bool   `tfsdk:"enable_single_logout"`
	RequestSigningCertificate types.String `tfsdk:"request_signing_certificate"`
	ServerProfile             types.String `tfsdk:"server_profile"`
}

type authenticationProfilesRsModelTacplusObject struct {
	Checkgroup    types.Bool   `tfsdk:"checkgroup"`
	ServerProfile types.String `tfsdk:"server_profile"`
}

type authenticationProfilesRsModelMultiFactorAuthObject struct {
	Factors   []types.String `tfsdk:"factors"`
	MfaEnable types.Bool     `tfsdk:"mfa_enable"`
}

type authenticationProfilesRsModelSingleSignOnObject struct {
	KerberosKeytab types.String `tfsdk:"kerberos_keytab"`
	Realm          types.String `tfsdk:"realm"`
}

// Metadata returns the data source type name.
func (r *authenticationProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_profiles"
}

// Schema defines the schema for this listing data source.
func (r *authenticationProfilesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

			"allow_list": rsschema.ListAttribute{
				Description: "",
				Optional:    true,
				ElementType: types.StringType,
			},
			"object_id": rsschema.StringAttribute{
				Description: "",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lockout": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"failed_attempts": rsschema.Int64Attribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(0, 10),
						},
					},
					"lockout_time": rsschema.Int64Attribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							DefaultInt64(0),
						},
						Validators: []validator.Int64{
							int64validator.Between(0, 60),
						},
					},
				},
			},
			"method": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"kerberos": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"realm": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"server_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"ldap": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"login_attribute": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
							"passwd_exp_days": rsschema.Int64Attribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Int64{
									DefaultInt64(0),
								},
							},
							"server_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"local_database": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
					},
					"radius": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"checkgroup": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"server_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
							},
						},
					},
					"saml_idp": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"attribute_name_usergroup": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 63),
								},
							},
							"attribute_name_username": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthBetween(1, 63),
								},
							},
							"certificate_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthAtMost(31),
								},
							},
							"enable_single_logout": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"request_signing_certificate": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthAtMost(64),
								},
							},
							"server_profile": rsschema.StringAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									DefaultString(""),
								},
								Validators: []validator.String{
									stringvalidator.LengthAtMost(63),
								},
							},
						},
					},
					"tacplus": rsschema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]rsschema.Attribute{
							"checkgroup": rsschema.BoolAttribute{
								Description: "",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									DefaultBool(false),
								},
							},
							"server_profile": rsschema.StringAttribute{
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
			"multi_factor_auth": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"factors": rsschema.ListAttribute{
						Description: "",
						Optional:    true,
						ElementType: types.StringType,
					},
					"mfa_enable": rsschema.BoolAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							DefaultBool(false),
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
			},
			"single_sign_on": rsschema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]rsschema.Attribute{
					"kerberos_keytab": rsschema.StringAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtMost(8192),
						},
					},
					"realm": rsschema.StringAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							DefaultString(""),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtMost(127),
						},
					},
				},
			},
			"user_domain": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(63),
				},
			},
			"username_modifier": rsschema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					DefaultString(""),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("%USERINPUT%", "%USERINPUT%@%USERDOMAIN%", "%USERDOMAIN%\\%USERINPUT%"),
				},
			},
		},
	}
}

// Configure prepares the struct.
func (r *authenticationProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*sase.Client)
}

// Create resource
func (r *authenticationProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state authenticationProfilesRsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Basic logging.
	tflog.Info(ctx, "performing resource create", map[string]any{
		"terraform_provider_function": "Create",
		"resource_name":               "sase_authentication_profiles",
		"folder":                      state.Folder.ValueString(),
	})

	// Prepare to create the config.
	svc := cUCsSiw.NewClient(r.client)
	input := cUCsSiw.CreateInput{
		Folder: state.Folder.ValueString(),
	}
	var var0 alljvhu.Config
	var0.AllowList = DecodeStringSlice(state.AllowList)
	var var1 *alljvhu.LockoutObject
	if state.Lockout != nil {
		var1 = &alljvhu.LockoutObject{}
		var1.FailedAttempts = state.Lockout.FailedAttempts.ValueInt64()
		var1.LockoutTime = state.Lockout.LockoutTime.ValueInt64()
	}
	var0.Lockout = var1
	var var2 *alljvhu.MethodObject
	if state.Method != nil {
		var2 = &alljvhu.MethodObject{}
		var var3 *alljvhu.KerberosObject
		if state.Method.Kerberos != nil {
			var3 = &alljvhu.KerberosObject{}
			var3.Realm = state.Method.Kerberos.Realm.ValueString()
			var3.ServerProfile = state.Method.Kerberos.ServerProfile.ValueString()
		}
		var2.Kerberos = var3
		var var4 *alljvhu.LdapObject
		if state.Method.Ldap != nil {
			var4 = &alljvhu.LdapObject{}
			var4.LoginAttribute = state.Method.Ldap.LoginAttribute.ValueString()
			var4.PasswdExpDays = state.Method.Ldap.PasswdExpDays.ValueInt64()
			var4.ServerProfile = state.Method.Ldap.ServerProfile.ValueString()
		}
		var2.Ldap = var4
		if state.Method.LocalDatabase.ValueBool() {
			var2.LocalDatabase = struct{}{}
		}
		var var5 *alljvhu.RadiusObject
		if state.Method.Radius != nil {
			var5 = &alljvhu.RadiusObject{}
			var5.Checkgroup = state.Method.Radius.Checkgroup.ValueBool()
			var5.ServerProfile = state.Method.Radius.ServerProfile.ValueString()
		}
		var2.Radius = var5
		var var6 *alljvhu.SamlIdpObject
		if state.Method.SamlIdp != nil {
			var6 = &alljvhu.SamlIdpObject{}
			var6.AttributeNameUsergroup = state.Method.SamlIdp.AttributeNameUsergroup.ValueString()
			var6.AttributeNameUsername = state.Method.SamlIdp.AttributeNameUsername.ValueString()
			var6.CertificateProfile = state.Method.SamlIdp.CertificateProfile.ValueString()
			var6.EnableSingleLogout = state.Method.SamlIdp.EnableSingleLogout.ValueBool()
			var6.RequestSigningCertificate = state.Method.SamlIdp.RequestSigningCertificate.ValueString()
			var6.ServerProfile = state.Method.SamlIdp.ServerProfile.ValueString()
		}
		var2.SamlIdp = var6
		var var7 *alljvhu.TacplusObject
		if state.Method.Tacplus != nil {
			var7 = &alljvhu.TacplusObject{}
			var7.Checkgroup = state.Method.Tacplus.Checkgroup.ValueBool()
			var7.ServerProfile = state.Method.Tacplus.ServerProfile.ValueString()
		}
		var2.Tacplus = var7
	}
	var0.Method = var2
	var var8 *alljvhu.MultiFactorAuthObject
	if state.MultiFactorAuth != nil {
		var8 = &alljvhu.MultiFactorAuthObject{}
		var8.Factors = DecodeStringSlice(state.MultiFactorAuth.Factors)
		var8.MfaEnable = state.MultiFactorAuth.MfaEnable.ValueBool()
	}
	var0.MultiFactorAuth = var8
	var0.Name = state.Name.ValueString()
	var var9 *alljvhu.SingleSignOnObject
	if state.SingleSignOn != nil {
		var9 = &alljvhu.SingleSignOnObject{}
		var9.KerberosKeytab = state.SingleSignOn.KerberosKeytab.ValueString()
		var9.Realm = state.SingleSignOn.Realm.ValueString()
	}
	var0.SingleSignOn = var9
	var0.UserDomain = state.UserDomain.ValueString()
	var0.UsernameModifier = state.UsernameModifier.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Create(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in create", err.Error())
		return
	}

	// Store the answer to state.
	var idBuilder strings.Builder
	idBuilder.WriteString(input.Folder)
	idBuilder.WriteString(IdSeparator)
	idBuilder.WriteString(ans.ObjectId)
	state.Id = types.StringValue(idBuilder.String())
	var var10 *authenticationProfilesRsModelLockoutObject
	if ans.Lockout != nil {
		var10 = &authenticationProfilesRsModelLockoutObject{}
		var10.FailedAttempts = types.Int64Value(ans.Lockout.FailedAttempts)
		var10.LockoutTime = types.Int64Value(ans.Lockout.LockoutTime)
	}
	var var11 *authenticationProfilesRsModelMethodObject
	if ans.Method != nil {
		var11 = &authenticationProfilesRsModelMethodObject{}
		var var12 *authenticationProfilesRsModelKerberosObject
		if ans.Method.Kerberos != nil {
			var12 = &authenticationProfilesRsModelKerberosObject{}
			var12.Realm = types.StringValue(ans.Method.Kerberos.Realm)
			var12.ServerProfile = types.StringValue(ans.Method.Kerberos.ServerProfile)
		}
		var var13 *authenticationProfilesRsModelLdapObject
		if ans.Method.Ldap != nil {
			var13 = &authenticationProfilesRsModelLdapObject{}
			var13.LoginAttribute = types.StringValue(ans.Method.Ldap.LoginAttribute)
			var13.PasswdExpDays = types.Int64Value(ans.Method.Ldap.PasswdExpDays)
			var13.ServerProfile = types.StringValue(ans.Method.Ldap.ServerProfile)
		}
		var var14 *authenticationProfilesRsModelRadiusObject
		if ans.Method.Radius != nil {
			var14 = &authenticationProfilesRsModelRadiusObject{}
			var14.Checkgroup = types.BoolValue(ans.Method.Radius.Checkgroup)
			var14.ServerProfile = types.StringValue(ans.Method.Radius.ServerProfile)
		}
		var var15 *authenticationProfilesRsModelSamlIdpObject
		if ans.Method.SamlIdp != nil {
			var15 = &authenticationProfilesRsModelSamlIdpObject{}
			var15.AttributeNameUsergroup = types.StringValue(ans.Method.SamlIdp.AttributeNameUsergroup)
			var15.AttributeNameUsername = types.StringValue(ans.Method.SamlIdp.AttributeNameUsername)
			var15.CertificateProfile = types.StringValue(ans.Method.SamlIdp.CertificateProfile)
			var15.EnableSingleLogout = types.BoolValue(ans.Method.SamlIdp.EnableSingleLogout)
			var15.RequestSigningCertificate = types.StringValue(ans.Method.SamlIdp.RequestSigningCertificate)
			var15.ServerProfile = types.StringValue(ans.Method.SamlIdp.ServerProfile)
		}
		var var16 *authenticationProfilesRsModelTacplusObject
		if ans.Method.Tacplus != nil {
			var16 = &authenticationProfilesRsModelTacplusObject{}
			var16.Checkgroup = types.BoolValue(ans.Method.Tacplus.Checkgroup)
			var16.ServerProfile = types.StringValue(ans.Method.Tacplus.ServerProfile)
		}
		var11.Kerberos = var12
		var11.Ldap = var13
		if ans.Method.LocalDatabase != nil {
			var11.LocalDatabase = types.BoolValue(true)
		}
		var11.Radius = var14
		var11.SamlIdp = var15
		var11.Tacplus = var16
	}
	var var17 *authenticationProfilesRsModelMultiFactorAuthObject
	if ans.MultiFactorAuth != nil {
		var17 = &authenticationProfilesRsModelMultiFactorAuthObject{}
		var17.Factors = EncodeStringSlice(ans.MultiFactorAuth.Factors)
		var17.MfaEnable = types.BoolValue(ans.MultiFactorAuth.MfaEnable)
	}
	var var18 *authenticationProfilesRsModelSingleSignOnObject
	if ans.SingleSignOn != nil {
		var18 = &authenticationProfilesRsModelSingleSignOnObject{}
		var18.KerberosKeytab = types.StringValue(ans.SingleSignOn.KerberosKeytab)
		var18.Realm = types.StringValue(ans.SingleSignOn.Realm)
	}
	state.AllowList = EncodeStringSlice(ans.AllowList)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Lockout = var10
	state.Method = var11
	state.MultiFactorAuth = var17
	state.Name = types.StringValue(ans.Name)
	state.SingleSignOn = var18
	state.UserDomain = types.StringValue(ans.UserDomain)
	state.UsernameModifier = types.StringValue(ans.UsernameModifier)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read resource.
func (r *authenticationProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	var state authenticationProfilesRsModel

	// Basic logging.
	tflog.Info(ctx, "performing resource read", map[string]any{
		"terraform_provider_function": "Read",
		"resource_name":               "sase_authentication_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	// Prepare to read the config.
	svc := cUCsSiw.NewClient(r.client)
	input := cUCsSiw.ReadInput{
		ObjectId: tokens[1],
		Folder:   tokens[0],
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
	var var0 *authenticationProfilesRsModelLockoutObject
	if ans.Lockout != nil {
		var0 = &authenticationProfilesRsModelLockoutObject{}
		var0.FailedAttempts = types.Int64Value(ans.Lockout.FailedAttempts)
		var0.LockoutTime = types.Int64Value(ans.Lockout.LockoutTime)
	}
	var var1 *authenticationProfilesRsModelMethodObject
	if ans.Method != nil {
		var1 = &authenticationProfilesRsModelMethodObject{}
		var var2 *authenticationProfilesRsModelKerberosObject
		if ans.Method.Kerberos != nil {
			var2 = &authenticationProfilesRsModelKerberosObject{}
			var2.Realm = types.StringValue(ans.Method.Kerberos.Realm)
			var2.ServerProfile = types.StringValue(ans.Method.Kerberos.ServerProfile)
		}
		var var3 *authenticationProfilesRsModelLdapObject
		if ans.Method.Ldap != nil {
			var3 = &authenticationProfilesRsModelLdapObject{}
			var3.LoginAttribute = types.StringValue(ans.Method.Ldap.LoginAttribute)
			var3.PasswdExpDays = types.Int64Value(ans.Method.Ldap.PasswdExpDays)
			var3.ServerProfile = types.StringValue(ans.Method.Ldap.ServerProfile)
		}
		var var4 *authenticationProfilesRsModelRadiusObject
		if ans.Method.Radius != nil {
			var4 = &authenticationProfilesRsModelRadiusObject{}
			var4.Checkgroup = types.BoolValue(ans.Method.Radius.Checkgroup)
			var4.ServerProfile = types.StringValue(ans.Method.Radius.ServerProfile)
		}
		var var5 *authenticationProfilesRsModelSamlIdpObject
		if ans.Method.SamlIdp != nil {
			var5 = &authenticationProfilesRsModelSamlIdpObject{}
			var5.AttributeNameUsergroup = types.StringValue(ans.Method.SamlIdp.AttributeNameUsergroup)
			var5.AttributeNameUsername = types.StringValue(ans.Method.SamlIdp.AttributeNameUsername)
			var5.CertificateProfile = types.StringValue(ans.Method.SamlIdp.CertificateProfile)
			var5.EnableSingleLogout = types.BoolValue(ans.Method.SamlIdp.EnableSingleLogout)
			var5.RequestSigningCertificate = types.StringValue(ans.Method.SamlIdp.RequestSigningCertificate)
			var5.ServerProfile = types.StringValue(ans.Method.SamlIdp.ServerProfile)
		}
		var var6 *authenticationProfilesRsModelTacplusObject
		if ans.Method.Tacplus != nil {
			var6 = &authenticationProfilesRsModelTacplusObject{}
			var6.Checkgroup = types.BoolValue(ans.Method.Tacplus.Checkgroup)
			var6.ServerProfile = types.StringValue(ans.Method.Tacplus.ServerProfile)
		}
		var1.Kerberos = var2
		var1.Ldap = var3
		if ans.Method.LocalDatabase != nil {
			var1.LocalDatabase = types.BoolValue(true)
		}
		var1.Radius = var4
		var1.SamlIdp = var5
		var1.Tacplus = var6
	}
	var var7 *authenticationProfilesRsModelMultiFactorAuthObject
	if ans.MultiFactorAuth != nil {
		var7 = &authenticationProfilesRsModelMultiFactorAuthObject{}
		var7.Factors = EncodeStringSlice(ans.MultiFactorAuth.Factors)
		var7.MfaEnable = types.BoolValue(ans.MultiFactorAuth.MfaEnable)
	}
	var var8 *authenticationProfilesRsModelSingleSignOnObject
	if ans.SingleSignOn != nil {
		var8 = &authenticationProfilesRsModelSingleSignOnObject{}
		var8.KerberosKeytab = types.StringValue(ans.SingleSignOn.KerberosKeytab)
		var8.Realm = types.StringValue(ans.SingleSignOn.Realm)
	}
	state.AllowList = EncodeStringSlice(ans.AllowList)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Lockout = var0
	state.Method = var1
	state.MultiFactorAuth = var7
	state.Name = types.StringValue(ans.Name)
	state.SingleSignOn = var8
	state.UserDomain = types.StringValue(ans.UserDomain)
	state.UsernameModifier = types.StringValue(ans.UsernameModifier)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update resource.
func (r *authenticationProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state authenticationProfilesRsModel
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
		"terraform_provider_function": "Update",
		"resource_name":               "sase_authentication_profiles",
		"object_id":                   state.ObjectId.ValueString(),
	})

	// Prepare to create the config.
	svc := cUCsSiw.NewClient(r.client)
	input := cUCsSiw.UpdateInput{
		ObjectId: state.ObjectId.ValueString(),
	}
	var var0 alljvhu.Config
	var0.AllowList = DecodeStringSlice(plan.AllowList)
	var var1 *alljvhu.LockoutObject
	if plan.Lockout != nil {
		var1 = &alljvhu.LockoutObject{}
		var1.FailedAttempts = plan.Lockout.FailedAttempts.ValueInt64()
		var1.LockoutTime = plan.Lockout.LockoutTime.ValueInt64()
	}
	var0.Lockout = var1
	var var2 *alljvhu.MethodObject
	if plan.Method != nil {
		var2 = &alljvhu.MethodObject{}
		var var3 *alljvhu.KerberosObject
		if plan.Method.Kerberos != nil {
			var3 = &alljvhu.KerberosObject{}
			var3.Realm = plan.Method.Kerberos.Realm.ValueString()
			var3.ServerProfile = plan.Method.Kerberos.ServerProfile.ValueString()
		}
		var2.Kerberos = var3
		var var4 *alljvhu.LdapObject
		if plan.Method.Ldap != nil {
			var4 = &alljvhu.LdapObject{}
			var4.LoginAttribute = plan.Method.Ldap.LoginAttribute.ValueString()
			var4.PasswdExpDays = plan.Method.Ldap.PasswdExpDays.ValueInt64()
			var4.ServerProfile = plan.Method.Ldap.ServerProfile.ValueString()
		}
		var2.Ldap = var4
		if plan.Method.LocalDatabase.ValueBool() {
			var2.LocalDatabase = struct{}{}
		}
		var var5 *alljvhu.RadiusObject
		if plan.Method.Radius != nil {
			var5 = &alljvhu.RadiusObject{}
			var5.Checkgroup = plan.Method.Radius.Checkgroup.ValueBool()
			var5.ServerProfile = plan.Method.Radius.ServerProfile.ValueString()
		}
		var2.Radius = var5
		var var6 *alljvhu.SamlIdpObject
		if plan.Method.SamlIdp != nil {
			var6 = &alljvhu.SamlIdpObject{}
			var6.AttributeNameUsergroup = plan.Method.SamlIdp.AttributeNameUsergroup.ValueString()
			var6.AttributeNameUsername = plan.Method.SamlIdp.AttributeNameUsername.ValueString()
			var6.CertificateProfile = plan.Method.SamlIdp.CertificateProfile.ValueString()
			var6.EnableSingleLogout = plan.Method.SamlIdp.EnableSingleLogout.ValueBool()
			var6.RequestSigningCertificate = plan.Method.SamlIdp.RequestSigningCertificate.ValueString()
			var6.ServerProfile = plan.Method.SamlIdp.ServerProfile.ValueString()
		}
		var2.SamlIdp = var6
		var var7 *alljvhu.TacplusObject
		if plan.Method.Tacplus != nil {
			var7 = &alljvhu.TacplusObject{}
			var7.Checkgroup = plan.Method.Tacplus.Checkgroup.ValueBool()
			var7.ServerProfile = plan.Method.Tacplus.ServerProfile.ValueString()
		}
		var2.Tacplus = var7
	}
	var0.Method = var2
	var var8 *alljvhu.MultiFactorAuthObject
	if plan.MultiFactorAuth != nil {
		var8 = &alljvhu.MultiFactorAuthObject{}
		var8.Factors = DecodeStringSlice(plan.MultiFactorAuth.Factors)
		var8.MfaEnable = plan.MultiFactorAuth.MfaEnable.ValueBool()
	}
	var0.MultiFactorAuth = var8
	var0.Name = plan.Name.ValueString()
	var var9 *alljvhu.SingleSignOnObject
	if plan.SingleSignOn != nil {
		var9 = &alljvhu.SingleSignOnObject{}
		var9.KerberosKeytab = plan.SingleSignOn.KerberosKeytab.ValueString()
		var9.Realm = plan.SingleSignOn.Realm.ValueString()
	}
	var0.SingleSignOn = var9
	var0.UserDomain = plan.UserDomain.ValueString()
	var0.UsernameModifier = plan.UsernameModifier.ValueString()
	input.Config = var0

	// Perform the operation.
	ans, err := svc.Update(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error in update", err.Error())
		return
	}

	// Store the answer to state.
	var var10 *authenticationProfilesRsModelLockoutObject
	if ans.Lockout != nil {
		var10 = &authenticationProfilesRsModelLockoutObject{}
		var10.FailedAttempts = types.Int64Value(ans.Lockout.FailedAttempts)
		var10.LockoutTime = types.Int64Value(ans.Lockout.LockoutTime)
	}
	var var11 *authenticationProfilesRsModelMethodObject
	if ans.Method != nil {
		var11 = &authenticationProfilesRsModelMethodObject{}
		var var12 *authenticationProfilesRsModelKerberosObject
		if ans.Method.Kerberos != nil {
			var12 = &authenticationProfilesRsModelKerberosObject{}
			var12.Realm = types.StringValue(ans.Method.Kerberos.Realm)
			var12.ServerProfile = types.StringValue(ans.Method.Kerberos.ServerProfile)
		}
		var var13 *authenticationProfilesRsModelLdapObject
		if ans.Method.Ldap != nil {
			var13 = &authenticationProfilesRsModelLdapObject{}
			var13.LoginAttribute = types.StringValue(ans.Method.Ldap.LoginAttribute)
			var13.PasswdExpDays = types.Int64Value(ans.Method.Ldap.PasswdExpDays)
			var13.ServerProfile = types.StringValue(ans.Method.Ldap.ServerProfile)
		}
		var var14 *authenticationProfilesRsModelRadiusObject
		if ans.Method.Radius != nil {
			var14 = &authenticationProfilesRsModelRadiusObject{}
			var14.Checkgroup = types.BoolValue(ans.Method.Radius.Checkgroup)
			var14.ServerProfile = types.StringValue(ans.Method.Radius.ServerProfile)
		}
		var var15 *authenticationProfilesRsModelSamlIdpObject
		if ans.Method.SamlIdp != nil {
			var15 = &authenticationProfilesRsModelSamlIdpObject{}
			var15.AttributeNameUsergroup = types.StringValue(ans.Method.SamlIdp.AttributeNameUsergroup)
			var15.AttributeNameUsername = types.StringValue(ans.Method.SamlIdp.AttributeNameUsername)
			var15.CertificateProfile = types.StringValue(ans.Method.SamlIdp.CertificateProfile)
			var15.EnableSingleLogout = types.BoolValue(ans.Method.SamlIdp.EnableSingleLogout)
			var15.RequestSigningCertificate = types.StringValue(ans.Method.SamlIdp.RequestSigningCertificate)
			var15.ServerProfile = types.StringValue(ans.Method.SamlIdp.ServerProfile)
		}
		var var16 *authenticationProfilesRsModelTacplusObject
		if ans.Method.Tacplus != nil {
			var16 = &authenticationProfilesRsModelTacplusObject{}
			var16.Checkgroup = types.BoolValue(ans.Method.Tacplus.Checkgroup)
			var16.ServerProfile = types.StringValue(ans.Method.Tacplus.ServerProfile)
		}
		var11.Kerberos = var12
		var11.Ldap = var13
		if ans.Method.LocalDatabase != nil {
			var11.LocalDatabase = types.BoolValue(true)
		}
		var11.Radius = var14
		var11.SamlIdp = var15
		var11.Tacplus = var16
	}
	var var17 *authenticationProfilesRsModelMultiFactorAuthObject
	if ans.MultiFactorAuth != nil {
		var17 = &authenticationProfilesRsModelMultiFactorAuthObject{}
		var17.Factors = EncodeStringSlice(ans.MultiFactorAuth.Factors)
		var17.MfaEnable = types.BoolValue(ans.MultiFactorAuth.MfaEnable)
	}
	var var18 *authenticationProfilesRsModelSingleSignOnObject
	if ans.SingleSignOn != nil {
		var18 = &authenticationProfilesRsModelSingleSignOnObject{}
		var18.KerberosKeytab = types.StringValue(ans.SingleSignOn.KerberosKeytab)
		var18.Realm = types.StringValue(ans.SingleSignOn.Realm)
	}
	state.AllowList = EncodeStringSlice(ans.AllowList)
	state.ObjectId = types.StringValue(ans.ObjectId)
	state.Lockout = var10
	state.Method = var11
	state.MultiFactorAuth = var17
	state.Name = types.StringValue(ans.Name)
	state.SingleSignOn = var18
	state.UserDomain = types.StringValue(ans.UserDomain)
	state.UsernameModifier = types.StringValue(ans.UsernameModifier)

	// Done.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete resource.
func (r *authenticationProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
		"terraform_provider_function": "Delete",
		"resource_name":               "sase_authentication_profiles",
		"locMap":                      map[string]int{"Folder": 0, "ObjectId": 1},
		"tokens":                      tokens,
	})

	svc := cUCsSiw.NewClient(r.client)
	input := cUCsSiw.DeleteInput{
		ObjectId: tokens[1],
	}

	// Perform the operation.
	if _, err := svc.Delete(ctx, input); err != nil && !IsObjectNotFound(err) {
		resp.Diagnostics.AddError("Error in delete", err.Error())
	}
}

func (r *authenticationProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
