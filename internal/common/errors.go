package common

type ErrUserNotExists struct{}

func (ErrUserNotExists) Error() string {
	return "user not exists"
}
