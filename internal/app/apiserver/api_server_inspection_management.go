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

		if user.Role == roleAdministrator {
			Admin = true
			LoggedIn = true
		} else if user.Role == roleStockkeeperWH {
			StockkeeperWH = true
			LoggedIn = true
			fmt.Println("кладовщик склада - ", StockkeeperWH)
		} else if user.Role == roleSuperIngenerQuality {
			SuperIngenerQuality = true
			LoggedIn = true
		} else if user.Role == roleWarehouseManager {
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
		s.mu.Lock()
		defer s.mu.Unlock()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
			
		}

		var rdata []req
		var slice []string
		rdata1 := []string{}
		json.Unmarshal(body, &rdata)
		json.Unmarshal(body, &rdata1)
		json.Unmarshal(body, &slice)
		fmt.Printf("test ininspection %s", body)
		//s.infoLog.Printf("Loading json body %s", body)
		s.logger.Infof("Loading json body %s", body)
		fmt.Println("\nall of the rdata ininspection", rdata)
		//s.infoLog.Printf("Loading json rdata %v", rdata)
		s.logger.Infof("Loading json rdata %v", rdata)
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

		if user.Groups == groupWarehouse || user.Groups == groupQuality {
			for _, v := range rdata {
				if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 45) {
					idMaterial := v.ScanID[0:45]

					//	fmt.Println("Пропускаем:\n" + idMaterial + "\n")
					//s.infoLog.Printf("Запись строки сканирования на входной контроль, П1: %v", idMaterial)
					s.logger.Infof("Запись строки сканирования на входной контроль, П1: %v", idMaterial)
					sapStr := v.ScanID[1:8]
					sap := v.SAP
					sap, err := strconv.Atoi(sapStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					idrollStr := v.ScanID[20:30]
					idrollIns := v.Roll
					idrollIns, err = strconv.Atoi(idrollStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					v.Lot = v.ScanID[9:19]
					qtyStr := v.ScanID[31:36]
					qtyIns := v.Qty
					qtyIns, err = strconv.Atoi(qtyStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
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
							s.logger.Errorf("Не верное сканирование на входной контроль, П1: %v", v)
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
					http.Redirect(w, r, "/operation/statusinspection", 303)
				}
				if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 35) {
					idMaterial := v.ScanID[0:35]

					//	fmt.Println("Пропускаем:\n" + idMaterial + "\n")
					//s.infoLog.Printf("Запись строки сканирования на входной контроль, П1: %v", idMaterial)
					s.logger.Infof("Запись строки сканирования на входной контроль, П1: %v", idMaterial)
					sapStr := v.ScanID[1:8]
					sap := v.SAP
					sap, err := strconv.Atoi(sapStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					idrollStr := v.ScanID[10:20]
					idrollIns := v.Roll
					idrollIns, err = strconv.Atoi(idrollStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					v.Lot = "без партии" //v.ScanID[9:19]
					qtyStr := v.ScanID[21:26]
					qtyIns := v.Qty
					qtyIns, err = strconv.Atoi(qtyStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					v.ProductionDate = v.ScanID[27:35]
					v.NumberVendor = "без поставщика" //v.ScanID[9:15]
					fmt.Println("v.NumberVendor", v.NumberVendor)
					if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 35) {
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
						if (strings.Contains(v.ScanID[0:1], "P") == false) && (len(v.ScanID) != 35) {
							fmt.Println("не верное сканирование :\n" + v.ScanID + "\n")
							s.logger.Errorf("Не верное сканирование на входной контроль, П1: %v", v)
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
					http.Redirect(w, r, "/operation/statusinspection", 303)
				}
			}
		}

		if user.Groups == groupWarehouseP5 || user.Groups == groupQualityP5 {
			for _, v := range rdata {
				if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 45) {
					idMaterial := v.ScanID[0:45]

					//	fmt.Println("Пропускаем:\n" + idMaterial + "\n")
					//s.infoLog.Printf("Запись строки сканирования на входной контроль, П5: %v", idMaterial)
					s.logger.Infof("Запись строки сканирования на входной контроль, П5: %v", idMaterial)
					sapStr := v.ScanID[1:8]
					sap := v.SAP
					sap, err := strconv.Atoi(sapStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					idrollStr := v.ScanID[20:30]
					idrollIns := v.Roll
					idrollIns, err = strconv.Atoi(idrollStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					v.Lot = v.ScanID[9:19]
					qtyStr := v.ScanID[31:36]
					qtyIns := v.Qty
					qtyIns, err = strconv.Atoi(qtyStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
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

						if err := s.store.Inspection().InInspectionP5(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)

							return
						}

					} else {
						if (strings.Contains(v.ScanID[0:1], "P") == false) && (len(v.ScanID) != 45) {
							fmt.Println("не верное сканирование :\n" + v.ScanID + "\n")
							s.logger.Errorf("Не верное сканирование на входной контроль, П5: %v", v)
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
					http.Redirect(w, r, "/operation/statusinspection", 303)
				}
				if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 35) {
					idMaterial := v.ScanID[0:35]

					//	fmt.Println("Пропускаем:\n" + idMaterial + "\n")
					//s.infoLog.Printf("Запись строки сканирования на входной контроль, П5: %v", idMaterial)
					s.logger.Infof("Запись строки сканирования на входной контроль, П5: %v", idMaterial)
					sapStr := v.ScanID[1:8]
					sap := v.SAP
					sap, err := strconv.Atoi(sapStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					idrollStr := v.ScanID[10:20]
					idrollIns := v.Roll
					idrollIns, err = strconv.Atoi(idrollStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					v.Lot = "без партии" //v.ScanID[9:19]
					qtyStr := v.ScanID[21:26]
					qtyIns := v.Qty
					qtyIns, err = strconv.Atoi(qtyStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					v.ProductionDate = v.ScanID[27:35]
					v.NumberVendor = "без поставщика" //v.ScanID[9:15]
					fmt.Println("v.NumberVendor", v.NumberVendor)
					if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 35) {
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

						if err := s.store.Inspection().InInspectionP5(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)

							return
						}

					} else {
						if (strings.Contains(v.ScanID[0:1], "P") == false) && (len(v.ScanID) != 35) {
							fmt.Println("не верное сканирование :\n" + v.ScanID + "\n")
							s.logger.Errorf("Не верное сканирование на входной контроль, П5: %v", v)
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
					http.Redirect(w, r, "/operation/statusinspection", 303)
				}
			}
		}
		/*
			err = tpl.ExecuteTemplate(w, "ininspection.html", nil)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}*/

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

		if user.Groups == groupQuality || user.Groups == groupWarehouse || user.Groups == groupQualityP5 || user.Groups == groupWarehouseP5 {
			GroupP1 = true
			GroupP5 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				LoggedIn = true
				fmt.Println("кладовщик склада - ", StockkeeperWH)
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
				LoggedIn = true
			} else if user.Groups == groupQuality {
				Quality = true
				Inspector = true
				LoggedIn = true
				fmt.Println("pageInspection quality - ", Quality)
			} else if user.Role == roleWarehouseManager {
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
				"GroupP1":             GroupP1,
				"GroupP5":             GroupP5,
				//	"GET":           get,
				"LoggedIn": LoggedIn,
				"User":     user.LastName,
				"Username": user.FirstName,
			}

			tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
		}

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
		s.mu.Lock()
		defer s.mu.Unlock()

		Admin := false
		StockkeeperWH := false
		WarehouseManager := false
		SuperIngenerQuality := false
		IngenerQuality := false
		Quality := false
		Inspector := false
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

		if user.Groups == groupQuality || user.Groups == groupWarehouse {
			GroupP1 = true

			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				LoggedIn = true
				fmt.Println("кладовщик склада - ", StockkeeperWH)
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
				LoggedIn = true
			} else if user.Role == roleInspector {
				Quality = true
				Inspector = true
				LoggedIn = true
				fmt.Println("pageInspection quality - ", Quality)
			} else if user.Role == roleWarehouseManager {
				WarehouseManager = true
				LoggedIn = true
			}

			search := &req{}
			materialInt, err := strconv.Atoi(r.FormValue("material"))
			if err != nil {
				fmt.Println(err)
				s.logger.Errorf(err.Error())
			}
			search.Date1 = r.FormValue("date1")
			fmt.Println("date1 - ", search.Date1)
			//s.infoLog.Printf("date1 - %v", search.Date1)
			s.logger.Infof("date1 - %v", search.Date1)
			search.Date2 = r.FormValue("date2")
			fmt.Println("date2 - ", search.Date2)
			//s.infoLog.Printf("date2 - %v", search.Date2)
			s.logger.Infof("date2 - %v", search.Date2)
			search.Material = materialInt
			fmt.Println("material - ", search.Material)
			//s.infoLog.Printf("material - %v", search.Material)
			s.logger.Infof("material - %v", search.Material)
			search.EO = r.FormValue("eo")
			//s.infoLog.Printf("eo - %v", search.EO)
			s.logger.Infof("eo - %v", search.EO)

			currentData := time.Now()
			searchDateNow := currentData.Format("2006-01-02")

			if search.Date1 == "" && search.Date2 == "" && search.Material == 0 && search.EO == "" {
				fmt.Println("Не заполнены поля ввода")
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
					"GroupP1":             GroupP1,
					//"GroupP5":             GroupP5,
					"LoggedIn": LoggedIn,
					//	"GET":                 get,
				}

				err = tpl.ExecuteTemplate(w, "errorSearchHistoryInspection.html", data)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}
			} else if search.Date1 == "" && search.Date2 == "" {
				if search.EO != "" {

					val, err := s.redis.Inspection().GetListShowDataByEO(r.Context(), search.EO)
					if err != nil {
						//RenderJSON(w, &val, http.StatusOK)
						fmt.Println(&val)
					}

					get, err := s.store.Inspection().ListShowDataByEO(search.EO)
					//get, err := s.store.Inspection().ListShowDataByEO(val)
					if err != nil {
						s.error(w, r, http.StatusUnprocessableEntity, err)
						return
					}

					_ = s.redis.Inspection().SetListShowDataByEO(r.Context(), get)
					/*
						count, _ := s.store.Inspection().CountInspection()
						fmt.Println(count)
						limit := 5
						page, begin := s.Pagination(r, limit)
						fmt.Printf("Current Page: %d, Begin: %d\n", page, begin)
					*/
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 1",
						//"GroupP5":             GroupP5,
						"GET": get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				} else if search.Material != 0 {

					count, _ := s.store.Inspection().CountInspection()
					fmt.Println(count)
					limit := 2
					page, begin := s.Pagination(r, limit)
					fmt.Printf("Current Page: %d, Begin: %d\n", page, begin)
					get, err := s.store.Inspection().ListShowDataBySap(search.Material, begin, limit)

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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 1",
						//	"GroupP5":             GroupP5,
						"GET": get,
					}
					// RenderJSON(w, get, http.StatusOK)
					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				}
			} else if search.Date1 != "" && search.Date2 == "" {

				if search.Material != 0 {
					fmt.Println("OK Material")

					get, err := s.store.Inspection().ListShowDataByDateAndSAP(search.Date1, searchDateNow, search.Material)
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 1",
						//	"GroupP5":             GroupP5,
						"GET": get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				} else if search.EO != "" {
					fmt.Println("OK EO")

					get, err := s.store.Inspection().ListShowDataByDateAndEO(search.Date1, searchDateNow, search.EO)
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 1",
						//	"GroupP5":             GroupP5,
						"GET": get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}

				} else {

					get, err := s.store.Inspection().ListShowDataByDate(search.Date1, searchDateNow)
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 1",
						//	"GroupP5":             GroupP5,
						"GET": get,
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 1",
						//	"GroupP5":             GroupP5,
						"GET": get,
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 1",
						//	"GroupP5":             GroupP5,
						"GET": get,
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 1",
						//	"GroupP5":             GroupP5,
						"GET": get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				}
			}
		}

		if user.Groups == groupQualityP5 || user.Groups == groupWarehouseP5 {

			GroupP5 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				LoggedIn = true
				fmt.Println("кладовщик склада - ", StockkeeperWH)
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
				LoggedIn = true
			} else if user.Role == roleInspector {
				Quality = true
				Inspector = true
				LoggedIn = true
				fmt.Println("pageInspection quality - ", Quality)
			} else if user.Role == roleWarehouseManager {
				WarehouseManager = true
				LoggedIn = true
			}

			search := &req{}
			materialInt, err := strconv.Atoi(r.FormValue("material"))
			if err != nil {
				fmt.Println(err)
				s.logger.Errorf(err.Error())
			}
			search.Date1 = r.FormValue("date1")
			fmt.Println("date1 - ", search.Date1)
			//s.infoLog.Printf("date1 - %v", search.Date1)
			s.logger.Infof("date1 - %v", search.Date1)
			search.Date2 = r.FormValue("date2")
			fmt.Println("date2 - ", search.Date2)
			//s.infoLog.Printf("date2 - %v", search.Date2)
			s.logger.Infof("date2 - %v", search.Date2)
			search.Material = materialInt
			fmt.Println("material - ", search.Material)
			//s.infoLog.Printf("material - %v", search.Material)
			s.logger.Infof("material - %v", search.Material)
			search.EO = r.FormValue("eo")
			//s.infoLog.Printf("eo - %v", search.EO)
			s.logger.Infof("eo - %v", search.EO)

			currentData := time.Now()
			searchDateNow := currentData.Format("2006-01-02")

			if search.Date1 == "" && search.Date2 == "" && search.Material == 0 && search.EO == "" {
				fmt.Println("Не заполнены поля ввода")
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
					//	"GroupP1":             GroupP1,
					"GroupP5":  GroupP5,
					"LoggedIn": LoggedIn,
					//	"GET":                 get,
				}

				err = tpl.ExecuteTemplate(w, "errorSearchHistoryInspection.html", data)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}
			} else if search.Date1 == "" && search.Date2 == "" {
				if search.EO != "" {

					get, err := s.store.Inspection().ListShowDataByEO(search.EO)
					if err != nil {
						s.error(w, r, http.StatusUnprocessableEntity, err)
						return
					}
					count, _ := s.store.Inspection().CountInspection()
					fmt.Println(count)
					limit := 5
					page, begin := s.Pagination(r, limit)
					fmt.Printf("Current Page: %d, Begin: %d\n", page, begin)

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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 1",
						"GET":     get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				} else if search.Material != 0 {

					count, _ := s.store.Inspection().CountInspection()
					fmt.Println(count)
					limit := 2
					page, begin := s.Pagination(r, limit)
					fmt.Printf("Current Page: %d, Begin: %d\n", page, begin)
					get, err := s.store.Inspection().ListShowDataBySap(search.Material, begin, limit)

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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 1",
						"GET":     get,
					}
					// RenderJSON(w, get, http.StatusOK)
					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				}
			} else if search.Date1 != "" && search.Date2 == "" {

				if search.Material != 0 {
					fmt.Println("OK Material")

					get, err := s.store.Inspection().ListShowDataByDateAndSAP(search.Date1, searchDateNow, search.Material)
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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 1",
						"GET":     get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				} else if search.EO != "" {
					fmt.Println("OK EO")

					get, err := s.store.Inspection().ListShowDataByDateAndEO(search.Date1, searchDateNow, search.EO)
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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 1",
						"GET":     get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}

				} else {

					get, err := s.store.Inspection().ListShowDataByDate(search.Date1, searchDateNow)
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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 1",
						"GET":     get,
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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 1",
						"GET":     get,
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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 1",
						"GET":     get,
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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 1",
						"GET":     get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
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
					roleSuperIngenerQuality: SuperIngenerQuality,
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

//historyInspectionP5
func (s *Server) PagehistoryInspectionP5() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
		StockkeeperWH := false
		WarehouseManager := false
		SuperIngenerQuality := false
		IngenerQuality := false
		Quality := false
		Inspector := false
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

		if user.Groups == groupQualityP5 || user.Groups == groupWarehouseP5 || user.Groups == groupQuality || user.Groups == groupWarehouse {
			GroupP5 = true
			GroupP1 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				LoggedIn = true
				fmt.Println("кладовщик склада - ", StockkeeperWH)
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
				LoggedIn = true
			} else if user.Groups == groupQuality {
				Quality = true
				Inspector = true
				LoggedIn = true
				fmt.Println("pageInspection quality - ", Quality)
			} else if user.Role == roleWarehouseManager {
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
				"GroupP5":             GroupP5,
				"GroupP1":             GroupP1,
				//	"GET":           get,
				"LoggedIn": LoggedIn,
				"User":     user.LastName,
				"Username": user.FirstName,
			}

			tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
		}
	}
}

func (s *Server) HistoryInspectionP5() http.HandlerFunc {
	type req struct {
		Date1    string `json:"date1"`
		Date2    string `json:"date2"`
		Material int    `json:"material"`
		EO       string `json:"eo"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		Admin := false
		StockkeeperWH := false
		WarehouseManager := false
		SuperIngenerQuality := false
		IngenerQuality := false
		Quality := false
		Inspector := false
		GroupP5 := false
		GroupP1 := false
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

		if user.Groups == groupQualityP5 || user.Groups == groupWarehouseP5 {
			GroupP5 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				LoggedIn = true
				fmt.Println("кладовщик склада - ", StockkeeperWH)
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
				LoggedIn = true
			} else if user.Role == roleInspector {
				Quality = true
				Inspector = true
				LoggedIn = true
				fmt.Println("pageInspection quality - ", Quality)
			} else if user.Role == roleWarehouseManager {
				WarehouseManager = true
				LoggedIn = true
			}

			search := &req{}
			materialInt, err := strconv.Atoi(r.FormValue("material"))
			if err != nil {
				fmt.Println(err)
				s.logger.Errorf(err.Error())
			}
			search.Date1 = r.FormValue("date1")
			fmt.Println("date1 - ", search.Date1)
			//s.infoLog.Printf("date1 - %v", search.Date1)
			search.Date2 = r.FormValue("date2")
			fmt.Println("date2 - ", search.Date2)
			//s.infoLog.Printf("date2 - %v", search.Date2)
			search.Material = materialInt
			fmt.Println("material - ", search.Material)
			//s.infoLog.Printf("material - %v", search.Material)
			search.EO = r.FormValue("eo")
			//s.infoLog.Printf("eo - %v", search.EO)

			currentData := time.Now()
			searchDateNow := currentData.Format("2006-01-02")

			if search.Date1 == "" && search.Date2 == "" && search.Material == 0 && search.EO == "" {
				fmt.Println("Не заполнены поля ввода")
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
					"GroupP5":             GroupP5,
					"Pobedit":             "Победит 5",
					//"GroupP1":             GroupP1,
					"LoggedIn": LoggedIn,
					//	"GET":                 get,
				}

				err = tpl.ExecuteTemplate(w, "errorSearchHistoryInspection.html", data)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}
			} else if search.Date1 == "" && search.Date2 == "" {
				if search.EO != "" {

					get, err := s.store.Inspection().ListShowDataByEOP5(search.EO)
					if err != nil {
						s.error(w, r, http.StatusUnprocessableEntity, err)
						return
					}
					//	count, _ := s.store.Inspection().CountInspection()
					//	fmt.Println(count)
					//	limit := 5
					//	page, begin := s.Pagination(r, limit)
					//	fmt.Printf("Current Page: %d, Begin: %d\n", page, begin)

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
						"Pobedit":             "Победит 5",
						//"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"GET":     get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				} else if search.Material != 0 {

					//	count, _ := s.store.Inspection().CountInspection()
					//	fmt.Println(count)
					//	limit := 2
					//	page, begin := s.Pagination(r, limit)
					//	fmt.Printf("Current Page: %d, Begin: %d\n", page, begin)
					//	get, err := s.store.Inspection().ListShowDataBySap(search.Material, begin, limit)
					get, err := s.store.Inspection().ListShowDataBySapP5(search.Material)

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
						"Pobedit":             "Победит 5",
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"GET":     get,
					}
					//	RenderJSON(w, get, http.StatusOK)
					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				}
			} else if search.Date1 != "" && search.Date2 == "" {

				if search.Material != 0 {
					fmt.Println("OK Material")

					get, err := s.store.Inspection().ListShowDataByDateAndSAPP5(search.Date1, searchDateNow, search.Material)
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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 5",
						"GET":     get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				} else if search.EO != "" {
					fmt.Println("OK EO")

					get, err := s.store.Inspection().ListShowDataByDateAndEOP5(search.Date1, searchDateNow, search.EO)
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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 5",
						"GET":     get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}

				} else {

					get, err := s.store.Inspection().ListShowDataByDateP5(search.Date1, searchDateNow)
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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 5",
						"GET":     get,
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

					get, err := s.store.Inspection().ListShowDataByDateAndSAPP5(search.Date1, search.Date2, search.Material)
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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 5",
						"GET":     get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				} else if search.EO != "" {
					fmt.Println("OK EO")

					get, err := s.store.Inspection().ListShowDataByDateAndEOP5(search.Date1, search.Date2, search.EO)
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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 5",
						"GET":     get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}

				} else {

					get, err := s.store.Inspection().ListShowDataByDateP5(search.Date1, search.Date2)
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
						//	"GroupP1":             GroupP1,
						"GroupP5": GroupP5,
						"Pobedit": "Победит 5",
						"GET":     get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				}
			}
		}

		if user.Groups == groupQuality || user.Groups == groupWarehouse {
			GroupP1 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				LoggedIn = true
				fmt.Println("кладовщик склада - ", StockkeeperWH)
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
				LoggedIn = true
			} else if user.Role == roleInspector {
				Quality = true
				Inspector = true
				LoggedIn = true
				fmt.Println("pageInspection quality - ", Quality)
			} else if user.Role == roleWarehouseManager {
				WarehouseManager = true
				LoggedIn = true
			}

			search := &req{}
			materialInt, err := strconv.Atoi(r.FormValue("material"))
			if err != nil {
				fmt.Println(err)
				s.logger.Errorf(err.Error())
			}
			search.Date1 = r.FormValue("date1")
			fmt.Println("date1 - ", search.Date1)
			//s.infoLog.Printf("date1 - %v", search.Date1)
			s.logger.Infof("date1 - %v", search.Date1)
			search.Date2 = r.FormValue("date2")
			fmt.Println("date2 - ", search.Date2)
			//s.infoLog.Printf("date2 - %v", search.Date2)
			s.logger.Infof("date2 - %v", search.Date2)
			search.Material = materialInt
			fmt.Println("material - ", search.Material)
			//s.infoLog.Printf("material - %v", search.Material)
			s.logger.Infof("material - %v", search.Material)
			search.EO = r.FormValue("eo")
			//s.infoLog.Printf("eo - %v", search.EO)
			s.logger.Infof("eo - %v", search.EO)

			currentData := time.Now()
			searchDateNow := currentData.Format("2006-01-02")

			if search.Date1 == "" && search.Date2 == "" && search.Material == 0 && search.EO == "" {
				fmt.Println("Не заполнены поля ввода")
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
					//	"GroupP5":             GroupP5,
					"GroupP1":  GroupP1,
					"Pobedit":  "Победит 5",
					"LoggedIn": LoggedIn,
					//	"GET":                 get,
				}

				err = tpl.ExecuteTemplate(w, "errorSearchHistoryInspection.html", data)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}
			} else if search.Date1 == "" && search.Date2 == "" {
				if search.EO != "" {

					get, err := s.store.Inspection().ListShowDataByEOP5(search.EO)
					if err != nil {
						s.error(w, r, http.StatusUnprocessableEntity, err)
						return
					}
					//	count, _ := s.store.Inspection().CountInspection()
					//	fmt.Println(count)
					//	limit := 5
					//	page, begin := s.Pagination(r, limit)
					//	fmt.Printf("Current Page: %d, Begin: %d\n", page, begin)

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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 5",
						//	"GroupP5":             GroupP5,
						"GET": get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				} else if search.Material != 0 {

					//	count, _ := s.store.Inspection().CountInspection()
					//	fmt.Println(count)
					//	limit := 2
					//	page, begin := s.Pagination(r, limit)
					//	fmt.Printf("Current Page: %d, Begin: %d\n", page, begin)
					//	get, err := s.store.Inspection().ListShowDataBySap(search.Material, begin, limit)
					get, err := s.store.Inspection().ListShowDataBySapP5(search.Material)

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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 5",
						//	"GroupP5":             GroupP5,
						"GET": get,
					}
					//	RenderJSON(w, get, http.StatusOK)
					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				}
			} else if search.Date1 != "" && search.Date2 == "" {

				if search.Material != 0 {
					fmt.Println("OK Material")

					get, err := s.store.Inspection().ListShowDataByDateAndSAPP5(search.Date1, searchDateNow, search.Material)
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 5",
						//	"GroupP5":             GroupP5,
						"GET": get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				} else if search.EO != "" {
					fmt.Println("OK EO")

					get, err := s.store.Inspection().ListShowDataByDateAndEOP5(search.Date1, searchDateNow, search.EO)
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 5",
						//	"GroupP5":             GroupP5,
						"GET": get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}

				} else {

					get, err := s.store.Inspection().ListShowDataByDateP5(search.Date1, searchDateNow)
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 5",
						//	"GroupP5":             GroupP5,
						"GET": get,
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

					get, err := s.store.Inspection().ListShowDataByDateAndSAPP5(search.Date1, search.Date2, search.Material)
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 5",
						//	"GroupP5":             GroupP5,
						"GET": get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				} else if search.EO != "" {
					fmt.Println("OK EO")

					get, err := s.store.Inspection().ListShowDataByDateAndEOP5(search.Date1, search.Date2, search.EO)
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 5",
						//	"GroupP5":             GroupP5,
						"GET": get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}

				} else {

					get, err := s.store.Inspection().ListShowDataByDateP5(search.Date1, search.Date2)
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
						"GroupP1":             GroupP1,
						"Pobedit":             "Победит 5",
						//	"GroupP5":             GroupP5,
						"GET": get,
					}

					err = tpl.ExecuteTemplate(w, "showhistoryinspection.html", data)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}
				}
			}
		}

	}
}

func (s *Server) PageInspection() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		//w.Header().Add("Content-Type", "application/json")
		//	Admin := true
		//	Warehouse := true
		StockkeeperWH := false
		WarehouseManager := false
		//	Quality := false
		Inspector := false
		SuperIngenerQuality := false
		IngenerQuality := false
		SuperIngenerQuality2 := false
		GroupP1 := false
		GroupP5 := false
		MixP1P5 := false
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

		if user.Groups == groupQuality || user.Groups == groupWarehouse {
			GroupP1 = true
			if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				SuperIngenerQuality2 = true
				MixP1P5 = true
				LoggedIn = true
				fmt.Println("pageInspection SuperIngenerQuality - ", SuperIngenerQuality)
			} else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				//	Warehouse = false
				//	WarehouseManager = true
				LoggedIn = true
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
			} else if user.Role == roleInspector {
				Inspector = true
			} else if user.Role == roleWarehouseManager {
				WarehouseManager = true
				LoggedIn = true
			} /* else if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			}*/ /**else if user.Groups == groupQuality {
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

			countVerifyComponents, err := s.store.Inspection().CountVerifyComponents()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			listVendor, err := s.store.Vendor().ListVendor()
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
				"StockkeeperWH":         StockkeeperWH,
				"SuperIngenerQuality":   SuperIngenerQuality,
				"SuperIngenerQuality2":  SuperIngenerQuality2,
				"GroupP1":               GroupP1,
				"MixP1P5":               MixP1P5,
				"GET":                   get,
				"ListVendor":            listVendor,
				"CountTotal":            countTotal,
				"HoldInspection":        holdInspection,
				"NotVerifyComponents":   notVerifyComponents,
				"CountVerifyComponents": countVerifyComponents,
				"GetStatic":             getStatic,
				"HoldCountDebitor":      holdCountDebitor,
				"NotVerifyDebitor":      notVerifyDebitor,
				"LoggedIn":              LoggedIn,
				"Pobedit":               "Победит 1",
			}

			tpl.ExecuteTemplate(w, "showinspection.html", groups)
			//	json.NewEncoder(w).Encode(get)
		}

		if user.Groups == groupQualityP5 || user.Groups == groupWarehouseP5 {
			GroupP5 = true
			if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				SuperIngenerQuality2 = true
				LoggedIn = true
				fmt.Println("pageInspection SuperIngenerQuality - ", SuperIngenerQuality)
			} else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				//	Warehouse = false
				//	WarehouseManager = true
				LoggedIn = true
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
			} else if user.Role == roleInspector {
				Inspector = true
			} else if user.Role == roleWarehouseManager {
				WarehouseManager = true
				LoggedIn = true
			} /* else if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			}*/ /**else if user.Groups == groupQuality {
				//	Quality = true
				Inspector = true
				IngenerQuality = true
				LoggedIn = true
				//	fmt.Println("pageInspection quality - ", Quality)
			} */

			get, err := s.store.Inspection().ListInspectionP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			countTotal, err := s.store.Inspection().CountTotalInspectionP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			holdInspection, err := s.store.Inspection().HoldInspectionP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			notVerifyComponents, err := s.store.Inspection().NotVerifyComponentsP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			getStatic, err := s.store.Inspection().CountDebitorP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			holdCountDebitor, err := s.store.Inspection().HoldCountDebitorP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			notVerifyDebitor, err := s.store.Inspection().NotVerifyDebitorP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			countVerifyComponents, err := s.store.Inspection().CountVerifyComponentsP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			listVendor, err := s.store.Vendor().ListVendor()
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
				"StockkeeperWH":         StockkeeperWH,
				"SuperIngenerQuality":   SuperIngenerQuality,
				"SuperIngenerQuality2":  SuperIngenerQuality2,
				"GroupP5":               GroupP5,
				"GET":                   get,
				"ListVendor":            listVendor,
				"CountTotal":            countTotal,
				"HoldInspection":        holdInspection,
				"NotVerifyComponents":   notVerifyComponents,
				"GetStatic":             getStatic,
				"HoldCountDebitor":      holdCountDebitor,
				"NotVerifyDebitor":      notVerifyDebitor,
				"CountVerifyComponents": countVerifyComponents,
				"LoggedIn":              LoggedIn,
				"Pobedit":               "Победит 5",
			}

			tpl.ExecuteTemplate(w, "showinspection.html", groups)
			//	json.NewEncoder(w).Encode(get)
		}
	}
}

func (s *Server) PageInspectionMix() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		//w.Header().Add("Content-Type", "application/json")
		//	Admin := true
		//	Warehouse := true
		StockkeeperWH := false
		WarehouseManager := false
		//	Quality := false
		Inspector := false
		SuperIngenerQuality := false
		IngenerQuality := false
		SuperIngenerQuality2 := false
		GroupP1 := false
		GroupP5 := false
		//MixP1P5 := false
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

		if user.Groups == groupQuality && user.Role == roleSuperIngenerQuality ||
			user.Groups == groupQuality && user.Role == roleIngenerQuality {
			//	MixP1P5 = true
			GroupP1 = true
			if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				SuperIngenerQuality2 = true
				LoggedIn = true
				fmt.Println("pageInspection SuperIngenerQuality - ", SuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
			} /* else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				//	Warehouse = false
				//	WarehouseManager = true
				LoggedIn = true
			} */
			/* else if user.Role == roleInspector {
				Inspector = true
			} else if user.Role == roleWarehouseManager {
				WarehouseManager = true
				LoggedIn = true
			} */
			/* else if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			}*/ /**else if user.Groups == groupQuality {
				//	Quality = true
				Inspector = true
				IngenerQuality = true
				LoggedIn = true
				//	fmt.Println("pageInspection quality - ", Quality)
			} */

			get, err := s.store.Inspection().ListInspectionP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			countTotal, err := s.store.Inspection().CountTotalInspectionP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			holdInspection, err := s.store.Inspection().HoldInspectionP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			notVerifyComponents, err := s.store.Inspection().NotVerifyComponentsP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			getStatic, err := s.store.Inspection().CountDebitorP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			holdCountDebitor, err := s.store.Inspection().HoldCountDebitorP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			notVerifyDebitor, err := s.store.Inspection().NotVerifyDebitorP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			countVerifyComponents, err := s.store.Inspection().CountVerifyComponentsP5()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			listVendor, err := s.store.Vendor().ListVendor()
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
				//"MixP1P5":              MixP1P5,
				"GroupP1":               GroupP1,
				"GET":                   get,
				"ListVendor":            listVendor,
				"CountTotal":            countTotal,
				"HoldInspection":        holdInspection,
				"NotVerifyComponents":   notVerifyComponents,
				"GetStatic":             getStatic,
				"HoldCountDebitor":      holdCountDebitor,
				"NotVerifyDebitor":      notVerifyDebitor,
				"CountVerifyComponents": countVerifyComponents,
				"LoggedIn":              LoggedIn,
				"Pobedit":               "Победит 5",
			}

			tpl.ExecuteTemplate(w, "showinspectionmix.html", groups)
			//	json.NewEncoder(w).Encode(get)
		}

		if user.Groups == groupQualityP5 && user.Role == roleSuperIngenerQuality ||
			user.Groups == groupQualityP5 && user.Role == roleIngenerQuality {
			//	MixP1P5 = true
			GroupP5 = true
			if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				SuperIngenerQuality2 = true
				LoggedIn = true
				fmt.Println("pageInspection SuperIngenerQuality - ", SuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
			} /* else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				//	Warehouse = false
				//	WarehouseManager = true
				LoggedIn = true
			} */
			/* else if user.Role == roleInspector {
				Inspector = true
			} else if user.Role == roleWarehouseManager {
				WarehouseManager = true
				LoggedIn = true
			} */
			/* else if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			}*/ /**else if user.Groups == groupQuality {
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

			countVerifyComponents, err := s.store.Inspection().CountVerifyComponents()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			listVendor, err := s.store.Vendor().ListVendor()
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
				//"MixP1P5":              MixP1P5,
				"GroupP5":               GroupP5,
				"GET":                   get,
				"ListVendor":            listVendor,
				"CountTotal":            countTotal,
				"HoldInspection":        holdInspection,
				"NotVerifyComponents":   notVerifyComponents,
				"GetStatic":             getStatic,
				"HoldCountDebitor":      holdCountDebitor,
				"NotVerifyDebitor":      notVerifyDebitor,
				"CountVerifyComponents": countVerifyComponents,
				"LoggedIn":              LoggedIn,
				"Pobedit":               "Победит 1",
			}

			tpl.ExecuteTemplate(w, "showinspectionmix.html", groups)
			//	json.NewEncoder(w).Encode(get)
		}

	}
}

func (s *Server) PageupdateInspection() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		//	w.Header().Set("Access-Control-Allow-Origin", "")
		//	if r.Method == http.MethodOptions {
		//		return
		//	}
		Admin := false
		SuperIngenerQuality := false
		IngenerQuality := false
		Inspector := false
		GroupP1 := false
		GroupP5 := false
		LoggedIn := false

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
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
		//s.infoLog.Printf("test user.Groups - %s", user.Groups)
		s.logger.Infof("test user.Groups - %s", user.Groups)

		if user.Groups == groupQuality || user.Groups == groupAdministrator {
			GroupP1 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
				fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
				LoggedIn = true
				fmt.Println("IngenerQuality - ", IngenerQuality)
			} else if user.Role == roleInspector {
				Inspector = true
				LoggedIn = true

			}
			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditInspection(id)
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
				"GroupP1":             GroupP1,
				"User":                user.LastName,
				"Username":            user.FirstName,
			}
			err = tpl.ExecuteTemplate(w, "updateinspection.html", data)
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
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
				LoggedIn = true
				fmt.Println("IngenerQuality - ", IngenerQuality)
			} else if user.Role == roleInspector {
				Inspector = true
				LoggedIn = true

			}
			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditInspectionP5(id)
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
				"GroupP5":             GroupP5,
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
}

func (s *Server) PageupdateInspectionJSON() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		//	w.Header().Set("Access-Control-Allow-Origin", "")
		//	if r.Method == http.MethodOptions {
		//		return
		//	}
		Admin := false
		SuperIngenerQuality := false
		IngenerQuality := false
		Inspector := false
		GroupP1 := false
		GroupP5 := false
		LoggedIn := false

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
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
		fmt.Println("Test json page update")
		//s.infoLog.Printf("user.Groups - %v", user.Groups)
		s.logger.Infof("user.Groups - %v", user.Groups)
		/*
			if (user.Groups == groupQuality && user.Role == roleSuperIngenerQuality) ||
				(user.Groups == groupQuality && user.Role == roleIngenerQuality) {
				fmt.Println("Test mix page update")
				GroupP1 = true
				if user.Role == roleAdministrator {
					Admin = true
					LoggedIn = true
				} else if user.Role == roleSuperIngenerQuality {
					SuperIngenerQuality = true
					LoggedIn = true
					fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
				} else if user.Role == roleIngenerQuality {
					IngenerQuality = true
					LoggedIn = true
					fmt.Println("IngenerQuality - ", IngenerQuality)
				} /*else if user.Role == roleInspector {
					Inspector = true
					LoggedIn = true

				}*/
		//fmt.Println("ID - ?", id)
		/*
				get, err := s.store.Inspection().EditInspectionP5(id)
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				data := map[string]interface{}{
					"Admin":               Admin,
					"SuperIngenerQuality": SuperIngenerQuality,
					"IngenerQuality":      IngenerQuality,
					"Inspector":           Inspector,
					"GroupP1":             GroupP1,
					"GET":                 get,
					"LoggedIn":            LoggedIn,
					"User":                user.LastName,
					"Username":            user.FirstName,
				}
				err = tpl.ExecuteTemplate(w, "updateinspectionjson.html", data)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}
			}
		*/
		if user.Groups == groupQuality {
			GroupP1 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
				fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
				LoggedIn = true
				fmt.Println("IngenerQuality - ", IngenerQuality)
			} else if user.Role == roleInspector {
				Inspector = true
				LoggedIn = true

			}
			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditInspection(id)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			data := map[string]interface{}{
				"Admin":               Admin,
				"SuperIngenerQuality": SuperIngenerQuality,
				"IngenerQuality":      IngenerQuality,
				"Inspector":           Inspector,
				"GroupP1":             GroupP1,
				"GET":                 get,
				"LoggedIn":            LoggedIn,
				"User":                user.LastName,
				"Username":            user.FirstName,
			}
			err = tpl.ExecuteTemplate(w, "updateinspectionjson.html", data)
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
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
				LoggedIn = true
				fmt.Println("IngenerQuality - ", IngenerQuality)
			} else if user.Role == roleInspector {
				Inspector = true
				LoggedIn = true

			}
			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditInspectionP5(id)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			data := map[string]interface{}{
				"Admin":               Admin,
				"SuperIngenerQuality": SuperIngenerQuality,
				"IngenerQuality":      IngenerQuality,
				"Inspector":           Inspector,
				"GroupP5":             GroupP5,
				"GET":                 get,
				"LoggedIn":            LoggedIn,
				"User":                user.LastName,
				"Username":            user.FirstName,
			}
			err = tpl.ExecuteTemplate(w, "updateinspectionjson.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}

	}
}

func (s *Server) PageupdateInspectionJSONmix() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		//	w.Header().Set("Access-Control-Allow-Origin", "")
		//	if r.Method == http.MethodOptions {
		//		return
		//	}
		Admin := false
		SuperIngenerQuality := false
		IngenerQuality := false
		Inspector := false
		GroupP1 := false
		GroupP5 := false
		LoggedIn := false

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
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
		fmt.Println("Test json page update")
		//s.infoLog.Printf("user.Groups - %v", user.Groups)
		s.logger.Infof("user.Groups - %v", user.Groups)

		if (user.Groups == groupQuality && user.Role == roleSuperIngenerQuality) ||
			(user.Groups == groupQuality && user.Role == roleIngenerQuality) {
			fmt.Println("Test mix page update")
			GroupP1 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
				fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
				LoggedIn = true
				fmt.Println("IngenerQuality - ", IngenerQuality)
			} /*else if user.Role == roleInspector {
				Inspector = true
				LoggedIn = true

			}*/
			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditInspectionP5(id)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			data := map[string]interface{}{
				"Admin":               Admin,
				"SuperIngenerQuality": SuperIngenerQuality,
				"IngenerQuality":      IngenerQuality,
				"Inspector":           Inspector,
				"GroupP1":             GroupP1,
				"GET":                 get,
				"LoggedIn":            LoggedIn,
				"User":                user.LastName,
				"Username":            user.FirstName,
			}
			err = tpl.ExecuteTemplate(w, "updateinspectionjsonmix.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}

		if (user.Groups == groupQualityP5 && user.Role == roleSuperIngenerQuality) ||
			(user.Groups == groupQualityP5 && user.Role == roleIngenerQuality) {
			fmt.Println("Test mix page update")
			GroupP5 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				SuperIngenerQuality = true
				LoggedIn = true
				fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				IngenerQuality = true
				LoggedIn = true
				fmt.Println("IngenerQuality - ", IngenerQuality)
			} /*else if user.Role == roleInspector {
				Inspector = true
				LoggedIn = true

			}*/
			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditInspection(id)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			data := map[string]interface{}{
				"Admin":               Admin,
				"SuperIngenerQuality": SuperIngenerQuality,
				"IngenerQuality":      IngenerQuality,
				"Inspector":           Inspector,
				"GroupP5":             GroupP5,
				"GET":                 get,
				"LoggedIn":            LoggedIn,
				"User":                user.LastName,
				"Username":            user.FirstName,
			}
			err = tpl.ExecuteTemplate(w, "updateinspectionjsonmix.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
	}
}

func (s *Server) UpdateInspection() http.HandlerFunc {

	// response format
	type response struct {
		ID      int64  `json:"id,omitempty"`
		Status  string `json:"status,omitempty"`
		Note    string `json:"note,omitempty"`
		Message string `json:"message,omitempty"`
	}

	type request struct {
		ID     int    `json:"ID"`
		Status string `json:"status"`
		Note   string `json:"note"`
	}

	type requestJSON struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	///	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		//	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		//GroupP1 = groupQuality
		//GroupP5 = groupQuality

		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
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

		/*	body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}
			var hdata []request
			json.Unmarshal(body, &hdata)
			fmt.Printf("body requestJSON: %s", body)
			fmt.Println("\njson  struct hdata requestJSON", hdata)*/
		u := &model.Inspection{
			ID:         req.ID,
			Status:     req.Status,
			Note:       req.Note,
			Update:     user.LastName, //
			Dateupdate: currentTime,   // Dateaccept
			Timeupdate: currentTime,   // Timeaccept
			Groups:     user.Groups,
		}

		//	for _, v := range hdata {
		/*	fmt.Println("проверка в цкле - requestJSON", v.ID, v.Status, v.Note)
			idRoll, err := strconv.Atoi(v.ID)
			if err != nil {
				log.Fatal(err)
			}*/
		/*	u := &model.Inspection{
			ID:         idRoll,
			Status:     v.Status,
			Note:       v.Note,
			Update:     user.LastName, //
			Dateupdate: currentTime,   // Dateaccept
			Timeupdate: currentTime,   // Timeaccept
			Groups:     user.Groups,
		}*/
		if user.Groups == groupQuality {

			if err := s.store.Inspection().UpdateInspection(u); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			// format the message string
			/*		msg := fmt.Sprintf("Inspection updated successfully. Total rows/record affected")
					//msg := "Inspection updated successfully. Total rows/record affected"

					// format the response message
					res := response{
						ID:      int64(idRoll),
						Status:  v.Status,
						Note:    v.Note,
						Message: msg,
					}
					// send the response
					json.NewEncoder(w).Encode(res)*/
			//	}

			/*	err = tpl.ExecuteTemplate(w, "layout", nil)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}*/
			http.Redirect(w, r, "/operation/statusinspection", 303)
		}
		if user.Groups == groupQualityP5 {

			if err := s.store.Inspection().UpdateInspectionP5(u); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			// format the message string
			/*		msg := fmt.Sprintf("Inspection updated successfully. Total rows/record affected")
					//msg := "Inspection updated successfully. Total rows/record affected"

					// format the response message
					res := response{
						ID:      int64(idRoll),
						Status:  v.Status,
						Note:    v.Note,
						Message: msg,
					}
					// send the response
					json.NewEncoder(w).Encode(res)*/
			//	}

			/*	err = tpl.ExecuteTemplate(w, "layout", nil)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}*/
			http.Redirect(w, r, "/operation/statusinspection", 303)
		}
	}
}

func (s *Server) UpdateInspectionJSON() http.HandlerFunc {
	// response format
	type response struct {
		ID       int64  `json:"id,omitempty"`
		Status   string `json:"status,omitempty"`
		Note     string `json:"note,omitempty"`
		Message  string `json:"message,omitempty"`
		Lastname string `json:"lastname,omitempty"`
	}

	type requestJSON struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		Note     string `json:"note"`
		Lastname string `json:"lastname"`
	}
	///	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		//	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}
		var hdata []requestJSON
		//var hdata map[string]interface{}
		json.Unmarshal(body, &hdata)
		//json.Unmarshal([]byte(body), &hdata)
		fmt.Printf("body requestJSON: %s", body)
		//s.infoLog.Printf("Loading body requestJSON: %v\n", body)
		s.logger.Infof("Loading body requestJSON: %v\n", body)
		fmt.Println("\njson  struct hdata requestJSON", hdata)
		//s.infoLog.Printf("Loading hdata requestJSON: %v\n", hdata)
		s.logger.Infof("Loading hdata requestJSON: %v\n", hdata)
		//hdata2 := hdata["hdata2"].(map[string]interface{})

		for _, v := range hdata {
			fmt.Println("проверка в цкле - requestJSON", v.ID, v.Status, v.Note)
			//	fmt.Println("проверка в цкле - requestJSON", v.(string), v.(string), v.(string))
			idRoll, err := strconv.Atoi(v.ID)
			lastname := v.Lastname
			//idRoll, err := strconv.Atoi(v.(string))
			if err != nil {
				log.Println(err)
				s.logger.Errorf(err.Error())
			}

			u := &model.Inspection{
				ID:         idRoll,
				Status:     v.Status,      // v.(string),
				Note:       v.Note,        // v.(string),
				Update:     user.LastName, //
				Dateupdate: currentTime,   // Dateaccept
				Timeupdate: currentTime,   // Timeaccept
				Groups:     user.Groups,
			}
			/*	if user.Groups == groupQuality && user.Role == roleSuperIngenerQuality ||
				user.Groups == groupQuality && user.Role == roleIngenerQuality {
				if err := s.store.Inspection().UpdateInspectionP5(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				// format the message string
				//msg := fmt.Sprintf("Inspection updated successfully. Total rows/record affected")
				msg := "Данные успешно отправлены."

				// format the response message
				res := response{
					ID:      int64(idRoll),
					Status:  v.Status, // v.(string),
					Note:    v.Note,   // v.(string),
					Message: msg,
				}
				// send the response
				json.NewEncoder(w).Encode(res)
			}*/
			if user.Groups == groupQuality {
				if err := s.store.Inspection().UpdateInspection(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				// format the message string
				//msg := fmt.Sprintf("Inspection updated successfully. Total rows/record affected")
				msg := "Данные успешно отправлены."

				// format the response message
				res := response{
					ID:       int64(idRoll),
					Status:   v.Status, // v.(string),
					Note:     v.Note,   // v.(string),
					Lastname: lastname,
					Message:  msg,
				}
				// send the response
				json.NewEncoder(w).Encode(res)
			}
			if user.Groups == groupQualityP5 {
				if err := s.store.Inspection().UpdateInspectionP5(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				// format the message string
				//msg := fmt.Sprintf("Inspection updated successfully. Total rows/record affected")
				msg := "Данные успешно отправлены."

				// format the response message
				res := response{
					ID:      int64(idRoll),
					Status:  v.Status, // v.(string),
					Note:    v.Note,   // v.(string),
					Message: msg,
				}
				// send the response
				json.NewEncoder(w).Encode(res)
			}
		}
		/*	err = tpl.ExecuteTemplate(w, "layout", nil)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}*/
		//	http.Redirect(w, r, "/operation/statusinspection", 303)
	}
}

func (s *Server) UpdateInspectionJSONmix() http.HandlerFunc {
	// response format
	type response struct {
		ID      int64  `json:"id,omitempty"`
		Status  string `json:"status,omitempty"`
		Note    string `json:"note,omitempty"`
		Message string `json:"message,omitempty"`
	}

	type requestJSON struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	///	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		//	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}
		var hdata []requestJSON
		//var hdata map[string]interface{}
		json.Unmarshal(body, &hdata)
		//json.Unmarshal([]byte(body), &hdata)
		fmt.Printf("body requestJSON: %s", body)
		//s.infoLog.Printf("Loading body requestJSON: %v\n", body)
		s.logger.Infof("Loading body requestJSON: %v\n", body)
		fmt.Println("\njson  struct hdata requestJSON", hdata)
		//s.infoLog.Printf("Loading hdata requestJSON: %v\n", hdata)
		s.logger.Infof("Loading hdata requestJSON: %v\n", hdata)
		//hdata2 := hdata["hdata2"].(map[string]interface{})

		for _, v := range hdata {
			fmt.Println("проверка в цкле - requestJSON", v.ID, v.Status, v.Note)
			//	fmt.Println("проверка в цкле - requestJSON", v.(string), v.(string), v.(string))
			idRoll, err := strconv.Atoi(v.ID)
			//idRoll, err := strconv.Atoi(v.(string))
			if err != nil {
				log.Println(err)
				s.logger.Errorf(err.Error())
			}
			u := &model.Inspection{
				ID:         idRoll,
				Status:     v.Status,      // v.(string),
				Note:       v.Note,        // v.(string),
				Update:     user.LastName, //
				Dateupdate: currentTime,   // Dateaccept
				Timeupdate: currentTime,   // Timeaccept
				Groups:     user.Groups,
			}
			if user.Groups == groupQuality && user.Role == roleSuperIngenerQuality ||
				user.Groups == groupQuality && user.Role == roleIngenerQuality {
				if err := s.store.Inspection().UpdateInspectionP5(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				// format the message string
				//msg := fmt.Sprintf("Inspection updated successfully. Total rows/record affected")
				msg := "Данные успешно отправлены."

				// format the response message
				res := response{
					ID:      int64(idRoll),
					Status:  v.Status, // v.(string),
					Note:    v.Note,   // v.(string),
					Message: msg,
				}
				// send the response
				json.NewEncoder(w).Encode(res)
			}

			if user.Groups == groupQualityP5 && user.Role == roleSuperIngenerQuality ||
				user.Groups == groupQualityP5 && user.Role == roleIngenerQuality {
				if err := s.store.Inspection().UpdateInspection(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				// format the message string
				//msg := fmt.Sprintf("Inspection updated successfully. Total rows/record affected")
				msg := "Данные успешно отправлены."

				// format the response message
				res := response{
					ID:      int64(idRoll),
					Status:  v.Status, // v.(string),
					Note:    v.Note,   // v.(string),
					Message: msg,
				}
				// send the response
				json.NewEncoder(w).Encode(res)
			}
		}
		/*	err = tpl.ExecuteTemplate(w, "layout", nil)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}*/
		//	http.Redirect(w, r, "/operation/statusinspection", 303)
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

		u := &model.Inspection{
			ID: req.ID,
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

		if user.Groups == groupWarehouse || user.Groups == groupQuality {
			//	s.Lock()
			fmt.Println("call")
			if err := s.store.Inspection().DeleteItemInspection(u); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			//	s.Unlock()
			http.Redirect(w, r, "/operation/statusinspection", 303)
		}
		if user.Groups == groupWarehouseP5 || user.Groups == groupQualityP5 {
			//	s.Lock()
			fmt.Println("call P5")
			if err := s.store.Inspection().DeleteItemInspectionP5(u); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			//	s.Unlock()
			http.Redirect(w, r, "/operation/statusinspection", 303)
		}
	}
}

//ListAcceptWHInspection
func (s *Server) PageListAcceptWHInspection() http.HandlerFunc { // acceptinspection.html showinspection.html
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "acceptinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

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

		if user.Groups == groupWarehouse {
			stockkeeperWH = true
			WarehouseManager = true
		}
		/*
			if user.Groups == groupQuality {
				quality = true
			} else if user.Groups == groupWarehouse {
				stockkeeperWH = true
			}
		*/
		if user.Role == roleAdministrator {
			Admin = true
			LoggedIn = true
		} else if user.Role == roleStockkeeperWH {
			StockkeeperWH = true
			LoggedIn = true
			fmt.Println("кладовщик склада - ", StockkeeperWH)
		} else if user.Role == roleWarehouseManager {
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
			//	roleSuperIngenerQuality: superIngenerQuality,
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
		s.mu.Lock()
		defer s.mu.Unlock()

		Admin := false
		StockkeeperWH := false
		WarehouseManager := false
		GroupP1 := false
		GroupP5 := false
		LoggedIn := false

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
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
		//s.infoLog.Printf("user.Groups - %v\n", user.Groups)
		s.logger.Infof("user.Groups - %v\n", user.Groups)

		if user.Groups == groupWarehouse {
			GroupP1 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true

			} else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				LoggedIn = true

			} else if user.Role == roleWarehouseManager {
				WarehouseManager = true
				LoggedIn = true
			}

			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditAcceptWarehouseInspection(id)
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
				"GroupP1":          GroupP1,
				"LoggedIn":         LoggedIn,
				"GET":              get,
			}
			err = tpl.ExecuteTemplate(w, "acceptWarehouseInspection.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
		if user.Groups == groupWarehouseP5 {
			GroupP5 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true

			} else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				LoggedIn = true

			} else if user.Role == roleWarehouseManager {
				WarehouseManager = true
				LoggedIn = true
			}

			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditAcceptWarehouseInspectionP5(id)
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
				"GroupP5":          GroupP5,
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
		s.mu.Lock()
		defer s.mu.Unlock()

		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
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
		if user.Groups == groupWarehouse {

			if err := s.store.Inspection().AcceptWarehouseInspection(u); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			/*	err = tpl.ExecuteTemplate(w, "layout", nil)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}*/
			//http.Redirect(w, r, "/statusinspectionforwh", 303)
			http.Redirect(w, r, "/operation/statusinspection", 303)
		}
		if user.Groups == groupWarehouseP5 {
			fmt.Println("Test accept")

			if err := s.store.Inspection().AcceptWarehouseInspectionP5(u); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			/*	err = tpl.ExecuteTemplate(w, "layout", nil)
				if err != nil {
					http.Error(w, err.Error(), 400)
					return
				}*/
			//http.Redirect(w, r, "/statusinspectionforwh", 303)
			http.Redirect(w, r, "/operation/statusinspection", 303)
		}
	}
}

func (s *Server) PageacceptWarehouseInspectionJSON() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "acceptWarehouseInspection.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		Admin := false
		StockkeeperWH := false
		WarehouseManager := false
		GroupP1 := false
		GroupP5 := false
		LoggedIn := false

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
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
		//s.infoLog.Printf("user.Groups - %v\n", user.Groups)
		s.logger.Info("user.Groups - %v\n", user.Groups)

		if user.Groups == groupWarehouse {
			GroupP1 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true

			} else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				LoggedIn = true

			} else if user.Role == roleWarehouseManager {
				WarehouseManager = true
				LoggedIn = true
			}

			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditAcceptWarehouseInspection(id)
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
				"GroupP1":          GroupP1,
				"LoggedIn":         LoggedIn,
				"GET":              get,
			}
			err = tpl.ExecuteTemplate(w, "acceptWarehouseInspectionjson.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
		if user.Groups == groupWarehouseP5 {
			GroupP5 = true
			if user.Role == roleAdministrator {
				Admin = true
				LoggedIn = true

			} else if user.Role == roleStockkeeperWH {
				StockkeeperWH = true
				LoggedIn = true

			} else if user.Role == roleWarehouseManager {
				WarehouseManager = true
				LoggedIn = true
			}

			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditAcceptWarehouseInspectionP5(id)
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
				"GroupP5":          GroupP5,
				"LoggedIn":         LoggedIn,
				"GET":              get,
			}
			err = tpl.ExecuteTemplate(w, "acceptWarehouseInspectionjson.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
	}
}

func (s *Server) AcceptWarehouseInspectionJSON() http.HandlerFunc {
	// response format
	type response struct {
		ID       int64  `json:"id,omitempty"`
		Location string `json:"location,omitempty"`
		Message  string `json:"message,omitempty"`
	}

	type requestJSON struct {
		ID       string `json:"id"`
		Location string `json:"location"`
	}
	///	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "acceptWarehouseInspection.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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
		/*
			req.ID = id
			req.Location = r.FormValue("location")
		*/
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}
		var hdata []requestJSON
		//var hdata map[string]interface{}
		json.Unmarshal(body, &hdata)
		//json.Unmarshal([]byte(body), &hdata)
		fmt.Printf("body requestJSON: %s", body)
		//s.infoLog.Printf("Accept body requestJSON: %v\n", body)
		s.logger.Infof("Accept body requestJSON: %v\n", body)
		fmt.Println("\njson  struct hdata requestJSON", hdata)
		//s.infoLog.Printf("Accept hdata requestJSON: %v\n", hdata)
		s.logger.Infof("Accept hdata requestJSON: %v\n", hdata)
		//hdata2 := hdata["hdata2"].(map[string]interface{})

		for _, v := range hdata {
			fmt.Println("проверка в цкле - requestJSON", v.ID, v.Location)
			idRoll, err := strconv.Atoi(v.ID)
			//idRoll, err := strconv.Atoi(v.(string))
			if err != nil {
				log.Println(err)
				s.logger.Errorf(err.Error())
			}
			u := &model.Inspection{
				ID:             idRoll,
				Location:       v.Location,
				Lastnameaccept: user.LastName, // Lastnameaccept
				Dateaccept:     currentTime,   // Dateaccept
				Timeaccept:     currentTime,   // Timeaccept
				Groups:         user.Groups,
			}
			if user.Groups == groupWarehouse {

				if err := s.store.Inspection().AcceptWarehouseInspection(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				/*	err = tpl.ExecuteTemplate(w, "layout", nil)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}*/
				//http.Redirect(w, r, "/statusinspectionforwh", 303)
				msg := "Данные успешно отправлены."

				// format the response message
				res := response{
					ID:       int64(idRoll),
					Location: v.Location, // v.(string),
					Message:  msg,
				}
				// send the response
				json.NewEncoder(w).Encode(res)
				//	http.Redirect(w, r, "/operation/statusinspection", 303)
			}
			if user.Groups == groupWarehouseP5 {
				fmt.Println("Test accept")

				if err := s.store.Inspection().AcceptWarehouseInspectionP5(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				/*	err = tpl.ExecuteTemplate(w, "layout", nil)
					if err != nil {
						http.Error(w, err.Error(), 400)
						return
					}*/

				//http.Redirect(w, r, "/statusinspectionforwh", 303)
				msg := "Данные успешно отправлены."

				// format the response message
				res := response{
					ID:       int64(idRoll),
					Location: v.Location, // v.(string),
					Message:  msg,
				}
				// send the response
				json.NewEncoder(w).Encode(res)
				// http.Redirect(w, r, "/operation/statusinspection", 303)
			}
		}
	}
}

func (s *Server) AcceptGroupsWarehouseInspection() http.HandlerFunc {
	type req struct {
		ScanID         string `json:"scanidAccept"`
		SAP            int
		Lot            string
		Roll           int
		Qty            int
		ProductionDate string
		NumberVendor   string
		Location       string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}

		var rdata []req
		json.Unmarshal(body, &rdata)
		fmt.Printf("test AcceptGroupsWarehouseInspection %s", body)
		//s.infoLog.Printf("Loading json accept list EO %s", body)
		s.logger.Infof("Loading json accept list EO %s", body)
		fmt.Println("\nall of the rdata AcceptGroupsWarehouseInspection", rdata)
		//s.infoLog.Printf("Loading json rdata list EO %v", rdata)
		s.logger.Infof("Loading json rdata list EO %v", rdata)

		currentTime := time.Now()

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

		const statusTransfer = "Принято на склад с ВК"

		if user.Groups == groupWarehouse {
			for _, v := range rdata {
				if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 45) {
					idMaterial := v.ScanID[0:45]

					if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 45) {
						u := &model.Inspection{
							Location:       statusTransfer,
							Lastnameaccept: user.LastName, // Lastnameaccept
							Dateaccept:     currentTime,   // Dateaccept
							Timeaccept:     currentTime,   // Timeaccept
							IdMaterial:     idMaterial,
							Groups:         user.Groups,
						}

						if err := s.store.Inspection().AcceptGroupsWarehouseInspection(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)

							return
						}

					} else {
						if (strings.Contains(v.ScanID[0:1], "P") == false) && (len(v.ScanID) != 45) {
							fmt.Println("не верное сканирование :\n" + v.ScanID + "\n")
							s.logger.Errorf("не верное сканирование :\n" + v.ScanID + "\n")
							//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
						}
						//	tpl.Execute(w, data)
						return
					}
				}
				if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 35) {
					idMaterial := v.ScanID[0:35]

					if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 35) {
						u := &model.Inspection{
							Location:       statusTransfer,
							Lastnameaccept: user.LastName, // Lastnameaccept
							Dateaccept:     currentTime,   // Dateaccept
							Timeaccept:     currentTime,   // Timeaccept
							IdMaterial:     idMaterial,
							Groups:         user.Groups,
						}

						if err := s.store.Inspection().AcceptGroupsWarehouseInspection(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)

							return
						}

					} else {
						if (strings.Contains(v.ScanID[0:1], "P") == false) && (len(v.ScanID) != 35) {
							fmt.Println("не верное сканирование :\n" + v.ScanID + "\n")
							s.logger.Errorf("не верное сканирование :\n" + v.ScanID + "\n")
							//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
						}
						//	tpl.Execute(w, data)
						return
					}
				}
			}
		}
		if user.Groups == groupWarehouseP5 {
			for _, v := range rdata {
				if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 45) {
					idMaterial := v.ScanID[0:45]

					if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 45) {
						u := &model.Inspection{
							Location:       statusTransfer,
							Lastnameaccept: user.LastName, // Lastnameaccept
							Dateaccept:     currentTime,   // Dateaccept
							Timeaccept:     currentTime,   // Timeaccept
							IdMaterial:     idMaterial,
							Groups:         user.Groups,
						}

						if err := s.store.Inspection().AcceptGroupsWarehouseInspectionP5(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)

							return
						}

					} else {
						if (strings.Contains(v.ScanID[0:1], "P") == false) && (len(v.ScanID) != 45) {
							fmt.Println("не верное сканирование :\n" + v.ScanID + "\n")
							s.logger.Errorf("не верное сканирование :\n" + v.ScanID + "\n")
							//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
						}
						//	tpl.Execute(w, data)
						return
					}
				}
				if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 35) {
					idMaterial := v.ScanID[0:35]

					if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 35) {
						u := &model.Inspection{
							Location:       statusTransfer,
							Lastnameaccept: user.LastName, // Lastnameaccept
							Dateaccept:     currentTime,   // Dateaccept
							Timeaccept:     currentTime,   // Timeaccept
							IdMaterial:     idMaterial,
							Groups:         user.Groups,
						}

						if err := s.store.Inspection().AcceptGroupsWarehouseInspectionP5(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)

							return
						}

					} else {
						if (strings.Contains(v.ScanID[0:1], "P") == false) && (len(v.ScanID) != 35) {
							fmt.Println("не верное сканирование :\n" + v.ScanID + "\n")
							s.logger.Errorf("не верное сканирование :\n" + v.ScanID + "\n")
							//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
						}
						//	tpl.Execute(w, data)
						return
					}
				}
			}
		}
	}

}
