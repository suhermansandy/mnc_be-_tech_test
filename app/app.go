package app

import (
	"mnc-be-tech-test/app/db"
	"mnc-be-tech-test/app/excel"
	"mnc-be-tech-test/app/handler"
	"mnc-be-tech-test/app/model"
	"mnc-be-tech-test/config"
	"mnc-be-tech-test/logger"
	"net/http"
	"time"

	"github.com/go-stomp/stomp"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// App has router and db instances
type App struct {
	Router  *mux.Router
	DB      map[string]db.DBHandler
	Artemis *model.Artemis
	// ArtemisQueuePath string
	// ArtemisNew       *stomp.Conn
	// DBNew            db.DBHandler
}

// log file to use
var log = logger.LogType.DefaultLog

// Initialize initializes the app with predefined configuration
func (a *App) Initialize() {
	log.Println("API Init")
	startTime := time.Now()
	a.DB = make(map[string]db.DBHandler)
	dbHandler, err := db.OpenDB("postgres", config.Env.DbConn)
	if err != nil {
		log.Println("Could not connect database")
		log.Fatal(err.Error())
	}
	a.DB["db"] = db.AutoMigrate(dbHandler)
	a.DB["db"].LogMode(true)

	mockDbHandler, errMock := db.OpenDB("postgres", config.Env.DbConnMock)
	if errMock != nil {
		log.Println("Could not connect database")
		log.Fatal(errMock.Error())
	}
	a.DB["mock"] = db.AutoMigrate(mockDbHandler)

	a.Artemis = &model.Artemis{
		Skip: config.Env.ArtemisSkip,
	}
	a.Artemis.Conn, err = stomp.Dial("tcp", config.Env.ArtemisConn, config.Env.ArtemisOptions...)
	if err != nil {
		log.Println("Could not connect artemis")
		log.Println(err.Error())
	}

	// set router
	a.Router = mux.NewRouter()
	a.setRouters()

	log.Println("Init finished in ", time.Now().Sub(startTime))
	log.Println("API Ready")
}

func (a *App) setRouter(urlName string, h handler.RESTHandler) {
	a.Get("/"+urlName, a.ihandleRequest(h, handler.GetList))
	a.Get("/"+urlName+"/excel", a.ihandleRequest(h, handler.CreateExcel))
	a.Get("/"+urlName+"/{id}", a.ihandleRequest(h, handler.GetByID))
	a.Post("/"+urlName, a.ihandleRequest(h, handler.Create))
	a.Put("/"+urlName+"/{id}", a.ihandleRequest(h, handler.Update))
	a.Delete("/"+urlName+"/{id}", a.ihandleRequest(h, handler.Delete))
}

func (a *App) setRouters() {
	userHandler := handler.UserHandler{RESTHandler: handler.New(a.DB, excel.New(), a.Artemis)}
	a.setRouter("user", userHandler)
	a.Post("/login", a.ihandleRequest(userHandler, handler.LoginUser))
	a.Post("/logout", a.ihandleRequest(userHandler, handler.LogOutUser))

	paymentHandler := handler.PaymentHandler{RESTHandler: handler.New(a.DB, excel.New(), a.Artemis)}
	a.setRouter("payment", paymentHandler)

	hostoryHandler := handler.HistoryHandler{RESTHandler: handler.New(a.DB, excel.New(), a.Artemis)}
	a.setRouter("history", hostoryHandler)
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *App) Run(host string) {
	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Accept", "Content-Length", "Access-Control-Allow-Headers", "Accept-Encoding", "X-Requested-With", "Content-Type", "X-CSRF-Token", "Authorization", "Bearer", "Origin", "Accept"},
		ExposedHeaders:     []string{"File-Name"},
		OptionsPassthrough: false,
	})
	log.Fatal(http.ListenAndServe(host, c.Handler(a.Router)))
}

// RequestIHandlerFunction is abbreviation for handler data, response write, request
type RequestIHandlerFunction func(h handler.RESTHandler, w http.ResponseWriter, r *http.Request)

func (a *App) ihandleRequest(h handler.RESTHandler, handler RequestIHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(h, w, r)
	}
}
