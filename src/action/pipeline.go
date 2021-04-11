package action

import (
	"log"
	"splitter/pkg"
)

var (
	validate                = Validate{}
	setPackagesDependencies = SetPackagesDependencies{}
	writeChanges            = WriteChanges{}
	tagRelease              = TagRelease{}
	splitPackages           = SplitPackages{}
)

type pipeline struct {
	actions []Action
}

func NewPipeline(names []string) *pipeline {
	p := &pipeline{actions: make([]Action, len(names))}
	var action Action
	for i, name := range names {
		switch name {
		case "validate":
			action = validate
		case setPackagesDependencies.String():
			action = setPackagesDependencies
		case writeChanges.String():
			action = writeChanges
		case tagRelease.String():
			action = tagRelease
		case splitPackages.String():
			action = splitPackages
		default:
			panic("unknown action: " + name)
		}
		p.actions[i] = action
	}
	return p
}

func (p pipeline) Run(collection *pkg.PackageCollection) {
	for _, action := range p.actions {
		log.Println(action.Description())
		action.Act(collection)
	}
}
