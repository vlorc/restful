package render

import (
	"fmt"
	"github.com/json-iterator/go"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var eventContentType = "text/event-stream"
var noCache = "no-cache"

var fieldReplacer = strings.NewReplacer(
	"\n", "\\n",
	"\r", "\\r")

var dataReplacer = strings.NewReplacer(
	"\n", "\ndata:",
	"\r", "\\r")

type Event struct {
	Event string
	Id    string
	Retry uint
	Data  interface{}
}

func WriteEvent(resp http.ResponseWriter, event Event) error {
	writeContentType(resp, eventContentType)
	writeNoCache(resp, noCache)

	w := checkWriter(resp)
	writeId(w, event.Id)
	writeEvent(w, event.Event)
	writeRetry(w, event.Retry)
	return writeData(w, event.Data)
}

func writeId(w stringWriter, id string) {
	if len(id) > 0 {
		w.WriteString("id:")
		fieldReplacer.WriteString(w, id)
		w.WriteString("\n")
	}
}

func writeEvent(w stringWriter, event string) {
	if len(event) > 0 {
		w.WriteString("event:")
		fieldReplacer.WriteString(w, event)
		w.WriteString("\n")
	}
}

func writeRetry(w io.StringWriter, retry uint) {
	if retry > 0 {
		w.WriteString("retry:")
		w.WriteString(strconv.FormatUint(uint64(retry), 10))
		w.WriteString("\n")
	}
}

func writeData(w stringWriter, data interface{}) error {
	w.WriteString("data:")
	switch kindOfData(data) {
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
		err := jsoniter.ConfigFastest.NewEncoder(w).Encode(data)
		if err != nil {
			return err
		}
		w.WriteString("\n")
	case reflect.String:
		if v, ok := data.(string); ok {
			dataReplacer.WriteString(w, v)
			w.WriteString("\n\n")
		}
	default:
		dataReplacer.WriteString(w, fmt.Sprint(data))
		w.WriteString("\n\n")
	}
	return nil
}

func (r Event) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, eventContentType)
}

func (r Event) Render(w http.ResponseWriter) error {
	return WriteEvent(w, r)
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
