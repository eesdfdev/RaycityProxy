package Network

import (
	"RCProxy/Common"
	"RCProxy/Crypto"
	"RCProxy/Logger"
	"RCProxy/Network/Server"
	"RCProxy/Stream"
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

type Client struct {
	ConnId       int
	Conn         net.Conn
	serverClient Common.IClient
	riv          uint32
	siv          uint32
}

func CalculateXorDword(a1 int64, shift byte) uint32 {
	first := a1 << (shift & 0x1F)
	second := (first >> 0x20) & 0xFFFFFFFF
	return uint32(second)
}
func (client *Client) RecvPacket() {
	r := bufio.NewReader(client.Conn)
	defer client.OnExit()
	for {
		lengthBuf := make([]byte, 4)
		_, err := io.ReadFull(r, lengthBuf)
		if err != nil {
			Logger.Errorf("Error in Packet::RecvPacket : %v", err)
			return
		}
		length := binary.LittleEndian.Uint32(lengthBuf)
		if client.riv != 0 {
			length = client.riv ^ length ^ 0xA05F33BA
		}
		//TODO : add xor encryption
		buf := make([]byte, length)
		_, err = io.ReadFull(r, buf)
		if err != nil {
			Logger.Errorf("Error in Packet::RecvPacket : %v", err)
			return
		}
		if client.riv != 0 {
			hash := Crypto.HashDecrypt(buf, length-4, client.riv)
			checksum := client.riv ^ binary.LittleEndian.Uint32(buf[length-4:length]) ^ 0x75AF23CB
			if hash != checksum {
				Logger.Errorf("checksum mismatch in Client!")
			}
			buf = buf[0 : length-4]
			client.riv += 0x1473F19
		}
		client.serverClient.RecvCallback(buf)
	}
}
func (clinet *Client) SendPacket(payload []byte) {
	if clinet.siv != 0 {
		length := uint32(len(payload))
		hash := Crypto.HashEncrypt(payload, length, clinet.siv)
		outStream := Stream.NewOutStream()
		outStream.WriteUInt(clinet.siv ^ (length + 4) ^ 0xA05F33BA)
		outStream.WriteBytes(payload)
		outStream.WriteUInt(clinet.siv ^ hash ^ 0x75AF23CB)
		clinet.Conn.Write(outStream.ToBytes())
		clinet.siv += 0x1473F19
	} else {
		length := make([]byte, 4)
		binary.LittleEndian.PutUint32(length, uint32(len(payload)))
		clinet.Conn.Write(append(length[:], payload[:]...))
	}
}
func HandleClient(client *Client) {
	Logger.Infof("New Connection! %v", client.Conn.RemoteAddr())
	client.serverClient = Server.NewServerClient(client)
	go client.RecvPacket()
}
func (client *Client) RecvCallback(payload []byte) {
	Logger.Infof("S2C : % X", payload)

	inStream := Stream.NewInStream(payload)
	rttiValue := inStream.ReadUInt()
	switch rttiValue {
	//PcFirstAccept
	case 0x221a050b:
		client.SendPacket(payload)
		xorOne := inStream.ReadInt64()
		xorTwo := inStream.ReadInt64()
		shift := inStream.ReadByte()
		client.riv = CalculateXorDword(xorOne, shift) ^ CalculateXorDword(xorTwo, shift) ^ 0xA815B623
		client.siv = client.riv
	default:
		client.SendPacket(payload)
	}
}
func (client *Client) OnExit() {
	_ = client.Conn.Close()
	client.serverClient.OnExit()
}
