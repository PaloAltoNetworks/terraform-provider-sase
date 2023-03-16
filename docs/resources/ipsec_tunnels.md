---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_ipsec_tunnels Resource - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_ipsec_tunnels (Resource)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `auto_key` (Attributes) (see [below for nested schema](#nestedatt--auto_key))
- `folder` (String) The folder of the entry
- `name` (String)

### Optional

- `anti_replay` (Boolean)
- `copy_tos` (Boolean)
- `enable_gre_encapsulation` (Boolean)
- `tunnel_monitor` (Attributes) (see [below for nested schema](#nestedatt--tunnel_monitor))

### Read-Only

- `id` (String) The object ID.
- `object_id` (String)

<a id="nestedatt--auto_key"></a>
### Nested Schema for `auto_key`

Required:

- `ike_gateway` (Attributes List) (see [below for nested schema](#nestedatt--auto_key--ike_gateway))
- `ipsec_crypto_profile` (String)

Optional:

- `proxy_id` (Attributes List) (see [below for nested schema](#nestedatt--auto_key--proxy_id))

<a id="nestedatt--auto_key--ike_gateway"></a>
### Nested Schema for `auto_key.ike_gateway`

Optional:

- `name` (String)


<a id="nestedatt--auto_key--proxy_id"></a>
### Nested Schema for `auto_key.proxy_id`

Required:

- `name` (String)

Optional:

- `local` (String)
- `protocol` (Attributes) (see [below for nested schema](#nestedatt--auto_key--proxy_id--protocol))
- `remote` (String)

<a id="nestedatt--auto_key--proxy_id--protocol"></a>
### Nested Schema for `auto_key.proxy_id.protocol`

Optional:

- `number` (Number)
- `tcp` (Attributes) (see [below for nested schema](#nestedatt--auto_key--proxy_id--protocol--tcp))
- `udp` (Attributes) (see [below for nested schema](#nestedatt--auto_key--proxy_id--protocol--udp))

<a id="nestedatt--auto_key--proxy_id--protocol--tcp"></a>
### Nested Schema for `auto_key.proxy_id.protocol.udp`

Optional:

- `local_port` (Number)
- `remote_port` (Number)


<a id="nestedatt--auto_key--proxy_id--protocol--udp"></a>
### Nested Schema for `auto_key.proxy_id.protocol.udp`

Optional:

- `local_port` (Number)
- `remote_port` (Number)





<a id="nestedatt--tunnel_monitor"></a>
### Nested Schema for `tunnel_monitor`

Required:

- `destination_ip` (String)

Optional:

- `enable` (Boolean)
- `proxy_id` (String)

