package blog

import (
	"strings"
	"testing"
)

func TestValidateAiChatRequest(t *testing.T) {
	cases := []struct {
		name    string
		req     AiChatRequest
		wantErr bool
	}{
		{"polish缺selection", AiChatRequest{Action: "polish"}, true},
		{"polish正常", AiChatRequest{Action: "polish", Selection: "某段文字"}, false},
		{"rewrite缺selection", AiChatRequest{Action: "rewrite"}, true},
		{"continue缺上下文", AiChatRequest{Action: "continue"}, true},
		{"continue正常", AiChatRequest{Action: "continue", Content: "正文"}, false},
		{"outline缺全部", AiChatRequest{Action: "outline"}, true},
		{"outline有标题", AiChatRequest{Action: "outline", Title: "标题"}, false},
		{"custom缺指令", AiChatRequest{Action: "custom"}, true},
		{"custom正常", AiChatRequest{Action: "custom", Instruction: "改写为问答体"}, false},
		{"非法action", AiChatRequest{Action: "hack"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateAiChatRequest(&tc.req)
			if (err != nil) != tc.wantErr {
				t.Fatalf("wantErr=%v, got err=%v", tc.wantErr, err)
			}
		})
	}
}

func TestBuildUserMessageTruncation(t *testing.T) {
	svc := &AiService{}
	long := strings.Repeat("字", 12000)
	req := &AiChatRequest{Action: aiActionCustom, Instruction: "总结", Content: long}
	msg := svc.buildUserMessage(req)
	// 截断后不应包含全部 12000 字（默认 context-limit 8000，从尾部截断）
	if strings.Count(msg, "字") >= 12000 {
		t.Fatalf("正文未按 context-limit 截断")
	}
	if !strings.Contains(msg, "总结") {
		t.Fatalf("自定义指令丢失")
	}
}

func TestTrimHistory(t *testing.T) {
	svc := &AiService{}
	history := make([]AiChatMessage, 0, 20)
	for i := 0; i < 10; i++ {
		history = append(history,
			AiChatMessage{Role: "user", Content: "u"},
			AiChatMessage{Role: "assistant", Content: "a"})
	}
	msgs := svc.trimHistory(history)
	if len(msgs) != aiMaxHistoryTurns {
		t.Fatalf("期望保留 %d 条，实际 %d 条", aiMaxHistoryTurns, len(msgs))
	}
	// 非法角色应被过滤
	filtered := svc.trimHistory([]AiChatMessage{{Role: "system", Content: "x"}, {Role: "user", Content: "y"}})
	if len(filtered) != 1 {
		t.Fatalf("期望过滤后 1 条，实际 %d 条", len(filtered))
	}
}
