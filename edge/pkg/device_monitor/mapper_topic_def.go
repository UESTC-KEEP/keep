package devicemonitor

//https://www.jianshu.com/p/f3f04876c968

func TopicDeviceDataUpdate(device_name string) string {
	//$hw/events/device/+/twin/update 	mapper 	edgecore 	通知设备属性的值更新 	是
	return "$hw/events/device/" + device_name + "/twin/update"
}
