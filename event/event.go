package event

type InformerEvent struct {
	Key          string
	EventType    string
	Namespace    string
	ResourceType string
}
