package flavors

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

// Flavor represents a simplified view of a Nova flavor.
type Flavor struct {
	ID    string
	Name  string
	VCPUs int
	RAM   int // in MB
	Disk  int // in GB
}

// fetchFlavors attempts to retrieve flavors from OpenStack.
// Falls back to mockFlavors on failure or no data.
func FetchFlavors() []flavors.Flavor {
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
		fmt.Println("Failed to authenticate: ", err)
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create compute client: ", err)
	}

	allPages, err := flavors.ListDetail(client, nil).AllPages()
	if err != nil {
		fmt.Println("Failed to list flavors: ", err)
	}

	flavorList, err := flavors.ExtractFlavors(allPages)
	if err != nil || len(flavorList) == 0 {
		fmt.Println("No flavors found or extract failed: ", err)
	}

	return flavorList
}

func FetchFlavorByID(flavorId string) *flavors.Flavor {
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
		fmt.Println("Failed to authenticate: ", err)
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create compute client: ", err)
	}

	flavorDetails, err := flavors.Get(client, flavorId).Extract()
	if err != nil {
		fmt.Println("Failed to get details for flavor: ", err)
	}

	return flavorDetails
}
