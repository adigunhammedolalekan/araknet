package araknet

import (
	"encoding/json"
	"fmt"
	"net"
)

//Peer - is used to communicate with a network client
type Peer struct {
	conn   net.Conn
	Name   string
	closed bool
}

//NewPeer creates a new peer instance
func NewPeer(conn net.Conn, name string) *Peer {

	return &Peer{
		conn: conn, Name: name, closed: false,
	}
}

//Send sends @Param data to this peer
func (p *Peer) Send(data []byte) {
	p.conn.Write(data)
}

//Close closes up this peer. No message would be sent after it is closed
func (p *Peer) Close() {

	if p.closed {
		return
	}

	if p.conn != nil {
		p.conn.Close()
		p.closed = true
	}
}

func (p *Peer) IsClosed() bool {
	return p.closed
}

//Run listens for messages and
func (p *Peer) Run() {

	decoder := json.NewDecoder(p.conn)

	for {

		message := &Message{}
		err := decoder.Decode(message)
		if err != nil {
			fmt.Println("Failed to decode message => ", err)
			fmt.Println("Closing Connection...")
			break //
		}

		if p.Name != message.Name {

			data := message.Format()
			fmt.Println(data)
		}
	}
}
