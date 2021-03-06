package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// IPv4 to Internet Gateway.
func TestAccAWSRoute_basic(t *testing.T) {
	var routeTable ec2.RouteTable
	resourceName := "aws_route.test"
	igwResourceName := "aws_internet_gateway.test"
	rtResourceName := "aws_route_table.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4InternetGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(rtResourceName, &routeTable),
					testAccCheckAWSRouteTableNumberOfRoutes(&routeTable, 2),
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "gateway_id", igwResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_disappears(t *testing.T) {
	resourceName := "aws_route.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4InternetGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceDisappears(testAccProvider, resourceAwsRoute(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSRoute_disappears_RouteTable(t *testing.T) {
	rtResourceName := "aws_route_table.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4InternetGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceDisappears(testAccProvider, resourceAwsRouteTable(), rtResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_EgressOnlyInternetGateway(t *testing.T) {
	resourceName := "aws_route.test"
	eoigwResourceName := "aws_egress_only_internet_gateway.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv6EgressOnlyInternetGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "egress_only_gateway_id", eoigwResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				// Verify that expanded form of the destination CIDR causes no diff.
				Config:   testAccAWSRouteConfigIpv6EgressOnlyInternetGateway(rName, "::0/0"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_InternetGateway(t *testing.T) {
	resourceName := "aws_route.test"
	igwResourceName := "aws_internet_gateway.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv6InternetGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "gateway_id", igwResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_Instance(t *testing.T) {
	resourceName := "aws_route.test"
	instanceResourceName := "aws_instance.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv6Instance(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", instanceResourceName, "id"),
					testAccCheckResourceAttrAccountID(resourceName, "instance_owner_id"),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "network_interface_id", instanceResourceName, "primary_network_interface_id"),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_NetworkInterface_Unattached(t *testing.T) {
	resourceName := "aws_route.test"
	eniResourceName := "aws_network_interface.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv6NetworkInterfaceUnattached(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "network_interface_id", eniResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateBlackhole),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_VpcPeeringConnection(t *testing.T) {
	resourceName := "aws_route.test"
	pcxResourceName := "aws_vpc_peering_connection.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv6VpcPeeringConnection(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_peering_connection_id", pcxResourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_VpnGateway(t *testing.T) {
	resourceName := "aws_route.test"
	vgwResourceName := "aws_vpn_gateway.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv6VpnGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "gateway_id", vgwResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv4_To_VpnGateway(t *testing.T) {
	resourceName := "aws_route.test"
	vgwResourceName := "aws_vpn_gateway.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4VpnGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "gateway_id", vgwResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv4_To_Instance(t *testing.T) {
	resourceName := "aws_route.test"
	instanceResourceName := "aws_instance.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4Instance(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", instanceResourceName, "id"),
					testAccCheckResourceAttrAccountID(resourceName, "instance_owner_id"),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "network_interface_id", instanceResourceName, "primary_network_interface_id"),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv4_To_NetworkInterface_Unattached(t *testing.T) {
	resourceName := "aws_route.test"
	eniResourceName := "aws_network_interface.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4NetworkInterfaceUnattached(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "network_interface_id", eniResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateBlackhole),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv4_To_NetworkInterface_Attached(t *testing.T) {
	resourceName := "aws_route.test"
	eniResourceName := "aws_network_interface.test"
	instanceResourceName := "aws_instance.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4NetworkInterfaceAttached(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", instanceResourceName, "id"),
					testAccCheckResourceAttrAccountID(resourceName, "instance_owner_id"),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "network_interface_id", eniResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

/*
func TestAccAWSRoute_IPv4_To_NetworkInterface_TwoAttachments(t *testing.T) {
	resourceName := "aws_route.test"
	eni1ResourceName := "aws_network_interface.test1"
	eni2ResourceName := "aws_network_interface.test2"
	instanceResourceName := "aws_instance.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4NetworkInterfaceTwoAttachments(rName, destinationCidr, eni1ResourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", instanceResourceName, "id"),
					testAccCheckResourceAttrAccountID(resourceName, "instance_owner_id"),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "network_interface_id", eni1ResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				Config: testAccAWSRouteConfigIpv4NetworkInterfaceTwoAttachments(rName, destinationCidr, eni2ResourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", instanceResourceName, "id"),
					testAccCheckResourceAttrAccountID(resourceName, "instance_owner_id"),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "network_interface_id", eni2ResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}
*/

func TestAccAWSRoute_IPv4_To_VpcPeeringConnection(t *testing.T) {
	resourceName := "aws_route.test"
	pcxResourceName := "aws_vpc_peering_connection.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4VpcPeeringConnection(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_peering_connection_id", pcxResourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv4_To_NatGateway(t *testing.T) {
	resourceName := "aws_route.test"
	ngwResourceName := "aws_nat_gateway.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4NatGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "nat_gateway_id", ngwResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_DoesNotCrashWithVpcEndpoint(t *testing.T) {
	var routeTable ec2.RouteTable
	resourceName := "aws_route.test"
	rtResourceName := "aws_route_table.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigWithVpcEndpoint(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableExists(rtResourceName, &routeTable),
					testAccCheckAWSRouteTableNumberOfRoutes(&routeTable, 3),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv4_To_TransitGateway(t *testing.T) {
	resourceName := "aws_route.test"
	tgwResourceName := "aws_ec2_transit_gateway.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv4TransitGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_id", tgwResourceName, "id"),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_TransitGateway(t *testing.T) {
	resourceName := "aws_route.test"
	tgwResourceName := "aws_ec2_transit_gateway.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigIpv6TransitGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_id", tgwResourceName, "id"),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv4_To_LocalGateway(t *testing.T) {
	resourceName := "aws_route.test"
	localGatewayDataSourceName := "data.aws_ec2_local_gateway.first"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "172.16.1.0/24"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSOutpostsOutposts(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteResourceConfigIpv4LocalGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "local_gateway_id", localGatewayDataSourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv6_To_LocalGateway(t *testing.T) {
	resourceName := "aws_route.test"
	localGatewayDataSourceName := "data.aws_ec2_local_gateway.first"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "2002:bc9:1234:1a00::/56"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSOutpostsOutposts(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteResourceConfigIpv6LocalGateway(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "local_gateway_id", localGatewayDataSourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "vpc_endpoint_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_ConditionalCidrBlock(t *testing.T) {
	resourceName := "aws_route.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	destinationCidr := "10.2.0.0/16"
	destinationIpv6Cidr := "::/0"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteConfigConditionalIpv4Ipv6(rName, destinationCidr, destinationIpv6Cidr, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
				),
			},
			{
				Config: testAccAWSRouteConfigConditionalIpv4Ipv6(rName, destinationCidr, destinationIpv6Cidr, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", destinationIpv6Cidr),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSRoute_IPv4_To_VpcEndpoint(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route.test"
	vpcEndpointResourceName := "aws_vpc_endpoint.test"
	destinationCidr := "172.16.1.0/24"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckElbv2GatewayLoadBalancer(t) },
		ErrorCheck:   testAccErrorCheck(t, ec2.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSRouteResourceConfigIpv4VpcEndpoint(rName, destinationCidr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destination_cidr_block", destinationCidr),
					resource.TestCheckResourceAttr(resourceName, "destination_ipv6_cidr_block", ""),
					resource.TestCheckResourceAttr(resourceName, "destination_prefix_list_id", ""),
					resource.TestCheckResourceAttr(resourceName, "egress_only_gateway_id", ""),
					// TODO resource.TestCheckResourceAttr(resourceName, "gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_id", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_owner_id", ""),
					resource.TestCheckResourceAttr(resourceName, "local_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_id", ""),
					resource.TestCheckResourceAttr(resourceName, "network_interface_id", ""),
					resource.TestCheckResourceAttr(resourceName, "origin", ec2.RouteOriginCreateRoute),
					resource.TestCheckResourceAttr(resourceName, "state", ec2.RouteStateActive),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_id", ""),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_endpoint_id", vpcEndpointResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "vpc_peering_connection_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccAWSRouteImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAWSRouteDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_route" {
			continue
		}

		conn := testAccProvider.Meta().(*AWSClient).ec2conn
		route, err := resourceAwsRouteFindRoute(
			conn,
			rs.Primary.Attributes["route_table_id"],
			rs.Primary.Attributes["destination_cidr_block"],
			rs.Primary.Attributes["destination_ipv6_cidr_block"],
		)

		if route == nil && err == nil {
			return nil
		}
	}

	return nil
}

func testAccAWSRouteImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		destination := rs.Primary.Attributes["destination_cidr_block"]
		if v, ok := rs.Primary.Attributes["destination_ipv6_cidr_block"]; ok && v != "" {
			destination = v
		}

		return fmt.Sprintf("%s_%s", rs.Primary.Attributes["route_table_id"], destination), nil
	}
}

func testAccAWSRouteConfigIpv4InternetGateway(rName, destinationCidr string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = %[2]q
  gateway_id             = aws_internet_gateway.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigIpv6InternetGateway(rName, destinationCidr string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_egress_only_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_ipv6_cidr_block = %[2]q
  gateway_id                  = aws_internet_gateway.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigIpv6NetworkInterfaceUnattached(rName, destinationCidr string) string {
	return composeConfig(
		testAccAvailableAZsNoOptInConfig(),
		fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block        = "10.1.1.0/24"
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.available.names[0]
  ipv6_cidr_block   = cidrsubnet(aws_vpc.test.ipv6_cidr_block, 8, 1)

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_interface" "test" {
  subnet_id = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_ipv6_cidr_block = %[2]q
  network_interface_id        = aws_network_interface.test.id
}
`, rName, destinationCidr))
}

func testAccAWSRouteConfigIpv6Instance(rName, destinationCidr string) string {
	return composeConfig(
		testAccLatestAmazonNatInstanceAmiConfig(),
		testAccAvailableAZsNoOptInConfig(),
		testAccAvailableEc2InstanceTypeForAvailabilityZone("data.aws_availability_zones.available.names[0]", "t3.micro", "t2.micro"),
		fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block        = "10.1.1.0/24"
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.available.names[0]
  ipv6_cidr_block   = cidrsubnet(aws_vpc.test.ipv6_cidr_block, 8, 1)

  tags = {
    Name = %[1]q
  }
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-nat-instance.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type
  subnet_id     = aws_subnet.test.id

  ipv6_address_count = 1

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_ipv6_cidr_block = %[2]q
  instance_id                 = aws_instance.test.id
}
`, rName, destinationCidr))
}

func testAccAWSRouteConfigIpv6VpcPeeringConnection(rName, destinationCidr string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc" "target" {
  cidr_block                       = "10.0.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc_peering_connection" "test" {
  vpc_id      = aws_vpc.test.id
  peer_vpc_id = aws_vpc.target.id
  auto_accept = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_ipv6_cidr_block = %[2]q
  vpc_peering_connection_id   = aws_vpc_peering_connection.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigIpv6EgressOnlyInternetGateway(rName, destinationCidr string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_egress_only_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_ipv6_cidr_block = %[2]q
  egress_only_gateway_id      = aws_egress_only_internet_gateway.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigWithVpcEndpoint(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = "10.3.0.0/16"
  gateway_id             = aws_internet_gateway.test.id

  # Forcing endpoint to create before route - without this the crash is a race.
  depends_on = [aws_vpc_endpoint.test]
}

data "aws_region" "current" {}

resource "aws_vpc_endpoint" "test" {
  vpc_id          = aws_vpc.test.id
  service_name    = "com.amazonaws.${data.aws_region.current.name}.s3"
  route_table_ids = [aws_route_table.test.id]
}
`, rName)
}

func testAccAWSRouteConfigIpv4TransitGateway(rName, destinationCidr string) string {
	return composeConfig(
		testAccAvailableAZsNoOptInDefaultExcludeConfig(),
		fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.1.1.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  destination_cidr_block = %[2]q
  route_table_id         = aws_route_table.test.id
  transit_gateway_id     = aws_ec2_transit_gateway_vpc_attachment.test.transit_gateway_id
}
`, rName, destinationCidr))
}

func testAccAWSRouteConfigIpv6TransitGateway(rName, destinationCidr string) string {
	return composeConfig(
		testAccAvailableAZsNoOptInDefaultExcludeConfig(),
		fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.1.1.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  destination_ipv6_cidr_block = %[2]q
  route_table_id              = aws_route_table.test.id
  transit_gateway_id          = aws_ec2_transit_gateway_vpc_attachment.test.transit_gateway_id
}
`, rName, destinationCidr))
}

func testAccAWSRouteConfigConditionalIpv4Ipv6(rName, destinationCidr, destinationIpv6Cidr string, ipv6Route bool) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

locals {
  ipv6             = %[4]t
  destination      = %[2]q
  destination_ipv6 = %[3]q
}

resource "aws_route" "test" {
  route_table_id = aws_route_table.test.id
  gateway_id     = aws_internet_gateway.test.id

  destination_cidr_block      = local.ipv6 ? "" : local.destination
  destination_ipv6_cidr_block = local.ipv6 ? local.destination_ipv6 : ""
}
`, rName, destinationCidr, destinationIpv6Cidr, ipv6Route)
}

func testAccAWSRouteConfigIpv4Instance(rName, destinationCidr string) string {
	return composeConfig(
		testAccLatestAmazonNatInstanceAmiConfig(),
		testAccAvailableAZsNoOptInConfig(),
		testAccAvailableEc2InstanceTypeForAvailabilityZone("data.aws_availability_zones.available.names[0]", "t3.micro", "t2.micro"),
		fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block        = "10.1.1.0/24"
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = %[1]q
  }
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-nat-instance.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type
  subnet_id     = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = %[2]q
  instance_id            = aws_instance.test.id
}
`, rName, destinationCidr))
}

func testAccAWSRouteConfigIpv4NetworkInterfaceUnattached(rName, destinationCidr string) string {
	return composeConfig(
		testAccAvailableAZsNoOptInConfig(),
		fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block        = "10.1.1.0/24"
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_interface" "test" {
  subnet_id = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = %[2]q
  network_interface_id   = aws_network_interface.test.id
}
`, rName, destinationCidr))
}

func testAccAWSRouteResourceConfigIpv4LocalGateway(rName, destinationCidr string) string {
	return fmt.Sprintf(`
data "aws_ec2_local_gateways" "all" {}

data "aws_ec2_local_gateway" "first" {
  id = tolist(data.aws_ec2_local_gateways.all.ids)[0]
}

data "aws_ec2_local_gateway_route_tables" "all" {}

data "aws_ec2_local_gateway_route_table" "first" {
  local_gateway_route_table_id = tolist(data.aws_ec2_local_gateway_route_tables.all.ids)[0]
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_local_gateway_route_table_vpc_association" "example" {
  local_gateway_route_table_id = data.aws_ec2_local_gateway_route_table.first.id
  vpc_id                       = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  depends_on = [aws_ec2_local_gateway_route_table_vpc_association.example]
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = %[2]q
  local_gateway_id       = data.aws_ec2_local_gateway.first.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteResourceConfigIpv6LocalGateway(rName, destinationCidr string) string {
	return fmt.Sprintf(`
data "aws_ec2_local_gateways" "all" {}

data "aws_ec2_local_gateway" "first" {
  id = tolist(data.aws_ec2_local_gateways.all.ids)[0]
}

data "aws_ec2_local_gateway_route_tables" "all" {}

data "aws_ec2_local_gateway_route_table" "first" {
  local_gateway_route_table_id = tolist(data.aws_ec2_local_gateway_route_tables.all.ids)[0]
}

resource "aws_vpc" "test" {
  cidr_block                       = "10.0.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_local_gateway_route_table_vpc_association" "example" {
  local_gateway_route_table_id = data.aws_ec2_local_gateway_route_table.first.id
  vpc_id                       = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  depends_on = [aws_ec2_local_gateway_route_table_vpc_association.example]
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_ipv6_cidr_block = %[2]q
  local_gateway_id            = data.aws_ec2_local_gateway.first.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigIpv4NetworkInterfaceAttached(rName, destinationCidr string) string {
	return composeConfig(
		testAccLatestAmazonNatInstanceAmiConfig(),
		testAccAvailableAZsNoOptInConfig(),
		testAccAvailableEc2InstanceTypeForAvailabilityZone("data.aws_availability_zones.available.names[0]", "t3.micro", "t2.micro"),
		fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block        = "10.1.1.0/24"
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_interface" "test" {
  subnet_id = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-nat-instance.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type

  network_interface {
    device_index         = 0
    network_interface_id = aws_network_interface.test.id
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = %[2]q
  network_interface_id   = aws_network_interface.test.id

  # Wait for the ENI attachment.
  depends_on = [aws_instance.test]
}
`, rName, destinationCidr))
}

/*
func testAccAWSRouteConfigIpv4NetworkInterfaceTwoAttachments(rName, destinationCidr, targetResourceName string) string {
	return composeConfig(
		testAccLatestAmazonNatInstanceAmiConfig(),
		testAccAvailableAZsNoOptInConfig(),
		testAccAvailableEc2InstanceTypeForAvailabilityZone("data.aws_availability_zones.available.names[0]", "t3.micro", "t2.micro"),
		fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block        = "10.1.1.0/24"
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_interface" "test1" {
  subnet_id = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_interface" "test2" {
  subnet_id = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-nat-instance.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type

  network_interface {
    device_index         = 0
    network_interface_id = aws_network_interface.test1.id
  }

  network_interface {
    device_index         = 1
    network_interface_id = aws_network_interface.test2.id
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = %[2]q
  network_interface_id   = %[3]s.id

  # Wait for the ENI attachment.
  depends_on = [aws_instance.test]
}
`, rName, destinationCidr, targetResourceName))
}
*/

func testAccAWSRouteConfigIpv4VpcPeeringConnection(rName, destinationCidr string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc" "target" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc_peering_connection" "test" {
  vpc_id      = aws_vpc.test.id
  peer_vpc_id = aws_vpc.target.id
  auto_accept = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id            = aws_route_table.test.id
  destination_cidr_block    = %[2]q
  vpc_peering_connection_id = aws_vpc_peering_connection.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigIpv4NatGateway(rName, destinationCidr string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  cidr_block = "10.1.1.0/24"
  vpc_id     = aws_vpc.test.id

  map_public_ip_on_launch = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_eip" "test" {
  vpc = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_nat_gateway" "test" {
  allocation_id = aws_eip.test.id
  subnet_id     = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }

  depends_on = [aws_internet_gateway.test]
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = %[2]q
  nat_gateway_id         = aws_nat_gateway.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigIpv4VpnGateway(rName, destinationCidr string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpn_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = %[2]q
  gateway_id             = aws_vpn_gateway.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteConfigIpv6VpnGateway(rName, destinationCidr string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block                       = "10.1.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpn_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id              = aws_route_table.test.id
  destination_ipv6_cidr_block = %[2]q
  gateway_id                  = aws_vpn_gateway.test.id
}
`, rName, destinationCidr)
}

func testAccAWSRouteResourceConfigIpv4VpcEndpoint(rName, destinationCidr string) string {
	return composeConfig(
		testAccAvailableAZsNoOptInConfig(),
		fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_vpc" "test" {
  cidr_block = "10.10.10.0/25"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = cidrsubnet(aws_vpc.test.cidr_block, 2, 0)
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_lb" "test" {
  load_balancer_type = "gateway"
  name               = %[1]q

  subnet_mapping {
    subnet_id = aws_subnet.test.id
  }
}

resource "aws_vpc_endpoint_service" "test" {
  acceptance_required        = false
  allowed_principals         = [data.aws_caller_identity.current.arn]
  gateway_load_balancer_arns = [aws_lb.test.arn]

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc_endpoint" "test" {
  service_name      = aws_vpc_endpoint_service.test.service_name
  subnet_ids        = [aws_subnet.test.id]
  vpc_endpoint_type = aws_vpc_endpoint_service.test.service_type
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = %[2]q
  vpc_endpoint_id        = aws_vpc_endpoint.test.id
}
`, rName, destinationCidr))
}
