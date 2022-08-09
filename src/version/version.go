package version

import "fmt"

type Version interface {
	fmt.Stringer
	GitTag() string
	CaretedMinorVersion() string
}
