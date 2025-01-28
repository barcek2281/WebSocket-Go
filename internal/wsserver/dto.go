package wsserver

type wsMessage struct {
	IpAddress string `json:"address"`
	Message string `json:"message"`
	Time string `json:"time"`
}