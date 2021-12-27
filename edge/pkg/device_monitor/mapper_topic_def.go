package devicemonitor

//device monitor 发起查询本机上的设备
func TopicInquireDeviceName(device_name string) string {
	return "$kp/events/edge/dm/inqurie_device"
}

//https://www.jianshu.com/p/f3f04876c968
func TopicDeviceDataUpdate(device_name string) string {
	//$hw/events/device/+/twin/update 	mapper 	edgecore 	通知设备属性的值更新 	是
	return "$hw/events/device/" + device_name + "/twin/update"
}
