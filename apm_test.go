package apm

import (
	"fmt"
	"github.com/arrietty-lang/apm/api"
	"testing"
)

var host string
var author string
var repoName string
var repoVersion string

func init() {
	host = "github.com"
	author = "x0y14"
	repoName = "arrietty_json"
	repoVersion = "v0.0.1"
}

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

func TestGetRepository(t *testing.T) {
	repo, err := GetRepositoryInstalled(host, author, fmt.Sprintf("%s@%s", repoName, repoVersion))
	if err != nil {
		t.Fatalf("failed to get repo: %v", err)
	}
	fmt.Println(repo)
}
