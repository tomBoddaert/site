package main

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const DEFAULT_ADDRESS = "localhost:8080"

type Resolve struct {
	Resolved bool
	Path     string
	ModTime  time.Time
}

type Handler struct {
	Config *Config
}

func getServeAddress() string {
	address := os.Getenv("SITE_ADDRESS")
	if address == "" {
		address = DEFAULT_ADDRESS
	}

	return address
}

func resolve(config *Config, url_path string) Resolve {
	base := path.Join(config.DstDir, url_path)

	if !strings.HasPrefix(base, config.DstDir) {
		return Resolve{Resolved: false}
	}

	entry, err := os.Stat(base)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			check(err)
		}
	} else {
		if !entry.IsDir() {
			return Resolve{
				Resolved: true,
				Path:     base,
				ModTime:  entry.ModTime(),
			}
		}
	}

	if !strings.HasSuffix(url_path, ".html") {
		resolution := resolve(config, url_path+".html")
		if resolution.Resolved {
			return resolution
		}
	}

	if entry != nil && entry.IsDir() {
		resolution := resolve(config, path.Join(url_path, "index"))
		if resolution.Resolved {
			return resolution
		}
	}

	return Resolve{Resolved: false}
}

func (handler Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	url_path := req.URL.Path

	logger.Debugf("Incoming request (%v)", url_path)

	resolution := resolve(handler.Config, url_path)
	if !resolution.Resolved {
		logger.Debugf("Request unresolved (%v -x 404)", url_path)
		res.WriteHeader(404)
		resolution = resolve(handler.Config, handler.Config.NotFoundPath)
	} else {
		logger.Debugf("Request resolved (%v -> %v)", url_path, resolution.Path)
	}

	if resolution.Resolved {
		file, err := os.Open(resolution.Path)
		check(err)
		defer file.Close()

		http.ServeContent(res, req, resolution.Path, resolution.ModTime, file)
	} else {
		res.Write([]byte("404 - Not Found"))
	}
}

func serve(config Config) {
	logger.SetReportTimestamp(true)
	logger.Info("Starting server...")
	logger.SetReportTimestamp(false)

	handler := Handler{Config: &config}

	address := getServeAddress()

	logger.Infof("Hosting on 'http://%v/' (THIS SERVER IS FOR TESTING ONLY! DO NOT EXPOSE IT)", address)
	logger.Print("Press '[Ctrl] + C' to exit")

	check(http.ListenAndServe(address, handler))
}
