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

var mockServers = map[string][]Server{
	"admin": {
		{"web-01-derpify-access-worker-10099-and-stuff-or-something", "123e4567", "ACTIVE", "m1.small", "ubuntu-20.04", "net1=192.168.1.10"},
		{"db-01", "789a0123", "SHUTOFF", "m1.medium", "ubuntu-18.04", "net1=192.168.1.20"},
		{"app-01", "456b7890", "ACTIVE", "m1.large", "centos-7", "net1=192.168.1.30"},
		{"cache-01", "234c5678", "ACTIVE", "m1.small", "alpine-3.16", "net1=192.168.1.40"},
	},
	"dev": {
		{"proxy-01-nginx-frontend-v2-prod", "345d6789", "BUILD", "m1.medium", "ubuntu-22.04", "net1=192.168.1.50"},
		{"monitor-01-prometheus-central-collector", "567e8901", "ERROR", "m1.large", "debian-11", "net1=192.168.1.60"},
		{"worker-02-processing-long-job-handler", "678f9012", "ACTIVE", "m1.xlarge", "ubuntu-20.04", "net1=192.168.1.70"},
	},
	"qa": {
		{"backup-01", "789g0123", "SHUTOFF", "m1.medium", "centos-8", "net1=192.168.1.80"},
		{"test-01", "890h1234", "ACTIVE", "m1.small", "fedora-35", "net1=192.168.1.90"},
		{"web-02-derpify-access-worker-10100-and-stuff", "901i2345", "ACTIVE", "m1.small", "ubuntu-20.04", "net1=192.168.1.100"},
		{"db-02", "012j3456", "ERROR", "m1.medium", "ubuntu-18.04", "net1=192.168.1.110"},
	},
	"ops": {
		{"batch-01", "345k6789", "ACTIVE", "m1.large", "centos-7", "net1=192.168.1.120"},
		{"test-web-01-derpify-access-worker-10099-and-stuff-or-something", "123e4567", "ACTIVE", "m1.small", "ubuntu-20.04", "net1=192.168.1.10"},
		{"test-db-01", "789a0123", "SHUTOFF", "m1.medium", "ubuntu-18.04", "net1=192.168.1.20"},
		{"test-app-01", "456b7890", "ACTIVE", "m1.large", "centos-7", "net1=192.168.1.30"},
	},
	"legacy": {
		{"test-cache-01", "234c5678", "ACTIVE", "m1.small", "alpine-3.16", "net1=192.168.1.40"},
		{"test-proxy-01-nginx-frontend-v2-prod", "345d6789", "BUILD", "m1.medium", "ubuntu-22.04", "net1=192.168.1.50"},
		{"test-monitor-01-prometheus-central-collector", "567e8901", "ERROR", "m1.large", "debian-11", "net1=192.168.1.60"},
	},
	"research": {
		{"test-worker-02-processing-long-job-handler", "678f9012", "ACTIVE", "m1.xlarge", "ubuntu-20.04", "net1=192.168.1.70"},
		{"test-backup-01", "789g0123", "SHUTOFF", "m1.medium", "centos-8", "net1=192.168.1.80"},
		{"test-test-01", "890h1234", "ACTIVE", "m1.small", "fedora-35", "net1=192.168.1.90"},
		{"test-web-02-derpify-access-worker-10100-and-stuff", "901i2345", "ACTIVE", "m1.small", "ubuntu-20.04", "net1=192.168.1.100"},

	},
	"staging": {
		{"test-santa-web-01-derpify-access-worker-10099-and-stuff-or-something", "123e4567", "ACTIVE", "m1.small", "ubuntu-20.04", "net1=192.168.1.10"},
		{"test-santa-db-01", "789a0123", "SHUTOFF", "m1.medium", "ubuntu-18.04", "net1=192.168.1.20"},
		{"test-santa-app-01", "456b7890", "ACTIVE", "m1.large", "centos-7", "net1=192.168.1.30"},
	},
	"prod": {
		{"test-db-02", "012j3456", "ERROR", "m1.medium", "ubuntu-18.04", "net1=192.168.1.110"},
		{"test-batch-01", "345k6789", "ACTIVE", "m1.large", "centos-7", "net1=192.168.1.120"},
		
		{"test-santa-cache-01", "234c5678", "ACTIVE", "m1.small", "alpine-3.16", "net1=192.168.1.40"},
		{"test-santa-proxy-01-nginx-frontend-v2-prod", "345d6789", "BUILD", "m1.medium", "ubuntu-22.04", "net1=192.168.1.50"},
		{"test-santa-monitor-01-prometheus-central-collector", "567e8901", "ERROR", "m1.large", "debian-11", "net1=192.168.1.60"},
		{"test-santa-worker-02-processing-long-job-handler", "678f9012", "ACTIVE", "m1.xlarge", "ubuntu-20.04", "net1=192.168.1.70"},
		{"test-santa-backup-01", "789g0123", "SHUTOFF", "m1.medium", "centos-8", "net1=192.168.1.80"},
		{"test-santa-test-01", "890h1234", "ACTIVE", "m1.small", "fedora-35", "net1=192.168.1.90"},
		{"test-santa-web-02-derpify-access-worker-10100-and-stuff", "901i2345", "ACTIVE", "m1.small", "ubuntu-20.04", "net1=192.168.1.100"},
		{"test-santa-db-02", "012j3456", "ERROR", "m1.medium", "ubuntu-18.04", "net1=192.168.1.110"},
		{"test-santa-batch-01", "345k6789", "ACTIVE", "m1.large", "centos-7", "net1=192.168.1.120"},
	},
}

func FetchServers() []Server {
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

	var result []Server
	for _, s := range serverList {
		result = append(result, Server{
			Name:     s.Name,
			ID:       s.ID,
			Status:   s.Status,
			Flavor:   fmt.Sprintf("%v", s.Flavor["id"]),
			Image:    fmt.Sprintf("%v", s.Image["id"]),
			Networks: fmt.Sprintf("%v", s.Addresses),
		})
	}
	return result
}