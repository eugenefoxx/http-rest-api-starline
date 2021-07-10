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

func (s *Server) PageinsertVendor() http.HandlerFunc {
	/*tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertvendor.html")
	if err != nil {
		panic(err)
	}*/
	///tpl = template.Must(template.New("base").ParseFiles(s.html+"layout1.html", s.html+"insertvendor1.html"))
	return func(w http.ResponseWriter, r *http.Request) {

		//	var body, _ = helper.LoadFile("./web/templates/insertsapbyship6.html")
		//	fmt.Fprintf(w, body)

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupQuality {

			if user.Role == roleAdministrator {
				params := []Param{
					{
						LoggedIn: true,
						GroupP1:  true,
						Admin:    true,
						User:     user.LastName,
						Username: user.LastName,
					},
				}

				data := map[string]interface{}{
					"GetParam": params,
				}

				RenderTemplate(w, "insertvendor.html", data)

			} else if user.Role == roleSuperIngenerQuality {
				params := []Param{
					{
						LoggedIn:            true,
						GroupP1:             true,
						SuperIngenerQuality: true,
						User:                user.LastName,
						Username:            user.LastName,
					},
				}

				data := map[string]interface{}{
					"GetParam": params,
				}

				RenderTemplate(w, "insertvendor.html", data)
			}

		}

		if user.Groups == groupQualityP5 {

			if user.Role == roleAdministrator {
				params := []Param{
					{
						LoggedIn: true,
						GroupP5:  true,
						Admin:    true,
						User:     user.LastName,
						Username: user.LastName,
					},
				}

				data := map[string]interface{}{
					"GetParam": params,
				}

				RenderTemplate(w, "insertvendor.html", data)

			} else if user.Role == roleSuperIngenerQuality {
				params := []Param{
					{
						LoggedIn:            true,
						GroupP5:             true,
						SuperIngenerQuality: true,
						User:                user.LastName,
						Username:            user.LastName,
					},
				}

				data := map[string]interface{}{
					"GetParam": params,
				}

				RenderTemplate(w, "insertvendor.html", data)

			}

		}
	}
}

func (s *Server) InsertVendor() http.HandlerFunc {
	type requestFrom struct {
		CodeDebitor string `json:"code_debitor"`
		NameDebitor string `json:"name_debitor"`
	}

	/*tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertvendor.html")
	if err != nil {
		panic(err)
	}*/
	///tpl = template.Must(template.New("base").ParseFiles(s.html+"layout1.html", s.html+"insertvendor1.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		fmt.Println("Check2 -")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}

		var hdata []requestFrom
		json.Unmarshal(body, &hdata)
		fmt.Printf("body json: %s", body)
		s.logger.Infof("Loading body json: %s", body)
		fmt.Println("\njson  struct hdata", hdata)
		s.logger.Infof("Loading hdata json: %v", hdata)

		user := r.Context().Value(ctxKeyUser).(*model.User)

		for _, v := range hdata {
			fmt.Println(v.CodeDebitor, v.NameDebitor)
			s.logger.Infof("Create vendor: %s, %s", v.CodeDebitor, v.NameDebitor)
			u := &model.Vendor{
				CodeDebitor: v.CodeDebitor,
				NameDebitor: v.NameDebitor,
			}
			if user.Role == roleAdministrator {
				params := []Param{
					{
						LoggedIn: true,
						Admin:    true,
						User:     user.LastName,
						Username: user.FirstName,
					},
				}
				if err := s.store.Vendor().InsertVendor(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				data := map[string]interface{}{
					"GetParam": params,
				}

				RenderTemplate(w, "insertvendor.html", data)

			} else if user.Role == roleSuperIngenerQuality {
				params := []Param{
					{
						LoggedIn:            true,
						SuperIngenerQuality: true,
						User:                user.LastName,
						Username:            user.FirstName,
					},
				}
				if err := s.store.Vendor().InsertVendor(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
				data := map[string]interface{}{
					"GetParam": params,
				}

				RenderTemplate(w, "insertvendor.html", data)

			}

		}

	}

}

func (s *Server) PageVendor() http.HandlerFunc {
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showvendor.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//	var body, _ = helper.LoadFile("./web/templates/insertsapbyship6.html")
		//	fmt.Fprintf(w, body)

		user := r.Context().Value(ctxKeyUser).(*model.User)
		get, err := s.store.Vendor().ListVendor()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if user.Groups == groupQuality {

			if user.Role == roleAdministrator {
				params := []Param{
					{GroupP1: true, LoggedIn: true, Admin: true, User: user.LastName, Username: user.FirstName},
				}
				data := map[string]interface{}{
					"GetParam": params,
					"GET":      get,
				}

				RenderTemplate(w, "showvendor.html", data)

			} else if user.Role == roleSuperIngenerQuality {
				params := []Param{
					{GroupP1: true, LoggedIn: true, SuperIngenerQuality: true, User: user.LastName, Username: user.FirstName},
				}

				data := map[string]interface{}{
					"GetParam": params,
					"GET":      get,
				}

				RenderTemplate(w, "showvendor.html", data)

			}

			// send all the vendors as response
			//json.NewEncoder(w).Encode(get)
			//fmt.Println("json.NewEncoder(w).Encode(get)")
		}
		if user.Groups == groupQualityP5 {

			if user.Role == roleAdministrator {
				params := []Param{
					{GroupP5: true, LoggedIn: true, Admin: true, User: user.LastName, Username: user.FirstName},
				}
				data := map[string]interface{}{
					"GetParam": params,
					"GET":      get,
				}

				err = tpl.ExecuteTemplate(w, "showvendor.html", data)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}

			} else if user.Role == roleSuperIngenerQuality {
				params := []Param{
					{GroupP5: true, LoggedIn: true, SuperIngenerQuality: true, User: user.LastName, Username: user.FirstName},
				}

				data := map[string]interface{}{
					"GetParam": params,
					"GET":      get,
				}

				err = tpl.ExecuteTemplate(w, "showvendor.html", data)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}

			}
			// send all the vendors as response
			//json.NewEncoder(w).Encode(get)
			//fmt.Println("json.NewEncoder(w).Encode(get)")
		}

	}
}

func (s *Server) PageupdateVendor() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updatevendor.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		user := r.Context().Value(ctxKeyUser).(*model.User)

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}

		get, err := s.store.Vendor().EditVendor(id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if user.Role == roleAdministrator {
			params := []Param{
				{
					LoggedIn: true,
					Admin:    true,
					User:     user.FirstName,
					Username: user.LastName,
				},
			}

			data := map[string]interface{}{
				"GetParam": params,
				"GET":      get,
			}

			RenderTemplate(w, "updatevendor.html", data)

		} else if user.Role == roleSuperIngenerQuality {
			params := []Param{
				{
					LoggedIn:            true,
					SuperIngenerQuality: true,
					User:                user.LastName,
					Username:            user.FirstName,
				},
			}

			data := map[string]interface{}{
				"GetParam": params,
				"GET":      get,
			}

			RenderTemplate(w, "updatevendor.html", data)

		}
	}
}

func (s *Server) UpdateVendor() http.HandlerFunc {
	type request struct {
		ID          int    `json:"ID"`
		CodeDebitor string `json:"codedebitor"`
		NameDebitor string `json:"namedebitor"`
	}
	///	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updatevendor.html")
	///	if err != nil {
	///		panic(err)
	///	}
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
		req.CodeDebitor = r.FormValue("codedebitor")
		req.NameDebitor = r.FormValue("namedebitor")
		fmt.Println("ID - ", req.ID)
		s.logger.Infof("Update vendor: %v, %v, %v\n", req.ID, req.CodeDebitor, req.NameDebitor)
		u := &model.Vendor{
			ID:          req.ID,
			CodeDebitor: req.CodeDebitor,
			NameDebitor: req.NameDebitor,
		}

		if err := s.store.Vendor().UpdateVendor(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		http.Redirect(w, r, "/operation/showvendor", 303)

	}

}

func (s *Server) DeleteVendor() http.HandlerFunc {
	type request struct {
		ID          int    `json:"ID"`
		CodeDebitor string `json:"codedebitor"`
		NameDebitor string `json:"namedebitor"`
		SPPElement  string `json:"sppelement"`
	}
	///	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updatevendor.html")
	///	if err != nil {
	///		panic(err)
	///	}
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
		s.logger.Infof("Deleted vendor id: %v\n", req.ID)
		//	req.CodeDebitor = r.FormValue("codedebitor")
		//	req.NameDebitor = r.FormValue("namedebitor")
		//	req.SPPElement = r.FormValue("sppelement")

		u := &model.Vendor{
			ID: req.ID,
			//	CodeDebitor: req.CodeDebitor,
			//	NameDebitor: req.NameDebitor,
			//	SPPElement:  req.SPPElement,
		}

		if err := s.store.Vendor().DeleteVendor(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		http.Redirect(w, r, "/operation/showvendor", 303)
	}
}
