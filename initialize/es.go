package initialize

import (
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"os"
	"server/global"
)

func ConnectEs() *elasticsearch.TypedClient {
	esCfg := global.Config.ES
	cfg := elasticsearch.Config{
		Addresses: []string{esCfg.URL},
		Username:  esCfg.Username,
		Password:  esCfg.Password,
	}

	if esCfg.IsConsolePrint {
		cfg.Logger = &elastictransport.ColorLogger{
			Output:             os.Stderr,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		}
	}

	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		global.Log.Error("Failed to create ES client", zap.Error(err))
		os.Exit(1)
	}
	return client
}
