package images

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

// Image represents a simplified image object for UI use.
type Image struct {
	ID   string
	Name string
	OS   string
	Size int64
}

// mockImages are used if OpenStack Glance fails or returns nothing.
var mockImages = map[string][]Image{
	"dev": {
		{"img-ubuntu-2004", "Ubuntu 20.04 LTS", "linux", 2048},
		{"img-ubuntu-2204", "Ubuntu 22.04 LTS", "linux", 2300},
		{"img-centos-7", "CentOS 7", "linux", 1800},
		{"img-centos-8", "CentOS 8", "linux", 2000},
		{"img-debian-11", "Debian 11", "linux", 2100},
	},
	"admin": {
		{"img-rocky-9", "Rocky Linux 9", "linux", 1900},
	},
	"qa": {
		{"img-alpine-3", "Alpine 3.16", "linux", 512},
	},
	"ops": {
		{"img-win2019", "Windows Server 2019", "windows", 4096},
		{"img-win2022", "Windows Server 2022", "windows", 5000},
	},
	"legacy": {
		{"img-debian-11", "Debian 11", "linux", 2100},
		{"img-fedora-35", "Fedora 35", "linux", 2200},
	},
	"research": {
		{"img-centos-7", "CentOS 7", "linux", 1800},
		{"img-centos-8", "CentOS 8", "linux", 2000},
	},
	"staging": {
		{"img-ubuntu-2004", "Ubuntu 20.04 LTS", "linux", 2048},
		{"img-ubuntu-2204", "Ubuntu 22.04 LTS", "linux", 2300},
	},
	"prod": {
		{"img-ubuntu-2004", "Ubuntu 20.04 LTS", "linux", 2048},
		{"img-ubuntu-2204", "Ubuntu 22.04 LTS", "linux", 2300},
		{"img-centos-7", "CentOS 7", "linux", 1800},
		{"img-centos-8", "CentOS 8", "linux", 2000},
		{"img-debian-11", "Debian 11", "linux", 2100},
		{"img-fedora-35", "Fedora 35", "linux", 2200},
		{"img-win2019", "Windows Server 2019", "windows", 4096},
		{"img-win2022", "Windows Server 2022", "windows", 5000},
		{"img-alpine-3", "Alpine 3.16", "linux", 512},
		{"img-rocky-9", "Rocky Linux 9", "linux", 1900},
	},
}

// fetchImages attempts to get a list of OpenStack images.
// If anything fails, it returns mockImages.
func FetchImages() []Image {
	opts, err := openstack.AuthOptionsFromEnv()
	project := os.Getenv("OS_PROJECT_NAME")
	if err != nil {
		fmt.Println("Auth error, falling back to mock image data:", err)
		return mockImages[project]
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

	var result []Image
	for _, img := range imageList {
		result = append(result, Image{
			ID:   img.ID,
			Name: img.Name,
			OS:   fmt.Sprintf("%v", img.Properties["os_distro"]),
			Size: img.SizeBytes / 1024 / 1024, // convert to MB
		})
	}
	return result
}