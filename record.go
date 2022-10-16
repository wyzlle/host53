package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type UnknownRecord struct {
	Domain  string
	QType   uint16
	DataLen uint16
	TTL     uint32
}

type ARecord struct {
	Domain string
	Addr   net.IP
	TTL    uint32
}

type NSRecord struct {
	Domain string
	Host   string
	TTL    uint32
}

type CNAMERecord struct {
	Domain string
	Host   string
	TTL    uint32
}

type MXRecord struct {
	Domain   string
	Priority uint16
	Host     string
	TTL      uint32
}

type AAAARecord struct {
	Domain string
	Addr   net.IP
	TTL    uint32
}

func read(buffer *BytePacketBuffer) interface{} {
	domain := ""
	buffer.readQname(&domain)

	qtype_num := buffer.readUint16()
	qtype := fromQueryTypeNum(qtype_num)
	_ = buffer.readUint16()
	ttl := buffer.readUint32()
	dataLen := buffer.readUint16()

	switch qtype {
	case A:
		raw_addr := buffer.readUint32()
		addr := net.IPAddr{
			IP: net.IPv4(
				uint8((raw_addr>>24)&0xFF),
				uint8((raw_addr>>16)&0xFF),
				uint8((raw_addr>>8)&0xFF),
				uint8((raw_addr>>0)&0xFF),
			),
		}
		return ARecord{
			Domain: domain,
			Addr:   addr.IP,
			TTL:    ttl,
		}
	case AAAA:
		raw_addr1 := buffer.readUint32()
		raw_addr2 := buffer.readUint32()
		raw_addr3 := buffer.readUint32()
		raw_addr4 := buffer.readUint32()
		addr := IPv6(
			uint16((raw_addr1>>16)&0xFFFF),
			uint16((raw_addr1>>0)&0xFFFF),
			uint16((raw_addr2>>16)&0xFFFF),
			uint16((raw_addr2>>0)&0xFFFF),
			uint16((raw_addr3>>16)&0xFFFF),
			uint16((raw_addr3>>0)&0xFFFF),
			uint16((raw_addr4>>16)&0xFFFF),
			uint16((raw_addr4>>0)&0xFFFF),
		)

		return AAAARecord{
			Domain: domain,
			Addr:   addr,
			TTL:    ttl,
		}
	case NS:
		ns := ""
		buffer.readQname(&ns)
		return NSRecord{
			Domain: domain,
			Host:   ns,
			TTL:    ttl,
		}
	case CNAME:
		cname := ""
		buffer.readQname(&cname)
		return CNAMERecord{
			Domain: domain,
			Host:   cname,
			TTL:    ttl,
		}
	case MX:
		priority := buffer.readUint16()
		mx := ""
		buffer.readQname(&mx)
		return MXRecord{
			Domain:   domain,
			Priority: priority,
			Host:     mx,
			TTL:      ttl,
		}
	default:
		buffer.step(uintptr(dataLen))
		return UnknownRecord{
			Domain:  domain,
			QType:   qtype_num,
			DataLen: dataLen,
			TTL:     ttl,
		}
	}
}

func (rec *ARecord) write(buffer *BytePacketBuffer) uintptr {
	startPos := buffer.pos()
	buffer.writeQname(&rec.Domain)
	buffer.writeUint16(A.toNum())
	buffer.writeUint16(1)
	buffer.writeUint32(rec.TTL)
	buffer.writeUint16(4)

	octets := rec.Addr.To4()
	buffer.writeUint8(octets[0])
	buffer.writeUint8(octets[1])
	buffer.writeUint8(octets[2])
	buffer.writeUint8(octets[3])

	return buffer.pos() - startPos
}

func (rec *NSRecord) write(buffer *BytePacketBuffer) uintptr {
	startPos := buffer.pos()
	buffer.writeQname(&rec.Domain)
	buffer.writeUint16(NS.toNum())
	buffer.writeUint16(1)
	buffer.writeUint32(rec.TTL)

	pos := buffer.pos()
	buffer.writeUint16(0)
	buffer.writeQname(&rec.Host)

	size := buffer.pos() - (pos + 2)
	buffer.setUint16(pos, uint16(size))

	return buffer.pos() - startPos
}

func (rec *CNAMERecord) write(buffer *BytePacketBuffer) uintptr {
	startPos := buffer.pos()

	buffer.writeQname(&rec.Domain)
	buffer.writeUint16(CNAME.toNum())
	buffer.writeUint16(1)
	buffer.writeUint32(rec.TTL)

	pos := buffer.pos()
	buffer.writeUint16(0)
	buffer.writeQname(&rec.Host)

	size := buffer.pos() - (pos + 2)
	buffer.setUint16(pos, uint16(size))

	return buffer.pos() - startPos
}

func (rec *MXRecord) write(buffer *BytePacketBuffer) uintptr {
	startPos := buffer.pos()

	buffer.writeQname(&rec.Domain)
	buffer.writeUint16(MX.toNum())
	buffer.writeUint16(1)
	buffer.writeUint32(rec.TTL)

	pos := buffer.pos()
	buffer.writeUint16(0)

	buffer.writeUint16(rec.Priority)
	buffer.writeQname(&rec.Host)

	size := buffer.pos() - (pos + 2)
	buffer.setUint16(pos, uint16(size))

	return buffer.pos() - startPos
}

func (rec *AAAARecord) write(buffer *BytePacketBuffer) uintptr {
	startPos := buffer.pos()

	buffer.writeQname(&rec.Domain)
	buffer.writeUint16(AAAA.toNum())
	buffer.writeUint16(1)
	buffer.writeUint32(rec.TTL)
	buffer.writeUint16(16)

	for _, octet := range rec.Addr.To16() {
		buffer.writeUint8(octet)
	}

	return buffer.pos() - startPos
}

func (rec *UnknownRecord) write(buffer *BytePacketBuffer) uintptr {
	startPos := buffer.pos()
	fmt.Println("skipping unkown record")

	return buffer.pos() - startPos
}

func IPv6(a, b, c, d, e, f, g, h uint16) net.IP {
	ip := make(net.IP, 16)
	binary.BigEndian.PutUint16(ip[0:2], a)
	binary.BigEndian.PutUint16(ip[2:4], b)
	binary.BigEndian.PutUint16(ip[4:6], c)
	binary.BigEndian.PutUint16(ip[6:8], d)
	binary.BigEndian.PutUint16(ip[8:10], e)
	binary.BigEndian.PutUint16(ip[10:12], f)
	binary.BigEndian.PutUint16(ip[12:14], g)
	binary.BigEndian.PutUint16(ip[14:16], h)
	return ip
}
