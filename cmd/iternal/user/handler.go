package user

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	config2 "restapi/cmd/iternal/config"
	"restapi/cmd/iternal/handlers"
	"restapi/cmd/iternal/user/db"
	"restapi/cmd/pkg/client/mongodb"
	"restapi/cmd/pkg/logging"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

const (
	usersURL = "/users"
	userURL  = "/users:uuid"
)

func (h *handler) Register(router *httprouter.Router) {
	router.GET(usersURL, h.GetList)
	router.GET(userURL, h.GetUserByUUID)
	router.POST(userURL, h.CreateUser)
	router.PUT(userURL, h.UpdateUser)
	router.PATCH(userURL, h.PartiallyUpdateUser)
	router.DELETE(userURL, h.DeleteUser)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("this is list of users"))
}
func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("this is user by uuid"))
}
func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("this is update user"))
}
func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("this is partially update user"))
}
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := logging.Getlogger()
	cfg := config2.GetConfig()
	cfgMongo := cfg.MongoDB
	client, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username, cfgMongo.Password, cfgMongo.Database, cfgMongo.Auth_db)
	if err != nil {
		panic(err)
	}
	user1 := User{
		ID:           "",
		Email:        r.FormValue("email"),
		Username:     r.FormValue("username"),
		PasswordHASH: r.FormValue("password"),
	}
	storage := db.NewStorage(client, cfg.MongoDB.Collection, logger)
	storage.Create(context.Background(), user1)
	logger.Info("пользователь успешно создан")
}
func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("this is delete user"))
}
