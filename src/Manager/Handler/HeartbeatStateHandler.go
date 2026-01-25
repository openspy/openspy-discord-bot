package Handler

import (
	"os"
	"time"

	"openspy.net/discord-gateway-bot/DiscordApi"
	"openspy.net/discord-gateway-bot/IManager"
	"openspy.net/discord-gateway-bot/Message"
)

type HeartbeatStateHandler struct {
	manager      IManager.IManager
	isRunning    bool
	sentIdentity bool

	lastHeartbeat time.Time
}

func (h *HeartbeatStateHandler) SetManager(manager IManager.IManager) {
	h.manager = manager
	h.isRunning = true
	h.sentIdentity = false

	h.sendHeartbeat()
}
func (h *HeartbeatStateHandler) HandleMessage(msg Message.Event) {
	// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-structure -> RESPOND

	switch msg.Op {

	case Message.HEARTBEAT_ACK:
		if !h.sentIdentity {
			h.sendIdentity()
		}
		//record last time, check if too high and exit
	case Message.DISPATCH_EVENT:
		//handle event, etc
		h.manager.GetDispatchChannel() <- msg
	}
}
func (h *HeartbeatStateHandler) sendHeartbeat() {
	h.lastHeartbeat = time.Now()
	var msg Message.Event
	msg.Op = Message.HEARTBEAT_EVENT
	h.manager.WriteMessage(msg)
}
func (h *HeartbeatStateHandler) Tick() {
	now := time.Now()

	diff := now.Sub(h.lastHeartbeat)

	var hbDuration = time.Duration(DiscordApi.HeartbeatInterval * int(time.Millisecond))

	if diff > hbDuration {
		h.sendHeartbeat()
	}
}

func (h *HeartbeatStateHandler) sendIdentity() {
	h.sentIdentity = true

	var msg Message.Event
	msg.Op = Message.IDENTIFY_EVENT
	var data Message.IdentifyEvent
	data.Token = os.Getenv("BOT_TOKEN")
	data.Intents = 0

	var props Message.IdentifyProperties
	props.Browser = "openspy"
	props.Device = "openspy"
	props.OS = "linux"
	data.Properties = props

	var presence Message.PresenceProperties
	presence.AFK = false
	presence.Status = "OpenSpy"
	presence.Status = "online"
	data.Presence = presence

	var activity Message.PresenceActivities
	activity.Name = "OpenSpy"
	activity.Type = 0
	data.Presence.Activities = append(data.Presence.Activities, activity)

	msg.Data = data
	h.manager.WriteMessage(msg)
}
