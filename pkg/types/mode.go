package types

type ConnMode int

const (
	DIRECT ConnMode = iota
	P2P
)

var modeName = map[ConnMode]string{
	DIRECT: "direct",
	P2P: "p2p",
}

func (cm ConnMode) String() string {
    return modeName[cm]
}
