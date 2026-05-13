package blog

import "testing"

func TestVisitBehaviorDefaults(t *testing.T) {
	tests := []struct {
		name     VisitBehavior
		behavior string
		content  string
		countPV  bool
	}{
		{name: VisitBehaviorIndex, behavior: "访问页面", content: "首页", countPV: true},
		{name: VisitBehaviorArchive, behavior: "访问页面", content: "归档", countPV: true},
		{name: VisitBehaviorBlog, behavior: "查看博客", countPV: true},
		{name: VisitBehaviorSearch, behavior: "搜索博客", countPV: false},
		{name: VisitBehaviorClickFriend, behavior: "点击友链", countPV: false},
		{name: VisitBehaviorCheckPassword, behavior: "校验博客密码", countPV: false},
	}

	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			got := visitBehaviorDefaults(tt.name)
			if got.Behavior != tt.behavior {
				t.Fatalf("Behavior = %q, want %q", got.Behavior, tt.behavior)
			}
			if got.Content != tt.content {
				t.Fatalf("Content = %q, want %q", got.Content, tt.content)
			}
			if got.CountPV != tt.countPV {
				t.Fatalf("CountPV = %v, want %v", got.CountPV, tt.countPV)
			}
		})
	}
}

func TestRequestValue(t *testing.T) {
	if got := requestValue([]byte(`{"nickname":"guitu"}`), "nickname"); got != "guitu" {
		t.Fatalf("requestValue json = %q, want guitu", got)
	}
	if got := requestValue([]byte(`nickname=guitu`), "nickname"); got != "guitu" {
		t.Fatalf("requestValue form = %q, want guitu", got)
	}
}
