package socket

import (
	"fmt"
	"pass/encryption"
	"strconv"
)

func Gen(service, username string, sizeStr string, v *encryption.Vault) (encryption.Entry, error) {

	size, err := strconv.Atoi(sizeStr)
	var entry encryption.Entry
	if err != nil {
		return entry, fmt.Errorf("Password length must be a number.")
	}
	if size > 33 {
		return entry, fmt.Errorf("Password length must be <=32.")
	}

	password, err := encryption.GenerateRandomPassword(size)
	if err != nil {
		return entry, err
	}
	entry = encryption.Entry{
		Service:  service,
		Username: username,
		Password: password,
	}

	v.AddEntry(entry)
	return entry, nil
}
