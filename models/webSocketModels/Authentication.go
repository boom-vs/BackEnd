package webSocketModels

type AuthenticationLoginRequest struct {
	Login    string `json:"Login"`
	Password string `json:"Password"`
	Keep     bool   `json:"Keep"`
}

type AuthenticationLoginResponce struct {
	Token string
}
