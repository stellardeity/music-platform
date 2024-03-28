package discover

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"rockwall/proto"
	"time"
)

var peers = make(map[string]string)

func StartDiscover(p *proto.Proto) {
	go startMeow("224.0.0.1:35035", p)
	go listenMeow("224.0.0.1:35035", p, connectToPeer)
}

func connectToPeer(p *proto.Proto, peerAddress string) {
	if _, exist := peers[peerAddress]; exist {
		log.Printf("peer %s already exist", peerAddress)
		return
	}
	peers[peerAddress] = peerAddress
	log.Printf("try to connect peer: %s", peerAddress)

	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		log.Printf("Dial ERROR: " + err.Error())
		return
	}

	defer conn.Close()

	peer := handShake(p, conn)

	if peer == nil {
		log.Printf("Fail on handshake")
		return
	}

	p.RegisterPeer(peer)
	p.ListenPeer(peer)
	p.UnregisterPeer(peer)
	delete(peers, peerAddress)
}

func startMeow(address string, p *proto.Proto) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Printf(err.Error())
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf(err.Error())
	}

	for {
		_, err := conn.Write([]byte(fmt.Sprintf("meow:%v:%v", hex.EncodeToString(p.PubKey), p.Port)))
		if err != nil {
			log.Printf(err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

func listenMeow(address string, p *proto.Proto, handler func(p *proto.Proto, peerAddress string)) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.SetReadBuffer(1024)
	if err != nil {
		log.Fatal(err)
	}

	for {
		buffer := make([]byte, 1024)
		_, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		trim := bytes.Trim(buffer, "\x00")

		peerPubKeyStr := string(trim[5 : 5+64])
		peerPubKey, err := hex.DecodeString(peerPubKeyStr)
		if err != nil {
			log.Printf("DecodeHexString failed: %s", err)
			continue
		}

		_, found := p.Peers.Get(string(peerPubKey))
		if found || bytes.Equal(p.PubKey, peerPubKey) {
			continue
		}

		peerAddress := src.IP.String() + string(trim[5+64:])

		handler(p, peerAddress)
	}
}

func handShake(p *proto.Proto, conn net.Conn) *proto.Peer {
	log.Printf("DISCOVERY: try handshake with %s", conn.RemoteAddr())
	peer := proto.NewPeer(conn)

	p.SendName(peer)

	envelope, err := proto.ReadEnvelope(bufio.NewReader(conn))
	if err != nil {
		log.Printf("Error on read Envelope: %s", err)
		return nil
	}

	if string(envelope.Cmd) == "HAND" {
		if _, found := p.Peers.Get(string(envelope.From)); found {
			log.Printf(" - - - - - - - - - - - - - - - --  -- - - - - Peer (%s) already exist", peer)
			return nil
		}
	}

	err = peer.UpdatePeer(envelope)
	if err != nil {
		log.Printf("HandShake error: %s", err)
	}

	return peer
}
