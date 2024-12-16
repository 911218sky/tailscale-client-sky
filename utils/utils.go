package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"tailscale/download"
	"tailscale/utils/drawer"
	"time"

	"github.com/nsf/termbox-go"
)

// HasTailscale checks if Tailscale is installed by executing the --version command.
// It returns true if Tailscale is installed and false otherwise.
func HasTailscale() bool {
	drawer.Clear(drawer.DefaultOption)
	cmd := exec.Command("tailscale", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		drawer.Print(fmt.Sprintf("Command execution error: %v", err), drawer.DefaultOptionNoFlush)
		return false
	}
	outputStr := string(output)
	if strings.Contains(outputStr, "go version") {
		drawer.Print(outputStr, drawer.DefaultOptionNoFlush)
		return true
	}
	return false
}

// OpenMstsc launches the Windows Remote Desktop Connection (mstsc.exe).
// This function only works on Windows systems.
func OpenMstsc() {
	if runtime.GOOS != "windows" {
		drawer.Print("System is not Windows, cannot start mstsc.exe", drawer.DefaultOptionNoFlush)
		return
	}

	cmdPath := "C:\\WINDOWS\\system32\\mstsc.exe"
	if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
		drawer.Print("mstsc.exe not found, please install Remote Desktop Connection", drawer.DefaultOptionNoFlush)
		return
	}

	if err := exec.Command(cmdPath).Start(); err != nil {
		drawer.Print(fmt.Sprintf("Error starting mstsc: %v", err), drawer.DefaultOptionNoFlush)
	}
}

// Execution runs a Tailscale subcommand with the provided arguments.
// It validates the subcommand against allowed commands and returns the command output.
func Execution(args ...string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("no subcommand provided")
	}

	subcommand := args[0]
	if !AllowedSubcommands[subcommand] {
		return "", fmt.Errorf("invalid subcommand: %s", subcommand)
	}

	cmd := exec.Command("tailscale", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	return string(output), nil
}

// GetUserInput displays a prompt and reads user input from the terminal.
// It handles special keys like Escape, Enter, and Backspace.
func GetUserInput(prompt string) string {
	var inputText strings.Builder
	cursorX := len(prompt)
	y := drawer.GetY()

	for {
		drawer.Render(y, 0, prompt+inputText.String())
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			switch event.Key {
			case termbox.KeyEsc:
				return KeyEsc
			case termbox.KeyEnter:
				drawer.NextLine()
				return inputText.String()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if cursorX > len(prompt) && inputText.Len() > 0 {
					cursorX--
					str := inputText.String()
					inputText.Reset()
					inputText.WriteString(str[:len(str)-1])
					drawer.Render(y, cursorX, " ")
				}
			default:
				if event.Ch != 0 {
					inputText.WriteRune(event.Ch)
					cursorX++
				}
			}
		}
	}
}

// CheckTailscale verifies if Tailscale is installed and installs it if not found.
func CheckTailscale() {
	drawer.Print("Checking for Tailscale...", drawer.DefaultOption)
	config := NewWaitAndExitConfig()
	if !HasTailscale() {
		installTailscale()
		waitAndExit("Installation is complete. Please run it again.", config)
	}
	waitAndExit("Environmental inspection complete.", config.WithShouldExit(false))
}

// installTailscale handles the installation of Tailscale based on the operating system.
func installTailscale() {
	if runtime.GOOS == "windows" {
		// Use a temporary directory instead of current directory
		tmpDir, err := os.MkdirTemp("", "tailscale-installer")
		if err != nil {
			drawer.Print(fmt.Sprintf("Failed to create temp dir: %v", err), drawer.DefaultOption)
			cleanupAndExit(tmpDir, 1)
		}

		exe := filepath.Join(tmpDir, "tailscale-setup-latest.exe")
		if _, notExist := os.Stat(exe); os.IsNotExist(notExist) {
			if err := download.DownloadTailscaleWindows(exe); err != nil {
				drawer.Print(fmt.Sprintf("Download error: %v", err), drawer.DefaultOption)
				cleanupAndExit(tmpDir, 1)
			}
		}
		if err := download.Install(exe); err != nil {
			drawer.Print(fmt.Sprintf("Installation error: %v", err), drawer.DefaultOption)
			cleanupAndExit(tmpDir, 1)
		}
		cleanupAndExit(tmpDir, 0)
	} else if runtime.GOOS == "linux" {
		if err := download.DownloadTailscaleLinux(); err != nil {
			drawer.Print(fmt.Sprintf("Installation error: %v", err), drawer.DefaultOption)
			os.Exit(1)
		}
	}
}

// cleanupAndExit removes the temporary directory and exits the program with the given code.
func cleanupAndExit(tmpDir string, code int) {
	drawer.Print(fmt.Sprintf("Removing temp dir: %s", tmpDir), drawer.DefaultOption)
	os.RemoveAll(tmpDir)
	os.Exit(code)
}

// NewWaitAndExitConfig creates a WaitAndExitConfig with default values
func NewWaitAndExitConfig() *WaitAndExitConfig {
	return &WaitAndExitConfig{
		ShouldExit: true,
		Countdown:  3,
	}
}

// WithShouldExit sets the ShouldExit option
func (c *WaitAndExitConfig) WithShouldExit(shouldExit bool) *WaitAndExitConfig {
	c.ShouldExit = shouldExit
	return c
}

// WithCountdown sets the Countdown duration
func (c *WaitAndExitConfig) WithCountdown(countdown int) *WaitAndExitConfig {
	c.Countdown = countdown
	return c
}

// waitAndExit displays a message and waits for user input or timeout before exiting.
// If config is nil, default configuration will be used.
func waitAndExit(msg string, config *WaitAndExitConfig) {
	drawer.Print(msg, drawer.DefaultOptionNoFlush)

	if config == nil {
		config = NewWaitAndExitConfig()
	}

	done := make(chan struct{})
	go func() {
		termbox.PollEvent()
		close(done)
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	x := drawer.GetX()
	y := drawer.GetY()
	countdown := config.Countdown

	for {
		drawer.Render(y+1, x, fmt.Sprintf("\rPress Enter to continue (auto exit in %ds)...", countdown))
		select {
		case <-done:
			drawer.NextLine()
			drawer.Clear(drawer.DefaultOption)
			if config.ShouldExit {
				os.Exit(0)
			}
			return
		case <-ticker.C:
			countdown--
			if countdown <= 0 {
				drawer.NextLine()
				drawer.Clear(drawer.DefaultOption)
				if config.ShouldExit {
					os.Exit(0)
				}
				return
			}
		}
	}
}

// Status executes the Tailscale status command and displays the result.
func Status() {
	output, err := Execution("status")
	if err != nil {
		drawer.Print(fmt.Sprintf("Error getting status: %v", err), drawer.DefaultOption)
		return
	}
	drawer.Print(output, drawer.DefaultOption)
}

// MyIP retrieves and displays the current Tailscale IP address.
func MyIP() {
	drawer.Print("My IP: ", drawer.DefaultOptionNoFlush)
	ip, err := Execution("ip")
	if err != nil {
		drawer.Print(fmt.Sprintf("Error getting IP: %v", err), drawer.DefaultOption)
		return
	}
	drawer.Print(ip, drawer.DefaultOption)
}

// SwitchAccount changes the active Tailscale account to the specified account.
func SwitchAccount(account string) {
	output, err := Execution("switch", account)
	if err != nil {
		drawer.Print(fmt.Sprintf("Error switching account: %v", err), drawer.DefaultOption)
		return
	}
	drawer.Print(output, drawer.DefaultOption)
}

// TailscaleAccount represents the structure for storing Tailscale account information.
type TailscaleAccount struct {
	AllAccounts    []string // List of all available Tailscale accounts
	CurrentAccount string   // Currently active Tailscale account
}

// ForEach executes the provided function for each account in the TailscaleAccount.
func (account *TailscaleAccount) ForEach(fn func(*string)) {
	for i := range account.AllAccounts {
		fn(&account.AllAccounts[i])
	}
}

// GetAccounts retrieves all available Tailscale accounts and the current active account.
func GetAccounts() (*TailscaleAccount, error) {
	output, err := Execution("switch", "--list")
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	accounts := &TailscaleAccount{}

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

// GetKey prompts for user credentials and retrieves a Tailscale authentication key.
func GetKey() (string, error) {
	account := GetUserInput("Enter your account: ")
	if account == KeyEsc {
		return account, nil
	}
	password := GetUserInput("Enter your password: ")
	if password == KeyEsc {
		return password, nil
	}

	data := map[string]string{
		"account":  account,
		"password": password,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal credentials: %w", err)
	}

	resp, err := http.Post(LoginAPIEndpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Key string `json:"key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Key, nil
}

// Login handles the Tailscale login process using an authentication key.
func Login() bool {
	for {
		key, err := GetKey()
		if key == KeyEsc {
			return false
		}
		if err != nil {
			drawer.Print(fmt.Sprintf("Login failed: %v", err), drawer.DefaultOption)
			continue
		}

		drawer.Print("Logging in...", drawer.DefaultOptionNoFlush)
		output, err := Execution("login", "--authkey", key)
		if err != nil {
			drawer.Print(fmt.Sprintf("Login error: %v", err), drawer.DefaultOption)
			continue
		}

		drawer.Print("Logged in successfully!", drawer.DefaultOptionNoFlush)
		drawer.Print(output, drawer.DefaultOption)
		return true
	}
}

// Logout performs the Tailscale logout operation.
func Logout() {
	output, err := Execution("logout")
	if err != nil {
		drawer.Print(fmt.Sprintf("Logout error: %v", err), drawer.DefaultOption)
		return
	}
	drawer.Print(output, drawer.DefaultOption)
}
