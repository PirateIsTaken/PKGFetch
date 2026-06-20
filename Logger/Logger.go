package Logger

import (
	"fmt"

	"pkgfetch/Globals"
)

const LOG_MISSUSE_COMMAND = "Type `" + Globals.PROGRAM_NAME_CMD + " help` to show the list of available commands and arguments"

func ShowHelpDialog() {
	fmt.Println("Usage:")
	fmt.Println("  pkgf <command>")
	fmt.Println("  pkgf <command> <args>")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  help")
	fmt.Println("    Show this dialog")
	println()
	fmt.Println("  search <package_name> | <package_owner/package_name>")
	fmt.Println("    Searches for the given package name")
	println()
	fmt.Println("  install <package_name> | <package_owner/package_name>")
	fmt.Println("    Installs given package name (if exists)")
}
