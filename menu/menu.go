package menu

import (
	"os"
	"strings"

	"tailscale/utils"

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
		utils.ClearMessage(utils.Option{NoFlush: true})
		RenderMenu(options, selectedIndex)

		event := termbox.PollEvent()
		isEnter := handleKeyEvent(event, &selectedIndex, options)

		if !isEnter {
			continue
		}

		action, found := optionToAction[selectedIndex]
		if !found {
			utils.PrintMessage("Unknown option")
			break
		}

		utils.ClearMessage()
		action()
		utils.ClearMessage()
	}
}

func RenderMenu(options []string, selectedIndex int) {
	for i, option := range options {
		if selectedIndex == i {
			option = ">  " + option
		}
		utils.PrintMessage(option, utils.Option{NoFlush: true})
	}
	termbox.Flush()
}

func handleKeyEvent(event termbox.Event, selectedIndex *int, options []string) bool {
	if event.Type == termbox.EventKey {
		switch event.Key {
		case termbox.KeyArrowUp:
			if *selectedIndex > 0 {
				*selectedIndex--
			}
		case termbox.KeyArrowDown:
			if *selectedIndex < len(options)-1 {
				*selectedIndex++
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
		utils.ClearMessage(utils.Option{NoFlush: true})
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
		utils.PrintMessage("It is not possible to select an account that is currently in use!")
		utils.PrintMessage("Press Enter to continue...")
		termbox.PollEvent()
		return ""
	}

	return tailscaleAccount.AllAccounts[selectedIndex]
}

func Connect() {
	utils.Login()
	utils.Status()
	utils.OpenMstsc()
	utils.PrintMessage("Press Enter to continue...")
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
	utils.PrintMessage("Press Enter to continue...")
	termbox.PollEvent()
}

func SignOut() {
	utils.Logout()
	utils.PrintMessage("Press Enter to continue...")
	termbox.PollEvent()
}

func ListInformation() {
	utils.MyIp()
	utils.Status()
	utils.PrintMessage("Press Enter to continue...")
	termbox.PollEvent()
}

func Quit() {
	utils.ClearMessage()
	os.Exit(0)
}
