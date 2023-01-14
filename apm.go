package apm

import (
	"fmt"
	"github.com/arrietty-lang/apm/api"
	"os"
	"path/filepath"
	"strings"
)

// GetApmPath パッケージマネージャのルートパスを取得
func GetApmPath() string {
	return os.Getenv("ARRIETTY_PM_PATH")
}

func GetApmPackagesPath() string {
	return filepath.Join(GetApmPath(), "packages")
}

func Get(repoUrl string) error {
	// バージョン指定があるかチェック
	var url string
	var version string
	if strings.Contains(repoUrl, "@") {
		ss := strings.Split(repoUrl, "@")
		url = ss[0]
		version = ss[1]
	} else {
		url = repoUrl
		version = ""
	}

	ss := strings.Split(url, "/")
	host := ss[0]
	author := ss[1]
	repoName := ss[2]

	var repo *api.SimpleRelease
	// 各ホストごとに分ける
	switch host {
	case "github.com":
		if version == "" {
			// latest
			r, err := api.GetGithubLatestRelease(author, repoName)
			if err != nil {
				return err
			}
			repo = r
		} else {
			r, err := api.GetGithubReleaseSpecificVersion(author, repoName, version)
			if err != nil {
				return err
			}
			repo = r
		}
	default:
		return fmt.Errorf("unsupported host: %v", url)
	}

	if IsRepositorySpecificVersionInstalled(host, author, repoName, version) {
		return fmt.Errorf("already installed")
	}
	err := InstallTarGz(repo.TarGzUrl, host, author, repoName, version)
	if err != nil {
		return err
	}

	return nil
}
