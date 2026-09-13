package blog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	modelService "github.com/flipped-aurora/gin-vue-admin/server/service/ai"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestChatValidatesBeforeQuota(t *testing.T) {
	oldLimit, oldRedis := global.GVA_CONFIG.AI.DailyLimit, global.GVA_REDIS
	defer func() { global.GVA_CONFIG.AI.DailyLimit, global.GVA_REDIS = oldLimit, oldRedis }()
	global.GVA_CONFIG.AI.DailyLimit, global.GVA_REDIS = 5, nil
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/blog/ai/chat", strings.NewReader(`{"action":"invalid"}`))
	(&AiApi{}).Chat(ctx)
	if !strings.Contains(recorder.Body.String(), "不支持的 action") {
		t.Fatal(recorder.Body.String())
	}
}

// 真实模型适配器 + Agent + API，供应商为本地模拟 HTTP 服务，不调用外部模型。
func TestChatCompletionReason(t *testing.T) {
	oldConfig, oldDB, oldLog := global.GVA_CONFIG.AI, global.GVA_DB, global.GVA_LOG
	defer func() {
		global.GVA_CONFIG.AI, global.GVA_DB, global.GVA_LOG = oldConfig, oldDB, oldLog
		modelService.Factory().Invalidate()
	}()
	global.GVA_DB, global.GVA_LOG = nil, zap.NewNop()
	for _, reason := range []string{"stop", "length", ""} {
		t.Run("reason="+reason, func(t *testing.T) {
			provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				fmt.Fprint(w, "data: {\"id\":\"test\",\"object\":\"chat.completion.chunk\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"测试结果\"},\"finish_reason\":null}]}\n\n")
				if reason != "" {
					data, _ := json.Marshal(map[string]any{"id": "test", "object": "chat.completion.chunk", "choices": []any{map[string]any{"index": 0, "delta": map[string]any{}, "finish_reason": reason}}})
					fmt.Fprintf(w, "data: %s\n\n", data)
				}
				fmt.Fprint(w, "data: [DONE]\n\n")
			}))
			defer provider.Close()
			global.GVA_CONFIG.AI.Enable = true
			global.GVA_CONFIG.AI.Provider = "openai"
			global.GVA_CONFIG.AI.APIKey = "local-test"
			global.GVA_CONFIG.AI.BaseURL = provider.URL
			global.GVA_CONFIG.AI.Model = "local-model"
			global.GVA_CONFIG.AI.DailyLimit = 0
			modelService.Factory().Invalidate()
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/blog/ai/chat", strings.NewReader(`{"action":"custom","instruction":"测试"}`))
			(&AiApi{}).Chat(ctx)
			expected := reason
			if expected == "" {
				expected = "unknown"
			}
			body := recorder.Body.String()
			if !strings.Contains(body, `"delta":"测试结果"`) || !strings.Contains(body, `"finishReason":"`+expected+`"`) {
				t.Fatal(body)
			}
		})
	}
}
