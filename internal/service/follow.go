package service

import (
	"errors"
	"github.com/fitenne/youthcampus-dousheng/internal/common/code"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
)

// FollowUser 关注一名用户，返回true表示关注成功（注：在已关注的情况下再次关注被认为是错误的）
func FollowUser(userID int, toUserID int) (bool, error) {
	followdealer := repository.GetDealerFollow()
	followed, err := followdealer.CheckHasFollowed(userID, toUserID)
	if err != nil {
		return false, err
	}
	if followed {
		return false, errors.New(code.UserFollowed.Msg())
	}
	err = followdealer.FollowUser(userID, toUserID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CancelFollowUser 取消关注用户，返回true表示取消成功（在未关注的情况下取消关注被认为是错误的）
func CancelFollowUser(userID int, toUserID int) (bool, error) {
	followdealer := repository.GetDealerFollow()
	followed, err := followdealer.CheckHasFollowed(userID, toUserID)
	if err != nil {
		return false, err
	}
	if !followed {
		return false, errors.New(code.UserUnfollowed.Msg())
	}
	err = followdealer.CancelFollowUser(userID, toUserID)
	if err != nil {
		return false, err
	}
	return true, nil
}
