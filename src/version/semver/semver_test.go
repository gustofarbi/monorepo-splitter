package semver

import (
	"testing"
)

func TestFromString(t *testing.T) {
	tag := "refs/head/v1.5.56"
	semver, err := FromTag(tag)
	if err != nil {
		t.Fatal(err)
	}
	if semver.IntVal() != 100050056 {
		t.Fatalf("wrong value")
	}
	if semver.String() != "1.5.56" {
		t.Fatalf("wrong value")
	}
}

func TestFromStringBeta(t *testing.T) {
	tag := "refs/head/v1.5.56-BETA"
	semver, err := FromTag(tag)
	if err != nil {
		t.Fatal(err)
	}
	if semver.IntVal() != 100050056 {
		t.Fatalf("wrong value")
	}
	if semver.String() != "1.5.56-BETA" {
		t.Fatalf("wrong value")
	}
}

func TestSemver_CaretedMinorVersion(t *testing.T) {
	tag := "refs/head/v1.5.56"
	semver, err := FromTag(tag)
	if err != nil {
		t.Fatal(err)
	}
	if semver.CaretedMinorVersion() != "^1.5" {
		t.Fatalf("wrong careted minor version: %s", semver.CaretedMinorVersion())
	}
}
