package main

import (
	"iguanazilla/logkit"
	virtualmachine "iguanazilla/virtual-machine"

	"os"
)

func main() {
	argc := len(os.Args)

	lk := logkit.NewLogkit("args-parsing")
	if argc == 2 {
		lk.Info("Iniciando Iguanazilla")

		iguanazilla := virtualmachine.NewVirtualMachine(os.Args[1])
		iguanazilla.Run()

	} else {
		lk.Error("Usage: iguanazilla [file]")
	}
}
