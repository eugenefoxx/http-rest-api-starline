package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"github.com/gorilla/mux"
)

func (s *Server) PageinInspection() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "ininspection.html")
	///	if err != nil {
	///		panic(err)
	///	}

	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
		StockkeeperWH := false
		WarehouseManager := false
		SuperIngenerQuality := false
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

		if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true
		} else if user.Role == "кладовщик склада" {
			StockkeeperWH = true
			LoggedIn = true
			fmt.Println("кладовщик склада - ", StockkeeperWH)
		} else if user.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true
		} else if user.Role == "старший кладовщик склада" {
			WarehouseManager = true
			LoggedIn = true
		}
		data := map[string]interface{}{
			"Admin":               Admin,
			"StockkeeperWH":       StockkeeperWH,
			"WarehouseManager":    WarehouseManager,
			"SuperIngenerQuality": SuperIngenerQuality,
			//	"GET":           get,
			"LoggedIn": LoggedIn,
			"User":     user.LastName,
			"Username": user.FirstName,
		}
		tpl.ExecuteTemplate(w, "ininspection.html", data)
		///	tpl.ExecuteTemplate(w, "layout", nil)
		///	err = tpl.ExecuteTemplate(w, "layout", nil)
		///	if err != nil {
		///		http.Error(w, err.Error(), 400)
		///		return
		///	}
	}
}

func (s *Server) InInspection() http.HandlerFunc {
	type req struct {
		ScanID         string `json:"scanid"`
		SAP            int
		Lot            string
		Roll           int
		Qty            int
		ProductionDate string
		NumberVendor   string
		Location       string
	}

	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "ininspection.html")
	///	if err != nil {
	///		panic(err)
	///	}

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		var rdata []req
		var slice []string
		rdata1 := []string{}
		json.Unmarshal(body, &rdata)
		json.Unmarshal(body, &rdata1)
		json.Unmarshal(body, &slice)
		fmt.Printf("test ininspection %s", body)
		fmt.Println("\nall of the rdata ininspection", rdata)
		//	rdata2 := removeDuplicates(rdata1)
		//	fmt.Print(rdata2)
		fmt.Printf("slice: %q\n", slice)

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

		const statusTransfer = "отгружено на ВК"

		for _, v := range rdata {
			idMaterial := v.ScanID[0:45]

			//	fmt.Println("Пропускаем:\n" + idMaterial + "\n")
			sapStr := v.ScanID[1:8]
			sap := v.SAP
			sap, err := strconv.Atoi(sapStr)
			if err != nil {
				fmt.Println(err)
			}
			idrollStr := v.ScanID[20:30]
			idrollIns := v.Roll
			idrollIns, err = strconv.Atoi(idrollStr)
			if err != nil {
				fmt.Println(err)
			}
			v.Lot = v.ScanID[9:19]
			qtyStr := v.ScanID[31:36]
			qtyIns := v.Qty
			qtyIns, err2 := strconv.Atoi(qtyStr)
			if err != nil {
				fmt.Println(err2)
			}
			v.ProductionDate = v.ScanID[37:45]
			v.NumberVendor = v.ScanID[9:15]
			fmt.Println("v.NumberVendor", v.NumberVendor)
			if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 45) {
				u := &model.Inspection{
					IdMaterial:     idMaterial,
					SAP:            sap,
					Lot:            v.Lot,
					IdRoll:         idrollIns,
					Qty:            qtyIns,
					ProductionDate: v.ProductionDate,
					NumberVendor:   v.NumberVendor,
					Location:       statusTransfer,
					Lastname:       user.LastName,
				}
				if err := s.store.Inspection().InInspection(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)

					return
				}
			} else {
				if (strings.Contains(v.ScanID[0:1], "P") == false) && (len(v.ScanID) != 45) {
					fmt.Println("не верное сканирование :\n" + v.ScanID + "\n")
					//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
				}
				//	tpl.Execute(w, data)
				return
			}
			//	} else {
			//		if idMaterial == idMaterial {
			//			fmt.Println("Значения совпадают:\n" + idMaterial + "\n")
			//		}
			//	}
		}

		err = tpl.ExecuteTemplate(w, "ininspection.html", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

	}
}

//historyInspection
func (s *Server) PagehistoryInspection() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
		StockkeeperWH := false
		WarehouseManager := false
		SuperIngenerQuality := false
		IngenerQuality := false
		Quality := false
		Inspector := false
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

		if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true
		} else if user.Role == "кладовщик склада" {
			StockkeeperWH = true
			LoggedIn = true
			fmt.Println("кладовщик склада - ", StockkeeperWH)
		} else if user.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true
		} else if user.Role == "инженер по качеству" {
			IngenerQuality = true
			LoggedIn = true
		} else if user.Groups == "качество" {
			Quality = true
			Inspector = true
			LoggedIn = true
			fmt.Println("pageInspection quality - ", Quality)
		} else if user.Role == "старший кладовщик склада" {
			WarehouseManager = true
			LoggedIn = true
		}
		data := map[string]interface{}{
			"Admin":               Admin,
			"StockkeeperWH":       StockkeeperWH,
			"SuperIngenerQuality": SuperIngenerQuality,
			"WarehouseManager":    WarehouseManager,
			"IngenerQuality":      IngenerQuality,
			"Quality":             Quality,
			"Inspector":           Inspector,
			//	"GET":           get,
			"LoggedIn": LoggedIn,
			"User":     user.LastName,
			"Username": user.FirstName,
		}

		tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
	}
}

func (s *Server) HistoryInspection() http.HandlerFunc {
	type req struct {
		Date1    string `json:"date1"`
		Date2    string `json:"date2"`
		Material int    `json:"material"`
		EO       string `json:"eo"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
		StockkeeperWH := false
		WarehouseManager := false
		SuperIngenerQuality := false
		IngenerQuality := false
		Quality := false
		Inspector := false
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

		if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true
		} else if user.Role == "кладовщик склада" {
			StockkeeperWH = true
			LoggedIn = true
			fmt.Println("кладовщик склада - ", StockkeeperWH)
		} else if user.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true
		} else if user.Role == "инженер по качеству" {
			IngenerQuality = true
			LoggedIn = true
		} else if user.Groups == "качество" {
			Quality = true
			Inspector = true
			LoggedIn = true
			fmt.Println("pageInspection quality - ", Quality)
		} else if user.Role == "старший кладовщик склада" {
			WarehouseManager = true
			LoggedIn = true
		}

		search := &req{}
		materialInt, err := strconv.Atoi(r.FormValue("material"))
		if err != nil {
			fmt.Println(err)
		}
		search.Date1 = r.FormValue("date1")
		fmt.Println("date1 - ", search.Date1)
		search.Date2 = r.FormValue("date2")
		fmt.Println("date2 - ", search.Date2)
		search.Material = materialInt
		fmt.Println("material - ", search.Material)
		search.EO = r.FormValue("eo")

		if search.Date1 == "" && search.Date2 == "" {
			if search.EO != "" {
				get, err := s.store.Inspection().ListShowDataByEO(search.EO)
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
				data := map[string]interface{}{
					"TitleDOC":            "Отчет по истроии ВК",
					"User":                user.LastName,
					"Username":            user.FirstName,
					"Admin":               Admin,
					"WarehouseManager":    WarehouseManager,
					"StockkeeperWH":       StockkeeperWH,
					"SuperIngenerQuality": SuperIngenerQuality,
					"IngenerQuality":      IngenerQuality,
					"Quality":             Quality,
					"Inspector":           Inspector,
					"LoggedIn":            LoggedIn,
					"GET":                 get,
				}

				err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}
			} else if search.Material != 0 {
				get, err := s.store.Inspection().ListShowDataBySap(search.Material)
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
				data := map[string]interface{}{
					"TitleDOC":            "Отчет по истроии ВК",
					"User":                user.LastName,
					"Username":            user.FirstName,
					"Admin":               Admin,
					"WarehouseManager":    WarehouseManager,
					"StockkeeperWH":       StockkeeperWH,
					"SuperIngenerQuality": SuperIngenerQuality,
					"IngenerQuality":      IngenerQuality,
					"Quality":             Quality,
					"Inspector":           Inspector,
					"LoggedIn":            LoggedIn,
					"GET":                 get,
				}

				err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}
			}
		} else {

			if search.Material != 0 {
				fmt.Println("OK Material")
				get, err := s.store.Inspection().ListShowDataByDateAndSAP(search.Date1, search.Date2, search.Material)
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
				data := map[string]interface{}{
					"TitleDOC":            "Отчет по истроии ВК",
					"User":                user.LastName,
					"Username":            user.FirstName,
					"Admin":               Admin,
					"WarehouseManager":    WarehouseManager,
					"StockkeeperWH":       StockkeeperWH,
					"SuperIngenerQuality": SuperIngenerQuality,
					"IngenerQuality":      IngenerQuality,
					"Quality":             Quality,
					"Inspector":           Inspector,
					"LoggedIn":            LoggedIn,
					"GET":                 get,
				}

				err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}
			} else if search.EO != "" {
				fmt.Println("OK EO")

				get, err := s.store.Inspection().ListShowDataByDateAndEO(search.Date1, search.Date2, search.EO)
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
				data := map[string]interface{}{
					"TitleDOC":            "Отчет по истроии ВК",
					"User":                user.LastName,
					"Username":            user.FirstName,
					"Admin":               Admin,
					"WarehouseManager":    WarehouseManager,
					"StockkeeperWH":       StockkeeperWH,
					"SuperIngenerQuality": SuperIngenerQuality,
					"IngenerQuality":      IngenerQuality,
					"Quality":             Quality,
					"Inspector":           Inspector,
					"LoggedIn":            LoggedIn,
					"GET":                 get,
				}

				err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}

			} else {
				get, err := s.store.Inspection().ListShowDataByDate(search.Date1, search.Date2)
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				data := map[string]interface{}{
					"TitleDOC":            "Отчет по истроии ВК",
					"User":                user.LastName,
					"Username":            user.FirstName,
					"Admin":               Admin,
					"WarehouseManager":    WarehouseManager,
					"StockkeeperWH":       StockkeeperWH,
					"SuperIngenerQuality": SuperIngenerQuality,
					"IngenerQuality":      IngenerQuality,
					"Quality":             Quality,
					"Inspector":           Inspector,
					"LoggedIn":            LoggedIn,
					"GET":                 get,
				}

				err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}
			}
		}
		/*
				data := map[string]interface{}{
					"TitleDOC":            "Отчет по истроии ВК",
					"User":                user.LastName,
					"Username":            user.FirstName,
					"Admin":               Admin,
					"StockkeeperWH":       StockkeeperWH,
					"главный инженер по качеству": SuperIngenerQuality,
					"Quality":             Quality,
					"LoggedIn":            LoggedIn,
					"GET":                 get,
				}

			err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", nil)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		*/
	}
}

func (s *Server) PageInspection() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	Admin := true
		//	Warehouse := true
		StockkeeperWH := false
		WarehouseManager := false
		//	Quality := false
		Inspector := false
		SuperIngenerQuality := false
		IngenerQuality := false
		SuperIngenerQuality2 := false
		LoggedIn := true

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

		if user.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			SuperIngenerQuality2 = true
			LoggedIn = true
			fmt.Println("pageInspection SuperIngenerQuality - ", SuperIngenerQuality)
		} else if user.Role == "кладовщик склада" {
			StockkeeperWH = true
			//	Warehouse = false
			//	WarehouseManager = true
			LoggedIn = true
		} else if user.Role == "инженер по качеству" {
			IngenerQuality = true
		} else if user.Role == "контролер качества" {
			Inspector = true
		} else if user.Role == "старший кладовщик склада" {
			WarehouseManager = true
			LoggedIn = true
		} /* else if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true
		}*/ /**else if user.Groups == "качество" {
			//	Quality = true
			Inspector = true
			IngenerQuality = true
			LoggedIn = true
			//	fmt.Println("pageInspection quality - ", Quality)
		} */

		get, err := s.store.Inspection().ListInspection()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		countTotal, err := s.store.Inspection().CountTotalInspection()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		holdInspection, err := s.store.Inspection().HoldInspection()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		notVerifyComponents, err := s.store.Inspection().NotVerifyComponents()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		getStatic, err := s.store.Inspection().CountDebitor()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		holdCountDebitor, err := s.store.Inspection().HoldCountDebitor()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		notVerifyDebitor, err := s.store.Inspection().NotVerifyDebitor()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		groups := map[string]interface{}{
			"TitleDOC": "Список ВК",
			"User":     user.LastName,
			"Username": user.FirstName,
			//	"Admin":                Admin,
			//	"Quality":              Quality,
			"Inspector":        Inspector,
			"IngenerQuality":   IngenerQuality,
			"WarehouseManager": WarehouseManager,
			//	"Warehouse":            Warehouse,
			"StockkeeperWH":        StockkeeperWH,
			"SuperIngenerQuality":  SuperIngenerQuality,
			"SuperIngenerQuality2": SuperIngenerQuality2,
			"GET":                  get,
			"CountTotal":           countTotal,
			"HoldInspection":       holdInspection,
			"NotVerifyComponents":  notVerifyComponents,
			"GetStatic":            getStatic,
			"HoldCountDebitor":     holdCountDebitor,
			"NotVerifyDebitor":     notVerifyDebitor,
			"LoggedIn":             LoggedIn,
		}

		tpl.ExecuteTemplate(w, "showinspection.html", groups)

	}
}

func (s *Server) PageupdateInspection() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
		SuperIngenerQuality := false
		IngenerQuality := false
		Inspector := false
		LoggedIn := false

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}

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
		fmt.Println("user.Groups - ?", user.Groups)

		if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true
		} else if user.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true
			fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
		} else if user.Role == "инженер по качеству" {
			IngenerQuality = true
			LoggedIn = true
			fmt.Println("IngenerQuality - ", IngenerQuality)
		} else if user.Role == "контролер качества" {
			Inspector = true
			LoggedIn = true

		}
		//fmt.Println("ID - ?", id)
		get, err := s.store.Inspection().EditInspection(id, user.Groups)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		data := map[string]interface{}{
			"Admin":               Admin,
			"SuperIngenerQuality": SuperIngenerQuality,
			"IngenerQuality":      IngenerQuality,
			"Inspector":           Inspector,
			"GET":                 get,
			"LoggedIn":            LoggedIn,
			"User":                user.LastName,
			"Username":            user.FirstName,
		}
		err = tpl.ExecuteTemplate(w, "updateinspection.html", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *Server) UpdateInspection() http.HandlerFunc {
	type request struct {
		ID     int    `json:"ID"`
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	///	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}

		currentTime := time.Now()

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

		req.ID = id
		req.Status = r.FormValue("status")
		req.Note = r.FormValue("note")

		u := &model.Inspection{
			ID:         req.ID,
			Status:     req.Status,
			Note:       req.Note,
			Update:     user.LastName, //
			Dateupdate: currentTime,   // Dateaccept
			Timeupdate: currentTime,   // Timeaccept
			Groups:     user.Groups,
		}

		if err := s.store.Inspection().UpdateInspection(u, user.Groups); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		/*	err = tpl.ExecuteTemplate(w, "layout", nil)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}*/
		http.Redirect(w, r, "/statusinspection", 303)
	}
}

func (s *Server) DeleteInspection() http.HandlerFunc {
	type request struct {
		ID int `json:"ID"`
	}
	///	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}

		req.ID = id

		u := &model.Inspection{
			ID: req.ID,
		}

		if err := s.store.Inspection().DeleteItemInspection(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		http.Redirect(w, r, "/statusinspection", 303)
	}
}

//ListAcceptWHInspection
func (s *Server) PageListAcceptWHInspection() http.HandlerFunc { // acceptinspection.html showinspection.html
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "acceptinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
		StockkeeperWH := false
		WarehouseManager := false
		LoggedIn := false
		stockkeeperWH := false
		//superIngenerQuality := true
		//quality := false
		//	statusStr := false

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

		if user.Groups == "склад" {
			stockkeeperWH = true
			WarehouseManager = true
		}
		/*
			if user.Groups == "качество" {
				quality = true
			} else if user.Groups == "склад" {
				stockkeeperWH = true
			}
		*/
		if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true
		} else if user.Role == "кладовщик склада" {
			StockkeeperWH = true
			LoggedIn = true
			fmt.Println("кладовщик склада - ", StockkeeperWH)
		} else if user.Role == "старший кладовщик склада" {
			WarehouseManager = true
			LoggedIn = true
		}
		get, err := s.store.Inspection().ListAcceptWHInspection()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		groups := map[string]interface{}{
			//	"quality":   quality,
			"Warehouse":        stockkeeperWH,
			"WarehouseManager": WarehouseManager,
			//	"главный инженер по качеству": superIngenerQuality,
			"GET": get,
			//	"status":    statusStr,
			"Admin":         Admin,
			"StockkeeperWH": StockkeeperWH,
			"LoggedIn":      LoggedIn,
			"User":          user.LastName,
			"Username":      user.FirstName,
		}

		err = tpl.ExecuteTemplate(w, "acceptinspection.html", groups)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *Server) PageacceptWarehouseInspection() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "acceptWarehouseInspection.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
		StockkeeperWH := false
		WarehouseManager := false
		LoggedIn := false

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}

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
		fmt.Println("user.Groups - ?", user.Groups)

		if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true

		} else if user.Role == "кладовщик склада" {
			StockkeeperWH = true
			LoggedIn = true

		} else if user.Role == "старший кладовщик склада" {
			WarehouseManager = true
			LoggedIn = true
		}

		//fmt.Println("ID - ?", id)
		get, err := s.store.Inspection().EditAcceptWarehouseInspection(id, user.Groups)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		data := map[string]interface{}{
			"User":             user.LastName,
			"Username":         user.FirstName,
			"Admin":            Admin,
			"WarehouseManager": WarehouseManager,
			"StockkeeperWH":    StockkeeperWH,
			"LoggedIn":         LoggedIn,
			"GET":              get,
		}
		err = tpl.ExecuteTemplate(w, "acceptWarehouseInspection.html", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *Server) AcceptWarehouseInspection() http.HandlerFunc {
	type request struct {
		ID       int    `json:"ID"`
		Location string `json:"location"`
	}
	///	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "acceptWarehouseInspection.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}

		currentTime := time.Now()

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

		req.ID = id
		req.Location = r.FormValue("location")

		u := &model.Inspection{
			ID:             req.ID,
			Location:       req.Location,
			Lastnameaccept: user.LastName, // Lastnameaccept
			Dateaccept:     currentTime,   // Dateaccept
			Timeaccept:     currentTime,   // Timeaccept
			Groups:         user.Groups,
		}

		if err := s.store.Inspection().AcceptWarehouseInspection(u, user.Groups); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		/*	err = tpl.ExecuteTemplate(w, "layout", nil)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}*/
		//http.Redirect(w, r, "/statusinspectionforwh", 303)
		http.Redirect(w, r, "/statusinspection", 303)
	}
}
