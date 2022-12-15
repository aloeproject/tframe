package tnet

import (
	"encoding/json"
	"github.com/aloeproject/tframe/iface"
)

var _ iface.IPacket = (*OnceQueryPack)(nil)

/*
不需要拆包的pack
*/
type OnceQueryPack struct {
}

func (o OnceQueryPack) Unpack(binaryData []byte) (iface.IMessage, error) {
	return &OnceQueryMessage{data: binaryData}, nil
}

func (o OnceQueryPack) Pack(msg iface.IMessage) ([]byte, error) {
	/*
		根据相应业务
	*/
	//例如需要 json 发送
	by, err := json.Marshal(msg.GetData())
	return by, err
}

func (o OnceQueryPack) GetHeadLen() int {
	return 0 //不需要拆包
}

func (o OnceQueryPack) GetMaxDataLen() int {
	return 2048
}

var _ iface.IMessage = (*OnceQueryMessage)(nil)

type OnceQueryMessage struct {
	data []byte
	m    int32
}

func (o *OnceQueryMessage) GetMID() int32 {
	return o.m
}

func (o *OnceQueryMessage) SetMID(i int32) {
	o.m = i
}

func (o *OnceQueryMessage) GetDataLen() int {
	return len(o.data)
}

func (o *OnceQueryMessage) GetData() []byte {
	return o.data
}

func (o *OnceQueryMessage) GetHeadData() []byte {
	return nil
}

func (o *OnceQueryMessage) SetHeadData(bytes []byte) {

}

func (o *OnceQueryMessage) SetData(bytes []byte) {
}

func (o *OnceQueryMessage) SetHeadLen(i int) {
}
