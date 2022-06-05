package controller

import "github.com/fitenne/youthcampus-dousheng/pkg/model"

var DemoVideos = []model.Video{
	{
		ID:            1,
		Author:        &DemoUser,
		PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
		CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
}

var DemoUser = model.User{
	ID:            1,
	Name:          "TestUser",
	FollowCount:   10,
	FollowerCount: 15,
	IsFollow:      true,
}
