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

	raw_data = bytes.ReplaceAll(raw_data, []byte("\t"), []byte(" ")) // REMOVE OS TABS
	raw_data = bytes.ReplaceAll(raw_data, []byte("\r"), []byte(" ")) // REMOVE OS CARRIAGE RETURN

	/*
		ANOTAÇÕES:
		faz a "tokenização" do raw_data em tipo do data, onde cada "primeiro slice"([esse][]string) é uma linha do arquivo, e cada "segundo slice"([]string) é uma frase da linha , e cada string é uma palavra/expressão (separada por espaço)
		o que tem dentro de '"' é mantido com aspas no inicio e fim
	*/

	var data [][][]byte = scan_data(raw_data)

	// PRINT
	fmt.Printf("====== %d linhas ======\n", len(data))
	for i_line := 0; i_line < len(data); i_line++ {
		fmt.Printf("Linha %d: ", i_line)
		for i_word := 0; i_word < len(data[i_line]); i_word++ {
			fmt.Printf("[%s] ", data[i_line][i_word])
		}
		fmt.Println()
	}
}

func (vm *VirtualMachine) Run() {
	vm.load_memory()
}

func terraform_line(line []byte) []byte {
	if len(line) == 0 {
		return []byte{}
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

			// se recentemente abriu aspas, adicionamos um espaço antes
			// isso é para evitar que as aspas fiquem grudadas com a palavra (word)
			if is_quote_opened && len(processed_line) > 0 && processed_line[len(processed_line)-1] != ' ' {
				processed_line = append(processed_line, ' ')
			}

			// enfim, adicionamos a aspas
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
				if last_char != ',' { // esse if consegue evitar duplos espaços entre duas virgulas
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
			}

		} else {
			processed_line = append(processed_line, char)
		}

	}

	return processed_line
}

func word_break_line(line []byte) [][]byte {
	var is_quote_opened bool = false
	var words_slice [][]byte
	var word []byte
	for i_char, char := range line {
		if char == '"' {
			is_quote_opened = !is_quote_opened
			word = append(word, char)
		} else if char == ' ' {
			if is_quote_opened {
				word = append(word, char)
			} else {
				words_slice = append(words_slice, word)
				word = []byte{} // reseta a WORD
			}
		} else { // se não for espaço, adiciona ao word

			word = append(word, char)

		}

		if i_char == len(line)-1 { // se for o ultimo char, adiciona a palavra
			if len(word) > 0 {
				words_slice = append(words_slice, word)
			}
		}
	}

	return words_slice
}

func scan_data(raw_data []byte) [][][]byte {

	var new_data [][][]byte = [][][]byte{}
	var lines [][]byte = bytes.Split(raw_data, []byte("\n"))
	for _, line := range lines {
		//PROCESSA CADA LINHA, REMOVENDO COMENTARIOS, ESPAÇOS, VIRGULAS MAL COLOCADAS, ETC
		var processed_line []byte = terraform_line(line)
		fmt.Printf("Processed Line: [%s]\n", string(processed_line))
		//BREAKS THE LINE INTO words_slice
		var broken_line [][]byte = word_break_line(processed_line)
		new_data = append(new_data, broken_line)
	}

	return new_data
}
