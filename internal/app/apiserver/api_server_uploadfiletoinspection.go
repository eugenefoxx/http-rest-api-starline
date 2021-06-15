package apiserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"

	_ "github.com/gorilla/mux"
)

func (s *Server) UploadFileToInspection() http.HandlerFunc {
	// response format
	type response struct {
		//	ID      int64  `json:"id,omitempty"`
		//	Status  string `json:"status,omitempty"`
		//	Note    string `json:"note,omitempty"`
		Message string `json:"message,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		/*
			w.Header().Set("Content-Type", "application/json")
		https://market-click2.yandex.ru/redir/GAkkM7lQwz62j9BQ6_qgZoep55uRD4T4A5aShDSjUBi68pDQCp1He6qhpWv2RvRNoamhkBdFBph70o8oA-5v9QGTf7JG5pnZ0SxlzCQYReb2bF_tn3F_Ed3CyF0oSFoCRAH-937jxBq2JUxNs50BMfViko7mGN-XBOFiSj2jFcNyuhBmwvRNLvkyl6K7tY_lXjoBWjQ6evL809qJ4bwvHy9VvSxSZLlSoFvVKVXrvH-sIfQQp11YyUsFkPCYgcmVVWOcD_CVF_MsIO8jgMsoqcrT8X4r6B_DrWhxHzw7twaYp1ewGltJCUtYQkIkLNfj6tAOowY8XHzzHe-Sf179_46iw9CxMGQaYL2w4wZJdrdkq3TqPcmSMQgnLbj8lPrP2F_FFxf8qCLalSOLFxV4dbnz5gJXL1Dop-5A37H2D1d2oJDLL48FE4YuLkfV4E_UliMzTaEO7rgs8B2AJJre7jh71xtJjctJfiIUsTBWBz82W8fNY3mqQp5hzJSVrLta9hLeJQ7_G6YsgHNvV4-yNT_5TzQgwmXsVGjkQ86__KIW1XIvUgCOMLkdeNKXHbPS-D1oQHfzN2VN44cM2ropq8t4vU5ztny575SvLUHYqxPoaU31BBMCdi_MQV-42Lgx3aa00knxMIAqcH1iN7RjKyi8uUvDGdLS8OFhfD1p2M8gWAIsENoG8i4HT3yCE_KcxCk7BLG1Bf6gJfLYEPmR-UpwJa_K4yjb62sJlapGVDm6ZcaFdeEuPtz2fRwXXT6VXgkcO5q4-kmnZMml_yETkyffFytNyzzb7fIqy1PurZxhDagk1svk3oTIT-9G7uQs1moXic4P9_3LNRpds4lGT2vpQa0f_TBzoGAO8wS0ruTT83B5cmj3JtMF4I4yuMa4ZNLmOAdY7pujn-ocnZXK1pfpKZoQ2uw46YU4XWio4WKUkeV1XCQ5TWhahIQ0bZ7klpYyqxLywoEpQEt40rByX_u--LOpsP6JSr10t3M5M7bFDPlSRk96FUOIKHIOsg8QPrZEOSU0J2V3JsuR-SyENcjTmPofsJh3YXPBBDbGNVwFJE6M5RPn_-TRx1jQWJ18007_BIONWolqwLIfESqcpt6idl2gICp85EtYS4lbqUNQwdZW08P_GImZNdqwDHebyCTpGxsbmqpit-uIbyYEo7yHURHlGYsp5g07C1SrQvGku2aN0s31D-qvfQVF7c6KdvaKUgklCEFqiAATId9QwL-MWVe_64tH-_PenXfqFExRBF_wB9o1m5KapQiBxFwoSr4m7tJLCvHJOcn3AmE7Xkiwxcn_rCIAUMDzhNvy8w0UVyQdOzIvbvp5r5e08WcasVQVuwbbLa-Ul96AD-umD7I5YGSpmOdB5kHdRIR3uKA2Z4e91W9UnA,,?data=QVyKqSPyGQwNvdoowNEPjcozEqKShkAf7BQtg3MVmxezzWyQ_AnAou0G0CoA2MaUbQhYPmbQ2kCA3q7BRlWBwyjfb1mkXd2WnEgq7VMx6hmefayHx1cpQQV-lWBPVSIopuSJ25hp-eekJl5mU5kqwKcM8lQ4OBdQu4hWqtlETYcUn0RoY7lrcP2C6oFce5VVg7CcEWjy0uo8I9pN3FpZqcKgVnxhgwxM_jeH4d_ZqLqp1UR8xN2dvE_WCkSQdOSO&b64e=1&sign=11c58b9a57be14e2b4bc6c1614ad669c&keyno=1			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		*/
		user := r.Context().Value(ctxKeyUser).(*model.User)

		const statusTransfer = "отгружено на ВК"

		// чтение данных из формы
		f, h, err := r.FormFile("q")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer f.Close()

		// for your information
		fmt.Println("\nfile:", f, "\nheader:", h, "\nerr:", err)

		// read
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Println("bs", bs)
		ss := string(bs)

		fmt.Println("ss", ss)
		// запись данных из формы
		f2, err := os.Create("data")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer f2.Close()
		n2, err := f2.WriteString(ss)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("n2", n2)
		f2.Sync()

		/*
			file, err := os.Open("data2")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close() */
		// читаем файл в массив строк
		p, err := readLines("data")
		if err != nil {
			fmt.Println(err)
		}
		// передаем на проверку дублей
		pp := removeDuplicatesinfile(p)

		// запись массива в файл
		file, err := os.OpenFile("data2", os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}

		datawriter := bufio.NewWriter(file)

		for _, data := range pp {
			_, _ = datawriter.WriteString(data + "\n")
		}

		datawriter.Flush()
		file.Close()
		// конечный результат
		file2, err := os.Open("data2")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file2)
		scanner.Split(bufio.ScanLines)

		// This is our buffer now
		var lines []string

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		err = os.Remove("data")
		if err != nil {
			log.Println(err)
		}
		err = os.Remove("data2")
		if err != nil {
			log.Println(err)
		}

		fmt.Println("read lines:")
		for i, line := range lines {
			//	fmt.Println(line)
			//	fmt.Printf("%v %v\n", i, line)
			s.logger.Infof("%v %v\n", i, line)
			//	fmt.Println(line[0:45])
			//	fmt.Println(line[0:1])
			//fmt.Println("Test", len(line))
			re := regexp.MustCompile(`P\d{7}LK\d{9}R\d{10}Q\d{5}D\d{8}`)
			re2 := regexp.MustCompile(`P\d{7}L\d{10}R\d{10}Q\d{5}D\d{8}`)
			re3 := regexp.MustCompile(`P\d{7}LR\d{10}Q\d{5}D\d{8}`)
			if re.MatchString(line) != true && re2.MatchString(line) != true && re3.MatchString(line) != true {
				fmt.Println("не верное сканирование MatchString :\n" + line + "\n")
				s.logger.Errorf(
					"Не верное сканирование из файла загрузки на входной контроль : %v," +
						"табель сканировавшего сотрудника %v", line , user.Tabel)
				//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
				w.Write([]byte("Запись не соответствует: " + line))

				return
			}

		}

		if user.Groups == groupWarehouse {
			for _, v := range lines {
				if (strings.Contains(v[0:1], "P") == true) && (len(v) == 45) {
					idMaterial := v[0:45]

					fmt.Println("Пропускаем:\n" + idMaterial + "\n")
					s.logger.Infof("Запись строки на загрузку с файла на входной контроль, П1: %v", idMaterial)
					sapStr := v[1:8]
					var sap int
					sap, err := strconv.Atoi(sapStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					idrollStr := v[20:30]
					var idrollIns int
					idrollIns, err = strconv.Atoi(idrollStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var lot string
					lot = v[9:19]
					qtyStr := v[31:36]
					var qtyIns int
					qtyIns, err = strconv.Atoi(qtyStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var productionDate string
					productionDate = v[37:45]
					var numberVendor string
					numberVendor = v[9:15]
					fmt.Println("v.NumberVendor", numberVendor)
					if (strings.Contains(v[0:1], "P") == true) && (len(v) == 45) {
						u := &model.Inspection{
							IdMaterial:     idMaterial,
							SAP:            sap,
							Lot:            lot,
							IdRoll:         idrollIns,
							Qty:            qtyIns,
							ProductionDate: productionDate,
							NumberVendor:   numberVendor,
							Location:       statusTransfer,
							Lastname:       user.LastName,
						}

						if err := s.store.Inspection().InInspection(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)
							return
						}

					}
					//	} else {
					//		if idMaterial == idMaterial {
					//			fmt.Println("Значения совпадают:\n" + idMaterial + "\n")
					//		}
					//	}
					//http.Redirect(w, r, "/operation/statusinspection", 303)
				} else {
					if (strings.Contains(v[0:1], "P") == false) && (len(v) != 45) {
						fmt.Println("не верное сканирование :\n" + v + "\n")
						s.logger.Errorf("Не верное сканирование из файла загрузки на входной контроль, П1: %v", v)
						//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
						w.Write([]byte("Запись не соответствует: " + v))
						return
					}
					//	tpl.Execute(w, data)

				}
				if (strings.Contains(v[0:1], "P") == true) && (len(v) == 35) {
					idMaterial := v[0:35]

					fmt.Println("Пропускаем:\n" + idMaterial + "\n")
					s.logger.Infof("Запись строки на загрузку с файла на входной контроль, П1: %v", idMaterial)
					sapStr := v[1:8]
					var sap int
					sap, err := strconv.Atoi(sapStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					idrollStr := v[10:20]
					var idrollIns int
					idrollIns, err = strconv.Atoi(idrollStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var lot string
					lot = "без партии" //v.ScanID[9:19]
					qtyStr := v[21:26]
					var qtyIns int
					qtyIns, err = strconv.Atoi(qtyStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var productionDate string
					productionDate = v[27:35]
					var numberVendor string
					numberVendor = "без поставщика" //v.ScanID[9:15]
					fmt.Println("v.NumberVendor", numberVendor)
					if (strings.Contains(v[0:1], "P") == true) && (len(v) == 35) {
						u := &model.Inspection{
							IdMaterial:     idMaterial,
							SAP:            sap,
							Lot:            lot,
							IdRoll:         idrollIns,
							Qty:            qtyIns,
							ProductionDate: productionDate,
							NumberVendor:   numberVendor,
							Location:       statusTransfer,
							Lastname:       user.LastName,
						}

						if err := s.store.Inspection().InInspection(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)
							return
						}

					}
					//	} else {
					//		if idMaterial == idMaterial {
					//			fmt.Println("Значения совпадают:\n" + idMaterial + "\n")
					//		}
					//	}
					//http.Redirect(w, r, "/operation/statusinspection", 303)
				} else {
					if (strings.Contains(v[0:1], "P") == false) && (len(v) != 35) {
						fmt.Println("не верное сканирование :\n" + v + "\n")
						s.logger.Errorf("Не верное сканирование из файла загрузки на входной контроль, П1: %v", v)
						//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
						w.Write([]byte("Запись не соответствует: " + v))
						return
					}
					//	tpl.Execute(w, data)

				}
			}
			w.Write([]byte("Файл успешно загружен"))
		}
		fmt.Println("Test1")
		if user.Groups == groupWarehouseP5 {
			fmt.Println("Test2")
			for _, v := range lines {
				fmt.Println("Test3")
				fmt.Println("Test4 - ", v)
				if (strings.Contains(v[0:1], "P") == true) && (len(v) == 45) {
					idMaterial := v[0:45]

					fmt.Println("Пропускаем:\n" + idMaterial + "\n")
					s.logger.Infof("Запись строки на загрузку с файла на входной контроль, П5: %v", idMaterial)
					sapStr := v[1:8]
					var sap int
					sap, err := strconv.Atoi(sapStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					idrollStr := v[20:30]
					var idrollIns int
					idrollIns, err = strconv.Atoi(idrollStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var lot string
					lot = v[9:19]
					qtyStr := v[31:36]
					var qtyIns int
					qtyIns, err = strconv.Atoi(qtyStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var productionDate string
					productionDate = v[37:45]
					var numberVendor string
					numberVendor = v[9:15]
					fmt.Println("v.NumberVendor", numberVendor)
					if (strings.Contains(v[0:1], "P") == true) && (len(v) == 45) {
						u := &model.Inspection{
							IdMaterial:     idMaterial,
							SAP:            sap,
							Lot:            lot,
							IdRoll:         idrollIns,
							Qty:            qtyIns,
							ProductionDate: productionDate,
							NumberVendor:   numberVendor,
							Location:       statusTransfer,
							Lastname:       user.LastName,
						}

						if err := s.store.Inspection().InInspectionP5(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)

							return
						}

					}
					//	} else {
					//		if idMaterial == idMaterial {
					//			fmt.Println("Значения совпадают:\n" + idMaterial + "\n")
					//		}
					//	}
					//http.Redirect(w, r, "/operation/statusinspection", 303)

				} else {
					if (strings.Contains(v[0:1], "P") == false) && (len(v) != 45) {
						fmt.Println("не верное сканирование :" + v)
						s.logger.Errorf("Не верное сканирование из файла загрузки на входной контроль, П5: %v", v)
						//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
						w.Write([]byte("Запись не соответствует: " + v))
						return
					}
					//	tpl.Execute(w, data)

				}
				if (strings.Contains(v[0:1], "P") == true) && (len(v) == 35) {
					idMaterial := v[0:35]

					fmt.Println("Пропускаем:\n" + idMaterial + "\n")
					s.logger.Infof("Запись строки на загрузку с файла на входной контроль, П5: %v", idMaterial)
					sapStr := v[1:8]
					var sap int
					sap, err := strconv.Atoi(sapStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					idrollStr := v[10:20]
					var idrollIns int
					idrollIns, err = strconv.Atoi(idrollStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var lot string
					lot = "без партии" //v.ScanID[9:19]
					qtyStr := v[21:26]
					var qtyIns int
					qtyIns, err = strconv.Atoi(qtyStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var productionDate string
					productionDate = v[27:35]
					var numberVendor string
					numberVendor = "без поставщика" //v.ScanID[9:15]
					fmt.Println("v.NumberVendor", numberVendor)
					if (strings.Contains(v[0:1], "P") == true) && (len(v) == 35) {
						u := &model.Inspection{
							IdMaterial:     idMaterial,
							SAP:            sap,
							Lot:            lot,
							IdRoll:         idrollIns,
							Qty:            qtyIns,
							ProductionDate: productionDate,
							NumberVendor:   numberVendor,
							Location:       statusTransfer,
							Lastname:       user.LastName,
						}

						if err := s.store.Inspection().InInspection(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)

							return
						}

					}
					//	} else {
					//		if idMaterial == idMaterial {
					//			fmt.Println("Значения совпадают:\n" + idMaterial + "\n")
					//		}
					//	}
					http.Redirect(w, r, "/operation/statusinspection", 303)
				} else {
					if (strings.Contains(v[0:1], "P") == false) && (len(v) != 35) {
						fmt.Println("не верное сканирование :\n" + v + "\n")
						s.logger.Errorf("Не верное сканирование из файла загрузки на входной контроль, П5: %v", v)
						//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
						w.Write([]byte("Запись не соответствует: " + v))
						return
					}
					//	tpl.Execute(w, data)

				}
			}
			w.Write([]byte("Файл успешно загружен"))
		}
		/*
			for i, y := range s {
				//ss := [i]
				//	fmt.Println(ss)
				fmt.Printf("%v %v", i, y)
			}*/
	}

}

// PageUploadFileToInspectionJSON
func (s *Server) PageUploadFileToInspectionJSON() http.HandlerFunc {
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
		//	SuperIngenerQuality := false
		//	IngenerQuality := false
		//	Inspector := false
		StockkeeperWH := false
		WarehouseManager := false
		GroupP1 := false
		GroupP5 := false
		LoggedIn := false
		/*
			vars := mux.Vars(r)
			id, err := strconv.Atoi(vars["ID"])
			if err != nil {
				log.Println(err)
				s.logger.Errorf(err.Error())
			}
		*/
		user := r.Context().Value(ctxKeyUser).(*model.User)

		fmt.Println("user.Groups - ?", user.Groups)
		fmt.Println("Test json page update")
		s.logger.Infof("user.Groups - %v", user.Groups)
		if user.Groups == groupWarehouse {
			StockkeeperWH = true
			WarehouseManager = true
		}

		if user.Groups == groupQuality {
			GroupP1 = true
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
			/*else if user.Role == roleSuperIngenerQuality {
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

			}*/
			//fmt.Println("ID - ?", id)
			/*
				get, err := s.store.Inspection().EditInspection(id)
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
			*/
			data := map[string]interface{}{
				"Admin":            Admin,
				"WarehouseManager": WarehouseManager,
				"StockkeeperWH":    StockkeeperWH,
				//	"SuperIngenerQuality": SuperIngenerQuality,
				//	"IngenerQuality":      IngenerQuality,
				//	"Inspector":           Inspector,
				"GroupP1": GroupP1,
				//"GET":      get,
				"LoggedIn": LoggedIn,
				"User":     user.LastName,
				"Username": user.FirstName,
			}
			err := tpl.ExecuteTemplate(w, "uploadfile.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
		if user.Groups == groupQualityP5 {
			GroupP5 = true
			fmt.Println("Test Upload Page")
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
			/*else if user.Role == roleSuperIngenerQuality {
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

			}*/
			//fmt.Println("ID - ?", id)
			/*
				get, err := s.store.Inspection().EditInspectionP5(id)
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
			*/
			data := map[string]interface{}{
				"Admin":            Admin,
				"WarehouseManager": WarehouseManager,
				"StockkeeperWH":    StockkeeperWH,
				//	"SuperIngenerQuality": SuperIngenerQuality,
				//	"IngenerQuality":      IngenerQuality,
				//	"Inspector":           Inspector,
				"GroupP5": GroupP5,
				//"GET":      get,
				"LoggedIn": LoggedIn,
				"User":     user.LastName,
				"Username": user.FirstName,
			}
			err := tpl.ExecuteTemplate(w, "uploadfile.html", data)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}

	}
}

func (s *Server) UploadFileToInspectionJSON() http.HandlerFunc {
	// response format
	type response struct {
		//	ID      int64  `json:"id,omitempty"`
		//	Status  string `json:"status,omitempty"`
		//	Note    string `json:"note,omitempty"`
		Message string `json:"message,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		user := r.Context().Value(ctxKeyUser).(*model.User)

		const statusTransfer = "отгружено на ВК"
		fmt.Println("test")
		// чтение данных из формы
		f, h, err := r.FormFile("q")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer f.Close()

		// for your information
		fmt.Println("\nfile:", f, "\nheader:", h, "\nerr:", err)

		// read
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Println(bs)
		ss := string(bs)

		fmt.Println(ss)
		// запись данных из формы
		f2, err := os.Create("data")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer f2.Close()
		n2, err := f2.WriteString(ss)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(n2)
		f2.Sync()

		/*
			file, err := os.Open("data2")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close() */
		// читаем файл в массив строк
		p, err := readLines("data")
		if err != nil {
			fmt.Println(err)
		}
		// передаем на проверку дублей
		pp := removeDuplicatesinfile(p)

		// запись массива в файл
		file, err := os.OpenFile("data2", os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}

		datawriter := bufio.NewWriter(file)

		for _, data := range pp {
			_, _ = datawriter.WriteString(data + "\n")
		}

		datawriter.Flush()
		file.Close()
		// конечный результат
		file2, err := os.Open("data2")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file2)
		scanner.Split(bufio.ScanLines)

		// This is our buffer now
		var lines []string

		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		err = os.Remove("data")
		if err != nil {
			log.Println(err)
		}
		err = os.Remove("data2")
		if err != nil {
			log.Println(err)
		}
		fmt.Println("read lines:")
		for i, line := range lines {
			//	fmt.Println(line)
			//	fmt.Printf("%v %v\n", i, line)
			s.logger.Infof("%v %v\n", i, line)
			//	fmt.Println(line[0:45])
			//	fmt.Println(line[0:1])
			//fmt.Println("Test", len(line))
			re := regexp.MustCompile(`P\d{7}LK\d{9}R\d{10}Q\d{5}D\d{8}`)
			re2 := regexp.MustCompile(`P\d{7}L\d{10}R\d{10}Q\d{5}D\d{8}`)
			re3 := regexp.MustCompile(`P\d{7}LR\d{10}Q\d{5}D\d{8}`)
			if re.MatchString(line) != true && re2.MatchString(line) != true && re3.MatchString(line) != true {
				fmt.Println("не верное сканирование MatchString :\n" + line + "\n")
				s.logger.Errorf("Не верное сканирование из файла загрузки на входной контроль, П1: %v", line)
				//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
				//	w.Write([]byte("Запись не соответствует: " + line))
				//	break
				//	return
				msg := "Запись не соответствует: " + line
				// format the response message
				res := response{
					//	ID:      int64(idRoll),
					//	Status:  v.Status, // v.(string),
					//	Note:    v.Note,   // v.(string),
					Message: msg,
				}
				// send the response
				json.NewEncoder(w).Encode(res)
				return
			}

		}

		if user.Groups == groupWarehouse || user.Groups == groupQuality {
			for _, v := range lines {
				if (strings.Contains(v[0:1], "P") == true) && (len(v) == 45) {
					idMaterial := v[0:45]

					//	fmt.Println("Пропускаем:\n" + idMaterial + "\n")
					sapStr := v[1:8]
					var sap int
					sap, err := strconv.Atoi(sapStr)
					if err != nil {
						fmt.Println(err)

						//	s.logger.Errorf(err.Error())

					}
					idrollStr := v[20:30]
					var idrollIns int
					idrollIns, err = strconv.Atoi(idrollStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var lot string
					lot = v[9:19]
					qtyStr := v[31:36]
					var qtyIns int
					qtyIns, err = strconv.Atoi(qtyStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var productionDate string
					productionDate = v[37:45]
					var numberVendor string
					numberVendor = v[9:15]
					fmt.Println("v.NumberVendor", numberVendor)
					if (strings.Contains(v[0:1], "P") == true) && (len(v) == 45) {
						u := &model.Inspection{
							IdMaterial:     idMaterial,
							SAP:            sap,
							Lot:            lot,
							IdRoll:         idrollIns,
							Qty:            qtyIns,
							ProductionDate: productionDate,
							NumberVendor:   numberVendor,
							Location:       statusTransfer,
							Lastname:       user.LastName,
						}

						if err := s.store.Inspection().InInspection(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)

							return
						}

					} else {
						if (strings.Contains(v[0:1], "P") == false) && (len(v) != 45) {
							fmt.Println("не верное сканирование :\n" + v + "\n")
							s.logger.Errorf("не верное сканирование :\n" + v + "\n")
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
				if (strings.Contains(v[0:1], "P") == true) && (len(v) == 35) {
					idMaterial := v[0:35]

					//	fmt.Println("Пропускаем:\n" + idMaterial + "\n")
					sapStr := v[1:8]
					var sap int
					sap, err := strconv.Atoi(sapStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					idrollStr := v[10:20]
					var idrollIns int
					idrollIns, err = strconv.Atoi(idrollStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var lot string
					lot = "без партии" //v.ScanID[9:19]
					qtyStr := v[21:26]
					var qtyIns int
					qtyIns, err = strconv.Atoi(qtyStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var productionDate string
					productionDate = v[27:35]
					var numberVendor string
					numberVendor = "без поставщика" //v.ScanID[9:15]
					fmt.Println("v.NumberVendor", numberVendor)
					if (strings.Contains(v[0:1], "P") == true) && (len(v) == 35) {
						u := &model.Inspection{
							IdMaterial:     idMaterial,
							SAP:            sap,
							Lot:            lot,
							IdRoll:         idrollIns,
							Qty:            qtyIns,
							ProductionDate: productionDate,
							NumberVendor:   numberVendor,
							Location:       statusTransfer,
							Lastname:       user.LastName,
						}

						if err := s.store.Inspection().InInspection(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)

							return
						}

					} else {
						if (strings.Contains(v[0:1], "P") == false) && (len(v) != 35) {
							fmt.Println("не верное сканирование :\n" + v + "\n")
							s.logger.Errorf("не верное сканирование :\n" + v + "\n")
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
			for _, v := range lines {
				if (strings.Contains(v[0:1], "P") == true) && (len(v) == 45) {
					idMaterial := v[0:45]

					//	fmt.Println("Пропускаем:\n" + idMaterial + "\n")
					sapStr := v[1:8]
					var sap int
					sap, err := strconv.Atoi(sapStr)
					if err != nil {
						fmt.Println(err)

						//	s.logger.Errorf(err.Error())

					}
					idrollStr := v[20:30]
					var idrollIns int
					idrollIns, err = strconv.Atoi(idrollStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var lot string
					lot = v[9:19]
					qtyStr := v[31:36]
					var qtyIns int
					qtyIns, err = strconv.Atoi(qtyStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var productionDate string
					productionDate = v[37:45]
					var numberVendor string
					numberVendor = v[9:15]
					fmt.Println("v.NumberVendor", numberVendor)
					if (strings.Contains(v[0:1], "P") == true) && (len(v) == 45) {
						u := &model.Inspection{
							IdMaterial:     idMaterial,
							SAP:            sap,
							Lot:            lot,
							IdRoll:         idrollIns,
							Qty:            qtyIns,
							ProductionDate: productionDate,
							NumberVendor:   numberVendor,
							Location:       statusTransfer,
							Lastname:       user.LastName,
						}

						if err := s.store.Inspection().InInspectionP5(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)
							msg := "Файл не соответствует."
							// format the response message
							res := response{
								//	ID:      int64(idRoll),
								//	Status:  v.Status, // v.(string),
								//	Note:    v.Note,   // v.(string),
								Message: msg,
							}
							// send the response
							json.NewEncoder(w).Encode(res)
							return
						}
						msg := "Файл успешно загружен."
						// format the response message
						res := response{
							//	ID:      int64(idRoll),
							//	Status:  v.Status, // v.(string),
							//	Note:    v.Note,   // v.(string),
							Message: msg,
						}
						// send the response
						json.NewEncoder(w).Encode(res)
					} else {
						if (strings.Contains(v[0:1], "P") == false) && (len(v) != 45) {
							fmt.Println("не верное сканирование :\n" + v + "\n")
							s.logger.Errorf("не верное сканирование :\n" + v + "\n")
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
					//	http.Redirect(w, r, "/operation/statusinspection", 303)
				}
				if (strings.Contains(v[0:1], "P") == true) && (len(v) == 35) {
					idMaterial := v[0:35]

					//	fmt.Println("Пропускаем:\n" + idMaterial + "\n")
					sapStr := v[1:8]
					var sap int
					sap, err := strconv.Atoi(sapStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					idrollStr := v[10:20]
					var idrollIns int
					idrollIns, err = strconv.Atoi(idrollStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var lot string
					lot = "без партии" //v.ScanID[9:19]
					qtyStr := v[21:26]
					var qtyIns int
					qtyIns, err = strconv.Atoi(qtyStr)
					if err != nil {
						fmt.Println(err)
						s.logger.Errorf(err.Error())
					}
					var productionDate string
					productionDate = v[27:35]
					var numberVendor string
					numberVendor = "без поставщика" //v.ScanID[9:15]
					fmt.Println("v.NumberVendor", numberVendor)
					if (strings.Contains(v[0:1], "P") == true) && (len(v) == 35) {
						u := &model.Inspection{
							IdMaterial:     idMaterial,
							SAP:            sap,
							Lot:            lot,
							IdRoll:         idrollIns,
							Qty:            qtyIns,
							ProductionDate: productionDate,
							NumberVendor:   numberVendor,
							Location:       statusTransfer,
							Lastname:       user.LastName,
						}

						if err := s.store.Inspection().InInspection(u); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)

							return
						}

						msg := "Файл успешно загружен."
						// format the response message
						res := response{
							//	ID:      int64(idRoll),
							//	Status:  v.Status, // v.(string),
							//	Note:    v.Note,   // v.(string),
							Message: msg,
						}
						// send the response
						json.NewEncoder(w).Encode(res)

					} else {
						if (strings.Contains(v[0:1], "P") == false) && (len(v) != 35) {
							fmt.Println("не верное сканирование :\n" + v + "\n")
							s.logger.Errorf("не верное сканирование :\n" + v + "\n")
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
					//	http.Redirect(w, r, "/operation/statusinspection", 303)
				}
			}
		}
		/*
			for i, y := range s {
				//ss := [i]
				//	fmt.Println(ss)
				fmt.Printf("%v %v", i, y)
			}*/
	}

}

func removeDuplicatesinfile(elements []string) []string { // change string to int here if required
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{} // change string to int here if required
	result := []string{}             // change string to int here if required

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
