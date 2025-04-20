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

// mockFlavors are used if OpenStack fails or returns nothing.
var mockFlavors = map[string][]Flavor{
	"admin": {
		{"flavor-1", "m1.tiny", 1, 512, 1},
		{"flavor-2", "m1.small", 1, 2048, 20},
		{"flavor-3", "m1.medium", 2, 4096, 40},
	},
	"staging": {
		{"flavor-1", "m1.tiny", 1, 512, 1},
		{"flavor-2", "m1.small", 1, 2048, 20},
		{"flavor-3", "m1.medium", 2, 4096, 40},
	},
	"dev": {
		{"flavor-4", "m1.large", 4, 8192, 80},
		{"flavor-5", "m1.xlarge", 8, 16384, 160},
	},
	"research": {
		{"flavor-4", "m1.large", 4, 8192, 80},
		{"flavor-5", "m1.xlarge", 8, 16384, 160},
	},
	"ops": {
		{"flavor-1", "m1.tiny", 1, 512, 1},
		{"flavor-2", "m1.small", 1, 2048, 20},
		{"flavor-3", "m1.medium", 2, 4096, 40},
	},
	"qa": {
		{"flavor-4", "m1.large", 4, 8192, 80},
		{"flavor-5", "m1.xlarge", 8, 16384, 160},
	},
	"prod": {
		{"flavor-6", "compute-optimized", 16, 8192, 40},
		{"flavor-7", "memory-optimized", 8, 32768, 80},
		{"flavor-8", "storage-optimized", 4, 8192, 320},
	},
	"legacy": {
		{"flavor-6", "compute-optimized", 16, 8192, 40},
		{"flavor-7", "memory-optimized", 8, 32768, 80},
		{"flavor-8", "storage-optimized", 4, 8192, 320},
	},
}

// fetchFlavors attempts to retrieve flavors from OpenStack.
// Falls back to mockFlavors on failure or no data.
func FetchFlavors() []Flavor {
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

	var result []Flavor
	for _, f := range flavorList {
		result = append(result, Flavor{
			ID:    f.ID,
			Name:  f.Name,
			VCPUs: f.VCPUs,
			RAM:   f.RAM,
			Disk:  f.Disk,
		})
	}
	return result
}