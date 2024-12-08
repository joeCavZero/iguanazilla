package virtualmachine

import (
	"bytes"
	"fmt"
	"os"
)

const (
	MAX_MEMORY = 4096
)

type VirtualMachine struct {
	source_file string

	accumulator    int16
	stack          [MAX_MEMORY]int16
	sp             uint32
	pc             uint32
	memory         [MAX_MEMORY]int16
	memory_pointer uint32
	symbol_table   SymbolTable
}

func NewVirtualMachine(source_file string) *VirtualMachine {
	return &VirtualMachine{
		source_file: source_file,

		accumulator:    0,
		stack:          [MAX_MEMORY]int16{},
		sp:             MAX_MEMORY,
		pc:             0,
		memory:         [MAX_MEMORY]int16{},
		memory_pointer: 0,
		symbol_table:   NewSymbolTable(),
	}
}

func (vm *VirtualMachine) load_memory() {
	raw_data, err := os.ReadFile(vm.source_file)
	if err != nil {
		panic(err)
	}

	/*
		======== PROCESS THE CONTENT OF THE FILE (REMOVE COMMENTS, MULTI SPACES, TABS, ETC)
	*/
	// remove \t
	raw_data = bytes.ReplaceAll(raw_data, []byte("\t"), []byte(" ")) // REMOVE OS TABS
	raw_data = bytes.ReplaceAll(raw_data, []byte("\r"), []byte(" ")) // REMOVE OS CARRIAGE RETURN
	//raw_data = bytes.ReplaceAll(raw_data, []byte("\n"), []byte(" ")) // REMOVE OS NEW LINE

	//fmt.Println("-->", string(raw_data))

	/*
		ANOTAÇÕES:
		faz a "tokenização" do raw_data em tipo do data, onde cada "primeiro slice"([esse][]string) é uma linha do arquivo, e cada "segundo slice"([]string) é uma frase da linha , e cada string é uma palavra/expressão (separada por espaço)
		o que tem dentro de '"' é mantido com aspas no inicio e fim

		e.g.:
		.data
		macaco: .asciiz "tome   " #COMENTARIO
		sagui: .asciiz "hi"
		chipanze: .word 20,30,40,50
		bonobo: .word 10, 11
		mico: .word 1 , 2, 3 , 4,   5
		orango: .byte 30
		#MACACOOOOO
		.text
		loop:
		    ADDD           20 #meo
		SUBD 10

		FICARIA
		[
			[".data"],
			["macaco:", ".asciiz", "\"tome  \""],
			["chipanze:", ".word", "20", ",", "30", ",", "40", ",", "50"],
			[...]
		]
	*/

	/*
		TERRAFORMING THE DATA
	*/

	var data [][]string
	lines := bytes.Split(raw_data, []byte("\n"))
	for _, line := range lines {
		/*
			PROCESSA CADA LINHA, REMOVENDO COMENTARIOS, ESPAÇOS, VIRGULAS MAL COLOCADAS, ETC
		*/
		if len(line) == 0 {
			data = append(data, []string{})
			continue
		}

		var is_quote_opened bool = false
		var processed_line []byte
		for i_char, char := range line {

			if char == ' ' && len(processed_line) == 0 { // ignora espaços no inicio da linha
				continue
			}
			if char == ' ' && i_char == len(line)-1 { // ignora espaços no final da linha
				continue
			}
			if char == '#' { // se for #, ignora o resto da linha (comentario)
				break
			} else if char == '"' { // se for aspas, inverte o estado de is_quote_opened
				is_quote_opened = !is_quote_opened
				processed_line = append(processed_line, char)
			} else if char == ' ' {
				if is_quote_opened {
					processed_line = append(processed_line, char)
				} else {
					var last_char byte = 0
					if len(processed_line) > 0 {
						last_char = processed_line[len(processed_line)-1]
					}
					if last_char == ' ' { // se o ultimo char for espaço, ignora
						continue
					} else { // se o ultimo char não for espaço, adiciona
						processed_line = append(processed_line, char)
					}
				}
			} else if char == ',' {
				if is_quote_opened {
					processed_line = append(processed_line, char)
				} else {
					var last_char byte = 0
					var next_char byte = 0
					if len(processed_line) > 0 {
						last_char = line[i_char-1]
					}
					if i_char < len(line)-1 {
						next_char = line[i_char+1]
					}

					if last_char == ' ' && next_char == ' ' {
						processed_line = append(processed_line, char)
					} else if last_char == ' ' && next_char != ' ' {
						processed_line = append(processed_line, char, ' ')
					} else if last_char != ' ' && next_char == ' ' {
						processed_line = append(processed_line, ' ', char)
					} else {
						processed_line = append(processed_line, ' ', char, ' ')
					}

				}

			} else {
				processed_line = append(processed_line, char)
			}

		}

		/*
			BREAKS THE LINE INTO WORDS
		*/

		is_quote_opened = false
		var words []string
		var word []byte
		for i_char, char := range processed_line {
			if char == '"' {
				is_quote_opened = !is_quote_opened
				word = append(word, char)
			} else if char == ' ' {
				if is_quote_opened {
					word = append(word, char)
				} else {
					words = append(words, string(word))
					word = []byte{} // reseta a WORD
				}
			} else { // se não for espaço, adiciona ao word

				word = append(word, char)

			}

			if i_char == len(processed_line)-1 { // se for o ultimo char, adiciona a palavra
				if len(word) > 0 {
					words = append(words, string(word))
				}
			}
		}

		data = append(data, words)
	}
	fmt.Printf("====== %d linhas ======\n", len(data))
	for i_line := 0; i_line < len(data); i_line++ {
		line := data[i_line]
		fmt.Printf("linha %d --> ", i_line+1)
		for i, word := range line {
			fmt.Printf("[%s] ", word)
			if len(line)-1 == i {
				print("<---")
			}
		}
		fmt.Println()
	}
}

func (vm *VirtualMachine) Run() {
	vm.load_memory()
}

func (vm *VirtualMachine) add_to_stack(number int16) {
	vm.sp--
	vm.stack[vm.sp] = number

}
