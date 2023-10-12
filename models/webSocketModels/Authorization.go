package webSocketModels

type AuthorizationGetRequest struct {
	Token string `json:"Token"`
}

type AuthorizationGetResponse struct {
	ID             uint
	FirstName      string
	LastName       string
	PatronymicName string
	Login          string
}
