package render

import "net/http"

type DATA struct {
	Type string
	Data []byte
}

func Data(typ string, data []byte) Render {
	return DATA{typ, data}
}

func (r DATA) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	_, err = w.Write(r.Data)
	return
}

func (r DATA) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, r.Type)
}
