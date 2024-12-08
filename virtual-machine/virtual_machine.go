package virtualmachine

import (
	"bytes"
	"iguanazilla/logkit"
	"os"
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
	lk := logkit.NewLogkit("load-memory:")

	raw_data, err := os.ReadFile(vm.source_file)
	if err != nil {
		panic(err)
	}

	/*
		======== PROCESS THE CONTENT OF THE FILE (REMOVE COMMENTS, MULTI SPACES, TABS, ETC)
	*/

	// remove comments
	/*
		for bytes.Contains(raw_data, []byte("#")) {
			comment_start := bytes.Index(raw_data, []byte("#"))
			comment_end := bytes.Index(raw_data[comment_start:], []byte("\n"))
			if comment_end == -1 {
				raw_data = raw_data[:comment_start]
			} else {
				raw_data = append(raw_data[:comment_start], raw_data[comment_start+comment_end:]...)
			}
		}
	*/
	// remove \t
	raw_data = bytes.ReplaceAll(raw_data, []byte("\t"), []byte(" "))
	//
	var data [][]string
	/*
		ANOTAÇÕES:
		faz a "tokenização" do raw_data em tipo do data, onde cada "primeiro slice"([esse][]string) é uma linha do arquivo, e cada "segundo slice"([]string) é uma frase da linha , e cada string é uma palavra/expressão (separada por espaço)
		o que tem dentro de '"' é mantido com aspas no inicio e fim

		e.g.:
		.data
		macaco: .asciiz "tome  "
		chipanze: .word 20
		orango: .byte 30
		#comentario de linha (será um slice vazio)
		.text
		loop:
		ADDD           20 #comentario

		FICARIA
		[
			[".data"],
			["macaco:", ".asciiz", "\"tome  \""],
			["chipanze:", ".word", "20"],
			[...]
	*/
	lines := bytes.Split(raw_data, []byte("\n"))
	for i_line, line := range lines {

		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			data = append(data, []string{})
			continue
		}

		var tokens []string
		var token []byte
		var is_in_quotes bool = false

		for i := 0; i < len(line); i++ {
			//checa se no final da linha tem aspas, se não tiver, dá erro
			if i == len(line)-1 && is_in_quotes && line[i] != '"' || i == len(line)-1 && is_in_quotes && line[i] == '"' {
				lk.LineError(uint16(i_line+1), "Missing closing quotes")
				os.Exit(1)
			}

			if line[i] == '#' {
				break
			}
			if line[i] == '"' {
				is_in_quotes = !is_in_quotes

				token = append(token, line[i])
			} else if line[i] == ' ' && !is_in_quotes {
				if len(token) > 0 {
					tokens = append(tokens, string(token))
					token = nil
				}
			} else {
				token = append(token, line[i])
			}
		}

		if len(token) > 0 {
			tokens = append(tokens, string(token))
		}

		data = append(data, tokens)
	}
	//printa
	for i_line, line := range data {
		for _, token := range line {
			println(i_line+1, "-->", token)
		}
	}

}

func (vm *VirtualMachine) first_pass() {
	/*
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
	*/
}

func (vm *VirtualMachine) Run() {
	vm.load_memory()
	vm.first_pass()
}
