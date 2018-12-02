package araknet

import (
	"fmt"
)

//Network houses peers, peers are authorized by comparing their secret key to network
//secret key. Sort of like using username/password pair to access a service.
//The creator of the network creates it with a secure secret key and share it with whomever he wants
//to communicate securely with. Secret should be shared in a face-to-face meeting
type Network struct {

	//The network name
	Name     string

	//Secret key. Will be used to encrypt all outgoing messages
	Secret   string

	//Connected peers
	Peers    []*Peer
	//Uses to encrypt outgoing messages
	enc      *Encryptor
}

//NewNetwork creates a new Network instance
func NewNetwork(name, secret string) *Network {

	return &Network{
		Name:     name,
		Secret:   secret,
		enc:      NewEncryptor(secret),
	}
}

//AddPeer Add a new peer to the network
func (n *Network) AddPeer(p *Peer) {
	n.Peers = append(n.Peers, p)
}

//Send send a {key} encrypted message to all connected peers
func (n *Network) Send(data interface{}) {

	ruffled, err := n.enc.Encrypt(data)
	if err != nil {
		fmt.Println("Failed to Encrypt data => ", err)
		return
	}

	for _, peer := range n.Peers {

		if peer != nil && !peer.IsClosed() {
			peer.Send(ruffled)
		}
	}
}
