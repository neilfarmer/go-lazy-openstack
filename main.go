package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/neilfarmer/internal/aggregates"
	"github.com/neilfarmer/internal/dns"
	"github.com/neilfarmer/internal/flavors"
	"github.com/neilfarmer/internal/hypervisors"
	projects "github.com/neilfarmer/internal/identity"
	"github.com/neilfarmer/internal/images"
	"github.com/neilfarmer/internal/loadbalancers"
	"github.com/neilfarmer/internal/networks"
	"github.com/neilfarmer/internal/servers"
	"github.com/neilfarmer/internal/volumes"
	"github.com/rivo/tview"
)

var pages *tview.Pages
var inputPrompt *tview.InputField
var headerFlex *tview.Flex
var detailsView *tview.TextView
var aggregatesList *tview.List
var hypervisorsList *tview.List
var serverList *tview.List
var imagesList *tview.List
var flavorsList *tview.List
var projectsList *tview.List
var networksList *tview.List
var volumesList *tview.List
var loadbalancersList *tview.List
var dnsList *tview.List

var currentServer *servers.Server
var currentView string

var knownCommands = []string{
	"servers",
	"aggregates",
	"hypervisors",
	"images",
	"flavors",
	"projects",
	"volumes",
	"loadbalancers",
	"dns",
	"networks",
}

var acceptShortcuts = true

func main() {
	// Root application
	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == ':' {
			acceptShortcuts = false
			headerFlex.ResizeItem(inputPrompt, 3, 1) // show prompt
			app.SetFocus(inputPrompt)
			return nil
		}

		if acceptShortcuts {
			switch event.Rune() {
			case 'a':
				populateAggregatesList()
				pages.SwitchToPage("aggregates")
				detailsView.Clear()
			case 'i':
				populateImagesList()
				pages.SwitchToPage("images")
				detailsView.Clear()
			case 'f':
				populateFlavorsList()
				pages.SwitchToPage("flavors")
				detailsView.Clear()
			case 'l':
				populateLoadbalancersList()
				pages.SwitchToPage("loadbalancers")
				detailsView.Clear()
			case 'h':
				populateHypervisorsList()
				pages.SwitchToPage("hypervisors")
				detailsView.Clear()
			case 'd':
				populateDnsList()
				pages.SwitchToPage("dns")
				detailsView.Clear()
			case 'v':
				populateVolumesList()
				pages.SwitchToPage("volumes")
				detailsView.Clear()
			case 's':
				populateServersList()
				pages.SwitchToPage("servers")
				detailsView.Clear()
			case 'n':
				populateNetworksList()
				pages.SwitchToPage("networks")
				detailsView.Clear()
			case 'p':
				pages.SwitchToPage("projects")
				detailsView.Clear()
			case 'x':
				if currentView == "servers" && currentServer != nil {
					go servers.SshToServer(*currentServer)
				}
			case 'q':
				app.Stop()
			default:
				return event
			}
		}

		return event
	})

	pages = tview.NewPages()

	header := tview.NewTextView()
	header.SetDynamicColors(true).SetTextAlign(tview.AlignLeft).SetBorder(true).SetTitle(" Lazy Openstack ")
	header.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
			return nil
		}
		return event
	})

	go func() {
		for {
			now := time.Now().Format("15:04:05")
			app.QueueUpdateDraw(func() {
				header.Clear()
				_, _, width, _ := header.GetInnerRect()

				shortcuts := "(a)ggregates (p)rojects (d)ns (i)mages (f)lavors (h)ypervisors (l)oadbalancers (s)ervers (n)etworks (v)olumes (q)uit"
				if currentView == "servers" && currentServer != nil {
					shortcuts += " (x)ssh"
				}
				text := fmt.Sprintf("%s%*s", shortcuts, width-len(shortcuts), now)
				fmt.Fprintf(header, "%s", text)
				fmt.Fprintf(header, "\n")
				fmt.Fprintf(header, "Current Project: %s\n", os.Getenv("OS_PROJECT_NAME"))
			})

			time.Sleep(100 * time.Millisecond)
		}
	}()

	detailsView = tview.NewTextView()
	detailsView.SetBorder(true).SetTitle(" Details ").SetTitleAlign(tview.AlignCenter)

	aggregatesList = tview.NewList()
	aggregatesList.SetBorder(true).SetTitle(" Aggregates ").SetTitleAlign(tview.AlignCenter)

	serverList = tview.NewList()
	serverList.SetBorder(true).SetTitle(" Servers ").SetTitleAlign(tview.AlignCenter)

	imagesList = tview.NewList()
	imagesList.SetBorder(true).SetTitle(" Images ").SetTitleAlign(tview.AlignCenter)

	flavorsList = tview.NewList()
	flavorsList.SetBorder(true).SetTitle(" Flavors ").SetTitleAlign(tview.AlignCenter)

	hypervisorsList = tview.NewList()
	hypervisorsList.SetBorder(true).SetTitle(" Hypervisors ").SetTitleAlign(tview.AlignCenter)

	volumesList = tview.NewList()
	volumesList.SetBorder(true).SetTitle(" Volumes ").SetTitleAlign(tview.AlignCenter)

	loadbalancersList = tview.NewList()
	loadbalancersList.SetBorder(true).SetTitle(" Loadbalancers ").SetTitleAlign(tview.AlignCenter)

	dnsList = tview.NewList()
	dnsList.SetBorder(true).SetTitle(" Dns ").SetTitleAlign(tview.AlignCenter)

	networksList = tview.NewList()
	networksList.SetBorder(true).SetTitle(" Networks ").SetTitleAlign(tview.AlignCenter)

	projectsList = tview.NewList()
	projectsList.SetBorder(true).SetTitle(" Projects ").SetTitleAlign(tview.AlignCenter)
	for _, project := range projects.FetchProjects() {
		projectsList.AddItem(project.Name, "", -1, func() {
			detailsView.Clear()
			os.Setenv("OS_PROJECT_NAME", project.Name)
			fmt.Fprintf(detailsView, "Current Project Set To:\nID: %s\nName: %s\nDescription: %s\nDomainID: %s\nEnabled: %t", project.ID, project.Name, project.Description, project.DomainID, project.Enabled)
		})
	}

	// This is our handy dandy input prompt
	inputPrompt = tview.NewInputField()
	inputPrompt.SetLabel("Command: ").
		SetFieldWidth(30).
		SetBorder(true).
		SetTitle(" Command Prompt ").
		SetTitleAlign(tview.AlignLeft)

	inputPrompt.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			command := inputPrompt.GetText()
			inputPrompt.SetText("")
			if command == "aggregates" {
				populateAggregatesList()
				pages.SwitchToPage(command)
				detailsView.Clear()
			}

			if command == "images" {
				populateImagesList()
				pages.SwitchToPage(command)
				detailsView.Clear()
			}
			if command == "servers" {
				populateServersList()
				pages.SwitchToPage(command)
				detailsView.Clear()
			}

			if command == "flavors" {
				populateFlavorsList()
				pages.SwitchToPage(command)
				detailsView.Clear()
			}
			if command == "hypervisors" {
				populateHypervisorsList()
				pages.SwitchToPage(command)
				detailsView.Clear()
			}
			if command == "volumes" {
				populateVolumesList()
				pages.SwitchToPage(command)
				detailsView.Clear()
			}
			if command == "loadbalancers" {
				populateLoadbalancersList()
				pages.SwitchToPage(command)
				detailsView.Clear()
			}
			if command == "dns" {
				populateDnsList()
				pages.SwitchToPage(command)
				detailsView.Clear()
			}
			if command == "networks" {
				populateNetworksList()
				pages.SwitchToPage(command)
				detailsView.Clear()
			}

			if command == "projects" {
				pages.SwitchToPage(command)
				detailsView.Clear()
			}

			headerFlex.ResizeItem(inputPrompt, 0, 0) // hide prompt
			acceptShortcuts = true
		}
	})

	inputPrompt.SetAutocompleteFunc(func(currentText string) (entries []string) {
		for _, cmd := range knownCommands {
			if len(currentText) > 0 && cmd != currentText && len(cmd) >= len(currentText) && cmd[:len(currentText)] == currentText {
				entries = append(entries, cmd)
			}
		}
		return entries
	})

	headerFlex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 0, 2, false).
		AddItem(inputPrompt, 0, 0, false)

	aggregateViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerFlex, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(aggregatesList, 0, 1, true).
			AddItem(detailsView, 0, 4, false),
			0, 5, true)

	serverViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerFlex, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(serverList, 0, 1, true).
			AddItem(detailsView, 0, 4, false),
			0, 5, true)

	imageViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerFlex, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(imagesList, 0, 1, true).
			AddItem(detailsView, 0, 4, false),
			0, 5, true)

	flavorViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerFlex, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(flavorsList, 0, 1, true).
			AddItem(detailsView, 0, 4, false),
			0, 5, true)

	hypervisorsViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerFlex, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(hypervisorsList, 0, 1, true).
			AddItem(detailsView, 0, 4, false),
			0, 5, true)

	volumeViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerFlex, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(volumesList, 0, 1, true).
			AddItem(detailsView, 0, 4, false),
			0, 5, true)

	loadbalancerViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerFlex, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(loadbalancersList, 0, 1, true).
			AddItem(detailsView, 0, 4, false),
			0, 5, true)

	dnsViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerFlex, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(dnsList, 0, 1, true).
			AddItem(detailsView, 0, 4, false),
			0, 5, true)

	networksViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerFlex, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(networksList, 0, 1, true).
			AddItem(detailsView, 0, 4, false),
			0, 5, true)

	projectsViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(headerFlex, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(projectsList, 0, 1, true).
			AddItem(detailsView, 0, 4, false),
			0, 5, true)

	populateServersList()

	pages.AddPage("prompt", inputPrompt, true, false)
	pages.AddPage("aggregates", aggregateViewFlex, true, true)
	pages.AddPage("servers", serverViewFlex, true, true)
	pages.AddPage("images", imageViewFlex, true, true)
	pages.AddPage("flavors", flavorViewFlex, true, true)
	pages.AddPage("hypervisors", hypervisorsViewFlex, true, true)
	pages.AddPage("volumes", volumeViewFlex, true, true)
	pages.AddPage("loadbalancers", loadbalancerViewFlex, true, true)
	pages.AddPage("dns", dnsViewFlex, true, true)
	pages.AddPage("networks", networksViewFlex, true, true)
	pages.AddPage("projects", projectsViewFlex, true, true)

	err := app.SetRoot(pages, true).Run()
	if err != nil {
		panic(err)
	}
}

func populateServersList() {
	serverList.Clear()
	for _, server := range servers.FetchServers() {
		serverList.AddItem(server.Name, "", -1, func() {
			detailsView.Clear()
			currentServer = &s
			currentView = "servers"

			flavorID, _ := server.Flavor["id"].(string)
			flavor := flavors.FetchFlavorByID(flavorID)
			flavorInfo := fmt.Sprintf("\n\tName: %s, \n\tRAM: %dMB, \n\tvCPUs: %d, \n\tDisk: %dGB", flavor.Name, flavor.RAM, flavor.VCPUs, flavor.Disk)

			imageID, _ := server.Image["id"].(string)
			image := images.FetchImageByID(imageID)
			imageInfo := fmt.Sprintf("\n\tName: %s,\n\tID: %s, \n\tSize: %dMB, \n\tTags: %s", image.Name, imageID, image.SizeBytes, image.Tags)

			addresses, err := json.Marshal(server.Addresses)
			if err != nil {
				addresses = []byte("unable to marshal addresses")
			}
			fmt.Fprintf(detailsView, "ID: %s\nStatus: %s\nFlavor: %s\nImage: %s\nNetworks: %s\nAttached Volumes: %s", server.ID, server.Status, flavorInfo, imageInfo, addresses, server.AttachedVolumes)
		})
	}
}

func populateAggregatesList() {
	aggregatesList.Clear()
	for _, aggregate := range aggregates.FetchAggregates() {
		var hosts string
		for _, host := range aggregate.Hosts {
			hosts += fmt.Sprintf("\n\t%s", host)
		}
		aggregatesList.AddItem(aggregate.Name, "", -1, func() {
			detailsView.Clear()
			fmt.Fprintf(detailsView, "ID: %d\nName: %s\nMetadata: %s\nHosts: %s", aggregate.ID, aggregate.Name, aggregate.Metadata, hosts)
		})
	}
}

func populateFlavorsList() {
	flavorsList.Clear()
	for _, flavor := range flavors.FetchFlavors() {
		flavorsList.AddItem(flavor.Name, "", -1, func() {
			detailsView.Clear()
			fmt.Fprintf(detailsView, "ID: %s\nName: %s\nVCPU: %d\nRAM: %d\nDisk: %d", flavor.ID, flavor.Name, flavor.VCPUs, flavor.RAM, flavor.Disk)
		})
	}
}

func populateVolumesList() {
	volumesList.Clear()
	for _, volume := range volumes.FetchVolumes() {
		volumesList.AddItem(volume.Name, "", -1, func() {
			detailsView.Clear()
			fmt.Fprintf(detailsView, "\n\tID: %s\n\tName: %s\n\tDescription: %s\n\tCreated at: %s\n\tSize: %d\n\tType: %s", volume.ID, volume.Name, volume.Description, volume.CreatedAt, volume.Size, volume.VolumeType)
		})
	}
}

func populateHypervisorsList() {
	hypervisorsList.Clear()
	for _, hypervisor := range hypervisors.FetchHypervisors() {
		hypervisorsList.AddItem(hypervisor.HypervisorHostname, "", -1, func() {
			detailsView.Clear()
			fmt.Fprintf(detailsView, "\n\tHostname: %s\n\tType: %s\n\tHost IP: %s\n\tState: %s\n\tCPU Info: %s", hypervisor.HypervisorHostname, hypervisor.HypervisorType, hypervisor.HostIP, hypervisor.State, hypervisor.CPUInfo)
		})
	}
}

func populateLoadbalancersList() {
	loadbalancersList.Clear()
	domain := projects.FetchDomainIDByName(os.Getenv("OS_USER_DOMAIN_NAME"))
	project := projects.FetchProjectByName(os.Getenv("OS_PROJECT_NAME"), domain.ID)
	for _, loadbalancer := range loadbalancers.FetchLoadbalancers(project.ID) {
		loadbalancersList.AddItem(loadbalancer.Name, "", -1, func() {
			detailsView.Clear()
			fmt.Fprintf(detailsView, "\n\tID: %s\n\tName: %s\n\tVIP Address: %s\n\tOperating Status: %s\n\tProvisioning Status: %s\n\t", loadbalancer.ID, loadbalancer.Name, loadbalancer.VipAddress, loadbalancer.OperatingStatus, loadbalancer.ProvisioningStatus)
		})
	}
}

func populateDnsList() {
	dnsList.Clear()
	domain := projects.FetchDomainIDByName(os.Getenv("OS_USER_DOMAIN_NAME"))
	project := projects.FetchProjectByName(os.Getenv("OS_PROJECT_NAME"), domain.ID)
	for _, zone := range dns.FetchZones(project.ID) {
		var recordsetListByZone string
		dnsList.AddItem(zone.Name, "", -1, func() {
			for _, recordset := range dns.FetchRecordsByZones(zone.ID, project.ID) {
				recordsetListByZone += fmt.Sprintf("\n\t\tName: %s,\n\t\tID: %s\n\t\tRecords: %s", recordset.Name, recordset.ID, recordset.Records)
			}
			detailsView.Clear()
			fmt.Fprintf(detailsView, "\n\tID: %s\n\tName: %s\n\tTTL: %d\n\tStatus: %s\n\tEmail: %s\n\tPool: %s\n\tRecords: %s", zone.ID, zone.Name, zone.TTL, zone.Status, zone.Email, zone.PoolID, recordsetListByZone)
		})
	}
}

func populateNetworksList() {
	networksList.Clear()
	domain := projects.FetchDomainIDByName(os.Getenv("OS_USER_DOMAIN_NAME"))
	project := projects.FetchProjectByName(os.Getenv("OS_PROJECT_NAME"), domain.ID)
	for _, network := range networks.FetchNetworks(project.ID) {
		networksList.AddItem(network.Name, "", -1, func() {
			detailsView.Clear()
			fmt.Fprintf(detailsView, "ID: %s\nName: %s", network.ID, network.Name)
		})
	}
}

func populateImagesList() {
	imagesList.Clear()
	for _, image := range images.FetchImages() {
		imagesList.AddItem(image.Name, "", -1, func() {
			detailsView.Clear()
			fmt.Fprintf(detailsView, "ID: %s\nName: %s\nSize: %d", image.ID, image.Name, image.SizeBytes)
		})
	}
}
