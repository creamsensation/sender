package sender

type Json struct {
	Result any `json:"result"`
}

type Error struct {
	Error string `json:"error"`
}
