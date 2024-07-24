package debug

import (
	"fmt"
	"os"
	"runtime/trace"
	"tailscale/menu"
	"tailscale/utils"
	"tailscale/utils/drawer"

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
