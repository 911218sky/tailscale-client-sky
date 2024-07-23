package menu

import (
	"strings"

	"tailscale/utils"
	"tailscale/utils/utilsTermbox"

	"github.com/nsf/termbox-go"
)

const (
	CONNECT          = iota // 0: Connect
	SWITCHACCOUNT           // 1: Switch Account
	SIGNOUT                 // 2: Sign Out
	LIST_INFORMATION        // 3: List Information
	OPEN_MSTSC              // 4: Open Remote Desktop
	QUIT                    // 5: Quit
)

var cm = utilsTermbox.Td.ClearMessage
var pm = utilsTermbox.Td.PrintMessage

// RunTermboxUI starts the Termbox user interface.
func RunTermboxUI() {
	options := []string{"Connect", "Switch Account", "Sign Out", "List Information", "Open Remote Desktop", "Quit"}
	optionToAction := map[int]func(){
		CONNECT:          Connect,
		SWITCHACCOUNT:    SwitchAccount,
		SIGNOUT:          SignOut,
		LIST_INFORMATION: ListInformation,
		OPEN_MSTSC:       utils.OpenMstsc,
	}
	selectedIndex := 0

	for {
		cm(utilsTermbox.Option{Flush: false})
		RenderMenu(options, selectedIndex)

		event := termbox.PollEvent()
		isEnter := handleKeyEvent(event, &selectedIndex, options)

		if !isEnter {
			continue
		}

		if selectedIndex == QUIT {
			cm()
			return
		}

		action, found := optionToAction[selectedIndex]
		if !found {
			pm("Unknown option")
			break
		}

		cm()
		action()
		cm()
	}
}

// RenderMenu displays menu items.
func RenderMenu(options []string, selectedIndex int) {
	for i, option := range options {
		if selectedIndex == i {
			option = ">  " + option
		}
		pm(option, utilsTermbox.MessageOption{Flush: false, NewLine: true})
	}
	termbox.Flush()
}

// handleKeyEvent handles keyboard events for selecting menu items.
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
			*selectedIndex = QUIT
		}
	}
	return false
}

// getAccount retrieves the Tailscale account to switch to.
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
		cm(utilsTermbox.Option{Flush: false})
		pm("Account : ")
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
		pm("It is not possible to select an account that is currently in use!")
		pm("Press Enter to continue...")
		termbox.PollEvent()
		return ""
	}
	return tailscaleAccount.AllAccounts[selectedIndex]
}

// Connect connects to Tailscale.
func Connect() {
	isLogin := utils.Login()
	if !isLogin {
		return
	}
	utils.Status()
	utils.OpenMstsc()
	pm("Press Enter to continue...")
	termbox.PollEvent()
}

// SwitchAccount switches Tailscale account.
func SwitchAccount() {
	account := getAccount()
	if account == "" {
		return
	}
	utils.SwitchAccount(account)
	utils.Status()
	utils.OpenMstsc()
	pm("Press Enter to continue...")
	termbox.PollEvent()
}

// SignOut logs out of the Tailscale account.
func SignOut() {
	utils.Logout()
	pm("Press Enter to continue...")
	termbox.PollEvent()
}

// ListInformation displays Tailscale-related information.
func ListInformation() {
	utils.MyIp()
	utils.Status()
	pm("Press Enter to continue...")
	termbox.PollEvent()
}
