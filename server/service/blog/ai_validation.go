package blog

import (
	"fmt"
	"unicode/utf8"
)

func ValidateAiChatRequest(req *AiChatRequest) error { return validateAiChatRequest(req) }

// 字符限制在构建提示词之前执行，HTTP 层另有 2 MiB 请求体限制。
func ValidateAiRequestSize(req *AiChatRequest) error {
	for _, field := range []struct {
		name, value string
		limit       int
	}{
		{"正文", req.Content, 200000}, {"选区", req.Selection, 200000},
		{"光标上下文", req.CursorContext, 200000}, {"标题", req.Title, 500},
		{"指令", req.Instruction, 8000},
	} {
		if utf8.RuneCountInString(field.value) > field.limit {
			return fmt.Errorf("%s超过 %d 字限制", field.name, field.limit)
		}
	}
	if len(req.History) > 20 {
		return fmt.Errorf("对话历史过长")
	}
	for _, msg := range req.History {
		if utf8.RuneCountInString(msg.Content) > 8000 {
			return fmt.Errorf("单条对话历史超过 8000 字限制")
		}
	}
	return nil
}
