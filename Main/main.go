package main

import (
	"fmt"
	"os"
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
	}
}

func HandleInstall() {
	if argument, ok := IsArgumentGiven(); ok {
		fmt.Printf("Installing: %s\n", argument)
	}
}

// Helpers
func IsArgumentGiven() (string, bool) {
	if len(os.Args) < 3 {
		fmt.Printf("No Argument Specified. \n%s", Logger.LOG_MISSUSE_COMMAND)
		return "", false
	}
	return os.Args[2], true
}

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
