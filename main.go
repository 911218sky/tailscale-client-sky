package main

import (
	"fmt"
	"os"
	"tailscale/menu"
	"tailscale/utils"
	"tailscale/utils/debug"
	"tailscale/utils/utilsTermbox"

	"github.com/nsf/termbox-go"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-d" {
		debug.Debug()
		return
	}

	defer termbox.Close()
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	utilsTermbox.InIt()
	utils.CheckTailscale()

	accounts, err := utils.GetAccounts()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else if len(accounts.AllAccounts) == 0 {
		menu.Connect()
	} else if len(accounts.AllAccounts) == 1 {
		utils.SwitchAccount(accounts.AllAccounts[0])
	}

	menu.RunTermboxUI()
}
