package common

type RawMessage struct {
	Msg string `json:"msg"`
	ID  int    `json:"id"`
}

type ServiceInstruction struct {
	Produce bool `json:"produce"`
	Consume bool `json:"consume"`
}
