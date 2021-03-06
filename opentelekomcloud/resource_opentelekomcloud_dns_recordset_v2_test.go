package opentelekomcloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/dns/v2/recordsets"
)

func randomZoneName() string {
	// TODO: why does back-end convert name to lowercase?
	return fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))
}

func TestAccDNSV2RecordSet_basic(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_basic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("opentelekomcloud_dns_recordset_v2.recordset_1", &recordset),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "description", "a record set"),
				),
			},
			{
				Config: testAccDNSV2RecordSet_update(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_dns_recordset_v2.recordset_1", "name", zoneName),
					resource.TestCheckResourceAttr("opentelekomcloud_dns_recordset_v2.recordset_1", "ttl", "6000"),
					resource.TestCheckResourceAttr("opentelekomcloud_dns_recordset_v2.recordset_1", "type", "A"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "description", "an updated record set"),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_updateTTL(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()
	newTTL := 1500

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_basic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("opentelekomcloud_dns_recordset_v2.recordset_1", &recordset),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "description", "a record set"),
				),
			},
			{
				Config: testAccDNSV2RecordSet_updateTTL(zoneName, newTTL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("opentelekomcloud_dns_recordset_v2.recordset_1", "name", zoneName),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_readTTL(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_readTTL(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("opentelekomcloud_dns_recordset_v2.recordset_1", &recordset),
					resource.TestMatchResourceAttr(
						"opentelekomcloud_dns_recordset_v2.recordset_1", "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_timeout(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_timeout(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("opentelekomcloud_dns_recordset_v2.recordset_1", &recordset),
				),
			},
		},
	})
}

func testAccCheckDNSV2RecordSetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dnsClient, err := config.dnsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_dns_recordset_v2" {
			continue
		}

		zoneID, recordsetID, err := parseDNSV2RecordSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err == nil {
			return fmt.Errorf("Record set still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2RecordSetExists(n string, recordset *recordsets.RecordSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		dnsClient, err := config.dnsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud DNS client: %s", err)
		}

		zoneID, recordsetID, err := parseDNSV2RecordSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err != nil {
			return err
		}

		if found.ID != recordsetID {
			return fmt.Errorf("Record set not found")
		}

		*recordset = *found

		return nil
	}
}

func testAccDNSV2RecordSet_basic(zoneName string) string {
	return fmt.Sprintf(`
		resource "opentelekomcloud_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "a zone"
			ttl = 6000
		}

		resource "opentelekomcloud_dns_recordset_v2" "recordset_1" {
			zone_id = "${opentelekomcloud_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "A"
			description = "a record set"
			ttl = 3000
			records = ["10.1.0.0"]
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSet_update(zoneName string) string {
	return fmt.Sprintf(`
		resource "opentelekomcloud_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "an updated zone"
			ttl = 6000
		}

		resource "opentelekomcloud_dns_recordset_v2" "recordset_1" {
			zone_id = "${opentelekomcloud_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "A"
			description = "an updated record set"
			ttl = 6000
			records = ["10.1.0.1"]
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSet_updateTTL(zoneName string, ttl int) string {
	return fmt.Sprintf(`
		resource "opentelekomcloud_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "a zone"
			ttl = 6000
		}

		resource "opentelekomcloud_dns_recordset_v2" "recordset_1" {
			zone_id = "${opentelekomcloud_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "A"
			description = "a record set"
			ttl = %d
			records = ["10.1.0.0"]
		}
	`, zoneName, zoneName, ttl)
}

func testAccDNSV2RecordSet_readTTL(zoneName string) string {
	return fmt.Sprintf(`
		resource "opentelekomcloud_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "an updated zone"
			ttl = 6000
		}

		resource "opentelekomcloud_dns_recordset_v2" "recordset_1" {
			zone_id = "${opentelekomcloud_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "A"
			records = ["10.1.0.2"]
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSet_timeout(zoneName string) string {
	return fmt.Sprintf(`
		resource "opentelekomcloud_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "an updated zone"
			ttl = 6000
		}

		resource "opentelekomcloud_dns_recordset_v2" "recordset_1" {
			zone_id = "${opentelekomcloud_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "A"
			ttl = 3000
			records = ["10.1.0.3", "10.1.0.2"]

			timeouts {
				create = "5m"
				update = "5m"
				delete = "5m"
			}
		}
	`, zoneName, zoneName)
}
