package DispatchEvents

const (
	INTERACTION_CREATE string = "INTERACTION_CREATE"
	READY                     = "READY"
	RESUMED                   = "RESUMED"
)

type DispatchEvent struct {
}

func GetDispatchData(t string, data interface{}) interface{} {
	switch t {
	case READY:
		return new(ReadyEvent)
	case INTERACTION_CREATE:
		var result = new(InteractionCreateEvent)
		result.Data = new(InteractionDataApplicationCommand)
		return result
	}
	return nil
}
