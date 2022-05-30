package code

//跟follow相关的错误码响应

type ResCode int64

const (
	UserFollowed ResCode = 9000 + iota
	UserUnfollowed
	ServeBusy
)

var FollowCodeMap = map[ResCode]string{
	UserFollowed:   "用户已关注",
	UserUnfollowed: "用户未关注",
	ServeBusy:      "服务繁忙",
}

func (c ResCode) Msg() string {
	msg, ok := FollowCodeMap[c]
	if !ok {
		msg = FollowCodeMap[ServeBusy]
	}
	return msg
}
