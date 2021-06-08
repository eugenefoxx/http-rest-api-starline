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

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Role == roleAdministrator {
			statusAdmin = true
			statusLoggedIn = true

		} else if user.Role == roleStockkeeperWH {
			statusStockkeeperWH = true
			statusLoggedIn = true
			fmt.Println("кладовщик склада - ", statusStockkeeperWH)
		} else if user.Role == roleSuperIngenerQuality {
			statusSuperIngenerQuality = true
			statusLoggedIn = true
		} else if user.Role == roleWarehouseManager {
			statusWarehouseManager = true
			statusLoggedIn = true
		}
		data := map[string]interface{}{
			"Admin":               statusAdmin,
			"StockkeeperWH":       statusStockkeeperWH,
			"WarehouseManager":    statusWarehouseManager,
			"SuperIngenerQuality": statusSuperIngenerQuality,
			//	"GET":           get,
			"LoggedIn": statusLoggedIn,
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

		user := r.Context().Value(ctxKeyUser).(*model.User)

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
		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupQuality || user.Groups == groupWarehouse || user.Groups == groupQualityP5 ||
			user.Groups == groupWarehouseP5 {
			statusGroupP1 = true
			statusGroupP5 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true
				fmt.Println("кладовщик склада - ", statusStockkeeperWH)
			} else if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
			} else if user.Groups == groupQuality {
				//statusQuality = true
				statusInspector = true
				statusLoggedIn = true
			//	fmt.Println("pageInspection quality - ", Quality)
			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
			}
			data := map[string]interface{}{
				"Admin":               statusAdmin,
				"StockkeeperWH":       statusStockkeeperWH,
				"SuperIngenerQuality": statusSuperIngenerQuality,
				"WarehouseManager":    statusWarehouseManager,
				"IngenerQuality":      statusIngenerQuality,
			//	"Quality":             statusQuality,
				"Inspector":           statusInspector,
				"GroupP1":             statusGroupP1,
				"GroupP5":             statusGroupP5,
				//	"GET":           get,
				"LoggedIn": statusLoggedIn,
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

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupQuality || user.Groups == groupWarehouse {
			statusGroupP1 = true

			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true
				fmt.Println("кладовщик склада - ", statusStockkeeperWH)
			} else if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
			} else if user.Role == roleInspector {
				//statusQuality = true
				statusInspector = true
				statusLoggedIn = true
			//	fmt.Println("pageInspection quality - ", Quality)
			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
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
					"Admin":               statusAdmin,
					"WarehouseManager":    statusWarehouseManager,
					"StockkeeperWH":       statusStockkeeperWH,
					"SuperIngenerQuality": statusSuperIngenerQuality,
					"IngenerQuality":      statusIngenerQuality,
					//"Quality":             statusQuality,
					"Inspector":           statusInspector,
					"GroupP1":             statusGroupP1,
					//"GroupP5":             GroupP5,
					"LoggedIn": statusLoggedIn,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
					//	"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":            statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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

			statusGroupP5 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true
				fmt.Println("кладовщик склада - ", statusStockkeeperWH)
			} else if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
			} else if user.Role == roleInspector {
				//statusQuality = true
				statusInspector = true
				statusLoggedIn = true
				//fmt.Println("pageInspection quality - ", statusQuality)
			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
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
					"Admin":               statusAdmin,
					"WarehouseManager":    statusWarehouseManager,
					"StockkeeperWH":       statusStockkeeperWH,
					"SuperIngenerQuality": statusSuperIngenerQuality,
					"IngenerQuality":      statusIngenerQuality,
					//"Quality":             statusQuality,
					"Inspector":           statusInspector,
					//	"GroupP1":             GroupP1,
					"GroupP5":  statusGroupP5,
					"LoggedIn": statusLoggedIn,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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

//historyInspectionP5
func (s *Server) PagehistoryInspectionP5() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupQualityP5 || user.Groups == groupWarehouseP5 || user.Groups == groupQuality ||
			user.Groups == groupWarehouse {
			statusGroupP5 = true
			statusGroupP1 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true
				fmt.Println("кладовщик склада - ", statusStockkeeperWH)
			} else if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
			} else if user.Groups == groupQuality {
				//Quality = true
				statusInspector = true
				statusLoggedIn = true
				//fmt.Println("pageInspection quality - ", Quality)
			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
			}
			data := map[string]interface{}{
				"Admin":               statusAdmin,
				"StockkeeperWH":       statusStockkeeperWH,
				"SuperIngenerQuality": statusSuperIngenerQuality,
				"WarehouseManager":    statusWarehouseManager,
				"IngenerQuality":      statusIngenerQuality,
				//"Quality":             statusQuality,
				"Inspector":           statusInspector,
				"GroupP5":             statusGroupP5,
				"GroupP1":             statusGroupP1,
				//	"GET":           get,
				"LoggedIn": statusLoggedIn,
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

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupQualityP5 || user.Groups == groupWarehouseP5 {
			statusGroupP5 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true
				fmt.Println("кладовщик склада - ", statusStockkeeperWH)
			} else if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
			} else if user.Role == roleInspector {
				//Quality = true
				statusInspector = true
				statusLoggedIn = true
				//fmt.Println("pageInspection quality - ", Quality)
			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
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
					"Admin":               statusAdmin,
					"WarehouseManager":    statusWarehouseManager,
					"StockkeeperWH":       statusStockkeeperWH,
					"SuperIngenerQuality": statusSuperIngenerQuality,
					"IngenerQuality":      statusIngenerQuality,
					//"Quality":             Quality,
					"Inspector":           statusInspector,
					"GroupP5":             statusGroupP5,
					"Pobedit":             "Победит 5",
					//"GroupP1":             GroupP1,
					"LoggedIn": statusLoggedIn,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"Pobedit":             "Победит 5",
						//"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"Pobedit":             "Победит 5",
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
					//	"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
					//	"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             statusQuality,
						"Inspector":           statusInspector,
						"LoggedIn":           statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						//	"GroupP1":             GroupP1,
						"GroupP5": statusGroupP5,
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
			statusGroupP1 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true
				fmt.Println("кладовщик склада - ", statusStockkeeperWH)
			} else if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
			} else if user.Role == roleInspector {
				//Quality = true
				statusInspector = true
				statusLoggedIn = true
			//	fmt.Println("pageInspection quality - ", Quality)
			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
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
					"Admin":               statusAdmin,
					"WarehouseManager":    statusWarehouseManager,
					"StockkeeperWH":       statusStockkeeperWH,
					"SuperIngenerQuality": statusSuperIngenerQuality,
					"IngenerQuality":      statusIngenerQuality,
					//"Quality":             statusQuality,
					"Inspector":           statusInspector,
					//	"GroupP5":             GroupP5,
					"GroupP1":  statusGroupP1,
					"Pobedit":  "Победит 5",
					"LoggedIn": statusLoggedIn,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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
						"Admin":               statusAdmin,
						"WarehouseManager":    statusWarehouseManager,
						"StockkeeperWH":       statusStockkeeperWH,
						"SuperIngenerQuality": statusSuperIngenerQuality,
						"IngenerQuality":      statusIngenerQuality,
						//"Quality":             Quality,
						"Inspector":           statusInspector,
						"LoggedIn":            statusLoggedIn,
						"GroupP1":             statusGroupP1,
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

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupQuality || user.Groups == groupWarehouse {
			statusGroupP1 = true
			if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				//statusSuperIngenerQuality2 = true
				//MixP1P5 = true
				statusLoggedIn = true
				fmt.Println("pageInspection SuperIngenerQuality - ", statusSuperIngenerQuality)
			} else if user.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				//	Warehouse = false
				//	WarehouseManager = true
				statusLoggedIn = true
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
			} else if user.Role == roleInspector {
				statusInspector = true
			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
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
				"Inspector":        statusInspector,
				"IngenerQuality":   statusIngenerQuality,
				"WarehouseManager": statusWarehouseManager,
				//	"Warehouse":            Warehouse,
				"StockkeeperWH":         statusStockkeeperWH,
				"SuperIngenerQuality":   statusSuperIngenerQuality,
				//"SuperIngenerQuality2":  statusSuperIngenerQuality2,
				"GroupP1":               statusGroupP1,
				//"MixP1P5":               MixP1P5,
				"GET":                   get,
				"ListVendor":            listVendor,
				"CountTotal":            countTotal,
				"HoldInspection":        holdInspection,
				"NotVerifyComponents":   notVerifyComponents,
				"CountVerifyComponents": countVerifyComponents,
				"GetStatic":             getStatic,
				"HoldCountDebitor":      holdCountDebitor,
				"NotVerifyDebitor":      notVerifyDebitor,
				"LoggedIn":              statusLoggedIn,
				"Pobedit":               "Победит 1",
			}

			tpl.ExecuteTemplate(w, "showinspection.html", groups)
			//	json.NewEncoder(w).Encode(get)
		}

		if user.Groups == groupQualityP5 || user.Groups == groupWarehouseP5 {
			statusGroupP5 = true
			if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				//statusSuperIngenerQuality2 = true
				statusLoggedIn = true
				fmt.Println("pageInspection SuperIngenerQuality - ", statusSuperIngenerQuality)
			} else if user.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				//	Warehouse = false
				//	WarehouseManager = true
				statusLoggedIn = true
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
			} else if user.Role == roleInspector {
				statusInspector = true
			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
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
				"Inspector":        statusInspector,
				"IngenerQuality":   statusIngenerQuality,
				"WarehouseManager": statusWarehouseManager,
				//	"Warehouse":            Warehouse,
				"StockkeeperWH":         statusStockkeeperWH,
				"SuperIngenerQuality":   statusSuperIngenerQuality,
				//"SuperIngenerQuality2":  SuperIngenerQuality2,
				"GroupP5":               statusGroupP5,
				"GET":                   get,
				"ListVendor":            listVendor,
				"CountTotal":            countTotal,
				"HoldInspection":        holdInspection,
				"NotVerifyComponents":   notVerifyComponents,
				"GetStatic":             getStatic,
				"HoldCountDebitor":      holdCountDebitor,
				"NotVerifyDebitor":      notVerifyDebitor,
				"CountVerifyComponents": countVerifyComponents,
				"LoggedIn":              statusLoggedIn,
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

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupQuality && user.Role == roleSuperIngenerQuality ||
			user.Groups == groupQuality && user.Role == roleIngenerQuality {
			//	MixP1P5 = true
			statusGroupP1 = true
			if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				//SuperIngenerQuality2 = true
				statusLoggedIn = true
				fmt.Println("pageInspection SuperIngenerQuality - ", statusSuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
			} /* else if user.Role == "кладовщик склада" {
				StockkeeperWH = true
				//	Warehouse = false
				//	WarehouseManager = true
				LoggedIn = true
			} */
			/* else if user.Role == "контролер качества" {
				Inspector = true
			} else if user.Role == "старший кладовщик склада" {
				WarehouseManager = true
				LoggedIn = true
			} */
			/* else if user.Role == "Administrator" {
				Admin = true
				LoggedIn = true
			}*/ /**else if user.Groups == "качество" {
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
				"Inspector":        statusInspector,
				"IngenerQuality":   statusIngenerQuality,
				"WarehouseManager": statusWarehouseManager,
				//	"Warehouse":            Warehouse,
				"StockkeeperWH":        statusStockkeeperWH,
				"SuperIngenerQuality":  statusSuperIngenerQuality,
				//"SuperIngenerQuality2": SuperIngenerQuality2,
				//"MixP1P5":              MixP1P5,
				"GroupP1":               statusGroupP1,
				"GET":                   get,
				"ListVendor":            listVendor,
				"CountTotal":            countTotal,
				"HoldInspection":        holdInspection,
				"NotVerifyComponents":   notVerifyComponents,
				"GetStatic":             getStatic,
				"HoldCountDebitor":      holdCountDebitor,
				"NotVerifyDebitor":      notVerifyDebitor,
				"CountVerifyComponents": countVerifyComponents,
				"LoggedIn":              statusLoggedIn,
				"Pobedit":               "Победит 5",
			}

			tpl.ExecuteTemplate(w, "showinspectionmix.html", groups)
			//	json.NewEncoder(w).Encode(get)
		}

		if user.Groups == groupQualityP5 && user.Role == roleSuperIngenerQuality ||
			user.Groups == groupQualityP5 && user.Role == roleIngenerQuality {
			//	MixP1P5 = true
			statusGroupP5 = true
			if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				//SuperIngenerQuality2 = true
				statusLoggedIn = true
				fmt.Println("pageInspection SuperIngenerQuality - ", statusSuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
			} /* else if user.Role == "кладовщик склада" {
				StockkeeperWH = true
				//	Warehouse = false
				//	WarehouseManager = true
				LoggedIn = true
			} */
			/* else if user.Role == "контролер качества" {
				Inspector = true
			} else if user.Role == "старший кладовщик склада" {
				WarehouseManager = true
				LoggedIn = true
			} */
			/* else if user.Role == "Administrator" {
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
				"Inspector":        statusInspector,
				"IngenerQuality":   statusIngenerQuality,
				"WarehouseManager": statusWarehouseManager,
				//	"Warehouse":            Warehouse,
				"StockkeeperWH":        statusStockkeeperWH,
				"SuperIngenerQuality":  statusSuperIngenerQuality,
				//"SuperIngenerQuality2": SuperIngenerQuality2,
				//"MixP1P5":              MixP1P5,
				"GroupP5":               statusGroupP5,
				"GET":                   get,
				"ListVendor":            listVendor,
				"CountTotal":            countTotal,
				"HoldInspection":        holdInspection,
				"NotVerifyComponents":   notVerifyComponents,
				"GetStatic":             getStatic,
				"HoldCountDebitor":      holdCountDebitor,
				"NotVerifyDebitor":      notVerifyDebitor,
				"CountVerifyComponents": countVerifyComponents,
				"LoggedIn":              statusLoggedIn,
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
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}

		user := r.Context().Value(ctxKeyUser).(*model.User)
		fmt.Println("user.Groups - ?", user.Groups)
		//s.infoLog.Printf("test user.Groups - %s", user.Groups)
		s.logger.Infof("test user.Groups - %s", user.Groups)

		if user.Groups == groupQuality || user.Groups == roleAdministrator {
			statusGroupP1 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
				fmt.Println("SuperIngenerQuality - ", statusSuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
				fmt.Println("IngenerQuality - ", statusIngenerQuality)
			} else if user.Role == roleInspector {
				statusInspector = true
				statusLoggedIn = true

			}
			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditInspection(id)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			data := map[string]interface{}{
				"Admin":               statusAdmin,
				"SuperIngenerQuality": statusSuperIngenerQuality,
				"IngenerQuality":      statusIngenerQuality,
				"Inspector":           statusInspector,
				"GET":                 get,
				"LoggedIn":            statusLoggedIn,
				"GroupP1":             statusGroupP1,
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
			statusGroupP5 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
				fmt.Println("SuperIngenerQuality - ", statusSuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
				fmt.Println("IngenerQuality - ", statusIngenerQuality)
			} else if user.Role == roleInspector {
				statusInspector = true
				statusLoggedIn = true

			}
			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditInspectionP5(id)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			data := map[string]interface{}{
				"Admin":               statusAdmin,
				"SuperIngenerQuality": statusSuperIngenerQuality,
				"IngenerQuality":      statusIngenerQuality,
				"Inspector":           statusInspector,
				"GET":                 get,
				"LoggedIn":            statusLoggedIn,
				"GroupP5":             statusGroupP5,
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
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}

		user := r.Context().Value(ctxKeyUser).(*model.User)
		fmt.Println("user.Groups - ?", user.Groups)
		fmt.Println("Test json page update")
		//s.infoLog.Printf("user.Groups - %v", user.Groups)
		s.logger.Infof("user.Groups - %v", user.Groups)
		/*
			if (user.Groups == "качество" && user.Role == "главный инженер по качеству") ||
				(user.Groups == "качество" && user.Role == "инженер по качеству") {
				fmt.Println("Test mix page update")
				GroupP1 = true
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
				} /*else if user.Role == "контролер качества" {
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
			statusGroupP1 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
				fmt.Println("SuperIngenerQuality - ", statusSuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
				fmt.Println("IngenerQuality - ", statusIngenerQuality)
			} else if user.Role == roleInspector {
				statusInspector = true
				statusLoggedIn = true

			}
			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditInspection(id)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			data := map[string]interface{}{
				"Admin":               statusAdmin,
				"SuperIngenerQuality": statusSuperIngenerQuality,
				"IngenerQuality":      statusIngenerQuality,
				"Inspector":           statusInspector,
				"GroupP1":             statusGroupP1,
				"GET":                 get,
				"LoggedIn":            statusLoggedIn,
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
			statusGroupP5 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
				fmt.Println("SuperIngenerQuality - ", statusSuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
				fmt.Println("IngenerQuality - ", statusIngenerQuality)
			} else if user.Role == roleInspector {
				statusInspector = true
				statusLoggedIn = true

			}
			//fmt.Println("ID - ?", id)

			get, err := s.store.Inspection().EditInspectionP5(id)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			data := map[string]interface{}{
				"Admin":               statusAdmin,
				"SuperIngenerQuality": statusSuperIngenerQuality,
				"IngenerQuality":      statusIngenerQuality,
				"Inspector":           statusInspector,
				"GroupP5":             statusGroupP5,
				"GET":                 get,
				"LoggedIn":            statusLoggedIn,
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
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}

		user := r.Context().Value(ctxKeyUser).(*model.User)
		fmt.Println("user.Groups - ?", user.Groups)
		fmt.Println("Test json page update")
		//s.infoLog.Printf("user.Groups - %v", user.Groups)
		s.logger.Infof("user.Groups - %v", user.Groups)

		if (user.Groups == groupQuality && user.Role == roleSuperIngenerQuality) ||
			(user.Groups == groupQuality && user.Role == roleSuperIngenerQuality) {
			fmt.Println("Test mix page update")
			statusGroupP1 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
				fmt.Println("SuperIngenerQuality - ", statusSuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
				fmt.Println("IngenerQuality - ", statusIngenerQuality)
			} /*else if user.Role == "контролер качества" {
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
				"Admin":               statusAdmin,
				"SuperIngenerQuality": statusSuperIngenerQuality,
				"IngenerQuality":      statusIngenerQuality,
				"Inspector":           statusInspector,
				"GroupP1":             statusGroupP1,
				"GET":                 get,
				"LoggedIn":            statusLoggedIn,
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
			statusGroupP5 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
			} else if user.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
				fmt.Println("SuperIngenerQuality - ", statusSuperIngenerQuality)
			} else if user.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
				fmt.Println("IngenerQuality - ", statusIngenerQuality)
			} /*else if user.Role == "контролер качества" {
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
				"Admin":               statusAdmin,
				"SuperIngenerQuality": statusSuperIngenerQuality,
				"IngenerQuality":      statusIngenerQuality,
				"Inspector":           statusInspector,
				"GroupP5":             statusGroupP5,
				"GET":                 get,
				"LoggedIn":            statusLoggedIn,
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

		//GroupP1 = "качество"
		//GroupP5 = "качество"

		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}

		currentTime := time.Now()

		user := r.Context().Value(ctxKeyUser).(*model.User)

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

		user := r.Context().Value(ctxKeyUser).(*model.User)

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
			/*	if user.Groups == "качество" && user.Role == "главный инженер по качеству" ||
				user.Groups == "качество" && user.Role == "инженер по качеству" {
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

		user := r.Context().Value(ctxKeyUser).(*model.User)

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

		user := r.Context().Value(ctxKeyUser).(*model.User)

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

		user := r.Context().Value(ctxKeyUser).(*model.User)

		if user.Groups == groupWarehouse {
			statusStockkeeperWH = true
			statusWarehouseManager = true
		}
		/*
			if user.Groups == "качество" {
				quality = true
			} else if user.Groups == "склад" {
				stockkeeperWH = true
			}
		*/
		if user.Role == roleAdministrator {
			statusAdmin = true
			statusLoggedIn = true
		} else if user.Role == roleStockkeeperWH {
			statusStockkeeperWH = true
			statusLoggedIn = true
			fmt.Println("кладовщик склада - ", statusStockkeeperWH)
		} else if user.Role == roleWarehouseManager {
			statusWarehouseManager = true
			statusLoggedIn = true
		}

		get, err := s.store.Inspection().ListAcceptWHInspection()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		groups := map[string]interface{}{
			//	"quality":   quality,
			//"Warehouse":        stockkeeperWH,
			"WarehouseManager": statusWarehouseManager,
			//	"главный инженер по качеству": superIngenerQuality,
			"GET": get,
			//	"status":    statusStr,
			"Admin":         statusAdmin,
			"StockkeeperWH": statusStockkeeperWH,
			"LoggedIn":      statusLoggedIn,
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

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}

		user := r.Context().Value(ctxKeyUser).(*model.User)
		fmt.Println("user.Groups - ?", user.Groups)
		//s.infoLog.Printf("user.Groups - %v\n", user.Groups)
		s.logger.Infof("user.Groups - %v\n", user.Groups)

		if user.Groups == groupWarehouse {
			statusGroupP1 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true

			} else if user.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true

			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
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
				"Admin":            statusAdmin,
				"WarehouseManager": statusWarehouseManager,
				"StockkeeperWH":    statusStockkeeperWH,
				"GroupP1":          statusGroupP1,
				"LoggedIn":         statusLoggedIn,
				"GET":              get,
			}
			err = tpl.ExecuteTemplate(w, "acceptWarehouseInspection.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
		if user.Groups == groupWarehouseP5 {
			statusGroupP5 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true

			} else if user.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true

			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
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
				"Admin":            statusAdmin,
				"WarehouseManager": statusWarehouseManager,
				"StockkeeperWH":    statusStockkeeperWH,
				"GroupP5":          statusGroupP5,
				"LoggedIn":         statusLoggedIn,
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

		user := r.Context().Value(ctxKeyUser).(*model.User)

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

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
			s.logger.Errorf(err.Error())
		}

		user := r.Context().Value(ctxKeyUser).(*model.User)
		fmt.Println("user.Groups - ?", user.Groups)
		//s.infoLog.Printf("user.Groups - %v\n", user.Groups)
		s.logger.Info("user.Groups - %v\n", user.Groups)

		if user.Groups == groupWarehouse {
			statusGroupP1 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true

			} else if user.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true

			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
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
				"Admin":            statusAdmin,
				"WarehouseManager": statusWarehouseManager,
				"StockkeeperWH":    statusStockkeeperWH,
				"GroupP1":          statusGroupP1,
				"LoggedIn":         statusLoggedIn,
				"GET":              get,
			}
			err = tpl.ExecuteTemplate(w, "acceptWarehouseInspectionjson.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
		if user.Groups == groupWarehouseP5 {
			statusGroupP5 = true
			if user.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true

			} else if user.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true

			} else if user.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
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
				"Admin":            statusAdmin,
				"WarehouseManager": statusWarehouseManager,
				"StockkeeperWH":    statusStockkeeperWH,
				"GroupP5":          statusGroupP5,
				"LoggedIn":         statusLoggedIn,
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

		user := r.Context().Value(ctxKeyUser).(*model.User)
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
			if user.Groups == "склад" {

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
			if user.Groups == "склад П5" {
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

		user := r.Context().Value(ctxKeyUser).(*model.User)

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
