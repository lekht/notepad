package auth

import "errors"

var ErrNoUser = errors.New("this user does not exist")

// Заглушка вместо redis
type DummyRedis map[int]string

func New() DummyRedis {
	dRedis := make(DummyRedis, 10)

	dRedis[10] = "dskc90832hcusad9"

	return dRedis
}

// Метод выполняет авторизацию пользователя
func (dr DummyRedis) Authorization(id int) (string, error) {
	token, ok := dr[id]
	if !ok {
		return "", ErrNoUser
	}

	return token, nil
}
