package Message

type ResumeEvent struct {
	Token     string `json:"token"`
	SessionId string `json:"session_id"`
	Sequence  int    `json:"seq"`
}
