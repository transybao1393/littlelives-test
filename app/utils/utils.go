package utils

import (
	"net/http"
	"strings"
)

var VideoContentType = []string{
	"video/mp4",
	"video/x-m4v",
	"video/x-matroska",
	"video/webm",
	"video/quicktime",
	"video/x-msvideo",
	"video/x-ms-wmv",
	"video/mpeg",
	"video/x-flv",
	"video/",
}

var FileContentType = []string{
	"application/pdf",
	"application/zip",
	"application/x-rar-compressed",
	"application/x-tar",
	"application/x-7z-compressed",
	"application/x-gzip",
	"application/x-bzip2",
	"application/x-bzip",
	"application/x-apple-diskimage",
	"application/x-deb",
	"application/x-rpm",
	"application/x-msdownload",
	"application/x-shockwave-flash",
	"application/octet-stream",
}

func UserIPHandling(r *http.Request) string {
	userIP := strings.Split(r.Header.Get("X-Forwarded-For"), ":")[0]
	if userIP == "" {
		userIP = strings.Split(r.RemoteAddr, ":")[0]
	} else {
		userIP = r.Header.Get("Proxy-Client-IP")
	}
	return userIP
}
