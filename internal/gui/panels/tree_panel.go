package panels

import (
	"fmt"
	"sort"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/user/pom-manager/internal/core/pom"
)

// TreePanel provides a hierarchical view of the POM structure
type TreePanel struct {
	// UI components
	tree          *widget.Tree
	mainContainer *fyne.Container

	// State
	project       *pom.Project
	treeData      map[string][]string // Parent UID -> Child UIDs
	labelCache    map[string]string   // UID -> Display Label
	labelCacheMux sync.RWMutex        // Protects labelCache from concurrent access

	// Callbacks
	onNodeSelected func(nodeType string, id string)
}

// NewTreePanel creates a new TreePanel
func NewTreePanel() *TreePanel {
	panel := &TreePanel{
		treeData:   make(map[string][]string),
		labelCache: make(map[string]string),
	}

	panel.createUI()
	return panel
}

// createUI creates the tree widget
func (p *TreePanel) createUI() {
	p.tree = widget.NewTree(
		func(uid string) []string {
			return p.treeData[uid]
		},
		func(uid string) bool {
			children, ok := p.treeData[uid]
			return ok && len(children) > 0
		},
		func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(uid string, branch bool, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			// Use cached label instead of computing dynamically
			// Thread-safe read with RLock
			p.labelCacheMux.RLock()
			cachedLabel, ok := p.labelCache[uid]
			p.labelCacheMux.RUnlock()

			if ok {
				label.SetText(cachedLabel)
			} else {
				label.SetText(uid) // Fallback
			}
		},
	)

	p.tree.OnSelected = func(uid string) {
		if p.onNodeSelected != nil {
			nodeType, id := p.parseUID(uid)
			p.onNodeSelected(nodeType, id)
		}
	}

	p.mainContainer = container.NewBorder(
		container.NewVBox(
			widget.NewLabel("POM Structure"),
			widget.NewSeparator(),
		),
		nil, nil, nil,
		p.tree,
	)
}

// LoadProject builds the tree from a project
func (p *TreePanel) LoadProject(project *pom.Project) {
	p.project = project
	p.treeData = make(map[string][]string)

	// Lock the label cache for writing
	p.labelCacheMux.Lock()
	p.labelCache = make(map[string]string)

	if project == nil {
		p.labelCacheMux.Unlock()
		fyne.Do(func() {
			p.tree.Refresh()
		})
		return
	}

	// Root node
	root := fmt.Sprintf("%s:%s:%s", project.GroupID, project.ArtifactID, project.Version)
	p.treeData[""] = []string{root}
	p.labelCache[root] = root

	// Add main sections
	sections := []string{"coordinates", "properties", "dependencies", "plugins", "profiles"}
	p.treeData[root] = sections

	// Cache section labels
	p.labelCache["coordinates"] = "📄 Coordinates"
	p.labelCache["properties"] = fmt.Sprintf("⚙️ Properties (%d)", len(project.Properties))
	p.labelCache["dependencies"] = fmt.Sprintf("📦 Dependencies (%d)", len(project.Dependencies))
	pluginCount := 0
	if project.Build != nil {
		pluginCount = len(project.Build.Plugins)
	}
	p.labelCache["plugins"] = fmt.Sprintf("🔧 Plugins (%d)", pluginCount)
	p.labelCache["profiles"] = fmt.Sprintf("👤 Profiles (%d)", len(project.Profiles))

	// Add dependencies
	if len(project.Dependencies) > 0 {
		depChildren := make([]string, 0, len(project.Dependencies))
		for i := range project.Dependencies {
			uid := fmt.Sprintf("dep:%d", i)
			depChildren = append(depChildren, uid)
			// Cache dependency label
			dep := project.Dependencies[i]
			p.labelCache[uid] = fmt.Sprintf("%s:%s", dep.ArtifactID, dep.Version)
		}
		p.treeData["dependencies"] = depChildren
	}

	// Add plugins
	if project.Build != nil && len(project.Build.Plugins) > 0 {
		pluginChildren := make([]string, 0, len(project.Build.Plugins))
		for i := range project.Build.Plugins {
			uid := fmt.Sprintf("plugin:%d", i)
			pluginChildren = append(pluginChildren, uid)
			// Cache plugin label
			plugin := project.Build.Plugins[i]
			p.labelCache[uid] = plugin.ArtifactID
		}
		p.treeData["plugins"] = pluginChildren
	}

	// Add properties (sorted alphabetically)
	if len(project.Properties) > 0 {
		// Extract and sort keys
		keys := make([]string, 0, len(project.Properties))
		for key := range project.Properties {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		// Build children in sorted order
		propChildren := make([]string, 0, len(keys))
		for _, key := range keys {
			uid := fmt.Sprintf("prop:%s", key)
			propChildren = append(propChildren, uid)
			// Cache property label
			p.labelCache[uid] = fmt.Sprintf("%s = %s", key, project.Properties[key])
		}
		p.treeData["properties"] = propChildren
	}

	// Add profiles
	if len(project.Profiles) > 0 {
		profileChildren := make([]string, 0, len(project.Profiles))
		for i, profile := range project.Profiles {
			uid := fmt.Sprintf("profile:%d", i)
			profileChildren = append(profileChildren, uid)
			// Cache profile label with activation status
			activationStatus := ""
			if profile.Activation != nil && profile.Activation.ActiveByDefault {
				activationStatus = " (active by default)"
			}
			p.labelCache[uid] = fmt.Sprintf("%s%s", profile.ID, activationStatus)
		}
		p.treeData["profiles"] = profileChildren
	}

	// Unlock the label cache before refreshing
	// This ensures all cache writes are complete before any reads
	p.labelCacheMux.Unlock()

	// Refresh must be called on UI thread
	fyne.Do(func() {
		p.tree.Refresh()
	})
}

// getNodeLabel returns the display label for a node
func (p *TreePanel) getNodeLabel(uid string) string {
	if p.project == nil {
		return uid
	}

	switch uid {
	case "coordinates":
		return "📄 Coordinates"
	case "properties":
		return fmt.Sprintf("⚙️ Properties (%d)", len(p.project.Properties))
	case "dependencies":
		return fmt.Sprintf("📦 Dependencies (%d)", len(p.project.Dependencies))
	case "plugins":
		pluginCount := 0
		if p.project.Build != nil {
			pluginCount = len(p.project.Build.Plugins)
		}
		return fmt.Sprintf("🔧 Plugins (%d)", pluginCount)
	case "profiles":
		return fmt.Sprintf("👤 Profiles (%d)", len(p.project.Profiles))
	}

	// Parse specific nodes
	nodeType, id := p.parseUID(uid)
	switch nodeType {
	case "dep":
		// Parse index
		var idx int
		fmt.Sscanf(id, "%d", &idx)
		if idx >= 0 && idx < len(p.project.Dependencies) {
			dep := p.project.Dependencies[idx]
			return fmt.Sprintf("%s:%s", dep.ArtifactID, dep.Version)
		}

	case "plugin":
		var idx int
		fmt.Sscanf(id, "%d", &idx)
		if p.project.Build != nil && idx >= 0 && idx < len(p.project.Build.Plugins) {
			plugin := p.project.Build.Plugins[idx]
			return plugin.ArtifactID
		}

	case "prop":
		if val, ok := p.project.Properties[id]; ok {
			return fmt.Sprintf("%s = %s", id, val)
		}

	case "profile":
		var idx int
		fmt.Sscanf(id, "%d", &idx)
		if idx >= 0 && idx < len(p.project.Profiles) {
			profile := p.project.Profiles[idx]
			activationStatus := ""
			if profile.Activation != nil && profile.Activation.ActiveByDefault {
				activationStatus = " (active by default)"
			}
			return fmt.Sprintf("%s%s", profile.ID, activationStatus)
		}
	}

	return uid
}

// parseUID extracts node type and ID from UID
func (p *TreePanel) parseUID(uid string) (nodeType string, id string) {
	if uid == "coordinates" || uid == "properties" || uid == "dependencies" || uid == "plugins" || uid == "profiles" {
		return uid, ""
	}

	// Format: "type:id"
	for i, c := range uid {
		if c == ':' {
			return uid[:i], uid[i+1:]
		}
	}

	return uid, ""
}

// OnNodeSelected sets the callback for node selection
func (p *TreePanel) OnNodeSelected(callback func(nodeType string, id string)) {
	p.onNodeSelected = callback
}

// GetContainer returns the main container for embedding
func (p *TreePanel) GetContainer() *fyne.Container {
	return p.mainContainer
}
