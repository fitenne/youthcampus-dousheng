package repository

/*
 * 此文件实现与 User 相关数据库的操作
 * 实现 UserCtl 接口
 * 对 User 表的操作应当通过 UserCtl 接口完成
 */

import (
	"errors"
	"time"

	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"gorm.io/gorm"
)

// user struct mapped to database
type User struct {
	ID            int64  `gorm:"primarykey"`
	UserName      string `gorm:"index:username,class:FULLTEXT,size:256"` // indexed for better authentication peformance
	Salt          []byte `gorm:"type:blob(32)"`
	Password      []byte `gorm:"type:blob(32)"`
	FollowCount   int64
	FollowerCount int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

type userCtl struct{}

var ctl userCtl

// 指定 User 对应的数据库表名
func (User) TableName() string {
	return "users"
}

func GetUserCtl() model.UserCtl {
	return &ctl
}

func (ctl *userCtl) QueryByID(id int64) (model.User, error) {
	user := make([]User, 0, 1)
	if res := dbProvider.GetDB().Limit(1).Find(&user, id); res.Error != nil {
		return model.User{}, res.Error
	}
	if len(user) == 0 {
		return model.User{}, errors.New("user not found")
	}

	return model.User{
		ID:            int64(user[0].ID),
		Name:          user[0].UserName,
		FollowCount:   user[0].FollowCount,
		FollowerCount: user[0].FollowerCount,
	}, nil
}

func (*userCtl) QueryByName(name string) (model.User, error) {
	user := make([]User, 0, 1)
	if res := dbProvider.GetDB().Limit(1).Find(&user, "user_name = ?", name); res.Error != nil {
		return model.User{}, res.Error
	}
	if len(user) == 0 {
		return model.User{}, errors.New("user not found")
	}

	return model.User{
		ID:            int64(user[0].ID),
		Name:          user[0].UserName,
		FollowCount:   user[0].FollowCount,
		FollowerCount: user[0].FollowerCount,
	}, nil
}

// 返回新用户的 ID
func (*userCtl) Create(name string, pass, salt []byte) (id int64, err error) {
	u := &User{
		UserName: name,
		Password: pass,
		Salt:     salt,
	}
	result := dbProvider.GetDB().Select("UserName", "Password", "Salt").Create(&u)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected != 1 {
		return 0, errors.New("failed to create user")
	}

	return u.ID, nil
}

func (*userCtl) QueryCredentialsByName(name string) (id int64, hashed []byte, salt []byte, err error) {
	user := make([]User, 0, 1)
	res := dbProvider.GetDB().Select("password", "salt").Limit(1).Find(&user, "user_name = ?", name)
	if res.Error != nil {
		return 0, nil, nil, res.Error
	}
	if len(user) == 0 {
		return 0, nil, nil, errors.New("user not exists")
	}

	return user[0].ID, user[0].Password, user[0].Salt, nil
}
