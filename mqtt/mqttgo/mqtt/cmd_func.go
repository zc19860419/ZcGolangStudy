package mqtt

import (
	"encoding/json"
	log "mqtt/seelog"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

const (
	SEND_WILL = uint8(iota)
	DONT_SEND_WILL
)

type SubInfo struct {
	TopicType int
	AppName   string
	AppKey    string
	AppValue  string
}

// Handle CONNECT
func HandleConnect(mqtt *Mqtt, conn *net.Conn, client **ClientRep) {
	//mqtt.Show()
	client_id := mqtt.ClientId

	log.Debugf("Handling CONNECT, client id:(%s)", client_id)

	if len(client_id) > 23 {
		log.Debugf("client id(%s) is longer than 23, will send IDENTIFIER_REJECTED", client_id)
		SendConnack(IDENTIFIER_REJECTED, conn, nil)
		return
	}

	if mqtt.ProtocolName != "MQIsdp" || mqtt.ProtocolVersion != 3 {
		log.Debugf("ProtocolName(%s) and/or version(%d) not supported, will send UNACCEPTABLE_PROTOCOL_VERSION",
			mqtt.ProtocolName, mqtt.ProtocolVersion)
		SendConnack(UNACCEPTABLE_PROTOCOL_VERSION, conn, nil)
		return
	}

	G_conn_clients_lock.Lock()
	client_rep, existed := G_conn_clients[client_id]
	if existed {
		log.Debugf("%s existed, will close old connection", client_id)
		ForceDisconnect(client_rep, DONT_SEND_WILL)

	} else {
		log.Debugf("Appears to be new client, will create ClientRep")
	}

	client_rep = CreateClientRep(client_id, conn, mqtt)
	G_conn_clients[client_id] = client_rep
	G_conn_clients_lock.Unlock()

	*client = client_rep

	if client_rep.Mqtt.KeepAliveTimer != 0 {
		go CheckTimeout(client_rep)
		log.Debugf("Timeout checker go-routine started")
	}

	if !client_rep.Mqtt.ConnectFlags.CleanSession {
		//Update device info
		remoteAddr := (*conn).RemoteAddr().String()
		remoteIp := strings.Split(remoteAddr, ":")[0]
		G_db_client.UpdateDeviceOnline(client_id, remoteIp)
		// deliver flying messages
		DeliverOnConnection(client_id)

	}
	SendConnack(ACCEPTED, conn, client_rep.WriteLock)
	log.Debugf("New client is all set and CONNACK is sent")
}

func SendConnack(rc uint8, conn *net.Conn, lock *sync.Mutex) {
	resp := CreateMqtt(CONNACK)
	resp.ReturnCode = rc

	bytes, _ := Encode(resp)
	MqttSendToClient(bytes, conn, lock)
}

/* Handle PUBLISH*/
// FIXME: support qos = 2
func HandlePublish(mqtt *Mqtt, conn *net.Conn, client **ClientRep) {
	if *client == nil {
		panic("client_resp is nil, that means we don't have ClientRep for this client sending PUBLISH")
		return
	}

	client_id := (*client).Mqtt.ClientId
	client_rep := *client
	client_rep.UpdateLastTime()
	topic := mqtt.TopicName
	payload := string(mqtt.Data)
	qos := mqtt.FixedHeader.QosLevel
	retain := mqtt.FixedHeader.Retain
	message_id := mqtt.MessageId
	timestamp := time.Now().Unix()
	log.Debugf("Handling PUBLISH, client_id: %s, topic:(%s), payload:(%s), qos=%d, retain=%t, message_id=%d",
		client_id, topic, payload, qos, retain, message_id)

	// Create new MQTT message
	mqtt_msg := CreateMqttMessage(topic, payload, client_id, qos, message_id, timestamp, retain)
	msg_internal_id := mqtt_msg.InternalId
	log.Debugf("Created new MQTT message, internal id:(", msg_internal_id, ")")

	PublishMessage(client_rep, mqtt_msg)

	// Send PUBACK if QOS is 1
	if qos == 1 {
		SendPuback(message_id, conn, client_rep.WriteLock)
		log.Debugf("PUBACK sent to client(%s)", client_id)
	}
}

func SendPuback(msg_id uint16, conn *net.Conn, lock *sync.Mutex) {
	resp := CreateMqtt(PUBACK)
	resp.MessageId = msg_id
	bytes, _ := Encode(resp)
	MqttSendToClient(bytes, conn, lock)

}

/* Handle SUBSCRIBE */

func HandleSubscribe(mqtt *Mqtt, conn *net.Conn, client **ClientRep) {
	if *client == nil {
		panic("client_resp is nil, that means we don't have ClientRep for this client sending SUBSCRIBE")
		return
	}

	client_id := (*client).Mqtt.ClientId
	log.Debugf("Handling SUBSCRIBE, client_id: %s", client_id)
	client_rep := *client
	client_rep.UpdateLastTime()
	for i := 0; i < len(mqtt.Topics); i++ {
		topic := mqtt.Topics[i]
		qos := mqtt.Topics_qos[i]
		log.Debugf("will subscribe client(%s) to topic(%s) with qos=%d", client_id, topic, qos)
		if !client_rep.Mqtt.ConnectFlags.CleanSession {
			//Update Subscriptions to SQL
			log.Debugf("storeSubTopics")
			storeSubTopics(client_id, topic)
		}
		log.Debugf("finding retained message for (%s)", topic)
		//FIXME Send Remain Message
	}

	log.Debugf("Subscriptions are all processed, will send SUBACK")
	SendSuback(mqtt.MessageId, mqtt.Topics_qos, conn, client_rep.WriteLock)
	showSubscriptions()
}

func SendSuback(msg_id uint16, qos_list []uint8, conn *net.Conn, lock *sync.Mutex) {
	resp := CreateMqtt(SUBACK)
	resp.MessageId = msg_id
	resp.Topics_qos = qos_list

	bytes, _ := Encode(resp)
	MqttSendToClient(bytes, conn, lock)
}

/* Handle UNSUBSCRIBE */

func HandleUnsubscribe(mqtt *Mqtt, conn *net.Conn, client **ClientRep) {
	if *client == nil {
		panic("client_resp is nil, that means we don't have ClientRep for this client sending UNSUBSCRIBE")
		return
	}

	client_id := (*client).Mqtt.ClientId
	log.Debugf("Handling UNSUBSCRIBE, client_id: %s", client_id)
	client_rep := *client
	client_rep.UpdateLastTime()

	for i := 0; i < len(mqtt.Topics); i++ {
		topic := mqtt.Topics[i]
		log.Debugf("unsubscribing client(%s) from topic(%s)", client_id, topic)
		delete(client_rep.Subscriptions, topic)
		//Clean Subscriptions from SQL
		cleanSubTopics(client_id, topic)
	}
	log.Debugf("unsubscriptions are all processed, will send UNSUBACK")
	SendUnsuback(mqtt.MessageId, conn, client_rep.WriteLock)
	showSubscriptions()
}

func SendUnsuback(msg_id uint16, conn *net.Conn, lock *sync.Mutex) {
	resp := CreateMqtt(UNSUBACK)
	resp.MessageId = msg_id
	bytes, _ := Encode(resp)
	MqttSendToClient(bytes, conn, lock)
}

/* Handle PINGREQ */

func HandlePingreq(mqtt *Mqtt, conn *net.Conn, client **ClientRep) {
	if *client == nil {
		panic("client_resp is nil, that means we don't have ClientRep for this client sending PINGREQ")
		return
	}

	client_id := (*client).Mqtt.ClientId
	log.Debugf("Handling PINGREQ, client_id: %s", client_id)
	client_rep := *client
	client_rep.UpdateLastTime()

	SendPingresp(conn, client_rep.WriteLock)
	log.Debugf("Sent PINGRESP, client_id: %s", client_id)
}

func SendPingresp(conn *net.Conn, lock *sync.Mutex) {
	resp := CreateMqtt(PINGRESP)
	bytes, _ := Encode(resp)
	MqttSendToClient(bytes, conn, lock)
}

/* Handle DISCONNECT */

func HandleDisconnect(mqtt *Mqtt, conn *net.Conn, client **ClientRep) {
	if *client == nil {
		panic("client_resp is nil, that means we don't have ClientRep for this client sending DISCONNECT")
		return
	}

	ForceDisconnect(*client, DONT_SEND_WILL)
}

/* Handle PUBACK */
func HandlePuback(mqtt *Mqtt, conn *net.Conn, client **ClientRep) {
	if *client == nil {
		panic("client_resp is nil, that means we don't have ClientRep for this client sending DISCONNECT")
		return
	}

	client_id := (*client).Mqtt.ClientId
	message_id := mqtt.MessageId
	log.Debugf("Handling PUBACK, client:(%s), message_id:(%d)", client_id, message_id)

	//FIXME Send Fly Message if needed
}

/* Helper functions */

// This is the main place to change if we need to use channel rather than lock
func MqttSendToClient(bytes []byte, conn *net.Conn, lock *sync.Mutex) {
	if lock != nil {
		lock.Lock()
		defer func() {
			lock.Unlock()
		}()
	}
	(*conn).Write(bytes)
}

/* Checking timeout */
func CheckTimeout(client *ClientRep) {
	defer func() {
		if r := recover(); r != nil {
			log.Debugf("got panic, will print stack")
			debug.PrintStack()
			panic(r)
		}
	}()

	interval := client.Mqtt.KeepAliveTimer
	client_id := client.ClientId
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	for {
		select {
		case <-ticker.C:
			now := time.Now().Unix()
			lastTimestamp := client.LastTime
			deadline := int64(float64(lastTimestamp) + float64(interval)*1.5)

			if deadline < now {
				ForceDisconnect(client, SEND_WILL)
				log.Debugf("client(%s) is timeout, kicked out",
					client_id)
			} else {
				log.Debugf("client(%s) will be kicked out in %d seconds",
					client_id,
					deadline-now)
			}
		case <-client.Shuttingdown:
			log.Debugf("client(%s) is being shutting down, stopped timeout checker", client_id)
			return
		}

	}
}

func ForceDisconnect(client *ClientRep, send_will uint8) {
	if client.Disconnected == true {
		return
	}

	client.Disconnected = true

	client_id := client.Mqtt.ClientId

	log.Debugf("Disconnecting client(%s), clean-session:%t", client_id, client.Mqtt.ConnectFlags.CleanSession)

	delete(G_conn_clients, client_id)

	if client.Mqtt.ConnectFlags.CleanSession {
		// remove her subscriptions

		// remove her flying messages
		log.Debugf("Removing all flying messages for (%s)", client_id)
	} else {
		G_db_client.UpdateDeviceOffline(client_id)
	}
	// FIXME: Send will if requested
	if send_will == SEND_WILL && client.Mqtt.ConnectFlags.WillFlag {
		will_topic := client.Mqtt.WillTopic
		will_payload := client.Mqtt.WillMessage
		will_qos := client.Mqtt.ConnectFlags.WillQos
		will_retain := client.Mqtt.ConnectFlags.WillRetain

		mqtt_msg := CreateMqttMessage(will_topic, will_payload, client_id, will_qos,
			0, // message id won't be used here
			time.Now().Unix(), will_retain)
		PublishMessage(client, mqtt_msg)

		log.Debugf("Sent will for %s, topic:(%s), payload:(%s)",
			client_id, will_topic, will_payload)
	}

	client.Shuttingdown <- 1
	log.Debugf("Sent 1 to shutdown channel")

	log.Debugf("Closing socket of %s", client_id)
	(*client.Conn).Close()
}

func PublishMessage(client *ClientRep, mqtt_msg *MqttMessage) {
	topic := mqtt_msg.Topic
	payload := mqtt_msg.Payload
	log.Debugf("Publishing job, topic(%s), payload(%s)", topic, payload)
	// Update global topic record

	if mqtt_msg.Retain {
		//FIXME Store Fly Retain Message
		log.Debugf("Set the message(%s) as the current retain content of topic:%s", payload, topic)
	}
	recvType := mqtt_msg.MessageId
	devices := G_db_client.FetchOnlineDeviceList(topic, recvType)
	// Dispatch delivering jobs
	for index := range devices {
		client, found := G_conn_clients[devices[index]]
		if found {
			go Deliver(client, mqtt_msg)
			log.Debugf("Started deliver job for %s", client.ClientId)
		}
	}
	log.Debugf("All delivering job dispatched")

	offline_devices := G_db_client.FetchOfflineDeviceList(topic, recvType)
	for _, offline_device := range offline_devices {
		go G_db_client.Cache(offline_device, mqtt_msg, G_db_client.CacheExceptionHandler)
		log.Debugf("Started cache job for %s", offline_device)
	}
	log.Debugf("All cache job done")

}

func DeliverOnConnection(client_id string) {
	log.Debugf("client(%s) just reconnected, delivering on the fly messages", client_id)
	//FIXME fetch remain message from redis

	//FIXME remove records in redis
	log.Debugf("client(%s), all flying messages put in pipeline, removed records in redis", client_id)
}

// Real heavy lifting jobs for delivering message
func DeliverMessage(client *ClientRep, msg *MqttMessage) {

	conn := client.Conn
	lock := client.WriteLock
	message_id := NextOutMessageIdForClient(client.ClientId)
	fly_msg := CreateFlyingMessage(client.ClientId, msg.InternalId, msg.Qos, PENDING_PUB, message_id)

	// FIXME: Add code to deal with failure
	resp := CreateMqtt(PUBLISH)
	resp.TopicName = msg.Topic
	if msg.Qos > 0 {
		resp.MessageId = message_id
	}
	resp.FixedHeader.QosLevel = msg.Qos
	resp.Data = []byte(msg.Payload)

	bytes, _ := Encode(resp)

	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	// FIXME: add write deatline
	(*conn).Write(bytes)
	log.Debugf("message sent by Write()")

	if msg.Qos == 1 {
		fly_msg.Status = PENDING_ACK
		//FIXME Store fly messge to redis
		log.Debugf("message(msg_id=%d) sent to client(%s), waiting for ACK, added to redis",
			message_id, client.ClientId)
	}
}

func Deliver(client *ClientRep, msg *MqttMessage) {
	defer func() {
		if r := recover(); r != nil {
			log.Debugf("got panic, will print stack")
			debug.PrintStack()
			panic(r)
		}
	}()

	log.Debugf("Delivering msg(internal_id=%d) to client(%s)", msg.InternalId, client.ClientId)

	DeliverMessage(client, msg)

	if msg.Qos > 0 {
		// Start retry
		go RetryDeliver(20, client, msg)
	}
}

func RetryDeliver(sleep uint64, client *ClientRep, msg *MqttMessage) {
	defer func() {
		if r := recover(); r != nil {
			log.Debugf("got panic, will print stack")
			debug.PrintStack()
			panic(r)
		}
	}()

	if sleep > 3600*4 {
		log.Debugf("too long retry delay(%s), abort retry deliver", sleep)
		return
	}

	time.Sleep(time.Duration(sleep) * time.Second)
	//FIXME retry send message here
}

func storeSubTopics(client_id string, topic string) {
	var s SubInfo
	err := json.Unmarshal([]byte(topic), &s)
	if err != nil {
		return
	}
	switch s.TopicType {
	case 0:
		G_db_client.UpdateDeviceRegisted(client_id, s.AppName, s.AppKey, s.AppValue)
	case 1:
		G_db_client.UpdateDeviceSubTopics(client_id, s.AppName, s.AppKey, s.AppValue)
	default:
	}
}
func cleanSubTopics(client_id string, topic string) {
	var s SubInfo
	err := json.Unmarshal([]byte(topic), &s)
	if err != nil {
		return
	}
	switch s.TopicType {
	case 0:
		G_db_client.UpdateDeviceUnregisted(client_id, s.AppName, s.AppKey, s.AppValue)
	case 1:
		G_db_client.UpdateDeviceUnsubTopics(client_id, s.AppName, s.AppKey, s.AppValue)
	default:
	}

}
func showSubscriptions() {
	// Disable for now
	return
}
