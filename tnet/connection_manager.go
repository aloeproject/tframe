package tnet

import (
	"github.com/aloeproject/tframe/iface"
	"sync"
)

var _ iface.IConnManager = (*ConnectionManager)(nil)

var MangerObj iface.IConnManager

func init() {
	MangerObj = NewConnectionManager()
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connList: make(map[string]iface.IConnection),
		lock:     new(sync.RWMutex),
	}
}

type ConnectionManager struct {
	connList map[string]iface.IConnection
	lock     *sync.RWMutex
}

func (c *ConnectionManager) Get(s string) iface.IConnection {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if conn, ok := c.connList[s]; ok {
		return conn
	}
	return nil
}

func (c *ConnectionManager) Add(conn iface.IConnection) {
	c.lock.Lock()
	c.connList[conn.GetConnId()] = conn
	c.lock.Unlock()
}

func (c *ConnectionManager) Close(conn iface.IConnection) {
	c.lock.Lock()
	if obj, ok := c.connList[conn.GetConnId()]; ok {
		obj.Stop()
		delete(c.connList, conn.GetConnId())
	}
	c.lock.Unlock()
}

func (c *ConnectionManager) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return len(c.connList)
}

func (c *ConnectionManager) ClearConn() {
	c.lock.Lock()
	for _, conn := range c.connList {
		conn.Stop()
	}
	c.connList = make(map[string]iface.IConnection)
	c.lock.Unlock()
}
