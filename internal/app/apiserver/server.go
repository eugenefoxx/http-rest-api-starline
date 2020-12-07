package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	//	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	//	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
	//	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store/helper"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"

	//	"github.com/gorilla/mux"
	//	"github.com/gorilla/sessions"
	//	_ "github.com/lib/pq"
	//	"github.com/sirupsen/logrus"
	//"github.com/eugenefoxx/http-rest-api-starline/internal/app/apiserver"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
	"github.com/jmoiron/sqlx"

	//	"github.com/gobuffalo/packr/v2/jam/store"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "starline"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID

//	web_DIR = "/web/"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
	tpl                         *template.Template
	LOGFILE                     = "/tmp/apiServer.log"

//	web_DIR                  = "/web/"
)

/*
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
*/
type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
	database     *sqlx.DB
	html         string
}

func init() {
	//	tpl = template.Must(template.ParseGlob("./web/templates/*.html"))
	//	tpl = template.Must(template.New("./web/templates/*.html").Delims("<<", ">>").ParseGlob("./web/templates/*.html"))
	//tpl = template.Must(template.ParseFiles("./web/templates/layout.html").Delims("<<", ">>").ParseGlob("./web/templates/*.html"))
	//	tpl = template.Must(template.ParseGlob("C:/Users/Евгений/templates/*.html"))
	//	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./web/images/"))))

}

func newServer(store store.Store, sessionStore sessions.Store, html string) *server {
	s := &server{
		router:       mux.NewRouter(), // mux.NewRouter()  NewRouter()
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
		html:         html,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)

}

/*
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true) //.StrictSlash(true)

	router.
		PathPrefix(web_DIR).
		Handler(http.StripPrefix(web_DIR, http.FileServer(http.Dir("."+web_DIR))))

	return router
}
*/
func (s *server) configureRouter() {
	//	connStr := "user=postgres password=123 host=localhost dbname=starline sslmode=disable"
	//	db, err := sql.Open("postgres", connStr)
	///	if err != nil {
	//		panic(err)
	//	}
	//	database = db
	//	defer db.Close()
	//	s.router.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./web/images"))))
	// http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./web/images"))))
	//http.Handle("/resource", http.FileServer(http.Dir("./web/images")))
	//	http.Handle("/", s.router)

	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/users", s.pagehandleUsersCreate()).Methods("GET")
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")

	s.router.HandleFunc("/updatepass", s.pageupdateSessionsCreate()).Methods("GET")
	s.router.HandleFunc("/updatepass", s.updateSessionsCreate()).Methods("POST")

	//	s.router.HandleFunc("/sessions", s.redirectMain(s.handleSessionsCreate())).Methods("POST")
	//s.router.HandleFunc("/sessions", s.handleSessionsCreate(s.redirectMain())).Methods("POST")
	s.router.HandleFunc("/login", s.pagehandleSessionsCreate()).Methods("GET")
	s.router.HandleFunc("/login", s.handleSessionsCreate()).Methods("POST")
	//	s.router.HandleFunc("/sessions", s.redirectMain())
	//	s.router.HandleFunc("/sessions", s.pageredirectMain())

	s.router.HandleFunc("/shipmentbysap", s.authMiddleware(s.pageshipmentBySAP())).Methods("GET")
	s.router.HandleFunc("/shipmentbysap", s.authMiddleware(s.shipmentBySAP())).Methods("POST")

	//	s.router.HandleFunc("/showdateshipmentbysap", s.pageshowShipmentBySAP()) //.Methods("GET")
	s.router.HandleFunc("/showdateshipmentbysap", s.authMiddleware(s.showShipmentBySAP())) //.Methods("POST")

	s.router.HandleFunc("/showdateshipmentbysapbysearch", s.authMiddleware(s.pageshowShipmentBySAPBySearch())).Methods("GET")
	s.router.HandleFunc("/showdateshipmentbysapbysearch", s.authMiddleware(s.showShipmentBySAPBySearch())).Methods("POST")
	// pageshowShipmentBySAPBySearchStatic
	s.router.HandleFunc("/showdateshipmentbysapbysearchstatic", s.authMiddleware(s.pageshowShipmentBySAPBySearchStatic())).Methods("GET")
	s.router.HandleFunc("/showdateshipmentbysapbysearchstatic", s.authMiddleware(s.showShipmentBySAPBySearchStatic())).Methods("POST")

	s.router.HandleFunc("/showShipmentbysapbysearchDateStatic", s.authMiddleware(s.pageshowShipmentBySAPBySearchStatic())).Methods("GET")
	s.router.HandleFunc("/showShipmentbysapbysearchDateStatic", s.authMiddleware(s.showShipmentBySAPBySearchDateStatic())).Methods("POST")

	s.router.HandleFunc("/insertIDReturn", s.authMiddleware(s.pageidReturn())).Methods("GET")
	s.router.HandleFunc("/insertIDReturn", s.authMiddleware(s.idReturn())).Methods("POST")

	s.router.HandleFunc("/showIDReturn", s.authMiddleware(s.pageshowIDReturnDataByDate())).Methods("GET")
	s.router.HandleFunc("/showIDReturn", s.authMiddleware(s.showIDReturnDataByDate())).Methods("POST")

	s.router.HandleFunc("/insertvendor", s.authMiddleware(s.pageinsertVendor())).Methods("GET")
	s.router.HandleFunc("/insertvendor", s.authMiddleware(s.insertVendor())).Methods("POST")

	s.router.HandleFunc("/showvendor", s.authMiddleware(s.pageVendor())).Methods("GET")
	s.router.HandleFunc("/updatevendor/{ID:[0-9]+}", s.authMiddleware(s.pageupdateVendor())).Methods("GET")
	s.router.HandleFunc("/updatevendor/{ID:[0-9]+}", s.authMiddleware(s.updateVendor())).Methods("POST")
	s.router.HandleFunc("/deletevendor/{ID:[0-9]+}", s.authMiddleware(s.deleteVendor())) //.Methods("DELETE")

	s.router.HandleFunc("/ininspection", s.authMiddleware(s.pageinInspection())).Methods("GET")
	s.router.HandleFunc("/ininspection", s.authMiddleware(s.inInspection())).Methods("POST")

	s.router.HandleFunc("/statusinspection", s.authMiddleware(s.pageInspection())).Methods("GET")
	s.router.HandleFunc("/updateinspection/{ID:[0-9]+}", s.authMiddleware(s.pageupdateInspection())).Methods("GET")
	s.router.HandleFunc("/updateinspection/{ID:[0-9]+}", s.authMiddleware(s.updateInspection())).Methods("POST")
	s.router.HandleFunc("/deleteinspection/{ID:[0-9]+}", s.authMiddleware(s.deleteInspection()))

	s.router.HandleFunc("/statusinspectionforwh", s.authMiddleware(s.pageListAcceptWHInspection())).Methods("GET")
	s.router.HandleFunc("/acceptinspectiontowh/{ID:[0-9]+}", s.authMiddleware(s.pageacceptWarehouseInspection())).Methods("GET")
	s.router.HandleFunc("/acceptinspectiontowh/{ID:[0-9]+}", s.authMiddleware(s.acceptWarehouseInspection())).Methods("POST")

	s.router.HandleFunc("/testPana", s.authMiddleware(s.testPana())).Methods("GET")
	s.router.HandleFunc("/testIDSAP", s.authMiddleware(s.testIDSAP())).Methods("GET")
	s.router.HandleFunc("/testMB52", s.authMiddleware(s.testMB52())).Methods("GET")
	//	s.router.HandleFunc("/main", s.authMiddleware(s.pageredirectMain())).Methods("GET")
	s.router.HandleFunc("/logout", s.signOut()).Methods("GET")

	s.router.HandleFunc("/hello", s.authMiddleware(s.handleHello()))
	s.router.HandleFunc("/main", s.authMiddleware(s.main())).Methods("GET")
	s.router.HandleFunc("/", s.upload()).Methods("GET")
	//	s.router.HandleFunc("/", s.loginPage())
	s.router.HandleFunc("/js", s.jsPage())

	// /private/***
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")

	s.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web"))))

	//	fmt.Println("Webserver StarlineProduction starting...")

	//	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./web/images"))))

	//	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./web/"))))
	//	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir(config.assets))))
	//	http.Handle("/", http.FileServer(http.Dir("./web/images")))
	//s.router.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./web/"))))
	//	http.Handle("/", s.router)
	fmt.Println("Webserver StarLine launch.")
	//	open.StartWith("http://localhost:3001/", "google-chrome-stable") // chromium

}

func (s *server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}

func (s *server) upload() http.HandlerFunc {
	tpl, err := template.New("base").ParseFiles(s.html + "layout1.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "base", nil)
	}
}

func (s *server) main() http.HandlerFunc {
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles(s.html + "layout.html"))
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles(s.html+"layout1.html", s.html+"login1.html"))
	tpl = template.Must(template.New("base").ParseFiles(s.html + "layout1.html"))
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/layout.html"))
	//tpl = template.Must(template.ParseFiles("web/templates/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		admin := false
		stockkeeper := false
		superIngenerQuality := false
		stockkeeperWH := false
		inspector := false
		LoggedIn := false

		//	u := s.authenticateUser()
		//	fmt.Println(u)
		//	var body, _ = helper.LoadFile("./web/templates/index.html")
		//	fmt.Fprintf(w, body)
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
		fmt.Println("/main - user:", u.Email, u.ID, u.Role)
		if u.Role == "Administrator" {
			admin = true
			LoggedIn = true

		} else if u.Role == "кладовщик" {
			stockkeeper = true
			LoggedIn = true

		} else if u.Role == "SuperIngenerQuality" {
			superIngenerQuality = true
			LoggedIn = true

		} else if u.Role == "кладовщик склада" {
			stockkeeperWH = true
			LoggedIn = true

		} else if u.Role == "контроллер качества" {
			inspector = true
			LoggedIn = true

		}
		GET := map[string]bool{
			"admin":               admin,
			"stockkeeper":         stockkeeper,
			"SuperIngenerQuality": superIngenerQuality,
			"stockkeeperWH":       stockkeeperWH,
			"inspector":           inspector,
		}

		data := map[string]interface{}{
			"user": u.LastName,
			"id":   u.FirstName,
			"GET":  GET,
			//	"admin":               admin,
			//	"stockkeeper":         stockkeeper,
			//	"SuperIngenerQuality": superIngenerQuality,
			//	"stockkeeperWH":       stockkeeperWH,
			//	"inspector":           inspector,
			"LoggedIn": LoggedIn,
		}

		//tpl.ExecuteTemplate(w, "home.html", data)
		//tpl.ExecuteTemplate(w, "layout", data)
		tpl.ExecuteTemplate(w, "base", data)
	}
}

func (s *server) jsPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "js.html", nil)
	}
}

func (s *server) loginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "register.html", nil) // register.html
	}
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// open a file
		f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Printf("error opening file: %v", err)
		}

		// don't forget to close it
		defer f.Close()

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

		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed witn %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
		/*
			var log = logrus.New()
			log.Out = os.Stdout
			file, err := os.OpenFile("/tmp/logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				log.Out = file
			} else {
				log.Info("Failed to log to file, using default stderr")
			}
		*/
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
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
		fmt.Println("???", u.Email)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) authMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id := session.Values["user_id"]
		fmt.Println("id is:", id)
		/*
			if !ok {
				//	s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
				http.Redirect(w, r, "/sessions", http.StatusSeeOther)
				return
			}
		*/
		/*
			//	email := session.Values["user_id"]
			u, err := s.store.User().Find(id.(int))
			fmt.Println("u is:", u)
			if err != nil {
				s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
				//http.Redirect(w, r, "/sessions", http.StatusSeeOther)
				return
			} */

		if id == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			//tpl.ExecuteTemplate(w, "login.html", nil)
		} else {
			h.ServeHTTP(w, r)
		}

	})
}

func ho(r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(*model.User)
	fmt.Println(user)
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser).(*model.User)
		fmt.Println(user)
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) pagehandleUsersCreate() http.HandlerFunc {
	//	tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/register.html"))
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/register.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "register.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/register.html")
		//	fmt.Fprintf(w, body)
		err = tpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	}

	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "register.html")
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/register.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		//	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		//		s.error(w, r, http.StatusBadRequest, err)
		//		return
		//	}
		fmt.Println("зашел на регистрацию")
		//	email := r.FormValue("email")
		//	password := r.FormValue("password")
		//	firstname := r.FormValue("firstname")
		//	lastname := r.FormValue("lastname")

		req.Email = r.FormValue("email")
		//fmt.Println(req.Email)

		req.Password = r.FormValue("password")

		req.FirstName = r.FormValue("firstname")

		req.LastName = r.FormValue("lastname")
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
		//tpl.ExecuteTemplate(w, "login.html", nil)
		err = tpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) pageupdateSessionsCreate() http.HandlerFunc {
	//	tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/login.html"))
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/login.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateUser.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/login.html") // "./web/templates/login.html"
		//	fmt.Fprintf(w, body)
		err = tpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) updateSessionsCreate() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
		//	Tabel    string `json:"tabel"`
		Password string `json:"password"`
	}
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateUser.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			fmt.Println(err)
			//	return
		}
		defer f.Close()
		iLog := log.New(f, "func updateSessionsCreate ", log.LstdFlags)
		iLog.SetFlags(log.LstdFlags | log.Lshortfile)

		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")

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
		http.Redirect(w, r, "/", 303)
		err = tpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) pageredirectMain() http.HandlerFunc {
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/layout.html"))
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/layout.html")
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "layout.html")
	//if err != nil {
	//	panic(err)
	//}
	tpl = template.Must(template.ParseFiles(s.html + "layout1.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/index.html")
		//	fmt.Fprintf(w, body)
		//	err = tpl.ExecuteTemplate(w, "layout", nil)
		//	if err != nil {
		//		http.Error(w, err.Error(), 400)
		//		return
		//	}
		tpl.Execute(w, nil)
	}
}

// h http.HandlerFunc
func (s *server) redirectMain() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/index.html")
		//	fmt.Fprintf(w, body)
		//s.router.HandleFunc("/main", s.main())
		http.Redirect(w, r, "/", 303)
	})
}

func (s *server) pagehandleSessionsCreate() http.HandlerFunc {
	//	tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/login.html"))
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/login.html")
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "login.html")
	tpl, err := template.New("base").ParseFiles(s.html+"layout1.html", s.html+"login1.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/login.html") // "./web/templates/login.html"
		//	fmt.Fprintf(w, body)
		//	err = tpl.ExecuteTemplate(w, "layout", nil)
		//	if err != nil {
		//		http.Error(w, err.Error(), 400)
		//		return
		//	}
		tpl.ExecuteTemplate(w, "base", nil)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/layout.html")
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "layout.html")
	//if err != nil {
	//	panic(err)
	//}
	tpl, err := template.New("base").ParseFiles(s.html + "layout1.html")
	if err != nil {
		panic(err)
	}
	//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			fmt.Println(err)
			//	return
		}
		defer f.Close()
		iLog := log.New(f, "func handleSessionsCreate ", log.LstdFlags)
		iLog.SetFlags(log.LstdFlags | log.Lshortfile)
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
		//		if r.Method == http.MethodPost {
		admin := false
		stockkeeper := false
		superIngenerQuality := false
		stockkeeperWH := false
		inspector := false
		LoggedIn := false

		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Println("email:   ", email)
		//	target := "/sessions"

		u, err := s.store.User().FindByEmail(email, email)
		//	fmt.Println("FindByEmail:   ", u)
		//	match := u.ComparePassword(password)
		//	fmt.Println("Match:   ", match)
		if err != nil || !u.ComparePassword(password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)

			return
		}

		if u.Role == "Administrator" {
			admin = true
			LoggedIn = true

		} else if u.Role == "кладовщик" {
			stockkeeper = true
			LoggedIn = true

		} else if u.Role == "SuperIngenerQuality" {
			superIngenerQuality = true
			LoggedIn = true
		} else if u.Role == "кладовщик склада" {
			stockkeeperWH = true
			LoggedIn = true

		} else if u.Role == "контроллер качества" {
			inspector = true
			LoggedIn = true

		}

		GET := map[string]bool{
			"admin":               admin,
			"stockkeeper":         stockkeeper,
			"SuperIngenerQuality": superIngenerQuality,
			"stockkeeperWH":       stockkeeperWH,
			"inspector":           inspector,
		}

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
		iLog.Println("залогинился пользователь", u.Email)
		//	s.redirectMain()
		//	s.pageredirectMain()
		//	_, ok := session.Values["user_id"]
		//	if !ok {

		//	return
		//	}
		data := map[string]interface{}{
			"user": u.LastName,
			"id":   u.FirstName,
			"GET":  GET,
			//	"admin":               admin,
			//	"stockkeeper":         stockkeeper,
			//	"SuperIngenerQuality": superIngenerQuality,
			//	"stockkeeperWH":       stockkeeperWH,
			//	"inspector":           inspector,
			"LoggedIn": LoggedIn,
		}
		tpl.ExecuteTemplate(w, "base", data) //  "index.html"
		//err = tpl.ExecuteTemplate(w, "base", data)
		//if err != nil {
		//	http.Error(w, err.Error(), 400)
		//	return
		//}
		http.Redirect(w, r, "/main", http.StatusFound)
	}
}

func (s *server) pageshipmentBySAP() http.HandlerFunc {
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertsapbyship6.html"))
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertsapbyship6.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertsapbyship6.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/insertsapbyship6.html")
		//	fmt.Fprintf(w, body)
		data := map[string]interface{}{
			"user": "Я тут",
		}
		err = tpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
			//	}
			//		if err = tpl.ExecuteTemplate(w, "layout", nil); err != nil {
			//			s.error(w, r, http.StatusUnprocessableEntity, err)
			//			return
			//		}
		}
	}
}

func (s *server) shipmentBySAP() http.HandlerFunc {
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
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertsapbyship6.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	hdata := reqA{}
		//	var hdata ReqA

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		var hdata []reqA

		json.Unmarshal(body, &hdata)
		//	json.Marshal(body)
		fmt.Printf("body json: %s", body)
		//	fmt.Printf("1 тест %s", &hdata.Material)
		//	fmt.Println("tect2 %s", hdata, "\n")
		fmt.Println("\njson  struct hdata", hdata)
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
		//	hdata.ID = user.ID
		//	hdata.LastName = user.LastName
		/*
			var record []request
			record.Material = hdata.Material
			//	fmt.Printf("2 тест %v", record.Material)
			record.Material = hdata.Qty
			record.Comment = hdata.Comment
			record.ID = user.ID
			record.LastName = user.LastName
			//	user2 := strconv.Atoi(user1)
		*/

		for _, v := range hdata {
			fmt.Println(v.Material, v.Qty, v.Comment, user.ID, user.LastName)
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
		err = tpl.ExecuteTemplate(w, "layout", nil)
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

/*
type rawTime []byte

func (t rawTime) Time() (time.Time, error) {
	return time.Parse("15:04:05", string(t))
}

type rawDate []byte

func (t rawDate) Time() (time.Time, error) {
	return time.Parse("2020-02-10", string(t))
}
*/
func (s *server) pageshowShipmentBySAP() http.HandlerFunc {
	// tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysap2.html"))
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysap2.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdatebysap2.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/showdatebysap2.html")
		//	fmt.Fprintf(w, body)
		err = tpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) showShipmentBySAP() http.HandlerFunc {
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysap2.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdatebysap2.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {

		get, err := s.store.Shipmentbysap().ShowDate()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		/*
			err = tpl.ExecuteTemplate(w, "showdatebysap2.html", get)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		*/
		err = tpl.ExecuteTemplate(w, "layout", get)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}

}

func (s *server) pageshowShipmentBySAPBySearch() http.HandlerFunc {
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysapbysearch.html"))
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysapbysearch.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdatebysapbysearch.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/showdatebysapbysearch.html")
		//	fmt.Fprintf(w, body)
		err = tpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

type SearchBy struct {
	LastName string     `json:"lastname, string"`
	Date1    CustomDate `json:"date1, string"`
	Date2    CustomDate `json:"date2, string"`
	Material int        `json:"material, string"`
}

/*
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}
*/
/*
// first create a type alias
type JsonDate time.Time

// imeplement Marshaler und Unmarshalere interface
func (j *JsonDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"") // s := strings.Trim(string(b), "\"")

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	fmt.Println("parse", t)
	*j = JsonDate(t)
	return nil
}
*/

type CustomDate struct {
	time.Time
}

const layout = "2006-01-02"

func (c *CustomDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // remove quotes
	if s == "null" {
		return
	}
	c.Time, err = time.Parse(layout, s)
	return
}

func (c CustomDate) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format(layout))), nil
}

func (s *server) showShipmentBySAPBySearch() http.HandlerFunc {
	/*
		type SearchBy struct {
			LastName string `json:"lastname"`
			Date1    string `json:"date1"`
			Date2    string `json:"date2"`
			Material int    `json:"material"`
		}
	*/
	type Request struct {
		Material int    `db:"material"`
		Qty      int    `db:"qty"`
		Comment  string `db:"comment"`
		ID       int    `db:"id"`
		LastName string `db:"lastname"`
	}

	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysap.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdatebysap.html")
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		//	body, err := r.Body  // ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		//	var searchdata []SearchBy
		//	searchdata := []SearchBy{}
		var searchdata []SearchBy
		const layoutISO = "2006-01-02"

		err1 := json.Unmarshal(body, &searchdata)
		if err1 != nil {
			fmt.Println("Error unmarshalling data: ", err)
			return
		}
		//	UnmarshalJSON(body)
		fmt.Printf("test searchData r.Body %s", body)
		fmt.Println("\nall of the data searchBy", searchdata)
		//	fmt.Println(searchdata.Date1.Format(time.RFC3339))
		for i, v := range searchdata {
			fmt.Println(i, v)
			fmt.Println("LastName", v.LastName)
			fmt.Println("ниже читаю формат Date1\n")
			//	fmt.Println("v.Date1.Format(layoutISO)", v.Date1.Format(time.RFC3339)) // fmt.Println(v.Date1.Format(layoutISO))
			//	fmt.Println(v.Date1.GobDecode)
			fmt.Println(v.Date1)

			//	fmt.Println("\t", v.Date1.Format(time.RFC3339))
			//	date := fmtdate.ParseTime("DD-MM-YYYY", v.Date1)
			//	fmt.Println(time.Parse(layoutISO, v.Date1))
			//	fmt.Println(date)

		}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&searchdata)
		fmt.Printf("%+v\n", searchdata)

		//	tpl.ExecuteTemplate(w, "showdatebysap.html", nil)
		err = tpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}

}

/*
// UnmarshalJSON ...
func (l *SearchBy) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "lastname" {
			l.LastName = v
		}

		if strings.ToLower(k) == "date1" {
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return err
			}
			l.Date1 = t
		}

		if strings.ToLower(k) == "date2" {
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return err
			}
			l.Date1 = t
		}

		if strings.ToLower(k) == "material" {
			l.Material, err = strconv.Atoi(v)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
*/
type MyTime struct {
	time.Time
}

const dateFormat = "2006-01-02"

/*
func (m *MyTime) UnmarshalJSON(p []byte) error {
	t, err := time.Parse(dateFormat, strings.Replace(
		string(p),
		"\"",
		"",
		-1,
	))

	if err != nil {
		return err
	}

	m.Time = t

	return nil
}
*/

func (s *server) pageshowShipmentBySAPBySearchStatic() http.HandlerFunc {
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysapbysearchstatic.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdatebysapbysearchstatic.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/showdatebysapbysearchstatic.html")
		//	fmt.Fprintf(w, body)
		data := map[string]interface{}{
			"user": "Я тут",
		}
		if err = tpl.ExecuteTemplate(w, "layout", data); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

	}
}

func (s *server) showShipmentBySAPBySearchStatic() http.HandlerFunc {
	type searchBy struct {
		LastName string `json:"lastname"`
		Date1    string `json:"date1"`
		Date2    string `json:"date2"`
		Material int    `json:"material"`
	}

	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysap2.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdatebysap2.html")
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		search := &searchBy{}
		//		r.ParseForm()
		materialInt, err := strconv.Atoi(r.FormValue("material"))
		if err != nil {
			fmt.Println(err)
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

		//	err = tpl.ExecuteTemplate(w, "showdatebysap2.html", get)
		err = tpl.ExecuteTemplate(w, "layout", get)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) showShipmentBySAPBySearchDateStatic() http.HandlerFunc {
	type searchBy struct {
		Date1 string `json:"date1"`
		Date2 string `json:"date2"`
	}
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysap2.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdatebysap2.html")
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		search := &searchBy{}
		//		r.ParseForm()

		search.Date1 = r.FormValue("date1")
		fmt.Println("date1 - ", search.Date1)
		search.Date2 = r.FormValue("date2")
		fmt.Println("date2 - ", search.Date2)

		get, err := s.store.Shipmentbysap().ShowDataByDate(search.Date1, search.Date2)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		//	err = tpl.ExecuteTemplate(w, "showdatebysap2.html", get)
		err = tpl.ExecuteTemplate(w, "layout", get)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) pageidReturn() http.HandlerFunc {
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertidreturn2.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertidreturn2.html")
	if err != nil {
		panic(err)
	}
	//	tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertidreturn2.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/insertidreturn2.html")
		//	fmt.Fprintf(w, body)
		err = tpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) idReturn() http.HandlerFunc {
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
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertidreturn2.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		var rdata []req

		json.Unmarshal(body, &rdata)
		fmt.Printf("test %s", body)
		fmt.Println("\nall of the rdata", rdata)

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

		for _, v := range rdata {
			sap := v.ScanID[1:8]
			material := v.Material
			material, err := strconv.Atoi(sap)
			if err != nil {
				fmt.Println(err)
			}
			idsap := v.ScanID[20:30]
			idroll := v.IDRoll
			idroll, err1 := strconv.Atoi(idsap)
			fmt.Println("idroll в 1-м цикле -", idroll)
			if err != nil {
				fmt.Println(err1)
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
					//	fmt.Fprintf(w, "не верное сканирование :"+v.ScanID)
				}
				//	tpl.Execute(w, data)
				return
			}

		}
		//	data := map[string]interface{}{
		//		"fuck": "OK",
		//	}
		err = tpl.ExecuteTemplate(w, "layout", nil)
		//	tpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

	}
}

func (s *server) pageshowIDReturnDataByDate() http.HandlerFunc {
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysapbysearchstatic.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdatebysapbysearchstatic.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/showdatebysapbysearchstatic.html")
		//	fmt.Fprintf(w, body)
		err = tpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) showIDReturnDataByDate() http.HandlerFunc {
	type searchBy struct {
		Date1 string `json:"date1"`
		Date2 string `json:"date2"`
	}

	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdateidreturn.html")
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdateidreturn.html")
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		s.store.HUMOSAPStock().ImportDate()
		s.store.MB52SAPStock().ImportDate()
		s.store.PanacimStock().ImportDate()

		search := &searchBy{}
		//		r.ParseForm()

		search.Date1 = r.FormValue("date1")
		fmt.Println("date1 - ", search.Date1)
		search.Date2 = r.FormValue("date2")
		fmt.Println("date2 - ", search.Date2)

		get, err := s.store.Showdateidreturn().ShowDataByDate(search.Date1, search.Date2)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		//	err = tpl.ExecuteTemplate(w, "showdateidreturn.html", get)
		err = tpl.ExecuteTemplate(w, "layout", get)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) testPana() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.store.PanacimStock().ImportDate()
		/*	_, err := s.store.PanacimStock().ImportDate()
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			//	tpl.ExecuteTemplate(w, "testpana.html", nil)

			//	err = tpl.ExecuteTemplate(w, "testpana.html", get)
			//	if err != nil {
			//		http.Error(w, err.Error(), 400)
			//		return
				} */
	}
}

func (s *server) testIDSAP() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.store.HUMOSAPStock().ImportDate()
	}
}

func (s *server) testMB52() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.store.MB52SAPStock().ImportDate()
	}
}

func (s *server) pageinsertVendor() http.HandlerFunc {
	/*tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertvendor.html")
	if err != nil {
		panic(err)
	}*/
	tpl = template.Must(template.New("base").ParseFiles(s.html+"layout1.html", s.html+"insertvendor1.html"))
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
		fmt.Println("Check -")
		tpl.ExecuteTemplate(w, "base", nil)
	}
}

func (s *server) insertVendor() http.HandlerFunc {
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
	tpl = template.Must(template.New("base").ParseFiles(s.html+"layout1.html", s.html+"insertvendor1.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Check2 -")
		admin := false
		superIngenerQuality := false
		LoggedIn := false

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		var hdata []requestFrom
		json.Unmarshal(body, &hdata)
		fmt.Printf("body json: %s", body)
		fmt.Println("\njson  struct hdata", hdata)

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

		GET := map[string]bool{
			"admin": admin,
			//	"stockkeeper":         stockkeeper,
			"SuperIngenerQuality": superIngenerQuality,
			//	"stockkeeperWH":       stockkeeperWH,
			//	"inspector":           inspector,
		}

		if user.Role == "Administrator" {
			admin = true
			LoggedIn = true
		} else if user.Role == "SuperIngenerQuality" {
			superIngenerQuality = true
			LoggedIn = true
			fmt.Println("SuperIngenerQuality - ", superIngenerQuality)
		}

		for _, v := range hdata {
			fmt.Println(v.CodeDebitor, v.NameDebitor)

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
			//	"admin":               admin,
			//	"SuperIngenerQuality": superIngenerQuality,
			"GET":      GET,
			"LoggedIn": LoggedIn,
		}

		//err = tpl.ExecuteTemplate(w, "layout", data)
		//	err = tpl.ExecuteTemplate(w, "layout", v)
		err = tpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

	}

}

func (s *server) pageVendor() http.HandlerFunc {
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showvendor.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {

		//	var body, _ = helper.LoadFile("./web/templates/insertsapbyship6.html")
		//	fmt.Fprintf(w, body)
		//data := map[string]interface{}{
		//	"user": "Я тут",
		//}
		get, err := s.store.Vendor().ListVendor()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		err = tpl.ExecuteTemplate(w, "layout", get)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
			//	}
			//		if err = tpl.ExecuteTemplate(w, "layout", nil); err != nil {
			//			s.error(w, r, http.StatusUnprocessableEntity, err)
			//			return
			//		}
		}
	}
}

func (s *server) pageupdateVendor() http.HandlerFunc {
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updatevendor.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}
		//fmt.Println("ID - ?", id)
		get, err := s.store.Vendor().EditVendor(id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		err = tpl.ExecuteTemplate(w, "layout", get)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) updateVendor() http.HandlerFunc {
	type request struct {
		ID          int    `json:"ID"`
		CodeDebitor string `json:"codedebitor"`
		NameDebitor string `json:"namedebitor"`
	}
	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updatevendor.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}

		req.ID = id
		req.CodeDebitor = r.FormValue("codedebitor")
		req.NameDebitor = r.FormValue("namedebitor")

		u := &model.Vendor{
			ID:          req.ID,
			CodeDebitor: req.CodeDebitor,
			NameDebitor: req.NameDebitor,
		}

		if err := s.store.Vendor().UpdateVendor(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		/*	err = tpl.ExecuteTemplate(w, "layout", nil)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}*/
		http.Redirect(w, r, "/showvendor", 303)
	}
}

func (s *server) deleteVendor() http.HandlerFunc {
	type request struct {
		ID          int    `json:"ID"`
		CodeDebitor string `json:"codedebitor"`
		NameDebitor string `json:"namedebitor"`
		SPPElement  string `json:"sppelement"`
	}
	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updatevendor.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}

		req.ID = id
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
		http.Redirect(w, r, "/showvendor", 303)
	}
}

func (s *server) pageinInspection() http.HandlerFunc {
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "ininspection.html")
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err = tpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) inInspection() http.HandlerFunc {
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

	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "ininspection.html")
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		var rdata []req
		json.Unmarshal(body, &rdata)
		fmt.Printf("test ininspection %s", body)
		fmt.Println("\nall of the rdata ininspection", rdata)

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
		}

		err = tpl.ExecuteTemplate(w, "layout", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

	}
}

func (s *server) pageInspection() http.HandlerFunc {
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showinspection.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		stockkeeperWH := true
		quality := false
		superIngenerQuality := false

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

		if user.Role == "SuperIngenerQuality" {
			superIngenerQuality = true
			fmt.Println("pageInspection SuperIngenerQuality - ", superIngenerQuality)
		} else if user.Groups == "качество" {
			quality = true
			fmt.Println("pageInspection quality - ", quality)
		} else if user.Groups == "склад" {
			stockkeeperWH = false
		}

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
			"quality":             quality,
			"warehouse":           stockkeeperWH,
			"SuperIngenerQuality": superIngenerQuality,
			"GET":                 get,
			"CountTotal":          countTotal,
			"HoldInspection":      holdInspection,
			"NotVerifyComponents": notVerifyComponents,
			"GetStatic":           getStatic,
			"HoldCountDebitor":    holdCountDebitor,
			"NotVerifyDebitor":    notVerifyDebitor,
		}

		err = tpl.ExecuteTemplate(w, "layout", groups)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) pageupdateInspection() http.HandlerFunc {
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateinspection.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {

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
		//fmt.Println("ID - ?", id)
		get, err := s.store.Inspection().EditInspection(id, user.Groups)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		err = tpl.ExecuteTemplate(w, "layout", get)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) updateInspection() http.HandlerFunc {
	type request struct {
		ID     int    `json:"ID"`
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateinspection.html")
	if err != nil {
		panic(err)
	}

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

func (s *server) deleteInspection() http.HandlerFunc {
	type request struct {
		ID int `json:"ID"`
	}
	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateinspection.html")
	if err != nil {
		panic(err)
	}
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
func (s *server) pageListAcceptWHInspection() http.HandlerFunc { // acceptinspection.html showinspection.html
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "acceptinspection.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
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
		}
		/*
			if user.Groups == "качество" {
				quality = true
			} else if user.Groups == "склад" {
				stockkeeperWH = true
			}
		*/
		get, err := s.store.Inspection().ListAcceptWHInspection()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		groups := map[string]interface{}{
			//	"quality":   quality,
			"warehouse": stockkeeperWH,
			//	"SuperIngenerQuality": superIngenerQuality,
			"GET": get,
			//	"status":    statusStr,
		}

		err = tpl.ExecuteTemplate(w, "layout", groups)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) pageacceptWarehouseInspection() http.HandlerFunc {
	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "acceptWarehouseInspection.html")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
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

		//fmt.Println("ID - ?", id)
		get, err := s.store.Inspection().EditAcceptWarehouseInspection(id, user.Groups)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		err = tpl.ExecuteTemplate(w, "layout", get)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) acceptWarehouseInspection() http.HandlerFunc {
	type request struct {
		ID       int    `json:"ID"`
		Location string `json:"location"`
	}
	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "acceptWarehouseInspection.html")
	if err != nil {
		panic(err)
	}
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
		http.Redirect(w, r, "/statusinspectionforwh", 303)
	}
}

func (s *server) signOut() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("Пытаюсь чистить куки")
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

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
