package clevergo

import (
	"github.com/clevergo/cache"
	"github.com/clevergo/ini"
	"github.com/clevergo/log"
	"github.com/clevergo/session"
	"strings"
)

const (
	ModeDev = 1
	ModePro = 2
)

// Configuration of application.
type Config struct {
	goPath  string
	srcPath string

	mode int

	// Server Configuration
	serverHost     string
	serverProtocol string
	serverCertFile string
	serverKeyFile  string

	// Controller Configuration
	controllerPrefix string
	controllerSuffix string

	// Action Configuration
	actionPrefix string
	actionSuffix string

	// View Configuration
	viewSuffix string

	// Session Configuration
	enableSession bool
	sessionStore  session.Store
	sessionName   string
	sessionMaxAge int

	// Log Configuration
	enableLog bool
	logger    *log.Logger
	logLevel  int
	logFlag   int
	// FileTartget
	logFileLevel    int
	logFileDir      string
	logFileName     string
	logFilePath     string
	logFileMaxSize  int64
	logFileInterval int
	// EmailTarget
	logMailLevel    int
	logMailHost     string
	logMailPort     string
	logMailUser     string
	logMailPassword string
	logMailFrom     string
	logMailTo       string
	logMailSubject  string

	// Router Configuration
	routeSuffix                  string
	routerRedirectTrailingSlash  bool
	routerRedirectFixedPath      bool
	routerHandleMethodNotAllowed bool
	routerHandleOPTIONS          bool

	// Redis Configuration
	enableCache      bool
	cache            *cache.RedisCache
	redisNetwork     string
	redisAddress     string
	redisPassword    string
	redisDb          string
	redisMaxIdle     int
	redisIdleTimeout int
}

func (c *Config) Load(filename string) {
	iniConfig := ini.NewConfig(filename)

	section, err := iniConfig.GetSection()
	if err != nil {
		panic(err)
	}

	// Get mode.
	mode, err := section.GetString("mode")
	if (err == nil) && (strings.EqualFold("PRO", mode)) {
		c.mode = ModePro
	} else {
		c.mode = ModeDev
	}

	// Get controller configuration.
	controllerPrefix, err := section.GetString("controller.prefix")
	if err == nil {
		c.controllerPrefix = controllerPrefix
	}
	controllerSuffix, err := section.GetString("controller.suffix")
	if err == nil {
		c.controllerSuffix = controllerSuffix
	}

	// Get action configuration.
	actionPrefix, err := section.GetString("action.prefix")
	if err == nil {
		c.actionPrefix = actionPrefix
	}
	actionSuffix, err := section.GetString("action.suffix")
	if err == nil {
		c.actionSuffix = actionSuffix
	}

	// Get session configuration
	enableSession, err := section.GetBool("session.enable")
	if err == nil {
		c.enableSession = enableSession
	}
	sessionName, err := section.GetString("session.name")
	if err == nil {
		c.sessionName = sessionName
	}
	sessionMaxAge, err := section.GetInt("session.max_age")
	if (err == nil) && (sessionMaxAge > 0) {
		c.sessionMaxAge = sessionMaxAge
	}

	// Get Redis Cache configuration
	redisMaxIdle, err := section.GetInt("redis.max_idle")
	if err == nil {
		c.redisMaxIdle = redisMaxIdle
	}
	redisIdleTimeout, err := section.GetInt("redis.idle_timeout")
	if (err == nil) && (redisIdleTimeout > 0) {
		c.redisIdleTimeout = redisIdleTimeout
	}
	redisNetwork, err := section.GetString("redis.network")
	if err == nil {
		c.redisNetwork = redisNetwork
	}
	redisAddress, err := section.GetString("redis.address")
	if err == nil {
		c.redisAddress = redisAddress
	}
	redisPassword, err := section.GetString("redis.password")
	if err == nil {
		c.redisPassword = redisPassword
	}
	redisDb, err := section.GetString("redis.db")
	if err == nil {
		c.redisDb = redisDb
	}

	// Get log configuration
	enableLog, err := section.GetBool("log.enable")
	if err == nil {
		c.enableLog = enableLog
	}
	logLevel, err := section.GetInt("log.level")
	if err == nil {
		c.logLevel = logLevel
	}
	logFlag, err := section.GetInt("log.flag")
	if err == nil {
		c.logFlag = logFlag
	}
	logFileDir, err := section.GetString("log.file_dir")
	if err == nil {
		c.logFileDir = logFileDir
	}
	logFileName, err := section.GetString("log.file_name")
	if err == nil {
		c.logFileName = logFileName
	}
	logFilePath, err := section.GetString("log.file_path")
	if err == nil {
		c.logFilePath = logFilePath
	}
	logFileMaxSize, err := section.GetInt("log.file_max_size")
	if err == nil {
		c.logFileMaxSize = int64(logFileMaxSize)
	}
	logFileInterval, err := section.GetInt("log.file_interval")
	if err == nil {
		c.logFileInterval = logFileInterval
	}
	logFileLevel, err := section.GetInt("log.file_level")
	if err == nil {
		c.logLevel = logFileLevel
	}
	logMailLevel, err := section.GetInt("log.mail_level")
	if err == nil {
		c.logMailLevel = logMailLevel
	}
	logMailHost, err := section.GetString("log.mail_host")
	if err == nil {
		c.logMailHost = logMailHost
	}
	logMailPort, err := section.GetString("log.mail_port")
	if err == nil {
		c.logMailPort = logMailPort
	}
	logMailUser, err := section.GetString("log.mail_user")
	if err == nil {
		c.logMailUser = logMailUser
	}
	logMailPassword, err := section.GetString("log.mail_password")
	if err == nil {
		c.logMailPassword = logMailPassword
	}
	logMailFrom, err := section.GetString("log.mail_from")
	if err == nil {
		c.logMailFrom = logMailFrom
	}
	logMailTo, err := section.GetString("log.mail_to")
	if err == nil {
		c.logMailTo = logMailTo
	}
	logMailSubject, err := section.GetString("log.mail_subject")
	if err == nil {
		c.logMailSubject = logMailSubject
	}

	// Get router configuration
	routerRedirectTrailingSlash, err := section.GetBool("router.redirect_trailing_slash")
	if err == nil {
		c.routerRedirectTrailingSlash = routerRedirectTrailingSlash
	}
	routerRedirectFixedPath, err := section.GetBool("router.redirect_fixed_path")
	if err == nil {
		c.routerRedirectFixedPath = routerRedirectFixedPath
	}
	routerHandleMethodNotAllowed, err := section.GetBool("router.handle_method_not_allowed")
	if err == nil {
		c.routerHandleMethodNotAllowed = routerHandleMethodNotAllowed
	}
	routerHandleOPTIONS, err := section.GetBool("router.handle_options")
	if err == nil {
		c.routerHandleOPTIONS = routerHandleOPTIONS
	}
}

func (c *Config) GoPath() string {
	return c.goPath
}

func (c *Config) SrcPath() string {
	return c.srcPath
}

func (c *Config) ServerHost() string {
	return c.serverHost
}

func (c *Config) ServerProtocol() string {
	return c.serverProtocol
}

func (c *Config) ServerCertFile() string {
	return c.serverCertFile
}

func (c *Config) ServerKeyFile() string {
	return c.serverKeyFile
}

func (c *Config) ControllerPrefix() string {
	return c.controllerPrefix
}

func (c *Config) ControllerSuffix() string {
	return c.controllerSuffix
}

func (c *Config) ActionPrefix() string {
	return c.actionPrefix
}

func (c *Config) ActionSuffix() string {
	return c.actionSuffix
}

func (c *Config) ViewSuffix() string {
	return c.viewSuffix
}

func (c *Config) EnableSession() bool {
	return c.enableSession
}

func (c *Config) SessionName() string {
	return c.sessionName
}

func (c *Config) EnableLog() bool {
	return c.enableLog
}

func (c *Config) LogFileDir() string {
	return c.logFileDir
}

func (c *Config) LogFilePath() string {
	return c.logFilePath
}

func (c *Config) LogFileName() string {
	return c.logFileName
}

func (c *Config) LogLevel() int {
	return c.logLevel
}

func (c *Config) LogFlag() int {
	return c.logFlag
}

func (c *Config) LogFileMaxSize() int64 {
	return c.logFileMaxSize
}

func (c *Config) LogFileInterval() int {
	return c.logFileInterval
}
