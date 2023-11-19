package main

import (
	"fmt"
	"log"
	"tailscale/menu"
	"tailscale/utils"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}

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

// go build -ldflags "-s -w" -o sky-tailscale.exe main.go
// 管理員身分執行
// https://hksanduo.github.io/2021/04/28/2021-04-28-run-go-windows-software-with-administrator-priviledge/
// 1、获取rsrc
// go install github.com/akavel/rsrc@latest
// 2、编译 *.syso
// rsrc -manifest ./.manifest -o test.syso
// 3、go编译打包
// go build

// https://stackoverflow.com/questions/25602600/how-do-you-set-the-application-icon-in-golang
// icon
// go install github.com/tc-hib/go-winres@latest
// go-winres simply --icon sky-tailscale.png
// go build
