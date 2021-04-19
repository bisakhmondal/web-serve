package core

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"strings"
)

var HTMLBlackListRoutes = []string{
	"/ping",
	"/admin",
}
var NoCustomisedRoutes = []string{
	"/custom",
}

func routeConfig(efs embed.FS) (*gin.Engine, error) {
	router := gin.New()
	filesys := fs.FS(efs)
	subtree, err := fs.Sub(filesys, "html")
	if err != nil {
		return nil, err
	}
	router.GET("/ping", func(g *gin.Context) {
		g.String(http.StatusOK, "pong")
	})
	router.Use(noTouch())
	router.Use(serve("/", subtree))
	router.NoRoute(func(g *gin.Context) {
		if g.Request.Method == "GET" {
			g.FileFromFS("404.html", http.FS(subtree))
			g.Abort()
		} else {
			g.String(http.StatusNotFound, "404 method not allowed")
		}
	})

	return router, nil
}

type spaFileSystem struct {
	http.FileSystem
}

func (fs *spaFileSystem) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	//Default failsafe page
	if err != nil {
		return fs.FileSystem.Open("build/index.html")
	}
	return f, err
}

func serve(urlPrefix string, fss fs.FS) gin.HandlerFunc {
	fileserver := http.FileServer(&spaFileSystem{http.FS(fss)})
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if c.Request.Method == "GET" && !isBlacklisted(c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func noTouch() gin.HandlerFunc {
	return func(g *gin.Context) {
		if noCustomised(g.Request.URL.Path) {
			g.JSON(http.StatusNotFound,
				gin.H{
					"code":    "PAGE_NOT_FOUND",
					"message": "Page not found",
				})
			g.Abort()
		}
		g.Next()

	}
}

func isBlacklisted(path string) bool {
	for _, checkpath := range HTMLBlackListRoutes {
		if strings.HasPrefix(path, checkpath) {
			return true
		}
	}
	return false
}

func noCustomised(path string) bool {
	for _, checkpath := range NoCustomisedRoutes {
		if strings.HasPrefix(path, checkpath) {
			return true
		}
	}
	return false
}

func exists(prefix string, filepath string, f *fs.FS) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		_, err := fs.Stat(*f, p)
		if err != nil {
			return false
		}
	}
	return true
}
