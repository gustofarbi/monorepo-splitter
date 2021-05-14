package action

import (
	"log"
	"splitter/pkg"
)

var actionMap map[string]Action

func init() {
	actionMap = make(map[string]Action)
	for _, action := range []Action{
		Validate{},
		SetPackagesDependencies{},
		WriteChanges{},
		CommitChanges{},
		SplitPackages{},
		Reset{},
		UpdateConfigs{},
	} {
		actionMap[action.String()] = action
	}
}

type pipeline struct {
	actions []Action
}

func NewPipeline(names []string) *pipeline {
	p := &pipeline{actions: make([]Action, len(names))}
	for i, name := range names {
		if action, ok := actionMap[name]; ok {
			p.actions[i] = action
		} else {
			panic("unknown action in config: " + name)
		}
	}
	return p
}

func (p pipeline) Run(collection *pkg.PackageCollection) {
	for _, action := range p.actions {
		log.Println(action.Description())
		action.Act(collection)
	}
}
