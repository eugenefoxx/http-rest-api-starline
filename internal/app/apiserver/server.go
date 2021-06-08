package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eugenefoxx/http-rest-api-starline/pkg/logging"
	"github.com/gorilla/handlers"
	"sync"
//	"unsafe"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store_redis"

	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "starline"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
	tpl                         *template.Template
	//LOGFILE                     = "/tmp/apiServer.log"
	// statement to frontend
	statusAdmin bool
	statusStockkeeper bool
	statusWarehouseManager bool
	statusSuperIngenerQuality bool
	statusIngenerQuality bool
	statusStockkeeperWH bool
	statusInspector bool
	statusGroupP1 bool
	statusGroupP5 bool
	statusLoggedIn bool



)

const (
	// variable for role and group
	roleAdministrator = "Administrator"
	roleWarehouseManager = "старший кладовщик склада"
	roleStockkeeper = "кладовщик"
	roleSuperIngenerQuality = "главный инженер по качеству"
	roleIngenerQuality = "инженер по качеству"
	roleStockkeeperWH = "кладовщик склада"
	roleInspector = "контролер качества"

	groupWarehouse = "склад"
	groupQuality = "качество"
	groupEmpty  = ""
	groupAdministrator = "администратор"
	groupWarehouseP5 = "склад П5"
	groupQualityP5 = "качество П5"
)

type ctxKey int8

type Server struct {
	router       *mux.Router
	//logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
	//database     *sqlx.DB
	//	html         string
	mu sync.Mutex
	//	httpServer *http.Server
//	errorLog *log.Logger
//	infoLog  *log.Logger
	redis    store_redis.Redis
	logger logging.Logger
}

func init() {
	tpl = template.Must(tpl.ParseGlob("web/templates/*.html"))

	//tpl = template.Must(tpl.ParseGlob("/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/web/templates/*.html"))
}

//func newServer(store store.Store, sessionStore sessions.Store, html string) *Server {
func newServer(store store.Store, sessionStore sessions.Store, redis store_redis.Redis) *Server {
	// "/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api/logfile.log"
/*	f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		//log.Fatal(err)
		fmt.Printf("error opening file: %v", err)
	}
	//	defer f.Close()
	fmt.Println(f.Name())

	infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(f, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)*/

	s := &Server{
		router:       mux.NewRouter(), // mux.NewRouter()  NewRouter()
	//	logger:       logrus.New(),
		logger: logging.GetLogger(),
		store:        store,
		sessionStore: sessionStore,
	//	errorLog:     errorLog,
	//	infoLog:      infoLog,
		mu:           sync.Mutex{},
		redis:        redis,
		/*	httpServer: &http.Server{
			WriteTimeout:   15 * time.Second,
			ReadTimeout:    15 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20,
			ErrorLog:       errorLog,
		},*/
	}

	s.configureRouter()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)

}

func (s *Server) configureRouter() {

	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	//	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"http://localhost:3001"}), handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})))
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/test", s.diaplayPage())

	s.router.HandleFunc("/users", s.pagehandleUsersCreate()).Methods("GET")
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")

	s.router.HandleFunc("/updatepass", s.pageupdateSessionsCreate()).Methods("GET")
	s.router.HandleFunc("/updatepass", s.updateSessionsCreate()).Methods("POST")

	s.router.HandleFunc("/login", s.pagehandleSessionsCreate()).Methods("GET")
	s.router.HandleFunc("/login", s.handleSessionsCreate()).Methods("POST")
	s.router.HandleFunc("/logout", s.signOut()).Methods("GET")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
	//private.HandleFunc("/ho", s.ho().Methods("GET"))
	// /operation/***
	operation := s.router.PathPrefix("/operation").Subrouter()
	operation.Use(s.authenticateUser)
	//operation.Use(s.AuthMiddleware)
	operation.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")

	// api_server_management_quality_personal.go
	operation.HandleFunc("/showusersquality", s.PageshowUsersQuality()).Methods("GET")
	//s.router.HandleFunc("/showusersquality", s.showUsersQuality()).Methods("POST")
	operation.HandleFunc("/createusersquality", s.CreateUserQuality()).Methods("POST")
	operation.HandleFunc("/updateuserquality/{ID:[0-9]+}", s.PageupdateUserQuality()).Methods("GET")
	operation.HandleFunc("/updateuserquality/{ID:[0-9]+}", s.UpdateUserQuality()).Methods("POST")
	operation.HandleFunc("/deleteuserquality/{ID:[0-9]+}", s.DeleteUserQuality())
	//////////////////////////////////////////////////////////////////////////////////////////
	// api_server_management_warehouse_personal.go
	operation.HandleFunc("/showuserswarehouse", s.PageshowUsersWarehouse()).Methods("GET")
	operation.HandleFunc("/createuserswarehouse", s.CreateUserWarehouse()).Methods("POST")
	operation.HandleFunc("/updateuserwarehouse/{ID:[0-9]+}", s.PageupdateUserWarehouse()).Methods("GET")
	operation.HandleFunc("/updateuserwarehouse/{ID:[0-9]+}", s.UpdateUserWarehouse()).Methods("POST")
	operation.HandleFunc("/deleteuserwarehouse/{ID:[0-9]+}", s.DeleteUserWarehouse())
	////////////////////////////////////////////////////////////////////////////////////////////

	operation.HandleFunc("/shipmentbysap", s.pageshipmentBySAP()).Methods("GET")
	operation.HandleFunc("/shipmentbysap", s.shipmentBySAP()).Methods("POST")
	operation.HandleFunc("/showdateshipmentbysapbysearchstatic", s.pageshowShipmentBySAPBySearchStatic()).Methods("GET")
	operation.HandleFunc("/showdateshipmentbysapbysearchstatic", s.showShipmentBySAPBySearchStatic()).Methods("POST")

	operation.HandleFunc("/insertIDReturn", s.pageidReturn()).Methods("GET")
	operation.HandleFunc("/insertIDReturn", s.idReturn()).Methods("POST")

	operation.HandleFunc("/insertvendor", s.PageinsertVendor()).Methods("GET")
	operation.HandleFunc("/insertvendor", s.InsertVendor()).Methods("POST")
	// api_server_vendor_management.go
	operation.HandleFunc("/showvendor", s.PageVendor()).Methods("GET", "OPTIONS")
	operation.HandleFunc("/updatevendor/{ID:[0-9]+}", s.PageupdateVendor()).Methods("GET")
	operation.HandleFunc("/updatevendor/{ID:[0-9]+}", s.UpdateVendor()).Methods("POST")
	operation.HandleFunc("/deletevendor/{ID:[0-9]+}", s.DeleteVendor()) //.Methods("DELETE")
	/////////////////////////////////////////////////////////////////////////////////////////////
	// api_server_inspection_management.go
	operation.HandleFunc("/ininspection", s.PageinInspection()).Methods("GET")
	operation.HandleFunc("/ininspection", s.InInspection()).Methods("POST")
	operation.HandleFunc("/uploadfile", s.UploadFileToInspection()).Methods("POST")
	//operation.HandleFunc("/uploadfile", s.PageUploadFileToInspectionJSON()).Methods("GET")
	//operation.HandleFunc("/uploadfile", s.UploadFileToInspectionJSON()).Methods("POST", "OPTIONS")
	operation.HandleFunc("/historyinspection", s.PagehistoryInspection()).Methods("GET")
	operation.HandleFunc("/historyinspection", s.HistoryInspection()).Methods("POST")
	operation.HandleFunc("/historyinspectionp5", s.PagehistoryInspectionP5()).Methods("GET")
	operation.HandleFunc("/historyinspectionp5", s.HistoryInspectionP5()).Methods("POST")

	operation.HandleFunc("/statusinspection", s.PageInspection()).Methods("GET")
	operation.HandleFunc("/statusinspectionmix", s.PageInspectionMix()).Methods("GET")
	//operation.HandleFunc("/updateinspection/{ID:[0-9]+}", s.PageupdateInspection()).Methods("GET")
	//operation.HandleFunc("/updateinspection/{ID:[0-9]+}", s.UpdateInspection()).Methods("POST", "OPTIONS")
	operation.HandleFunc("/updateinspection/{ID:[0-9]+}", s.PageupdateInspectionJSON()).Methods("GET")
	operation.HandleFunc("/updateinspection", s.UpdateInspectionJSON()).Methods("POST", "OPTIONS")
	operation.HandleFunc("/updateinspectionmix/{ID:[0-9]+}", s.PageupdateInspectionJSONmix()).Methods("GET")
	operation.HandleFunc("/updateinspectionmix", s.UpdateInspectionJSONmix()).Methods("POST", "OPTIONS")
	operation.HandleFunc("/deleteinspection/{ID:[0-9]+}", s.DeleteInspection())

	operation.HandleFunc("/statusinspectionforwh", s.PageListAcceptWHInspection()).Methods("GET")
	//operation.HandleFunc("/acceptinspectiontowh/{ID:[0-9]+}", s.PageacceptWarehouseInspection()).Methods("GET")
	//	operation.HandleFunc("/acceptinspectiontowh/{ID:[0-9]+}", s.AcceptWarehouseInspection()).Methods("POST")
	operation.HandleFunc("/acceptinspectiontowh/{ID:[0-9]+}", s.PageacceptWarehouseInspectionJSON()).Methods("GET")
	operation.HandleFunc("/acceptinspectiontowh", s.AcceptWarehouseInspectionJSON()).Methods("POST", "OPTIONS")
	operation.HandleFunc("/acceptgroupsinspectiontowh", s.AcceptGroupsWarehouseInspection()).Methods("POST")
	//////////////////////////////////////////////////////////////////////////////////////////////////

	// /operationP5/***
	//	operationp5 := s.router.PathPrefix("/operationP5").Subrouter()
	//	operationp5.Use(s.AuthMiddleware)

	// api_server_p5_management_quality_personal.go
	//operationp5.HandleFunc("/showusersquality", s.PageshowUsersQuality()).Methods("GET")
	//operationp5.HandleFunc("/createusersquality", s.CreateUserQuality()).Methods("POST")

	//s.router.HandleFunc("/main", s.AuthMiddleware(s.main())).Methods("GET")
	operation.HandleFunc("/main", s.main()).Methods("GET")
	//operationp5.HandleFunc("/main", s.main()).Methods("GET")

	s.router.HandleFunc("/", s.upload()).Methods("GET")
	s.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	fmt.Println("Webserver StarLine launch.")
	//s.infoLog.Printf("Webserver StarLine launch.")

}

type Page struct {
	TitleDOC, Title, Content, Navbar, User, Username                 string
	LoggedIn, Admin, Stockkeeper, SuperIngenerQuality, StockkeeperWH bool
	Inspector, Quality                                               bool
}

func (s *Server) diaplayPage() http.HandlerFunc {
	//	type Page struct {
	//		TitleDOC, Title, Content, Navbar string
	//	Check                            bool
	//	}
	//	t = template.Must(template.ParseFiles("templates/index.html", "templates/head.html", "templates/f
	return func(w http.ResponseWriter, r *http.Request) {

		p := &Page{
			TitleDOC: "Chapter",
			Title:    "An Example",
			Content:  "Have fun stormin' da castle.",
			Navbar:   "ContentNav",
			//	Check:    false,
		}
		tpl.ExecuteTemplate(w, "index1.html", p) // ← Обработка шаблона с передачей данных
	}

}

func (s *Server) ho(r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(*model.User)
	fmt.Println("user ho", user)
}

func (s *Server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser).(*model.User)
		fmt.Println("handleWhoami", user)
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}



func (s *Server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *Server) logRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
		// open a file
		f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("error opening file: %v", err)
		}

		// don't forget to close it
		//defer f.Close()

		// Log as JSON instead of the default ASCII formatter.
		s.logger.SetFormatter(&logrus.TextFormatter{}) //(&s.logger.JSONFormatter{})

		// Output to stderr instead of stdout, could also be a file.
		s.logger.SetOutput(f)

		// Only log the warning severity or above.
		s.logger.SetLevel(s.logger.Level)

		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})

		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		// переопределение ResponseWriter
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed witn %d %s in %v",
			rw.code,
			http.StatusText(rw.code), //2bit.ru/media/images/items/248/247158-413538.jpgtp.StatusText(rw.code),
			time.Now().Sub(start),
		)

		*/
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)
		start := time.Now()
		// переопределение ResponseWriter
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed witn %d %s in %v",
			rw.code,
			http.StatusText(rw.code), //2bit.ru/media/images/items/248/247158-413538.jpgtp.StatusText(rw.code),
			time.Now().Sub(start),
		)
		//	var log = logrus.New()
		//	log.Out = os.Stdout
		//	file, err := os.OpenFile("/tmp/logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		//	if err == nil {
		//		log.Out = file
		//	} else {
		//		log.Info("Failed to log to file, using default stderr")
		//	}

	})
}

func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}

func (s *Server) pagehandleUsersCreate() http.HandlerFunc {
	//	tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/register.html"))
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/register.html")
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "register.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/register.html")
		//	fmt.Fprintf(w, body)
		tpl.ExecuteTemplate(w, "register.html", nil)
		///	err = tpl.ExecuteTemplate(w, "layout", nil)
		///	if err != nil {
		///		http.Error(w, err.Error(), 400)
		///		return
		///	}
	}
}

func (s *Server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	}

	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "register.html")
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/register.html")
	///	if err != nil {
	///	panic(err)
	///}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		req := &request{}
		//	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		//		s.error(w, r, http.StatusBadRequest, err)
		//		return
		//	}
		fmt.Println("зашел на регистрацию")
		s.logger.Info("Starting registration")
		//	email := r.FormValue("email")
		//	password := r.FormValue("password")
		//	firstname := r.FormValue("firstname")
		//	lastname := r.FormValue("lastname")

		req.Email = r.FormValue("email")
		//fmt.Println(req.Email)

		req.Password = r.FormValue("password")

		req.FirstName = r.FormValue("firstname")

		req.LastName = r.FormValue("lastname")

		//s.infoLog.Printf("Create account: %v, %v, %v, %v\n", req.Email, req.Password, req.FirstName, req.LastName)
		s.logger.Infof("Create account: %v, %v, %v, %v\n", req.Email, req.Password, req.FirstName, req.LastName)
		//	target := "/users"

		//	fmt.Println(req.Password)
		//	fmt.Println(req.FirstName)
		//	fmt.Println(req.LastName)
		u := &model.User{
			Email:     req.Email,     //req.Email, email
			Password:  req.Password,  //req.Password, password
			FirstName: req.FirstName, // firstname req.FirstName
			LastName:  req.LastName,  // lastname req.LastName
		}
		//json.NewEncoder(w).Encode(u)

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()

		//	s.respond(w, r, http.StatusCreated, u)
		//	target = "/one"
		//	http.Redirect(w, r, target, 302)

		//	}
		fmt.Println("регистрируюсь")
		s.logger.Info("End registration")
		//tpl.ExecuteTemplate(w, "login.html", nil)
		tpl.ExecuteTemplate(w, "register.html", nil)
		///	err = tpl.ExecuteTemplate(w, "layout", nil)
		///	if err != nil {
		///		http.Error(w, err.Error(), 400)
		///		return
		///	}
	}
}

func (s *Server) pageupdateSessionsCreate() http.HandlerFunc {
	//	tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/login.html"))
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/login.html")
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateUser.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/login.html") // "./web/templates/login.html"
		//	fmt.Fprintf(w, body)
		tpl.ExecuteTemplate(w, "updateUser.html", nil)
		///	err = tpl.ExecuteTemplate(w, "layout", nil)
		///	if err != nil {
		///		http.Error(w, err.Error(), 400)
		///		return
		///	}
	}
}

func (s *Server) updateSessionsCreate() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
		//	Tabel    string `json:"tabel"`
		PasswordOld string `json:"passwordold"`
		Password    string `json:"password"`
	}
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateUser.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		req := &request{}
		/*	f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				fmt.Println(err)
				//	return
			}
			defer f.Close()
			iLog := log.New(f, "func updateSessionsCreate ", log.LstdFlags)
			iLog.SetFlags(log.LstdFlags | log.Lshortfile)*/

		req.Email = r.FormValue("email")
		req.PasswordOld = r.FormValue("passwordold")
		req.Password = r.FormValue("password")

		//s.infoLog.Printf("Update password: email - %v, passwordold - %v, password - %v\n",
		//	req.Email, req.PasswordOld, req.Password)
		s.logger.Infof("Update password: email - %v, passwordold - %v, password - %v\n",
			req.Email, req.PasswordOld, req.Password)

		up, err := s.store.User().FindByEmail(req.Email, req.Email)
		if err != nil || !up.ComparePassword(req.PasswordOld) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)

			return
		}

		u := &model.User{
			Email: req.Email, //req.Email, email
			//	Tabel:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().UpdatePass(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u.Sanitize()

		fmt.Println("обновляю пароль")
		s.logger.Info("Update password")
		http.Redirect(w, r, "/", 303)
		err = tpl.ExecuteTemplate(w, "updateUser.html", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *Server) upload() http.HandlerFunc {
	//	type Page struct {
	//		TitleDOC, Navbar string
	//		LoggedIn         bool
	//	}
	return func(w http.ResponseWriter, r *http.Request) {
		//tpl.ExecuteTemplate(w, "base", nil)
		//fmt.Print("Test Upload")
		p := &Page{
			TitleDOC: "Главная",
			//	Navbar:   "ContentNav",
			LoggedIn: false,
		}
		tpl.ExecuteTemplate(w, "index0.html", p)
		//tpl.Execute(w, nil)
	}
}
/*
func (s *Server) RecognationUser(h http.Handler)  *model.User {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		return *model.User, nil
	}
}
*/
func (s *Server) main() http.HandlerFunc {

	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles(s.html + "layout.html"))
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles(s.html+"layout1.html", s.html+"login1.html"))
	//tpl = template.Must(template.New("base").ParseFiles(s.html + "layout1.html"))
	//tpl = template.Must(template.ParseFiles(s.html+"header.html", s.html+"footer.html"))
	//tpl = template.ParseFiles(s.html + "*.html")
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/layout.html"))
	//tpl = template.Must(template.ParseFiles("web/templates/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		u:= r.Context().Value(ctxKeyUser).(*model.User)

		fmt.Println("/main - user:", u.Email, u.ID, u.Role)
		//s.infoLog.Printf("main, user - %s, %d, %s", u.Email, u.ID, u.Role)
		s.logger.Infof("main, user - %s, %d, %s", u.Email, u.ID, u.Role)
		if u.Groups == groupWarehouse || u.Groups == groupQuality || u.Groups == groupEmpty || u.Groups == groupAdministrator {

			if u.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true
				statusGroupP1 = true

			} else if u.Role == roleStockkeeper {
				statusStockkeeper = true
				statusLoggedIn = true
				statusGroupP1 = true

			} else if u.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
				statusGroupP1 = true

			} else if u.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
				statusGroupP1 = true

			} else if u.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true
				statusGroupP1 = true

			} else if u.Role == roleInspector {
				statusInspector = true
				statusLoggedIn = true
				statusGroupP1 = true

			} else if u.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
				statusGroupP1 = true
				fmt.Println("yes yes %v, %v", u.Groups, u.Role)
			}

			//	data := &Page{
			//		LoggedIn:            LoggedIn,
			//		Admin:               Admin,
			//		Stockkeeper:         Stockkeeper,
			//		SuperIngenerQuality: SuperIngenerQuality,
			//		StockkeeperWH:       StockkeeperWH,
			//		Inspector:           Inspector,
			//		User:                u.LastName,
			//		ID:                  u.FirstName,
			//	}
			//	GET := map[string]bool{
			//		"Admin":               Admin,
			//		"Stockkeeper":         Stockkeeper,
			//		"главный инженер по качеству": SuperIngenerQuality,
			//		"StockkeeperWH":       StockkeeperWH,
			//		"Inspector":           Inspector,
			//	}

			data := map[string]interface{}{
				"TitleDOC": "Главная",
				"User":     u.LastName,
				"Username": u.FirstName,
				//"GET":                 GET,
				"Admin":               statusAdmin,
				"Stockkeeper":         statusStockkeeper,
				"SuperIngenerQuality": statusSuperIngenerQuality,
				"WarehouseManager":    statusWarehouseManager,
				"IngenerQuality":      statusIngenerQuality,
				"StockkeeperWH":       statusStockkeeperWH,
				"Inspector":           statusInspector,
				"GroupP1":             statusGroupP1,
				"LoggedIn":            statusLoggedIn,
			}

			//tpl.ExecuteTemplate(w, "home.html", data)
			//tpl.ExecuteTemplate(w, "layout", data)
			//tpl.ExecuteTemplate(w, "base", data)
			tpl.ExecuteTemplate(w, "index.html", data) // index3.html
		}

		if u.Groups == groupWarehouseP5 || u.Groups == groupQualityP5 {
			statusGroupP5 = true
			if u.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true

			} else if u.Role == roleStockkeeper {
				statusStockkeeper = true
				statusLoggedIn = true

			} else if u.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true

			} else if u.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true

			} else if u.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true

			} else if u.Role == roleInspector {
				statusInspector = true
				statusLoggedIn = true

			} else if u.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
			}

			//	data := &Page{
			//		LoggedIn:            LoggedIn,
			//		Admin:               Admin,
			//		Stockkeeper:         Stockkeeper,
			//		SuperIngenerQuality: SuperIngenerQuality,
			//		StockkeeperWH:       StockkeeperWH,
			//		Inspector:           Inspector,
			//		User:                u.LastName,
			//		ID:                  u.FirstName,
			//	}
			//	GET := map[string]bool{
			//		"Admin":               Admin,
			//		"Stockkeeper":         Stockkeeper,
			//		"главный инженер по качеству": SuperIngenerQuality,
			//		"StockkeeperWH":       StockkeeperWH,
			//		"Inspector":           Inspector,
			//	}

			data := map[string]interface{}{
				"TitleDOC": "Главная",
				"User":     u.LastName,
				"Username": u.FirstName,
				//"GET":                 GET,
				"Admin":               statusAdmin,
				"Stockkeeper":         statusStockkeeper,
				"SuperIngenerQuality": statusSuperIngenerQuality,
				"WarehouseManager":    statusWarehouseManager,
				"IngenerQuality":      statusIngenerQuality,
				"StockkeeperWH":       statusStockkeeperWH,
				"Inspector":           statusInspector,
				"GroupP5":             statusGroupP5,
				"LoggedIn":            statusLoggedIn,
			}

			//tpl.ExecuteTemplate(w, "home.html", data)
			//tpl.ExecuteTemplate(w, "layout", data)
			//tpl.ExecuteTemplate(w, "base", data)
			tpl.ExecuteTemplate(w, "index.html", data) // index3.html
		}
	}
}

func (s *Server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		//	s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		fmt.Println("authenticateUser: ", u.Email, u.Role, u.Tabel)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *Server) AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id := session.Values["user_id"]
		fmt.Println("id is:", id)

		//	if !ok {
		//		//	s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
		//		http.Redirect(w, r, "/sessions", http.StatusSeeOther)
		//		return
		//	}

		//	email := session.Values["user_id"]
		//	u, err := s.store.User().Find(id.(int))
		//	fmt.Println("u is:", u)
		//	if err != nil {
		//		s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
		//		//http.Redirect(w, r, "/sessions", http.StatusSeeOther)
		//		return
		//	}

		if id == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			//w.Write([]byte("Ooops"))
			//tpl.ExecuteTemplate(w, "login.html", nil)
		} else {
			h.ServeHTTP(w, r)
		}

	})
}

func (s *Server) pagehandleSessionsCreate() http.HandlerFunc {
	//	type Page struct {
	//		TitleDOC string
	//LoggedIn         bool
	//	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/login.html") // "./web/templates/login.html"
		//	fmt.Fprintf(w, body)
		//	err = tpl.ExecuteTemplate(w, "layout", nil)
		//	if err != nil {
		//		http.Error(w, err.Error(), 400)
		//		return
		//	}
		//tpl.ExecuteTemplate(w, "base", nil)
		p := &Page{
			TitleDOC: "Login",
			//	Navbar:   "ContentNav",
		}
		fmt.Println("Test")
		tpl.ExecuteTemplate(w, "login.html", p)
	}
}

func (s *Server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/layout.html")
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "layout.html")
	//if err != nil {
	//	panic(err)
	//}
	//tpl, err := template.New("base").ParseFiles(s.html + "layout1.html")
	//if err != nil {
	//	panic(err)
	//}
	//tpl = template.Must(template.ParseFiles(s.html + "/*.html"))
	//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//tpl = template.Must(template.ParseFiles(s.html+"header.html", s.html+"footer.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		/*	f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				fmt.Println(err)
				//	return
			}
			defer f.Close()
			iLog := log.New(f, "func handleSessionsCreate ", log.LstdFlags)
			iLog.SetFlags(log.LstdFlags | log.Lshortfile)*/
		//	reqBody, err := ioutil.ReadAll(r.Body)
		//	if err != nil {
		//		log.Fatal(err)
		//	}
		//	fmt.Printf("%s", reqBody)

		//	req := &request{}
		//	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		//		s.error(w, r, http.StatusBadRequest, err)
		//		return
		//	}

		//	fmt.Printf("%s", req)
		//		if r.Method == http.MethodPost

		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Println("email:   ", email)
		//s.infoLog.Printf("Loggin: %v, %v\n", email, password)
		s.logger.Infof("Loggin: %v, %v\n", email, password)
		//	target := "/sessions"

		//	s.Lock()
		u, err := s.store.User().FindByEmail(email, email)
		//	fmt.Println("FindByEmail:   ", u)
		//	match := u.ComparePassword(password)
		//	fmt.Println("Match:   ", match)
		if err != nil || !u.ComparePassword(password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)

			return
		}
		//	s.Unlock()
		if u.Groups == groupWarehouse || u.Groups == groupQuality || u.Groups == groupEmpty || u.Groups == groupAdministrator {
			statusGroupP1 = true
			if u.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true

			} else if u.Role == roleStockkeeper {
				statusStockkeeper = true
				statusLoggedIn = true

			} else if u.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
			} else if u.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
			} else if u.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true

			} else if u.Role == roleInspector {
				statusInspector = true
				statusLoggedIn = true

			} else if u.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
			}

			//	GET := map[string]bool{
			//		"admin":               admin,
			//		"stockkeeper":         stockkeeper,
			//		"главный инженер по качеству": superIngenerQuality,
			//		"stockkeeperWH":       stockkeeperWH,
			//		"inspector":           inspector,
			//	}

			session, err := s.sessionStore.Get(r, sessionName)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			session.Values["user_id"] = u.ID
			if err := s.sessionStore.Save(r, w, session); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, nil)

			fmt.Println("handleSessionsCreate()", u.Email, u.Role)
			//	iLog.Println("залогинился пользователь", u.Email)
			//	s.redirectMain()
			//	s.pageredirectMain()
			//	_, ok := session.Values["user_id"]
			//	if !ok {

			//	return
			//	}
			//	data := &Page{
			//		Admin:               Admin,
			//		Stockkeeper:         Stockkeeper,
			//		SuperIngenerQuality: SuperIngenerQuality,
			//		StockkeeperWH:       StockkeeperWH,
			//		Inspector:           Inspector,
			//		User:                u.LastName,
			//		ID:                  u.FirstName,
			//		LoggedIn:            LoggedIn,
			//	}
			data := map[string]interface{}{
				"User":     u.LastName,
				"Username": u.FirstName,
				//		"GET":  GET,
				"Admin":               statusAdmin,
				"Stockkeeper":         statusStockkeeper,
				"SuperIngenerQuality": statusSuperIngenerQuality,
				"WarehouseManager":    statusWarehouseManager,
				"IngenerQuality":      statusIngenerQuality,
				"StockkeeperWH":       statusStockkeeperWH,
				"Inspector":           statusInspector,
				"GroupP1":             statusGroupP1,
				"LoggedIn":            statusLoggedIn,
			}
			//	tpl.ExecuteTemplate(w, "base", data) //  "index.html"
			tpl.ExecuteTemplate(w, "index.html", data)
			//err = tpl.ExecuteTemplate(w, "base", data)
			//if err != nil {
			//	http.Error(w, err.Error(), 400)
			//	return
			//}
			http.Redirect(w, r, "/", http.StatusFound)
		}

		if u.Groups == groupWarehouseP5 || u.Groups == groupQualityP5 {
			statusGroupP5 = true
			if u.Role == roleAdministrator {
				statusAdmin = true
				statusLoggedIn = true

			} else if u.Role == roleStockkeeper {
				statusStockkeeper = true
				statusLoggedIn = true

			} else if u.Role == roleSuperIngenerQuality {
				statusSuperIngenerQuality = true
				statusLoggedIn = true
			} else if u.Role == roleIngenerQuality {
				statusIngenerQuality = true
				statusLoggedIn = true
			} else if u.Role == roleStockkeeperWH {
				statusStockkeeperWH = true
				statusLoggedIn = true

			} else if u.Role == roleInspector {
				statusInspector = true
				statusLoggedIn = true

			} else if u.Role == roleWarehouseManager {
				statusWarehouseManager = true
				statusLoggedIn = true
			}

			//	GET := map[string]bool{
			//		"admin":               admin,
			//		"stockkeeper":         stockkeeper,
			//		"главный инженер по качеству": superIngenerQuality,
			//		"stockkeeperWH":       stockkeeperWH,
			//		"inspector":           inspector,
			//	}

			session, err := s.sessionStore.Get(r, sessionName)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			session.Values["user_id"] = u.ID
			if err := s.sessionStore.Save(r, w, session); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, nil)

			fmt.Println("handleSessionsCreate()", u.Email, u.Role)
			//	iLog.Println("залогинился пользователь", u.Email)
			//	s.redirectMain()
			//	s.pageredirectMain()
			//	_, ok := session.Values["user_id"]
			//	if !ok {

			//	return
			//	}
			//	data := &Page{
			//		Admin:               Admin,
			//		Stockkeeper:         Stockkeeper,
			//		SuperIngenerQuality: SuperIngenerQuality,
			//		StockkeeperWH:       StockkeeperWH,
			//		Inspector:           Inspector,
			//		User:                u.LastName,
			//		ID:                  u.FirstName,
			//		LoggedIn:            LoggedIn,
			//	}
			data := map[string]interface{}{
				"User":     u.LastName,
				"Username": u.FirstName,
				//		"GET":  GET,
				"Admin":               statusAdmin,
				"Stockkeeper":         statusStockkeeper,
				"SuperIngenerQuality": statusSuperIngenerQuality,
				"WarehouseManager":    statusWarehouseManager,
				"IngenerQuality":      statusIngenerQuality,
				"StockkeeperWH":       statusStockkeeperWH,
				"Inspector":           statusInspector,
				"GroupP5":             statusGroupP5,
				"LoggedIn":            statusLoggedIn,
			}
			//	tpl.ExecuteTemplate(w, "base", data) //  "index.html"
			tpl.ExecuteTemplate(w, "index.html", data)
			//err = tpl.ExecuteTemplate(w, "base", data)
			//if err != nil {
			//	http.Error(w, err.Error(), 400)
			//	return
			//}
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func (s *Server) signOut() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		session, err := s.sessionStore.Get(r, sessionName)
		//	if err != nil {
		//		s.error(w, r, http.StatusInternalServerError, err)
		//		return
		//	}
		//	session.Values["user_id"] = false

		session.Options.MaxAge = -1
		/*
			session, _ := s.sessionStore.Get(r, sessionName)
			//	session.Options.MaxAge = -1
			u, err := s.store.User().Find(id)
			session.Values["user_id"] = u.ID
			session.Values["user_id"] = false */
		session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Очистка куки")
		/*
			c, err := r.Cookie(sessionName)
			if err != nil {
				panic(err.Error())
			}
			c.Name = "Deleted"
			c.Value = "Unuse"
			c.Expires = time.Unix(1414414788, 1414414788000)
		*/
		http.Redirect(w, r, "/", http.StatusSeeOther)
		//	tpl.ExecuteTemplate(w, "login.html", nil)
	})
}

func (s *Server) pageshipmentBySAP() http.HandlerFunc {
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertsapbyship6.html"))
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertsapbyship6.html")
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertsapbyship6.html")
	///if err != nil {
	///	panic(err)
	///}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/insertsapbyship6.html")
		//	fmt.Fprintf(w, body)
		//	data := map[string]interface{}{
		//		"user": "Я тут",
		//	}

		u:= r.Context().Value(ctxKeyUser).(*model.User)
		if u.Role == roleAdministrator {
			statusAdmin = true
			statusLoggedIn = true

		} else if u.Role == roleStockkeeper {
			statusStockkeeper = true
			statusLoggedIn = true

		}
		data := map[string]interface{}{
			"TitleDOC":    "Отгрузка изделий",
			"User":        u.LastName,
			"Username":    u.FirstName,
			"Admin":       statusAdmin,
			"Stockkeeper": statusStockkeeper,
			"LoggedIn":    statusLoggedIn,
		}
		tpl.ExecuteTemplate(w, "insertsapbyship6.html", data)
		///	err = tpl.ExecuteTemplate(w, "layout", data)
		///	if err != nil {
		///		http.Error(w, err.Error(), 400)
		///		return
		//	}
		//		if err = tpl.ExecuteTemplate(w, "layout", nil); err != nil {
		//			s.error(w, r, http.StatusUnprocessableEntity, err)
		//			return
		//		}
		///	}
	}
}

func (s *Server) shipmentBySAP() http.HandlerFunc {
	type reqA struct {
		Material string //int
		Qty      int
		Comment  string
		//	ID       int
		//	LastName string
	}

	type request struct {
		Material     int       `db:"material"`
		Qty          int       `db:"qty"`
		Comment      string    `db:"comment"`
		ID           int       `db:"id"`
		LastName     string    `db:"lastname"`
		ShipmentDate time.Time `db:"shipment_date"`
		ShipmentTime time.Time `db:"shipment_time"`
	}

	//	tpl = template.Must(template.ParseFiles("web/templates/insertsapbyship6.html"))
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertsapbyship6.html"))
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertsapbyship6.html")
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertsapbyship6.html")
	///if err != nil {
	///	panic(err)
	///}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		//	hdata := reqA{}
		//	var hdata ReqA

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			s.logger.Infof(err.Error())
			s.logger.Errorf(err.Error())
		}

		var hdata []reqA

		json.Unmarshal(body, &hdata)
		//	json.Marshal(body)
		fmt.Printf("body json: %s", body)
		//s.infoLog.Printf("Loading body json shipmentBySAP: %s\n", body)
		s.logger.Infof("Loading body json shipmentBySAP: %s\n", body)
		//	fmt.Printf("1 тест %s", &hdata.Material)
		//	fmt.Println("tect2 %s", hdata, "\n")
		fmt.Println("\njson  struct hdata", hdata)
		//s.infoLog.Printf("Loading hdata json shipmentBySAP: %v\n", hdata)
		s.logger.Infof("Loading hdata json shipmentBySAP: %v\n", hdata)
		//req := &request{}
		//	if err := json.NewDecoder(r.Body).Decode(&hdata); err != nil {
		//		s.error(w, r, http.StatusBadRequest, err)
		//		return
		//	}
		//	fmt.Printf("%s", hdata)

		//	material := r.FormValue("material")
		//	material, err := strconv.Atoi(r.FormValue("Material[]"))
		//	material, err := strconv.Atoi(r.FormValue("material"))
		//	material := req.Material
		//	использую	material, err := strconv.Atoi(r.FormValue("material"))
		//	material := r.FormValue("material")
		//	material, err := strconv.ParseInt(r.FormValue("material")[0:], 10, 64)
		//		if err != nil {
		//			fmt.Println(err)
		//		}
		// проверка кол-ва символом в номере материала SAP
		//	checkmaterial := strconv.FormatInt(material, 10)
		//	checkmaterial := strconv.Itoa(material)
		///checkmaterial := material
		//	if len(checkmaterial) != 7 {
		//		tpl.ExecuteTemplate(w, "error.html", nil)
		//		return
		//	}
		//	req.Qty = r.ParseForm("qty")
		// qty, err := strconv.ParseInt(r.FormValue("qty")[0:], 10, 64)
		//	qty := req.Qty
		//	использую	qty, err := strconv.Atoi(r.FormValue("qty"))
		//	qty := r.FormValue("qty")
		//		if err != nil {
		//			fmt.Println(err)
		//		}
		//checkqty := strconv.Itoa(qty)
		//if len(checkqty) > 1 {
		//	tpl.ExecuteTemplate(w, "error.html", nil)
		//		return
		//	}
		//req.Comment = r.ParseForm("comment")
		//	comment := req.Comment
		//	использую	comment := r.FormValue("comment")

		user := r.Context().Value(ctxKeyUser).(*model.User)
		//	hdata.ID = user.ID
		//	hdata.LastName = user.LastName

		//	var record []request
		//	record.Material = hdata.Material
		//	//	fmt.Printf("2 тест %v", record.Material)
		//	record.Material = hdata.Qty
		//	record.Comment = hdata.Comment
		//	record.ID = user.ID
		//	record.LastName = user.LastName
		//	user2 := strconv.Atoi(user1)

		for _, v := range hdata {
			fmt.Println(v.Material, v.Qty, v.Comment, user.ID, user.LastName)
			//s.infoLog.Printf("shipmentBySAP: %v, %v, %v, %v, %v\n", v.Material, v.Qty, v.Comment, user.ID, user.LastName)
			s.logger.Infof("shipmentBySAP: %v, %v, %v, %v, %v\n", v.Material, v.Qty, v.Comment, user.ID, user.LastName)
			checkmaterial, _ := strconv.Atoi(v.Material)

			//	checkmaterial := v.Material
			if len(v.Material) == 7 {
				u := &model.Shipmentbysap{
					Material: checkmaterial, //v.Material,    //record.Material, //material,
					Qty:      v.Qty,         //record.Qty,      // qty,
					Comment:  v.Comment,     //record.Comment,  //comment,
					ID:       user.ID,       //record.ID,       //user.ID,
					LastName: user.LastName, //record.LastName, // user.LastName,

				}

				//	var send bool
				if err := s.store.Shipmentbysap().InterDate(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}

				//	}

				//	fm := session.Flashes("message")
				//	session.Save(r, w)
				//	fmt.Fprintf(w, "%v", fm[0])
			} else {
				//if checkmaterial != 7 {
				if len(v.Material) != 7 {
					fmt.Println("кол-во не равно 7", v.Material)
					s.logger.Errorf("кол-во не равно 7 %v", v.Material)
					s.logger.Infof("кол-во не равно 7 %v", v.Material)
					strN, err := json.Marshal("JSON кол-во не равно 7.")
					fmt.Println(string(strN))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					w.Write(strN)
					//	break
				}

				/*	session, err := s.sessionStore.Get(r, "flash-session")
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					session.AddFlash("This is a flashed message!", "message")
					fm := session.Flashes("message")
					session.Save(r, w)
					fmt.Fprintf(w, "%v", fm[0]) */
				//	tpl.ExecuteTemplate(w, "error.html", nil)
				//	return
			}

			//	http.Redirect(w, r, "/main", 303)
			//	http.Redirect(w, r, "/main", 303)
			//	return
			//	}
			//	u :=
			//	fmt.Println("Material print:", u.Material)
			//	fmt.Println("Qty print:", u.Qty)
			//	fmt.Println("Comment print:", u.Comment)

			//	} else {
			//		tpl.ExecuteTemplate(w, "error.htmp", nil)
			//		return
			//	}

		}
		strB, err := json.Marshal("Данные отправлены.")
		fmt.Println(string(strB))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(strB)
		err = tpl.ExecuteTemplate(w, "insertsapbyship6.html", nil)
		//	err = tpl.ExecuteTemplate(w, "layout", v)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		//	tpl.ExecuteTemplate(w, "insertsapbyship6.html", nil)
		//	sessionF, err := s.sessionStore.Get(r, "flash-session")
		//	if err != nil {
		//		http.Error(w, err.Error(), http.StatusInternalServerError)
		//	}
		//	sessionF.AddFlash("This is a flashed message!", "FlashedMessages")
		//	v := map[string]interface{}{
		//		"FlashedMessages": session.Flashes(),
		//	}

	}

}

func (s *Server) pageshowShipmentBySAPBySearchStatic() http.HandlerFunc {
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysapbysearchstatic.html")
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdatebysapbysearchstatic.html")
	///if err != nil {
	///	panic(err)
	///}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/showdatebysapbysearchstatic.html")
		//	fmt.Fprintf(w, body)

		u:= r.Context().Value(ctxKeyUser).(*model.User)
		if u.Role == roleAdministrator {
			statusAdmin = true
			statusLoggedIn = true

		} else if u.Role == roleStockkeeper {
			statusStockkeeper = true
			statusLoggedIn = true

		}
		data := map[string]interface{}{
			"TitleDOC":    "Отгрузка изделий",
			"User":        u.LastName,
			"Username":    u.FirstName,
			"Admin":       statusAdmin,
			"Stockkeeper": statusStockkeeper,
			"LoggedIn":    statusLoggedIn,
		}
		tpl.ExecuteTemplate(w, "showdatebysapbysearchstatic.html", data)
		///	if err = tpl.ExecuteTemplate(w, "layout", data); err != nil {
		///		s.error(w, r, http.StatusUnprocessableEntity, err)
		///		return
		///	}

	}
}

func (s *Server) showShipmentBySAPBySearchStatic() http.HandlerFunc {
	type searchBy struct {
		LastName string `json:"lastname"`
		Date1    string `json:"date1"`
		Date2    string `json:"date2"`
		Material int    `json:"material"`
	}

	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysap2.html")
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdatebysap2.html")
	///if err != nil {
	///	panic(err)
	///}

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		u:= r.Context().Value(ctxKeyUser).(*model.User)
		if u.Role == roleAdministrator {
			statusAdmin = true
			statusLoggedIn = true

		} else if u.Role == roleStockkeeper {
			statusStockkeeper = true
			statusLoggedIn = true

		}

		search := &searchBy{}
		//		r.ParseForm()
		materialInt, err := strconv.Atoi(r.FormValue("material"))
		if err != nil {
			log.Println(err)
			s.logger.Infof(err.Error())
			s.logger.Errorf(err.Error())
		}

		search.LastName = r.FormValue("lastname")
		fmt.Println("lastname - ", search.LastName)
		search.Date1 = r.FormValue("date1")
		fmt.Println("date1 - ", search.Date1)
		search.Date2 = r.FormValue("date2")
		fmt.Println("date2 - ", search.Date2)
		search.Material = materialInt
		fmt.Println("material - ", search.Material)

		//	search.Material = r.FormValue("material")
		/*
			ss := &model.Shipmentbysap{
				LastName:      search.LastName,
				ShipmentDate2: search.Date1,
				ShipmentDate3: search.Date2,
				Material:      search.Material,
			}
		*/

		get, err := s.store.Shipmentbysap().ShowDateBySearch(search.LastName, search.Date1, search.Date2, search.Material)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		data := map[string]interface{}{
			"TitleDOC":    "Отгрузка изделий",
			"User":        u.LastName,
			"Username":    u.FirstName,
			"Admin":       statusAdmin,
			"Stockkeeper": statusStockkeeper,
			"LoggedIn":    statusLoggedIn,
			"GET":         get,
		}

		//	err = tpl.ExecuteTemplate(w, "showdatebysap2.html", get)
		err = tpl.ExecuteTemplate(w, "showdatebysap2.html", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *Server) pageidReturn() http.HandlerFunc {
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertidreturn2.html")
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertidreturn2.html")
	///if err != nil {
	///	panic(err)
	///}
	//	tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertidreturn2.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/insertidreturn2.html")
		//	fmt.Fprintf(w, body)
		u:= r.Context().Value(ctxKeyUser).(*model.User)
		if u.Role == roleAdministrator {
			statusAdmin = true
			statusLoggedIn = true

		} else if u.Role == roleStockkeeper {
			statusStockkeeper = true
			statusLoggedIn = true

		}
		data := map[string]interface{}{
			"TitleDOC":    "Проверка катушек",
			"User":        u.LastName,
			"Username":    u.FirstName,
			"Admin":       statusAdmin,
			"Stockkeeper": statusStockkeeper,
			"LoggedIn":    statusLoggedIn,
		}
		tpl.ExecuteTemplate(w, "insertidreturn2.html", data)
		///	err = tpl.ExecuteTemplate(w, "layout", nil)
		///	if err != nil {
		///		http.Error(w, err.Error(), 400)
		///		return
		///	}
	}
}

func (s *Server) idReturn() http.HandlerFunc {
	type req struct {
		ScanID     string `json:"scanid"`
		Material   int
		IDRoll     int
		Lot        string
		QtyFact    int `json:"qtyfact"`
		QtySAP     int `json:"qtysap"`
		QtyPanacim int `json:"qtypanacim"`
	}

	//tpl = template.Must(template.ParseFiles("web/templates/insertidreturn2.html"))
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertidreturn2.html"))
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertidreturn2.html")
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertidreturn2.html")
	///if err != nil {
	///	panic(err)
	///}
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
		fmt.Printf("test %s", body)
		//s.infoLog.Printf("idReturn loadin body json %s\n", body)
		s.logger.Infof("idReturn loadin body json %s\n", body)
		fmt.Println("\nall of the rdata", rdata)
		//s.infoLog.Printf("idReturn loadin rdata json %v\n", rdata)
		s.logger.Infof("idReturn loadin rdata json %v\n", rdata)

		user := r.Context().Value(ctxKeyUser).(*model.User)

		for _, v := range rdata {
			sap := v.ScanID[1:8]
			material := v.Material
			material, err := strconv.Atoi(sap)
			if err != nil {
				log.Println(err)
				s.logger.Errorf(err.Error())
				s.logger.Infof(err.Error())
			}
			idsap := v.ScanID[20:30]
			idroll := v.IDRoll
			idroll, err = strconv.Atoi(idsap)
			fmt.Println("idroll в 1-м цикле -", idroll)
			if err != nil {
				log.Println(err)
				s.logger.Errorf(err.Error())
				s.logger.Infof(err.Error())
			}
			v.Lot = v.ScanID[9:19]
			if (strings.Contains(v.ScanID[0:1], "P") == true) && (len(v.ScanID) == 45) {

				u := &model.IDReturn{
					Material:   material,
					IDRoll:     idroll,
					Lot:        v.Lot,
					QtyFact:    v.QtyFact,
					QtySAP:     v.QtySAP,
					QtyPanacim: v.QtyPanacim,
					ID:         user.ID, //record.ID,       //user.ID,
					LastName:   user.LastName,
				}
				fmt.Println("idroll - ", idroll)

				if err := s.store.IDReturn().InterDate(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)

					return
				}

				//	fmt.Fprintf(w, "Date of ID uploaded successfully")
				//	return

			} else {
				if (strings.Contains(v.ScanID[0:1], "P") == false) && (len(v.ScanID) != 45) {
					fmt.Println("не верное сканирование :\n" + v.ScanID + "\n")
					s.logger.Errorf("не верное сканирование :\n" + v.ScanID + "\n")
					s.logger.Infof("не верное сканирование :\n" + v.ScanID + "\n")
					//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
				}
				//	tpl.Execute(w, data)
				return
			}

		}
		//	data := map[string]interface{}{
		//		"fuck": "OK",
		//	}
		tpl.ExecuteTemplate(w, "insertidreturn2.html", nil)
		///	err = tpl.ExecuteTemplate(w, "layout", nil)
		//	tpl.ExecuteTemplate(w, "layout", data)
		///	if err != nil {
		///		http.Error(w, err.Error(), 400)
		///		return
		///	}

	}
}

/*
func (s *Server) getStatic() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		getStatic, err := s.store.Inspection().CountDebitor()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		static := map[string]string{
			"Static": getStatic,
		}
		pie := charts.NewPie()
		pie.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{Title: "Radius style"}),
		)

		pie.AddSeries("pie", static).
			SetSeriesOptions(
				charts.WithLabelOpts(opts.Label{
					Show:      true,
					Formatter: "{b}: {c}",
				}),
				charts.WithPieChartOpts(opts.PieChart{
					Radius: []string{"40%", "75%"},
				}),
			)

		f, err := os.Create("./web/templates/test.html")
		if err != nil {
			panic(err)
		}
		pie.Render(io.Writer(f))
		//	Examples()
		get := map[string]interface{}{
			"GetStatic": getStatic,
		}
		tpl.ExecuteTemplate(w, "test.html", get)
	}

}
*/
/*
func pieRadius(s *Server) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Radius style"}),
	)
	getStatic, err := s.store.Inspection().CountDebitor()
	if err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	pie.AddSeries("pie", getStatic).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      true,
				Formatter: "{b}: {c}",
			}),
			charts.WithPieChartOpts(opts.PieChart{
				Radius: []string{"40%", "75%"},
			}),
		)
	return pie
}

type PieExamples struct{}

func (PieExamples) Examples() {
	page := components.NewPage()
	page.AddCharts(
		//	pieBase(),
		//	pieShowLabel(),
		pieRadius(),
	//	pieRoseArea(),
	//	pieRoseRadius(),
	//	pieRoseAreaRadius(),
	//	pieInPie(),
	)
	f, err := os.Open("./web/templates/test.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
*/

func removeDuplicates(elements []string) []string { // change string to int here if required
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

// showUsersQuality

func (s *Server) showUsersQuality() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		user := r.Context().Value(ctxKeyUser).(*model.User)
		if user.Role == roleAdministrator {
			statusAdmin = true
			statusLoggedIn = true
		} else if user.Role == roleSuperIngenerQuality {
			statusSuperIngenerQuality = true
			statusLoggedIn = true
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
			"Admin":               statusAdmin,
			"SuperIngenerQuality": statusSuperIngenerQuality,
			"LoggedIn":            statusLoggedIn,
			"GET":                 get,
		}
		err = tpl.ExecuteTemplate(w, "showUsersQuality.html", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

// showUsersQuality

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *Server) Pagination(r *http.Request, limit int) (int, int) {
	keys := r.URL.Query()
	if keys.Get("page") == "" {
		return 1, 0
	}

	page, _ := strconv.Atoi(keys.Get("page"))
	if page < 1 {
		return 1, 0
	}

	begin := (limit * page) - limit
	return page, begin
}

func RenderJSON(w http.ResponseWriter, val interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset-UTF8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(val)
	if err != nil {
		log.Println(err)
		ErrorLogger.Printf(err.Error())
	}
}
