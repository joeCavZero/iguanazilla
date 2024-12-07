package utils

import "bytes"

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
