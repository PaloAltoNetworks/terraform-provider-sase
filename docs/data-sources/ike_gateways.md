---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_ike_gateways Data Source - sase"
subcategory: ""
description: |-
  Retrieves config for a specific item.
---

# sase_ike_gateways (Data Source)

Retrieves config for a specific item.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `object_id` (String) The uuid of the resource

### Read-Only

- `authentication` (Attributes) (see [below for nested schema](#nestedatt--authentication))
- `id` (String) The object ID.
- `local_id` (Attributes) (see [below for nested schema](#nestedatt--local_id))
- `name` (String)
- `peer_address` (Attributes) (see [below for nested schema](#nestedatt--peer_address))
- `peer_id` (Attributes) (see [below for nested schema](#nestedatt--peer_id))
- `protocol` (Attributes) (see [below for nested schema](#nestedatt--protocol))
- `protocol_common` (Attributes) (see [below for nested schema](#nestedatt--protocol_common))

<a id="nestedatt--authentication"></a>
### Nested Schema for `authentication`

Read-Only:

- `allow_id_payload_mismatch` (Boolean)
- `certificate_profile` (String)
- `local_certificate` (Attributes) (see [below for nested schema](#nestedatt--authentication--local_certificate))
- `pre_shared_key` (Attributes) (see [below for nested schema](#nestedatt--authentication--pre_shared_key))
- `strict_validation_revocation` (Boolean)
- `use_management_as_source` (Boolean)

<a id="nestedatt--authentication--local_certificate"></a>
### Nested Schema for `authentication.local_certificate`

Read-Only:

- `local_certificate_name` (String)


<a id="nestedatt--authentication--pre_shared_key"></a>
### Nested Schema for `authentication.pre_shared_key`

Read-Only:

- `key` (String)



<a id="nestedatt--local_id"></a>
### Nested Schema for `local_id`

Read-Only:

- `object_id` (String)
- `type` (String)


<a id="nestedatt--peer_address"></a>
### Nested Schema for `peer_address`

Read-Only:

- `dynamic_value` (Boolean)
- `fqdn` (String)
- `ip` (String)


<a id="nestedatt--peer_id"></a>
### Nested Schema for `peer_id`

Read-Only:

- `object_id` (String)
- `type` (String)


<a id="nestedatt--protocol"></a>
### Nested Schema for `protocol`

Read-Only:

- `ikev1` (Attributes) (see [below for nested schema](#nestedatt--protocol--ikev1))
- `ikev2` (Attributes) (see [below for nested schema](#nestedatt--protocol--ikev2))
- `version` (String)

<a id="nestedatt--protocol--ikev1"></a>
### Nested Schema for `protocol.ikev1`

Read-Only:

- `dpd` (Attributes) (see [below for nested schema](#nestedatt--protocol--ikev1--dpd))
- `ike_crypto_profile` (String)

<a id="nestedatt--protocol--ikev1--dpd"></a>
### Nested Schema for `protocol.ikev1.dpd`

Read-Only:

- `enable` (Boolean)



<a id="nestedatt--protocol--ikev2"></a>
### Nested Schema for `protocol.ikev2`

Read-Only:

- `dpd` (Attributes) (see [below for nested schema](#nestedatt--protocol--ikev2--dpd))
- `ike_crypto_profile` (String)

<a id="nestedatt--protocol--ikev2--dpd"></a>
### Nested Schema for `protocol.ikev2.dpd`

Read-Only:

- `enable` (Boolean)




<a id="nestedatt--protocol_common"></a>
### Nested Schema for `protocol_common`

Read-Only:

- `fragmentation` (Attributes) (see [below for nested schema](#nestedatt--protocol_common--fragmentation))
- `nat_traversal` (Attributes) (see [below for nested schema](#nestedatt--protocol_common--nat_traversal))
- `passive_mode` (Boolean)

<a id="nestedatt--protocol_common--fragmentation"></a>
### Nested Schema for `protocol_common.fragmentation`

Read-Only:

- `enable` (Boolean)


<a id="nestedatt--protocol_common--nat_traversal"></a>
### Nested Schema for `protocol_common.nat_traversal`

Read-Only:

- `enable` (Boolean)

