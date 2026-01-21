package DiscordApi

type CreateCommandOption struct {
	Name        string `json:"name"`
	Type        int    `json:"type"`
	Description string `json:"description"`
	MinLength   int    `json:"min_length"`
	MaxLength   int    `json:"max_length"`
	Required    bool   `json:"required"`
}
