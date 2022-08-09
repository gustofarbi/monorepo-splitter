package version

type StringVersion struct {
	Version string
}

func (s StringVersion) String() string {
	return s.Version
}

func (s StringVersion) GitTag() string {
	return s.Version
}

func (s StringVersion) CaretedMinorVersion() string {
	return s.Version
}
