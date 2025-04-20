package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/neilfarmer/internal/flavors"
	projects "github.com/neilfarmer/internal/identity"
	"github.com/neilfarmer/internal/images"
	"github.com/neilfarmer/internal/servers"
	"github.com/rivo/tview"
)

var pages *tview.Pages
var inputPrompt *tview.InputField
var headerFlex *tview.Flex
var detailsView *tview.TextView
var serverList *tview.List
var imagesList *tview.List
var flavorsList *tview.List
var projectsList *tview.List

var knownCommands = []string{
	"servers",
	"images",
	"flavors",
	"projects",
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
			case 'i':
				populateImagesList()
				pages.SwitchToPage("images")
				detailsView.Clear()
			case 'f':
				populateFlavorsList()
				pages.SwitchToPage("flavors")
				detailsView.Clear()
			case 's':
				populateServersList()
				pages.SwitchToPage("servers")
				detailsView.Clear()
			case 'p':
				pages.SwitchToPage("projects")
				detailsView.Clear()
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

	go func () {
		for {
			now := time.Now().Format("15:04:05")
			app.QueueUpdateDraw(func() {
				header.Clear()
				_, _, width, _ := header.GetInnerRect()
				
				shortcuts := "(p)rojects (i)mages (f)flavors (s)ervers (q)uit"
				text := fmt.Sprintf("%s%*s", shortcuts, width-len(shortcuts), now)
				fmt.Fprintf(header, text)
				fmt.Fprintf(header, "\n")
				fmt.Fprintf(header, "Current Project: %s\n", os.Getenv("OS_PROJECT_NAME"))
			})

			time.Sleep(100 * time.Millisecond)
		}
	}()

	detailsView = tview.NewTextView()
	detailsView.SetBorder(true).SetTitle(" Details ").SetTitleAlign(tview.AlignCenter)
	
	serverList = tview.NewList()
	serverList.SetBorder(true).SetTitle(" Servers ").SetTitleAlign(tview.AlignCenter)
	

	imagesList = tview.NewList()
	imagesList.SetBorder(true).SetTitle(" Images ").SetTitleAlign(tview.AlignCenter)
	for _, image := range images.FetchImages() {
		imagesList.AddItem(image.Name, "", -1, func() {
			detailsView.Clear()
			fmt.Fprintf(detailsView, "ID: %s\nName: %s\nSize: %d", image.ID, image.Name, image.SizeBytes)
		})
	}

	flavorsList = tview.NewList()
	flavorsList.SetBorder(true).SetTitle(" Flavors ").SetTitleAlign(tview.AlignCenter)
	for _, flavor := range flavors.FetchFlavors() {
		flavorsList.AddItem(flavor.Name, "", -1, func() {
			detailsView.Clear()
			fmt.Fprintf(detailsView, "ID: %s\nName: %s\nVCPU: %d\nRAM: %d\nDisk: %d", flavor.ID, flavor.Name, flavor.VCPUs, flavor.RAM, flavor.Disk)
		})
	}

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
	
	projectsViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
	AddItem(headerFlex, 0, 1, false).
	AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(projectsList, 0, 1, true).
		AddItem(detailsView, 0, 4, false), 
	0, 5, true)

	populateServersList()

	pages.AddPage("prompt", inputPrompt, true, false)
	pages.AddPage("servers", serverViewFlex, true, true)
	pages.AddPage("images", imageViewFlex, true, true)
	pages.AddPage("flavors", flavorViewFlex, true, true)
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
			flavor, err := json.Marshal(server.Flavor)
			if err != nil {
				flavor = []byte("unable to marshal flavor")
			}
			image, err := json.Marshal(server.Image)
			if err != nil {
				image = []byte("unable to marshal image")
			}
			addresses, err := json.Marshal(server.Addresses)
			if err != nil {
				addresses = []byte("unable to marshal addresses")
			}
			fmt.Fprintf(detailsView, "ID: %s\nStatus: %s\nFlavor: %s\nImage: %s\nNetworks: %s", server.ID, server.Status, flavor, image, addresses)
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

func populateImagesList() {
	imagesList.Clear()
	for _, image := range images.FetchImages() {
		imagesList.AddItem(image.Name, "", -1, func() {
			detailsView.Clear()
			fmt.Fprintf(detailsView, "ID: %s\nName: %s\nSize: %d", image.ID, image.Name, image.SizeBytes)
		})
	}
}
