package mediawiki

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func AddLineNumber(lineNumber int, data []byte) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, int32(lineNumber))
	buf.Write(data)
	return buf.Bytes()
}

func ParseLineNumber(data []byte) (int, []byte, error) {
	if len(data) < 4 {
		return 0, nil, fmt.Errorf("data too short to contain line number")
	}
	lineNumber := int(binary.LittleEndian.Uint32(data[:4]))
	return lineNumber, data[4:], nil
}
