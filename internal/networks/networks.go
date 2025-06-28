package networks

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

// FetchNetworks retrieves a list of OpenStack networks.
func FetchNetworks(projectId string) []networks.Network {
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

	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create network client:", err)
		return nil
	}

	allPages, err := networks.List(client, networks.ListOpts{
		ProjectID: projectId,
	}).AllPages()
	if err != nil {
		fmt.Println("Failed to list networks:", err)
		return nil
	}

	networkList, err := networks.ExtractNetworks(allPages)
	if err != nil {
		fmt.Println("Failed to extract networks:", err)
		return nil
	}

	return networkList
}

// FetchNetworkByID retrieves a single network by its ID.
func FetchNetworkByID(networkID string) *networks.Network {
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

	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create network client:", err)
		return nil
	}

	network, err := networks.Get(client, networkID).Extract()
	if err != nil {
		fmt.Println("Failed to get network details:", err)
		return nil
	}

	return network
}
