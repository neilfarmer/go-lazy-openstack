package images

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	openstack_images "github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

// Image represents a simplified image object for UI use.
type Image struct {
	ID   string
	Name string
	OS   string
	Size int64
}

var mockImages = map[string][]images.Image{
	"dev": {
		{
			ID:        "img-ubuntu-2004",
			Name:      "Ubuntu 20.04 LTS",
			SizeBytes: 2300,
		},
		{
			ID:        "img-centos-7",
			Name:      "CentOS 7",
			SizeBytes: 1800,
		},
		{
			ID:        "img-centos-8",
			Name:      "CentOS 8",
			SizeBytes: 2000,
		},
		{
			ID:        "img-debian-11",
			Name:      "Debian 11",
			SizeBytes: 2100,
		},
	},
	"admin": {
		{
			ID:        "img-rocky-9",
			Name:      "Rocky Linux 9",
			SizeBytes: 1900,
		},
	},
	"qa": {
		{
			ID:        "img-alpine-3",
			Name:      "Alpine 3.16",
			SizeBytes: 512,
		},
	},
	"ops": {
		{
			ID:        "img-win2019",
			Name:      "Windows Server 2019",
			SizeBytes: 4096,
		},
		{
			ID:        "img-win2022",
			Name:      "Windows Server 2022",
			SizeBytes: 5000,
		},
	},
	"legacy": {
		{
			ID:        "img-debian-11",
			Name:      "Debian 11",
			SizeBytes: 2100,
		},
		{
			ID:        "img-fedora-35",
			Name:      "Fedora 35",
			SizeBytes: 2200,
		},
	},
	"research": {
		{
			ID:        "img-centos-7",
			Name:      "CentOS 7",
			SizeBytes: 1800,
		},
		{
			ID:        "img-centos-8",
			Name:      "CentOS 8",
			SizeBytes: 2000,
		},
	},
	"staging": {
		{
			ID:        "img-ubuntu-2004",
			Name:      "Ubuntu 20.04 LTS",
			SizeBytes: 2048,
		},
		{
			ID:        "img-ubuntu-2204",
			Name:      "Ubuntu 22.04 LTS",
			SizeBytes: 2300,
		},
	},
	"prod": {
		{
			ID:        "img-ubuntu-2004",
			Name:      "Ubuntu 20.04 LTS",
			SizeBytes: 2048,
		},
		{
			ID:        "img-ubuntu-2204",
			Name:      "Ubuntu 22.04 LTS",
			SizeBytes: 2300,
		},
		{
			ID:        "img-centos-7",
			Name:      "CentOS 7",
			SizeBytes: 1800,
		},
		{
			ID:        "img-centos-8",
			Name:      "CentOS 8",
			SizeBytes: 2000,
		},
		{
			ID:        "img-debian-11",
			Name:      "Debian 11",
			SizeBytes: 2100,
		},
		{
			ID:        "img-fedora-35",
			Name:      "Fedora 35",
			SizeBytes: 2200,
		},
		{
			ID:        "img-win2019",
			Name:      "Windows Server 2019",
			SizeBytes: 4096,
		},
		{
			ID:        "img-win2022",
			Name:      "Windows Server 2022",
			SizeBytes: 5000,
		},
		{
			ID:        "img-alpine-3",
			Name:      "Alpine 3.16",
			SizeBytes: 512,
		},
		{
			ID:        "img-rocky-9",
			Name:      "Rocky Linux 9",
			SizeBytes: 1900,
		},
	},
}


// fetchImages attempts to get a list of OpenStack images.
// If anything fails, it returns mockImages.
func FetchImages() []openstack_images.Image {
	opts, err := openstack.AuthOptionsFromEnv()
	project := os.Getenv("OS_PROJECT_NAME")
	if project == "" {
		panic("OS_PROJECT_NAME environment variable is not set")
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println("Failed to authenticate, using mock image data:", err)
		return mockImages[project]
	}

	client, err := openstack.NewImageServiceV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create image service client, using mock image data:", err)
		return mockImages[project]
	}

	allPages, err := images.List(client, nil).AllPages()
	if err != nil {
		fmt.Println("Failed to list images, using mock image data:", err)
		return mockImages[project]
	}

	imageList, err := images.ExtractImages(allPages)
	if err != nil || len(imageList) == 0 {
		fmt.Println("No images found or extract failed, using mock image data:", err)
		return mockImages[project]
	}

	return imageList
}