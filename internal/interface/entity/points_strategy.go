package entity

type IStrategyArg interface {
	Type() string
}

// StrategyArgSignInDaily 每日签到
type StrategyArgSignInDaily struct {
}

func (s StrategyArgSignInDaily) Type() string {
	return PointsTypeSignIn
}

// StrategyArgInvite 邀请好友
type StrategyArgInvite struct {
}

func (s StrategyArgInvite) Type() string {
	return PointsTypeInvite
}

// StrategyArgBeInvited 被邀请
type StrategyArgBeInvited struct {
}

func (s StrategyArgBeInvited) Type() string {
	return PointsTypeBeInvited
}

// StrategyArgNewRegister 新用户注册
type StrategyArgNewRegister struct{}

func (s StrategyArgNewRegister) Type() string {
	return PointsTypeNewRegister
}

// StrategyArgRecharge 积分充值
type StrategyArgRecharge struct {
	Points int
}

func (s StrategyArgRecharge) Type() string {
	return PointsTypePointsRecharge
}

// StrategyArgDuplicatingApp 使用程序
type StrategyArgDuplicatingApp struct {
}

func (s StrategyArgDuplicatingApp) Type() string {
	return PointsTypeDuplicatingApp
}

// StrategyArgAppDuplicated 使用程序
type StrategyArgAppDuplicated struct {
}

func (s StrategyArgAppDuplicated) Type() string {
	return PointsTypeAppDuplicated
}

// StrategyArgUsingApp 使用程序
type StrategyArgUsingApp struct {
	Form MiniAppFormFields
}

func (s StrategyArgUsingApp) Type() string {
	return PointsTypeUsingApp
}

// StrategyArgAppCreated 创建程序
type StrategyArgAppCreated struct {
}

func (s StrategyArgAppCreated) Type() string {
	return PointsTypeAppCreated
}

// StrategyArgAppUsed 使用程序
type StrategyArgAppUsed struct {
	App       MiniApp
	RunUserId int64
}

func (s StrategyArgAppUsed) Type() string {
	return PointsTypeAppUsed
}

// StrategyArgGptConversation GPT 对话
type StrategyArgGptConversation struct {
}

func (s StrategyArgGptConversation) Type() string {
	return PointsTypeGptConversation
}
