package mqtt

import (
	"testing"
)

var g_ssdb_client *SsdbClient = StartSsdbClient("127.0.0.1",8888)

func TestSsdbStoreDevInfo(t *testing.T) {

	key := "IMEI:1234567890"
	val := "deviceinfo"
	
	g_ssdb_client.SetDevInfo(key,val)
	
	ret := g_ssdb_client.GetDevInfo(key)
	if ret != val {
		t.Fatal(ret)
	}
	
	g_ssdb_client.DelDevInfo(key)
	ret = g_ssdb_client.GetDevInfo(key)
	if ret != "" {
		t.Fatal(ret)
	}
}

