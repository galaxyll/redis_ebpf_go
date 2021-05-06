package plug

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/galaxyll/redis_ebpf_go/bpf"
	"github.com/galaxyll/redis_ebpf_go/config"
	"github.com/galaxyll/redis_ebpf_go/db"
	ent "github.com/galaxyll/redis_ebpf_go/event"
	"github.com/iovisor/gobpf/bcc"
)

var binaryProg string = config.Conf.PlugConf.BinaryPath

// func init() {
// 	flag.StringVar(&binaryProg, "binary", "", "the binary to probe")
// }

func Duration(cmd string, seconds int64) error {
	// flag.Parse()
	// if len(binaryProg) == 0 {
	// 	panic("argument --binary must be specified")
	// }
	bccMode := bcc.NewModule(bpf.Get_src, []string{})
	defer bccMode.Close()
	uprobeFD, err := bccMode.LoadUprobe("trace_start_time")
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("/%s/\n", config.COMMAND[cmd])
	err = bccMode.AttachUprobe(binaryProg, config.COMMAND[cmd], uprobeFD, -1)
	if err != nil {
		fmt.Println(err)
		return err
	}

	uretprobeFD, err := bccMode.LoadUprobe("send_duration")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = bccMode.AttachUretprobe(binaryProg, config.COMMAND[cmd], uretprobeFD, -1)
	if err != nil {
		fmt.Println(err)
		return err
	}

	table := bcc.NewTable(bccMode.TableId("duration_events"), bccMode)
	ch := make(chan []byte)

	pm, err := bcc.InitPerfMap(table, ch, nil)
	if err != nil {
		fmt.Println(err)
		return err
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
			key := make([]byte, 128)
			err = binary.Read(bf, binary.LittleEndian, key)
			if err != nil {
				fmt.Printf("faild to parse receive data: %s\n", err)
				continue
			}
			copy(event.Key[:], key)
			db.InsertGetEv(event)
			fmt.Printf("[log] Key:%s duration:%d\n", event.Key, event.Duration)
		}
	}()

	pm.Start()
	time.Sleep(time.Duration(seconds) * time.Second)
	pm.Stop()

	return nil
}
