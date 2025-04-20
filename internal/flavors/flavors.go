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

var mockFlavors = map[string][]flavors.Flavor{
	"admin": {
		{Name: "flavor-1", ID: "m1.tiny", VCPUs: 1, RAM: 512, Disk: 1},
		{Name: "flavor-2", ID: "m1.small", VCPUs: 1, RAM: 2048, Disk: 20},
		{Name: "flavor-3", ID: "m1.medium", VCPUs: 2, RAM: 4096, Disk: 40},
	},
	"staging": {
		{Name: "flavor-1", ID: "m1.tiny", VCPUs: 1, RAM: 512, Disk: 1},
		{Name: "flavor-2", ID: "m1.small", VCPUs: 1, RAM: 2048, Disk: 20},
		{Name: "flavor-3", ID: "m1.medium", VCPUs: 2, RAM: 4096, Disk: 40},
	},
	"dev": {
		{Name: "flavor-4", ID: "m1.large", VCPUs: 4, RAM: 8192, Disk: 80},
		{Name: "flavor-5", ID: "m1.xlarge", VCPUs: 8, RAM: 16384, Disk: 160},
	},
	"research": {
		{Name: "flavor-4", ID: "m1.large", VCPUs: 4, RAM: 8192, Disk: 80},
		{Name: "flavor-5", ID: "m1.xlarge", VCPUs: 8, RAM: 16384, Disk: 160},
	},
	"ops": {
		{Name: "flavor-1", ID: "m1.tiny", VCPUs: 1, RAM: 512, Disk: 1},
		{Name: "flavor-2", ID: "m1.small", VCPUs: 1, RAM: 2048, Disk: 20},
		{Name: "flavor-3", ID: "m1.medium", VCPUs: 2, RAM: 4096, Disk: 40},
	},
	"qa": {
		{Name: "flavor-4", ID: "m1.large", VCPUs: 4, RAM: 8192, Disk: 80},
		{Name: "flavor-5", ID: "m1.xlarge", VCPUs: 8, RAM: 16384, Disk: 160},
	},
	"prod": {
		{Name: "flavor-6", ID: "compute-optimized", VCPUs: 16, RAM: 8192, Disk: 40},
		{Name: "flavor-7", ID: "memory-optimized", VCPUs: 8, RAM: 32768, Disk: 80},
		{Name: "flavor-8", ID: "storage-optimized", VCPUs: 4, RAM: 8192, Disk: 320},
	},
	"legacy": {
		{Name: "flavor-6", ID: "compute-optimized", VCPUs: 16, RAM: 8192, Disk: 40},
		{Name: "flavor-7", ID: "memory-optimized", VCPUs: 8, RAM: 32768, Disk: 80},
		{Name: "flavor-8", ID: "storage-optimized", VCPUs: 4, RAM: 8192, Disk: 320},
	},
}


// fetchFlavors attempts to retrieve flavors from OpenStack.
// Falls back to mockFlavors on failure or no data.
func FetchFlavors() []flavors.Flavor {
	opts, err := openstack.AuthOptionsFromEnv()
	project := os.Getenv("OS_PROJECT_NAME")
	if project == "" {
		panic("OS_PROJECT_NAME environment variable is not set")
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println("Failed to authenticate, using mock flavors:", err)
		return mockFlavors[project]
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create compute client, using mock flavors:", err)
		return mockFlavors[project]
	}

	allPages, err := flavors.ListDetail(client, nil).AllPages()
	if err != nil {
		fmt.Println("Failed to list flavors, using mock flavors:", err)
		return mockFlavors[project]
	}

	flavorList, err := flavors.ExtractFlavors(allPages)
	if err != nil || len(flavorList) == 0 {
		fmt.Println("No flavors found or extract failed, using mock flavors:", err)
		return mockFlavors[project]
	}

	return flavorList
}