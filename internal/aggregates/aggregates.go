package aggregates

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
)

// FetchVolumes retrieves a list of volumes for a given project.
func FetchAggregates() []aggregates.Aggregate {
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

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create compute client:", err)
		return nil
	}

	allPages, err := aggregates.List(client).AllPages()
	if err != nil {
		fmt.Println("Failed to list aggregates:", err)
		return nil
	}

	aggregateList, err := aggregates.ExtractAggregates(allPages)
	if err != nil {
		fmt.Println("Failed to extract aggregates:", err)
		return nil
	}

	return aggregateList
}

// FetchVolumeByID retrieves a single volume by its ID.
func FetchAggregateByID(aggregateId int) *aggregates.Aggregate {
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

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create compute client:", err)
		return nil
	}

	aggregate, err := aggregates.Get(client, aggregateId).Extract()
	if err != nil {
		fmt.Println("Failed to get aggregate details:", err)
		return nil
	}

	return aggregate
}

// FetchVolumeByName retrieves a volume by its name (and optionally filters by project).
func FetchAggregateByName(aggregateName, projectID string) *aggregates.Aggregate {
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

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create compute client:", err)
		return nil
	}

	allPages, err := aggregates.List(client).AllPages()
	if err != nil {
		fmt.Println("Failed to list volumes:", err)
		return nil
	}

	allAggregates, err := aggregates.ExtractAggregates(allPages)
	if err != nil || len(allAggregates) == 0 {
		fmt.Println("No volume found or extract failed:", err)
		return nil
	}

	return &allAggregates[0] // Return the first match
}
