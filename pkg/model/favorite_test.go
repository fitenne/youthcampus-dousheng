package model_test

import (
	// "log"
	"testing"

	"github.com/fitenne/youthcampus-dousheng/internal/common/settings"
	"github.com/fitenne/youthcampus-dousheng/internal/repository"
	// "github.com/fitenne/youthcampus-dousheng/pkg/model"
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

func TestCreateTable(t *testing.T) {
	t.Run("TestCreateTable", func(t *testing.T) {

		repository.GetFavoriteCtl().CreateTableTest()
	})
}
