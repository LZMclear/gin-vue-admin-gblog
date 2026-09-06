package ai

// Factory 返回包级唯一的模型工厂实例。
// 数据库配置变更时通过 Invalidate 失效缓存，模型实例随之重建。
func Factory() *ModelFactory { return defaultFactory }

var defaultFactory = &ModelFactory{}

type ServiceGroup struct {
	ModelConfigService ModelConfigService
}
