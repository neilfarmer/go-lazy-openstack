package projects

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	openstack_projects "github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

// Project represents a simplified project for UI usage.
type Project struct {
	ID          string
	Name        string
	Description string
	DomainID    string
	Enabled     bool
}

// mockProjects are used if Keystone fails or returns nothing.
var mockProjects = []openstack_projects.Project{
	{
		ID:          "proj-001",
		Name:        "admin",
		Description: "Administrative project",
		DomainID:    "default",
		Enabled:     true,
	},
	{
		ID:          "proj-002",
		Name:        "dev",
		Description: "Development team",
		DomainID:    "default",
		Enabled:     true,
	},
	{
		ID:          "proj-003",
		Name:        "qa",
		Description: "Quality assurance",
		DomainID:    "default",
		Enabled:     true,
	},
	{
		ID:          "proj-004",
		Name:        "ops",
		Description: "Operations team",
		DomainID:    "default",
		Enabled:     true,
	},
	{
		ID:          "proj-005",
		Name:        "legacy",
		Description: "Legacy systems",
		DomainID:    "default",
		Enabled:     false,
	},
	{
		ID:          "proj-006",
		Name:        "research",
		Description: "R&D team sandbox",
		DomainID:    "default",
		Enabled:     true,
	},
	{
		ID:          "proj-007",
		Name:        "staging",
		Description: "Staging environment",
		DomainID:    "default",
		Enabled:     true,
	},
	{
		ID:          "proj-008",
		Name:        "prod",
		Description: "Production workloads",
		DomainID:    "default",
		Enabled:     true,
	},
}


// fetchProjects tries to retrieve projects from Keystone.
// Falls back to mockProjects on error or empty response.
func FetchProjects() []openstack_projects.Project {
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		fmt.Println("Auth error, falling back to mock projects:", err)
		return mockProjects
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println("Failed to authenticate, using mock projects:", err)
		return mockProjects
	}

	client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create identity client, using mock projects:", err)
		return mockProjects
	}

	allPages, err := projects.List(client, nil).AllPages()
	if err != nil {
		fmt.Println("Failed to list projects, using mock projects:", err)
		return mockProjects
	}

	projectList, err := projects.ExtractProjects(allPages)
	if err != nil || len(projectList) == 0 {
		fmt.Println("No projects found or extract failed, using mock projects:", err)
		return mockProjects
	}

	return projectList
}