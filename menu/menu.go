// Package menu provides functionality for a terminal-based menu interface
package menu

import (
	"strings"
	"tailscale/utils"
	"tailscale/utils/drawer"

	"github.com/nsf/termbox-go"
)

// Menu item constants representing different menu options
const (
	CONNECT          = iota // Connect to Tailscale
	SWITCHACCOUNT           // Switch between Tailscale accounts
	SIGNOUT                 // Sign out of current account
	LIST_INFORMATION        // Display Tailscale information
	OPEN_MSTSC              // Open Remote Desktop Connection
	QUIT                    // Exit the application
)

// RunTermboxUI starts the Termbox user interface and handles the main menu loop.
// It displays menu options and executes corresponding actions based on user input.
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
// Parameters:
//   - options: slice of strings containing menu options
//   - selectedIndex: index of currently selected menu item
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
// Parameters:
//   - event: termbox keyboard event
//   - selectedIndex: pointer to current selection index
//   - options: available menu options
//
// Returns:
//   - bool: true if Enter key was pressed, false otherwise
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
// It displays a list of available accounts and handles user selection.
// Returns selected account name or empty string if selection is cancelled.
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
		drawer.Clear(drawer.DefaultOptionNoFlush)
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
// It handles the login process, checks status, and opens Remote Desktop connection.
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
// It allows switching between different Tailscale accounts and updates the connection.
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
// It performs the logout operation and waits for user acknowledgment.
func SignOut() {
	utils.Logout()
	drawer.Print("Press Enter to continue...", drawer.DefaultOption)
	termbox.PollEvent()
}

// ListInformation displays Tailscale-related information to the user.
// It shows IP address and status information, then waits for user acknowledgment.
func ListInformation() {
	utils.MyIP()
	utils.Status()
	drawer.Print("Press Enter to continue...", drawer.DefaultOption)
	termbox.PollEvent()
}
