package action

import (
	"fmt"
	"splitter/pkg"
)

type Action interface {
	Act(collection *pkg.PackageCollection)
	Description() string
	fmt.Stringer
}
