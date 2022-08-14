package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	config2 "restapi/cmd/iternal/config"
	"restapi/cmd/iternal/user"
	"restapi/cmd/pkg/logging"
	"time"
)

func main() {
	logger := logging.Getlogger()
	logger.Info("я ебу собак")
	router := httprouter.New()

	cfg := config2.GetConfig()

	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config2.Config) {
	logger := logging.Getlogger()
	logger.Info("start app")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("creating socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("create unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Info("listen unix socket")
		logger.Infof("server is listening unix socket :%s", socketPath)
		if err != nil {
			logger.Fatal(err)
		}

	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
		logger.Infof("server listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
		if listenErr != nil {
			panic(listenErr)
		}
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Fatal(server.Serve(listener))
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("fuckshitfuck"))
}

