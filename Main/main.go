package main

import (
	"fmt"
	"os"
	"sort"

	"pkgfetch/Globals"
	"pkgfetch/Logger"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  pkgf <command>")
		fmt.Println("  pkgf <command> <args>")
		fmt.Println(Logger.LOG_MISSUSE_COMMAND)
		return
	}

	command := os.Args[1]

	switch command {
	case "help":
		ShowHelpDialog()
	case "search":
		HandleSearch()
	case "install":
		HandleInstall()
	default:
		fmt.Printf("Unknown Command: %s \n%s\n", command, Logger.LOG_MISSUSE_COMMAND)
	}
}

// Mechanics
func HandleSearch() {
	if argument, ok := IsArgumentGiven(); ok {
		fmt.Printf("Searching For: %s\n", argument)
		repos := GetPkgListGithub(argument)

		if len(repos) == 0 {
			fmt.Printf("Didn't Find Any Packages With The Name: %s\n", argument)
		}

		for index, repo := range repos {
			fmt.Printf("Name: %s | Score: %d\n", repo.Name, repo.Score)
			fmt.Printf("  IsFork: %t | IsArvhived: %t\n", repo.IsFork, repo.IsArchived)
			if !repo.HasReleases {
				fmt.Printf("!!! This Repo Has No Releases. Meaning, It Can't Be Installed Using %s !!!\n", Globals.PROGRAM_NAME)
			}

			if index != len(repos) {
				fmt.Println()
			}
		}
	}
}

func HandleInstall() {
	if argument, ok := IsArgumentGiven(); ok {
		fmt.Printf("Installing: %s\n", argument)
	}
}

func GetPkgListGithub(pkgName string) []Repository {
	repos := SearchGithub(pkgName)
	askedRepos := SearchForPkg(pkgName, repos)

	// Score All The Repos
	for i := range askedRepos {
		CalculateScore(&askedRepos[i])
	}
	SortPkgOnScore(askedRepos)

	return askedRepos
}

func SortPkgOnScore(repos []Repository) {
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Score > repos[j].Score
	})
}

// Helpers
func IsArgumentGiven() (string, bool) {
	if len(os.Args) < 3 {
		fmt.Printf("No Argument Specified. \n%s", Logger.LOG_MISSUSE_COMMAND)
		return "", false
	}
	return os.Args[2], true
}

// @TODO: Put this in Logger package
func ShowHelpDialog() {
	fmt.Println("Usage:")
	fmt.Println("  pkgf <command>")
	fmt.Println("  pkgf <command> <args>")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  help")
	fmt.Println("    Show this dialog")
	println()
	fmt.Println("  install <package>")
	fmt.Println("    Installs given package name (if exists)")
	println()
	fmt.Println("  search <package name>")
	fmt.Println("    Searches for the given package name")
}
