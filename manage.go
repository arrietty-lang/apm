package apm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// インストールしたことのあるホスティングサービスの一覧を取得
// ex) github.com, ...
func getHostingServicesInstalled() ([]string, error) {
	apmPath := GetApmPackagesPath()

	var hosts []string

	files, err := os.ReadDir(apmPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() && f.Name() != apmPath {
			hosts = append(hosts, f.Name())
		}
	}

	return hosts, nil
}

// ホスティングサービスにインストールされているリポジトリ作者一覧を取得
// getHostingServicesInstalledで取得したものを使用する
func getRepositoryAuthorsInstalled(hostingService string) ([]string, error) {
	hostDir := filepath.Join(GetApmPackagesPath(), hostingService)

	var authors []string

	files, err := os.ReadDir(hostDir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() && f.Name() != hostDir {
			authors = append(authors, f.Name())
		}
	}

	return authors, nil
}

// ホスティングサービスと作者からリポジトリ一覧を取得する
// `repoName@version`というフォーマットになってると思う
// getRepositoryAuthorsInstalledで取得したものを使う
func getRepositoryNameAtVersionsInstalled(hostingService, author string) ([]string, error) {
	hostDir := filepath.Join(GetApmPackagesPath(), hostingService)
	authorDir := filepath.Join(hostDir, author)

	files, err := os.ReadDir(authorDir)
	if err != nil {
		return nil, err
	}

	var repos []string
	for _, repo := range files {
		if repo.IsDir() {
			repos = append(repos, repo.Name())
		}
	}

	return repos, nil
}

// GetRepositoryInstalled リポジトリの詳細を取得
func GetRepositoryInstalled(hostingService, author, repoNameAtVersion string) (*Repository, error) {
	// 適切なフォーマットか確認
	if !strings.Contains(repoNameAtVersion, "@") {
		return nil, fmt.Errorf("invalid repository id format, expect: 'repoName@repoVersion', but reserve %s", repoNameAtVersion)
	}

	// 存在確認
	pkgJsonPath := filepath.Join(GetApmPackagesPath(), hostingService, author, repoNameAtVersion, "pkg.json")
	if !existsFile(pkgJsonPath) {
		return nil, fmt.Errorf("pkg.json not found: %s", pkgJsonPath)
	}

	// pkg.jsonを解読
	bytes, err := os.ReadFile(pkgJsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %v: %v", pkgJsonPath, err)
	}

	// Pkgとして読み込む
	pkgJson, err := UnmarshalPkgJson(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %v", pkgJsonPath, err)
	}

	// Pkgとして読み込んだものと追加情報を返却
	nameVersion := strings.Split(repoNameAtVersion, "@")
	repoName := nameVersion[0]
	repoVersion := nameVersion[1]
	return &Repository{
		Host:    hostingService,
		Author:  author,
		Name:    repoName,
		Version: repoVersion,
		Deps:    pkgJson.Deps,
	}, nil
}

// GetRepositoriesInstalledByAuthor リポジトリの詳細の一覧を作者から取得
func GetRepositoriesInstalledByAuthor(hostingService, author string) ([]*Repository, error) {
	repoNameAtVersions, err := getRepositoryNameAtVersionsInstalled(hostingService, author)
	if err != nil {
		return nil, fmt.Errorf("failed to get installed repositories: %v", err)
	}

	var repos []*Repository

	for _, repoNV := range repoNameAtVersions {
		repo, err := GetRepositoryInstalled(hostingService, author, repoNV)
		if err != nil {
			return nil, fmt.Errorf("failed to get repository data: %v", err)
		}
		repos = append(repos, repo)
	}

	return repos, nil
}

// GetMultipleVersionRepositoriesInstalledByRepoName リポジトリ名から、複数のバージョンのリポジトリを検索する
// x@0.0.1
// x@0.0.2
// x@0.0.3
// の一覧をxから検索
func GetMultipleVersionRepositoriesInstalledByRepoName(hostingService, author, repoName string) ([]*Repository, error) {
	repos, err := GetRepositoriesInstalledByAuthor(hostingService, author)
	if err != nil {
		return nil, err
	}

	var multipleVersionRepos []*Repository
	for _, repo := range repos {
		if repo.Name == repoName {
			multipleVersionRepos = append(multipleVersionRepos, repo)
		}
	}

	return multipleVersionRepos, nil
}

// IsRepositoryInstalledByRepoName リポジトリ名からインストールされているかを取得する
// バージョンは問わない
func IsRepositoryInstalledByRepoName(hostingService, author, repoName string) bool {
	hostDir := filepath.Join(GetApmPackagesPath(), hostingService)
	if !existsFile(hostDir) {
		return false
	}

	authorDir := filepath.Join(hostDir, author)
	if !existsFile(authorDir) {
		return false
	}

	repos, err := GetRepositoriesInstalledByAuthor(hostingService, author)
	if err != nil {
		return false
	}

	for _, repo := range repos {
		if repo.Name == repoName {
			return true
		}
	}

	return false
}

// IsRepositorySpecificVersionInstalled 特定のバージョンのリポジトリがインストールされているか
func IsRepositorySpecificVersionInstalled(hostingService, author, repoName, repoVersion string) bool {
	if !IsRepositoryInstalledByRepoName(hostingService, author, repoName) {
		return false
	}

	multipleVersions, err := GetMultipleVersionRepositoriesInstalledByRepoName(hostingService, author, repoName)
	if err != nil {
		return false
	}

	for _, repo := range multipleVersions {
		if repo.Version == repoVersion {
			return true
		}
	}

	return false
}

//func UninstallRepository(hostingService, author, repoName string) error {
//
//}
