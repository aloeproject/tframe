package tnet

import "github.com/aloeproject/tframe/iface"

var _ iface.IMessage = (*DefaultMessage)(nil)

type DefaultMessage struct {
	head []byte
	body []byte
	mid  int32

	bodyLength int
}

func (d *DefaultMessage) GetMID() int32 {
	return d.mid
}

func (d *DefaultMessage) SetMID(i int32) {
	d.mid = i
}

func (d *DefaultMessage) SetHeadLen(i int) {
	d.bodyLength = i
}

func (d *DefaultMessage) SetHeadData(bytes []byte) {
	d.head = bytes
}

func (d *DefaultMessage) GetHeadData() []byte {
	return d.head
}

func (d *DefaultMessage) GetDataLen() int {
	return d.bodyLength
}

func (d *DefaultMessage) GetData() []byte {
	return d.body
}

func (d *DefaultMessage) SetData(bytes []byte) {
	d.body = bytes
}
