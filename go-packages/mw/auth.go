package mw

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/Shonei/student-information-system/go-packages/dba"
	"github.com/Shonei/student-information-system/go-packages/dbc"
)

// BasicAuth is the authorization middleware for get routes
// that accept a /{user} in the url. It compares the user from the url and the
// user in the token and data in the database. It grants access to people that
// have a hign enought level or are the owner of that information.
func BasicAuth(db dba.DBAbstraction, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		lvl := dbc.CheckToken(db, token)

		switch lvl {
		case -1:
			http.Error(w, "Wrong credentials send.", http.StatusUnauthorized)
			return
		case 1:
			user := strings.Split(token, ":")[0]
			vars := mux.Vars(r)
			if user == vars["user"] {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "You don't have the authority to access that resource", http.StatusUnauthorized)
				return
			}
		case 2, 3:
			next.ServeHTTP(w, r)
			return
		}
	})
}
