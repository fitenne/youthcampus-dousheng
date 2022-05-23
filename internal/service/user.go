package service

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"errors"
	"log"

	"github.com/fitenne/youthcampus-dousheng/internal/common"
	"github.com/fitenne/youthcampus-dousheng/internal/common/jwt"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

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
	if n, err := rand.Reader.Read(salt); err != nil || n != saltSize {
		return 0, "", err
	}

	hasher := crypto.SHA256.New()
	_, err = hasher.Write(append(salt, []byte(password)...))
	if err != nil {
		return 0, "", err
	}
	dgst := make([]byte, 0, 32)
	dgst = hasher.Sum(dgst)

	id, err = repository.GetUserCtl().Create(username, dgst, salt)
	if err != nil {
		return 0, "", err
	}

	token, err = jwt.GenToken(id)
	if err != nil {
		return 0, "", err
	}

	return id, token, nil
}

// 返回 id, token，若登陆凭证无效，返回 (0, "", err)
func UserLogin(username, password string) (id int64, token string, err error) {
	log.Printf("UserLogin: %v : %v", username, password)
	id, p, s, err := repository.GetUserCtl().QueryCredentialsByName(username)
	if err != nil {
		return 0, "", err
	}

	hasher := crypto.SHA256.New()
	_, err = hasher.Write(append(s, []byte(password)...))
	if err != nil {
		return 0, "", err
	}

	dgst := make([]byte, 0, 32)
	dgst = hasher.Sum(dgst)
	if !hmac.Equal(dgst, p) {
		log.Printf("UserLogin: %v =? %v", dgst, p)
		return 0, "", errors.New("invalid credentials")
	}

	token, err = jwt.GenToken(id)
	if err != nil {
		return 0, "", err
	}

	return id, token, nil
}

func UserInfo(id int64) (model.User, error) {
	u, err := repository.GetUserCtl().QueryByID(id)
	if err != nil {
		return model.User{}, nil
	}

	return u, nil
}
