package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"pkgfetch/Globals"
	"pkgfetch/Logger"
	"strings"
)

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

func ExpandHome(path string) string {
	if !strings.HasPrefix(path, "~/") {
		Logger.LogWarning(string("Given String To ExpandHome Func Doesn't Start With '~/' Which Is Recognised As Home Dir." +
			"\nIf You Changed The Config File With Your Absolute Home Path, Ignore This Warning." +
			"\nIf Not, Please Re-Check The Config File (Default Path: ~/.config/pkgfetch/pkgf.toml) To Make Sure The Paths Are Correct"))
		return path
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		Logger.LogError("Cannot Find User Home Dir \nBecause: %v", err)
		Logger.LogNewLine()
		return ""
	}

	return strings.Replace(
		path,
		"~",
		userHome,
		1,
	)
}

func FindDesktopFile(root string) string {
	desktopFile := ""

	err := filepath.WalkDir(root,
		func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}

			if strings.HasSuffix(entry.Name(), ".desktop") {
				desktopFile = path
				return filepath.SkipAll
			}

			return nil
		},
	)
	if err != nil {
		Logger.LogWarning("Failed To Fetch The '.desktop' File Of This AppImage. \nBecause: %v. \nNo Desktop Entry Will Be Created.", err)
		Logger.LogNewLine()
		return ""
	}

	return desktopFile
}

func FetchIconFromDesktopFile(filePath string, fileRoot string) string {
	iconName := FindVariableInFile(filePath, "Icon=")
	if iconName == "" {
		return ""
	}

	// Search For Icon
	var iconFilePath string
	err := filepath.WalkDir(fileRoot,
		func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}

			base := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			if strings.EqualFold(base, iconName) {
				for _, suffix := range Globals.SupportedIconSuffix {
					if entry.Name() == string(iconName+suffix) {
						iconFilePath = path
						return filepath.SkipAll
					}
				}
			}

			return nil
		},
	)

	if err != nil {
		Logger.LogWarning("No Icon File Was Found. \nBecause: %v \nDesktop Entry Will Be Created But Without An Icon \nDesktop Entry Path: %s", err, filePath)
		Logger.LogNewLine()
		return ""
	}

	return iconFilePath
}

func FindVariableInFile(filePath string, varName string) string {
	file, err := os.Open(filePath)
	if err != nil {
		Logger.LogError("Couldn't Open File At Path: %s \nBecause: %v", filePath, err)
		Logger.LogNewLine()
		return ""
	}
	defer file.Close()

	variable := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, varName) {
			variable = strings.TrimPrefix(line, varName)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		Logger.LogError("Couldn't Find Anything That Starts With '%s' In File %s. \nBecause: %v", varName, filePath, err)
		Logger.LogNewLine()
		return ""
	}

	return variable
}

func SetVariableInFile(filePath string, varName string, value string) error {
	prefix := varName + "="
	varNotFound := true

	content, err := os.ReadFile(filePath)
	if err != nil {
		Logger.LogError("Failed To Read File: %s \nBecause: %v", filePath, err)
		Logger.LogNewLine()
		return err
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, prefix) {
			lines[i] = prefix + value
			varNotFound = false
			break
		}
	}
	newContent := strings.Join(lines, "\n")
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		Logger.LogError("Failed To Write File Into: %s \bBecause: %v", filePath, err)
		Logger.LogNewLine()
		return err
	}

	if varNotFound {
		return fmt.Errorf("Variable: %s Was Not Found In File: %s", varName, filePath)
	}

	return nil
}
