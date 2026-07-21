package socket

import (
	"fmt"
	"net"
	"pass/encryption"
	"pass/protocol"
	"pass/storage"
)

var CURRENTPATH = "~/"

func Handle(conn net.Conn, masterKey string, vault *encryption.Vault, vaultPath string) {
	defer conn.Close()
	var req protocol.Request
	req, err := protocol.Deserialize(conn)
	if err != nil {
		protocol.Serialize(conn, true, err)
	}
	switch req.Action {
	case "get":
		if req.Service == "" {
			protocol.Serialize(conn, true, "Missing Required Value(Service).")
		}
		entry, err := vault.GetEntry(req.Service)
		if err != nil {
			protocol.Serialize(conn, true, err)
		}
		protocol.Serialize(conn, false, entry)

	case "gen":
		if req.Service == "" || req.Username == "" {
			protocol.Serialize(conn, true, "Missing Required Value(Service).")
			return
		}
		entry, err := Gen(req.Service, req.Username, "32", vault)
		if err != nil {
			protocol.Serialize(conn, true, err)
			return
		}
		storage.SaveVault(*vault, masterKey, vaultPath)
		err = protocol.Serialize(conn, false, entry)
		if err != nil {
			protocol.Serialize(conn, true, err)
			return
		}
	case "add":
		if req.Service == "" || req.Password == "" || req.Username == "" {
			protocol.Serialize(conn, true, "Missing Required Value(Service or Username or Password).")
		}
		Add(req.Service, req.Username, req.Password, vault)

		storage.SaveVault(*vault, masterKey, vaultPath)
		response := fmt.Sprintf(req.Service, " ", req.Username, " ", req.Password, " added to vault.")
		err = protocol.Serialize(conn, false, response)
		if err != nil {
			protocol.Serialize(conn, true, err)
		}
	case "init":
		err = Init(CURRENTPATH)
		if err != nil {
			protocol.Serialize(conn, true, err)
		}
	case "del":
		if req.Service == "" {
			protocol.Serialize(conn, true, "Missing Required Values(Service).")
		}
		err = DeleteByService(req.Service, vault)
		if err != nil {
			protocol.Serialize(conn, true, err)
		}
	case "list":
		entries := List(vault)
		protocol.Serialize(conn, false, entries)
	case "export":
		if req.ExportPath == "" {
			protocol.Serialize(conn, true, "Missing Required Values(ExportPath).")
		}
		Export(req.ExportPath, CURRENTPATH)
	default:
		protocol.Serialize(conn, true, "Action is not supported.")
	}
}
