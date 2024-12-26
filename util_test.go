package mediawiki

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddLineNumber(t *testing.T) {
	lineNumber := 10
	data := []byte("Hello, World!")

	result := AddLineNumber(lineNumber, data)

	expected := make([]byte, 4)
	binary.LittleEndian.PutUint32(expected, uint32(lineNumber))
	expected = append(expected, data...)

	assert.Equal(t, expected, result, "The result should match the expected format.")
}

func TestParseLineNumber(t *testing.T) {
	lineNumber := 113287550
	data := []byte("This is a test message.")
	combined := AddLineNumber(lineNumber, data)

	parsedLineNumber, parsedData, err := ParseLineNumber(combined)

	assert.Nil(t, err, "There should be no error when parsing.")
	assert.Equal(t, lineNumber, parsedLineNumber, "The parsed line number should match the expected one.")
	assert.Equal(t, data, parsedData, "The parsed data should match the original data.")
}

func TestParseLineNumber_Error(t *testing.T) {
	shortData := []byte("s")

	parsedLineNumber, parsedData, err := ParseLineNumber(shortData)

	assert.Error(t, err, "There should be an error because the data is too short.")
	assert.Equal(t, 0, parsedLineNumber, "The parsed line number should be 0 due to the error.")
	assert.Nil(t, parsedData, "Parsed data should be nil due to the error.")
}
