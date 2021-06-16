package apiserver

import (
	//"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"net/http"
)

func (s *Server) PageshowUsersQualityP5() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		Admin := false
		SuperIngenerQuality := false
		GroupP5 := false
		LoggedIn := false

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		user, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		if user.Groups == "качество П5" {
			GroupP5 = true
			if user.Role == "Administrator" {
				Admin = true
				LoggedIn = true
			} else if user.Role == "главный инженер по качеству" {
				SuperIngenerQuality = true
				LoggedIn = true
			}
		}

		get, err := s.store.User().ListUsersQualityP5()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		data := map[string]interface{}{
			"TitleDOC":            "Сотрудники качества",
			"User":                user.LastName,
			"Username":            user.FirstName,
			"Admin":               Admin,
			"SuperIngenerQuality": SuperIngenerQuality,
			"GroupP5":             GroupP5,
			"LoggedIn":            LoggedIn,
			"GET":                 get,
		}
		err = tpl.ExecuteTemplate(w, "showUsersQuality.html", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}
