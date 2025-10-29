package pom

import "fmt"

// Organizer interface for organizing plugin executions
type Organizer interface {
	ByPhase(project *Project) map[string][]PluginExecution
	ByGoal(project *Project) map[string][]PluginExecution
	ByPlugin(project *Project) map[string][]PluginExecution
	GetPhaseOrder() []string
}

// defaultOrganizer implements Organizer
type defaultOrganizer struct{}

// NewOrganizer creates a new Organizer instance
func NewOrganizer() Organizer {
	return &defaultOrganizer{}
}

// ByPhase organizes plugin executions by Maven lifecycle phase
func (o *defaultOrganizer) ByPhase(project *Project) map[string][]PluginExecution {
	result := make(map[string][]PluginExecution)

	if project == nil || project.Build == nil {
		return result
	}

	for _, plugin := range project.Build.Plugins {
		for _, exec := range plugin.Executions {
			if exec.Phase != "" {
				result[exec.Phase] = append(result[exec.Phase], exec)
			}
		}
	}

	return result
}

// ByGoal organizes plugin executions by goal (format: "plugin:goal")
func (o *defaultOrganizer) ByGoal(project *Project) map[string][]PluginExecution {
	result := make(map[string][]PluginExecution)

	if project == nil || project.Build == nil {
		return result
	}

	for _, plugin := range project.Build.Plugins {
		for _, exec := range plugin.Executions {
			for _, goal := range exec.Goals {
				// Format as "plugin:goal"
				goalKey := fmt.Sprintf("%s:%s", plugin.ArtifactID, goal)
				result[goalKey] = append(result[goalKey], exec)
			}
		}
	}

	return result
}

// ByPlugin organizes plugin executions by plugin artifactId
func (o *defaultOrganizer) ByPlugin(project *Project) map[string][]PluginExecution {
	result := make(map[string][]PluginExecution)

	if project == nil || project.Build == nil {
		return result
	}

	for _, plugin := range project.Build.Plugins {
		if len(plugin.Executions) > 0 {
			result[plugin.ArtifactID] = append(result[plugin.ArtifactID], plugin.Executions...)
		}
	}

	return result
}

// GetPhaseOrder returns Maven lifecycle phases in execution order
func (o *defaultOrganizer) GetPhaseOrder() []string {
	return MavenLifecyclePhases
}
