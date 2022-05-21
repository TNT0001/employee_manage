package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
	"tungnt/emmployee_manage/pkg/share/utils"
)

func NewServer(logger *logrus.Logger, db *gorm.DB) *gin.Engine {
	router := gin.New()
	if os.Getenv("STAGE") == "DEV" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(
		NewlogginMiddleware(logger),
		gin.Recovery(),
		accessControlMiddleware,
	)

	router.GET("/", HelloHandler)
	router.GET("/health-check", NewHealthCheckHandler(db))
	router.GET("/metrics", prometheus())

	return router
}

func prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}

func accessControlMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Credentials", "true")
	// Required to handler the cache
	// Reference: https://dw-ml-nfc.atlassian.net/browse/MLI-404
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")

	// Required to handler the token/cookies properly on the frontend side
	// Reference: https://monstarlab.slack.com/archives/G01FDH423RQ/p1620952537280600
	// The wildcard value does not work when allowing credentials
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin#directives
	if origin := c.Request.Header.Get("Origin"); origin != "" {
		c.Header("Access-Control-Allow-Origin", origin)
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
	}

	// When allowing credentials, the wildcard does not work with this header
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers#directives
	if requestedHeaders := c.Request.Header.Get("Access-Control-Request-Headers"); requestedHeaders != "" {
		c.Header("Access-Control-Allow-Headers", requestedHeaders)
	} else {
		c.Header("Access-Control-Allow-Headers", "*")
	}

	// When allowing credentials, the wildcard does not work with this header
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods#directives
	if requestedMethod := c.Request.Header.Get("Access-Control-Request-Method"); requestedMethod != "" {
		c.Header("Access-Control-Allow-Methods", requestedMethod)
	} else {
		c.Header("Access-Control-Allow-Methods", "*")
	}

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	} else {
		c.Set(utils.RequestIDKey, uuid.NewString())
		c.Next()
	}
}

func NewlogginMiddleware(logger *logrus.Logger) func(c *gin.Context) {
	loggingMiddleware := func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		field := logrus.Fields{
			"Host":       c.Request.Host,
			"Durations":  duration.String(),
			"Method":     c.Request.Method,
			"Url":        c.Request.URL,
			"Status":     c.Writer.Status(),
			"Referer":    c.Request.Referer(),
			"User-agent": c.Request.UserAgent(),
			"IP":         utils.GetClientIP(c),
		}

		if logFiled, exist := c.Keys[utils.LoggerField].(map[string]interface{}); exist {
			for k, v := range logFiled {
				field[k] = v
			}
		}

		if c.Writer.Status()/100 == 4 || c.Writer.Status()/100 == 5 {
			logger.WithFields(field).Error()
			return
		} else {
			logger.WithFields(field).Info()
		}
	}

	return loggingMiddleware
}

func NewHealthCheckHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ping(db); err != nil {
			ctx.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}

func HelloHandler(ctx *gin.Context) {
	_, err := ctx.Writer.WriteString("Hello")
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}
