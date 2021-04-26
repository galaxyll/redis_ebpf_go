package db

import (
	"fmt"

	"github.com/galaxyll/redis_ebpf_go/event"
	client "github.com/influxdata/influxdb1-client/v2"
)

func NewClient() client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "nacl",
		Password: "170607",
	})

	if err != nil {
		fmt.Println("create new client faild：", err)
	}
	return cli
}

func Insert(event event.GetEvent) {
	c := NewClient()
	defer c.Close()

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "redis_info",
		Precision: "ns",
	})
	tags := map[string]string{
		"host":    "39.104.13.134",
		"service": "duration",
	}
	fileds := map[string]interface{}{}
	fileds["pid"] = int64(event.Pid)
	fileds["duration"] = int64(event.Duration)
	fileds["key"] = string(event.Key[:event.Klen])
	fileds["klen"] = int32(event.Klen)
	//	pt, err := client.NewPoint("duration_info", tags, fileds, time.Unix(0, int64(event.Start_time_ns)))
	pt, err := client.NewPoint("duration_info", tags, fileds)
	if err != nil {
		fmt.Println("insert point to influxdb faild: ", err)
	}
	bp.AddPoint(pt)
	err = c.Write(bp)
	if err != nil {
		fmt.Println("write db faild: ", err)
	}
}
