package action

import (
	"splitter/utils/pkg"
)

type UpdateReplaceRelease struct{}

func (u UpdateReplaceRelease) Act(collection *pkg.PackageCollection) {
	return
	// todo do we need this?
	newVersion := collection.Conf.Semver.String()
	for name := range collection.RootPackage.Composer.Items.Replace {
		if _, ok := collection.Packages[name]; ok {
			collection.RootPackage.Composer.Items.Replace[name] = newVersion
		}
	}
}

func (u UpdateReplaceRelease) Description() string {
	return "updates replace in root composer.json"
}
