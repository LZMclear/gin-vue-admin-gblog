package blog

import "strings"

var aiTones = map[string]string{
	"": "", "natural": "表达自然，保留作者原有文风", "formal": "语气正式、表述严谨", "friendly": "语气亲切，表达易懂",
}
var aiLengths = map[string]string{
	"": "", "original": "保持与原文相近的篇幅", "shorter": "压缩冗余表达，篇幅约为原文的七成，保留关键信息", "longer": "适度补充解释，篇幅约为原文的一点三倍，不编造事实",
}

func writingPreferences(req *AiChatRequest) string {
	switch req.Action {
	case aiActionPolish, aiActionRewrite, aiActionContinue, aiActionCustom:
	default:
		return ""
	}
	var preferences []string
	if tone := aiTones[req.Tone]; tone != "" {
		preferences = append(preferences, tone)
	}
	if length := aiLengths[req.Length]; length != "" {
		preferences = append(preferences, length)
	}
	if len(preferences) == 0 {
		return ""
	}
	return "\n写作偏好（服从作者明确指令，续写仍不超过300字）：" + strings.Join(preferences, "；")
}
