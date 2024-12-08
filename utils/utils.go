package utils

import (
	"bytes"
	"strconv"
)

func BytesEndsWith(data []byte, suffix []byte) bool {
	if len(data) < len(suffix) {
		return false
	}
	return string(data[len(data)-len(suffix):]) == string(suffix)
}

func BytesStartsWith(data []byte, prefix []byte) bool {
	if len(data) < len(prefix) {
		return false
	}
	return string(data[:len(prefix)]) == string(prefix)
}

func BytesProcessReplace(str []byte) []byte {
	aux := bytes.Replace(str, []byte(" "), []byte(""), -1)
	aux = bytes.Replace(str, []byte("\n"), []byte(""), -1)
	aux = bytes.Replace(aux, []byte("\t"), []byte(""), -1)
	aux = bytes.Replace(aux, []byte("\r"), []byte(""), -1)
	return aux
}

func Processint16(data []byte) (int16, bool) { // enter: "123" -> 123, true; "abc" -> 0, false
	aux, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, false
	}
	return int16(aux), true
}

func StringToDotByte(str string) (int16, bool) {
	aux, err := strconv.Atoi(str)
	if err != nil {
		return 0, false
	}
	if aux < 0 || aux > 255 {
		return 0, false
	}
	return int16(aux), true

}

func StringToDotWord(str string) (int16, bool) {
	//.word size is 16 bits, so it can be from 0 to 65535
	// enter: "123" -> 123, true; "abc" -> 0, false
	// enter: "65535" -> 65535, true; "65536" -> 0, false
	/* enter:
	0x0000	0
	0x0001	1
	0x7FFF	32,767
	0x8000	-32,768
	0xFFFF	-1
	0xFFFE	-2
	0x8001	-32,767
	0x7FFE	32,766
	*/

	val, err := strconv.ParseInt(str, 0, 16)
	if err != nil {
		return 0, false
	}

	// Verifica se está no intervalo válido para um int16 (-32768 a 32767).
	if val < -32768 || val > 32767 {
		return 0, false
	}

	// Converte o valor para int16 e retorna true.
	return int16(val), true

}
