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

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupQuality {

			if user.Role == roleAdministrator {
				get, err := s.store.User().ListUsersQuality()
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
				params := []Param{
					{
						GroupP1:  true,
						Admin:    true,
						LoggedIn: true,
						User:     user.LastName,
						Username: user.FirstName,
						GET:      get,
					},
				}

				data := map[string]interface{}{
					"TitleDOC": "Сотрудники качества",
					"GetParam": params,
				}

				RenderTemplate(w, "showUsersQuality.html", data)

			} else if user.Role == roleSuperIngenerQuality {
				get, err := s.store.User().ListUsersQuality()
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
				params := []Param{
					{
						GroupP1:             true,
						SuperIngenerQuality: true,
						LoggedIn:            true,
						User:                user.LastName,
						Username:            user.FirstName,
						GET:                 get,
					},
				}

				data := map[string]interface{}{
					"TitleDOC": "Сотрудники качества",
					"GetParam": params,
				}

				RenderTemplate(w, "showUsersQuality.html", data)

			}

			/*
				getRoleQuality, err := s.store.RoleQuality().ListRoleQuality()
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
			*/
		}

		if user.Groups == groupQualityP5 {

			if user.Role == roleAdministrator {
				get, err := s.store.User().ListUsersQualityP5()
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
				params := []Param{
					{
						GroupP5:  true,
						Admin:    true,
						LoggedIn: true,
						User:     user.LastName,
						Username: user.FirstName,
						GET:      get,
					},
				}

				data := map[string]interface{}{
					"TitleDOC": "Сотрудники качества",
					"GetParam": params,
				}

				RenderTemplate(w, "showUsersQuality.html", data)

			} else if user.Role == roleSuperIngenerQuality {
				get, err := s.store.User().ListUsersQualityP5()
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
				params := []Param{
					{
						GroupP5:             true,
						SuperIngenerQuality: true,
						LoggedIn:            true,
						User:                user.LastName,
						Username:            user.FirstName,
						GET:                 get,
					},
				}

				data := map[string]interface{}{
					"TitleDOC": "Сотрудники качества",
					"GetParam": params,
				}

				RenderTemplate(w, "showUsersQuality.html", data)

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

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())

		}
		var hdata []requestFrom
		json.Unmarshal(body, &hdata)
		fmt.Printf("body json: %s", body)
		s.logger.Infof("Loading body json: %s\n", body)
		fmt.Println("\njson  struct hdata", hdata)
		s.logger.Infof("Loading hdata json: %v\n", hdata)

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupQuality {
			if user.Role == roleAdministrator {
				//	Admin = true
				//	LoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {

				for _, v := range hdata {
					fmt.Println(v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)
					s.logger.Infof("P1 create quality employee: %v, %v, %v, %v, %v, %v\n", v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)
					u := &model.User{
						Email:     v.Email,
						Password:  v.Password,
						FirstName: v.FirstName,
						LastName:  v.LastName,
						Role:      v.Role,
						Groups:    groupQuality,
						Tabel:     v.Tabel,
					}

					if err := s.store.User().CreateUserByManager(u); err != nil {
						s.error(w, r, http.StatusUnprocessableEntity, err)
						return
					}
				}
			}
		}
		if user.Groups == groupQualityP5 {
			if user.Role == roleAdministrator {
				//	Admin = true
				//	LoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				for _, v := range hdata {
					fmt.Println(v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)
					s.logger.Infof("P5 create quality employee: %v, %v, %v, %v, %v, %v\n", v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)

					u := &model.User{
						Email:     v.Email,
						Password:  v.Password,
						FirstName: v.FirstName,
						LastName:  v.LastName,
						Role:      v.Role,
						Groups:    groupQualityP5,
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
}

func (s *Server) PageupdateUserQuality() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		user := r.Context().Value(ctxKeyUser).(*model.User)

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		fmt.Println("var id - ", id)
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}

		get, err := s.store.User().EditUserByManager(id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if user.Role == roleAdministrator {

		} else if user.Role == roleSuperIngenerQuality {
			params := []Param{
				{
					LoggedIn:            true,
					SuperIngenerQuality: true,
					GET:                 get,
					User:                user.LastName,
					Username:            user.FirstName,
				},
			}

			data := map[string]interface{}{
				"GetParam": params,
			}

			RenderTemplate(w, "updateuserquality.html", data)
		}
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
			s.logger.Errorf(err.Error())
		}
		req.ID = id
		req.Email = r.FormValue("email")
		req.Firstname = r.FormValue("firstname")
		req.Lastname = r.FormValue("lastname")
		req.Role = r.FormValue("role")
		//fmt.Println("Роль - ", req.Role)
		req.Tabel = r.FormValue("tabel")
		fmt.Println("ID - ", req.ID)
		s.logger.Infof("Update quality employee: %v, %v, %v, %v, %v, %v\n", req.ID, req.Email, req.Firstname, req.Lastname, req.Role, req.Tabel)
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
			s.logger.Errorf(err.Error())
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

