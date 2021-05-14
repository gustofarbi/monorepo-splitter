package version

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	width     = 4
	partCount = 3
)

var format = fmt.Sprintf("%%0%dd%%0%dd%%0%dd", width, width, width)

type Semver struct {
	major int
	minor int
	patch int
}

func (s Semver) String() string {
	return fmt.Sprintf("%d.%d.%d", s.major, s.minor, s.patch)
}

func (s Semver) CaretedMinorVersion() string {
	return fmt.Sprintf("^%d.%d", s.major, s.minor)
}
func (s Semver) GitTag() string {
	return "v" + s.String()
}
func (s Semver) IntVal() int {
	value, _ := strconv.Atoi(fmt.Sprintf(format, s.major, s.minor, s.patch))
	return value
}

func FromString(s string) Semver {
	strList := strings.Split(s, ".")
	if len(strList) != partCount {
		panic("wrong semantic version format: " + s)
	}
	intList := make([]int, partCount)
	var res int64
	var err error

	for i, v := range strList {
		if res, err = strconv.ParseInt(v, 10, 64); err != nil {
			log.Fatalf("error parsing version from %s: %s", s, err)
		}
		intList[i] = int(res)
	}

	return Semver{
		major: intList[0],
		minor: intList[1],
		patch: intList[2],
	}
}

func FromTag(tag string) Semver {
	base := filepath.Base(tag)
	if strings.HasPrefix(base, "v") {
		base = base[1:]
	}
	return FromString(base)
}

func (s Semver) IsGreater(b Semver) bool {
	if s.major == b.major {
		if s.minor == b.minor {
			return s.patch > b.patch
		}

		return s.minor > b.minor
	}

	return s.major > b.major
}

type SemverCollection struct {
	versions map[Semver]bool
}

func NewSemverCollection() *SemverCollection {
	return &SemverCollection{
		versions: make(map[Semver]bool),
	}
}

func (c SemverCollection) Add(s Semver) {
	c.versions[s] = true
}

func (c SemverCollection) GetHighest() Semver {
	highest := Semver{}
	for s := range c.versions {
		if s.IsGreater(highest) {
			highest = s
		}
	}

	return highest
}
