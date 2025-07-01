package servers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	openstack_servers "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func FetchServers() []openstack_servers.Server {
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

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		fmt.Println("Failed to create compute client: ", err)
	}

	allPages, err := openstack_servers.List(client, nil).AllPages()
	if err != nil {
		fmt.Println("Failed to list servers: ", err)
	}

	serverList, err := openstack_servers.ExtractServers(allPages)
	if err != nil || len(serverList) == 0 {
		fmt.Println("No servers found or extract failed: ", err)
	}

	return serverList
}

func SshToServer(server openstack_servers.Server) {
	user := "test"
	keyPath := "~/.ssh/test"
	cmd := exec.Command("ssh", "-i", keyPath, fmt.Sprintf("%s@%s", user, server.AccessIPv4))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("SSH Failed: ", err)
	}
}
