package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"tailscale/download"
	"tailscale/utils/drawer"

	"github.com/nsf/termbox-go"
)

// HasTailscale checks if Tailscale is installed.
func HasTailscale() bool {
	drawer.Clear(drawer.DefaultOption)
	cmd := exec.Command("tailscale", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		drawer.Print(fmt.Sprintf("Command execution error: %v\n", err), drawer.DefaultOptionNoFlush)
		return false
	}
	outputStr := string(output)
	if strings.Contains(outputStr, "go version") {
		drawer.Print(outputStr, drawer.DefaultOptionNoFlush)
		return true
	}
	return false
}

// OpenMstsc opens the Remote Desktop Connection (mstsc.exe).
func OpenMstsc() {
	if runtime.GOOS != "windows" {
		drawer.Print("System is not Windows, cannot start mstsc.exe\n", drawer.DefaultOptionNoFlush)
		return
	}
	cmdPath := "C:\\WINDOWS\\system32\\mstsc.exe"
	err := exec.Command(cmdPath).Start()
	if err != nil {
		drawer.Print(fmt.Sprintf("Error: %v\n", err), drawer.DefaultOptionNoFlush)
	}
}

// Execution runs a Tailscale subcommand and returns the output.
func Execution(args ...string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("error: no subcommand provided")
	}

	subcommand := args[0]
	if !allowedSubcommands[subcommand] {
		return "", fmt.Errorf("error: invalid subcommand: %s", subcommand)
	}

	cmd := exec.Command("tailscale", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// GetUserInput prompts the user for input.
func GetUserInput(prompt string) string {
	var inputText string
	cursorX := len(prompt)
	y := drawer.GetY()

	for {
		drawer.Render(y, 0, prompt+inputText)
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			switch event.Key {
			case termbox.KeyEsc:
				return ESC
			case termbox.KeyEnter:
				drawer.NextLine()
				return inputText
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cursorX > len(prompt) {
					cursorX--
					inputText = inputText[:len(inputText)-1]
					drawer.Render(y, cursorX, " ")
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

// CheckTailscale checks for Tailscale installation and installs it if not found.
func CheckTailscale() {
	drawer.Print("Checking for Tailscale...", drawer.DefaultOption)
	if !HasTailscale() {
		if runtime.GOOS == "windows" {
			downloadFileName := "tailscale-setup-latest.exe"
			download.DownloadTailscaleWindows(downloadFileName)
			err := download.Install("./" + downloadFileName)
			if err != nil {
				drawer.Print(fmt.Sprintf("Error: %v", err), drawer.DefaultOption)
				os.Exit(0)
			}
			os.Remove("./" + downloadFileName)
		} else if runtime.GOOS == "linux" {
			download.DownloadTailscaleLinux()
		}
		drawer.Print("Installation is complete. Please run it again.", drawer.DefaultOptionNoFlush)
		drawer.Print("Press Enter to continue...", drawer.DefaultOption)
		termbox.PollEvent()
		drawer.Clear(drawer.DefaultOption)
		os.Exit(0)
	}
	drawer.Print("Environmental inspection complete.", drawer.DefaultOptionNoFlush)
	drawer.Print("Press Enter to continue...", drawer.DefaultOption)
	termbox.PollEvent()
	drawer.Clear(drawer.DefaultOption)
}

// Status runs the Tailscale status command.
func Status() {
	output, err := Execution("status")
	if err != nil {
		drawer.Print(fmt.Sprintf("Error: %v", err), drawer.DefaultOption)
		return
	}
	drawer.Print(output, drawer.DefaultOption)
}

// MyIp runs the Tailscale IP command to get the IP address.
func MyIp() {
	drawer.Print("My IP : ", drawer.DefaultOptionNoFlush)
	ip, err := Execution("ip")
	if err != nil {
		drawer.Print(fmt.Sprintf("Error: %v", err), drawer.DefaultOption)
		return
	}
	drawer.Print(ip, drawer.DefaultOption)
}

// SwitchAccount switches Tailscale accounts.
func SwitchAccount(account string) {
	output, err := Execution("switch", account)
	if err != nil {
		drawer.Print(fmt.Sprintf("Error: %v", err), drawer.DefaultOption)
		return
	}
	drawer.Print(output, drawer.DefaultOption)
}

// TailscaleAccount holds Tailscale account information.
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

	lines := strings.Split(string(output), "\n")
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

// GetKey retrieves the Tailscale account key from user input.
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

	resp, err := http.Post(GET_KEY_URL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to log in: %d", resp.StatusCode)
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
			drawer.Print("Login failed!", drawer.DefaultOption)
			continue
		}
		drawer.Print("Please wait, logging in...", drawer.DefaultOptionNoFlush)
		output, err := Execution("login", "--authkey", key)
		if err != nil {
			drawer.Print(fmt.Sprintf("Error: %v", err), drawer.DefaultOption)
		} else {
			drawer.Print("Logged in successfully!", drawer.DefaultOptionNoFlush)
			drawer.Print(output, drawer.DefaultOption)
			return true
		}
	}
}

// Logout logs out of the Tailscale account.
func Logout() {
	output, err := Execution("logout")
	if err != nil {
		drawer.Print(fmt.Sprintf("Error: %v", err), drawer.DefaultOption)
		return
	}
	drawer.Print(output, drawer.DefaultOption)
}
