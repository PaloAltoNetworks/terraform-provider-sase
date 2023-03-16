---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_url_access_profiles_list Data Source - sase"
subcategory: ""
description: |-
  Retrieves a listing of config items.
---

# sase_url_access_profiles_list (Data Source)

Retrieves a listing of config items.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String) The folder of the entry

### Optional

- `limit` (Number) The max count in result entry (count per page)
- `name` (String) The name of the entry
- `offset` (Number) The offset of the result entry

### Read-Only

- `data` (Attributes List) (see [below for nested schema](#nestedatt--data))
- `id` (String) The object ID.
- `total` (Number)

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `alert` (List of String)
- `allow` (List of String)
- `block` (List of String)
- `continue` (List of String)
- `credential_enforcement` (Attributes) (see [below for nested schema](#nestedatt--data--credential_enforcement))
- `description` (String)
- `log_container_page_only` (Boolean)
- `log_http_hdr_referer` (Boolean)
- `log_http_hdr_user_agent` (Boolean)
- `log_http_hdr_xff` (Boolean)
- `mlav_category_exception` (List of String)
- `mlav_engine_urlbased_enabled` (Attributes List) (see [below for nested schema](#nestedatt--data--mlav_engine_urlbased_enabled))
- `name` (String)
- `object_id` (String)
- `safe_search_enforcement` (Boolean)

<a id="nestedatt--data--credential_enforcement"></a>
### Nested Schema for `data.credential_enforcement`

Read-Only:

- `alert` (List of String)
- `allow` (List of String)
- `block` (List of String)
- `continue` (List of String)
- `log_severity` (String)
- `mode` (Attributes) (see [below for nested schema](#nestedatt--data--credential_enforcement--mode))

<a id="nestedatt--data--credential_enforcement--mode"></a>
### Nested Schema for `data.credential_enforcement.mode`

Read-Only:

- `disabled` (Boolean)
- `domain_credentials` (Boolean)
- `group_mapping` (String)
- `ip_user` (Boolean)



<a id="nestedatt--data--mlav_engine_urlbased_enabled"></a>
### Nested Schema for `data.mlav_engine_urlbased_enabled`

Read-Only:

- `mlav_policy_action` (String)
- `name` (String)

