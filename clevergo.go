package clevergo

import (
	"crypto"
	"github.com/clevergo/cache"
	"github.com/clevergo/clevergo/utils/string"
	"github.com/clevergo/jwt"
	"github.com/clevergo/log"
	"github.com/clevergo/session"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/smtp"
	"os"
	"path"
	"strings"
)

var (
	apps          Applications
	Configuration *Config
	defaultApp    *Application
	goPath        string
	srcPath       string
)

func init() {
	goPath = os.Getenv("GOPATH")
	if len(goPath) == 0 {
		panic("GOPATH is not set.")
	}
	srcPath = path.Join(goPath, "src")

	Configuration = &Config{
		goPath:  goPath,
		srcPath: srcPath,
		mode:    ModeDev,
		// Server configuration
		serverHost:     ":8080",
		serverProtocol: "HTTP",
		serverCertFile: "",
		serverKeyFile:  "",

		// Controller configuration
		controllerPrefix: "",
		controllerSuffix: "Controller",

		// Action configuration
		actionPrefix: "Action",
		actionSuffix: "",
		actionMethod: "_method",

		// View configuration
		viewSuffix: ".html",

		// JSON WEB TOKEN Configuration
		enableJWT:        true,
		JWT:              nil,
		jwtIssuer:        "CleverGO",
		jwtTTL:           int64(3600 * 24 * 7),
		jwtHMACSecretKey: stringutil.GenerateRandomString(32),
		jwtRSAPrivateKey: "",
		jwtRSAPublicKey:  "",

		// Session configuration
		enableSession: false,
		sessionStore:  nil,
		sessionName:   "GOSESSION",
		sessionMaxAge: 10 * 24 * 3600,

		// Log configuration
		enableLog:       true,
		logger:          nil,
		logFlag:         log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile,
		logLevel:        log.LevelDebug | log.LevelInfo | log.LevelWarn | log.LevelError | log.LevelFatal,
		logFileLevel:    log.LevelInfo | log.LevelWarn | log.LevelError | log.LevelFatal,
		logFileDir:      "logs",
		logFilePath:     path.Join(goPath, "logs"),
		logFileName:     "app.log",
		logFileMaxSize:  int64(20 * 1024 * 1024),
		logFileInterval: 3600,
		logMailLevel:    log.LevelError | log.LevelFatal,
		logMailHost:     "",
		logMailPort:     "",
		logMailUser:     "",
		logMailPassword: "",
		logMailFrom:     "",
		logMailTo:       "",
		logMailSubject:  "Application Log",

		// Route configuration
		routerRedirectTrailingSlash:  true,
		routerRedirectFixedPath:      true,
		routerHandleMethodNotAllowed: true,
		routerHandleOPTIONS:          true,

		// Cache configuration
		enableCache:      true,
		cache:            nil,
		redisNetwork:     "tcp",
		redisAddress:     ":6379",
		redisPassword:    "",
		redisDb:          "0",
		redisMaxIdle:     1000,
		redisIdleTimeout: 300,
	}

	apps = make(Applications, 0)
}

func LoadConfig(filename string) {
	Configuration.Load(filename)
}

func Init() {
	// Check configuration.
	checkConfiguration()

	// Initialize logger.
	if Configuration.enableLog {
		Configuration.logger = log.NewLogger(
			Configuration.logLevel,
			Configuration.logFlag,
		)

		// Add FileTarget
		logFile, err := log.OpenFile(path.Join(Configuration.logFilePath, Configuration.logFileName))
		if err != nil {
			panic(err.Error())
		}

		if len(Configuration.logFilePath) > 0 {
			fileTarget := log.NewFileTarget(Configuration.logger, Configuration.logFileLevel, logFile)

			go fileTarget.Crontab()

			Configuration.logger.AddTarget(fileTarget)
		}

		if len(Configuration.logMailHost) > 0 {
			auth := smtp.PlainAuth("", Configuration.logMailUser, Configuration.logMailPassword, Configuration.logMailHost)

			mailTarget := log.NewMailTarget(
				Configuration.logMailLevel,
				Configuration.logMailHost+":"+Configuration.logMailPort,
				Configuration.logMailFrom,
				Configuration.logMailTo,
				auth,
			)

			mailTarget.SetSubject(Configuration.logMailSubject)
			Configuration.logger.AddTarget(mailTarget)
		}
	}

	// Initialize JWT.
	if Configuration.enableJWT {
		// Create JWT instance.
		Configuration.JWT = jwt.NewJWT(Configuration.jwtIssuer, Configuration.jwtTTL)

		// Add HMAC Algorithm.
		hs256, err := jwt.NewHMACAlgorithm(crypto.SHA256, []byte(Configuration.jwtHMACSecretKey))
		if err != nil {
			panic(err)
		}
		Configuration.JWT.AddAlgorithm("HS256", hs256)

		// Add RSA Algorithm.
		var publicKey []byte
		var privateKey []byte
		// Read byte from file.
		publicKey, err = jwt.ReadBytes(Configuration.jwtRSAPublicKey)
		if err == nil {
			privateKey, err = jwt.ReadBytes(Configuration.jwtRSAPrivateKey)
			if err == nil {
				rsa256, err := jwt.NewRSAAlgorithm(crypto.SHA256, publicKey, privateKey)
				if err != nil {
					panic(err)
				}
				Configuration.JWT.AddAlgorithm("RS256", rsa256)
			}
		}
	}

	// Initialize redis cache.
	if Configuration.enableCache {
		redisPool := cache.NewRedisPool(
			Configuration.redisMaxIdle,
			Configuration.redisIdleTimeout,
			Configuration.redisNetwork,
			Configuration.redisAddress,
			Configuration.redisPassword,
			Configuration.redisDb,
		)

		Configuration.cache = cache.NewRedisCache(redisPool)
		_, err := Configuration.cache.GetConn().Do("PING")
		if err != nil {
			panic(err)
		}
	}

	// Initialize session store.
	if Configuration.enableSession {
		if !Configuration.enableCache {
			panic("The session depends on redis cache, please enable the cache component.")
		}

		store := session.NewRedisStore(Configuration.cache.GetPool(), session.Options{Path: "/"})

		store.SetMaxAge(Configuration.sessionMaxAge)

		Configuration.sessionStore = store
	}

}

func checkConfiguration() {
	if len(Configuration.actionPrefix) == 0 && len(Configuration.actionSuffix) == 0 {
		panic("You should set action's prefix or suffix.")
	}
}

func NewApp(domain string) *Application {
	apps[domain] = NewApplication()
	if len(domain) == 0 {
		SetDefaultApp(apps[domain])
	}

	if Configuration.enableSession {
		apps[domain].sessionStore = Configuration.sessionStore
	}

	if Configuration.enableLog {
		apps[domain].logger = Configuration.logger
	}

	if Configuration.enableCache {
		apps[domain].cache = Configuration.cache
	}

	if Configuration.enableJWT {
		apps[domain].jwt = Configuration.JWT
	}

	return apps[domain]
}

func NewRouter() *httprouter.Router {
	return &httprouter.Router{
		RedirectTrailingSlash:  Configuration.routerRedirectTrailingSlash,
		RedirectFixedPath:      Configuration.routerRedirectFixedPath,
		HandleMethodNotAllowed: Configuration.routerHandleMethodNotAllowed,
		HandleOPTIONS:          Configuration.routerHandleOPTIONS,
		NotFound:               &NotFoundHandler{},
		MethodNotAllowed:       &MethodNotAllowedHandler{},
		PanicHandler:           nil,
	}
}

func SetDefaultApp(app *Application) {
	defaultApp = app
}

func Close() {
	if Configuration.enableCache {
		Configuration.cache.GetPool().Close()
	}
	if Configuration.enableLog {
		Configuration.logger.Close()
	}
}

func Run() {
	if defaultApp == nil {
		defaultApp = NewApplication()
	}

	for _, app := range apps {
		app.Run()
	}

	err := http.ListenAndServe(Configuration.serverHost, apps)
	if err != nil {
		panic(err)
	}
}

// ============================== Helper ==============================
// Pretty name, for example, "PostComments" will be formated as "post-comments".
func PrettyName(name string) string {
	if len(name) == 0 {
		return ""
	}
	prettyRoute := strings.ToLower(string(name[0]))
	for i := 1; i < len(name); i++ {
		c := name[i]
		if ('A' <= c) && (c <= 'Z') {
			prettyRoute += "-" + string(rune(int(c)+32))
		} else {
			prettyRoute += string(c)
		}
	}
	return prettyRoute
}

// Remove the controller's prefix and suffix from action's name.
func getControllerName(name string) string {
	// remove prefix.
	if len(Configuration.controllerPrefix) > 0 {
		if 0 != strings.Index(name, Configuration.controllerPrefix) {
			return ""
		}

		prefixLen := len(Configuration.controllerPrefix)
		name = stringutil.SubString(name, prefixLen, len(name)-prefixLen)
	}
	// remove suffix.
	if len(Configuration.controllerSuffix) > 0 {
		pos := len(name) - len(Configuration.controllerSuffix)

		if (pos == -1) || (pos != strings.Index(name, Configuration.controllerSuffix)) {
			return ""
		}

		name = stringutil.SubString(name, 0, pos)
	}
	return name
}

// Remove the action's prefix and suffix from action's name.
func getActionName(name string) string {
	// remove prefix.
	if len(Configuration.actionPrefix) > 0 {
		if 0 == strings.Index(name, Configuration.actionPrefix) {
			prefixLen := len(Configuration.actionPrefix)
			name = stringutil.SubString(name, prefixLen, len(name)-prefixLen)
		}
	}
	// remove suffix.
	if len(Configuration.actionSuffix) > 0 {
		pos := len(name) - len(Configuration.actionSuffix)

		if (pos != -1) || (pos != strings.Index(name, Configuration.actionSuffix)) {
			name = stringutil.SubString(name, 0, pos)
		}
	}

	return name
}

func GoPath() string {
	return goPath
}

func SrcPath() string {
	return srcPath
}
