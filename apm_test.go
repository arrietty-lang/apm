package apm

import (
	"github.com/arrietty-lang/apm/api"
	"testing"
)

func TestInstallTarGz(t *testing.T) {
	release, err := api.GetGithubLatestRelease("x0y14", "arrietty_json")
	if err != nil {
		t.Fatalf("%v", err)
	}
	err = InstallTarGz(release.TarGzUrl, "github.com", "x0y14", "arrietty_json", release.TagName)
	if err != nil {
		t.Fatalf("%v", err)
	}
}
