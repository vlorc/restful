package render

import "io"

type stringWriter interface {
	io.Writer
	io.StringWriter
}

type stringWrapper struct {
	io.Writer
	sw io.StringWriter
}

func (w stringWrapper) WriteString(str string) (int, error) {
	if nil != w.sw {
		return w.sw.WriteString(str)
	}
	return w.Writer.Write([]byte(str))
}

func checkWriter(writer io.Writer) stringWriter {
	if w, ok := writer.(stringWriter); ok {
		return w
	}
	sw, _ := writer.(io.StringWriter)
	return stringWrapper{writer, sw}
}
