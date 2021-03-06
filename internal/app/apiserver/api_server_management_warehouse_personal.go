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

func (s *Server) PageshowUsersWarehouse() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupWarehouse {
			group := groupWarehouse
			if user.Role == roleAdministrator {
				//	Admin = true
				//	LoggedIn = true
			} else if user.Role == roleWarehouseManager {
				get, err := s.store.User().ListUsersWarehouse(group)
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
				params := []Param{
					{
						GroupP1:          true,
						LoggedIn:         true,
						WarehouseManager: true,
						User:             user.LastName,
						Username:         user.FirstName,
						GET:              get,
					},
				}

				data := map[string]interface{}{
					"TitleDOC": "Сотрудники склада",
					"GetParam": params,
				}

				RenderTemplate(w, "showUsersWarehouse.html", data)

			}
		}

		if user.Groups == groupWarehouseP5 {
			group := groupWarehouseP5
			if user.Role == roleAdministrator {
				//Admin = true
				//LoggedIn = true
			} else if user.Role == roleWarehouseManager {
				get, err := s.store.User().ListUsersWarehouse(group)
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
				params := []Param{
					{
						GroupP1:          true,
						LoggedIn:         true,
						WarehouseManager: true,
						User:             user.LastName,
						Username:         user.FirstName,
						GET:              get,
					},
				}

				data := map[string]interface{}{
					"TitleDOC": "Сотрудники склада",
					"GetParam": params,
				}

				RenderTemplate(w, "showUsersWarehouse.html", data)

			}
		}
	}
}

func (s *Server) CreateUserWarehouse() http.HandlerFunc {
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

		if user.Groups == groupWarehouse {
			if user.Role == roleWarehouseManager {

				for _, v := range hdata {
					fmt.Println(v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)
					s.logger.Infof("P1 create warehouse employee:: %v, %v, %v, %v, %v, %v\n", v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)
					u := &model.User{
						Email:     v.Email,
						Password:  v.Password,
						FirstName: v.FirstName,
						LastName:  v.LastName,
						Role:      v.Role,
						Groups:    groupWarehouse,
						Tabel:     v.Tabel,
					}

					if err := s.store.User().CreateUserByManager(u); err != nil {
						s.error(w, r, http.StatusUnprocessableEntity, err)
						return
					}

				}
			}
		}

		if user.Groups == groupWarehouseP5 {
			if user.Role == roleWarehouseManager {

				for _, v := range hdata {
					fmt.Println("create stockkeeper P5", v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)
					s.logger.Infof("P5 create warehouse employee: %v, %v, %v, %v, %v, %v\n", v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)

					u := &model.User{
						Email:     v.Email,
						Password:  v.Password,
						FirstName: v.FirstName,
						LastName:  v.LastName,
						Role:      v.Role,
						Groups:    groupWarehouseP5,
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

func (s *Server) PageupdateUserWarehouse() http.HandlerFunc {

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

		} else if user.Role == roleWarehouseManager {
			params := []Param{
				{
					LoggedIn:         true,
					WarehouseManager: true,
					GET:              get,
					User:             user.LastName,
					Username:         user.FirstName,
				},
			}

			data := map[string]interface{}{
				"GetParam": params,
			}

			RenderTemplate(w, "updateuserwarehouse.html", data)

		}
	}
}

func (s *Server) UpdateUserWarehouse() http.HandlerFunc {
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
		s.logger.Infof("Update warehouse employee: %v, %v, %v, %v, %v, %v\n", req.ID, req.Email, req.Firstname, req.Lastname, req.Role, req.Tabel)
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

		http.Redirect(w, r, "/operation/showuserswarehouse", 303)
	}

}

func (s *Server) DeleteUserWarehouse() http.HandlerFunc {
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

		http.Redirect(w, r, "/operation/showuserswarehouse", 303)
	}
}
