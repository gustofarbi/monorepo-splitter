package version

import "testing"

func TestFromString(t *testing.T) {
	tag := "refs/head/v1.5.56"
	semver := FromTag(tag)
	if semver.IntVal() != 100050056 {
		t.Fatalf("wrong value")
	}
	if semver.String() != "1.5.56" {
		t.Fatalf("wrong value")
	}
}

func TestSemver_CaretedMinorVersion(t *testing.T) {
	tag := "refs/head/v1.5.56"
	semver := FromTag(tag)
	if semver.CaretedMinorVersion() != "^1.5" {
		t.Fatalf("wrong careted minor version: %s", semver.CaretedMinorVersion())
	}
}
