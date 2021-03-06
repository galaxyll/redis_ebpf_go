package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/galaxyll/redis_ebpf_go/config"
	"github.com/galaxyll/redis_ebpf_go/plug"
)

func main() {
	addr := config.Conf.ServerConf.Addr
	log.Printf("addr %s", addr)

	http.HandleFunc("/duration", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		cmd := req.URL.Query().Get("command")
		if !in(cmd) {
			w.Write([]byte("command is not support at this time"))
			return
		}
		secstr := req.URL.Query().Get("seconds")
		seconds, err := strconv.ParseInt(secstr, 10, 64)
		if err != nil {
			fmt.Println("second parse faild:", err)
			w.Write([]byte("second parse faild"))
			return
		}
		if seconds > 1800 {
			seconds = 1800
		}
		if seconds <= 0 {
			w.Write([]byte("invaild parameter seconds"))
			return
		}
		fmt.Println("the parm after deal: ", cmd, " ", seconds)
		go plug.Duration(cmd, seconds)
		http.Redirect(w, req, config.Conf.ServerConf.Grafana, http.StatusTemporaryRedirect)
	})

	fmt.Println("Server start...")
	err := http.ListenAndServe(addr, nil)
	if err != nil || err != http.ErrServerClosed {
		fmt.Printf("server start faild: %s\n", err)
	}

}

func in(target string) bool {
	for key := range config.COMMAND {
		if key == target {
			return true
		}
	}
	return false
}
