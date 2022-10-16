package main

import (
	"fmt"
	"net"
)

func main() {
	server := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 2054,
	}
	socket, err := net.ListenUDP("udp", &server)
	if err != nil {
		panic(err)
	}
	for {
		handleQuery(*socket)
	}
}

func lookup(qname *string, qtype QueryType) DnsPacket {
	server := net.UDPAddr{
		IP:   net.ParseIP("8.8.8.8"),
		Port: 53,
	}
	addr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 42069,
	}
	socket, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}

	packet := newDnsPacket()

	packet.Header.Id = 6666
	packet.Header.Questions = 1
	packet.Header.RecursionDesired = true
	packet.Questions = append(packet.Questions, (newDnsQuestion(*qname, qtype)))

	reqBuffer := newBytePacketBuffer()
	packet.write(&reqBuffer)

	_, err = socket.WriteToUDP(reqBuffer.Buf[0:reqBuffer.Pos], &server)
	if err != nil {
		panic(err)
	}

	resBuffer := newBytePacketBuffer()
	bytes := make([]byte, 512)
	_, _, err = socket.ReadFromUDP(bytes)

	if err != nil {
		panic(err)
	}
	socket.Close()
	copy(resBuffer.Buf[0:512], bytes[0:512])
	resPacket := fromBuffer(&resBuffer)
	return resPacket
}

func handleQuery(socket net.UDPConn) {
	reqBuffer := newBytePacketBuffer()

	_, src, err := socket.ReadFromUDP(reqBuffer.Buf[:])
	if err != nil {
		panic(err)
	}

	reqPacket := fromBuffer(&reqBuffer)

	resPacket := newDnsPacket()
	resPacket.Header.Id = reqPacket.Header.Id
	resPacket.Header.RecursionDesired = true
	resPacket.Header.RecursionAvailable = true
	resPacket.Header.Response = true

	for _, question := range reqPacket.Questions {
		fmt.Println("Received query: ", question)

		result := lookup(&question.Name, question.QType)
		resPacket.Questions = append(resPacket.Questions, question)
		resPacket.Header.Rescode = result.Header.Rescode

		for _, answer := range result.Answers {
			fmt.Println("Answer: ", answer)
			resPacket.Answers = append(resPacket.Answers, answer)
		}

		for _, authority := range result.Authorities {
			fmt.Println("Authority: ", authority)
			resPacket.Authorities = append(resPacket.Authorities, authority)
		}

		for _, resource := range result.Resources {
			fmt.Println("Resource: ", resource)
			resPacket.Resources = append(resPacket.Resources, resource)
		}
	}

	if len(resPacket.Questions) == 0 {
		resPacket.Header.Rescode = FORMERR
	}

	resBuffer := newBytePacketBuffer()
	resPacket.write(&resBuffer)

	len := resBuffer.pos()
	data := resBuffer.getRange(0, len)

	_, err = socket.WriteToUDP(data, src)
	if err != nil {
		panic(err)
	}

}
