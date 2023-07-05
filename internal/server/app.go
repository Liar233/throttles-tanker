package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Liar233/throttles-tank/internal/server/action"
	"github.com/Liar233/throttles-tank/internal/server/command"
	"github.com/Liar233/throttles-tank/internal/server/service"
	"github.com/Liar233/throttles-tank/internal/server/storage"
	"github.com/Liar233/throttles-tank/pkg/protocol"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type AppServerConfig struct {
	HttpConfig     HttpServerConfig              `mapstructure:"http"`
	TankerConfig   service.TankerConfig          `mapstructure:"tanker"`
	StorageConfig  storage.PostgresStorageConfig `mapstructure:"storage"`
	UpdatesTimeout time.Duration                 `mapstructure:"update_timeout"`
}

type AppServer struct {
	config      AppServerConfig
	http        HttpServerAdapter
	ruleStorage storage.RuleStorageInterface
	tanker      service.TankerInterface
	whiteList   service.ListManagerInterface
	blackList   service.ListManagerInterface
	ticker      *time.Ticker
}

func (app *AppServer) Bootstrap() {

	initLog()

	app.ruleStorage = storage.NewPostgresRuleStorage(app.config.StorageConfig)

	app.ticker = time.NewTicker(app.config.UpdatesTimeout)

	app.whiteList = service.NewListManager()
	app.blackList = service.NewListManager()

	app.tanker = service.NewTanker(app.config.TankerConfig, app.ticker.C)

	app.http = NewHttpServerAdapter(app.config.HttpConfig)

	router := app.BootstrapActions()

	app.http.SetHandler(router)
}

func (app *AppServer) BootstrapActions() http.Handler {

	addCmd := command.NewAddRuleCommand(app.ruleStorage, app.whiteList, app.blackList)
	deleteCmd := command.NewDeleteRuleCommand(app.ruleStorage, app.whiteList, app.blackList)
	checkCmd := command.NewCheckCommand(app.whiteList, app.blackList, app.tanker)
	flushCmd := command.NewFlushCommand(app.tanker)

	router := mux.NewRouter()

	router.Handle(protocol.AddWhiteRuleUrl, action.NewAddRuleAction("white", addCmd)).
		Methods(http.MethodPost)
	router.Handle(protocol.AddBlackRuleUrl, action.NewAddRuleAction("black", addCmd)).
		Methods(http.MethodPost)
	router.Handle(protocol.DeleteWhiteRuleUrl, action.NewDeleteRuleAction("white", deleteCmd)).
		Methods(http.MethodDelete)
	router.Handle(protocol.DeleteBlackRuleUrl, action.NewDeleteRuleAction("black", deleteCmd)).
		Methods(http.MethodDelete)
	router.Handle(protocol.CheckUrl, action.NewCheckAction(checkCmd)).
		Methods(http.MethodPost)
	router.Handle(protocol.FlushUrl, action.NewFlushAction(flushCmd)).
		Methods(http.MethodPost)

	return router
}

func (app *AppServer) Run() (err error) {

	log.Infoln("Running server...")

	if errStorage := app.ruleStorage.Connect(); errStorage != nil {

		return errStorage
	}

	log.Infoln("Storage ok.")

	go func() {

		if errHttp := app.http.ListenAndServe(); errHttp != nil {

			err = errHttp
			return
		}
	}()

	log.Infoln("Http-server ok.")

	return app.Stop()
}

func (app *AppServer) Stop() error {

	sigintChan := make(chan os.Signal, 1)

	defer close(sigintChan)

	signal.Notify(sigintChan, syscall.SIGINT, syscall.SIGTERM)

	s := <-sigintChan

	log.Infof("Stopping server. Got signal %s", s)

	app.ticker.Stop()

	if err := app.http.Close(); err != nil {

		return err
	}

	log.Infoln("Http-server off.")

	if err := app.ruleStorage.Close(); err != nil {

		return err
	}

	log.Infoln("Storage off.")

	log.Infoln("Server stopped.")

	return nil
}

func NewAppServer(config AppServerConfig) *AppServer {

	return &AppServer{
		config: config,
	}
}

func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}
