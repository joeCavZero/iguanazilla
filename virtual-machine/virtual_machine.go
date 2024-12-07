package virtualmachine

import (
	"bytes"
	"fmt"
	"iguanazilla/logkit"
	"iguanazilla/utils"
	"os"
	"strconv"
)

const (
	MAX_MEMORY = 4096
)

type VirtualMachine struct {
	source_file string

	accumulator  int16
	stack        [MAX_MEMORY]int16
	sp           uint32
	pc           uint32
	memory       [MAX_MEMORY]Instruction
	memory_size  uint32
	symbol_table SymbolTable
}

func NewVirtualMachine(source_file string) *VirtualMachine {
	return &VirtualMachine{
		source_file: source_file,
		pc:          0,
		memory:      [MAX_MEMORY]Instruction{},
	}
}

func (vm *VirtualMachine) load_memory() {
	raw_data, err := os.ReadFile(vm.source_file)
	if err != nil {
		panic(err)
	}

	/*
		PROCESS THE CONTENT OF THE FILE (REMOVE COMMENTS, MULTI SPACES, TABS, ETC)
	*/

	// remove comments
	for bytes.Contains(raw_data, []byte("#")) {
		comment_start := bytes.Index(raw_data, []byte("#"))
		comment_end := bytes.Index(raw_data[comment_start:], []byte("\n"))
		if comment_end == -1 {
			raw_data = raw_data[:comment_start]
		} else {
			raw_data = append(raw_data[:comment_start], raw_data[comment_start+comment_end:]...)
		}
	}

	// remove \t
	raw_data = bytes.ReplaceAll(raw_data, []byte("\t"), []byte(" "))
	// decrease ' ' queue to 1 and
	for bytes.Contains(raw_data, []byte("  ")) {
		raw_data = bytes.ReplaceAll(raw_data, []byte("  "), []byte(" "))
	}
	// remove ' ' from start and end of lines
	raw_data = bytes.ReplaceAll(raw_data, []byte(" \n"), []byte("\n"))
	raw_data = bytes.ReplaceAll(raw_data, []byte("\n "), []byte("\n"))

	//var data [][]byte
	for i_line, line := range bytes.Split(raw_data, []byte("\n")) {
		for _, expr := range bytes.Split(line, []byte(" ")) {
			new_instruction := Instruction{
				Line:     uint16(i_line),
				Operator: expr,
			}
			vm.memory[vm.memory_size] = new_instruction
			vm.memory_size++
		}
	}

	for i := 0; i < int(vm.memory_size); i++ {
		fmt.Println(vm.memory[i].Line, "-->", string(vm.memory[i].Operator))
	}
}

func (vm *VirtualMachine) first_pass() {
	lk := logkit.NewLogkit("first-pass")

	var label []byte = nil
	var value_type byte = '0'
	var value_ptr int16 = 0

	for i := 0; i < int(vm.memory_size); i++ {

		if label != nil && value_type != '0' {
			vm.symbol_table.Add(label, value_type, value_ptr)
		}

		if utils.BytesEndsWith(vm.memory[i].Operator, []byte(":")) {
			label = vm.memory[i].Operator
		}

		if bytes.Equal(vm.memory[i].Operator, []byte(".byte")) {
			value_type = 'b'
		} else if bytes.Equal(vm.memory[i].Operator, []byte(".word")) {
			value_type = 'w'
		} else if bytes.Equal(vm.memory[i].Operator, []byte(".asciiz")) {
			value_type = 'a'
		} else {
			interg, err := strconv.Atoi(string(vm.memory[i].Operator))
			if err != nil {
				lk.LineError(vm.memory[i].Line, "Invalid instruction")
				os.Exit(1)
			} else {
				vm.stack[vm.sp] = int16(interg)
				vm.sp++
			}
		}

	}
}

func (vm *VirtualMachine) add_stack_value(value int16) {
	vm.stack[vm.sp] = value
	vm.sp++
}

func (vm *VirtualMachine) Run() {
	vm.load_memory()
	vm.first_pass()
}
