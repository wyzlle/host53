package main

type DnsHeader struct {
	Id                  uint16
	RecursionDesired    bool
	TruncatedMessage    bool
	AuthoritativeAnswer bool
	Opcode              uint8
	Response            bool

	Rescode            ResultCode
	CheckingDisabled   bool
	AuthedData         bool
	Z                  bool
	RecursionAvailable bool

	Questions            uint16
	Answers              uint16
	AuthoritativeEntries uint16
	ResourceEntries      uint16
}

func newDnsHeader() DnsHeader {
	return DnsHeader{
		Id:                  0,
		RecursionDesired:    false,
		TruncatedMessage:    false,
		AuthoritativeAnswer: false,
		Opcode:              0,
		Response:            false,

		Rescode:            NOERROR,
		CheckingDisabled:   false,
		AuthedData:         false,
		Z:                  false,
		RecursionAvailable: false,

		Questions:            0,
		Answers:              0,
		AuthoritativeEntries: 0,
		ResourceEntries:      0,
	}
}

func (header *DnsHeader) read(buffer *BytePacketBuffer) {
	header.Id = buffer.readUint16()
	flags := buffer.readUint16()

	a := uint8(flags >> 8)
	b := uint8(flags & 0xFF)
	header.RecursionDesired = (a & (1 << 0)) > 0
	header.TruncatedMessage = (a & (1 << 1)) > 0
	header.AuthoritativeAnswer = (a & (1 << 2)) > 0
	header.Opcode = (a >> 3) & 0x0F
	header.Response = (a & (1 << 7)) > 0

	header.Rescode = fromResultCodeNum(b & 0x0F)
	header.CheckingDisabled = (b & (1 << 4)) > 0
	header.AuthedData = (b & (1 << 5)) > 0
	header.Z = (b & (1 << 6)) > 0
	header.RecursionAvailable = (b & (1 << 7)) > 0

	header.Questions = buffer.readUint16()
	header.Answers = buffer.readUint16()
	header.AuthoritativeEntries = buffer.readUint16()
	header.ResourceEntries = buffer.readUint16()
}

func (header *DnsHeader) write(buffer *BytePacketBuffer) {
	buffer.writeUint16(header.Id)

	buffer.writeUint8(
		boolToUint8(header.RecursionDesired) |
			boolToUint8(header.TruncatedMessage)<<1 |
			boolToUint8(header.AuthoritativeAnswer)<<2 |
			uint8(header.Opcode)<<3 |
			boolToUint8(header.Response)<<7,
	)

	buffer.writeUint8(
		uint8(header.Rescode) |
			boolToUint8(header.CheckingDisabled)<<4 |
			boolToUint8(header.AuthedData)<<5 |
			boolToUint8(header.Z)<<6 |
			boolToUint8(header.RecursionAvailable)<<7,
	)

	buffer.writeUint16(header.Questions)
	buffer.writeUint16(header.Answers)
	buffer.writeUint16(header.AuthoritativeEntries)
	buffer.writeUint16(header.ResourceEntries)
}

func boolToUint8(variable bool) uint8 {
	var variable8 uint8
	if variable {
		variable8 = 1
	}
	return variable8
}
