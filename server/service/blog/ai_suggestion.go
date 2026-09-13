package blog

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

func normalizeTagSuggestion(output string, categories, tags map[string]uint) (*TagSuggestion, error) {
	text := strings.TrimSpace(strings.ReplaceAll(output, "\r\n", "\n"))
	if strings.HasPrefix(text, "```json\n") && strings.HasSuffix(text, "```") {
		text = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(text, "```json\n"), "```"))
	}
	if strings.HasPrefix(text, "```\n") && strings.HasSuffix(text, "```") {
		text = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(text, "```\n"), "```"))
	}
	var raw struct {
		Category string   `json:"category"`
		Tags     []string `json:"tags"`
		NewTags  []string `json:"newTags"`
	}
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&raw); err != nil {
		return nil, fmt.Errorf("模型未返回有效的分类标签 JSON")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF || !strings.HasPrefix(text, "{") || raw.Tags == nil || raw.NewTags == nil {
		return nil, fmt.Errorf("模型推荐结果格式不完整")
	}
	if len(raw.Tags) > 100 || len(raw.NewTags) > 100 {
		return nil, fmt.Errorf("模型返回了过多标签")
	}
	result := &TagSuggestion{Tags: []string{}, NewTags: []string{}, TagIDs: []uint{}, Warnings: []string{}}
	canonical := func(items map[string]uint) map[string]string {
		keys := make([]string, 0, len(items))
		for key := range items {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		names := make(map[string]string)
		for _, key := range keys {
			normalized := strings.ToLower(strings.TrimSpace(key))
			if _, ok := names[normalized]; !ok {
				names[normalized] = key
			}
		}
		return names
	}
	categoryNames, tagNames := canonical(categories), canonical(tags)
	if name := strings.TrimSpace(raw.Category); name != "" {
		if actual, ok := categoryNames[strings.ToLower(name)]; ok {
			result.Category = actual
			result.CategoryID = categories[actual]
		} else {
			result.Warnings = append(result.Warnings, "忽略了不在现有目录中的分类")
		}
	}
	seen := make(map[string]bool)
	for _, name := range append(raw.Tags, raw.NewTags...) {
		name = strings.TrimSpace(name)
		key := strings.ToLower(name)
		if name == "" || seen[key] {
			continue
		}
		seen[key] = true
		if actual, ok := tagNames[key]; ok {
			if len(result.Tags) < 3 {
				result.Tags = append(result.Tags, actual)
				result.TagIDs = append(result.TagIDs, tags[actual])
			}
			continue
		}
		if utf8.RuneCountInString(name) > 32 || strings.ContainsAny(name, "<>\r\n") || strings.IndexFunc(name, unicode.IsControl) >= 0 {
			result.Warnings = append(result.Warnings, "忽略了名称无效的新标签")
			continue
		}
		if len(result.NewTags) < 2 {
			result.NewTags = append(result.NewTags, name)
		}
	}
	return result, nil
}
