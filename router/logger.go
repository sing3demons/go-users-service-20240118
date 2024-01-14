package router

import (
	"encoding/json"
	"io"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sing3demons/users/constant"
	"github.com/sing3demons/users/utils"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Starting time request
		startTime := time.Now()
		// Processing request
		body, _ := io.ReadAll(ctx.Request.Body)
		ctx.Request.Body = io.NopCloser(strings.NewReader(string(body)))
		reqId := ctx.Writer.Header().Get(constant.XSessionId)
		if reqId == "" {
			reqId = uuid.NewString()
			ctx.Writer.Header().Set(constant.XSessionId, reqId)
		}

		sessionId := ctx.Writer.Header().Get(constant.XSessionId)
		if sessionId == "" {
			ctx.Writer.Header().Set(constant.XSessionId, uuid.NewString())
		}
		ctx.Next()
		// End Time request
		endTime := time.Now()
		// Request method
		reqMethod := ctx.Request.Method
		// Request route
		path := ctx.Request.RequestURI
		// status code
		statusCode := ctx.Writer.Status()
		// Request host
		host := ctx.Request.Host
		// Request user agent
		userID, exists := ctx.Get("userId")
		if exists {
			userID = userID.(string)
		} else {
			userID = ""
		}

		bodySize := ctx.Writer.Size()
		// execution time
		latencyTime := endTime.Sub(startTime)

		headers := GetHeaders(ctx)

		var bodyJson any
		json.Unmarshal(body, &bodyJson)

		logrus.WithFields(logrus.Fields{
			"uuid":          reqId,
			"headers":       headers,
			"body":          utils.MaskSensitiveData(bodyJson),
			"method":        reqMethod,
			"status":        statusCode,
			"latency":       latencyTime,
			"error":         ctx.Errors.ByType(gin.ErrorTypePrivate).String(),
			"request":       ctx.Request.PostForm.Encode(),
			"body_size":     bodySize,
			"host":          host,
			"protocol":      ctx.Request.Proto,
			"path":          path,
			"query":         ctx.Request.URL.RawQuery,
			"response_size": ctx.Writer.Size(),
			"timezone":      time.Now().Location().String(),
			"ISOTime":       startTime,
			"UnixTime":      startTime.UnixNano(),
		}).Info("HTTP::REQUEST")
		ctx.Next()
	}
}

func GetReqId() string {
	return uuid.NewString()
}

func GetHeaders(ctx *gin.Context) map[string]any {
	// Request user agent
	userAgent := ctx.Request.UserAgent()
	platform := strings.Split(ctx.Request.Header.Get("sec-ch-ua"), ",")
	mobile := ctx.Request.Header.Get("sec-ch-ua-mobile")
	operatingSystem := ctx.Request.Header.Get("sec-ch-ua-platform")
	clientIP := ctx.ClientIP()
	reqId := ctx.Writer.Header().Get("X-Session-Id")
	if reqId == "" {
		reqId = uuid.NewString()
	}

	macIp := getMACAndIP()

	return map[string]any{
		"user_agent": userAgent,
		"Platform":   platform,
		"Mobile":     mobile,
		"OS":         operatingSystem,
		"client_ip":  clientIP,
		"request_id": reqId,
		"remote_ip":  ctx.Request.RemoteAddr,
		"mac_ip":     macIp,
	}
}

func getMACAndIP() MacIP {
	interfaces, _ := net.Interfaces()
	macAddr := MacIP{}
	for _, iface := range interfaces {

		if iface.Name != "" {
			macAddr.InterfaceName = iface.Name
		}

		if iface.HardwareAddr != nil {
			macAddr.HardwareAddr = iface.HardwareAddr.String()
		}

		var ips []string
		addrs, _ := iface.Addrs()

		for _, addr := range addrs {
			ips = append(ips, addr.String())
		}

		if len(ips) > 0 {
			macAddr.IPs = ips
		}
	}

	return macAddr
}

type MacIP struct {
	InterfaceName string   `json:"interface_name"`
	HardwareAddr  string   `json:"hardware_addr"`
	IPs           []string `json:"ips"`
}
