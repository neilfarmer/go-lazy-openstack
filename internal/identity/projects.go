package projects

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
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

// fetchProjects tries to retrieve projects from Keystone.
// Falls back to mockProjects on error or empty response.
func FetchProjects() []openstack_projects.Project {
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

	client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create identity client: ", err)
	}

	allPages, err := projects.List(client, nil).AllPages()
	if err != nil {
		fmt.Println("Failed to list projects: ", err)
	}

	projectList, err := projects.ExtractProjects(allPages)
	if err != nil || len(projectList) == 0 {
		fmt.Println("No projects found or extract failed: ", err)
	}

	return projectList
}

// FetchNetworkByID retrieves a single network by its ID.
func FetchProjectByID(projectID string) *projects.Project {
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

	client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create project client:", err)
		return nil
	}

	project, err := projects.Get(client, projectID).Extract()
	if err != nil {
		fmt.Println("Failed to get project details:", err)
		return nil
	}

	return project
}

func FetchProjectByName(projectName, domainId string) *projects.Project {
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

	client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create identity client:", err)
		return nil
	}

	listOpts := projects.ListOpts{
		Name:     projectName,
		DomainID: domainId,
	}

	allPages, err := projects.List(client, listOpts).AllPages()
	if err != nil {
		fmt.Println("Failed to list projects:", err)
		return nil
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil || len(allProjects) == 0 {
		fmt.Println("No project found or extract failed:", err)
		return nil
	}

	return &allProjects[0] // Return the first match
}

func FetchDomainIDByName(domainName string) domains.Domain {
	provider, err := openstack.AuthenticatedClient(gophercloud.AuthOptions{
		IdentityEndpoint: os.Getenv("OS_AUTH_URL"),
		Username:         os.Getenv("OS_USERNAME"),
		Password:         os.Getenv("OS_PASSWORD"),
		DomainName:       os.Getenv("OS_USER_DOMAIN_NAME"),
		TenantName:       os.Getenv("OS_PROJECT_NAME"),
	})
	if err != nil {
		fmt.Println("No domain found or extract failed:", err)
	}

	client, _ := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})

	allPages, _ := domains.List(client, domains.ListOpts{Name: domainName}).AllPages()

	allDomains, _ := domains.ExtractDomains(allPages)

	return allDomains[0]
}
