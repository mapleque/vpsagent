package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mapleque/vpsagent/core"
)

func checkIp(ip string, list []string) error {
	// This can be optimize with hash map
	for _, allowIp := range list {
		if allowIp == ip {
			return nil
		}
	}
	return fmt.Errorf("client ip %s not in white list")
}

func (s *Server) processRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	status := 400
	start := time.Now()
	ip := getRemoteIp(r)
	message := ""
	timestamp := r.Header.Get("Timestamp")
	sign := r.Header.Get("Signature")
	key := sign + timestamp

	defer func() {
		end := time.Now()
		s.log.Log(
			status,
			ip,
			start.Format("2006-01-02 15:04:05.006"),
			end.Format("2006-01-02 15:04:05.006"),
			end.Sub(start),
			message,
		)
	}()

	if err := s.checkConflict(key); err != nil {
		message = s.badRequest(w, err)
		return
	}

	go func() {
		s.mux.Lock()
		defer s.mux.Unlock()
		s.cache[key] = true
	}()
	defer func() {
		s.mux.Lock()
		defer s.mux.Unlock()
		if _, ok := s.cache[key]; ok {
			s.cache[key] = false
		}
	}()
	time.AfterFunc(1*time.Second, func() {
		s.mux.Lock()
		defer s.mux.Unlock()
		delete(s.cache, key)
	})

	if err := checkIp(ip, s.ipWhiteList); err != nil {
		message = s.badRequest(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		message = s.badRequest(w, err)
		return
	}
	r.Body.Close()

	if err := checkTimestamp(timestamp); err != nil {
		message = s.badRequest(w, err)
		return
	}

	if err := checkSign(sign, s.token, timestamp, body); err != nil {
		message = s.badRequest(w, err)
		return
	}

	res, err := executeScript(body)
	if err != nil {
		message = s.badRequest(w, err)
		return
	}
	status = 200
	w.Write(res)
}

func checkSign(sign, token, timestamp string, body []byte) error {
	tarSign := core.MakeSignature(token, timestamp, body)
	if sign != tarSign {
		return errors.New("invalid signature")
	}
	return nil
}

func checkTimestamp(timestamp string) error {
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp %s", timestamp)
	}
	then := time.Unix(ts, 0)
	if time.Now().Sub(then) > 5*time.Second {
		return errors.New("request is expired")
	}
	return nil
}

func (s *Server) checkConflict(key string) error {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if _, exist := s.cache[key]; exist {
		return errors.New("request is conflict")
	}
	return nil
}

func getRemoteIp(r *http.Request) string {
	ips := r.Header.Get("X-Forwarded-For")
	ip := ""
	if ips != "" {
		ip = strings.Split(ips, ",")[0]
	}
	if ip == "" {
		ip = r.Header.Get("X-Real-Ip")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	if xffHost, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		ip = xffHost
	}

	return ip
}

func (s *Server) badRequest(w http.ResponseWriter, err error) string {
	w.WriteHeader(400)
	w.Write([]byte(err.Error()))
	s.log.Error(err)
	return err.Error()
}
