package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"

	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
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
	LOGFILE                     = "/tmp/apiServer.log"
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
	database     *sqlx.DB
	//	html         string
}

func init() {
	tpl = template.Must(tpl.ParseGlob("web/templates/*.html"))
	//tpl = template.Must(tpl.ParseGlob("/home/eugenearch/Code/github.com/eugenefoxx/http-rest-api-starline/web/templates/*.html"))
}

//func newServer(store store.Store, sessionStore sessions.Store, html string) *server {
func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(), // mux.NewRouter()  NewRouter()
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
		//	html:         html,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)

}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/test", s.diaplayPage())

	s.router.HandleFunc("/users", s.pagehandleUsersCreate()).Methods("GET")
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")

	s.router.HandleFunc("/showusersquality", s.pageshowUsersQuality()).Methods("GET")
	//s.router.HandleFunc("/showusersquality", s.showUsersQuality()).Methods("POST")
	s.router.HandleFunc("/createusersquality", s.createUserQuality()).Methods("POST")
	s.router.HandleFunc("/updateuserquality/{ID:[0-9]+}", s.authMiddleware(s.pageupdateUserQuality())).Methods("GET")
	s.router.HandleFunc("/updateuserquality/{ID:[0-9]+}", s.authMiddleware(s.updateUserQuality())).Methods("POST")
	s.router.HandleFunc("/deleteuserquality/{ID:[0-9]+}", s.authMiddleware(s.deleteUserQuality()))

	s.router.HandleFunc("/showuserswarehouse", s.pageshowUsersWarehouse()).Methods("GET")
	s.router.HandleFunc("/createuserswarehouse", s.createUserWarehouse()).Methods("POST")
	s.router.HandleFunc("/updateuserwarehouse/{ID:[0-9]+}", s.authMiddleware(s.pageupdateUserWarehouse())).Methods("GET")
	s.router.HandleFunc("/updateuserwarehouse/{ID:[0-9]+}", s.authMiddleware(s.updateUserWarehouse())).Methods("POST")
	s.router.HandleFunc("/deleteuserwarehouse/{ID:[0-9]+}", s.authMiddleware(s.deleteUserWarehouse()))

	s.router.HandleFunc("/updatepass", s.pageupdateSessionsCreate()).Methods("GET")
	s.router.HandleFunc("/updatepass", s.updateSessionsCreate()).Methods("POST")

	s.router.HandleFunc("/login", s.pagehandleSessionsCreate()).Methods("GET")
	s.router.HandleFunc("/login", s.handleSessionsCreate()).Methods("POST")
	s.router.HandleFunc("/logout", s.signOut()).Methods("GET")

	s.router.HandleFunc("/shipmentbysap", s.authMiddleware(s.pageshipmentBySAP())).Methods("GET")
	s.router.HandleFunc("/shipmentbysap", s.authMiddleware(s.shipmentBySAP())).Methods("POST")
	s.router.HandleFunc("/showdateshipmentbysapbysearchstatic", s.authMiddleware(s.pageshowShipmentBySAPBySearchStatic())).Methods("GET")
	s.router.HandleFunc("/showdateshipmentbysapbysearchstatic", s.authMiddleware(s.showShipmentBySAPBySearchStatic())).Methods("POST")

	s.router.HandleFunc("/insertIDReturn", s.authMiddleware(s.pageidReturn())).Methods("GET")
	s.router.HandleFunc("/insertIDReturn", s.authMiddleware(s.idReturn())).Methods("POST")

	s.router.HandleFunc("/insertvendor", s.authMiddleware(s.pageinsertVendor())).Methods("GET")
	s.router.HandleFunc("/insertvendor", s.authMiddleware(s.insertVendor())).Methods("POST")

	s.router.HandleFunc("/showvendor", s.authMiddleware(s.pageVendor())).Methods("GET")
	s.router.HandleFunc("/updatevendor/{ID:[0-9]+}", s.authMiddleware(s.pageupdateVendor())).Methods("GET")
	s.router.HandleFunc("/updatevendor/{ID:[0-9]+}", s.authMiddleware(s.updateVendor())).Methods("POST")
	s.router.HandleFunc("/deletevendor/{ID:[0-9]+}", s.authMiddleware(s.deleteVendor())) //.Methods("DELETE")

	s.router.HandleFunc("/ininspection", s.authMiddleware(s.pageinInspection())).Methods("GET")
	s.router.HandleFunc("/ininspection", s.authMiddleware(s.inInspection())).Methods("POST")
	s.router.HandleFunc("/historyinspection", s.authMiddleware(s.pagehistoryInspection())).Methods("GET")
	s.router.HandleFunc("/historyinspection", s.authMiddleware(s.historyInspection())).Methods("POST")

	s.router.HandleFunc("/statusinspection", s.authMiddleware(s.pageInspection())).Methods("GET")
	s.router.HandleFunc("/updateinspection/{ID:[0-9]+}", s.authMiddleware(s.pageupdateInspection())).Methods("GET")
	s.router.HandleFunc("/updateinspection/{ID:[0-9]+}", s.authMiddleware(s.updateInspection())).Methods("POST")
	s.router.HandleFunc("/deleteinspection/{ID:[0-9]+}", s.authMiddleware(s.deleteInspection()))

	s.router.HandleFunc("/statusinspectionforwh", s.authMiddleware(s.pageListAcceptWHInspection())).Methods("GET")
	s.router.HandleFunc("/acceptinspectiontowh/{ID:[0-9]+}", s.authMiddleware(s.pageacceptWarehouseInspection())).Methods("GET")
	s.router.HandleFunc("/acceptinspectiontowh/{ID:[0-9]+}", s.authMiddleware(s.acceptWarehouseInspection())).Methods("POST")

	s.router.HandleFunc("/main", s.authMiddleware(s.main())).Methods("GET")
	s.router.HandleFunc("/", s.upload()).Methods("GET")
	s.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	fmt.Println("Webserver StarLine launch.")

}

type Page struct {
	TitleDOC, Title, Content, Navbar, User, Username                 string
	LoggedIn, Admin, Stockkeeper, SuperIngenerQuality, StockkeeperWH bool
	Inspector, Quality                                               bool
}

func (s *server) diaplayPage() http.HandlerFunc {
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

func (s *server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}

func (s *server) pagehandleUsersCreate() http.HandlerFunc {
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

func (s *server) handleUsersCreate() http.HandlerFunc {
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
		tpl.ExecuteTemplate(w, "register.html", nil)
		///	err = tpl.ExecuteTemplate(w, "layout", nil)
		///	if err != nil {
		///		http.Error(w, err.Error(), 400)
		///		return
		///	}
	}
}

func (s *server) pageupdateSessionsCreate() http.HandlerFunc {
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

func (s *server) updateSessionsCreate() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
		//	Tabel    string `json:"tabel"`
		Password string `json:"password"`
	}
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updateUser.html")
	///	if err != nil {
	///		panic(err)
	///	}
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
		err = tpl.ExecuteTemplate(w, "updateUser.html", nil)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) upload() http.HandlerFunc {
	//	type Page struct {
	//		TitleDOC, Navbar string
	//		LoggedIn         bool
	//	}
	return func(w http.ResponseWriter, r *http.Request) {
		//tpl.ExecuteTemplate(w, "base", nil)
		//fmt.Print("Test Upload")
		p := &Page{
			TitleDOC: "Start",
			//	Navbar:   "ContentNav",
			LoggedIn: false,
		}
		tpl.ExecuteTemplate(w, "index0.html", p)
		//tpl.Execute(w, nil)
	}
}

func (s *server) main() http.HandlerFunc {
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles(s.html + "layout.html"))
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles(s.html+"layout1.html", s.html+"login1.html"))
	//tpl = template.Must(template.New("base").ParseFiles(s.html + "layout1.html"))
	//tpl = template.Must(template.ParseFiles(s.html+"header.html", s.html+"footer.html"))
	//tpl = template.ParseFiles(s.html + "*.html")
	//tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/layout.html"))
	//tpl = template.Must(template.ParseFiles("web/templates/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
		Stockkeeper := false
		WarehouseManager := false
		SuperIngenerQuality := false
		IngenerQuality := false
		StockkeeperWH := false
		Inspector := false
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
			Admin = true
			LoggedIn = true

		} else if u.Role == "кладовщик" {
			Stockkeeper = true
			LoggedIn = true

		} else if u.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true

		} else if u.Role == "инженер по качеству" {
			IngenerQuality = true
			LoggedIn = true

		} else if u.Role == "кладовщик склада" {
			StockkeeperWH = true
			LoggedIn = true

		} else if u.Role == "контролер качества" {
			Inspector = true
			LoggedIn = true

		} else if u.Role == "старший кладовщик склада" {
			WarehouseManager = true
			LoggedIn = true
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
			"TitleDOC": "MAIN",
			"User":     u.LastName,
			"Username": u.FirstName,
			//"GET":                 GET,
			"Admin":               Admin,
			"Stockkeeper":         Stockkeeper,
			"SuperIngenerQuality": SuperIngenerQuality,
			"WarehouseManager":    WarehouseManager,
			"IngenerQuality":      IngenerQuality,
			"StockkeeperWH":       StockkeeperWH,
			"Inspector":           Inspector,
			"LoggedIn":            LoggedIn,
		}

		//tpl.ExecuteTemplate(w, "home.html", data)
		//tpl.ExecuteTemplate(w, "layout", data)
		//tpl.ExecuteTemplate(w, "base", data)
		tpl.ExecuteTemplate(w, "index.html", data) // index3.html
	}
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
			//tpl.ExecuteTemplate(w, "login.html", nil)
		} else {
			h.ServeHTTP(w, r)
		}

	})
}

func (s *server) pagehandleSessionsCreate() http.HandlerFunc {
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
	//tpl, err := template.New("base").ParseFiles(s.html + "layout1.html")
	//if err != nil {
	//	panic(err)
	//}
	//tpl = template.Must(template.ParseFiles(s.html + "/*.html"))
	//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//tpl = template.Must(template.ParseFiles(s.html+"header.html", s.html+"footer.html"))
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
		//		if r.Method == http.MethodPost
		Admin := false
		Stockkeeper := false
		SuperIngenerQuality := false
		WarehouseManager := false
		IngenerQuality := false
		StockkeeperWH := false
		Inspector := false
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
			Admin = true
			LoggedIn = true

		} else if u.Role == "кладовщик" {
			Stockkeeper = true
			LoggedIn = true

		} else if u.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true
		} else if u.Role == "инженер по качеству" {
			IngenerQuality = true
			LoggedIn = true
		} else if u.Role == "кладовщик склада" {
			StockkeeperWH = true
			LoggedIn = true

		} else if u.Role == "контролер качества" {
			Inspector = true
			LoggedIn = true

		} else if u.Role == "старший кладовщик склада" {
			WarehouseManager = true
			LoggedIn = true
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
		iLog.Println("залогинился пользователь", u.Email)
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
			"Admin":               Admin,
			"Stockkeeper":         Stockkeeper,
			"SuperIngenerQuality": SuperIngenerQuality,
			"WarehouseManager":    WarehouseManager,
			"IngenerQuality":      IngenerQuality,
			"StockkeeperWH":       StockkeeperWH,
			"Inspector":           Inspector,
			"LoggedIn":            LoggedIn,
		}
		//	tpl.ExecuteTemplate(w, "base", data) //  "index.html"
		tpl.ExecuteTemplate(w, "index.html", data)
		//err = tpl.ExecuteTemplate(w, "base", data)
		//if err != nil {
		//	http.Error(w, err.Error(), 400)
		//	return
		//}
		http.Redirect(w, r, "/main", http.StatusFound)
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

func (s *server) pageshipmentBySAP() http.HandlerFunc {
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

		Admin := false
		Stockkeeper := false
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

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		if u.Role == "Administrator" {
			Admin = true
			LoggedIn = true

		} else if u.Role == "кладовщик" {
			Stockkeeper = true
			LoggedIn = true

		}
		data := map[string]interface{}{
			"TitleDOC":    "Отгрузка изделий",
			"User":        u.LastName,
			"Username":    u.FirstName,
			"Admin":       Admin,
			"Stockkeeper": Stockkeeper,
			"LoggedIn":    LoggedIn,
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
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertsapbyship6.html")
	///if err != nil {
	///	panic(err)
	///}
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

func (s *server) pageshowShipmentBySAPBySearchStatic() http.HandlerFunc {
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysapbysearchstatic.html")
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdatebysapbysearchstatic.html")
	///if err != nil {
	///	panic(err)
	///}
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/showdatebysapbysearchstatic.html")
		//	fmt.Fprintf(w, body)
		Admin := false
		Stockkeeper := false
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

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		if u.Role == "Administrator" {
			Admin = true
			LoggedIn = true

		} else if u.Role == "кладовщик" {
			Stockkeeper = true
			LoggedIn = true

		}
		data := map[string]interface{}{
			"TitleDOC":    "Отгрузка изделий",
			"User":        u.LastName,
			"Username":    u.FirstName,
			"Admin":       Admin,
			"Stockkeeper": Stockkeeper,
			"LoggedIn":    LoggedIn,
		}
		tpl.ExecuteTemplate(w, "showdatebysapbysearchstatic.html", data)
		///	if err = tpl.ExecuteTemplate(w, "layout", data); err != nil {
		///		s.error(w, r, http.StatusUnprocessableEntity, err)
		///		return
		///	}

	}
}

func (s *server) showShipmentBySAPBySearchStatic() http.HandlerFunc {
	type searchBy struct {
		LastName string `json:"lastname"`
		Date1    string `json:"date1"`
		Date2    string `json:"date2"`
		Material int    `json:"material"`
	}

	Admin := false
	Stockkeeper := false
	LoggedIn := false
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/showdatebysap2.html")
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showdatebysap2.html")
	///if err != nil {
	///	panic(err)
	///}

	return func(w http.ResponseWriter, r *http.Request) {
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
		if u.Role == "Administrator" {
			Admin = true
			LoggedIn = true

		} else if u.Role == "кладовщик" {
			Stockkeeper = true
			LoggedIn = true

		}

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

		data := map[string]interface{}{
			"TitleDOC":    "Отгрузка изделий",
			"User":        u.LastName,
			"Username":    u.FirstName,
			"Admin":       Admin,
			"Stockkeeper": Stockkeeper,
			"LoggedIn":    LoggedIn,
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

func (s *server) pageidReturn() http.HandlerFunc {
	//tpl, err := template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertidreturn2.html")
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertidreturn2.html")
	///if err != nil {
	///	panic(err)
	///}
	//	tpl = template.Must(template.New("").Delims("<<", ">>").ParseFiles("web/templates/insertidreturn2.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/insertidreturn2.html")
		//	fmt.Fprintf(w, body)
		Admin := false
		Stockkeeper := false
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

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		if u.Role == "Administrator" {
			Admin = true
			LoggedIn = true

		} else if u.Role == "кладовщик" {
			Stockkeeper = true
			LoggedIn = true

		}
		data := map[string]interface{}{
			"TitleDOC":    "Проверка катушек",
			"User":        u.LastName,
			"Username":    u.FirstName,
			"Admin":       Admin,
			"Stockkeeper": Stockkeeper,
			"LoggedIn":    LoggedIn,
		}
		tpl.ExecuteTemplate(w, "insertidreturn2.html", data)
		///	err = tpl.ExecuteTemplate(w, "layout", nil)
		///	if err != nil {
		///		http.Error(w, err.Error(), 400)
		///		return
		///	}
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
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "insertidreturn2.html")
	///if err != nil {
	///	panic(err)
	///}
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
		tpl.ExecuteTemplate(w, "insertidreturn2.html", nil)
		///	err = tpl.ExecuteTemplate(w, "layout", nil)
		//	tpl.ExecuteTemplate(w, "layout", data)
		///	if err != nil {
		///		http.Error(w, err.Error(), 400)
		///		return
		///	}

	}
}

func (s *server) pageInspection() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showinspection.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		//	Admin := true
		Warehouse := true
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
		} else if user.Groups == "склад" {
			StockkeeperWH = true
			Warehouse = false
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
			"Inspector":            Inspector,
			"IngenerQuality":       IngenerQuality,
			"WarehouseManager":     WarehouseManager,
			"Warehouse":            Warehouse,
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

/*
func (s *server) getStatic() http.HandlerFunc {

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
func pieRadius(s *server) *charts.Pie {
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
func (s *server) pageinsertVendor() http.HandlerFunc {
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
		//	GET := map[string]bool{
		//		"admin": admin,
		//		//	"stockkeeper":         stockkeeper,
		//		"главный инженер по качеству": superIngenerQuality,
		//	"stockkeeperWH":       stockkeeperWH,
		//	"inspector":           inspector,
		//	}
		if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true
		} else if user.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true
			fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
		}
		data := map[string]interface{}{
			"Admin":               Admin,
			"SuperIngenerQuality": SuperIngenerQuality,
			//"GET":      GET,
			"LoggedIn": LoggedIn,
			"User":     user.LastName,
			"Username": user.FirstName,
		}
		fmt.Println("Check -")
		tpl.ExecuteTemplate(w, "insertvendor.html", data)
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
	///tpl = template.Must(template.New("base").ParseFiles(s.html+"layout1.html", s.html+"insertvendor1.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Check2 -")
		Admin := false
		SuperIngenerQuality := false
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

		//	GET := map[string]bool{
		//		"admin": admin,
		//	"stockkeeper":         stockkeeper,
		//		"главный инженер по качеству": superIngenerQuality,
		//	"stockkeeperWH":       stockkeeperWH,
		//	"inspector":           inspector,
		//	}

		if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true
		} else if user.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true
			fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
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
			"Admin":               Admin,
			"SuperIngenerQuality": SuperIngenerQuality,
			//"GET":      GET,
			"LoggedIn": LoggedIn,
			"User":     user.LastName,
			"Username": user.FirstName,
		}

		//err = tpl.ExecuteTemplate(w, "layout", data)
		//	err = tpl.ExecuteTemplate(w, "layout", v)
		tpl.ExecuteTemplate(w, "insertvendor.html", data)

	}

}

func (s *server) pageVendor() http.HandlerFunc {
	///tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "showvendor.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {

		//	var body, _ = helper.LoadFile("./web/templates/insertsapbyship6.html")
		//	fmt.Fprintf(w, body)
		//data := map[string]interface{}{
		//	"user": "Я тут",
		//}
		Admin := false
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
		} else if user.Role == "главный инженер по качеству" {
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
			"GET":                 get,
			"LoggedIn":            LoggedIn,
			"User":                user.LastName,
			"Username":            user.FirstName,
		}
		tpl.ExecuteTemplate(w, "showvendor.html", data)

	}
}

func (s *server) pageupdateVendor() http.HandlerFunc {
	///	tpl, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updatevendor.html")
	///	if err != nil {
	///		panic(err)
	///	}
	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
		SuperIngenerQuality := false
		LoggedIn := false

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

		if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true
		} else if user.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true
			fmt.Println("SuperIngenerQuality - ", SuperIngenerQuality)
		}

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
		data := map[string]interface{}{
			"Admin":               Admin,
			"SuperIngenerQuality": SuperIngenerQuality,
			"GET":                 get,
			"LoggedIn":            LoggedIn,
			"User":                user.LastName,
			"Username":            user.FirstName,
		}
		tpl.ExecuteTemplate(w, "updatevendor.html", data)

	}
}

func (s *server) updateVendor() http.HandlerFunc {
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

		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}

		req.ID = id
		req.CodeDebitor = r.FormValue("codedebitor")
		req.NameDebitor = r.FormValue("namedebitor")
		fmt.Println("ID - ", req.ID)
		u := &model.Vendor{
			ID:          req.ID,
			CodeDebitor: req.CodeDebitor,
			NameDebitor: req.NameDebitor,
		}

		if err := s.store.Vendor().UpdateVendor(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
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
	///	_, err := template.New("").Delims("<<", ">>").ParseFiles(s.html + "updatevendor.html")
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
		rdata2 := removeDuplicates(rdata1)
		fmt.Print(rdata2)
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

func (s *server) pageupdateInspection() http.HandlerFunc {
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

func (s *server) updateInspection() http.HandlerFunc {
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

func (s *server) deleteInspection() http.HandlerFunc {
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

//historyInspection
func (s *server) pagehistoryInspection() http.HandlerFunc {

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

func (s *server) historyInspection() http.HandlerFunc {
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

//ListAcceptWHInspection
func (s *server) pageListAcceptWHInspection() http.HandlerFunc { // acceptinspection.html showinspection.html
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

func (s *server) pageacceptWarehouseInspection() http.HandlerFunc {
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

func (s *server) acceptWarehouseInspection() http.HandlerFunc {
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

// showUsersQuality
func (s *server) pageshowUsersQuality() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
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
		} else if user.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true
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
			"Admin":               Admin,
			"SuperIngenerQuality": SuperIngenerQuality,
			"LoggedIn":            LoggedIn,
			"GET":                 get,
		}
		err = tpl.ExecuteTemplate(w, "showUsersQuality.html", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) showUsersQuality() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
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
		} else if user.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true
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
			"Admin":               Admin,
			"SuperIngenerQuality": SuperIngenerQuality,
			"LoggedIn":            LoggedIn,
			"GET":                 get,
		}
		err = tpl.ExecuteTemplate(w, "showUsersQuality.html", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) createUserQuality() http.HandlerFunc {
	type requestFrom struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
		Tabel     string `json:"tabel"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var hdata []requestFrom
		json.Unmarshal(body, &hdata)
		fmt.Printf("body json: %s", body)
		fmt.Println("\njson  struct hdata", hdata)

		Group := "качество"

		for _, v := range hdata {
			fmt.Println(v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)

			u := &model.User{
				Email:     v.Email,
				Password:  v.Password,
				FirstName: v.FirstName,
				LastName:  v.LastName,
				Role:      v.Role,
				Groups:    Group,
				Tabel:     v.Tabel,
			}

			if err := s.store.User().CreateUserByManager(u); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
		}

	}
}

func (s *server) pageupdateUserQuality() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
		SuperIngenerQuality := false
		LoggedIn := false

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

		if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true
		} else if user.Role == "главный инженер по качеству" {
			SuperIngenerQuality = true
			LoggedIn = true
			fmt.Println("SuperIngenerQuality pageupdateUserQuality - ", SuperIngenerQuality)
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		fmt.Println("var id - ", id)
		if err != nil {
			log.Println(err)
		}
		get, err := s.store.User().EditUserByManager(id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		data := map[string]interface{}{
			"GET":                 get,
			"Admin":               Admin,
			"SuperIngenerQuality": SuperIngenerQuality,
			"LoggedIn":            LoggedIn,
			"User":                user.LastName,
			"Username":            user.FirstName,
		}

		//	fmt.Println("Get.email" - get.id)
		tpl.ExecuteTemplate(w, "updateuserquality.html", data)
	}
}

func (s *server) updateUserQuality() http.HandlerFunc {
	type request struct {
		ID        int    `json:"ID"`
		Email     string `json:"email"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Role      string `json:"role"`
		Tabel     string `json:"tabel"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}
		req.ID = id
		req.Email = r.FormValue("email")
		req.Firstname = r.FormValue("firstname")
		req.Lastname = r.FormValue("lastname")
		req.Role = r.FormValue("role")
		//fmt.Println("Роль - ", req.Role)
		req.Tabel = r.FormValue("tabel")
		fmt.Println("ID - ", req.ID)
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
		http.Redirect(w, r, "/showusersquality", 303)
	}

}

func (s *server) deleteUserQuality() http.HandlerFunc {
	type request struct {
		ID int `json:"ID"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}
		req.ID = id

		u := &model.User{
			ID: req.ID,
		}

		if err := s.store.User().DeleteUserByManager(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		http.Redirect(w, r, "/showusersquality", 303)
	}
}

// showUsersQuality
func (s *server) pageshowUsersWarehouse() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
		WarehouseManager := false
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
		} else if user.Role == "старший кладовщик склада" {
			WarehouseManager = true
			LoggedIn = true
		}

		get, err := s.store.User().ListUsersWarehouse()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		data := map[string]interface{}{
			"TitleDOC":         "Сотрудники склада",
			"User":             user.LastName,
			"Username":         user.FirstName,
			"Admin":            Admin,
			"WarehouseManager": WarehouseManager,
			"LoggedIn":         LoggedIn,
			"GET":              get,
		}
		err = tpl.ExecuteTemplate(w, "showUsersWarehouse.html", data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
}

func (s *server) createUserWarehouse() http.HandlerFunc {
	type requestFrom struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Role      string `json:"role"`
		Tabel     string `json:"tabel"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var hdata []requestFrom
		json.Unmarshal(body, &hdata)
		fmt.Printf("body json: %s", body)
		fmt.Println("\njson  struct hdata", hdata)

		group := "склад"

		for _, v := range hdata {
			fmt.Println(v.Email, v.FirstName, v.LastName, v.Password, v.Role, v.Tabel)

			u := &model.User{
				Email:     v.Email,
				Password:  v.Password,
				FirstName: v.FirstName,
				LastName:  v.LastName,
				Role:      v.Role,
				Groups:    group,
				Tabel:     v.Tabel,
			}

			if err := s.store.User().CreateUserByManager(u); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
		}

	}
}

func (s *server) pageupdateUserWarehouse() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		Admin := false
		WarehouseManager := false
		LoggedIn := false

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

		if user.Role == "Administrator" {
			Admin = true
			LoggedIn = true
		} else if user.Role == "старший кладовщик склада" {
			WarehouseManager = true
			LoggedIn = true
			fmt.Println("SuperIngenerQuality pageupdateUserQuality - ", WarehouseManager)
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		fmt.Println("var id - ", id)
		if err != nil {
			log.Println(err)
		}
		get, err := s.store.User().EditUserByManager(id)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		data := map[string]interface{}{
			"GET":              get,
			"Admin":            Admin,
			"WarehouseManager": WarehouseManager,
			"LoggedIn":         LoggedIn,
			"User":             user.LastName,
			"Username":         user.FirstName,
		}

		//	fmt.Println("Get.email" - get.id)
		tpl.ExecuteTemplate(w, "updateuserwarehouse.html", data)
	}
}

func (s *server) updateUserWarehouse() http.HandlerFunc {
	type request struct {
		ID        int    `json:"ID"`
		Email     string `json:"email"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Role      string `json:"role"`
		Tabel     string `json:"tabel"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}
		req.ID = id
		req.Email = r.FormValue("email")
		req.Firstname = r.FormValue("firstname")
		req.Lastname = r.FormValue("lastname")
		req.Role = r.FormValue("role")
		//fmt.Println("Роль - ", req.Role)
		req.Tabel = r.FormValue("tabel")
		fmt.Println("ID - ", req.ID)
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
		http.Redirect(w, r, "/showuserswarehouse", 303)
	}

}

func (s *server) deleteUserWarehouse() http.HandlerFunc {
	type request struct {
		ID int `json:"ID"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["ID"])
		if err != nil {
			log.Println(err)
		}
		req.ID = id

		u := &model.User{
			ID: req.ID,
		}

		if err := s.store.User().DeleteUserByManager(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		http.Redirect(w, r, "/showuserswarehouse", 303)
	}
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
