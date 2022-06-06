package repository

import (
	"os"
	"reflect"
	"testing"

	"github.com/fitenne/youthcampus-dousheng/internal/common/settings"
	"github.com/spf13/viper"

	"github.com/fitenne/youthcampus-dousheng/pkg/model"
)

func TestMain(m *testing.M) {
	if err := settings.Init("../../test.config.yaml"); err != nil {
		panic(err.Error())
	}

	c := DBConfig{
		Driver:   viper.GetString("db.driver"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.pass"),
		DBname:   viper.GetString("db.database"),
		LogLevel: viper.GetInt("db.loglevel"),
	}
	dbProvider = &MysqlProdiver{}
	if err := dbProvider.Connect(c); err != nil {
		panic(err.Error())
	}
	if err := dbProvider.GetDB().AutoMigrate(&User{}); err != nil {
		panic(err.Error())
	}

	os.Exit(m.Run())
}

func TestQueryByID(t *testing.T) {
	cleanUpTable(t)
	ctl := GetUserCtl()

	{
		name := "query should find"
		want := model.User{
			Name:          "alice",
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		wantErr := false

		t.Run(name, func(t *testing.T) {
			u := User{
				ID:            0,
				UserName:      want.Name,
				Salt:          []byte{},
				Password:      []byte{},
				FollowCount:   0,
				FollowerCount: 0,
			}
			if err := dbProvider.GetDB().Create(&u).Error; err != nil {
				t.Fatal(err.Error())
			}

			got, err := ctl.QueryByID(u.ID)
			if (err != nil) != wantErr {
				t.Errorf("userCtl.QueryByID() error = %v, wantErr %v", err, wantErr)
			}

			want.ID = got.ID
			if !reflect.DeepEqual(got, want) {
				t.Errorf("userCtl.QueryByID() = %v, want %v", got, want)
			}
		})
	}

	{
		name := "query not found"
		wantErr := true

		t.Run(name, func(t *testing.T) {
			_, err := ctl.QueryByID(114514)
			if (err != nil) != wantErr {
				t.Errorf("userCtl.QueryByID() error = %v, wantErr %v", err, wantErr)
			}
		})
	}
}

func TestQueryByName(t *testing.T) {
	cleanUpTable(t)
	ctl := GetUserCtl()

	{
		name := "query should find"
		want := model.User{
			Name:          "alice",
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		wantErr := false

		t.Run(name, func(t *testing.T) {
			u := User{
				ID:            0,
				UserName:      want.Name,
				Salt:          []byte{},
				Password:      []byte{},
				FollowCount:   0,
				FollowerCount: 0,
			}
			if err := dbProvider.GetDB().Create(&u).Error; err != nil {
				t.Fatal(err.Error())
			}

			got, err := ctl.QueryByName(u.UserName)
			if (err != nil) != wantErr {
				t.Errorf("userCtl.QueryByID() error = %v, wantErr %v", err, wantErr)
			}

			want.ID = got.ID
			if !reflect.DeepEqual(got, want) {
				t.Errorf("userCtl.QueryByID() = %v, want %v", got, want)
			}
		})
	}

	{
		name := "query not found"
		wantErr := true

		t.Run(name, func(t *testing.T) {
			_, err := ctl.QueryByName("114514")
			if (err != nil) != wantErr {
				t.Errorf("userCtl.QueryByID() error = %v, wantErr %v", err, wantErr)
			}
		})
	}
}

func TestQueryCredentialsByName(t *testing.T) {
	cleanUpTable(t)
	ctl := GetUserCtl()

	tt := struct {
		name       string
		argsName   string
		wantId     int64
		wantHashed []byte
		wantSalt   []byte
		wantErr    bool
	}{
		name:       "query credentials by name",
		argsName:   "alice",
		wantHashed: []byte{0x1, 0x1, 0x0},
		wantSalt:   []byte{0x1, 0x1, 0x0},
		wantErr:    false,
	}

	t.Run(tt.name, func(t *testing.T) {
		dbProvider.GetDB().Create(&User{
			UserName: tt.argsName,
			Salt:     tt.wantSalt,
			Password: tt.wantHashed,
		})

		_, gotHashed, gotSalt, err := ctl.QueryCredentialsByName(tt.argsName)
		if (err != nil) != tt.wantErr {
			t.Errorf("userCtl.QueryCredentialsByName() error = %v, wantErr %v", err, tt.wantErr)
		}
		if !reflect.DeepEqual(gotHashed, tt.wantHashed) {
			t.Errorf("userCtl.QueryCredentialsByName() gotHashed = %v, want %v", gotHashed, tt.wantHashed)
		}
		if !reflect.DeepEqual(gotSalt, tt.wantSalt) {
			t.Errorf("userCtl.QueryCredentialsByName() gotSalt = %v, want %v", gotSalt, tt.wantSalt)
		}
	})
}

func TestCreate(t *testing.T) {
	cleanUpTable(t)
	ctl := GetUserCtl()

	type args struct {
		name string
		pass []byte
		salt []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal create",
			args: args{
				name: "12345678901234567890123456789012",
				pass: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2},
				salt: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2},
			},
			wantErr: false,
		},
		{
			name: "empty pass",
			args: args{
				name: "1",
				pass: []byte{},
				salt: []byte("10086"),
			},
			wantErr: true,
		},
		{
			name: "empty salt",
			args: args{
				name: "2",
				pass: []byte("10086"),
				salt: []byte(""),
			},
			wantErr: true,
		},
		{
			name: "too long name",
			args: args{
				name: "123456789012345678901234567890123",
				pass: []byte("10086"),
				salt: []byte("10086"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ctl.Create(tt.args.name, tt.args.pass, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("userCtl.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func cleanUpTable(t *testing.T) {
	res := dbProvider.GetDB().Model(&User{}).Unscoped().Where("TRUE").Delete(&User{})
	if res.Error != nil {
		t.Fatal(res.Error.Error())
	}
}
