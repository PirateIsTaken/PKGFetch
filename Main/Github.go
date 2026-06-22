package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"

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
		Logger.LogError("Failed To Contact Github: %v", err)
		return nil
	}

	defer resp.Body.Close()

	var result GithubSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		Logger.LogError("Failed To Get Info From Github Response: %v", err)
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
	if repo.HasReleases {
		score += 1000
	}

	repo.Score = score
}

func CheckReleasesGithub(repo *GithubRepo) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/releases/latest",
		repo.Name,
	)

	resp, err := http.Get(url)
	if err != nil {
		repo.HasReleases = false
		return
	}
	defer resp.Body.Close()

	repo.HasReleases = (resp.StatusCode == http.StatusOK)
}

// Helpers
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
