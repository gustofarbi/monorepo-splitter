package action

import (
	"fmt"
	"path/filepath"
	"splitter/pkg"
)

type WriteChanges struct{}

func (w WriteChanges) Act(collection *pkg.PackageCollection) error {
	for _, singlePkg := range collection.Packages {
		if err := singlePkg.Composer.WriteToFile(filepath.Join(
			collection.RootPackage.Path,
			singlePkg.Path,
			"composer.json",
		)); err != nil {
			return fmt.Errorf("writing changes to singlePkg %s failed: %s", singlePkg.Path, err)
		}
	}

	return nil
}

func (w WriteChanges) Description() string {
	return "write changes made to composer.jsons"
}

func (w WriteChanges) String() string {
	return "write-changes"
}
