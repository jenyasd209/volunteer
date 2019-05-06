package data

import (
	"fmt"
	"net/http"
)

func SessionChek(r *http.Request, session *Session) (err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		// (*session).SetUUID(cookie.Value)
		(*session).UUID = cookie.Value
		if ok, err := session.Check(); !ok {
			fmt.Println(err)
		}
	}

	return
}
