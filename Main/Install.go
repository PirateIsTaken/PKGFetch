package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"pkgfetch/Globals"
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
	Logger.LogNewLine()
	Logger.LogMessage("Copying Your AppImage To Directory %s", Globals.AppConfig.AppImagePath)
	Logger.LogNewLine()

	userHome, err := os.UserHomeDir()
	if err != nil {
		Logger.LogError("Cannot Find User Home Dir")
		Logger.LogNewLine()
		return
	}

	appImagePath := strings.Replace(
		Globals.AppConfig.AppImagePath,
		"~",
		userHome,
		1,
	)

	appName := filepath.Base(file)
	newPlace := filepath.Join(appImagePath, appName)

	err = CopyFile(file, newPlace)

	if err == nil {
		Logger.LogNewLine()
		Logger.LogMessage("Package Copied To %s", appImagePath)
		Logger.LogMessage("To Delete Package, Just Remove It From This Folder: %s", appImagePath)
		Logger.LogNewLine()

		SetupExecForAppImage(appName)
		AskToDeleteCache(file)
	} else {
		Logger.LogError("Failed To Copy Downloaded File From: %s To: %s \nBecause: %v",
			file, newPlace, err)
		Logger.LogNewLine()
	}
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

func SetupExecForAppImage(app string) {
	Logger.LogMessage("Creating Exec For %s", app)
	file := string(Globals.AppConfig.AppImagePath + "/" + app)

	appNoSuffix := strings.ToLower(app)
	appNoSuffix = strings.TrimSuffix(appNoSuffix, ".appimage")

	homeDir, err := os.UserHomeDir()
	file = strings.Replace(
		file,
		"~",
		homeDir,
		1,
	)

	if err != nil {
		Logger.LogError("Cannot Find User Home Dir")
		Logger.LogNewLine()
		return
	}
	appDesktopDir := filepath.Join(
		homeDir,
		".local",
		"share",
		"applications",
		string(appNoSuffix+".desktop"),
	)

	content := fmt.Sprintf(
		`[Desktop Entry]
		Version=1.0
		Type=Application
		Name=%s
		Exec=%s
		Terminal=false`,
		appNoSuffix,
		file,
	)

	err = os.WriteFile(
		appDesktopDir,
		[]byte(content),
		0644,
	)

	os.Chmod(appDesktopDir, 0755)

	Logger.LogMessage("Created A .desktop File In: %s", appDesktopDir)
	Logger.LogMessage("If It's Not Being Shown, Run `update-desktop-database <your .desktop dir.>`. Or Log Out and Log Back In")
	Logger.LogNewLine()
}

func CopyFile(source string, dest string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}

	dst, err := os.Create(dest)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	defer dst.Close()
	defer src.Close()
	return err
}
