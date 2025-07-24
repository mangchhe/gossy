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
	Command   = color.New(color.FgMagenta)
	
	Running   = color.New(color.FgGreen)
	Stopped   = color.New(color.FgRed)
	Pending   = color.New(color.FgYellow)
)

const HeaderHeight = 15
var HeaderDisplayed = false
var CurrentProfile = "None"
var CurrentServiceType = "None"
var CurrentCommand = ""
var ShowCommandHistory = false

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

func PrintCommand(command string) {
	Command.Printf("Command: %s\n", command)
}

func PrintLiveCommand(baseCommand string, parts ...string) {
	fullCommand := baseCommand
	for _, part := range parts {
		if part != "" {
			fullCommand += " " + part
		}
	}
	Command.Printf("Command: %s\n", fullCommand)
}

func PrintFixedCommand(baseCommand string, parts ...string) {
	// Command history display disabled
}

func InitializeScreen() {
	fmt.Print("\033[2J\033[H")
	PrintGossyBanner()
	HeaderDisplayed = true
}

func ClearContent() {
	if HeaderDisplayed {
		fmt.Printf("\033[%d;1H", HeaderHeight+1)
		fmt.Print("\033[0J")
	} else {
		InitializeScreen()
	}
}

func UpdateProfile(profile string) {
	CurrentProfile = profile
	if HeaderDisplayed {
		fmt.Print("\033[1;1H")
		PrintGossyBanner()
	}
}

func UpdateServiceType(serviceType string) {
	CurrentServiceType = serviceType
	if HeaderDisplayed {
		fmt.Print("\033[1;1H")
		PrintGossyBanner()
	}
}

func UpdateCommand(command string) {
	CurrentCommand = command
	if HeaderDisplayed {
		fmt.Print("\033[1;1H")
		PrintGossyBanner()
	}
}

func SetShowCommandHistory(show bool) {
	ShowCommandHistory = show
}

func PrintCommandLine() {
	if ShowCommandHistory {
		if CurrentCommand != "" {
			Command.Printf("Command: %s\n", CurrentCommand)
		} else {
			Muted.Println("Select an option to see command...")
		}
	}
}



func PrintGossyBanner() {
	banner := `
                 ██████╗  ██████╗ ███████╗███████╗██╗   ██╗
                ██╔════╝ ██╔═══██╗██╔════╝██╔════╝╚██╗ ██╔╝
                ██║  ███╗██║   ██║███████╗███████╗ ╚████╔╝ 
                ██║   ██║██║   ██║╚════██║╚════██║  ╚██╔╝  
                ╚██████╔╝╚██████╔╝███████║███████║   ██║   
                 ╚═════╝  ╚═════╝ ╚══════╝╚══════╝   ╚═╝   
`
	
	Primary.Print(banner)
	
	fmt.Println()
	Primary.Print(strings.Repeat("━", 80))
	fmt.Println()
	
	fmt.Printf(" Profile: %s │ Service: %s │ Mode: %s │ Version: %s \n",
		color.New(color.FgGreen).Sprint(CurrentProfile),
		color.New(color.FgYellow).Sprint(CurrentServiceType),
		color.New(color.FgCyan).Sprint("Interactive"),
		color.New(color.FgMagenta).Sprint("v0.1.0"))
	fmt.Println()
	
	Primary.Print(strings.Repeat("━", 80))
	fmt.Println()
	
	PrintCommandLine()
	fmt.Println()
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
	return serviceType
}

func FormatConnectionType(connType string) string {
	return connType
}

func PrintWelcome() {
	fmt.Println()
	Primary.Println("Welcome to Gossy!")
	Muted.Println("Your AWS connection companion")
	fmt.Println()
}
