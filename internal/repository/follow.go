package repository

import (
	"time"
)

//  这个文件中的函数仅提供了相关操作的接口，未对业务逻辑做过判断

type Follow struct {
	UserID     int       `json:"user_id,omitempty" gorm:"column:user_id;not null"`
	FollowedID int       `json:"followed_id,omitempty" gorm:"column:followed_id;not null"`
	ID         uint      `gorm:"primaryKey"`
	CreateAt   time.Time `gorm:"autoCreateTime;not null"`
}

func (Follow) TableName() string {
	return "follow"
}

type FollowCtl interface {
	CheckHasFollowed(userID int, ToUserID int) (bool, error)
	FollowUser(userID int, ToUserID int) error
	CancelFollowUser(userID int, ToUserID int) error
	SelectAllFollower(userID int) (*[]int, error) //查询所有粉丝
	SelectAllFollowed(userID int) (*[]int, error) //查询所有已经关注的
}

type FollowCtlDealer struct{}

func (f *FollowCtlDealer) SelectAllFollower(userID int) (*[]int, error) {
	var ids []int
	err := dbProvider.GetDB().Raw("select followed_id from follow where user_id = ?", userID).Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	return &ids, err
}

func (f *FollowCtlDealer) SelectAllFollowed(userID int) (*[]int, error) {
	var ids []int
	err := dbProvider.GetDB().Raw("select user_id from follow where followed_id = ?", userID).Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	return &ids, err
}

func (f *FollowCtlDealer) FollowUser(userID int, ToUserID int) error {
	follow := &Follow{
		UserID:     userID,
		FollowedID: ToUserID,
	}
	if err := dbProvider.GetDB().Create(follow).Error; err != nil {
		return err
	}
	return nil
}

func (f *FollowCtlDealer) CancelFollowUser(userID int, ToUserID int) error {
	err := dbProvider.GetDB().Where("user_id = ? and followed_id = ?", userID, ToUserID).Delete(&Follow{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (f *FollowCtlDealer) CheckHasFollowed(userID int, ToUserID int) (bool, error) {
	var ids []int
	err := dbProvider.GetDB().Raw("select followed_id from follow where user_id = ?", userID).Scan(&ids).Error
	if err != nil {
		return false, err
	}
	for _, id := range ids {
		if id == ToUserID {
			return true, nil
		}
	}
	return false, nil
}

// GetDealerFollow 这个函数显式的实现了接口，若未实现接口会报错
func GetDealerFollow() FollowCtl {
	return &FollowCtlDealer{}
}
