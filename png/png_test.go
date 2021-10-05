package png

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPng_NextChunk(t *testing.T) {
	input := []byte{0, 0, 0, 4,
					0, 0, 0, 8,
					9, 1, 2, 3,
					4, 5, 6, 7}
	r := bytes.NewReader(input)
	var a = Png{A:r}
	assert := assert.New(t)

	chunk, err := a.NextChunk()
	assert.NoError(err)
	assert.Equal(uint32(4), chunk.Length)
	assert.Equal(uint32(8), chunk.Type)
	assert.Equal([]byte{9, 1, 2, 3}, chunk.Data)
	assert.Equal([]byte{4, 5, 6, 7}, chunk.CRC)
}

func TestSanity(t *testing.T) {
	assert := assert.New(t)
	input := []byte{1, 2, 3, 4, 5}
	r := bytes.NewReader(input)
	p := make([]byte, 1024)

	num, err := r.Read(p)
	assert.NoError(err)
	assert.Equal(5, num)
}

func TestNewPng(t *testing.T) {
	assert := assert.New(t)
	input := []byte{0x89, 0x50, 0x4E, 0x47,
					0x0D, 0x0A, 0x1A, 0x0A,
					0, 0, 0, 13,	//Length
					0, 0, 0, 8,		//Type
					0, 0, 0, 9,		//Data
					0, 0, 0, 9,
					8, 2, 1, 1, 1,
					4, 5, 6, 7}		//CRC
	r := bytes.NewReader(input)
	p, err := NewPng(r)
	assert.NoError(err)
	assert.Equal(Png{A:r,Parameters: &IHDR{
		Width:         9,
		Height:        9,
		Depth:         8,
		Color:         2,
		CompMeth:      1,
		FilterMeth:    1,
		InterfaceMeth: 1,
	}}, *p)
}