package apiserver

import (
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
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

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupQualityP5 {
			GroupP5 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
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
