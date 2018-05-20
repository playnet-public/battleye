package battleye

// Protocol offers an interface representation for BattlEye communications
//go:generate counterfeiter -o ../mocks/protocol.go --fake-name Protocol . Protocol
type Protocol interface {
	//BuildPacket creates a new packet with data and type
	BuildPacket([]byte, Type) Packet

	//BuildLoginPacket creates a login packet with password
	BuildLoginPacket(string) Packet

	//BuildCmdPacket creates a packet with cmd and seq
	BuildCmdPacket([]byte, Sequence) Packet

	//BuildKeepAlivePacket creates a keepAlivePacket with seq
	BuildKeepAlivePacket(Sequence) Packet

	//BuildMsgAckPacket creates a server message packet with seq
	BuildMsgAckPacket(Sequence) Packet

	// Verify if the packet is valid
	Verify(Packet) error

	// Sequence extracts the seq number from a packet
	Sequence(Packet) (Sequence, error)

	// Type determines the kind of response from a packet
	Type(Packet) (Type, error)

	// Data returns the actual data inside the packet
	Data(Packet) ([]byte, error)

	// VerifyLogin returns nil on successful login
	// and a respective error on failed login
	VerifyLogin(d Packet) error

	// Multi checks whether a packet is part of a multiPacketResponse
	// Returns: packetCount, currentPacket and isSingle
	Multi(Packet) (byte, byte, bool)
}

type protocol struct{}

// New provides a real and tested implementation of Protocol
func New() Protocol {
	return &protocol{}
}

// Sequence is used by BattlEye to keep track of transactions
type Sequence uint32

//BuildPacket creates a new packet with data and type
func (p *protocol) BuildPacket(data []byte, t Type) Packet {
	data = append([]byte{0xFF, byte(t)}, data...)
	checksum := makeChecksum(data)
	header := buildHeader(checksum)

	return append(header, data...)
}

//BuildLoginPacket creates a login packet with password
func (p *protocol) BuildLoginPacket(pw string) Packet {
	return p.BuildPacket([]byte(pw), Login)
}

//BuildCmdPacket creates a packet with cmd and seq
func (p *protocol) BuildCmdPacket(cmd []byte, seq Sequence) Packet {
	return p.BuildPacket(append([]byte{byte(seq)}, cmd...), Command)
}

//BuildKeepAlivePacket creates a keepAlivePacket with seq
func (p *protocol) BuildKeepAlivePacket(seq Sequence) Packet {
	return p.BuildPacket([]byte{byte(seq)}, Command)
}

//BuildMsgAckPacket creates a server message packet with seq
func (p *protocol) BuildMsgAckPacket(seq Sequence) Packet {
	return p.BuildPacket([]byte{byte(seq)}, ServerMessage)
}
