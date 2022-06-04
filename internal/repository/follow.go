package repository

import (
	"time"

	"gorm.io/gorm"
)

//  这个文件中的函数仅提供了相关操作的接口，未对业务逻辑做过判断

type Follow struct {
	UserID     int64     `json:"user_id,omitempty" gorm:"column:user_id;not null"`
	FollowedID int64     `json:"followed_id,omitempty" gorm:"column:followed_id;not null"`
	ID         uint      `gorm:"primaryKey"`
	CreateAt   time.Time `gorm:"autoCreateTime;not null"`
}

func (Follow) TableName() string {
	return "follow"
}

type FollowCtl interface {
	CheckHasFollowed(userID int64, ToUserID int64) (bool, error)
	FollowUser(userID int64, ToUserID int64) error
	CancelFollowUser(userID int64, ToUserID int64) error
	SelectAllFollower(userID int64) (*[]int64, error) //查询所有粉丝
	SelectAllFollowed(userID int64) (*[]int64, error) //查询所有已经关注的
	AddFollower(userID int64) error
	SubFollower(userID int64) error
	AddFollowed(userID int64) error
	SubFollowed(userID int64) error
}

type FollowCtlDealer struct{}

func (f *FollowCtlDealer) AddFollower(userID int64) error {
	err := dbProvider.GetDB().Exec("update users set follower_count = follower_count+1 where id = ?", userID).Error
	if err != nil {
		return err
	}
	return err
}

func (f *FollowCtlDealer) SubFollower(userID int64) error {
	err := dbProvider.GetDB().Exec("update users set follower_count = follower_count-1 where id = ?", userID).Error
	if err != nil {
		return err
	}
	return err
}

func (f *FollowCtlDealer) AddFollowed(userID int64) error {
	err := dbProvider.GetDB().Exec("update users set follow_count = follow_count+1 where id = ?", userID).Error
	if err != nil {
		return err
	}
	return err
}

func (f *FollowCtlDealer) SubFollowed(userID int64) error {
	err := dbProvider.GetDB().Exec("update users set follow_count = follow_count-1 where id = ?", userID).Error
	if err != nil {
		return err
	}
	return err
}

func (f *FollowCtlDealer) SelectAllFollower(userID int64) (*[]int64, error) {
	var ids []int64
	err := dbProvider.GetDB().Raw("select user_id from follow where followed_id = ?", userID).Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	return &ids, err
}

func (f *FollowCtlDealer) SelectAllFollowed(userID int64) (*[]int64, error) {
	var ids []int64
	//查询关注列表的用户id
	err := dbProvider.GetDB().Raw("select followed_id from follow where user_id = ?", userID).Scan(&ids).Error
	if err != nil {
		return nil, err
	}
	return &ids, err
}

func (f *FollowCtlDealer) FollowUser(userID int64, ToUserID int64) error {
	follow := &Follow{
		UserID:     userID,
		FollowedID: ToUserID,
	}
	//开启事务  任何错误都将回滚   在事务中应该使用tx而不是dbProvider
	return dbProvider.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(follow).Error; err != nil {
			return err
		}
		//粉丝列表增加一位
		if err := tx.Exec("update users set follower_count = follower_count+1 where id = ?", ToUserID).Error; err != nil {
			return err
		}
		//关注列表增加一位
		if err := tx.Exec("update users set follow_count = follow_count+1 where id = ?", userID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (f *FollowCtlDealer) CancelFollowUser(userID int64, ToUserID int64) error {
	return dbProvider.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? and followed_id = ?", userID, ToUserID).Delete(&Follow{}).Error; err != nil {
			return err
		}
		//粉丝列表减少一位
		if err := tx.Exec("update users set follower_count = follower_count-1 where id = ?", ToUserID).Error; err != nil {
			return err
		}
		//关注列表减少一位
		if err := tx.Exec("update users set follow_count = follow_count-1 where id = ?", userID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (f *FollowCtlDealer) CheckHasFollowed(userID int64, ToUserID int64) (bool, error) {
	var ids []int64
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
