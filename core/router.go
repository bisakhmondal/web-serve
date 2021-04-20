// Package core
/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* 	http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package core

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func routeConfig(efs embed.FS, w io.Writer) (*gin.Engine, error) {
	if conf.Configuration.Server.Mode == "DEBUG" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	f, _ := os.OpenFile(filepath.Join(conf.Configuration.Logging.Logdir,
		conf.Configuration.Logging.Logfile),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	router := gin.New()
	//logger
	router.Use(logFilter())
	router.Use(gin.Recovery())
	filesys := fs.FS(efs)
	subtree, err := fs.Sub(filesys, "html")
	if err != nil {
		return nil, err
	}
	router.GET("/ping", func(g *gin.Context) {
		g.String(http.StatusOK, "pong")
	})
	router.Use(blacklistJSON())
	router.Use(serve("/", subtree))
	router.NoRoute(func(g *gin.Context) {
		if g.Request.Method == "GET" {
			g.FileFromFS(strings.TrimPrefix(conf.Configuration.Server.Page404, "html/"),
				http.FS(subtree))
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

	if conf.Configuration.Server.Failsafe {
		//Default failsafe page
		if err != nil {
			return fs.FileSystem.Open("build/index.html")
		}
	}

	return f, err
}

func serve(urlPrefix string, fss fs.FS) gin.HandlerFunc {
	fileserver := http.FileServer(&spaFileSystem{http.FS(fss)})
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if c.Request.Method == "GET" && (exists(urlPrefix, c.Request.URL.Path, &fss) || !blacklistHTML(c.Request.URL.Path)) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func blacklistJSON() gin.HandlerFunc {
	return func(g *gin.Context) {
		if noHTMLCustomised(g.Request.URL.Path) {
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

func blacklistHTML(path string) bool {
	for _, checkpath := range conf.Configuration.Server.HtmlBlackList {
		if strings.HasPrefix(path, checkpath) {
			return true
		}
	}
	return false
}

func noHTMLCustomised(path string) bool {
	for _, checkpath := range conf.Configuration.Server.Blacklist {
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
