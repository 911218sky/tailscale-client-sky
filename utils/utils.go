package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"tailscale/download"
	"tailscale/utilsTermbox"

	"github.com/nsf/termbox-go"
)

var allowedSubcommands = map[string]bool{
	"up":        true,
	"down":      true,
	"set":       true,
	"login":     true,
	"logout":    true,
	"switch":    true,
	"configure": true, // ALPHA
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
	"update":    true, // ALPHA
}

func HasTailscale() bool {
	utilsTermbox.ClearMessage()
	cmdObj := exec.Command("cmd", "/C", "where", "tailscale.exe")
	output, err := cmdObj.CombinedOutput()
	if err != nil {
		utilsTermbox.PrintMessage(fmt.Sprintf("Command execution error: %v\n", err))
		return false
	}
	outputStr := string(output)
	if strings.Contains(outputStr, "tailscale.exe") {
		return true
	} else {
		return false
	}
}

func OpenMstsc() {
	cmdPath := "C:\\WINDOWS\\system32\\mstsc.exe"
	cmd := exec.Command(cmdPath)
	err := cmd.Start()
	if err != nil {
		utilsTermbox.PrintMessage(fmt.Sprintf("Cannot start mstsc.exe: %v\n", err))
		return
	}
}

func Execution(args ...string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("Error: No subcommand provided.")
	}

	subcommand := args[0]
	if !allowedSubcommands[subcommand] {
		return "", fmt.Errorf("Error: Invalid subcommand: %s", subcommand)
	}

	cmd := exec.Command("tailscale", args...)
	// 捕获标准输出和标准错误
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	outputString := string(output)
	return outputString, nil
}

func GetUserInput(prompt string) string {
	inputText := ""
	cursorX := len(prompt)

	for {
		utilsTermbox.DrawString(0, utilsTermbox.YTermbox, prompt+inputText)
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			if event.Key == termbox.KeyEnter {
				utilsTermbox.YTermbox++
				return inputText
			} else if event.Key == termbox.KeyEsc {
				inputText = ""
				return inputText
			} else if event.Key == termbox.KeyBackspace || event.Key == termbox.KeyBackspace2 {
				if cursorX > len(prompt) {
					cursorX--
					inputText = inputText[:len(inputText)-1]
					utilsTermbox.DrawString(cursorX, utilsTermbox.YTermbox, " ")
				}
			} else if event.Ch != 0 {
				inputText = inputText[:cursorX-len(prompt)] + string(event.Ch) + inputText[cursorX-len(prompt):]
				cursorX++
			}
		}
	}
}

func CheckTailscale() {
	utilsTermbox.PrintMessage("Checking for Tailscale...")
	if !HasTailscale() {
		downloadFileName := "tailscale-setup-latest.exe"
		download.DownloadTailscale(downloadFileName)
		download.Install("./" + downloadFileName)
		os.Remove("./" + downloadFileName)
		utilsTermbox.PrintMessage("The installation is complete, please run it again.")
		utilsTermbox.PrintMessage("Press Enter to continue...")
		termbox.PollEvent()
		utilsTermbox.ClearMessage()
		os.Exit(0)
	}
}

func Status() {
	output, err := Execution("status")
	if err != nil {
		utilsTermbox.PrintMessage(fmt.Sprintf("Error: %v", err))
	} else {
		utilsTermbox.PrintMessage(output)
	}
}

func MyIp() {
	utilsTermbox.PrintMessage("My IP : ")
	ip, err := Execution("ip")
	if err != nil {
		utilsTermbox.PrintMessage(fmt.Sprintf("Error: %v", err))
	} else {
		utilsTermbox.PrintMessage(ip)
	}
}

func SwitchAccount(account string) {
	output, err := Execution("switch", account)
	if err != nil {
		utilsTermbox.PrintMessage(fmt.Sprintf("Error: %v", err))
	} else {
		utilsTermbox.PrintMessage(output)
	}
}

type TailscaleAccount struct {
	AllAccounts    []string
	CurrentAccount string
}

func (account *TailscaleAccount) ForEach(fn func(*string)) {
	for i := range account.AllAccounts {
		fn(&account.AllAccounts[i])
	}
}

func GetAccounts() (TailscaleAccount, error) {
	output, err := Execution("switch", "--list")
	if err != nil {
		return TailscaleAccount{}, err
	}

	outputString := string(output)
	lines := strings.Split(outputString, "\n")

	accounts := TailscaleAccount{}
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			accountName := parts[0]
			accounts.AllAccounts = append(accounts.AllAccounts, accountName)
			if strings.HasSuffix(line, "*") {
				accounts.CurrentAccount = accountName
			}
		}
	}
	return accounts, nil
}

func GetKey() (string, error) {
	account := GetUserInput("Enter your account: ")
	password := GetUserInput("Enter your password: ")

	data := map[string]string{"account": account, "password": password}
	payload, _ := json.Marshal(data)

	// 发起 POST 请求
	resp, err := http.Post("https://sky-tailscale.sky1218.com/api/logIn", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 解析响应
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

func Login() {
	for {
		key, err := GetKey()
		if err != nil {
			utilsTermbox.PrintMessage("Login failed!")
			continue
		}
		utilsTermbox.PrintMessage("Landed successfully!")
		output, err := Execution("login", "--authkey", key)
		if err != nil {
			utilsTermbox.PrintMessage(fmt.Sprintf("Error: %v", err))
		} else {
			utilsTermbox.PrintMessage(output)
			break
		}
	}
}

func Logout() {
	output, err := Execution("logout")
	if err != nil {
		utilsTermbox.PrintMessage(fmt.Sprintf("Error: %v", err))
	} else {
		utilsTermbox.PrintMessage(output)
	}
}
