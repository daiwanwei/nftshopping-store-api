package messages

type DeliverItemMessage struct {
	CreationId string `json:"creationId"`
	Contract   string `json:"contract"`
	Token      string `json:"token"`
}

type OrderItemMessage struct {
	CreationId string `json:"creationId"`
	Contract   string `json:"contract"`
	Amount     int    `json:"amount"`
}
