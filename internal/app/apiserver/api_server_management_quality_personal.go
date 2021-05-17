package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"github.com/gorilla/mux"
)

func (s *Server) PageshowUsersQuality() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		Admin := false
		SuperIngenerQuality := false
		GroupP1 := false
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

		if user.Groups == "качество" {
			GroupP1 = true
			if user.Role == "Administrator" {
				Admin = true
				LoggedIn = true
			} else if user.Role == "главный инженер по качеству" {
				SuperIngenerQuality = true
				LoggedIn = true
			}

			get, err := s.store.User().ListUsersQuality()
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
				"GroupP1":             GroupP1,
				"LoggedIn":            LoggedIn,
				"GET":                 get,
			}
			err = tpl.ExecuteTemplate(w, "showUsersQuality.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
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
}

func (s *Server) CreateUserQuality() http.HandlerFunc {
	type requestFrom struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
		Tabel     string `json:"tabel"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		//	Admin := false
		SuperIngenerQuality := false
		//	LoggedIn := false

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			s.errorLog.Printf(err.Error())
		}
		var hdata []requestFrom
		json.Unmarshal(body, &hdata)
		fmt.Printf("body json: %s", body)
		s.infoLog.Printf("Loading body json: %s\n", body)
		fmt.Println("\njson  struct hdata", hdata)
		s.infoLog.Printf("Loading hdata json: %v\n", hdata)

		Groupp1 := "качество"
		Groupp5 := "качество П5"

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		idd, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		user, err := s.store.User().Find(idd.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		if user.Groups == "качество" {
			if user.Role == "Administrator" {
				//	Admin = true
				//	LoggedIn = true
			} else if user.Role == "главный инженер по качеству" {
				SuperIngenerQuality = true
				//	LoggedIn = true
				fmt.Println("SuperIngenerQuality pageupdateUserQuality - ", SuperIngenerQuality)
			}

			for _, v := range hdata {
				fmt.Println(v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)
				s.infoLog.Printf("P1 create quality employee: %v, %v, %v, %v, %v, %v\n", v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)
				u := &model.User{
					Email:     v.Email,
					Password:  v.Password,
					FirstName: v.FirstName,
					LastName:  v.LastName,
					Role:      v.Role,
					Groups:    Groupp1,
					Tabel:     v.Tabel,
				}

				if err := s.store.User().CreateUserByManager(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

			}
		}
		if user.Groups == "качество П5" {
			if user.Role == "Administrator" {
				//	Admin = true
				//	LoggedIn = true
			} else if user.Role == "главный инженер по качеству" {
				SuperIngenerQuality = true
				//	LoggedIn = true
				fmt.Println("SuperIngenerQuality pageupdateUserQuality - ", SuperIngenerQuality)
			}

			for _, v := range hdata {
				fmt.Println(v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)
				s.infoLog.Printf("P5 create quality employee: %v, %v, %v, %v, %v, %v\n", v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)

				u := &model.User{
					Email:     v.Email,
					Password:  v.Password,
					FirstName: v.FirstName,
					LastName:  v.LastName,
					Role:      v.Role,
					Groups:    Groupp5,
					Tabel:     v.Tabel,
				}

				if err := s.store.User().CreateUserByManager(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

			}
		}

	}
}

func (s *Server) PageupdateUserQuality() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		Admin := false
		SuperIngenerQuality := false
		LoggedIn := false

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		idd, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		user, err := s.store.User().Find(idd.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true
		} else if user.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true
			fmt.Println("SuperIngenerQuality pageupdateUserQuality - ", SuperIngenerQuality)
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		fmt.Println("var id - ", id)
		if err != nil {
			log.Println(err)
			s.errorLog.Printf(err.Error())
		}

		get, err := s.store.User().EditUserByManager(id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		data := map[string]interface{}{
			"GET":                 get,
			"Admin":               Admin,
			"SuperIngenerQuality": SuperIngenerQuality,
			"LoggedIn":            LoggedIn,
			"User":                user.LastName,
			"Username":            user.FirstName,
		}

		//	fmt.Println("Get.email" - get.id)
		tpl.ExecuteTemplate(w, "updateuserquality.html", data)
	}
}

func (s *Server) UpdateUserQuality() http.HandlerFunc {
	type request struct {
		ID        int    `json:"ID"`
		Email     string `json:"email"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Role      string `json:"role"`
		Tabel     string `json:"tabel"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.errorLog.Printf(err.Error())
		}
		req.ID = id
		req.Email = r.FormValue("email")
		req.Firstname = r.FormValue("firstname")
		req.Lastname = r.FormValue("lastname")
		req.Role = r.FormValue("role")
		//fmt.Println("Роль - ", req.Role)
		req.Tabel = r.FormValue("tabel")
		fmt.Println("ID - ", req.ID)
		s.infoLog.Printf("Update quality employee: %v, %v, %v, %v, %v, %v\n", req.ID, req.Email, req.Firstname, req.Lastname, req.Role, req.Tabel)
		u := &model.User{
			ID:        req.ID,
			Email:     req.Email,
			FirstName: req.Firstname,
			LastName:  req.Lastname,
			Role:      req.Role,
			Tabel:     req.Tabel,
		}

		if err := s.store.User().UpdateUserByManager(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		http.Redirect(w, r, "/operation/showusersquality", 303)
	}

}

func (s *Server) DeleteUserQuality() http.HandlerFunc {
	type request struct {
		ID int `json:"ID"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.errorLog.Printf(err.Error())
		}
		req.ID = id

		u := &model.User{
			ID: req.ID,
		}

		if err := s.store.User().DeleteUserByManager(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		http.Redirect(w, r, "/operation/showusersquality", 303)
	}
}
