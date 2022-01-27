package web

import (
	"net/http"
)

type Query func(string) string
type Param func(string) string
type Header func(string) string
type Cookie func(string) *http.Cookie

type Method string
type Host string
type RemoteAddr string
type RequestURI string
type ContentLength int64
type ContentType string
type UserAgent string
type Path string
type RawPath string
type RawQuery string
type TransferEncoding []string

type Location string

type Redirect struct {
	Code     int
	Location string
}

type Code int

type HTML []byte

type TEXT []byte

type DATA []byte

type STRING string
