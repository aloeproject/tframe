package tnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func getData() []byte {
	p := DefaultPack{}

	type data struct {
		head   [2]byte
		length int32
	}

	raw := []byte("12345")
	length := len(raw)
	d := data{
		head:   [2]byte{'a', 'b'},
		length: int32(length),
	}

	headBuff := bytes.NewBuffer([]byte{})
	err := binary.Write(headBuff, binary.BigEndian, d)

	dataBuff := bytes.NewBuffer([]byte{})
	err = binary.Write(dataBuff, binary.BigEndian, raw)

	msg := DefaultMessage{
		head:       headBuff.Bytes(),
		body:       dataBuff.Bytes(),
		bodyLength: length,
	}
	pData, err := p.Pack(&msg)
	if err != nil {
		return []byte{}
	}
	return pData
}

func TestDefaultPack_Pack(t *testing.T) {
	t.Log(getData())
}

func TestDefaultPack_Unpack(t *testing.T) {
	str := getData()
	p := DefaultPack{}
	t.Log(p.Unpack(str))
	t.Log(fmt.Sprintf("%x",str))
}
