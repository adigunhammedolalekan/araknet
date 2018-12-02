package araknet

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

var (
	TYPE_JOIN_NETWORK   = "join"
	TYPE_SEND_MESSAGE   = "message"
	TYPE_CREATE_NETWORK = "create_network"
)

//Message is a message sent from connected peers. Message has 3 types
// - TYPE_JOIN_NETWORK -> This is sent immediately a peer connects to the service.
// It allows Peer to specify which network he wants to join
//- TYPE_CREATE_NETWORK -> Command sent by peers to create a network
//
type Message struct {
	Type    string      `json:"type"`

	//Network name
	Network string      `json:"network"`

	//Name of the peer that sent this message
	Name    string      `json:"name"`

	//Content of the message
	Data    interface{} `json:"data"`

	Secret string `json:"secret"`
}

//String - Just to debug
func (m *Message) String() string {
	return m.Network + " " + m.Secret
}

//
func (m *Message) Format() string {

	data, ok := m.Data.(map[string]interface{})
	if !ok {
		return ""
	}

	return fmt.Sprintf("%s > %s", m.Name, data["message"])
}

//ArakNet is the house of all network.
type ArakNet struct {

	//All networks in ArakNet mapped with their name
	Networks map[string]*Network
	mtx      sync.Mutex

	//Server Address
	Address  string
}

//New creates a new ArakNet instance
func New(address string) *ArakNet {

	return &ArakNet{
		Networks: make(map[string]*Network, 0),
		Address:  address,
	}

}

//Connect open network to the outside world and allows connections from peers
func (a *ArakNet) Connect() {

	listener, err := net.Listen("tcp", a.Address)
	if err != nil {
		log.Fatalf("Failed to connect to => %s, Reason %s ", a.Address, err.Error())
	}

	defer listener.Close()
	a.Listen(listener)
}

//Listen watches and process new Connections
func (a *ArakNet) Listen(listener net.Listener) {

	for {

		conn, err := listener.Accept()
		if err != nil {

			continue
		}

		go a.HandleConnection(conn)
	}
}

//HandleConnection processes connection.
func (a *ArakNet) HandleConnection(conn net.Conn) {

	for {

		decoder := json.NewDecoder(conn)
		message := &Message{}
		err := decoder.Decode(message)
		if err != nil {
			fmt.Println("Error => ", err)
			break
		}

		fmt.Println("Received => ", message.String())
		switch message.Type {

		case TYPE_JOIN_NETWORK:
			p := NewPeer(conn, message.Name)
			err := a.AuthorizePeer(message.Network, message.Secret, p)
			if err != nil {
				p.Send([]byte(err.Error()))
				p.Close()
			}
		case TYPE_SEND_MESSAGE:
			a.SendMessage(message)

		case TYPE_CREATE_NETWORK:
			a.CreateNetwork(conn, message)
		}
	}
}

//AuthorizePeer checks if peer has the secret key to join a network given by name @Param network
//and then add this peer to the network or return error if the secret is incorrect or network does not exists
func (a *ArakNet) AuthorizePeer(network, secret string, peer *Peer) error {

	a.Lock()
	defer a.UnLock()

	n := a.Networks[network]
	if n == nil {
		return errors.New("Network not found")
	}

	if n.Secret != secret {
		return errors.New("Unauthorized access, invalid credentials")
	}

	n.AddPeer(peer)
	return nil
}

//CreateNetwork - creates a new network.
// A network can allow unlimited number of peers to join and communicate/send message securely
func (a *ArakNet) CreateNetwork(conn net.Conn, message *Message) {

	a.Lock()
	defer a.UnLock()

	n := a.Networks[message.Network]
	if n != nil {

		conn.Write([]byte("Network already exists"))
		return
	}

	network := NewNetwork(message.Network, message.Secret)
	peer := NewPeer(conn, message.Name)

	//add network creator to peers
	network.AddPeer(peer)
	a.Networks[message.Network] = network
}

//SendMessage - Send message to all peers in a network
func (a *ArakNet) SendMessage(message *Message) {

	a.Lock()
	defer a.UnLock()

	network := a.Networks[message.Network]
	if network != nil {
		network.Send(message.Data)
	}
}

func (a *ArakNet) Lock() {
	a.mtx.Lock()
}

func (a *ArakNet) UnLock() {
	a.mtx.Unlock()
}
