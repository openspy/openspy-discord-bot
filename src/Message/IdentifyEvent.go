package Message

const (
	ACTIVITY_TYPE_PLAYING   int = 0
	ACTIVITY_TYPE_STREAMING     = 1
	ACTIVITY_TYPE_LISTENING     = 2
	ACTIVITY_TYPE_WATCHING      = 3
	ACTIVITY_TYPE_CUSTOM        = 4
	ACTIVITY_TYPE_COMPETING     = 5
)

type PresenceActivities struct {
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Url   string `json:"url"`
	State string `json:"state"`
}
type PresenceProperties struct {
	Status     string               `json:"status"`
	AFK        bool                 `json:"afk"`
	Since      *int                 `json:"since"`
	Activities []PresenceActivities `json:"activities"`
}
type IdentifyProperties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}
type IdentifyEvent struct {
	Token      string             `json:"token"`
	Properties IdentifyProperties `json:"properties"`
	Presence   PresenceProperties `json:"presence"`
	Intents    int                `json:"intents"`
}
