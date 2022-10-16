package main

type DnsQuestion struct {
	Name  string
	QType QueryType
}

func newDnsQuestion(name string, QType QueryType) DnsQuestion {
	return DnsQuestion{
		Name:  name,
		QType: QType,
	}
}

func (q *DnsQuestion) read(buffer *BytePacketBuffer) {
	buffer.readQname(&q.Name)
	q.QType = fromQueryTypeNum(buffer.readUint16())
	_ = buffer.readUint16()
}

func (q *DnsQuestion) write(buffer *BytePacketBuffer) {
	buffer.writeQname(&q.Name)

	typeNum := q.QType.toNum()
	buffer.writeUint16(typeNum)
	buffer.writeUint16(1)
}
