---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sase_objects_external_dynamic_lists_list Data Source - sase"
subcategory: ""
description: |-
  Retrieves a listing of config items.
---

# sase_objects_external_dynamic_lists_list (Data Source)

Retrieves a listing of config items.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `folder` (String) The folder of the entry. Value must be one of: `"Shared"`, `"Mobile Users"`, `"Remote Networks"`, `"Service Connections"`, `"Mobile Users Container"`, `"Mobile Users Explicit Proxy"`.

### Optional

- `limit` (Number) The max count in result entry (count per page).
- `name` (String) The name of the entry.
- `offset` (Number) The offset of the result entry.

### Read-Only

- `data` (Attributes List) The `data` parameter. (see [below for nested schema](#nestedatt--data))
- `id` (String) The object ID.
- `total` (Number) The `total` parameter.

<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `name` (String) The `name` parameter.
- `object_id` (String) The `object_id` parameter.
- `type` (Attributes) The `type` parameter. (see [below for nested schema](#nestedatt--data--type))

<a id="nestedatt--data--type"></a>
### Nested Schema for `data.type`

Read-Only:

- `domain` (Attributes) The `domain` parameter. (see [below for nested schema](#nestedatt--data--type--domain))
- `imei` (Attributes) The `imei` parameter. (see [below for nested schema](#nestedatt--data--type--imei))
- `imsi` (Attributes) The `imsi` parameter. (see [below for nested schema](#nestedatt--data--type--imsi))
- `ip` (Attributes) The `ip` parameter. (see [below for nested schema](#nestedatt--data--type--ip))
- `predefined_ip` (Attributes) The `predefined_ip` parameter. (see [below for nested schema](#nestedatt--data--type--predefined_ip))
- `predefined_url` (Attributes) The `predefined_url` parameter. (see [below for nested schema](#nestedatt--data--type--predefined_url))
- `url` (Attributes) The `url` parameter. (see [below for nested schema](#nestedatt--data--type--url))

<a id="nestedatt--data--type--domain"></a>
### Nested Schema for `data.type.domain`

Read-Only:

- `auth` (Attributes) The `auth` parameter. (see [below for nested schema](#nestedatt--data--type--domain--auth))
- `certificate_profile` (String) The `certificate_profile` parameter.
- `description` (String) The `description` parameter.
- `exception_list` (List of String) The `exception_list` parameter.
- `expand_domain` (Boolean) The `expand_domain` parameter.
- `recurring` (Attributes) The `recurring` parameter. (see [below for nested schema](#nestedatt--data--type--domain--recurring))
- `url` (String) The `url` parameter.

<a id="nestedatt--data--type--domain--auth"></a>
### Nested Schema for `data.type.domain.url`

Read-Only:

- `password` (String) The `password` parameter.
- `username` (String) The `username` parameter.


<a id="nestedatt--data--type--domain--recurring"></a>
### Nested Schema for `data.type.domain.url`

Read-Only:

- `daily` (Attributes) The `daily` parameter. (see [below for nested schema](#nestedatt--data--type--domain--url--daily))
- `five_minute` (Boolean) The `five_minute` parameter.
- `hourly` (Boolean) The `hourly` parameter.
- `monthly` (Attributes) The `monthly` parameter. (see [below for nested schema](#nestedatt--data--type--domain--url--monthly))
- `weekly` (Attributes) The `weekly` parameter. (see [below for nested schema](#nestedatt--data--type--domain--url--weekly))

<a id="nestedatt--data--type--domain--url--daily"></a>
### Nested Schema for `data.type.domain.url.daily`

Read-Only:

- `at` (String) The `at` parameter.


<a id="nestedatt--data--type--domain--url--monthly"></a>
### Nested Schema for `data.type.domain.url.monthly`

Read-Only:

- `at` (String) The `at` parameter.
- `day_of_month` (Number) The `day_of_month` parameter.


<a id="nestedatt--data--type--domain--url--weekly"></a>
### Nested Schema for `data.type.domain.url.weekly`

Read-Only:

- `at` (String) The `at` parameter.
- `day_of_week` (String) The `day_of_week` parameter.




<a id="nestedatt--data--type--imei"></a>
### Nested Schema for `data.type.imei`

Read-Only:

- `auth` (Attributes) The `auth` parameter. (see [below for nested schema](#nestedatt--data--type--imei--auth))
- `certificate_profile` (String) The `certificate_profile` parameter.
- `description` (String) The `description` parameter.
- `exception_list` (List of String) The `exception_list` parameter.
- `recurring` (Attributes) The `recurring` parameter. (see [below for nested schema](#nestedatt--data--type--imei--recurring))
- `url` (String) The `url` parameter.

<a id="nestedatt--data--type--imei--auth"></a>
### Nested Schema for `data.type.imei.url`

Read-Only:

- `password` (String) The `password` parameter.
- `username` (String) The `username` parameter.


<a id="nestedatt--data--type--imei--recurring"></a>
### Nested Schema for `data.type.imei.url`

Read-Only:

- `daily` (Attributes) The `daily` parameter. (see [below for nested schema](#nestedatt--data--type--imei--url--daily))
- `five_minute` (Boolean) The `five_minute` parameter.
- `hourly` (Boolean) The `hourly` parameter.
- `monthly` (Attributes) The `monthly` parameter. (see [below for nested schema](#nestedatt--data--type--imei--url--monthly))
- `weekly` (Attributes) The `weekly` parameter. (see [below for nested schema](#nestedatt--data--type--imei--url--weekly))

<a id="nestedatt--data--type--imei--url--daily"></a>
### Nested Schema for `data.type.imei.url.daily`

Read-Only:

- `at` (String) The `at` parameter.


<a id="nestedatt--data--type--imei--url--monthly"></a>
### Nested Schema for `data.type.imei.url.monthly`

Read-Only:

- `at` (String) The `at` parameter.
- `day_of_month` (Number) The `day_of_month` parameter.


<a id="nestedatt--data--type--imei--url--weekly"></a>
### Nested Schema for `data.type.imei.url.weekly`

Read-Only:

- `at` (String) The `at` parameter.
- `day_of_week` (String) The `day_of_week` parameter.




<a id="nestedatt--data--type--imsi"></a>
### Nested Schema for `data.type.imsi`

Read-Only:

- `auth` (Attributes) The `auth` parameter. (see [below for nested schema](#nestedatt--data--type--imsi--auth))
- `certificate_profile` (String) The `certificate_profile` parameter.
- `description` (String) The `description` parameter.
- `exception_list` (List of String) The `exception_list` parameter.
- `recurring` (Attributes) The `recurring` parameter. (see [below for nested schema](#nestedatt--data--type--imsi--recurring))
- `url` (String) The `url` parameter.

<a id="nestedatt--data--type--imsi--auth"></a>
### Nested Schema for `data.type.imsi.url`

Read-Only:

- `password` (String) The `password` parameter.
- `username` (String) The `username` parameter.


<a id="nestedatt--data--type--imsi--recurring"></a>
### Nested Schema for `data.type.imsi.url`

Read-Only:

- `daily` (Attributes) The `daily` parameter. (see [below for nested schema](#nestedatt--data--type--imsi--url--daily))
- `five_minute` (Boolean) The `five_minute` parameter.
- `hourly` (Boolean) The `hourly` parameter.
- `monthly` (Attributes) The `monthly` parameter. (see [below for nested schema](#nestedatt--data--type--imsi--url--monthly))
- `weekly` (Attributes) The `weekly` parameter. (see [below for nested schema](#nestedatt--data--type--imsi--url--weekly))

<a id="nestedatt--data--type--imsi--url--daily"></a>
### Nested Schema for `data.type.imsi.url.daily`

Read-Only:

- `at` (String) The `at` parameter.


<a id="nestedatt--data--type--imsi--url--monthly"></a>
### Nested Schema for `data.type.imsi.url.monthly`

Read-Only:

- `at` (String) The `at` parameter.
- `day_of_month` (Number) The `day_of_month` parameter.


<a id="nestedatt--data--type--imsi--url--weekly"></a>
### Nested Schema for `data.type.imsi.url.weekly`

Read-Only:

- `at` (String) The `at` parameter.
- `day_of_week` (String) The `day_of_week` parameter.




<a id="nestedatt--data--type--ip"></a>
### Nested Schema for `data.type.ip`

Read-Only:

- `auth` (Attributes) The `auth` parameter. (see [below for nested schema](#nestedatt--data--type--ip--auth))
- `certificate_profile` (String) The `certificate_profile` parameter.
- `description` (String) The `description` parameter.
- `exception_list` (List of String) The `exception_list` parameter.
- `recurring` (Attributes) The `recurring` parameter. (see [below for nested schema](#nestedatt--data--type--ip--recurring))
- `url` (String) The `url` parameter.

<a id="nestedatt--data--type--ip--auth"></a>
### Nested Schema for `data.type.ip.url`

Read-Only:

- `password` (String) The `password` parameter.
- `username` (String) The `username` parameter.


<a id="nestedatt--data--type--ip--recurring"></a>
### Nested Schema for `data.type.ip.url`

Read-Only:

- `daily` (Attributes) The `daily` parameter. (see [below for nested schema](#nestedatt--data--type--ip--url--daily))
- `five_minute` (Boolean) The `five_minute` parameter.
- `hourly` (Boolean) The `hourly` parameter.
- `monthly` (Attributes) The `monthly` parameter. (see [below for nested schema](#nestedatt--data--type--ip--url--monthly))
- `weekly` (Attributes) The `weekly` parameter. (see [below for nested schema](#nestedatt--data--type--ip--url--weekly))

<a id="nestedatt--data--type--ip--url--daily"></a>
### Nested Schema for `data.type.ip.url.daily`

Read-Only:

- `at` (String) The `at` parameter.


<a id="nestedatt--data--type--ip--url--monthly"></a>
### Nested Schema for `data.type.ip.url.monthly`

Read-Only:

- `at` (String) The `at` parameter.
- `day_of_month` (Number) The `day_of_month` parameter.


<a id="nestedatt--data--type--ip--url--weekly"></a>
### Nested Schema for `data.type.ip.url.weekly`

Read-Only:

- `at` (String) The `at` parameter.
- `day_of_week` (String) The `day_of_week` parameter.




<a id="nestedatt--data--type--predefined_ip"></a>
### Nested Schema for `data.type.predefined_ip`

Read-Only:

- `description` (String) The `description` parameter.
- `exception_list` (List of String) The `exception_list` parameter.
- `url` (String) The `url` parameter.


<a id="nestedatt--data--type--predefined_url"></a>
### Nested Schema for `data.type.predefined_url`

Read-Only:

- `description` (String) The `description` parameter.
- `exception_list` (List of String) The `exception_list` parameter.
- `url` (String) The `url` parameter.


<a id="nestedatt--data--type--url"></a>
### Nested Schema for `data.type.url`

Read-Only:

- `auth` (Attributes) The `auth` parameter. (see [below for nested schema](#nestedatt--data--type--url--auth))
- `certificate_profile` (String) The `certificate_profile` parameter.
- `description` (String) The `description` parameter.
- `exception_list` (List of String) The `exception_list` parameter.
- `recurring` (Attributes) The `recurring` parameter. (see [below for nested schema](#nestedatt--data--type--url--recurring))
- `url` (String) The `url` parameter.

<a id="nestedatt--data--type--url--auth"></a>
### Nested Schema for `data.type.url.url`

Read-Only:

- `password` (String) The `password` parameter.
- `username` (String) The `username` parameter.


<a id="nestedatt--data--type--url--recurring"></a>
### Nested Schema for `data.type.url.url`

Read-Only:

- `daily` (Attributes) The `daily` parameter. (see [below for nested schema](#nestedatt--data--type--url--url--daily))
- `five_minute` (Boolean) The `five_minute` parameter.
- `hourly` (Boolean) The `hourly` parameter.
- `monthly` (Attributes) The `monthly` parameter. (see [below for nested schema](#nestedatt--data--type--url--url--monthly))
- `weekly` (Attributes) The `weekly` parameter. (see [below for nested schema](#nestedatt--data--type--url--url--weekly))

<a id="nestedatt--data--type--url--url--daily"></a>
### Nested Schema for `data.type.url.url.daily`

Read-Only:

- `at` (String) The `at` parameter.


<a id="nestedatt--data--type--url--url--monthly"></a>
### Nested Schema for `data.type.url.url.monthly`

Read-Only:

- `at` (String) The `at` parameter.
- `day_of_month` (Number) The `day_of_month` parameter.


<a id="nestedatt--data--type--url--url--weekly"></a>
### Nested Schema for `data.type.url.url.weekly`

Read-Only:

- `at` (String) The `at` parameter.
- `day_of_week` (String) The `day_of_week` parameter.


