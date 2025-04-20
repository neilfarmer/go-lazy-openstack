package projects

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
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
var mockProjects = []Project{
	{"proj-001", "admin", "Administrative project", "default", true},
	{"proj-002", "dev", "Development team", "default", true},
	{"proj-003", "qa", "Quality assurance", "default", true},
	{"proj-004", "ops", "Operations team", "default", true},
	{"proj-005", "legacy", "Legacy systems", "default", false},
	{"proj-006", "research", "R&D team sandbox", "default", true},
	{"proj-007", "staging", "Staging environment", "default", true},
	{"proj-008", "prod", "Production workloads", "default", true},
}

// fetchProjects tries to retrieve projects from Keystone.
// Falls back to mockProjects on error or empty response.
func FetchProjects() []Project {
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

	var result []Project
	for _, p := range projectList {
		result = append(result, Project{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			DomainID:    p.DomainID,
			Enabled:     p.Enabled,
		})
	}
	return result
}