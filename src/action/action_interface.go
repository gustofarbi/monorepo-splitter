package action

import (
	"splitter/utils/pkg"
)

type Action interface {
	Act(collection *pkg.PackageCollection)
	Description() string
	//HandleFunc() func(singlePackage *pkg.Package)
}
