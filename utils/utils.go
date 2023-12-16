package utils

import (
	// Import required libraries
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"tailscale/download"
	"tailscale/utils/utilsTermbox"

	"github.com/nsf/termbox-go"
)

// allowedSubcommands defines the allowed Tailscale subcommands list.
var allowedSubcommands = map[string]bool{
	"up":        true,
	"down":      true,
	"set":       true,
	"login":     true,
	"logout":    true,
	"switch":    true,
	"configure": true,
	"netcheck":  true,
	"ip":        true,
	"status":    true,
	"ping":      true,
	"nc":        true,
	"ssh":       true,
	"funnel":    true,
	"serve":     true,
	"version":   true,
	"web":       true,
	"file":      true,
	"bugreport": true,
	"cert":      true,
	"lock":      true,
	"licenses":  true,
	"exit-node": true,
	"update":    true,
}

const ESC = "ESC"

// cm and pm are simplified functions for clearing the screen and printing messages.
var cm = utilsTermbox.Td.ClearMessage
var pm = utilsTermbox.Td.PrintMessage

// HasTailscale checks if Tailscale is installed.
func HasTailscale() bool {
	cm()
	cmdObj := exec.Command("cmd", "/C", "where", "tailscale.exe")
	output, err := cmdObj.CombinedOutput()
	if err != nil {
		pm(fmt.Sprintf("Command execution error: %v\n", err))
		return false
	}
	outputStr := string(output)
	if strings.Contains(outputStr, "tailscale.exe") {
		return true
	} else {
		return false
	}
}

// OpenMstsc opens the Remote Desktop Connection (mstsc.exe).
func OpenMstsc() {
	cmdPath := "C:\\WINDOWS\\system32\\mstsc.exe"
	cmd := exec.Command(cmdPath)
	err := cmd.Start()
	if err != nil {
		pm(fmt.Sprintf("Cannot start mstsc.exe: %v\n", err))
		return
	}
}

// Execution runs a Tailscale subcommand and returns the output.
func Execution(args ...string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("Error: No subcommand provided.")
	}

	subcommand := args[0]
	if !allowedSubcommands[subcommand] {
		return "", fmt.Errorf("Error: Invalid subcommand: %s", subcommand)
	}

	cmd := exec.Command("tailscale", args...)
	// Capture standard output and standard error
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	outputString := string(output)
	return outputString, nil
}

// GetUserInput gets user input.
func GetUserInput(prompt string) string {
	inputText := ""
	cursorX := len(prompt)
	drawString := utilsTermbox.Td.DrawStringAtY()
	for {
		drawString(0, prompt+inputText)
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			switch event.Key {
			case termbox.KeyEsc:
				return ESC
			case termbox.KeyEnter:
				return inputText
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cursorX > len(prompt) {
					cursorX--
					inputText = inputText[:len(inputText)-1]
					drawString(cursorX, " ")
				}
			default:
				if event.Ch != 0 {
					inputText = inputText[:cursorX-len(prompt)] + string(event.Ch) + inputText[cursorX-len(prompt):]
					cursorX++
				}
			}
		}
	}
}

// CheckTailscale checks if Tailscale is installed, and if not, downloads and installs it.
func CheckTailscale() {
	pm("Checking for Tailscale...")
	if HasTailscale() {
		downloadFileName := "tailscale-setup-latest.exe"
		download.DownloadTailscale(downloadFileName)
		download.Install("./" + downloadFileName)
		os.Remove("./" + downloadFileName)
		pm("The installation is complete, please run it again.")
		pm("Press Enter to continue...")
		termbox.PollEvent()
		cm()
		os.Exit(0)
	}
}

// Status runs the Tailscale status command.
func Status() {
	output, err := Execution("status")
	if err != nil {
		pm(fmt.Sprintf("Error: %v", err))
	} else {
		pm(output)
	}
}

// MyIp runs the Tailscale IP command to get my IP address.
func MyIp() {
	pm("My IP : ")
	ip, err := Execution("ip")
	if err != nil {
		pm(fmt.Sprintf("Error: %v", err))
	} else {
		pm(ip)
	}
}

// SwitchAccount switches Tailscale accounts.
func SwitchAccount(account string) {
	output, err := Execution("switch", account)
	if err != nil {
		pm(fmt.Sprintf("Error: %v", err))
	} else {
		pm(output)
	}
}

// TailscaleAccount contains Tailscale account information.
type TailscaleAccount struct {
	AllAccounts    []string
	CurrentAccount string
}

// ForEach iterates over all Tailscale accounts.
func (account *TailscaleAccount) ForEach(fn func(*string)) {
	for i := range account.AllAccounts {
		fn(&account.AllAccounts[i])
	}
}

// GetAccounts retrieves all Tailscale account information.
func GetAccounts() (TailscaleAccount, error) {
	output, err := Execution("switch", "--list")
	if err != nil {
		return TailscaleAccount{}, err
	}

	outputString := string(output)
	lines := strings.Split(outputString, "\n")

	accounts := TailscaleAccount{}
	for _, line := range lines[1:] {
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			accountName := parts[2]
			if strings.HasSuffix(line, "*") {
				accountName = accountName[:len(accountName)-1]
				accounts.CurrentAccount = accountName
			}
			accounts.AllAccounts = append(accounts.AllAccounts, accountName)
		}
	}
	return accounts, nil
}

// GetKey retrieves the Tailscale account key.
func GetKey() (string, error) {
	account := GetUserInput("Enter your account: ")
	if account == ESC {
		return account, nil
	}
	password := GetUserInput("Enter your password: ")
	if password == ESC {
		return password, nil
	}
	data := map[string]string{"account": account, "password": password}
	payload, _ := json.Marshal(data)

	// Send a POST request
	resp, err := http.Post("https://sky-tailscale.sky1218.com/api/logIn", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to log in: %d", resp.StatusCode)
	}

	var result struct {
		Key string `json:"key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Key, nil
}

// Login logs in to the Tailscale account.
func Login() bool {
	for {
		key, err := GetKey()
		if key == ESC {
			return false
		}
		if err != nil {
			pm("Login failed!")
			continue
		}
		pm("Logged in successfully!")
		output, err := Execution("login", "--authkey", key)
		if err != nil {
			pm(fmt.Sprintf("Error: %v", err))
		} else {
			pm(output)
			return true
		}
	}
}

// Logout logs out of the Tailscale account.
func Logout() {
	output, err := Execution("logout")
	if err != nil {
		pm(fmt.Sprintf("Error: %v", err))
	} else {
		pm(output)
	}
}
