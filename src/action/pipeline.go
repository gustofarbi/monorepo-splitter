package action

import (
	"log"
	"splitter/utils/pkg"
)

type pipeline struct {
	actions []Action
}

func NewPipeline(actions ...Action) *pipeline {
	p := &pipeline{actions: make([]Action, len(actions))}
	for i, action := range actions {
		p.actions[i] = action
	}
	return p
}

func (p pipeline) Act(collection *pkg.PackageCollection) {
	for _, action := range p.actions {
		log.Println(action.Description())
		action.Act(collection)
	}
}

func (p pipeline) Description() string {
	return ""
}

