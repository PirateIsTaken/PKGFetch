package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"pkgfetch/Globals"
	"pkgfetch/Logger"
)

type GithubRepo struct {
	Name        string `json:"full_name"`
	Description string `json:"description"`
	ReleaseURL  string `json:"releases_url"`

	Stars uint `json:"stargazers_count"`
	Forks uint `json:"forks_count"`

	IsFork      bool `json:"fork"`
	IsArchived  bool `json:"archived"`
	HasReleases bool `json:"-"`

	Score uint `json:"-"`
}

type GithubAsset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        uint64 `json:"size"`
}

type GithubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`

	Assets []GithubAsset `json:"assets"`
}

type GithubSearchResponse struct {
	Items []GithubRepo `json:"items"`
}

func SearchGithub(searchName string) []GithubRepo {
	searchURL := fmt.Sprintf(
		"https://api.github.com/search/repositories?q=%s",
		url.QueryEscape(searchName),
	)
	resp, err := http.Get(searchURL)

	if err != nil {
		Logger.LogError("Failed To Contact Github: \n%v", err)
		Logger.LogNewLine()
		return nil
	}

	defer resp.Body.Close()

	var result GithubSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		Logger.LogError("Failed To Get Info From Github Response: \n%v", err)
		Logger.LogNewLine()
		return nil
	}

	Logger.LogMessageSameLine("Total Packages Found: %d.", len(result.Items))
	Logger.LogMessageSameLine(" Use --show-all to show all the packages")
	Logger.LogNewLine()

	return result.Items
}

// Mechanics
func CalculateScoreGithub(repo *GithubRepo) {
	var score uint = 0

	score += repo.Stars / 100
	score += repo.Forks

	if !repo.IsArchived {
		score += 500
	}
	if !repo.IsFork {
		score += 500
	}

	repo.Score = score
}

func InstallPkgGithub(repo GithubRepo) {
	Logger.LogMessage("Installing: %s", repo.Name)
	Logger.LogNewLine()
	releases := CheckReleasesGithub(&repo)

	if !repo.HasReleases {
		Logger.LogMessage("The Selected Repo Doesn't Have Any Releases. Meaning, It Can't Be Installed With %s", Globals.PROGRAM_NAME)
		return
	}

	// SELECTED RELEASES //
	Logger.LogMessage("Available Releases: ")
	for index, rel := range releases {
		Logger.LogMessage("%d. Version: %s | Tag: %s", index+1, rel.Name, rel.TagName)
	}
	Logger.LogNewLine()
	Logger.LogMessage("Select Version (1...<last_num>): ")
	choice, ok := Logger.ChooseDialog(uint(len(releases)))
	if !ok {
		return
	}

	selectedVersion := releases[choice-1]

	Logger.LogMessage("Selected Version: %s | Tag: %s", selectedVersion.Name, selectedVersion.TagName)
	Logger.LogNewLine()

	// SELLECTING ASSET //
	Logger.LogMessage("Available Assets: ")
	for index, asset := range selectedVersion.Assets {
		Logger.LogMessage("%d. %s | Size: %dMB", index+1, asset.Name, asset.Size/(1024*1024))
	}
	Logger.LogNewLine()
	Logger.LogMessage("Select Asset (1...<last_num>): ")
	choice, ok = Logger.ChooseDialog(uint(len(selectedVersion.Assets)))
	if !ok {
		return
	}

	selectedAsset := selectedVersion.Assets[choice-1]

	Logger.LogMessage("Selected Asset: %s", selectedAsset.Name)
	Logger.LogMessage("Downloading From: %s", selectedAsset.DownloadURL)
	Logger.LogMessage("Size: %dMB", selectedAsset.Size/(1024*1024))

	Logger.LogMessage("Press Enter To Confirm ")
	fmt.Scanln()

	Logger.LogMessage("WAIAJWIHNAW")
}

// Helpers
func CheckReleasesGithub(repo *GithubRepo) []GithubRelease {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/releases",
		repo.Name,
	)

	resp, err := http.Get(url)
	if err != nil {
		repo.HasReleases = false
		return nil
	}

	repo.HasReleases = (resp.StatusCode == http.StatusOK)

	var relResult []GithubRelease
	err = json.NewDecoder(resp.Body).Decode(&relResult)
	if err != nil {
		Logger.LogError("Failed To Fetch Releases Information From Github: \n%v", err)
		Logger.LogNewLine()
		return nil
	}

	defer resp.Body.Close()

	relResult = relResult[:10]
	return relResult
}

func GetPkgListGithub(pkgName string) []GithubRepo {
	repos := SearchGithub(pkgName)

	// Score All The Repos
	for i := range repos {
		CalculateScoreGithub(&repos[i])
	}
	SortPkgOnScoreGithub(repos)

	if len(repos) > 10 {
		repos = repos[:10]
	}

	return repos
}

func SortPkgOnScoreGithub(repos []GithubRepo) {
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Score > repos[j].Score
	})
}
