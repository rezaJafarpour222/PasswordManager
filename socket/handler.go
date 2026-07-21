package socket

import (
	"fmt"
	"net"
	"pass/encryption"
	"pass/storage"
)

var CURRENTPATH = "~/"

func Handle(conn net.Conn, masterKey string, vault *encryption.Vault, vaultPath string) {
	defer conn.Close()
	var req Request
	req, err := Deserialize(conn)
	if err != nil {
		Serialize(conn, Response{Error: true, Payload: err})
	}
	switch req.Action {
	case "get":
		if req.Service == "" {
			Serialize(conn, Response{Error: true, Payload: "Missing Required Value(Service)."})
		}
		entry, err := vault.GetEntry(req.Service)
		if err != nil {
			Serialize(conn, Response{Error: true, Payload: err})
		}
		Serialize(conn, Response{Error: false, Payload: entry})

	case "gen":
		if req.Service == "" || req.Username == "" {
			Serialize(conn, Response{Error: true, Payload: "Missing Required Value(Service)."})
		}
		entry, err := Gen(req.Service, req.Username, "32", vault)
		if err != nil {
			Serialize(conn, Response{Error: true, Payload: err})
		}
		storage.SaveVault(*vault, masterKey, vaultPath)
		err = Serialize(conn, Response{Error: false, Payload: entry})
		if err != nil {
			Serialize(conn, Response{Error: true, Payload: err})
		}
	case "add":
		if req.Service == "" || req.Password == "" || req.Username == "" {
			Serialize(conn, Response{Error: true, Payload: "Missing Required Value(Service or Username or Password)."})
		}
		Add(req.Service, req.Username, req.Password, vault)

		storage.SaveVault(*vault, masterKey, vaultPath)
		response := fmt.Sprintf(req.Service, " ", req.Username, " ", req.Password, " added to vault.")
		err = Serialize(conn, Response{Error: false, Payload: response})
		if err != nil {
			Serialize(conn, Response{Error: true, Payload: err})
		}
	case "init":
		err = Init(CURRENTPATH)
		if err != nil {
			Serialize(conn, Response{Error: true, Payload: err})
		}
	case "del":
		if req.Service == "" {
			Serialize(conn, Response{Error: true, Payload: "Missing Required Values(Service)."})
		}
		err = DeleteByService(req.Service, vault)
		if err != nil {
			Serialize(conn, Response{Error: true, Payload: err})
		}
	case "list":
		entries := List(vault)
		Serialize(conn, Response{Error: false, Payload: entries})
	case "export":
		if req.ExportPath == "" {
			Serialize(conn, Response{Error: true, Payload: "Missing Required Values(ExportPath)."})
		}
		Export(req.ExportPath, CURRENTPATH)
	default:
		Serialize(conn, Response{Error: true, Payload: "Action is not supported."})
	}
}
