package action

import (
	"fmt"
	"splitter/pkg"
)

type Action interface {
	Act(collection *pkg.PackageCollection) error
	Description() string
	fmt.Stringer
}
