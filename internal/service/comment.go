package service

import (
	"errors"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/fitenne/youthcampus-dousheng/pkg/model"
	"log"
)

var commentCtl = repository.GetCommentCtl()

func Publish(userId int64, comment *model.Comment) error {

	// 获取用户信息
	user, err := repository.GetUserCtl().QueryByID(userId)
	if err != nil {
		log.Println("service.Publish|server error|repository.GetUserCtl.QueryByID: " + err.Error())
		return err
	}

	// 发布
	err = commentCtl.Publish(comment)
	if err != nil {
		return err
	}

	// 填充用户信息
	comment.User = &model.UserEntity{
		ID:            user.ID,
		UserName:      user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      false,
	}

	return nil
}

func DeleteById(userId, commentId int64) error {

	comment, err := commentCtl.QueryById(commentId)
	if err != nil {
		return err
	}

	if userId != comment.UserID {
		return errors.New("用户无删除权限")
	}

	return commentCtl.DeleteById(commentId)
}

func QueryListByVideoId(videoId, userId int64) ([]model.Comment, error) {

	// 获取评论list，异常处理
	comments, err := commentCtl.QueryListByVideoId(videoId)
	if err != nil {
		return nil, err
	}

	// userId 不合法，直接返回
	if userId <= 0 {
		return comments, nil
	}

	// 查看用户是否关注 (对客户端业务没有可见的影响，为了提高效率，暂时注释掉)
	//pass := true
	//for i := range comments {
	//	followed, err := repository.GetDealerFollow().CheckHasFollowed(userId, comments[i].UserID)
	//	if err != nil {
	//		pass = false
	//		comments[i].User.IsFollow = false
	//		continue
	//	}
	//
	//	comments[i].User.IsFollow = followed
	//}

	// 异常处理
	//if !pass {
	//	log.Println("用户关注评论者校验异常：" + strconv.Itoa(int(userId)))
	//}

	return comments, nil
}
