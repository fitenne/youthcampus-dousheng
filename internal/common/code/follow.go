package code

//跟follow相关的错误码响应

type ResCode int64

const (
	Success      ResCode = 0
	UserFollowed ResCode = 9000 + iota
	UserUnfollowed
	ServeBusy
	InvalidParameter
	RepeatFollow
)

var FollowCodeMap = map[ResCode]string{
	Success:          "success",
	UserFollowed:     "用户已关注",
	UserUnfollowed:   "用户未关注",
	ServeBusy:        "服务繁忙",
	InvalidParameter: "参数错误",
	RepeatFollow:     "不能关注自己",
}

func (c ResCode) Msg() string {
	msg, ok := FollowCodeMap[c]
	if !ok {
		msg = FollowCodeMap[ServeBusy]
	}
	return msg
}
