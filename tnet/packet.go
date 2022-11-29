package tnet

import (
	"bytes"
	"encoding/binary"
	"github.com/aloeproject/tframe/iface"
)

/*
协议1
包头 2 字节
长度 4 字节 int32
数据内容
*/

var _ iface.IPacket = (*DefaultPack)(nil)

type DefaultPack struct {
}

func (d DefaultPack) GetMaxDataLen() int {
	return 0 //拆包时是0
}

func (d DefaultPack) Unpack(binaryData []byte) (iface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &DefaultMessage{}

	var head [2]byte
	//读包头
	if err := binary.Read(dataBuff, binary.BigEndian, &head); err != nil {
		return nil, err
	}

	//读长度
	var length int32
	if err := binary.Read(dataBuff, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	msg.SetHeadData(binaryData)
	msg.SetHeadLen(int(length))

	return msg, nil
}

func (d DefaultPack) Pack(msg iface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写包头
	if err := binary.Write(dataBuff, binary.BigEndian, msg.GetHeadData()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (d DefaultPack) GetHeadLen() int {
	return 6
}
