package _interface

// 分享提醒
const (
	ShareHintCreateAppThreshold = 1 // APP 创建数触发阈值
	ShareHintUseAppThreshold    = 3 // APP 运行数触发阈值
)

// 积分策略
const (
	PointsStrategyUsingAppMaxCost = 10
	PointsStrategyUsingAppMinCost = 5
	PointsStrategyCreateAppEarn   = 5
	PointsStrategyFirstLogin      = 30
)

// 风控限制
const (
	LimitMaxAppUsedEarnPointsEveryDay    = 1000 // 运行 APP每日最大积分收益
	LimitMaxAppCreatedEarnPointsEveryday = 50   // 创建 APP每日最大积分收益
)
