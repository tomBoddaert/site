package main

import (
	"errors"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const defaultAddress = "localhost:8080"

type resolution struct {
	Resolved bool
	Path     string
	ModTime  time.Time
}

type handler struct {
	Config *Config
}

func getServeAddress() string {
	address := os.Getenv("SITE_ADDRESS")
	if address == "" {
		address = defaultAddress
	}

	return address
}

func resolve(config *Config, urlPath string) resolution {
	base := path.Join(config.DstDir, urlPath)

	if !strings.HasPrefix(base, config.DstDir) {
		return resolution{Resolved: false}
	}

	entry, err := os.Stat(base)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			check(err)
		}
	} else {
		if !entry.IsDir() {
			return resolution{
				Resolved: true,
				Path:     base,
				ModTime:  entry.ModTime(),
			}
		}
	}

	if !strings.HasSuffix(urlPath, ".html") {
		resolution := resolve(config, urlPath+".html")
		if resolution.Resolved {
			return resolution
		}
	}

	if entry != nil && entry.IsDir() {
		resolution := resolve(config, path.Join(urlPath, "index"))
		if resolution.Resolved {
			return resolution
		}
	}

	return resolution{Resolved: false}
}

func (handler handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	urlPath := req.URL.Path

	logger.Debugf("Incoming request (%v)", urlPath)

	resolution := resolve(handler.Config, urlPath)
	if !resolution.Resolved {
		logger.Debugf("Request unresolved (%v -x 404)", urlPath)
		res.WriteHeader(404)
		resolution = resolve(handler.Config, handler.Config.NotFoundPath)
	} else {
		logger.Debugf("Request resolved (%v -> %v)", urlPath, resolution.Path)
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

func Serve(config Config) {
	logger.SetReportTimestamp(true)
	logger.Info("Starting server...")
	logger.SetReportTimestamp(false)

	handler := handler{Config: &config}

	address := getServeAddress()

	logger.Infof("Hosting on 'http://%v/'\n  (THIS SERVER IS FOR TESTING ONLY! DO NOT EXPOSE IT)", address)
	logger.Print("Press '[Ctrl] + C' to exit")

	check(http.ListenAndServe(address, handler))
}
