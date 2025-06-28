package volumes

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
)

// FetchVolumes retrieves a list of volumes for a given project.
func FetchVolumes() []volumes.Volume {
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

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create block storage client:", err)
		return nil
	}

	allPages, err := volumes.List(client, volumes.ListOpts{}).AllPages()
	if err != nil {
		fmt.Println("Failed to list volumes:", err)
		return nil
	}

	volumeList, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		fmt.Println("Failed to extract volumes:", err)
		return nil
	}

	return volumeList
}

// FetchVolumeByID retrieves a single volume by its ID.
func FetchVolumeByID(volumeID string) *volumes.Volume {
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

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create block storage client:", err)
		return nil
	}

	volume, err := volumes.Get(client, volumeID).Extract()
	if err != nil {
		fmt.Println("Failed to get volume details:", err)
		return nil
	}

	return volume
}

// FetchVolumeByName retrieves a volume by its name (and optionally filters by project).
func FetchVolumeByName(volumeName, projectID string) *volumes.Volume {
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

	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create block storage client:", err)
		return nil
	}

	listOpts := volumes.ListOpts{
		Name:     volumeName,
		TenantID: projectID,
	}

	allPages, err := volumes.List(client, listOpts).AllPages()
	if err != nil {
		fmt.Println("Failed to list volumes:", err)
		return nil
	}

	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil || len(allVolumes) == 0 {
		fmt.Println("No volume found or extract failed:", err)
		return nil
	}

	return &allVolumes[0] // Return the first match
}
