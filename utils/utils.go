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
	cmdObj := exec.Command("cmd", "/C", "where", "tailscale.exe")
	output, err := cmdObj.CombinedOutput()
	if err != nil {
		PrintMessage(fmt.Sprintf("Command execution error: %v\n", err))
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
		PrintMessage(fmt.Sprintf("Cannot start mstsc.exe: %v\n", err))
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
		DrawString(0, y, prompt+inputText)
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			if event.Key == termbox.KeyEnter {
				y++
				return inputText
			} else if event.Key == termbox.KeyEsc {
				inputText = ""
				return inputText
			} else if event.Key == termbox.KeyBackspace || event.Key == termbox.KeyBackspace2 {
				if cursorX > len(prompt) {
					cursorX--
					inputText = inputText[:len(inputText)-1]
					DrawString(cursorX, y, " ")
				}
			} else if event.Ch != 0 {
				inputText = inputText[:cursorX-len(prompt)] + string(event.Ch) + inputText[cursorX-len(prompt):]
				cursorX++
			}
		}
	}
}

func CheckTailscale() {
	PrintMessage("Checking for Tailscale...")
	if !HasTailscale() {
		downloadFileName := "tailscale-setup-latest.exe"
		download.DownloadTailscale(downloadFileName)
		download.Install("./" + downloadFileName)
		PrintMessage("The installation is complete, please run it again.")
		PrintMessage("Press Enter to continue...")
		termbox.PollEvent()
		os.Exit(0)
	}
}

func Status() {
	output, err := Execution("status")
	if err != nil {
		PrintMessage(fmt.Sprintf("Error: %v", err))
	} else {
		PrintMessage(output)
	}
}

func MyIp() {
	PrintMessage("My IP : ")
	ip, err := Execution("ip")
	if err != nil {
		PrintMessage(fmt.Sprintf("Error: %v", err))
	} else {
		PrintMessage(ip)
	}
}

func SwitchAccount(account string) {
	output, err := Execution("switch", account)
	if err != nil {
		PrintMessage(fmt.Sprintf("Error: %v", err))
	} else {
		PrintMessage(output)
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
			PrintMessage("Login failed!")
			continue
		}
		PrintMessage("Landed successfully!")
		output, err := Execution("login", "--authkey", key)
		if err != nil {
			PrintMessage(fmt.Sprintf("Error: %v", err))
		} else {
			PrintMessage(output)
			break
		}
	}
}

func Logout() {
	output, err := Execution("logout")
	if err != nil {
		PrintMessage(fmt.Sprintf("Error: %v", err))
	} else {
		PrintMessage(output)
	}
}

var (
	x = 2
	y = 0
)

type Option struct {
	NoNewLine bool
	NoFlush   bool
}

func DrawString(x, y int, str string) {
	for i, ch := range str {
		termbox.SetCell(x+i, y, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()
}

func PrintMessage(message string, options ...Option) {
	option := Option{}
	if len(options) > 0 {
		if options[0].NoFlush {
			option.NoFlush = true
		}
		if options[0].NoNewLine {
			option.NoNewLine = true
		}
	}

	lines := strings.Split(message, "\n")
	for _, line := range lines {
		for _, ch := range line {
			termbox.SetCell(x, y, ch, termbox.ColorDefault, termbox.ColorDefault)
			x++
		}
		y++
		x = 2
	}

	if option.NoNewLine {
		y--
	}

	if !option.NoFlush {
		termbox.Flush()
	}
}

func ClearMessage(options ...Option) {
	x = 2
	y = 0
	isFlush := true
	if len(options) > 0 {
		if options[0].NoFlush {
			isFlush = false
		}
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if isFlush {
		termbox.Flush()
	}
}
