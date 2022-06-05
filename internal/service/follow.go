package service

import (
	"errors"

	"github.com/fitenne/youthcampus-dousheng/internal/common/code"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

// FollowUser 关注一名用户，返回true表示关注成功（注：在已关注的情况下再次关注被认为是错误的，且不能关注自己）
func FollowUser(userID int64, toUserID int64) error {
	followdealer := repository.GetDealerFollow()
	userCTL := repository.GetUserCtl()

	//校验用户是否存在
	if _, err := userCTL.QueryByID(userID); err != nil {
		return err
	}
	if _, err := userCTL.QueryByID(toUserID); err != nil {
		return err
	}
	//不能自己关注自己
	if userID == toUserID {
		return errors.New(code.RepeatFollow.Msg())
	}
	//校验是否已关注
	followed, err := followdealer.CheckHasFollowed(userID, toUserID)
	if err != nil {
		return err
	}
	if followed {
		return errors.New(code.UserFollowed.Msg())
	}

	err = followdealer.FollowUser(userID, toUserID)
	if err != nil {
		return err
	}
	return nil
}

// CancelFollowUser 取消关注用户，返回true表示取消成功（在未关注的情况下取消关注被认为是错误的）
func CancelFollowUser(userID int64, toUserID int64) error {
	followdealer := repository.GetDealerFollow()
	userCTL := repository.GetUserCtl()
	//校验用户是否存在
	if _, err := userCTL.QueryByID(userID); err != nil {
		return err
	}
	if _, err := userCTL.QueryByID(toUserID); err != nil {
		return err
	}

	//校验是否未关注
	followed, err := followdealer.CheckHasFollowed(userID, toUserID)
	if err != nil {
		return err
	}
	if !followed {
		return errors.New(code.UserUnfollowed.Msg())
	}
	err = followdealer.CancelFollowUser(userID, toUserID)
	if err != nil {
		return err
	}
	return nil
}

// GetFollowList 获取关注列表的所有用户
func GetFollowList(userID int64) (*[]model.User, error) {
	userCTL := repository.GetUserCtl()
	//校验用户是否存在
	if _, err := userCTL.QueryByID(userID); err != nil {
		return nil, err
	}
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

// GetFollowerList 获取粉丝列表所有用户
func GetFollowerList(userID int64) (*[]model.User, error) {
	userCTL := repository.GetUserCtl()
	//校验用户是否存在
	if _, err := userCTL.QueryByID(userID); err != nil {
		return nil, err
	}
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
