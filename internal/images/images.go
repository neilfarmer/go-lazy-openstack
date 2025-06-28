package images

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	openstack_images "github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

// fetchImages attempts to get a list of OpenStack images.
// If anything fails, it returns mockImages.
func FetchImages() []openstack_images.Image {
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

	client, err := openstack.NewImageServiceV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create image service client: ", err)
	}

	allPages, err := images.List(client, nil).AllPages()
	if err != nil {
		fmt.Println("Failed to list images: ", err)
	}

	imageList, err := images.ExtractImages(allPages)
	if err != nil || len(imageList) == 0 {
		fmt.Println("No images found or extract failed: ", err)
	}

	return imageList
}

func FetchImageByID(imageId string) *images.Image {
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

	client, err := openstack.NewImageServiceV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create compute client: ", err)
	}

	imageDetails, err := images.Get(client, imageId).Extract()
	if err != nil {
		fmt.Println("Failed to get details for image: ", err)
	}

	return imageDetails
}
