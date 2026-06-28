package main

import (
	"os"
	"path/filepath"
	"time"

	"pkgfetch/Globals"
	"pkgfetch/Logger"

	"github.com/BurntSushi/toml"
)

func main() {
	if !Setup() {
		Logger.LogError("Failed To Setup PKGFetch. Stopping The Program Here.")
		Logger.LogNewLine()
		return
	}

	if Globals.ConfigPath == "" {
		Logger.LogError("Failed To Apply Config Path, ReRun The Program By Checking If You Have A Config File With The Name `pkgfetch/pkgf.toml` In Your Config Dir")
	}

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
		if os.Args[2] == "--local" {
			if len(os.Args) < 3 {
				Logger.LogError("You Didn't Provide A Local File To Be Installed. \nStopping Installation Process...")
				Logger.LogNewLine()
				return
			}
			CheckAndInstall(os.Args[3])
		} else {
			HandleInstall()
		}
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
			Logger.LogMessage("Couldn't Find Any Packages With The Name: %s", argument)
			Logger.LogNewLine()
			return
		}

		for index, repo := range repos {
			Logger.LogMessage("%s", repo.Name)

			// Repo Status
			archiveStatus := ""
			if repo.IsArchived {
				archiveStatus = "[ARCHIVED] "
			}
			forkStatus := ""
			if repo.IsFork {
				forkStatus = "[FORK] "
			}
			if repo.IsArchived || repo.IsFork {
				Logger.LogWarningSameLine(" | %s%s", archiveStatus, forkStatus)
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
	Logger.LogMessageSameLine("Selected PKGType: %s", Globals.AppConfig.PKGType)

	if argument, ok := IsArgumentGiven(); ok {
		Logger.LogMessage("Searching Package: %s", argument)
		Logger.LogNewLine()

		repos := GetPkgListGithub(argument)

		if len(repos) == 0 {
			Logger.LogMessage("Couldn't Find Any Packages With The Name: %s", argument)
			Logger.LogNewLine()
			return
		}

		for index, repo := range repos {
			Logger.LogMessage("%d. %s", index+1, repo.Name)

			// Repo Status
			archiveStatus := ""
			if repo.IsArchived {
				archiveStatus = "[ARCHIVED] "
			}
			forkStatus := ""
			if repo.IsFork {
				forkStatus = "[FORK] "
			}
			if repo.IsArchived || repo.IsFork {
				Logger.LogWarningSameLine(" | %s%s", archiveStatus, forkStatus)
			}
		}
		Logger.LogNewLine()
		choice := Logger.ChooseDialog(uint(len(repos)), "Select Package (1...<last_num>): ")
		selectedRepo := repos[choice-1]

		InstallPkgGithub(selectedRepo)
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

// return if setup was successful
func Setup() bool {
	configDir, err := os.UserConfigDir()
	if err != nil {
		Logger.LogErrorSameLine("Failed To Load Config Dir: %s \nBecause: %v", configDir, err)
		Logger.LogNewLine()
		return false
	}
	Globals.ConfigPath = filepath.Join(configDir, "pkgfetch", "pkgf.toml")

	// Load Config
	var config Globals.Config

	_, err = toml.DecodeFile(
		Globals.ConfigPath,
		&config,
	)

	if err != nil {
		Logger.LogErrorSameLine("Failed To Load Config From: %s \nBecause: %v", Globals.ConfigPath, err)
		Logger.LogNewLine()
		return false
	}
	Globals.AppConfig = config

	Globals.SupportedAssets = append(Globals.SupportedAssets, string("."+Globals.AppConfig.PKGType))

	return true
}
