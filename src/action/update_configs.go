package action

import (
	"splitter/composer"
	"splitter/pkg"
)

type UpdateConfigs struct{}

func (u UpdateConfigs) Act(collection *pkg.PackageCollection) {
	for _, singlePkg := range collection.Packages {
		delete(singlePkg.Composer.Items.Config, composer.VendorDir)
	}
}

func (u UpdateConfigs) Description() string {
	return "update package composer.jsons (remove .config.vendor-dir)"
}

func (u UpdateConfigs) String() string {
	return "update-configs"
}
