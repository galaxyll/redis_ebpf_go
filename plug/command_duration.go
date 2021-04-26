package plug

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/galaxyll/redis_ebpf_go/bpf"
	"github.com/galaxyll/redis_ebpf_go/db"
	ent "github.com/galaxyll/redis_ebpf_go/event"
	"github.com/iovisor/gobpf/bcc"
)

var COMMAND = map[string]string{
	"GET":   "lookupKeyReadOrReply",
	"SET":   "setKey",
	"INCR":  "incrDecrCommand",
	"DECR":  "incrDecrCommand",
	"LPUSH": "pushGenericCommand",
	"RPUSH": "pushGenericCommand",
	"LPOP":  "popGenericCommand",
	"RPOP":  "popGenericCommand",
	"SADD":  "saddCommand",
	"HSET":  "hsetCommand",
	"SPOP":  "spopCommand",
	"MSET":  "msetGenericCommand",
}

var binaryProg string = "/usr/local/bin/redis-server"

// func init() {
// 	flag.StringVar(&binaryProg, "binary", "", "the binary to probe")
// }

func Duration(cmd string, seconds int64) {
	// flag.Parse()
	// if len(binaryProg) == 0 {
	// 	panic("argument --binary must be specified")
	// }
	bccMode := bcc.NewModule(bpf.Bpf_source, []string{})
	uprobeFD, err := bccMode.LoadUprobe("trace_start_time")
	if err != nil {
		panic(err)
	}
	err = bccMode.AttachUprobe(binaryProg, COMMAND[cmd], uprobeFD, -1)
	if err != nil {
		panic(err)
	}

	uretprobeFD, err := bccMode.LoadUprobe("send_duration")
	if err != nil {
		panic(err)
	}
	err = bccMode.AttachUretprobe(binaryProg, COMMAND[cmd], uretprobeFD, -1)
	if err != nil {
		panic(err)
	}

	table := bcc.NewTable(bccMode.TableId("duration_events"), bccMode)
	ch := make(chan []byte)

	pm, err := bcc.InitPerfMap(table, ch, nil)
	if err != nil {
		panic(err)
	}

	intCh := make(chan os.Signal, 1)
	signal.Notify(intCh, os.Interrupt)

	go func() {
		var event ent.GetEvent
		for {
			data := <-ch
			bf := bytes.NewBuffer(data)
			err := binary.Read(bf, binary.LittleEndian, &event.Pid)
			if err != nil {
				fmt.Printf("faild to parse receive data: %s\n", err)
				continue
			}
			fmt.Println(event.Pid)
			err = binary.Read(bf, binary.LittleEndian, &event.Pad)
			if err != nil {
				fmt.Printf("faild to parse receive data: %s\n", err)
				continue
			}
			err = binary.Read(bf, binary.LittleEndian, &event.Start_time_ns)
			if err != nil {
				fmt.Printf("faild to parse receive data: %s\n", err)
				continue
			}
			fmt.Println(event.Start_time_ns)
			err = binary.Read(bf, binary.LittleEndian, &event.Duration)
			if err != nil {
				fmt.Printf("faild to parse receive data: %s\n", err)
				continue
			}
			fmt.Println(event.Duration)
			err = binary.Read(bf, binary.LittleEndian, &event.Klen)
			if err != nil {
				fmt.Printf("faild to parse receive data: %s\n", err)
				continue
			}
			fmt.Println(event.Klen)
			key := make([]byte, 128)
			err = binary.Read(bf, binary.LittleEndian, key)
			if err != nil {
				fmt.Printf("faild to parse receive data: %s\n", err)
				continue
			}
			copy(event.Key[:], key)
			db.Insert(event)
			fmt.Println(string(key))
			fmt.Println(event)
		}
	}()

	pm.Start()
	time.Sleep(time.Duration(seconds) * time.Second)
	pm.Stop()
}
