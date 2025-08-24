package entity



// Sign 表示能指/所指实体
type Sign struct {
	ID        int64  `json:"id"`
	Signifier string `json:"signifier"` // 能指
	Signified string `json:"signified"` // 所指
}

type Onton struct {
	ID   int64 `json:"id"`
	Adic int8  `json:"adic"`
}

