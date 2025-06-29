package dns

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/recordsets"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
)

func FetchZones(projectId string) []zones.Zone {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: os.Getenv("OS_AUTH_URL"),
		Username:         os.Getenv("OS_USERNAME"),
		Password:         os.Getenv("OS_PASSWORD"),
		DomainName:       os.Getenv("OS_USER_DOMAIN_NAME"),
		TenantName:       os.Getenv("OS_PROJECT_NAME"),
		AllowReauth:      true,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println("Failed to authenticate:", err)
		return nil
	}

	client, err := openstack.NewDNSV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create dns client:", err)
		return nil
	}

	allPages, err := zones.List(client, zones.ListOpts{}).AllPages()
	if err != nil {
		fmt.Println("Failed to list zones:", err)
		return nil
	}

	zoneList, err := zones.ExtractZones(allPages)
	if err != nil {
		fmt.Println("Failed to extract zones:", err)
		return nil
	}

	return zoneList
}

func FetchZoneByID(zoneId string) *zones.Zone {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: os.Getenv("OS_AUTH_URL"),
		Username:         os.Getenv("OS_USERNAME"),
		Password:         os.Getenv("OS_PASSWORD"),
		DomainName:       os.Getenv("OS_USER_DOMAIN_NAME"),
		TenantName:       os.Getenv("OS_PROJECT_NAME"),
		AllowReauth:      true,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println("Failed to authenticate:", err)
		return nil
	}

	client, err := openstack.NewDNSV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create dns client:", err)
		return nil
	}

	zone, err := zones.Get(client, zoneId).Extract()
	if err != nil {
		fmt.Println("Failed to get zone details:", err)
		return nil
	}

	return zone
}

func FetchZoneByName(zoneName, projectID string) *zones.Zone {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: os.Getenv("OS_AUTH_URL"),
		Username:         os.Getenv("OS_USERNAME"),
		Password:         os.Getenv("OS_PASSWORD"),
		DomainName:       os.Getenv("OS_USER_DOMAIN_NAME"),
		TenantName:       os.Getenv("OS_PROJECT_NAME"),
		AllowReauth:      true,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println("Failed to authenticate:", err)
		return nil
	}

	client, err := openstack.NewDNSV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create zone client:", err)
		return nil
	}

	allPages, err := zones.List(client, zones.ListOpts{}).AllPages()
	if err != nil {
		fmt.Println("Failed to list zones:", err)
		return nil
	}

	allZones, err := zones.ExtractZones(allPages)
	if err != nil || len(allZones) == 0 {
		fmt.Println("No zone found or extract failed:", err)
		return nil
	}

	return &allZones[0] // Return the first match
}

func FetchRecordsByZones(zoneId string, projectId string) []recordsets.RecordSet {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: os.Getenv("OS_AUTH_URL"),
		Username:         os.Getenv("OS_USERNAME"),
		Password:         os.Getenv("OS_PASSWORD"),
		DomainName:       os.Getenv("OS_USER_DOMAIN_NAME"),
		TenantName:       os.Getenv("OS_PROJECT_NAME"),
		AllowReauth:      true,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println("Failed to authenticate:", err)
		return nil
	}

	client, err := openstack.NewDNSV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create dns client:", err)
		return nil
	}

	allPages, err := recordsets.ListByZone(client, zoneId, recordsets.ListOpts{}).AllPages()
	if err != nil {
		fmt.Println("Failed to list zones:", err)
		return nil
	}

	recordsetList, err := recordsets.ExtractRecordSets(allPages)
	if err != nil {
		fmt.Println("Failed to extract zones:", err)
		return nil
	}

	return recordsetList
}
