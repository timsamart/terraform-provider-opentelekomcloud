---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dns_zone_v2"
sidebar_current: "docs-opentelekomcloud-resource-dns-zone-v2"
description: |-
  Manages a DNS zone in the OpenTelekomCloud DNS Service
---

# opentelekomcloud\_dns\_zone_v2

Manages a DNS zone in the OpenTelekomCloud DNS Service.

## Example Usage

### Public Zone Configuration

```hcl
resource "opentelekomcloud_dns_zone_v2" "public_example_com" {
  name = "public.example.com."
  email = "public@example.com"
  description = "An example for public zone"
  ttl = 3000
  type = "public"
}
```

### Private Zone Configuration

```hcl
resource "opentelekomcloud_dns_zone_v2" "private_example_com" {
  name = "private.example.com."
  email = "private@example.com"
  description = "An example for private zone"
  ttl = 3000
  type = "private"
  router {
     router_id = "${var.vpc_id}"
     router_region = "${var.region}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the zone. Note the `.` at the end of the name.
  Changing this creates a new DNS zone.

* `email` - (Optional) The email contact for the zone record.

* `type` - (Optional) The type of zone. Can either be `public` or `private`.
  Changing this creates a new zone.

* `ttl` - (Optional) The time to live (TTL) of the zone.

* `description` - (Optional) A description of the zone.

* `router` (Optional) The Router(VPC) configuration for the private zone.
    it is required when type is `private`.

* `masters` - (Optional) An array of master DNS servers. 

* `value_specs` - (Optional) Map of additional options. Changing this creates a
  new zone.

The `router` block supports:

* `router_id` (Required) The Router(VPC) ID. which VPC network will assicate with.

* `router_region` (Required) The Region name for this private zone.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `email` - See Argument Reference above.
* `type` - See Argument Reference above.
* `ttl` - See Argument Reference above.
* `description` - See Argument Reference above.
* `masters` - See Argument Reference above.
* `value_specs` - See Argument Reference above.

## Import

This resource can be imported by specifying the zone ID:

```
$ terraform import opentelekomcloud_dns_zone_v2.zone_1 <zone_id>
```
