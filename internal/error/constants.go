package businesserrors

import "github.com/gone-io/gone/goner/gin"

// not found
var (
	ErrorChannelNotFound  = gin.NewBusinessError("channel not found", 404001)
	ErrorApprovalNotFound = gin.NewBusinessError("approval not found", 404002)
	ErrorNoticeNotFound   = gin.NewBusinessError("notice not found", 404003)
	ErrorPermissionDenied = gin.NewBusinessError("permission denied", 403000)
)

var (
	ErrorPointsNotEnough         = gin.NewBusinessError("points not enough", 500000)
	ErrorPointsNotAbleToWithdraw = gin.NewBusinessError("points not able to withdraw", 500100)
)

var (
	ErrorInvalidLoginUser = gin.NewBusinessError("invalid login user", 401001)
)

var (
	ErrorPointsOrderNewDealDenied = gin.NewBusinessError("you are not a new user", 403010)
)

var (
	ErrorAlreadySignIn = gin.NewBusinessError("you have already sign in", 403020)
)
