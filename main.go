package main

import (
	"fmt"
	"os"
	"tailscale/menu"
	"tailscale/utils"
	"tailscale/utils/debug"
	"tailscale/utils/drawer"
)

func main() {
	if err := drawer.Init(); err != nil {
		drawer.Print(fmt.Sprintf("Error: %v\n", err), drawer.DefaultOption)
		return
	}
	defer drawer.Close()

	if len(os.Args) > 1 && os.Args[1] == "-d" {
		debug.Debug()
		return
	}

	utils.CheckTailscale()

	accounts, err := utils.GetAccounts()
	if err != nil {
		drawer.Print(fmt.Sprintf("Error: %v\n", err), drawer.DefaultOption)
		return
	} else if len(accounts.AllAccounts) == 0 {
		menu.Connect()
	} else if len(accounts.AllAccounts) == 1 {
		utils.SwitchAccount(accounts.AllAccounts[0])
	}

	drawer.Clear(drawer.DefaultOption)
	menu.RunTermboxUI()
}
