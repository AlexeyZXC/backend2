package red

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	StatusOk              = "ok"
	StatusConnectionError = "connect_error"
	StatusReadError       = "read_error"
	StatusParseError      = "parse_error"
)

type (
	measurable interface {
		http.ResponseWriter
		Status() int
	}

	responseWriter struct {
		http.ResponseWriter
		statusCode  int
		wroteHeader bool
	}
)

func newMeasurableWriter(w http.ResponseWriter) measurable {
	var id int
	if _, ok := w.(http.Flusher); ok {
		id += flusher
	}
	if _, ok := w.(http.Hijacker); ok {
		id += hijacker
	}
	if _, ok := w.(io.ReaderFrom); ok {
		id += readerFrom
	}
	if _, ok := w.(http.Pusher); ok {
		id += pusher
	}

	return builders[id](&responseWriter{w, 0, false})
}

func (w *responseWriter) Status() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}

	return w.statusCode
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(code)
}

const (
	flusher = 1 << iota
	hijacker
	readerFrom
	pusher
)

type (
	flusherWriter    struct{ *responseWriter }
	hijackerWriter   struct{ *responseWriter }
	readerFromWriter struct{ *responseWriter }
	pusherWriter     struct{ *responseWriter }
)

func (w flusherWriter) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w hijackerWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w readerFromWriter) ReadFrom(r io.Reader) (int64, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.(io.ReaderFrom).ReadFrom(r)
}

func (w pusherWriter) Push(target string, opts *http.PushOptions) error {
	return w.ResponseWriter.(http.Pusher).Push(target, opts)
}

var builders = make([]func(w *responseWriter) measurable, 1<<4)

func init() {
	builders[0] = func(w *responseWriter) measurable { return w }
	builders[flusher] = func(w *responseWriter) measurable { return flusherWriter{w} }
	builders[hijacker] = func(w *responseWriter) measurable { return hijackerWriter{w} }
	builders[readerFrom] = func(w *responseWriter) measurable { return readerFromWriter{w} }
	builders[pusher] = func(w *responseWriter) measurable { return pusherWriter{w} }

	builders[flusher+hijacker] = func(w *responseWriter) measurable {
		return struct {
			*responseWriter
			http.Hijacker
			http.Flusher
		}{w, hijackerWriter{w}, flusherWriter{w}}
	}

	builders[flusher+readerFrom] = func(w *responseWriter) measurable {
		return struct {
			*responseWriter
			io.ReaderFrom
			http.Flusher
		}{w, readerFromWriter{w}, flusherWriter{w}}
	}

	builders[flusher+pusher] = func(w *responseWriter) measurable {
		return struct {
			*responseWriter
			http.Pusher
			http.Flusher
		}{w, pusherWriter{w}, flusherWriter{w}}
	}

	builders[hijacker+pusher] = func(w *responseWriter) measurable {
		return struct {
			*responseWriter
			http.Hijacker
			http.Flusher
		}{w, hijackerWriter{w}, flusherWriter{w}}
	}

	builders[hijacker+readerFrom] = func(w *responseWriter) measurable {
		return struct {
			*responseWriter
			http.Hijacker
			io.ReaderFrom
		}{w, hijackerWriter{w}, readerFromWriter{w}}
	}

	builders[pusher+readerFrom] = func(w *responseWriter) measurable {
		return struct {
			*responseWriter
			http.Pusher
			io.ReaderFrom
		}{w, pusherWriter{w}, readerFromWriter{w}}
	}
	//3
	builders[hijacker+pusher+readerFrom] = func(w *responseWriter) measurable {
		return struct {
			*responseWriter
			http.Pusher
			io.ReaderFrom
			http.Hijacker
		}{w, pusherWriter{w}, readerFromWriter{w}, hijackerWriter{w}}
	}

	builders[flusher+pusher+readerFrom] = func(w *responseWriter) measurable {
		return struct {
			*responseWriter
			http.Pusher
			io.ReaderFrom
			http.Flusher
		}{w, pusherWriter{w}, readerFromWriter{w}, flusherWriter{w}}
	}

	builders[flusher+hijacker+readerFrom] = func(w *responseWriter) measurable {
		return struct {
			*responseWriter
			io.ReaderFrom
			http.Hijacker
			http.Flusher
		}{w, readerFromWriter{w}, hijackerWriter{w}, flusherWriter{w}}
	}

	builders[flusher+hijacker+pusher] = func(w *responseWriter) measurable {
		return struct {
			*responseWriter
			http.Pusher
			http.Hijacker
			http.Flusher
		}{w, pusherWriter{w}, hijackerWriter{w}, flusherWriter{w}}
	}

	//4
	builders[flusher+hijacker+pusher+readerFrom] = func(w *responseWriter) measurable {
		return struct {
			*responseWriter
			http.Pusher
			io.ReaderFrom
			http.Hijacker
			http.Flusher
		}{w, pusherWriter{w}, readerFromWriter{w}, hijackerWriter{w}, flusherWriter{w}}
	}
}

var MeasurableHandler = func(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		m := r.Method
		p := r.URL.Path
		requestsTotal.
			WithLabelValues(p, m).
			Inc()

		mw := newMeasurableWriter(w)
		h(mw, r)
		if mw.Status()/100 > 3 {
			errorsTotal.
				WithLabelValues(p, m, strconv.Itoa(mw.Status())).
				Inc()
		}

		duration.
			WithLabelValues(p, m, strconv.Itoa(mw.Status())).
			Observe(time.Since(t).Seconds())
	}
}
