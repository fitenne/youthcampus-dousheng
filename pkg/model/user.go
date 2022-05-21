package model

type User struct {
	ID            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

// 对数据库的修改应通过 model.UserCtl 完成
type UserCtl interface {
	// 找不到返回零值
	QueryByID(id int64) (User, error)
	QueryByName(name string) (User, error)

	// 获取id, 密码hash，盐
	QueryCredentialsByName(name string) (id int64, hashed []byte, salt []byte, err error)

	// 返回 id
	Create(name string, pass, salt []byte) (id int64, err error)
}
