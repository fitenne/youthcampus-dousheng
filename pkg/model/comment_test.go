package model_test

import (
	"fmt"
	"testing"

	"github.com/fitenne/youthcampus-dousheng/internal/common/settings"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	"github.com/spf13/viper"
)

func init() {
	if err := settings.Init("../../config.yaml"); err != nil {
		panic(err.Error())
	}

	repository.Init(repository.DBConfig{
		Driver:   viper.GetString("db.driver"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.pass"),
		DBname:   viper.GetString("db.database"),
	})
}

func TestQueryListByVideoID(t *testing.T) {
	t.Run("TestComment", func(t *testing.T) {
		comCtl := repository.GetCommentCtl()
		comments, _ := comCtl.QueryListByVideoId(1)
		fmt.Println(comments)
	})
}
