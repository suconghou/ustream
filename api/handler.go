package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/suconghou/videoproxy/route"
)

var (
	startTime = time.Now()
	logger    = log.New(os.Stdout, "", 0)
)

var sysStatus struct {
	Uptime       string
	GoVersion    string
	Hostname     string
	MemAllocated uint64 // bytes allocated and still in use
	MemTotal     uint64 // bytes allocated (even if freed)
	MemSys       uint64 // bytes obtained from system
	NumGoroutine int
	CPUNum       int
	Pid          int
}

// Handler hand
func Handler(w http.ResponseWriter, r *http.Request) {
	routeMatch(w, r)
}

func routeMatch(w http.ResponseWriter, r *http.Request) {
	for _, p := range route.Route {
		if p.Reg.MatchString(r.URL.Path) {
			if err := p.Handler(w, r, p.Reg.FindStringSubmatch(r.URL.Path)); err != nil {
				logger.Print(err)
			}
			return
		}
	}
	if r.URL.Path == "/video/api/status" {
		status(w, r)
		return
	}
	fallback(w, r)
}

func status(w http.ResponseWriter, r *http.Request) {
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	sysStatus.Uptime = time.Since(startTime).String()
	sysStatus.NumGoroutine = runtime.NumGoroutine()
	sysStatus.MemAllocated = memStat.Alloc / 1024  // 当前内存使用量
	sysStatus.MemTotal = memStat.TotalAlloc / 1024 // 所有被分配的内存
	sysStatus.MemSys = memStat.Sys / 1024          // 内存占用量
	sysStatus.CPUNum = runtime.NumCPU()
	sysStatus.GoVersion = runtime.Version()
	sysStatus.Hostname, _ = os.Hostname()
	sysStatus.Pid = os.Getpid()
	if bs, err := json.Marshal(&sysStatus); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(bs)
	}
}

func fallback(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
