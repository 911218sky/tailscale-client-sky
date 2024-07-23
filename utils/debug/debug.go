package debug

import (
	"fmt"
	"os"
	"runtime/trace"
	"tailscale/menu"
	"tailscale/utils"
	"tailscale/utils/utilsTermbox"

	"github.com/nsf/termbox-go"
)

func Debug() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	defer trace.Stop()
	defer termbox.Close()

	if err := trace.Start(f); err != nil {
		panic(err)
	}

	err = termbox.Init()
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
