package service

import (
	"errors"

	"github.com/fitenne/youthcampus-dousheng/internal/common/code"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

// FollowUser 关注一名用户，返回true表示关注成功（注：在已关注的情况下再次关注被认为是错误的）
func FollowUser(userID int64, toUserID int64) (bool, error) {
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
func CancelFollowUser(userID int64, toUserID int64) (bool, error) {
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

// GetFollowList 获取关注列表的所有用户
func GetFollowList(userID int64) (*[]model.User, error) {
	followdealer := repository.GetDealerFollow()
	follows, err := followdealer.SelectAllFollowed(userID)
	if err != nil {
		return nil, err
	}
	userCTl := repository.GetUserCtl()
	var users []model.User
	for _, v := range *follows {
		user, err := userCTl.QueryByID(v)
		if err != nil {
			return nil, err
		}
		user.IsFollow = true // 关注列表的用户可以确认已关注
		users = append(users, user)
	}
	return &users, nil
}

func GetFollowerList(userID int64) (*[]model.User, error) {
	followdealer := repository.GetDealerFollow()
	follows, err := followdealer.SelectAllFollower(userID)
	if err != nil {
		return nil, err
	}
	userCTl := repository.GetUserCtl()
	var users []model.User
	for _, v := range *follows {
		user, err := userCTl.QueryByID(v)
		if err != nil {
			return nil, err
		}
		isfollowed, err := followdealer.CheckHasFollowed(userID, v)
		if err != nil {
			return nil, err
		}
		user.IsFollow = isfollowed
		users = append(users, user)
	}
	return &users, nil
}
