package db

import (
	"fmt"

	"github.com/galaxyll/redis_ebpf_go/config"
	"github.com/galaxyll/redis_ebpf_go/event"
	client "github.com/influxdata/influxdb1-client/v2"
)

func NewClient() client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Conf.InfluxdbConf.Addr,
		Username: config.Conf.InfluxdbConf.Username,
		Password: config.Conf.InfluxdbConf.Password,
	})

	if err != nil {
		fmt.Println("create new client faild：", err)
	}
	return cli
}

func insert(event event.GetEvent, tb string) error {
	c := NewClient()
	defer c.Close()

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  config.Conf.InfluxdbConf.Database,
		Precision: config.Conf.InfluxdbConf.Precision,
	})
	tags := map[string]string{
		"host":    config.Conf.InfluxdbConf.Tag.Host,
		"service": config.Conf.InfluxdbConf.Tag.Service,
	}
	fileds := map[string]interface{}{}
	fileds["pid"] = int64(event.Pid)
	fileds["duration"] = int64(event.Duration)
	fileds["key"] = string(event.Key[:event.Klen])
	fileds["klen"] = int32(event.Klen - 1)
	//	pt, err := client.NewPoint("duration_info", tags, fileds, time.Unix(0, int64(event.Start_time_ns)))
	pt, err := client.NewPoint(tb, tags, fileds)
	if err != nil {
		fmt.Println("insert point to influxdb faild: ", err)
		return err
	}
	bp.AddPoint(pt)
	err = c.Write(bp)
	if err != nil {
		fmt.Println("write db faild: ", err)
		return err
	}
	return nil
}

func InsertGetEv(event event.GetEvent) error {
	err := insert(event, "get")
	if err != nil {
		return err
	}
	return nil
}

func InsertSetEv(event event.GetEvent) error {
	err := insert(event, "set")
	if err != nil {
		return err
	}
	return nil
}
