package main

import (
	"fmt"
	"myruscoint/internal/emulator"
	"myruscoint/internal/ruscoin"
)

func main() {
	if err := ruscoin.InitRuscoinSettings(); err != nil {
		fmt.Println(err)
	}
	if err := emulator.LoadSettingsFromEnv(); err != nil {
		fmt.Println(err)
	}
	wb := emulator.NewEmulatorWeb().DefaultRcManager()
	if emulator.WITH_LOG {
		wb.StartWithLogger()
	} else {
		wb.Start()
	}
}
