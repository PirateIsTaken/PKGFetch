package main

import (
	"os"
	"time"

	"pkgfetch/Globals"
	"pkgfetch/Logger"
)

func main() {
	if len(os.Args) < 2 {
		Logger.LogMessageSameLine("Usage:")
		Logger.LogMessage("  %s <command>", Globals.PROGRAM_NAME_CMD)
		Logger.LogMessage("  %s <command> <args>", Globals.PROGRAM_NAME_CMD)
		Logger.LogMessage(Logger.LOG_MISSUSE_COMMAND)
		return
	}

	command := os.Args[1]

	switch command {
	case "help":
		Logger.ShowHelpDialog()
	case "search":
		HandleSearch()
	case "install":
		HandleInstall()
	default:
		Logger.LogWarningSameLine("Unknown Command: %s \n%s", command, Logger.LOG_MISSUSE_COMMAND)
	}
}

// Mechanics
func HandleSearch() {
	start := time.Now()
	if argument, ok := IsArgumentGiven(); ok {
		Logger.LogMessageSameLine("Searching For: %s", argument)
		Logger.LogNewLine()
		repos := GetPkgListGithub(argument)

		if len(repos) == 0 {
			Logger.LogMessage("Didn't Find Any Packages With The Name: %s", argument)
			Logger.LogNewLine()
			return
		}

		for index, repo := range repos {
			Logger.LogMessage("%s [Github Link: %s]", repo.Name, repo.Link)

			// Repo Status
			releaseStatus := ""
			if !repo.HasReleases {
				releaseStatus = "[NO RELEASES] "
			}
			archiveStatus := ""
			if repo.IsArchived {
				archiveStatus = "[ARCHIVED] "
			}
			forkStatus := ""
			if repo.IsFork {
				forkStatus = "[FORK] "
			}
			if repo.IsArchived || !repo.HasReleases || repo.IsFork {
				Logger.LogWarning("%s%s%s", releaseStatus, archiveStatus, forkStatus)
			}

			if index != len(repos) {
				Logger.LogNewLine()
			}
		}
	}

	Logger.LogMessage("Time Took To Search: %v", time.Since(start))
	Logger.LogNewLine()
}

func HandleInstall() {
	if argument, ok := IsArgumentGiven(); ok {
		Logger.LogMessageSameLine("Installing: %s", argument)
	}
}

// Helpers
func IsArgumentGiven() (string, bool) {
	if len(os.Args) < 3 {
		Logger.LogWarningSameLine("No Argument Specified For The Command Used. \n%s", Logger.LOG_MISSUSE_COMMAND)
		Logger.LogNewLine()
		return "", false
	}
	return os.Args[2], true
}
