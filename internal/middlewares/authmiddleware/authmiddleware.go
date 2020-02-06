package authmiddleware

/*
var store = sessions.NewCookieStore([]byte("starline"))

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		session, _ := store.Get(request, "starline")
		session.Values["user_id"] = u.ID
	//	email := session.Values["user_id"]
		if email == nil {
			http.Redirect(response, request, "/sessions", http.StatusSeeOther)
		} else {
			h.ServeHTTP(response, request)
		}
	})
}
*/
