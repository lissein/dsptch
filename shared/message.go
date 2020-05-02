package shared

type SourceMessage struct {
	Source  string
	Content map[string]interface{}
}

type DestinationMessage struct {
	Source  string
	Content map[string]interface{}
}
