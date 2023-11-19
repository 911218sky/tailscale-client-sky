package menu

import (
	"os"
	"strings"

	"tailscale/utils"
	"tailscale/utilsTermbox"

	"github.com/nsf/termbox-go"
)

const (
	CONNECT          = iota // 0: 連線
	SWITCHACCOUNT           // 1: 切換帳戶
	SIGNOUT                 // 2: 登出
	LIST_INFORMATION        // 3: 列出資訊
	OPEN_MSTSC              // 4: 開啟遠端桌面
	QUIT                    // 5: 退出
)

func RunTermboxUI() {

	defer termbox.Close()

	options := []string{"Connect", "SwitchAccount", "SignOut", "ListInformation", "OpenMstsc", "Quit"}
	optionToAction := map[int]func(){
		CONNECT:          Connect,
		SWITCHACCOUNT:    SwitchAccount,
		SIGNOUT:          SignOut,
		LIST_INFORMATION: ListInformation,
		OPEN_MSTSC:       utils.OpenMstsc,
		QUIT:             Quit,
	}
	selectedIndex := 0

	for {
		utilsTermbox.ClearMessage(utilsTermbox.Option{NoFlush: true})
		RenderMenu(options, selectedIndex)

		event := termbox.PollEvent()
		isEnter := handleKeyEvent(event, &selectedIndex, options)

		if !isEnter {
			continue
		}

		action, found := optionToAction[selectedIndex]
		if !found {
			utilsTermbox.PrintMessage("Unknown option")
			break
		}

		utilsTermbox.ClearMessage()
		action()
		utilsTermbox.ClearMessage()
	}
}

func RenderMenu(options []string, selectedIndex int) {
	for i, option := range options {
		if selectedIndex == i {
			option = ">  " + option
		}
		utilsTermbox.PrintMessage(option, utilsTermbox.Option{NoFlush: true})
	}
	termbox.Flush()
}

func handleKeyEvent(event termbox.Event, selectedIndex *int, options []string) bool {
	if event.Type == termbox.EventKey {
		switch event.Key {
		case termbox.KeyArrowUp:
			if *selectedIndex > 0 {
				*selectedIndex--
			} else {
				*selectedIndex = len(options) - 1
			}
		case termbox.KeyArrowDown:
			if *selectedIndex < len(options)-1 {
				*selectedIndex++
			} else {
				*selectedIndex = 0
			}
		case termbox.KeyEnter:
			return true
		case termbox.KeyEsc:
			termbox.Close()
		}
	}
	return false
}

func getAccount() string {
	tailscaleAccount, _ := utils.GetAccounts()
	selectedIndex := 0
	tailscaleAccount.ForEach(func(account *string) {
		if *account == tailscaleAccount.CurrentAccount {
			*account = "*" + *account
		}
	})
	tailscaleAccount.AllAccounts = append(tailscaleAccount.AllAccounts, "QUIT")
	for {
		utilsTermbox.ClearMessage(utilsTermbox.Option{NoFlush: true})
		RenderMenu(tailscaleAccount.AllAccounts, selectedIndex)
		event := termbox.PollEvent()
		isEnter := handleKeyEvent(event, &selectedIndex, tailscaleAccount.AllAccounts)
		if isEnter {
			break
		}
	}

	if tailscaleAccount.AllAccounts[selectedIndex] == "QUIT" {
		return ""
	}

	if strings.HasPrefix(tailscaleAccount.AllAccounts[selectedIndex], "*") {
		utilsTermbox.PrintMessage("It is not possible to select an account that is currently in use!")
		utilsTermbox.PrintMessage("Press Enter to continue...")
		termbox.PollEvent()
		return ""
	}

	return tailscaleAccount.AllAccounts[selectedIndex]
}

func Connect() {
	utils.Login()
	utils.Status()
	utils.OpenMstsc()
	utilsTermbox.PrintMessage("Press Enter to continue...")
	termbox.PollEvent()
}

func SwitchAccount() {
	account := getAccount()
	if account == "" {
		return
	}
	utils.SwitchAccount(account)
	utils.Status()
	utils.OpenMstsc()
	utilsTermbox.PrintMessage("Press Enter to continue...")
	termbox.PollEvent()
}

func SignOut() {
	utils.Logout()
	utilsTermbox.PrintMessage("Press Enter to continue...")
	termbox.PollEvent()
}

func ListInformation() {
	utils.MyIp()
	utils.Status()
	utilsTermbox.PrintMessage("Press Enter to continue...")
	termbox.PollEvent()
}

func Quit() {
	utilsTermbox.ClearMessage()
	os.Exit(0)
}
