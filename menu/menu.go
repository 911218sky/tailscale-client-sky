package menu

import (
	"strings"
	"tailscale/utils"
	"tailscale/utils/drawer"

	"github.com/nsf/termbox-go"
)

// Menu Item Constants
const (
	CONNECT          = iota // 0: Connect
	SWITCHACCOUNT           // 1: Switch Account
	SIGNOUT                 // 2: Sign Out
	LIST_INFORMATION        // 3: List Information
	OPEN_MSTSC              // 4: Open Remote Desktop
	QUIT                    // 5: Quit
)

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
		RenderMenu(options, selectedIndex)

		event := termbox.PollEvent()
		isEnter := handleKeyEvent(event, &selectedIndex, options)

		if !isEnter {
			continue
		}

		if selectedIndex == QUIT {
			drawer.Clear(drawer.DefaultOption)
			return
		}

		action, found := optionToAction[selectedIndex]
		if found {
			// Clear the screen before executing the action menu
			drawer.Clear(drawer.DefaultOptionNoFlush)
			action()
			// Clear the screen after executing the action menu
			drawer.Clear(drawer.DefaultOptionNoFlush)
		}
	}
}

// RenderMenu displays the menu items with the selected option highlighted.
func RenderMenu(options []string, selectedIndex int) {
	drawer.Clear(drawer.DefaultOptionNoFlush)
	for i, option := range options {
		if selectedIndex == i {
			option = ">  " + option // Highlight selected option
		}
		drawer.Print(option, drawer.DefaultOptionNoFlush) // Use no-flush option for performance
	}
	drawer.Flush()
}

// handleKeyEvent processes keyboard events for selecting menu items.
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

	// Mark the current account with an asterisk (*)
	tailscaleAccount.ForEach(func(account *string) {
		if *account == tailscaleAccount.CurrentAccount {
			*account = "*" + *account
		}
	})

	tailscaleAccount.AllAccounts = append(tailscaleAccount.AllAccounts, "QUIT")

	for {
		drawer.Clear(drawer.Option{Flush: false})
		drawer.Print("Account : ", drawer.DefaultOption)
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
		drawer.Print("It is not possible to select an account that is currently in use!", drawer.DefaultOption)
		drawer.Print("Press Enter to continue...", drawer.DefaultOption)
		termbox.PollEvent()
		return ""
	}
	return tailscaleAccount.AllAccounts[selectedIndex]
}

// Connect initiates the connection to Tailscale.
func Connect() {
	isLogin := utils.Login()
	if !isLogin {
		return
	}
	utils.Status()
	utils.OpenMstsc()
	drawer.Print("Press Enter to continue...", drawer.DefaultOption)
	termbox.PollEvent()
}

// SwitchAccount changes the current Tailscale account.
func SwitchAccount() {
	account := getAccount()
	if account == "" {
		return
	}
	utils.SwitchAccount(account)
	utils.Status()
	utils.OpenMstsc()
	drawer.Print("Press Enter to continue...", drawer.DefaultOption)
	termbox.PollEvent()
}

// SignOut logs the user out of the Tailscale account.
func SignOut() {
	utils.Logout()
	drawer.Print("Press Enter to continue...", drawer.DefaultOption)
	termbox.PollEvent()
}

// ListInformation displays Tailscale-related information to the user.
func ListInformation() {
	utils.MyIp()
	utils.Status()
	drawer.Print("Press Enter to continue...", drawer.DefaultOption)
	termbox.PollEvent()
}
