package main

import (
	"flag"
	"os"
	"strings"

	"github.com/joho/godotenv"
	logrus "github.com/sirupsen/logrus"
	"github.com/star-39/moe-sticker-bot/core"
)

// Common abbr. in this project:
// S : Sticker
// SS : StickerSet

func main() {
	// ✅ Cargar variables desde .env
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("No se pudo cargar .env, usando variables del sistema")
	}

	conf := parseCmdLine()

	// ✅ Inicializar el bot con la configuración
	core.Init(conf)
}

func parseCmdLine() core.ConfigTemplate {
	var help = flag.Bool("help", false, "Show help")
	var adminUid = flag.Int64("admin_uid", -1, "Admin's UID(optional)")
	var botToken = flag.String("bot_token", "", "Telegram Bot Token")
	var dataDir = flag.String("data_dir", "", "Overwrites the working directory where msb puts data.")
	var webappUrl = flag.String("webapp_url", "", "Public HTTPS URL to WebApp, in unset, webapp will be disabled.")
	var WebappApiListenAddr = flag.String("webapp_listen_addr", "", "Webapp API server listen address(IP:PORT)")
	var webappDataDir = flag.String("webapp_data_dir", "", "Where to put webapp data to share with ReactApp ")
	var dbAddr = flag.String("db_addr", "", "mariadb(mysql) address, if unset, database will be disabled.")
	var dbUser = flag.String("db_user", "", "mariadb(mysql) username")
	var dbPass = flag.String("db_pass", "", "mariadb(mysql) password")
	var logLevel = flag.String("log_level", "debug", "Log level")
	flag.Parse()

	if *help {
		flag.Usage()
		println("Only --bot_token is required to run.")
		os.Exit(0)
	}

	conf := core.ConfigTemplate{}

	// ✅ Leer BOT_TOKEN desde .env si no se pasa por argumento
	conf.BotToken = *botToken
	if conf.BotToken == "" {
		conf.BotToken = os.Getenv("BOT_TOKEN")
	}
	if conf.BotToken == "" {
		logrus.Error("Please set BOT_TOKEN in .env")
		logrus.Fatal("No bot token provided!")
	}
	if !strings.Contains(conf.BotToken, ":") {
		logrus.Fatal("Bad bot token!")
	}

	// ✅ Leer datos de la DB desde .env
	if *dbAddr == "" {
		conf.DbAddr = os.Getenv("DB_ADDR")
	} else {
		conf.DbAddr = *dbAddr
	}
	if *dbUser == "" {
		conf.DbUser = os.Getenv("DB_USER")
	} else {
		conf.DbUser = *dbUser
	}
	if *dbPass == "" {
		conf.DbPass = os.Getenv("DB_PASS")
	} else {
		conf.DbPass = *dbPass
	}

	conf.WebappUrl = *webappUrl
	conf.WebappDataDir = *webappDataDir
	conf.WebappApiListenAddr = *WebappApiListenAddr

	conf.LogLevel = *logLevel
	conf.AdminUid = *adminUid
	conf.DataDir = *dataDir

	return conf
}
