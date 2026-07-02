package encryption

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

func (v *Vault) DeleteEntry(service string) {
	newEntries := make([]Entry, 0, len(v.Entries))

	for _, e := range v.Entries {
		if e.Service != service {
			newEntries = append(newEntries, e)
		}
	}

	v.Entries = newEntries
}
