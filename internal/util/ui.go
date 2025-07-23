package util

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	Primary   = color.New(color.FgCyan, color.Bold)
	Success   = color.New(color.FgGreen, color.Bold)
	Warning   = color.New(color.FgYellow, color.Bold)
	Error     = color.New(color.FgRed, color.Bold)
	Info      = color.New(color.FgBlue)
	Muted     = color.New(color.FgHiBlack)
	
	Running   = color.New(color.FgGreen)
	Stopped   = color.New(color.FgRed)
	Pending   = color.New(color.FgYellow)
)

func PrintHeader(title string, emoji string) {
	fmt.Println()
	Primary.Printf("%s %s\n", emoji, title)
	Primary.Println(strings.Repeat("─", len(title)+4))
	fmt.Println()
}

func PrintSubHeader(title string) {
	fmt.Println()
	Info.Printf("▶ %s\n", title)
	fmt.Println()
}

func PrintSuccess(message string) {
	Success.Printf("[SUCCESS] %s\n", message)
}

func PrintError(message string) {
	Error.Printf("[ERROR] %s\n", message)
}

func PrintWarning(message string) {
	Warning.Printf("[WARNING] %s\n", message)
}

func PrintInfo(message string) {
	Info.Printf("[INFO] %s\n", message)
}

func FormatStatus(status string) string {
	switch strings.ToLower(status) {
	case "running":
		return Running.Sprintf("🟢 %s", status)
	case "stopped", "terminated":
		return Stopped.Sprintf("🔴 %s", status)
	case "pending", "starting":
		return Pending.Sprintf("🟡 %s", status)
	default:
		return Muted.Sprintf("⚪ %s", status)
	}
}

func FormatServiceType(serviceType string) string {
	return serviceType // Just return as-is, no emojis
}

func FormatConnectionType(connType string) string {
	return connType // Just return as-is, no emojis
}

func PrintWelcome() {
	fmt.Println()
	Primary.Println("Welcome to Gossy!")
	Muted.Println("Your AWS connection companion")
	fmt.Println()
}
