// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/comradesharf/terraform-provider-hostinger/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-nettypes/iptypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AgencyHostingWebsiteSSLCertModel struct {
	Names     []types.String    `tfsdk:"names"`
	ExpiresAt timetypes.RFC3339 `tfsdk:"expires_at"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
}

func (m *AgencyHostingWebsiteSSLCertModel) Merge(item client.AgencyHostingV1WebsitesSslCertResource) {
	m.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)
	m.ExpiresAt = timetypes.NewRFC3339TimePointerValue(item.ExpiresAt)

	if item.Names != nil {
		m.Names = make([]types.String, len(*item.Names))
		for j, name := range *item.Names {
			m.Names[j] = types.StringValue(name)
		}
	}
}

type AgencyHostingWebsiteCustomSSLCertModel struct {
	IsExpired types.Bool        `tfsdk:"is_expired"`
	ExpiresAt timetypes.RFC3339 `tfsdk:"expires_at"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
}

func (m *AgencyHostingWebsiteCustomSSLCertModel) Merge(item client.AgencyHostingV1WebsitesCustomSslCertResource) {
	m.IsExpired = types.BoolPointerValue(item.IsExpired)
	m.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)
	m.ExpiresAt = timetypes.NewRFC3339TimePointerValue(item.ExpiresAt)
}

type AgencyHostingWebsitePreviewDomainModel struct {
	FQDN      types.String      `tfsdk:"fqdn"`
	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
}

func (m *AgencyHostingWebsitePreviewDomainModel) Merge(item client.AgencyHostingV1WebsitesWebsitePreviewDomainResource) {
	m.FQDN = types.StringPointerValue(item.Fqdn)
	m.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)
}

type AgencyHostingWebsiteDomainModel struct {
	FQDN          types.String                            `tfsdk:"fqdn"`
	ParentFQDN    types.String                            `tfsdk:"parent_fqdn"`
	IPv6          iptypes.IPv6Address                     `tfsdk:"ipv6"`
	CreatedAt     timetypes.RFC3339                       `tfsdk:"created_at"`
	Nameservers   []types.String                          `tfsdk:"nameservers"`
	SSLCert       *AgencyHostingWebsiteSSLCertModel       `tfsdk:"ssl_cert"`
	CustomSSLCert *AgencyHostingWebsiteCustomSSLCertModel `tfsdk:"custom_ssl_cert"`
}

func (m *AgencyHostingWebsiteDomainModel) Merge(item client.AgencyHostingV1WebsitesWebsiteDomainDetailsResource) {
	m.FQDN = types.StringPointerValue(item.Fqdn)
	m.ParentFQDN = types.StringPointerValue(item.ParentFqdn)
	m.IPv6 = iptypes.NewIPv6AddressPointerValue(item.Ipv6)
	m.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)

	if item.Nameservers != nil {
		m.Nameservers = make([]types.String, len(*item.Nameservers))
		for j, nameserver := range *item.Nameservers {
			m.Nameservers[j] = types.StringValue(nameserver)
		}
	}

	if item.SslCert != nil {
		v, err := item.SslCert.AsAgencyHostingV1WebsitesSslCertResource()
		if err == nil {
			var s AgencyHostingWebsiteSSLCertModel
			s.Merge(v)
			m.SSLCert = &s
		}
	}

	if item.CustomSslCert != nil {
		v, err := item.CustomSslCert.AsAgencyHostingV1WebsitesCustomSslCertResource()
		if err == nil {
			var s AgencyHostingWebsiteCustomSSLCertModel
			s.Merge(v)
			m.CustomSSLCert = &s
		}
	}
}

type AgencyHostingWebsiteSettingsPHPModel struct {
	Version types.String `tfsdk:"version"`
	Workers types.Int32  `tfsdk:"workers"`
}

func (m *AgencyHostingWebsiteSettingsPHPModel) Merge(item client.AgencyHostingV1WebsitesWebsitePhpSettingsResource) {
	m.Version = types.StringPointerValue(item.Version)
	m.Workers = int32Value(item.Workers)
}

type AgencyHostingWebsiteSettingsModel struct {
	PHP AgencyHostingWebsiteSettingsPHPModel `tfsdk:"php"`
}

func (m *AgencyHostingWebsiteSettingsModel) Merge(item client.AgencyHostingV1WebsitesWebsiteSettingsResource) {
	if item.Php != nil {
		var php AgencyHostingWebsiteSettingsPHPModel
		v, err := item.Php.AsAgencyHostingV1WebsitesWebsitePhpSettingsResource()
		if err == nil {
			php.Merge(v)
		}
		m.PHP = php
	}
}

type AgencyHostingWebsiteWordpressModel struct {
	Domain         types.String      `tfsdk:"domain"`
	Title          types.String      `tfsdk:"title"`
	Language       types.String      `tfsdk:"language"`
	IsConfigLocked types.Bool        `tfsdk:"is_config_locked"`
	CreatedAt      timetypes.RFC3339 `tfsdk:"created_at"`
}

func (m *AgencyHostingWebsiteWordpressModel) Merge(item client.AgencyHostingV1WebsitesWordPressInstallResource) {
	m.Domain = types.StringPointerValue(item.Domain)
	m.Title = types.StringPointerValue(item.Title)
	m.Language = types.StringPointerValue(item.Language)
	m.IsConfigLocked = types.BoolPointerValue(item.IsConfigLocked)
	m.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)
}

type AgencyHostingWebsiteRemoteAccessSSHModel struct {
	Username          types.String `tfsdk:"username"`
	Host              types.String `tfsdk:"host"`
	Port              types.Int32  `tfsdk:"port"`
	IsEnabled         types.Bool   `tfsdk:"is_enabled"`
	IsPasswordEnabled types.Bool   `tfsdk:"is_password_enabled"`
}

func (m *AgencyHostingWebsiteRemoteAccessSSHModel) Merge(item client.AgencyHostingV1WebsitesWebsiteSshDetailsResource) {
	m.Username = types.StringPointerValue(item.Username)
	m.Host = types.StringPointerValue(item.Host)
	m.Port = int32Value(item.Port)
	m.IsEnabled = types.BoolPointerValue(item.IsEnabled)
	m.IsPasswordEnabled = types.BoolPointerValue(item.IsPasswordEnabled)
}

type AgencyHostingWebsiteRemoteAccessSFTPModel struct {
	Username  types.String `tfsdk:"username"`
	Host      types.String `tfsdk:"host"`
	Port      types.Int32  `tfsdk:"port"`
	IsEnabled types.Bool   `tfsdk:"is_enabled"`
}

func (m *AgencyHostingWebsiteRemoteAccessSFTPModel) Merge(item client.AgencyHostingV1WebsitesWebsiteSftpDetailsResource) {
	m.Username = types.StringPointerValue(item.Username)
	m.Host = types.StringPointerValue(item.Host)
	m.Port = int32Value(item.Port)
	m.IsEnabled = types.BoolPointerValue(item.IsEnabled)
}

type AgencyHostingWebsiteRemoteAccessModel struct {
	Mode types.String                              `tfsdk:"mode"`
	SSH  AgencyHostingWebsiteRemoteAccessSSHModel  `tfsdk:"ssh"`
	SFTP AgencyHostingWebsiteRemoteAccessSFTPModel `tfsdk:"sftp"`
}

func (m *AgencyHostingWebsiteRemoteAccessModel) Merge(item client.AgencyHostingV1WebsitesWebsiteRemoteAccessResource) {
	m.Mode = types.StringPointerValue(item.Mode)

	if item.Ssh != nil {
		var ssh AgencyHostingWebsiteRemoteAccessSSHModel
		ssh.Merge(*item.Ssh)
		m.SSH = ssh
	}

	if item.Sftp != nil {
		var sftp AgencyHostingWebsiteRemoteAccessSFTPModel
		sftp.Merge(*item.Sftp)
		m.SFTP = sftp
	}
}

type AgencyHostingWebsiteServerModel struct {
	Hostname    types.String `tfsdk:"hostname"`
	CountryCode types.String `tfsdk:"country_code"`
}

func (m *AgencyHostingWebsiteServerModel) Merge(item client.AgencyHostingV1WebsitesWebsiteServerResource) {
	m.Hostname = types.StringPointerValue(item.Hostname)
	m.CountryCode = types.StringPointerValue(item.CountryCode)
}

type AgencyHostingWebsiteOrderPlanParametersModel struct {
	DiskQuotaBytes         types.Int64 `tfsdk:"disk_quota_bytes"`
	InodeQuota             types.Int64 `tfsdk:"inode_quota"`
	CPUCores               types.Int32 `tfsdk:"cpu_cores"`
	MemoryQuotaBytes       types.Int64 `tfsdk:"memory_quota_bytes"`
	DiskIOPSQuota          types.Int64 `tfsdk:"disk_iops_quota"`
	ProcessQuota           types.Int32 `tfsdk:"process_quota"`
	WebsiteQuota           types.Int32 `tfsdk:"website_quota"`
	MaxDatabasesPerWebsite types.Int32 `tfsdk:"max_databases_per_website"`
	IsCDNAvailable         types.Bool  `tfsdk:"is_cdn_available"`
}

func (m *AgencyHostingWebsiteOrderPlanParametersModel) Merge(item client.AgencyHostingV1WebsitesWebsiteOrderPlanParametersResource) {
	m.DiskQuotaBytes = int64Value(item.DiskQuotaBytes)
	m.InodeQuota = int64Value(item.InodeQuota)
	m.CPUCores = int32Value(item.CpuCores)
	m.MemoryQuotaBytes = int64Value(item.MemoryQuotaBytes)
	m.DiskIOPSQuota = int64Value(item.DiskIopsQuota)
	m.ProcessQuota = int32Value(item.ProcessQuota)
	m.WebsiteQuota = int32Value(item.WebsiteQuota)
	m.MaxDatabasesPerWebsite = int32Value(item.MaxDatabasesPerWebsite)
	m.IsCDNAvailable = types.BoolPointerValue(item.IsCdnAvailable)
}

type AgencyHostingWebsiteOrderPlanModel struct {
	Name       types.String                                 `tfsdk:"name"`
	Parameters AgencyHostingWebsiteOrderPlanParametersModel `tfsdk:"parameters"`
}

func (m *AgencyHostingWebsiteOrderPlanModel) Merge(item client.AgencyHostingV1WebsitesWebsiteOrderPlanResource) {
	m.Name = types.StringPointerValue(item.Name)
	if item.Parameters != nil {
		var pp AgencyHostingWebsiteOrderPlanParametersModel
		pp.Merge(*item.Parameters)
		m.Parameters = pp
	}
}

type AgencyHostingWebsiteOrderModel struct {
	ID        types.Int64                        `tfsdk:"id"`
	Status    types.String                       `tfsdk:"status"`
	CreatedAt timetypes.RFC3339                  `tfsdk:"created_at"`
	Plan      AgencyHostingWebsiteOrderPlanModel `tfsdk:"plan"`
}

func (m *AgencyHostingWebsiteOrderModel) Merge(item client.AgencyHostingV1WebsitesWebsiteOrderResource) {
	m.ID = int64Value(item.Id)
	m.Status = types.StringPointerValue(item.Status)
	m.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)

	if item.Plan != nil {
		var p AgencyHostingWebsiteOrderPlanModel
		p.Merge(*item.Plan)
		m.Plan = p
	}
}

type AgencyHostingWebsiteUserModel struct {
	Username types.String `tfsdk:"username"`
	State    types.String `tfsdk:"state"`
}

func (m *AgencyHostingWebsiteUserModel) Merge(item client.AgencyHostingV1WebsitesWebsiteUserResource) {
	m.Username = types.StringPointerValue(item.Username)
	m.State = types.StringPointerValue(item.State)
}

type AgencyHostingWebsiteStagingRootModel struct {
	UID types.String `tfsdk:"uid"`
}

func (m *AgencyHostingWebsiteStagingRootModel) Merge(item client.AgencyHostingV1WebsitesWebsiteStagingRootResource) {
	m.UID = types.StringPointerValue(item.Uid)
}

type AgencyHostingWebsiteModel struct {
	UID           types.String                            `tfsdk:"uid"`
	IPv4          iptypes.IPv4Address                     `tfsdk:"ipv4"`
	Flavor        types.String                            `tfsdk:"flavor"`
	Type          types.String                            `tfsdk:"type"`
	Description   types.String                            `tfsdk:"description"`
	State         types.String                            `tfsdk:"state"`
	CreatedAt     timetypes.RFC3339                       `tfsdk:"created_at"`
	Domains       []AgencyHostingWebsiteDomainModel       `tfsdk:"domains"`
	PreviewDomain *AgencyHostingWebsitePreviewDomainModel `tfsdk:"preview_domain"`
	Settings      *AgencyHostingWebsiteSettingsModel      `tfsdk:"settings"`
	Wordpress     *AgencyHostingWebsiteWordpressModel     `tfsdk:"wordpress"`
	RemoteAccess  *AgencyHostingWebsiteRemoteAccessModel  `tfsdk:"remote_access"`
	Server        *AgencyHostingWebsiteServerModel        `tfsdk:"server"`
	Order         *AgencyHostingWebsiteOrderModel         `tfsdk:"order"`
	User          *AgencyHostingWebsiteUserModel          `tfsdk:"user"`
	StagingRoot   *AgencyHostingWebsiteStagingRootModel   `tfsdk:"staging_root"`
}

func (m *AgencyHostingWebsiteModel) Merge(item client.AgencyHostingV1WebsitesWebsiteResource) {
	m.UID = types.StringPointerValue(item.Uid)
	m.Flavor = types.StringPointerValue(item.Flavor)
	m.IPv4 = iptypes.NewIPv4AddressPointerValue(item.Ipv4)
	m.Flavor = types.StringPointerValue(item.Flavor)
	m.Type = types.StringPointerValue(item.Type)
	m.Description = types.StringPointerValue(item.Description)
	m.State = types.StringPointerValue((*string)(item.State))
	m.CreatedAt = timetypes.NewRFC3339TimePointerValue(item.CreatedAt)

	if item.Domains != nil {
		m.Domains = make([]AgencyHostingWebsiteDomainModel, len(*item.Domains))
		for i, domain := range *item.Domains {
			var d AgencyHostingWebsiteDomainModel
			d.Merge(domain)
			m.Domains[i] = d
		}
	}

	if item.PreviewDomain != nil {
		var d AgencyHostingWebsitePreviewDomainModel
		v, err := item.PreviewDomain.AsAgencyHostingV1WebsitesWebsitePreviewDomainResource()
		if err == nil {
			d.Merge(v)
		}
		m.PreviewDomain = &d
	}

	if item.Settings != nil {
		var d AgencyHostingWebsiteSettingsModel
		d.Merge(*item.Settings)
		m.Settings = &d
	}

	if item.Wordpress != nil {
		var d AgencyHostingWebsiteWordpressModel
		v, err := item.Wordpress.AsAgencyHostingV1WebsitesWordPressInstallResource()
		if err == nil {
			d.Merge(v)
		}
		m.Wordpress = &d
	}

	if item.RemoteAccess != nil {
		var d AgencyHostingWebsiteRemoteAccessModel
		d.Merge(*item.RemoteAccess)
		m.RemoteAccess = &d
	}

	if item.Server != nil {
		var s AgencyHostingWebsiteServerModel
		s.Merge(*item.Server)
		m.Server = &s
	}

	if item.Order != nil {
		var o AgencyHostingWebsiteOrderModel
		o.Merge(*item.Order)
		m.Order = &o
	}

	if item.User != nil {
		var u AgencyHostingWebsiteUserModel
		u.Merge(*item.User)
		m.User = &u
	}

	if item.StagingRoot != nil {
		var s AgencyHostingWebsiteStagingRootModel
		v, err := item.StagingRoot.AsAgencyHostingV1WebsitesWebsiteStagingRootResource()
		if err == nil {
			s.Merge(v)
		}
		m.StagingRoot = &s
	}
}
