package socket

type Request struct {
	Action   string `json:"action"`
	Service  string `json:"service,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Length   int    `json:"length,omitempty"`
}
