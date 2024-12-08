package main

import (
	"iguanazilla/logkit"
	virtualmachine "iguanazilla/virtual-machine"

	"os"
)

func main() {
	argc := len(os.Args)

	lk := logkit.NewLogkit("virtual-machine")
	if argc == 2 {
		lk.Info("initializing iguanazilla virtual machine")

		iguanazilla := virtualmachine.NewVirtualMachine(os.Args[1])
		iguanazilla.Run()

	} else {
		lk.Error("Usage: iguanazilla [file]")
	}
}
