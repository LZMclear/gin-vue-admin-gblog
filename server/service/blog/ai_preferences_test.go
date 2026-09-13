package blog

import (
	"strings"
	"testing"
)

func TestWritingPreferenceValidation(t *testing.T) {
	for _, req := range []*AiChatRequest{{Tone: "injected"}, {Length: "invalid"}} {
		if ValidateAiRequestSize(req) == nil {
			t.Fatal("invalid preferences accepted")
		}
	}
	if err := ValidateAiRequestSize(&AiChatRequest{Tone: "formal", Length: "shorter"}); err != nil {
		t.Fatal(err)
	}
}

func TestWritingPreferencesAppliedToBodyActions(t *testing.T) {
	svc := &AiService{}
	req := &AiChatRequest{Action: "rewrite", Selection: "原选区", Tone: "formal", Length: "shorter"}
	message := svc.buildUserMessage(req)
	for _, part := range []string{"原选区", "正式", "七成"} {
		if !strings.Contains(message, part) {
			t.Fatal(message)
		}
	}
	req.Action = "summary"
	if writingPreferences(req) != "" {
		t.Fatal("body preferences leaked into structured task")
	}
	req.Action = "custom"
	req.Instruction = "保持原篇幅"
	if !strings.Contains(svc.buildUserMessage(req), "服从作者明确指令") {
		t.Fatal("instruction precedence missing")
	}
}

func TestTitlePromptSupportsTitleOnly(t *testing.T) {
	req := &AiChatRequest{Action: "title", Title: "仅有主题的文章"}
	if err := validateAiChatRequest(req); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains((&AiService{}).buildUserMessage(req), req.Title) {
		t.Fatal("title-only topic omitted")
	}
}
