package action

import (
	"fmt"
	"splitter/pkg"
)

var actionMap map[string]Action

type pipeline struct {
	actions []Action
}

func NewPipeline(names []string, dryRun bool) (*pipeline, error) {
	initActionMap(dryRun)
	p := &pipeline{actions: make([]Action, len(names))}
	for i, name := range names {
		if action, ok := actionMap[name]; ok {
			p.actions[i] = action
		} else {
			return nil, fmt.Errorf("unknown action in config: %s", name)
		}
	}
	return p, nil
}

func (p pipeline) Run(collection *pkg.PackageCollection) error {
	for _, action := range p.actions {
		fmt.Println(action.Description())
		if err := action.Act(collection); err != nil {
			return fmt.Errorf("error running action %s: %s", action, err)
		}
	}

	return nil
}

func initActionMap(dryRun bool) {
	actionMap = make(map[string]Action)
	for _, action := range []Action{
		Validate{},
		SetPackagesDependencies{},
		WriteChanges{},
		CommitChanges{},
		SplitPackages{dryRun: dryRun},
		Reset{dryRun: dryRun},
		UpdateConfigs{},
	} {
		actionMap[action.String()] = action
	}
}
