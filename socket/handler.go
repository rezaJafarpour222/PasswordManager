package socket

import (
	"encoding/json"
	"fmt"
	"net"
	"pass/encryption"
	"pass/storage"
)

func Handle(conn net.Conn, masterKey string, vault *encryption.Vault, vaultPath string) error {
	defer conn.Close()
	var req Request
	err := json.NewDecoder(conn).Decode(&req)
	if err != nil {
		return err
	}
	switch req.Action {
	case "get":
		if req.Service == "" {
			return fmt.Errorf("service is empty")
		}
		entry, err := vault.GetEntry(req.Service)
		if err != nil {
			return err
		}
		return json.NewEncoder(conn).Encode(entry)

	case "gen":
		if req.Service == "" || req.Username == "" {
			return fmt.Errorf("service is empty")
		}
		entry, err := Gen(req.Service, req.Username, "32", vault)
		if err != nil {
			return err
		}
		storage.SaveVault(*vault, masterKey, vaultPath)
		err = json.NewEncoder(conn).Encode(entry)
		if err != nil {
			return err
		}
	default:
		return json.NewEncoder(conn).Encode("action is not supported.")
	}
	return nil
}
