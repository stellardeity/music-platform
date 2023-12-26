package main

import (
	"crypto/ed25519"
	"net"
)

type Peer struct {
	Name   string
	Conn   *net.Conn
	PubKey ed25519.PublicKey
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		PubKey: nil,
		Conn:   &conn,
		Name:   conn.RemoteAddr().String(),
	}
}
