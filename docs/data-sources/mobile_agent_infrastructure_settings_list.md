---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_mobile_agent_infrastructure_settings_list Data Source - sase"
subcategory: ""
description: |-
  Retrieves a listing of config items.
---

# sase_mobile_agent_infrastructure_settings_list (Data Source)

Retrieves a listing of config items.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String)

### Read-Only

- `data` (Attributes List) (see [below for nested schema](#nestedatt--data))
- `id` (String) The object ID.
- `limit` (Number)
- `offset` (Number)
- `total` (Number)

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `dns_servers` (Attributes List) (see [below for nested schema](#nestedatt--data--dns_servers))
- `enable_wins` (Attributes) (see [below for nested schema](#nestedatt--data--enable_wins))
- `ip_pools` (Attributes List) (see [below for nested schema](#nestedatt--data--ip_pools))
- `ipv6` (Boolean)
- `name` (String)
- `portal_hostname` (Attributes) (see [below for nested schema](#nestedatt--data--portal_hostname))
- `region_ipv6` (Attributes) (see [below for nested schema](#nestedatt--data--region_ipv6))
- `udp_queries` (Attributes) (see [below for nested schema](#nestedatt--data--udp_queries))

<a id="nestedatt--data--dns_servers"></a>
### Nested Schema for `data.dns_servers`

Read-Only:

- `dns_suffix` (List of String)
- `internal_dns_match` (Attributes List) (see [below for nested schema](#nestedatt--data--dns_servers--internal_dns_match))
- `name` (String)
- `primary_public_dns` (Attributes) (see [below for nested schema](#nestedatt--data--dns_servers--primary_public_dns))
- `secondary_public_dns` (Attributes) (see [below for nested schema](#nestedatt--data--dns_servers--secondary_public_dns))

<a id="nestedatt--data--dns_servers--internal_dns_match"></a>
### Nested Schema for `data.dns_servers.internal_dns_match`

Read-Only:

- `domain_list` (List of String)
- `name` (String)
- `primary` (Attributes) (see [below for nested schema](#nestedatt--data--dns_servers--internal_dns_match--primary))
- `secondary` (Attributes) (see [below for nested schema](#nestedatt--data--dns_servers--internal_dns_match--secondary))

<a id="nestedatt--data--dns_servers--internal_dns_match--primary"></a>
### Nested Schema for `data.dns_servers.internal_dns_match.secondary`

Read-Only:

- `dns_server` (Boolean)
- `use_cloud_default` (Boolean)


<a id="nestedatt--data--dns_servers--internal_dns_match--secondary"></a>
### Nested Schema for `data.dns_servers.internal_dns_match.secondary`

Read-Only:

- `dns_server` (Boolean)
- `use_cloud_default` (Boolean)



<a id="nestedatt--data--dns_servers--primary_public_dns"></a>
### Nested Schema for `data.dns_servers.primary_public_dns`

Read-Only:

- `dns_server` (String)


<a id="nestedatt--data--dns_servers--secondary_public_dns"></a>
### Nested Schema for `data.dns_servers.secondary_public_dns`

Read-Only:

- `dns_server` (String)



<a id="nestedatt--data--enable_wins"></a>
### Nested Schema for `data.enable_wins`

Read-Only:

- `no` (Boolean)
- `yes` (Attributes) (see [below for nested schema](#nestedatt--data--enable_wins--yes))

<a id="nestedatt--data--enable_wins--yes"></a>
### Nested Schema for `data.enable_wins.yes`

Read-Only:

- `wins_servers` (Attributes List) (see [below for nested schema](#nestedatt--data--enable_wins--yes--wins_servers))

<a id="nestedatt--data--enable_wins--yes--wins_servers"></a>
### Nested Schema for `data.enable_wins.yes.wins_servers`

Read-Only:

- `name` (String)
- `primary` (String)
- `secondary` (String)




<a id="nestedatt--data--ip_pools"></a>
### Nested Schema for `data.ip_pools`

Read-Only:

- `ip_pool` (List of String)
- `name` (String)


<a id="nestedatt--data--portal_hostname"></a>
### Nested Schema for `data.portal_hostname`

Read-Only:

- `custom_domain` (Attributes) (see [below for nested schema](#nestedatt--data--portal_hostname--custom_domain))
- `default_domain` (Attributes) (see [below for nested schema](#nestedatt--data--portal_hostname--default_domain))

<a id="nestedatt--data--portal_hostname--custom_domain"></a>
### Nested Schema for `data.portal_hostname.custom_domain`

Read-Only:

- `cname` (String)
- `hostname` (String)
- `ssl_tls_service_profile` (String)


<a id="nestedatt--data--portal_hostname--default_domain"></a>
### Nested Schema for `data.portal_hostname.default_domain`

Read-Only:

- `hostname` (String)



<a id="nestedatt--data--region_ipv6"></a>
### Nested Schema for `data.region_ipv6`

Read-Only:

- `region` (Attributes List) (see [below for nested schema](#nestedatt--data--region_ipv6--region))

<a id="nestedatt--data--region_ipv6--region"></a>
### Nested Schema for `data.region_ipv6.region`

Read-Only:

- `locations` (List of String)
- `name` (String)



<a id="nestedatt--data--udp_queries"></a>
### Nested Schema for `data.udp_queries`

Read-Only:

- `retries` (Attributes) (see [below for nested schema](#nestedatt--data--udp_queries--retries))

<a id="nestedatt--data--udp_queries--retries"></a>
### Nested Schema for `data.udp_queries.retries`

Read-Only:

- `attempts` (Number)
- `interval` (Number)

