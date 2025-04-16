package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var loginFormTmpl = []byte(`
<html>
© geekbrains.ru
18<body>
<form action="/login" method="post">
Login: <input type="text" name="login">
Password: <input type="password" name="password">
<input type="submit" value="Login">
</form>
</body>
</html>
`)

const (
	loginValue    = "login"
	passwordValue = "password"
)

var welcome = "Welcome, %s <br />\nSession User-Agent: %s <br />\n<a href=\"/logout\">logout</a>"

func (c RedisClient) RootHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := c.checkSession(r)
	if err != nil {
		err = fmt.Errorf("check session: %w", err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if sess == nil {
		_, _ = w.Write(loginFormTmpl)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprintln(w, fmt.Sprintf(welcome, sess.Login, sess.Useragent))
}

var users = map[string]string{
	"geek": "brains",
}

const cookieName = "session_id"

func (c RedisClient) LoginHandler(w http.ResponseWriter, r *http.Request) {
	inputLogin := r.FormValue(loginValue)
	inputPass := r.FormValue(passwordValue)
	// common map!!! dont make the same in production
	pass, exist := users[inputLogin]
	if !exist || pass != inputPass {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sess, err := c.Create(Session{
		Login:     inputLogin,
		Useragent: r.UserAgent(),
	})
	if err != nil {
		err = fmt.Errorf("create session: %w", err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:    cookieName,
		Value:   sess.ID,
		Expires: time.Now().Add(c.TTL),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (c RedisClient) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie(cookieName)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if err != nil {
		err = fmt.Errorf("read cookie %q: %w", cookieName, err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.Delete(SessionID{ID: session.Value})
	if err != nil {
		err = fmt.Errorf("delete session value %q: %w", session.Value, err)
		log.Printf("[ERR] %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (c RedisClient) checkSession(r *http.Request) (*Session, error) {
	cookieSessionID, err := r.Cookie(cookieName)
	if err == http.ErrNoCookie {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	sess, err := c.Check(SessionID{ID: cookieSessionID.Value})
	if err != nil {
		return nil, fmt.Errorf("check session value %q: %w", cookieSessionID.Value, err)
	}
	return sess, nil
}
