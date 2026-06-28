package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"pkgfetch/Globals"
	"pkgfetch/Logger"
	"strings"
)

type DesktopEntryProperties struct {
	AppName  string
	IconPath string
	ExecPath string
}

// Mechanics
func CheckAndInstall(file string) {
	name := strings.ToLower(file)

	switch {
	case strings.HasSuffix(name, ".rpm"):
		InstallRpm(file)
		return
	case strings.HasSuffix(name, ".deb"):
		InstallDeb(file)
		return

	case strings.HasSuffix(name, ".appimage"):
		InstallAppimage(file)
		return

	case strings.HasSuffix(name, ".tar.gz"):
		InstallArchive(file)
		return
	case strings.HasSuffix(name, ".tar.xz"):
		InstallArchive(file)
		return
	case strings.HasSuffix(name, ".zip"):
		InstallArchive(file)
		return
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
	Logger.LogMessage(".deb Installation Not Yet... Sorry :)")
	Logger.LogNewLine()
}

func InstallAppimage(file string) {
	appImagePath := ExpandHome(Globals.AppConfig.AppImagePath)
	if appImagePath == "" {
		Logger.LogError("Something Went Wrong While Trying To Expand Home Dir... \nStopping Installation Process...")
		Logger.LogNewLine()
		return
	}

	Logger.LogNewLine()
	Logger.LogMessage("Installing AppImage...")
	Logger.LogNewLine()

	appLocation := filepath.Join(
		appImagePath,
		filepath.Base(file),
	)

	// Create dir if doesn't exist
	err := os.MkdirAll(appImagePath, 0755)
	if err != nil {
		Logger.LogError("Failed To Create Directory: %s \nBecause: %v", appImagePath, err)
		Logger.LogNewLine()
		return
	}
	err = CopyFile(file, appLocation)
	if err != nil {
		Logger.LogError("Failed To Copy AppImage From %s To %s \nBecause: %v", file, appImagePath, err)
		Logger.LogNewLine()
		return
	}

	err = os.Chmod(appLocation, 0755)
	if err != nil {
		Logger.LogError("Failed To Give Executing Permissions To AppImage At: %s \nBecause: %v", appLocation, err)
		Logger.LogNewLine()
		return
	}

	ExtractAppImage(appLocation)
	AskToDeleteCache(file)
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

// HELPERS
func ExtractAppImage(appLocation string) {
	extractPath := ExpandHome(Globals.AppConfig.ExtractPath)
	err := os.MkdirAll(extractPath, 0755)
	if err != nil {
		Logger.LogError("Failed To Create Dirs For Extracting AppImage \nBecause: %v. \nStopping Installation Process...", err)
		Logger.LogNewLine()
		return
	}

	cmd := exec.Command(
		appLocation,
		"--appimage-extract",
	)
	cmd.Dir = extractPath // this will give squashroot-fs
	err = cmd.Run()
	if err != nil {
		Logger.LogError(string("Failed To Extract AppImage Content By Runnign The Command `--appimage-extract`. \nBecause: %v"+
			"\nCheck If The File %s Has +x(0755) Perms With `ls -l`"), err, appLocation)
		Logger.LogNewLine()
		return
	}

	appImgExtractDir := filepath.Join(
		extractPath,
		"squashfs-root",
	)
	desktopFile := FindDesktopFile(appImgExtractDir)
	if desktopFile == "" {
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		Logger.LogError("Failed To Get User Home Dir. \nBecause: %v", err)
		Logger.LogNewLine()
		Logger.LogError("Stopping Installation Process.")
		Logger.LogNewLine()
		return
	}
	applicationsDir := filepath.Join(
		homeDir,
		".local",
		"share",
		"applications",
	)
	newDesktopFile := filepath.Join(
		applicationsDir,
		filepath.Base(desktopFile),
	)
	CopyFile(desktopFile, newDesktopFile)

	desktopEntryInfo := ExtractInfoFromDesktopFile(newDesktopFile, appImgExtractDir)
	desktopEntryInfo.ExecPath = appLocation
	if desktopEntryInfo.AppName == "" {
		desktopEntryInfo.AppName = filepath.Base(appLocation)
	}

	// @TODO: UPDATE SYSTEM. Check If There Is A .desktop File With The Same desktopEntryInfo.AppName

	// Modifying The Desktop File
	SetVariableInFile(newDesktopFile, "Exec", desktopEntryInfo.ExecPath)
	SetVariableInFile(newDesktopFile, "Icon", desktopEntryInfo.IconPath)

	Logger.LogMessage("Cleaning Up...")
	Logger.LogNewLine()

	err = os.RemoveAll(appImgExtractDir)
	if err != nil {
		Logger.LogError("Failed To Delete Extracted Cache At %s. \nBecause: %v", appImgExtractDir, err)
		Logger.LogWarning("You Should Care About That Error And Go Delete It Yourself.")
		Logger.LogNewLine()
		return
	}

	Logger.LogMessage("Successfully Installed %s To Path: %s", desktopEntryInfo.AppName, appLocation)
	Logger.LogNewLine()
}

func ExtractInfoFromDesktopFile(file string, fileRoot string) DesktopEntryProperties {
	var desktopEntryInfo DesktopEntryProperties = DesktopEntryProperties{}

	iconPath := FetchIconFromDesktopFile(file, fileRoot)

	if iconPath != "" {
		configIconPath := ExpandHome(Globals.AppConfig.IconPath)
		newIconPath := filepath.Join(
			configIconPath,
			filepath.Base(iconPath),
		)

		err := os.MkdirAll(configIconPath, 0755)
		if err != nil {
			Logger.LogError("Failed To Create Dir: %s. \nBecause: %v", newIconPath, err)
			Logger.LogWarning("Desktop Entry Will Be Created But With No Icon...")
			Logger.LogNewLine()
		}

		if err == nil {
			err = CopyFile(iconPath, newIconPath)
			if err != nil {
				Logger.LogError(string("Failed To Copy Icon From: %s, To: %s"+
					"\nBecause: %v"+
					"\nDesktop Entry Will Be Created But Without An Icon."),
					iconPath, newIconPath, err,
				)
				Logger.LogNewLine()
			}

			desktopEntryInfo.IconPath = newIconPath
		}
	}

	// Finding The App Name
	appName := FindVariableInFile(file, "Name=")
	if appName == "" {
		Logger.LogWarning("No App Name Was Found... This Is Unusual. Though The Desktop Entry Will Be Created... It Will Have A Complex Name. \nDesktop Entry Path: %s", file)
		Logger.LogNewLine()
	} else {
		desktopEntryInfo.AppName = appName
	}

	return desktopEntryInfo
}
