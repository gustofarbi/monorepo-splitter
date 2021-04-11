package action

import (
	"log"
	"splitter/pkg"
)

type SetPackagesDependencies struct{}

func (s SetPackagesDependencies) Act(collection *pkg.PackageCollection) {
	versionString := collection.Conf.Semver.CaretedVersion()
	for _, singlePkg := range collection.Packages {
		for name := range singlePkg.Composer.Items.Require {
			if _, ok := collection.Packages[name]; ok {
				singlePkg.Composer.Items.Require[name] = versionString
			} else if currentVersion, ok := collection.RootPackage.Composer.Items.Require[name]; ok {
				singlePkg.Composer.Items.Require[name] = currentVersion
			} else {
				log.Fatalf("singlePkg %s not found locally or in root", name)
			}
		}
	}
}

func (s SetPackagesDependencies) Description() string {
	return "set versions of mutual dependencies to current version"
}

func (s SetPackagesDependencies) String() string {
return "set-packages-dependencies"
}
