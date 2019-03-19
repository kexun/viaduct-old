package main

import (
	"strings"

	"github.com/kexun/viaduct/cmd/config"
	"github.com/kexun/viaduct/pkg/conn"
	"github.com/kexun/viaduct/pkg/mux"
	"github.com/kubeedge/beehive/pkg/common/log"
)

func main() {
	cfg := config.InitConfig()

	var err error
	if strings.Compare(cfg.CmdType, "server") == 0 {
		err = StartServer(cfg)
	} else if strings.Compare(cfg.CmdType, "client") == 0 {
		err = StartClient(cfg)
	} else {
		panic("wrong cmd type")
	}

	if err != nil {
		log.LOGGER.Errorf("Failed to start &s", cfg.CmdType)
	} else {
		log.LOGGER.Infof("%s is finished.", cfg.CmdType)
	}
}

func ConnNotify(connection conn.Conn) {
	log.LOGGER.Infof("notify connection(%s)", connection)
}

func handle(container *mux.MessageContainer, writer mux.ResponseWriter) {
	log.LOGGER.Infof("receive a message, message container: %+v", container.Message.Content)
}

func initEntries() {
	// todo
}

func StartServer(cfg *config.Config) error {
	// todo

	return nil
}

func StartClient(cfg *config.Config) error {
	// todo

	return nil
}
