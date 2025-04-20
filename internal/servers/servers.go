package servers

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	openstack_servers "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type Server struct {
	Name     string
	ID       string
	Status   string
	Flavor   string
	Image    string
	Networks string
}

var mockServers = map[string][]openstack_servers.Server{
	"admin": {
		{
			Name: "web-01-derpify-access-worker-10099-and-stuff-or-something",
			ID: "123e4567", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.small"},
			Image: map[string]interface{}{"image": "ubuntu-20.04"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.10"},
		},
		{
			Name: "db-01", ID: "789a0123", Status: "SHUTOFF",
			Flavor: map[string]interface{}{"flavor": "m1.medium"},
			Image: map[string]interface{}{"image": "ubuntu-18.04"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.20"},
		},
		{
			Name: "app-01", ID: "456b7890", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.large"},
			Image: map[string]interface{}{"image": "centos-7"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.30"},
		},
		{
			Name: "cache-01", ID: "234c5678", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.small"},
			Image: map[string]interface{}{"image": "alpine-3.16"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.40"},
		},
	},
	"dev": {
		{
			Name: "proxy-01-nginx-frontend-v2-prod", ID: "345d6789", Status: "BUILD",
			Flavor: map[string]interface{}{"flavor": "m1.medium"},
			Image: map[string]interface{}{"image": "ubuntu-22.04"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.50"},
		},
		{
			Name: "monitor-01-prometheus-central-collector", ID: "567e8901", Status: "ERROR",
			Flavor: map[string]interface{}{"flavor": "m1.large"},
			Image: map[string]interface{}{"image": "debian-11"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.60"},
		},
		{
			Name: "worker-02-processing-long-job-handler", ID: "678f9012", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.xlarge"},
			Image: map[string]interface{}{"image": "ubuntu-20.04"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.70"},
		},
	},
	"qa": {
		{
			Name: "backup-01", ID: "789g0123", Status: "SHUTOFF",
			Flavor: map[string]interface{}{"flavor": "m1.medium"},
			Image: map[string]interface{}{"image": "centos-8"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.80"},
		},
		{
			Name: "test-01", ID: "890h1234", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.small"},
			Image: map[string]interface{}{"image": "fedora-35"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.90"},
		},
		{
			Name: "web-02-derpify-access-worker-10100-and-stuff", ID: "901i2345", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.small"},
			Image: map[string]interface{}{"image": "ubuntu-20.04"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.100"},
		},
		{
			Name: "db-02", ID: "012j3456", Status: "ERROR",
			Flavor: map[string]interface{}{"flavor": "m1.medium"},
			Image: map[string]interface{}{"image": "ubuntu-18.04"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.110"},
		},
		{
			Name: "batch-01", ID: "345k6789", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.large"},
			Image: map[string]interface{}{"image": "centos-7"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.120"},
		},
	},
	"test": {
		{
			Name: "test-web-01-derpify-access-worker-10099-and-stuff-or-something", ID: "123e4567", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.small"},
			Image: map[string]interface{}{"image": "ubuntu-20.04"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.10"},
		},
	},
	"prod": {
		{
			Name: "test-db-01", ID: "789a0123", Status: "SHUTOFF",
			Flavor: map[string]interface{}{"flavor": "m1.medium"},
			Image: map[string]interface{}{"image": "ubuntu-18.04"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.20"},
		},
		{
			Name: "test-app-01", ID: "456b7890", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.large"},
			Image: map[string]interface{}{"image": "centos-7"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.30"},
		},
		{
			Name: "test-cache-01", ID: "234c5678", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.small"},
			Image: map[string]interface{}{"image": "alpine-3.16"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.40"},
		},
		{
			Name: "test-proxy-01-nginx-frontend-v2-prod", ID: "345d6789", Status: "BUILD",
			Flavor: map[string]interface{}{"flavor": "m1.medium"},
			Image: map[string]interface{}{"image": "ubuntu-22.04"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.50"},
		},
	},
	"legacy": {
		{
			Name: "test-monitor-01-prometheus-central-collector", ID: "567e8901", Status: "ERROR",
			Flavor: map[string]interface{}{"flavor": "m1.large"},
			Image: map[string]interface{}{"image": "debian-11"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.60"},
		},
		{
			Name: "test-worker-02-processing-long-job-handler", ID: "678f9012", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.xlarge"},
			Image: map[string]interface{}{"image": "ubuntu-20.04"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.70"},
		},
		{
			Name: "test-backup-01", ID: "789g0123", Status: "SHUTOFF",
			Flavor: map[string]interface{}{"flavor": "m1.medium"},
			Image: map[string]interface{}{"image": "centos-8"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.80"},
		},
		{
			Name: "test-test-01", ID: "890h1234", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.small"},
			Image: map[string]interface{}{"image": "fedora-35"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.90"},
		},
		{
			Name: "test-web-02-derpify-access-worker-10100-and-stuff", ID: "901i2345", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.small"},
			Image: map[string]interface{}{"image": "ubuntu-20.04"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.100"},
		},
		{
			Name: "test-db-02", ID: "012j3456", Status: "ERROR",
			Flavor: map[string]interface{}{"flavor": "m1.medium"},
			Image: map[string]interface{}{"image": "ubuntu-18.04"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.110"},
		},
		{
			Name: "test-batch-01", ID: "345k6789", Status: "ACTIVE",
			Flavor: map[string]interface{}{"flavor": "m1.large"},
			Image: map[string]interface{}{"image": "centos-7"},
			Addresses: map[string]interface{}{"addresses": "net1=192.168.1.120"},
		},
	},
}


func FetchServers() []openstack_servers.Server {
	opts, err := openstack.AuthOptionsFromEnv()
	project := os.Getenv("OS_PROJECT_NAME")
	if project == "" {
		panic("OS_PROJECT_NAME environment variable is not set")
	}

	if err != nil {
		fmt.Println("Auth error, falling back to mock data:", err)
		return mockServers[project]
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println("Failed to authenticate, using mock data:", err)
		return mockServers[project]
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		fmt.Println("Failed to create compute client, using mock data:", err)
		return mockServers[project]
	}

	allPages, err := openstack_servers.List(client, nil).AllPages()
	if err != nil {
		fmt.Println("Failed to list servers, using mock data:", err)
		return mockServers[project]
	}

	serverList, err := openstack_servers.ExtractServers(allPages)
	if err != nil || len(serverList) == 0 {
		fmt.Println("No servers found or extract failed, using mock data:", err)
		return mockServers[project]
	}

	return serverList
}