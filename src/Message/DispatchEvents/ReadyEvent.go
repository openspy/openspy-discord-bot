package DispatchEvents

type ReadyEvent struct {
	APIVersion       int    `json:"v"`
	ResumeGatewayUrl string `json:"resume_gateway_url"`
	SessionId        string `json:"session_id"`
}
