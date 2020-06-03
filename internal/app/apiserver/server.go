package apiserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	//	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	//	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
	//	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store/helper"
	//	"github.com/google/uuid"
	//	"github.com/gorilla/handlers"
	//	"github.com/gorilla/mux"
	//	"github.com/gorilla/sessions"
	//	_ "github.com/lib/pq"
	//	"github.com/sirupsen/logrus"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/model"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store"
	"github.com/eugenefoxx/http-rest-api-starline/internal/app/store/helper"

	//	"github.com/gobuffalo/packr/v2/jam/store"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
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
	database     *sql.DB
}

func init() {
	//	tpl = template.Must(template.ParseGlob("./web/templates/*.html"))
	tpl = template.Must(template.New("./web/templates/*.html").Delims("<<", ">>").ParseGlob("./web/templates/*.html"))
	//tpl = template.Must(template.ParseFiles("./web/templates/*.html"))
	//	tpl = template.Must(template.ParseGlob("C:/Users/Евгений/templates/*.html"))
	//	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./web/images/"))))

}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(), // mux.NewRouter()  NewRouter()
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
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

	//	s.router.HandleFunc("/sessions", s.redirectMain(s.handleSessionsCreate())).Methods("POST")
	//s.router.HandleFunc("/sessions", s.handleSessionsCreate(s.redirectMain())).Methods("POST")
	s.router.HandleFunc("/", s.pagehandleSessionsCreate()).Methods("GET")
	s.router.HandleFunc("/", s.handleSessionsCreate()).Methods("POST")
	//	s.router.HandleFunc("/sessions", s.redirectMain())
	//	s.router.HandleFunc("/sessions", s.pageredirectMain())

	s.router.HandleFunc("/shipmentbysap", s.authMiddleware(s.pageshipmentBySAP())).Methods("GET")
	s.router.HandleFunc("/shipmentbysap", s.authMiddleware(s.shipmentBySAP())).Methods("POST")

	//	s.router.HandleFunc("/showdateshipmentbysap", s.pageshowShipmentBySAP()) //.Methods("GET")
	s.router.HandleFunc("/showdateshipmentbysap", s.authMiddleware(s.showShipmentBySAP())) //.Methods("POST")

	//	s.router.HandleFunc("/main", s.authMiddleware(s.pageredirectMain())).Methods("GET")
	s.router.HandleFunc("/logout", s.signOut()).Methods("POST")

	s.router.HandleFunc("/hello", s.authMiddleware(s.handleHello()))
	s.router.HandleFunc("/main", s.authMiddleware(s.main())).Methods("GET")
	//	s.router.HandleFunc("/", s.loginPage())
	s.router.HandleFunc("/js", s.jsPage())

	// /private/***
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")

	s.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web"))))

	//	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./web/images"))))

	//	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./web/"))))
	//	http.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir(config.assets))))
	//	http.Handle("/", http.FileServer(http.Dir("./web/images")))
	//s.router.Handle("/resources/", http.StripPrefix("/resources", http.FileServer(http.Dir("./web/"))))
	//	http.Handle("/", s.router)

	open.StartWith("http://localhost:3000/", "firefox") // chromium

}

func (s *server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}

func (s *server) main() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("/main - user:", u.Email, u.ID)
		data := map[string]interface{}{
			"user": u.LastName,
			"id":   u.FirstName,
		}

		tpl.ExecuteTemplate(w, "index.html", data)
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
	return func(w http.ResponseWriter, r *http.Request) {
		var body, _ = helper.LoadFile("./web/templates/register.html")
		fmt.Fprintf(w, body)
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		//	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		//		s.error(w, r, http.StatusBadRequest, err)
		//		return
		//	}

		//	email := r.FormValue("email")
		//	password := r.FormValue("password")
		//	firstname := r.FormValue("firstname")
		//	lastname := r.FormValue("lastname")

		//req.Email = r.FormValue("email")
		//fmt.Println(req.Email)

		//req.Password = r.FormValue("password")

		//req.FirstName = r.FormValue("firstname")

		//req.LastName = r.FormValue("lastname")
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
		tpl.ExecuteTemplate(w, "login.html", nil)
	}
}

func (s *server) pagehandleSessionsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body, _ = helper.LoadFile("./web/templates/login.html") // "./web/templates/login.html"
		fmt.Fprintf(w, body)
	}
}

func (s *server) pageredirectMain() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body, _ = helper.LoadFile("./web/templates/index.html")
		fmt.Fprintf(w, body)
	}
}

// h http.HandlerFunc
func (s *server) redirectMain() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	var body, _ = helper.LoadFile("./web/templates/index.html")
		//	fmt.Fprintf(w, body)
		//s.router.HandleFunc("/main", s.main())
		http.Redirect(w, r, "/main", 303)
	})
}

func (s *server) handleSessionsCreate() http.HandlerFunc {

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		r.ParseForm()
		email := r.FormValue("email")
		password := r.FormValue("password")
		//	target := "/sessions"

		u, err := s.store.User().FindByEmail(email)
		if err != nil || !u.ComparePassword(password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)

			return
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

		fmt.Println("handleSessionsCreate()", u.Email)
		//	s.redirectMain()
		//	s.pageredirectMain()
		//	_, ok := session.Values["user_id"]
		//	if !ok {

		http.Redirect(w, r, "/main", 303)
		//	return
		//	}
		data := map[string]interface{}{
			"user": u.LastName,
			"id":   u.FirstName,
		}
		tpl.ExecuteTemplate(w, "index.html", data) //  "index.html"

	}
}

func (s *server) pageshipmentBySAP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body, _ = helper.LoadFile("./web/templates/insertsapbyship4.html")
		fmt.Fprintf(w, body)
	}
}

func (s *server) shipmentBySAP() http.HandlerFunc {
	type reqA struct {
		Material int
		Qty      int
		Comment  string
		//	ID       int
		//	LastName string
	}

	type request struct {
		Material int    `db:"material"`
		Qty      int    `db:"qty"`
		Comment  string `db:"comment"`
		ID       int    `db:"id"`
		LastName string `db:"lastname"`
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
		fmt.Printf("tect %s", body)
		//	fmt.Printf("1 тест %s", &hdata.Material)
		//	fmt.Println("tect2 %s", hdata, "\n")
		fmt.Println("\nall of the data", hdata)
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
			checkmaterial := strconv.Itoa(v.Material)
			//	checkmaterial := v.Material
			if len(checkmaterial) == 7 {
				u := &model.Shipmentbysap{
					Material: v.Material,    //record.Material, //material,
					Qty:      v.Qty,         //record.Qty,      // qty,
					Comment:  v.Comment,     //record.Comment,  //comment,
					ID:       user.ID,       //record.ID,       //user.ID,
					LastName: user.LastName, //record.LastName, // user.LastName,
				}

				if err := s.store.Shipmentbysap().InterDate(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
			} else {
				fmt.Println("кол-во не равно 7", v.Material)

				tpl.ExecuteTemplate(w, "error.html", nil)
				return
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
		tpl.ExecuteTemplate(w, "insertsapbyship4.html", nil)
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
	return func(w http.ResponseWriter, r *http.Request) {
		var body, _ = helper.LoadFile("./web/templates/showdatebysap.html")
		fmt.Fprintf(w, body)
	}
}

func (s *server) showShipmentBySAP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		get, err := s.store.Shipmentbysap().ShowDate()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		err = tpl.ExecuteTemplate(w, "showdatebysap.html", get)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
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
