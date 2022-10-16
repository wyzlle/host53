package main

type QueryType uint16

const (
	UNKNOWN QueryType = 0
	A       QueryType = 1
	NS      QueryType = 2
	CNAME   QueryType = 5
	MX      QueryType = 15
	AAAA    QueryType = 28
)

func (q QueryType) toNum() uint16 {
	switch q {
	case A:
		return 1
	case NS:
		return 2
	case CNAME:
		return 5
	case MX:
		return 15
	case AAAA:
		return 28
	default:
		return 0
	}
}

func fromQueryTypeNum(num uint16) QueryType {
	switch num {
	case 1:
		return A
	case 2:
		return NS
	case 5:
		return CNAME
	case 15:
		return MX
	case 28:
		return AAAA
	default:
		return UNKNOWN
	}
}
