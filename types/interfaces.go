package types

type ControllerInterface interface {
	GetMethods() map[string]func(context *WebSocketContext)
}
