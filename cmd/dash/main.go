package main

import (
	"context"
	"ctf_dashboard/internal/common"
	"ctf_dashboard/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	setupConfig()
	initLogger()
	setLogLevel()
	setWebServerMode()
	setConfigDefaults()
	parseConfig()

	dashCfg := parseDashboardConfig()

	s := web.NewServer(dashCfg)
	httpServer := startHttpServer(s)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	shutdownHttpServer(httpServer, ctx)

	logrus.Info("Shutdown successful")
}

func setupConfig() {
	pflag.BoolP("verbose", "v", false, "Enable verbose logging")
	pflag.StringP("config", "c", "config.yml", "Path to config file")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		logrus.Fatalf("Error binding flags: %v", err)
	}

	viper.SetEnvPrefix("DASH")
	viper.AutomaticEnv()
}

func initLogger() {
	mainFormatter := &logrus.TextFormatter{}
	mainFormatter.FullTimestamp = true
	mainFormatter.ForceColors = true
	mainFormatter.PadLevelText = true
	mainFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(mainFormatter)
}

func setLogLevel() {
	if viper.GetBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func setWebServerMode() {
	if viper.GetBool("verbose") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func setConfigDefaults() {
	viper.SetDefault("web.static_dir", "front/dist")
	viper.SetDefault("web.listen", "0.0.0.0:8000")
	viper.SetDefault("start_sploit", "resources/start_sploit.py")
	viper.SetDefault("neo.runner_path", "resources/run_neo.sh")
	viper.SetDefault("neo.version", "latest")
	viper.SetDefault("key_file", "resources/ssh_key")
}

func parseConfig() {
	viper.SetConfigFile(viper.GetString("config"))
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal("Error reading config from yaml: ", err)
	}
}

func parseDashboardConfig() common.Config {
	cfg := common.Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.Fatal("Error parsing proxy config from file: ", err)
	}
	logrus.Debugf("Parsed dashboard config: %+v", cfg)
	return cfg
}

func startHttpServer(s *web.Server) *http.Server {
	srv := &http.Server{
		Addr:         viper.GetString("web.listen"),
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 30,
		Handler:      s,
	}

	go func() {
		logrus.Infof("Serving on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatal("Error serving api server: ", err)
		}
	}()

	return srv
}

func shutdownHttpServer(srv *http.Server, ctx context.Context) {
	logrus.Info("Shutting down http server")
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Error shutting down http server: %v", err)
	}
}
