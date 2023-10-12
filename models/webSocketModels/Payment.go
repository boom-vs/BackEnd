package webSocketModels

type PaymentSplit struct {
	Order         uint32  `json:"Order"`
	Employee      uint32  `json:"Employee"`
	Object        uint32  `json:"Object"`
	TypeOfPayment uint32  `json:"TypeOfPayment"`
	Credit        float64 `json:"Credit"`
}

type PaymentScans struct {
	MediaTypes uint16 `json:"MediaTypes"`
	Documents  uint16 `json:"Documents"`
	Payments   uint16 `json:"Payments"`
}

type Payment struct {
	Period      string         `json:"Period"`
	Date        string         `json:"Date"`
	Allocate    float64        `json:"Allocate"`
	Destination string         `json:"Destination"`
	Payer       string         `json:"Payer"`
	Scans       PaymentScans   `json:"Scans"`
	Comment     string         `json:"Comment"`
	Split       []PaymentSplit `json:"Split"`
}
