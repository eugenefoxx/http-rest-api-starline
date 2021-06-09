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

		//user := r.Context().Value(ctxKeyUser).(*model.User)
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		user, err := s.store.User().Find(id.(int))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			//	s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		if user.Groups == groupWarehouse {
			statusGroupP1 = true
			group := "склад"
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
			}

			get, err := s.store.User().ListUsersWarehouse(group)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			data := map[string]interface{}{
				"TitleDOC":         "Сотрудники склада",
				"User":             user.LastName,
				"Username":         user.FirstName,
				"Admin":            statusAdmin,
				"WarehouseManager": statusWarehouseManager,
				"GroupP1":          statusGroupP1,
				"LoggedIn":         statusLoggedIn,
				"GET":              get,
			}
			err = tpl.ExecuteTemplate(w, "showUsersWarehouse.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}

		if user.Groups == groupWarehouseP5 {
			statusGroupP5 = true
			group := "склад П5"
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
			}

			get, err := s.store.User().ListUsersWarehouse(group)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			data := map[string]interface{}{
				"TitleDOC":         "Сотрудники склада",
				"User":             user.LastName,
				"Username":         user.FirstName,
				"Admin":            statusAdmin,
				"WarehouseManager": statusWarehouseManager,
				"GroupP5":          statusGroupP5,
				"LoggedIn":         statusLoggedIn,
				"GET":              get,
			}
			err = tpl.ExecuteTemplate(w, "showUsersWarehouse.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
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
		//	WarehouseManager := false

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

		Groupp1 := "склад"
		Groupp5 := "склад П5"

		//user := r.Context().Value(ctxKeyUser).(*model.User)
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		user, err := s.store.User().Find(id.(int))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			//	s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		//group := "склад"
		if user.Groups == groupWarehouse {
			if user.Role == roleWarehouseManager {
				//	WarehouseManager = true
				for _, v := range hdata {
					fmt.Println(v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)
					s.logger.Infof("P1 create warehouse employee:: %v, %v, %v, %v, %v, %v\n", v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)
					u := &model.User{
						Email:     v.Email,
						Password:  v.Password,
						FirstName: v.FirstName,
						LastName:  v.LastName,
						Role:      v.Role,
						Groups:    Groupp1,
						Tabel:     v.Tabel,
					}
					//	s.Lock()
					if err := s.store.User().CreateUserByManager(u); err != nil {
						s.error(w, r, http.StatusUnprocessableEntity, err)
						return
					}
					//	s.Unlock()
				}
			}
		}

		if user.Groups == groupWarehouseP5 {
			if user.Role == roleWarehouseManager {
				//	WarehouseManager = true
				for _, v := range hdata {
					fmt.Println("create stockkeeper P5", v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)
					s.logger.Infof("P5 create warehouse employee: %v, %v, %v, %v, %v, %v\n", v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)

					u := &model.User{
						Email:     v.Email,
						Password:  v.Password,
						FirstName: v.FirstName,
						LastName:  v.LastName,
						Role:      v.Role,
						Groups:    Groupp5,
						Tabel:     v.Tabel,
					}
					//	s.Lock()
					if err := s.store.User().CreateUserByManager(u); err != nil {
						s.error(w, r, http.StatusUnprocessableEntity, err)
						return
					}
					//	s.Unlock()
				}
			}
		}
	}
}

func (s *Server) PageupdateUserWarehouse() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		//user := r.Context().Value(ctxKeyUser).(*model.User)
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		idd, ok := session.Values["user_id"]
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		user, err := s.store.User().Find(idd.(int))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			//	s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		if user.Role == roleAdministrator {
			statusAdmin = true
			statusLoggedIn = true
		} else if user.Role == roleWarehouseManager {
			statusWarehouseManager = true
			statusLoggedIn = true
			fmt.Println("SuperIngenerQuality pageupdateUserQuality - ", statusWarehouseManager)
		}

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

		data := map[string]interface{}{
			"GET":              get,
			"Admin":            statusAdmin,
			"WarehouseManager": statusWarehouseManager,
			"LoggedIn":         statusLoggedIn,
			"User":             user.LastName,
			"Username":         user.FirstName,
		}

		//	fmt.Println("Get.email" - get.id)
		tpl.ExecuteTemplate(w, "updateuserwarehouse.html", data)
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
