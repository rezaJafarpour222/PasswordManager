package protocol

import (
	"encoding/json"
	"net"
	"pass/encryption"
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
	Error   bool            `json:"error"`
	Payload json.RawMessage `json:"response,omitempty"`
}

var SOCKET = "/tmp/pass.sock"

func Serialize(conn net.Conn, error bool, payload any) error {
	p, err := json.Marshal(payload)
	if err != nil {
		return json.NewEncoder(conn).Encode(err)
	}
	res := Response{
		Error:   error,
		Payload: p,
	}

	return json.NewEncoder(conn).Encode(res)
}
func Deserialize(conn net.Conn) (Request, error) {
	var req Request
	err := json.NewDecoder(conn).Decode(&req)
	if err != nil {
		return req, err
	}
	return req, nil
}

func DeserializeEntries(res Response) ([]encryption.Entry, error) {
	var entries []encryption.Entry
	if err := json.Unmarshal(res.Payload, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func DeserializeEntry(res Response) (encryption.Entry, error) {
	var entry encryption.Entry
	if err := json.Unmarshal(res.Payload, &entry); err != nil {
		return entry, err
	}
	return entry, nil
}

func Dial(req Request) (Response, error) {
	conn, err := net.Dial("unix", SOCKET)
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return Response{}, err
	}

	var resp Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return Response{}, err
	}
	return resp, nil
}
