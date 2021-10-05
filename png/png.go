package png

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log"
)

type Png struct {
	A          io.Reader
	Parameters *Ihdr
}

type Chunk struct {
	Length uint32
	Type   []byte
	Data   []byte
	CRC    uint32
}

type Ihdr struct {
	Width         uint32
	Height        uint32
	Depth         uint8
	Color         uint8
	CompMeth      uint8
	FilterMeth    uint8
	InterfaceMeth uint8
}

type Text struct {
	Keyword   string
	Separator []byte
	Text      string
}

const (
	IHDR = "IHDR"
	TEXT = "tEXt"
	ZTXT = "zTXt"
)

func ChunkType(c *Chunk) (interface{}, error){

	switch string(c.Type){
	case IHDR:
		return c.parseIhdr()
	case TEXT:
	case ZTXT:
	}
	return nil, errors.New("Unknown type!!!")
}

//func checkHash(c *Chunk) bool {
//	hash := crc32.ChecksumIEEE(c.Data)
//	if hash == c.CRC {
//		return true
//	}
//	return false
//}

func checkHash(c *Chunk) bool {
	return c.CRC == crc32.ChecksumIEEE(c.Data)
}

func (c *Chunk)parseText() (*Text, error) {

	return nil, nil
}

func (c *Chunk)parseIhdr() (*Ihdr, error) {
	a := Ihdr{}
	var err error

	r := bytes.NewReader(c.Data)
	a.Width, err = readInt32(r)
	if err != nil {
		return nil, err
	}

	a.Height, err = readInt32(r)
	if err != nil {
		return nil, err
	}

	a.Depth, err = readInt8(r)
	if err != nil {
		return nil, err
	}

	a.Color, err = readInt8(r)
	if err != nil {
		return nil, err
	}

	a.CompMeth, err = readInt8(r)
	if err != nil {
		return nil, err
	}

	a.FilterMeth, err = readInt8(r)
	if err != nil {
		return nil, err
	}

	a.InterfaceMeth, err = readInt8(r)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (p *Png) NextChunk() (*Chunk, error) {
	length, err := readInt32(p.A)
	if err != nil {
		return nil, err
	}

	a := Chunk{}
	a.Length = length

	buf := make([]byte, 4)
	_, err = p.A.Read(buf)
	if err != nil {
		return nil, err
	}
	a.Type = buf


	buf = make([]byte, length)
	_, err = p.A.Read(buf)
	if err != nil {
		return nil, err
	}
	a.Data = buf

	a.CRC, err = readInt32(p.A)
	if err != nil {
		return nil, err
	}

	if checkHash(&a) {
		return nil, fmt.Errorf("Invalid HASH sum!!!")
	}

	return &a, nil
}

func readInt32(r io.Reader) (uint32, error) {
	buf := make([]byte, 4)
	_, err := r.Read(buf)
	if err != nil {
		return 0, err
	}
	num := binary.BigEndian.Uint32(buf)
	return num, nil
}

func readInt8(r io.Reader) (uint8, error) {
	buf := make([]byte, 1)
	_, err := r.Read(buf)
	if err != nil {
		return 0, err
	}

	return buf[0], nil
}

func NewPng(r io.Reader) (*Png, error) {
	v := isPng(r)
	if !v {
		return nil, fmt.Errorf("It's not a PNG!!!")
	}
	a := Png{A: r}

	c, err := a.NextChunk()
	if err != nil {
		return nil, err
	}

	cht, err := ChunkType(c)
	if err != nil {
		return nil, err
	}
	switch v := cht.(type) {
	case *Ihdr:
		a.Parameters = v
	default:
		return nil, errors.New("It isn't of IHDR type!!!")
	}

	return &a, nil
}

func isPng(r io.Reader) bool {
	val := []byte{0x89, 0x50, 0x4E, 0x47,
		0x0D, 0x0A, 0x1A, 0x0A}
	p := make([]byte, 8)

	_, err := r.Read(p)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(val); i++ {
		if p[i] != val[i] {
			return false
		}
	}

	return true
}
