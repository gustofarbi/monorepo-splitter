package action

import (
	"log"
	"path/filepath"
	"splitter/utils/pkg"
)

type WriteChanges struct{}

func (w WriteChanges) Act(collection *pkg.PackageCollection) {
	for _, singlePkg := range collection.Packages {
		err := singlePkg.Composer.WriteToFile(filepath.Join(
			collection.RootPackage.Path,
			singlePkg.Path,
			"composer.json",
		))
		if err != nil {
			log.Fatalf("writing changes to singlePkg %s failed: %s", singlePkg.Path, err)
		}
	}
}

func (w WriteChanges) Description() string {
	return "write changes made to composer.jsons"
}
