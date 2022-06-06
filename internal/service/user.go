package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"

	"github.com/fitenne/youthcampus-dousheng/internal/common"
	"github.com/fitenne/youthcampus-dousheng/internal/common/jwt"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

func getDgst(salt []byte, T string) ([]byte, error) {
	hasher := sha256.New()
	_, err := hasher.Write(append(salt, []byte(T)...))
	if err != nil {
		return nil, err
	}

	dgst := make([]byte, 0, 32)
	dgst = hasher.Sum(dgst)
	return dgst, nil
}

func UserExists(username string) (bool, error) {
	_, err := repository.GetUserCtl().QueryByName(username)
	if err != nil {
		if errors.Is(err, common.ErrUserNotExists{}) {
			err = nil
		}
		return false, err
	}

	return true, nil
}

func UserRegister(username, password string) (id int64, token string, err error) {
	const saltSize = 32
	salt := make([]byte, saltSize)
	if n, err := rand.Reader.Read(salt); err != nil {
		if n != saltSize {
			return 0, "", errors.New("filaed to generate salt")
		}
		return 0, "", err
	}

	dgst, err := getDgst(salt, password)
	if err != nil {
		return 0, "", err
	}

	id, err = repository.GetUserCtl().Create(username, dgst, salt)
	if err != nil {
		return 0, "", err
	}

	token, err = jwt.GenToken(id)
	return id, token, nil
}

// 返回 id, token，若登陆凭证无效，返回 (0, "", nil)
func UserLogin(username, password string) (id int64, token string, err error) {
	id, p, s, err := repository.GetUserCtl().QueryCredentialsByName(username)
	if err != nil {
		if errors.Is(err, common.ErrUserNotExists{}) {
			return 0, "", nil
		}
		return 0, "", err
	}

	dgst, err := getDgst(s, password)
	if err != nil || !hmac.Equal(dgst, p) {
		return 0, "", err
	}

	token, err = jwt.GenToken(id)
	return id, token, nil
}

func UserInfo(myID, id int64) (model.User, error) {
	u, err := repository.GetUserCtl().QueryByID(id)
	if err != nil {
		return model.User{}, nil
	}

	u.IsFollow, err = repository.GetDealerFollow().CheckHasFollowed(myID, id)
	return u, err
}
