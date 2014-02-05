package mqtt

import (
	"fmt"
	"mqtt/db/ssdb"
	log "mqtt/seelog"
	"strings"
	"sync"
	"time"
)

type SsdbClient struct {
	Client *ssdb.Client
	host   string
	port   int
}

var g_ssdb_lock *sync.Mutex = new(sync.Mutex)

func StartSsdbClient(host string, port int) *SsdbClient {

	conn, err := ssdb.Connect(host, port)
	if err != nil {
		panic("Failed to connect to ssdb")
	} else {
		log.Info("ssdb client started")
	}
	client := new(SsdbClient)
	client.Client = conn
	client.host = host
	client.port = port
	go ping_pong_ssdb(client, 240)
	return client
}

func ping_pong_ssdb(client *SsdbClient, interval int) {
	c := time.Tick(time.Duration(interval) * time.Second)
	for _ = range c {
		g_ssdb_lock.Lock()
		//(*client.Conn).Do("PING")
		g_ssdb_lock.Unlock()
		log.Debug("sent PING to ssdb...")
	}
}
func (client *SsdbClient) Reconnect() {

	conn, err := ssdb.Connect(client.host, client.port)
	if err != nil {
		panic("Failed to connect to ssdb at port 6379")
	} else {
		log.Info("ssdb client reconnected")
	}
	client.Client = conn
}
func (client *SsdbClient) do(args ...interface{}) []string {

	g_ssdb_lock.Lock()
	defer g_ssdb_lock.Unlock()

	return client.doWithNoLock(args...)
}

func (client *SsdbClient) doWithNoLock(args ...interface{}) []string {

	ret, err := (*client.Client).Do(args...)
	if err != nil {
		if err.Error() == "use of closed network connection" {
			client.Reconnect()
			client.doWithNoLock(args...)
			return nil
		} else {
			panic(err)
		}
	}
	log.Debugf("send cmds to ssdb, key(some key), val(some bytes), returned=%s", ret)
	return ret
}

func (client *SsdbClient) FetchDeviceSubTopics(app_key string) []string {
	var topics []string
	key := "Topic-Id:" + app_key + "/"

	ret := client.do("hlist", key, "", -1)

	log.Debugf("fetch  sub topics from ssdb, key=%s, val(some bytes), returned=%s", key, ret)
	if ret[0] != "ok" {
		return nil
	}

	for i := 1; i < len(ret); i++ {
		if strings.HasPrefix(ret[i], key) {
			topics = append(topics, ret[i])
		}
	}
	return topics
}

func (client *SsdbClient) FetchOfflineDeviceList(topic string, recvType uint16) []string {

	var key []string
	var kv map[string]string = make(map[string]string)
	var ret []string
	var list []string

	log.Debugf("Topic: %s,RecvType: %d", topic, recvType)
	switch recvType {
	case 1:
		ret = client.do("get", "Bind-Id:"+topic)
		if ret[0] != "ok" {
			return nil
		} else {
			status := client.do("multi_get", "Reg-Id:"+topic, "Conn-Id:"+ret[1])
			//1:do error 2:app unregister 3:device offline
			if status[0] != "ok" || status[2] != "1" || status[4] != "0" {
				return nil
			}
			return ret[1:]
		}
	case 2:
		key = make([]string, 1)
		key[0] = "Topic-Id:" + topic
	case 4:
		key = client.FetchDeviceSubTopics(topic)
		if key == nil {
			return nil
		}
	default:
		return nil
	}
	log.Debugf("keys: %v", key)
	for _, v := range key {
		ret = client.do("hscan", v, "", "", -1)
		if ret[0] != "ok" || len(ret)%2 != 1 {
			return nil
		}
		for i := 1; i < len(ret); i += 2 {
			status := client.do("multi_get", "Reg-Id:"+ret[i], "Conn-Id:"+ret[i+1])
			//1:do error 2:app unregister 3:device offline
			if status[0] != "ok" || status[2] != "1" || status[4] != "0" {
				continue
			}
			kv[ret[i]] = ret[i+1]
		}
	}

	for _, v := range kv {
		list = append(list, v)
	}
	log.Debugf("return offline device list: %v", list)
	return list
}

func (client *SsdbClient) CacheExceptionHandler(e interface{}) {
	print(e)
}

func (client *SsdbClient) Cache(offline_client string, mqtt_msg *MqttMessage, exception_handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			exception_handler(err)
		}
	}()

	recvType := mqtt_msg.MessageId
	switch recvType {
	case 1:
		//topic=>appname/appkey ,where appkey is IMEI actually
		topic := mqtt_msg.Topic
		delimiter := strings.Index(topic, "/")
		appname, _ := topic[0:delimiter], topic[delimiter+1:]

		//get timestamp
		time_stamp := get_current_timestamp()
		life_time := fmt.Sprintf("%d", time_stamp)

		//store by hset: appname/offline_client=>{time_stamp/life_time:cache message}
		key := appname + "/" + offline_client
		ret := client.do("hset", key, time_stamp+"/"+life_time, mqtt_msg.Payload)
		log.Debugf("cache p2p message to ssdb, key=%s, val(message Payload), returned=%s", key, ret)

	case 2:
		// key := "Topic-Id:" + app_name + "/" + app_val
		// ret := client.do("hset", key, app_name+"/"+app_key, device_id)
		// log.Debugf("UpdateDeviceSubTopics to ssdb, key=%s, val(some bytes), returned=%s", key, ret)
	case 4:
	default:
		panic("Unknown recv type " + string(recvType))
	}

}

func get_current_timestamp() (time_stamp string) {
	var cur int64 = time.Now().Unix()
	time_stamp = fmt.Sprintf("%d", cur)
	fmt.Println("current time stamp:", time_stamp)
	return
}

func (client *SsdbClient) FetchOnlineDeviceList(topic string, recvType uint16) []string {

	var key []string
	var kv map[string]string = make(map[string]string)
	var ret []string
	var list []string

	log.Debugf("Topic: %s,RecvType: %d", topic, recvType)
	switch recvType {
	case 1:
		ret = client.do("get", "Bind-Id:"+topic)
		if ret[0] != "ok" {
			return nil
		} else {
			status := client.do("multi_get", "Reg-Id:"+topic, "Conn-Id:"+ret[1])
			//1:do error 2:app unregister 3:device offline
			if status[0] != "ok" || status[2] != "1" || status[4] != "1" {
				return nil
			}
			return ret[1:]
		}
	case 2:
		key = make([]string, 1)
		key[0] = "Topic-Id:" + topic
	case 4:
		key = client.FetchDeviceSubTopics(topic)
		if key == nil {
			return nil
		}
	default:
		return nil
	}
	log.Debugf("keys: %v", key)
	for _, v := range key {
		ret = client.do("hscan", v, "", "", -1)
		if ret[0] != "ok" || len(ret)%2 != 1 {
			return nil
		}
		for i := 1; i < len(ret); i += 2 {
			status := client.do("multi_get", "Reg-Id:"+ret[i], "Conn-Id:"+ret[i+1])
			//1:do error 2:app unregister 3:device offline
			if status[0] != "ok" || status[2] != "1" || status[4] != "1" {
				continue
			}
			kv[ret[i]] = ret[i+1]
		}
	}

	for _, v := range kv {
		list = append(list, v)
	}
	log.Debugf("return online device list: %v", list)
	return list
}

func (client *SsdbClient) UpdateDeviceOnline(device_id, login_ip string) {
	key := "Conn-Id:" + device_id

	ret := client.do("set", key, "1")

	log.Debugf("UpdateDeviceOnline to ssdb, key=%s, val(some bytes), returned=%s", key, ret)

}
func (client *SsdbClient) UpdateDeviceOffline(device_id string) {

	key := "Conn-Id:" + device_id
	ret := client.do("set", key, "0")

	log.Debugf("UpdateDeviceOffline to ssdb, key=%s, val(some bytes), returned=%s", key, ret)

}
func (client *SsdbClient) UpdateDeviceSubTopics(device_id, app_name, app_key, app_val string) {
	key := "Topic-Id:" + app_name + "/" + app_val
	ret := client.do("hset", key, app_name+"/"+app_key, device_id)
	log.Debugf("UpdateDeviceSubTopics to ssdb, key=%s, val(some bytes), returned=%s", key, ret)
}

func (client *SsdbClient) UpdateDeviceUnsubTopics(device_id, app_name, app_key, app_val string) {
	key := "Topic-Id:" + app_name + "/" + app_val
	ret := client.do("hdel", key, app_name+"/"+app_key)
	log.Debugf("UpdateDeviceUnsubTopics to ssdb, key=%s, val(some bytes), returned=%s", key, ret)
}

func (client *SsdbClient) UpdateDeviceRegisted(device_id, app_name, app_key, app_val string) {
	key := "Reg-Id:" + app_name + "/" + app_key

	ret := client.do("set", key, "1")
	if ret[0] != "ok" {
		return
	}
	key = "Bind-Id:" + app_name + "/" + app_key
	ret = client.do("get", key)
	if ret[0] != "ok" {
		client.do("set", key, device_id)
	}

	log.Debugf("UpdateDeviceRegisted to ssdb, key=%s, val(some bytes), returned=%s", key, ret)

}

func (client *SsdbClient) UpdateDeviceUnregisted(device_id, app_name, app_key, app_val string) {
	key := "Reg-Id:" + app_name + "/" + app_key

	ret := client.do("set", key, "0")

	log.Debugf("UpdateDeviceUnregisted to ssdb, key=%s, val(some bytes), returned=%s", key, ret)
}
