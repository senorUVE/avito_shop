package main

import (
	"log/slog"
	"os"

	authhandler "auth/internal/api/authentication"
	infohandler "auth/internal/api/info"
	transhandler "auth/internal/api/transaction"
	"auth/internal/config"
	"auth/internal/middleware"

	"auth/internal/repository"
	authsrv "auth/internal/services/authentication"
	infosrv "auth/internal/services/info"
	passwordsrv "auth/internal/services/password"
	tokenizersrv "auth/internal/services/tokenizer"
	transsrv "auth/internal/services/transaction"

	pkgconfig "github.com/WantBeASleep/goooool/config"
	"github.com/WantBeASleep/goooool/loglib"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	successExitCode = 0
	failExitCode    = 1
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	loglib.InitLogger(loglib.WithDevEnv())
	cfg, err := pkgconfig.Load[config.Config]()
	if err != nil {
		slog.Error("init config", "err", err)
		return failExitCode
	}

	pubKey, privKey, err := cfg.ParseRsaKeys()
	if err != nil {
		slog.Error("parse rsa keys", "err", err)
		return failExitCode
	}
	db, err := sqlx.Open("postgres", cfg.DB.Dsn)
	if err != nil {
		slog.Error("init db", "err", err)
		return failExitCode
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("ping db", "err", err)
		return failExitCode
	}

	dao := repository.NewRepository(db)
	tokenizerSrv := tokenizersrv.New(
		cfg.JWT.AccessTokenTime,
		cfg.JWT.RefreshTokenTime,
		privKey,
		pubKey,
	)
	passwordSrv := passwordsrv.New()
	authSrv := authsrv.New(dao, tokenizerSrv, passwordSrv)
	infoSrv := infosrv.New(dao)
	transSrv := transsrv.New(dao)
	authHandler := authhandler.New(authSrv)
	infoHandler := infohandler.New(infoSrv)
	transHandler := transhandler.New(transSrv)
	router := gin.Default()
	router.POST("/api/auth", authHandler.Authenticate)

	protected := router.Group("/api")
	protected.Use(middleware.CheckJWT(tokenizerSrv))
	{
		protected.GET("/info", infoHandler.GetUserInfo)
		protected.GET("/buy/:item", transHandler.BuyItem)
		protected.POST("/sendCoins", transHandler.TransferCoins)

	}

	slog.Info("Starting HTTP server", slog.String("url", cfg.App.Url))
	if err := router.Run(cfg.App.Url); err != nil {
		slog.Error("Failed to start HTTP server", "err", err)
		return failExitCode
	}

	return successExitCode
}
