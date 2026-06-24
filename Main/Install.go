package main

import (
	"fmt"
	"os"
	"os/exec"
	"pkgfetch/Logger"
	"strings"
)

// Mechanics
func CheckAndInstall(file string) {
	name := strings.ToLower(file)

	switch {
	case strings.HasSuffix(name, ".rpm"):
		InstallRpm(file)
	case strings.HasSuffix(name, ".deb"):
		InstallDeb(file)

	case strings.HasSuffix(name, ".appimage"):
		InstallAppimage(file)

	case strings.HasSuffix(name, ".tar.gz"):
		InstallArchive(file)
	case strings.HasSuffix(name, ".tar.xz"):
		InstallArchive(file)
	case strings.HasSuffix(name, ".zip"):
		InstallArchive(file)
	}
}

// Helpers
func InstallRpm(file string) {
	Logger.LogNewLine()
	Logger.LogMessage("Installing Package Through dnf...")
	Logger.LogNewLine()

	cmd := exec.Command(
		"sudo",
		"dnf",
		"install",
		file,
	)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()

	if err == nil {
		Logger.LogNewLine()
		Logger.LogMessage("Package Installed Successfully.")
		Logger.LogMessage("To Remove The Package Use `dnf remove <package_name>`")
		Logger.LogNewLine()
		AskToDeleteCache(file)
	} else {
		Logger.LogError("Failed To Install Package From: %s \nBecause %s", file, err)
	}
}

func InstallDeb(file string) {

}

func InstallAppimage(file string) {

}

func InstallArchive(file string) {

}

func AskToDeleteCache(file string) {
	var deleteCache bool = false
	var optionSet bool = false

	for optionSet == false {
		Logger.LogMessage("Do You Want To Delete The Downloaded Cache %s? \n[Y/N]: ", file)

		var choice string
		fmt.Scanln(&choice)
		choice = strings.ToLower(choice)

		switch choice {
		case "y":
			deleteCache = true
			optionSet = true
		case "n":
			deleteCache = false
			optionSet = true
		default:
			Logger.LogError("Please Enter Either 'Y' or 'N'")
		}
	}

	if deleteCache {
		err := os.Remove(file)

		if err != nil {
			Logger.LogError("Failed To Delete Downloaded Cache %s. \nBecause: %v", file, err)
		}
	}
}
