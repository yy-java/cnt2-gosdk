package main

import (
	"fmt"
	"github.com/yy-java/cnt2-gosdk"
	"log"
	"time"
)

func main() {
	cnt2Service, err := cnt2.Start(&cnt2.ClientConfig{
		Endpoints: []string{"1.1.1.1:2379", "2.2.2.2:2379", "3.3.3.3:2379"},
		//Endpoints: []string{"61.147.187.152:2379", "61.147.187.142:2379", "61.147.187.150:2379"},
		App:     "demo",
		Profile: "development",
	})
	if err != nil {
		log.Printf("start cnt2 error")
	}
	if err == nil {
		listener := &CommonListenter{}
		cnt2Service.RegisterListener(listener, "my_config_1", "test_config")

		time.Sleep(time.Minute * 10)
		cnt2Service.Close()
	}
	//if can't start cnt2 or initialize，you can use the lastly config
	result, _ := cnt2Service.GetConfig("my_config_2")
	log.Printf("get from local: %s", result)
}

type CommonListenter struct{}

//通知key的增加或者修改事件
func (t *CommonListenter) HandlePutEvent(config *cnt2.Config) error {
	fmt.Printf("put key: %s ; newValue: %s; version:%s \n", config.Key, config.Value, config.Version)
	return nil
}

//通知key的删除事件
func (t *CommonListenter) HandleDeleteEvent(config *cnt2.Config) error {
	fmt.Printf("delete app:%s, profile:%s, key: %s \n", config.App, config.Profile, config.Key)
	return nil
}
