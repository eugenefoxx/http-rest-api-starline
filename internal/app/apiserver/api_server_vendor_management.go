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
		//data := map[string]interface{}{
		//	"user": "Я тут",
		//}
		/*err = tpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return*/
		//	}
		//		if err = tpl.ExecuteTemplate(w, "layout", nil); err != nil {
		//			s.error(w, r, http.StatusUnprocessableEntity, err)
		//			return
		//		}
		Admin := false
		SuperIngenerQuality := false
		GroupP1 := false
		GroupP5 := false
		LoggedIn := false

		user := r.Context().Value(ctxKeyUser).(*model.User)
		//	GET := map[string]bool{
		//		"admin": admin,
		//		//	"stockkeeper":         stockkeeper,
		//		roleSuperIngenerQuality: superIngenerQuality,
		//	"stockkeeperWH":       stockkeeperWH,
		//	"inspector":           inspector,
		//	}
		if user.Groups == groupQuality {
			GroupP1 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
				fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
			}
			data := map[string]interface{}{
				"Admin":               Admin,
				"SuperIngenerQuality": SuperIngenerQuality,
				"GroupP1":             GroupP1,
				//"GET":      GET,
				"LoggedIn": LoggedIn,
				"User":     user.LastName,
				"Username": user.FirstName,
			}
			fmt.Println("Check -")
			err := tpl.ExecuteTemplate(w, "insertvendor.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}

		if user.Groups == groupQualityP5 {
			GroupP5 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
				fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
			}
			data := map[string]interface{}{
				"Admin":               Admin,
				"SuperIngenerQuality": SuperIngenerQuality,
				"GroupP5":             GroupP5,
				//"GET":      GET,
				"LoggedIn": LoggedIn,
				"User":     user.LastName,
				"Username": user.FirstName,
			}
			fmt.Println("Check -")
			err := tpl.ExecuteTemplate(w, "insertvendor.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
	}
}

func (s *Server) InsertVendor() http.HandlerFunc {
	type requestFrom struct {
		CodeDebitor string `json:"code_debitor"`
		NameDebitor string `json:"name_debitor"`
	}
	/*
		type requestDB struct {
			CodeDebitor string `db:"code_debitor"`
			NameDebitor string `db:"name_debitor"`
			SPPElement  string `db:"spp_element"`
		}
	*/
	/*tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertvendor.html")
	if err != nil {
		panic(err)
	}*/
	///tpl = template.Must(template.New("base").ParseFiles(s.html+"layout1.html", s.html+"insertvendor1.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		fmt.Println("Check2 -")
		Admin := false
		SuperIngenerQuality := false
		LoggedIn := false

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

		//	GET := map[string]bool{
		//		"admin": admin,
		//	"stockkeeper":         stockkeeper,
		//		roleSuperIngenerQuality: superIngenerQuality,
		//	"stockkeeperWH":       stockkeeperWH,
		//	"inspector":           inspector,
		//	}

		if user.Role == roleAdministrator {
			Admin = true
			LoggedIn = true
		} else if user.Role == roleSuperIngenerQuality {
			SuperIngenerQuality = true
			LoggedIn = true
			fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
		}

		for _, v := range hdata {
			fmt.Println(v.CodeDebitor, v.NameDebitor)
			s.logger.Infof("Create vendor: %s, %s", v.CodeDebitor, v.NameDebitor)
			u := &model.Vendor{
				CodeDebitor: v.CodeDebitor,
				NameDebitor: v.NameDebitor,
			}

			if err := s.store.Vendor().InsertVendor(u); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

		}

		data := map[string]interface{}{
			"Admin":               Admin,
			"SuperIngenerQuality": SuperIngenerQuality,
			//"GET":      GET,
			"LoggedIn": LoggedIn,
			"User":     user.LastName,
			"Username": user.FirstName,
		}

		//err = tpl.ExecuteTemplate(w, "layout", data)
		//	err = tpl.ExecuteTemplate(w, "layout", v)
		err = tpl.ExecuteTemplate(w, "insertvendor.html", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
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
		//data := map[string]interface{}{
		//	"user": "Я тут",
		//}
		Admin := false
		SuperIngenerQuality := false
		GroupP1 := false
		GroupP5 := false
		LoggedIn := false

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupQuality {
			GroupP1 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
				fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
			}

			get, err := s.store.Vendor().ListVendor()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			data := map[string]interface{}{
				"Admin":               Admin,
				"SuperIngenerQuality": SuperIngenerQuality,
				"GroupP1":             GroupP1,
				"GET":                 get,
				"LoggedIn":            LoggedIn,
				"User":                user.LastName,
				"Username":            user.FirstName,
			}

			err = tpl.ExecuteTemplate(w, "showvendor.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			// send all the vendors as response
			//json.NewEncoder(w).Encode(get)
			//fmt.Println("json.NewEncoder(w).Encode(get)")
		}
		if user.Groups == groupQualityP5 {
			GroupP5 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
				fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
			}

			get, err := s.store.Vendor().ListVendor()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			data := map[string]interface{}{
				"Admin":               Admin,
				"SuperIngenerQuality": SuperIngenerQuality,
				"GroupP5":             GroupP5,
				"GET":                 get,
				"LoggedIn":            LoggedIn,
				"User":                user.LastName,
				"Username":            user.FirstName,
			}

			err = tpl.ExecuteTemplate(w, "showvendor.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
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

		Admin := false
		SuperIngenerQuality := false
		LoggedIn := false

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Role == roleAdministrator {
			Admin = true
			LoggedIn = true
		} else if user.Role == roleSuperIngenerQuality {
			SuperIngenerQuality = true
			LoggedIn = true
			fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}
		//fmt.Println("ID - ?", id)

		get, err := s.store.Vendor().EditVendor(id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		data := map[string]interface{}{
			"Admin":               Admin,
			"SuperIngenerQuality": SuperIngenerQuality,
			"GET":                 get,
			"LoggedIn":            LoggedIn,
			"User":                user.LastName,
			"Username":            user.FirstName,
		}
		err = tpl.ExecuteTemplate(w, "updatevendor.html", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
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

		/*	err = tpl.ExecuteTemplate(w, "layout", nil)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}*/
		http.Redirect(w, r, "/operation/showvendor", 303)
	}
}
