package loadbalancers

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
)

func FetchLoadbalancers(projectId string) []loadbalancers.LoadBalancer {
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

	client, err := openstack.NewLoadBalancerV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create loadbalancer client:", err)
		return nil
	}

	allPages, err := loadbalancers.List(client, loadbalancers.ListOpts{
		ProjectID: projectId,
	}).AllPages()
	if err != nil {
		fmt.Println("Failed to list loadbalancers:", err)
		return nil
	}

	loadbalancerList, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil {
		fmt.Println("Failed to extract loadbalancers:", err)
		return nil
	}

	return loadbalancerList
}

func FetchLoadbalancerByID(loadbalancerID string) *loadbalancers.LoadBalancer {
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

	client, err := openstack.NewLoadBalancerV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create loadbalancer client:", err)
		return nil
	}

	loadbalancer, err := loadbalancers.Get(client, loadbalancerID).Extract()
	if err != nil {
		fmt.Println("Failed to get loadbalancer details:", err)
		return nil
	}

	return loadbalancer
}

func FetchLoadbalancerByName(loadbalancerName, projectID string) *loadbalancers.LoadBalancer {
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

	client, err := openstack.NewLoadBalancerV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create block storage client:", err)
		return nil
	}

	allPages, err := loadbalancers.List(client, loadbalancers.ListOpts{}).AllPages()
	if err != nil {
		fmt.Println("Failed to list loadbalancers:", err)
		return nil
	}

	allLoadbalancers, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil || len(allLoadbalancers) == 0 {
		fmt.Println("No loadbalancers found or extract failed:", err)
		return nil
	}

	return &allLoadbalancers[0] // Return the first match
}
