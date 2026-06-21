package Logger

import (
	"fmt"

	"pkgfetch/Globals"
)

const LOG_MISSUSE_COMMAND = "Type `" + Globals.PROGRAM_NAME_CMD + " help` to show the list of available commands and arguments"

func ShowHelpDialog() {
	LogMessageSameLine("Usage:")
	LogMessage("  %s <command>", Globals.PROGRAM_NAME_CMD)
	LogMessage("  %s <command> <args>", Globals.PROGRAM_NAME_CMD)
	LogNewLine()
	LogMessage("Available Commands:")
	LogMessage("  - help")
	LogMessage("      Show this dialog")
	LogNewLine()
	LogMessage("  - search <package_name> | <package_owner/package_name>")
	LogMessage("      Searches for the given package name")
	LogNewLine()
	LogMessage("  - install <package_name> | <package_owner/package_name>")
	LogMessage("      Installs given package name (if exists)")
}

// Message Loggins
func LogMessage(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("\n%s", msg)
}

func LogWarning(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("\nWARNING: %s", msg)
}

func LogError(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("\nERROR: %s", msg)
}

func LogMessageSameLine(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s", msg)
}

func LogWarningSameLine(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("WARNING: %s", msg)
}

func LogErrorSameLine(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("ERROR: %s", msg)
}

func LogNewLine() {
	fmt.Println()
}
