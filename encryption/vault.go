package encryption

import "fmt"

type Entry struct {
	Service  string `json:"service"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Vault struct {
	Entries []Entry `json:"entries"`
}

type VaultFile struct {
	Salt       string `json:"salt"`
	Nonce      string `json:"nonce"`
	CipherText string `json:"ciphertext"`
}

func NewVault() Vault {
	return Vault{
		Entries: []Entry{},
	}
}

func (v *Vault) AddEntry(e Entry) {
	v.Entries = append(v.Entries, e)
}

func (v *Vault) DeleteEntry(service string) error {
	newEntries := make([]Entry, 0, len(v.Entries))
	for _, e := range v.Entries {
		if e.Service != service {
			newEntries = append(newEntries, e)
		}
	}
	if len(v.Entries) == len(newEntries) {
		return fmt.Errorf("\n%s not found", service)
	}
	v.Entries = newEntries
	return nil
}
func (v *Vault) GetEntry(service string) ([]Entry, error) {
	var entries []Entry
	for _, e := range v.Entries {
		if e.Service == service {
			entries = append(entries, e)
		}
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("No entry found for %s", service)
	}
	return entries, nil
}
