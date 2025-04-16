package red

import (
	"net/http"
	"strconv"
	"time"
)

func (w pusherWriter) Push(target string, opts *http.PushOptions) error {
	return w.ResponseWriter.(http.Pusher).Push(target, opts)
}

var builders = make([]func(w *responseWriter) measurable, 1<<4) //45.04

func init() {
	builders[0] = func(w *responseWriter) measurable { return w }
	builders[flusher] = func(w *responseWriter) measurable { return flusherWriter{w} }
	//...
	builders[hijacker] = func(w *responseWriter) measurable { return hijackerWriter{w} }
	builders[readerFrom] = func(w *responseWriter) measurable { return readerFromWriter{w} }
	builders[pusher] = func(w *responseWriter) measurable { return pusherWriter{w} }

}

var MeasurableHandler = func(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		m := r.Method
		p := r.URL.Path
		requestsTotal.WithLabelValues(p, m).Inc()
		mw := newMeasurableWriter(w)
		h(mw, r)
		if mw.Status()/100 > 3 {
			errorsTotal.WithLabelValues(p, m, strconv.Itoa(mw.Status())).Inc()
		}
		duration.WithLabelValues(p, m,
			strconv.Itoa(mw.Status())).Observe(time.Since(t).Seconds())
	}
}
