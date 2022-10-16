package main

type DnsPacket struct {
	Header      DnsHeader
	Questions   []DnsQuestion
	Answers     []interface{}
	Authorities []interface{}
	Resources   []interface{}
}

func newDnsPacket() DnsPacket {
	return DnsPacket{
		Header:      newDnsHeader(),
		Questions:   make([]DnsQuestion, 0),
		Answers:     make([]interface{}, 0),
		Authorities: make([]interface{}, 0),
		Resources:   make([]interface{}, 0),
	}
}

func fromBuffer(buffer *BytePacketBuffer) DnsPacket {
	result := newDnsPacket()
	result.Header.read(buffer)

	for i := uint16(0); i < result.Header.Questions; i++ {
		question := newDnsQuestion("", UNKNOWN)
		question.read(buffer)
		result.Questions = append(result.Questions, question)
	}

	for i := uint16(0); i < result.Header.Answers; i++ {
		rec := read(buffer)
		result.Answers = append(result.Answers, rec)
	}

	for i := uint16(0); i < result.Header.AuthoritativeEntries; i++ {
		rec := read(buffer)
		result.Authorities = append(result.Authorities, rec)
	}

	for i := uint16(0); i < result.Header.ResourceEntries; i++ {
		rec := read(buffer)
		result.Resources = append(result.Resources, rec)
	}

	return result
}

func (packet *DnsPacket) write(buffer *BytePacketBuffer) {
	packet.Header.Questions = uint16(len(packet.Questions))
	packet.Header.Answers = uint16(len(packet.Answers))
	packet.Header.AuthoritativeEntries = uint16(len(packet.Authorities))
	packet.Header.ResourceEntries = uint16(len(packet.Resources))

	packet.Header.write(buffer)
	for _, question := range packet.Questions {
		question.write(buffer)
	}

	for _, answer := range packet.Answers {
		if answer, ok := answer.(ARecord); ok {
			answer.write(buffer)
		}

		if answer, ok := answer.(CNAMERecord); ok {
			answer.write(buffer)
		}

		if answer, ok := answer.(NSRecord); ok {
			answer.write(buffer)
		}

		if answer, ok := answer.(AAAARecord); ok {
			answer.write(buffer)
		}

		if answer, ok := answer.(MXRecord); ok {
			answer.write(buffer)
		}

		if answer, ok := answer.(UnknownRecord); ok {
			answer.write(buffer)
		}
	}

	for _, authority := range packet.Authorities {
		if answer, ok := authority.(ARecord); ok {
			answer.write(buffer)
		}

		if answer, ok := authority.(CNAMERecord); ok {
			answer.write(buffer)
		}

		if answer, ok := authority.(NSRecord); ok {
			answer.write(buffer)
		}

		if answer, ok := authority.(AAAARecord); ok {
			answer.write(buffer)
		}

		if answer, ok := authority.(MXRecord); ok {
			answer.write(buffer)
		}

		if answer, ok := authority.(UnknownRecord); ok {
			answer.write(buffer)
		}
	}

	for _, resource := range packet.Resources {
		if answer, ok := resource.(ARecord); ok {
			answer.write(buffer)
		}

		if answer, ok := resource.(CNAMERecord); ok {
			answer.write(buffer)
		}

		if answer, ok := resource.(NSRecord); ok {
			answer.write(buffer)
		}

		if answer, ok := resource.(AAAARecord); ok {
			answer.write(buffer)
		}

		if answer, ok := resource.(MXRecord); ok {
			answer.write(buffer)
		}

		if answer, ok := resource.(UnknownRecord); ok {
			answer.write(buffer)
		}
	}
}
