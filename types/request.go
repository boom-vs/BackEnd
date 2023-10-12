package types

type WebSocketPackage struct {
	Controller string
	Action     string
	Serial     string
	Hash       string
	Error      string
	Data       []interface{}
}
