package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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
	Logger.LogMessage("Selected Package: %s", repo.Name)
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
	// Trim assets to show distro specific and common files
	trimmed := TrimAssetsGithub(selectedVersion.Assets)

	if len(trimmed) <= 0 {
		Logger.LogMessage("This Repo Has 0 Releases Supported For Your Platform. Meaning, It Can't Be Downloaded From PKGFetch")
		Logger.LogNewLine()
		return
	}

	Logger.LogMessage("Available Assets: ")
	for index, asset := range trimmed {
		Logger.LogMessage("%d. %s | Size: %dMB", index+1, asset.Name, asset.Size/(1024*1024))
	}
	Logger.LogNewLine()
	Logger.LogMessage("Select Asset (1...<last_num>): ")
	choice, ok = Logger.ChooseDialog(uint(len(trimmed)))
	if !ok {
		return
	}

	selectedAsset := trimmed[choice-1]

	Logger.LogMessage("Selected Asset: %s", selectedAsset.Name)
	Logger.LogMessage("Downloading From: %s", selectedAsset.DownloadURL)
	Logger.LogMessage("Size: %dMB", selectedAsset.Size/(1024*1024))

	Logger.LogMessage("Press Enter To Confirm ")
	fmt.Scanln()

	DownloadAssetGithub(selectedAsset)
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

func TrimAssetsGithub(assets []GithubAsset) []GithubAsset {
	var filtered []GithubAsset

	for _, asset := range assets {
		// Checking For Supported Assets
		supported := false
		for _, ext := range Globals.SupportedAssets {
			if strings.HasSuffix(asset.Name, ext) {
				supported = true
				break
			}
		}
		if !supported {
			continue
		}

		// Checking For UnSupported Assets
		unsupported := false
		for _, word := range Globals.UnsupportedKeywords {
			if strings.Contains(asset.Name, word) {
				unsupported = true
				break
			}
		}
		if unsupported {
			continue
		}

		filtered = append(filtered, asset)
	}
	return filtered
}

func DownloadAssetGithub(asset GithubAsset) {
	resp, err := http.Get(asset.DownloadURL)

	if err != nil {
		Logger.LogError("Failed To Download Asset %s \nBecause: %v", asset.Name, err)
		Logger.LogNewLine()
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		Logger.LogError("Failed To Get User Home Dir")
		Logger.LogNewLine()
		return
	}

	downloadPath := strings.Replace(
		Globals.AppConfig.DownloadPath,
		"~",
		homeDir,
		1,
	)

	filePath := filepath.Join(
		downloadPath,
		asset.Name,
	)

	file, err := os.Create(filePath)
	if err != nil {
		Logger.LogError("Failed To Create File At %s \nBecause: %v", filePath, err)
		Logger.LogNewLine()
		return
	}

	readBuffer := make([]byte, 32*1024)
	var downloadedBytes int64 = 0
	var speedMBPerSec float64
	var speedBytesPerSec float64

	totalBytes := resp.ContentLength
	var remainingBytes int64

	var downloadedMB float64
	var totalMB float64
	var downloadPercentage float64
	var etaTime float64

	lastUpdatedTUI := time.Now()
	startTime := time.Now()
	for {
		n, err := resp.Body.Read(readBuffer)

		if n > 0 {
			file.Write(readBuffer[:n])
			downloadedBytes += int64(n)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			Logger.LogError("Failed To Download File: %s \nBecause: %v", asset.Name, err)
			Logger.LogNewLine()
			return
		}

		elapsedTime := time.Since(startTime).Seconds()
		speedBytesPerSec = float64(downloadedBytes) / elapsedTime
		speedMBPerSec = speedBytesPerSec / (1024 * 1024)

		remainingBytes = totalBytes - downloadedBytes

		downloadedMB = float64(downloadedBytes) / (1024 * 1024)
		totalMB = float64(totalBytes) / (1024 * 1024)
		downloadPercentage = float64(downloadedBytes) * 100 / float64(totalBytes)
		etaTime = float64(remainingBytes) / speedBytesPerSec
		eta := time.Duration(etaTime) * time.Second

		if time.Since(lastUpdatedTUI) >= time.Second {
			lastUpdatedTUI = time.Now()

			fmt.Printf("\rDownloading... %.1f%% | %.1f MB / %.1f MB | %.1f MB/s | ETA: %s",
				downloadPercentage,
				downloadedMB, totalMB,
				speedMBPerSec,
				eta,
			)
		}
	}

	Logger.LogNewLine()
	Logger.LogMessage("Successfully Downloaded \nAsset: %s \nTo Filepath: %s", asset.Name, filePath)

	defer file.Close()
	defer resp.Body.Close()
}
