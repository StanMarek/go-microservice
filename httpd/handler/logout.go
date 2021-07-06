package handler

import "net/http"

func Logout(writer http.ResponseWriter, request *http.Request) {
	Extract := func(request *http.Request) (*AccessDetails, error) {

	}
}

func DeleteAuthorization(authId string) (int64, error) {
	del, err := clientRedis.Del(ctxRedis, authId).Result()
	if err != nil {
		return 0, err
	}
	return del, nil
}
