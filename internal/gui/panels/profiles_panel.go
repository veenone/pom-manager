package panels

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
)

// ProfilesPanel displays Maven build profiles
type ProfilesPanel struct {
	// UI components
	profilesList  *widget.List
	detailsCard   *widget.Card
	detailsText   *widget.RichText
	mainContainer *fyne.Container

	// State
	profiles      []pom.Profile
	selectedIndex int
}

// NewProfilesPanel creates a new ProfilesPanel
func NewProfilesPanel() *ProfilesPanel {
	panel := &ProfilesPanel{
		selectedIndex: -1,
	}

	panel.createUI()
	return panel
}

// createUI creates the panel layout
func (p *ProfilesPanel) createUI() {
	// Create profiles list
	p.profilesList = widget.NewList(
		func() int {
			return len(p.profiles)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if int(id) < len(p.profiles) {
				profile := p.profiles[id]
				activationStatus := ""
				if profile.Activation != nil && profile.Activation.ActiveByDefault {
					activationStatus = " ✓"
				}
				label.SetText(fmt.Sprintf("%s%s", profile.ID, activationStatus))
			}
		},
	)

	p.profilesList.OnSelected = func(id widget.ListItemID) {
		p.selectedIndex = int(id)
		p.showProfileDetails(int(id))
	}

	// Create details view
	p.detailsText = widget.NewRichText()
	p.detailsCard = widget.NewCard("Profile Details", "", p.detailsText)

	// Create split layout
	split := container.NewHSplit(
		container.NewBorder(
			container.NewVBox(
				widget.NewLabel("Profiles"),
				widget.NewSeparator(),
			),
			nil, nil, nil,
			p.profilesList,
		),
		container.NewScroll(p.detailsCard),
	)
	split.SetOffset(0.3)

	p.mainContainer = container.NewMax(split)
}

// LoadProfiles updates the panel with profiles
func (p *ProfilesPanel) LoadProfiles(profiles []pom.Profile) {
	p.profiles = profiles
	p.selectedIndex = -1

	fyne.Do(func() {
		p.profilesList.Refresh()
		p.detailsText.ParseMarkdown("*Select a profile to view details*")
		p.detailsCard.SetSubTitle("")
	})
}

// showProfileDetails displays details for the selected profile
func (p *ProfilesPanel) showProfileDetails(index int) {
	if index < 0 || index >= len(p.profiles) {
		return
	}

	profile := p.profiles[index]

	// Build details markdown
	details := fmt.Sprintf("# Profile: %s\n\n", profile.ID)

	// Activation info
	if profile.Activation != nil {
		details += "## Activation\n\n"
		if profile.Activation.ActiveByDefault {
			details += "- **Active by Default**: Yes\n"
		} else {
			details += "- **Active by Default**: No\n"
		}

		if profile.Activation.JDK != "" {
			details += fmt.Sprintf("- **JDK**: %s\n", profile.Activation.JDK)
		}

		if profile.Activation.Property != nil {
			details += fmt.Sprintf("- **Property**: %s", profile.Activation.Property.Name)
			if profile.Activation.Property.Value != "" {
				details += fmt.Sprintf(" = %s", profile.Activation.Property.Value)
			}
			details += "\n"
		}

		if profile.Activation.OS != nil {
			details += "- **OS**:\n"
			if profile.Activation.OS.Family != "" {
				details += fmt.Sprintf("  - Family: %s\n", profile.Activation.OS.Family)
			}
			if profile.Activation.OS.Name != "" {
				details += fmt.Sprintf("  - Name: %s\n", profile.Activation.OS.Name)
			}
			if profile.Activation.OS.Arch != "" {
				details += fmt.Sprintf("  - Arch: %s\n", profile.Activation.OS.Arch)
			}
		}

		details += "\n"
	}

	// Properties
	if len(profile.Properties) > 0 {
		details += fmt.Sprintf("## Properties (%d)\n\n", len(profile.Properties))
		for key, value := range profile.Properties {
			details += fmt.Sprintf("- **%s**: `%s`\n", key, value)
		}
		details += "\n"
	}

	// Dependencies
	if len(profile.Dependencies) > 0 {
		details += fmt.Sprintf("## Dependencies (%d)\n\n", len(profile.Dependencies))
		for _, dep := range profile.Dependencies {
			details += fmt.Sprintf("- %s:%s:%s", dep.GroupID, dep.ArtifactID, dep.Version)
			if dep.Scope != "" && dep.Scope != "compile" {
				details += fmt.Sprintf(" (%s)", dep.Scope)
			}
			details += "\n"
		}
		details += "\n"
	}

	// Build
	if profile.Build != nil && len(profile.Build.Plugins) > 0 {
		details += fmt.Sprintf("## Build - Plugins (%d)\n\n", len(profile.Build.Plugins))
		for _, plugin := range profile.Build.Plugins {
			details += fmt.Sprintf("- %s:%s", plugin.GroupID, plugin.ArtifactID)
			if plugin.Version != "" {
				details += fmt.Sprintf(":%s", plugin.Version)
			}
			if len(plugin.Executions) > 0 {
				details += fmt.Sprintf(" (%d executions)", len(plugin.Executions))
			}
			details += "\n"
		}
		details += "\n"
	}

	// Modules
	if len(profile.Modules) > 0 {
		details += fmt.Sprintf("## Modules (%d)\n\n", len(profile.Modules))
		for _, module := range profile.Modules {
			details += fmt.Sprintf("- %s\n", module)
		}
		details += "\n"
	}

	fyne.Do(func() {
		p.detailsText.ParseMarkdown(details)
		subtitle := fmt.Sprintf("%d properties, %d dependencies, %d modules",
			len(profile.Properties), len(profile.Dependencies), len(profile.Modules))
		p.detailsCard.SetSubTitle(subtitle)
	})
}

// GetContainer returns the main container for embedding
func (p *ProfilesPanel) GetContainer() *fyne.Container {
	return p.mainContainer
}
