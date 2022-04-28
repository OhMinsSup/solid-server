package web

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// RoutedService 는 새 끝점을 제공하기 위해
// 모든 서비스가 웹 서버에 자신을 등록하는 데 필요한 인터페이스를 정의합니다.
type RoutedService interface {
	RegisterRoutes(*mux.Router)
}

// Server 는 http 웹 서버를 관리하는 구조입니다.
type Server struct {
	http.Server

	baseURL  string
	rootPath string
	port     int
	ssl      bool
	logger   *mlog.Logger
}

func NewServer(rootPath string, serverRoot string, port int, ssl, localOnly bool, logger *mlog.Logger) *Server {
	r := mux.NewRouter()

	var addr string
	if localOnly {
		addr = fmt.Sprintf(`localhost:%d`, port)
	} else {
		addr = fmt.Sprintf(`:%d`, port)
	}

	baseURL := ""
	url, err := url.Parse(serverRoot)
	if err != nil {
		logger.Error("Invalid ServerRoot setting", mlog.Err(err))
	}
	baseURL = url.Path

	ws := &Server{
		Server: http.Server{
			Addr:    addr,
			Handler: r,
		},
		baseURL:  baseURL,
		rootPath: rootPath,
		port:     port,
		ssl:      ssl,
		logger:   logger,
	}

	return ws
}

func (ws *Server) Router() *mux.Router {
	return ws.Server.Handler.(*mux.Router)
}

// AddRoutes 를 사용하면 서비스가 웹 서버 라우터에 자체 등록하고 새 끝점을 제공할 수 있습니다.
func (ws *Server) AddRoutes(rs RoutedService) {
	rs.RegisterRoutes(ws.Router())
}

// Start 웹 서버를 실행하고 연결 수신을 시작합니다.
func (ws *Server) Start() {
	//ws.registerRoutes()
	if ws.port == -1 {
		ws.logger.Debug("server not bind to any port")
		return
	}

	isSSL := ws.ssl && fileExists("./cert/cert.pem") && fileExists("./cert/key.pem")
	if isSSL {
		ws.logger.Info("https server started", mlog.Int("port", ws.port))
		go func() {
			if err := ws.ListenAndServeTLS("./cert/cert.pem", "./cert/key.pem"); err != nil {
				ws.logger.Fatal("ListenAndServeTLS", mlog.Err(err))
			}
		}()

		return
	}

	ws.logger.Info("http server started", mlog.Int("port", ws.port))
	go func() {
		if err := ws.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			ws.logger.Fatal("ListenAndServeTLS", mlog.Err(err))
		}
		ws.logger.Info("http server stopped")
	}()
}

func (ws *Server) Shutdown() error {
	return ws.Close()
}

// fileExists returns true if a file exists at the path.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

// errorOrWarn 은 이 서버 인스턴스가 단위 테스트를 실행 중이면 '경고' 수준을 반환하고, 그렇지 않으면 '오류'를 반환합니다.
func errorOrWarn() mlog.Level {
	unitTesting := strings.ToLower(strings.TrimSpace(os.Getenv("FB_UNIT_TESTING")))
	if unitTesting == "1" || unitTesting == "y" || unitTesting == "t" {
		return mlog.LvlWarn
	}
	return mlog.LvlError
}
