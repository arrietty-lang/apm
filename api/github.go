package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GithubUser struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type GithubAssets struct {
	Url                string      `json:"url"`
	Id                 int         `json:"id"`
	NodeId             string      `json:"node_id"`
	Name               string      `json:"name"`
	Label              string      `json:"label"`
	Uploader           *GithubUser `json:"uploader"`
	ContentType        string      `json:"content_type"`
	State              string      `json:"state"`
	Size               int         `json:"size"`
	DownloadCount      int         `json:"download_count"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	BrowserDownloadUrl string      `json:"browser_download_url"`
}

type GithubRelease struct {
	Url             string          `json:"url"`
	AssetsUrl       string          `json:"assets_url"`
	UploadUrl       string          `json:"upload_url"`
	HtmlUrl         string          `json:"html_url"`
	Id              int             `json:"id"`
	Author          *GithubUser     `json:"author"`
	NodeId          string          `json:"node_id"`
	TagName         string          `json:"tag_name"`
	TargetCommitish string          `json:"target_commitish"`
	Name            string          `json:"name"`
	Draft           bool            `json:"draft"`
	Prerelease      bool            `json:"prerelease"`
	CreatedAt       time.Time       `json:"created_at"`
	PublishedAt     time.Time       `json:"published_at"`
	Assets          []*GithubAssets `json:"assets"`
	TarballUrl      string          `json:"tarball_url"`
	ZipballUrl      string          `json:"zipball_url"`
	Body            string          `json:"body"`
}

type SimpleRelease struct {
	TagName  string
	TarGzUrl string
}

func GetGitHubReleases(author, repoName string) ([]*SimpleRelease, error) {
	var githubReleases []*GithubRelease
	var releases []*SimpleRelease
	for i := 1; i < 100; i++ {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?page=%d", author, repoName, i)

		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(bodyBytes, &githubReleases)
		if err != nil {
			return nil, err
		}

		if len(githubReleases) == 0 {
			break
		}

		for _, gRelease := range githubReleases {
			releases = append(releases, &SimpleRelease{
				TagName:  gRelease.TagName,
				TarGzUrl: gRelease.TarballUrl,
			})
		}
		_ = resp.Body.Close()
	}

	return releases, nil
}

func GetGithubLatestRelease(author, repoName string) (*SimpleRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", author, repoName)

	var githubLatestRelease GithubRelease

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &githubLatestRelease)
	if err != nil {
		return nil, err
	}

	return &SimpleRelease{
		TagName:  githubLatestRelease.TagName,
		TarGzUrl: githubLatestRelease.TarballUrl,
	}, nil
}

func GetGithubReleaseSpecificVersion(author, repoName, version string) (*SimpleRelease, error) {
	exist, err := ExistsGithubRelease(author, repoName, version)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("repository of specific version not found: %v/%v @ %v", author, repoName, version)
	}
	// 特定のバージョンを探すAPIを調べるのがめんどいので、いったん全部取得して、一致するものを返すようにします
	// todo : まともなものを作る
	releases, err := GetGitHubReleases(author, repoName)
	if err != nil {
		return nil, err
	}
	for _, rel := range releases {
		if rel.TagName == version {
			return rel, nil
		}
	}

	return nil, fmt.Errorf("unknown err")
}

func ExistsGithubRelease(author, repoName, tagName string) (bool, error) {
	repos, err := GetGitHubReleases(author, repoName)
	if err != nil {
		return false, err
	}
	for _, repo := range repos {
		if repo.TagName == tagName {
			return true, nil
		}
	}
	return false, nil
}
