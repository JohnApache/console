package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/wzshiming/console"
	"golang.org/x/net/websocket"
)

var (
	sessionsMu sync.RWMutex
	session    = map[string]console.Sessions{}
)

func addSession(id string, s console.Sessions) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	session[id] = s
}

func getSession(id string) console.Sessions {
	sessionsMu.RLock()
	defer sessionsMu.RUnlock()
	return session[id]
}

func delSession(id string) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	delete(session, id)
}

func requests(rc io.ReadCloser, i interface{}) error {
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return err
	}
	defer rc.Close()
	return json.Unmarshal(data, i)
}

type errMsg struct {
	Msg string `json:"msg,omitempty"`
}

func execRouter() (*mux.Router, error) {
	// 路由
	mux0 := mux.NewRouter()
	rend := render.New()

	// 创建连接
	mux0.HandleFunc("/create_exec", func(w http.ResponseWriter, r *http.Request) {
		req := &console.ReqCreateExec{}
		requests(r.Body, &req)

		// 获取驱动
		sesss, err := console.GetDrivers(req.Name, req.Host)
		if err != nil {
			rend.JSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		// 创建连接
		exec, err := sesss.CreateExec(req)

		if err != nil {
			rend.JSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		addSession(exec.EId, sesss)

		rend.JSON(w, http.StatusOK, exec)
		return
	})

	// 开始连接
	mux0.Handle("/start_exec", websocket.Handler(func(ws *websocket.Conn) {
		req := ws.Request()

		eid := req.FormValue("eid")

		client := getSession(eid)
		defer delSession(eid)

		// 执行连接
		err := client.StartExec(eid, ws)
		if err != nil {
			rend.JSON(ws, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		rend.Data(ws, http.StatusSwitchingProtocols, nil)
		return
	}))

	// 窗口大小调整
	mux0.HandleFunc("/resize_exec_tty", func(w http.ResponseWriter, r *http.Request) {
		req := &console.ReqResizeExecTTY{}
		requests(r.Body, &req)

		client := getSession(req.EId)

		err := client.ResizeExecTTY(req)
		if err != nil {
			rend.JSON(w, http.StatusBadRequest, errMsg{err.Error()})
			return
		}

		rend.JSON(w, http.StatusOK, nil)
		return
	})

	return mux0, nil
}
