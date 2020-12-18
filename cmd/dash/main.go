package main

import (
	"context"
	"ctf_dashboard/internal/common"
	"ctf_dashboard/internal/web"
	"flag"
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
	pflag.String("log_level", "INFO", "Log level {INFO|DEBUG|WARNING|ERROR}")
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
	ll := viper.GetString("log_level")
	switch ll {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARNING":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		viper.Set("debug", true)
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.Errorf("Invalid log level provided: %s", ll)
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func setWebServerMode() {
	level := logrus.StandardLogger().GetLevel()
	if level == logrus.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func setConfigDefaults() {
	viper.SetDefault("web.static_dir", "front/dist")
	viper.SetDefault("web.listen", "0.0.0.0:8000")
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
