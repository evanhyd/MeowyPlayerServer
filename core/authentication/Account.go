package authentication

type account struct {
	ID   string `json:"id"`
	Salt []byte `json:"salt"`
	Hash []byte `json:"hash"`
}
