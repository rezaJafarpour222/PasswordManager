package socket

import (
	"encoding/json"
	"net"
)

type Request struct {
	Action     string `json:"action"`
	Service    string `json:"service,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Length     int    `json:"length,omitempty"`
	ExportPath string `json:"exportPath,omitempty"`
}
type Response struct {
	Error   bool `json:"error"`
	Payload any  `json:"response,omitempty"`
}

func Serialize(conn net.Conn, any Response) error {
	return json.NewEncoder(conn).Encode(any)
}
func Deserialize(conn net.Conn) (Request, error) {
	var req Request
	err := json.NewDecoder(conn).Decode(&req)
	if err != nil {
		return req, err
	}
	return req, nil
}
