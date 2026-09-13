package request

type AiChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AiChatRequest struct {
	Action        string          `json:"action"`
	Content       string          `json:"content"`
	Selection     string          `json:"selection"`
	CursorContext string          `json:"cursorContext"`
	CursorOffset  *int            `json:"cursorOffset,omitempty"`
	Title         string          `json:"title"`
	Instruction   string          `json:"instruction"`
	History       []AiChatMessage `json:"history"`
}
