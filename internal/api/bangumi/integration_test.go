package bangumi

import (
	"context"
	"os"
	"testing"
)

// TestIntegrationMe 真实网络集成用例：校验 token 与连通性。
// CI / 本地默认跳过，须显式设置：
//
//	BANGUMI_INTEGRATION=1 BANGUMI_TOKEN=<个人令牌> go test ./internal/api/bangumi/... -run TestIntegrationMe -v
//
// 注意：api.bgm.tv 在中国大陆可能需要代理（用户环境走 Clash，见 M1 摸底报告）。
func TestIntegrationMe(t *testing.T) {
	if os.Getenv("BANGUMI_INTEGRATION") == "" {
		t.Skip("跳过真实网络集成测试：设置 BANGUMI_INTEGRATION=1 且 BANGUMI_TOKEN=<token> 后启用")
	}
	token := os.Getenv("BANGUMI_TOKEN")
	if token == "" {
		t.Fatal("BANGUMI_INTEGRATION 已开启但缺少 BANGUMI_TOKEN")
	}

	c := New(token)
	defer c.Close()

	user, err := c.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GET /v0/me 失败: %v", err)
	}
	t.Logf("当前用户: id=%d username=%s nickname=%s", user.ID, user.Username, user.Nickname)
}
