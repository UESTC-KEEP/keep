package docker

import (
	"fmt"
	"testing"
)

func TestGetAllDockerImages(t *testing.T) {
	fmt.Println(GetAllDockerImages())
}

func TestListContainerMetrics(t *testing.T) {
	metrics, err := ListContainerMetrics()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("ID\tName\tCPU%%\tMEM\tMEM%%\tNET\tDISK\n")
	for _, container := range *metrics {
		fmt.Printf("%s\t%s\t%f\t%f\t%f\trx:%f, tx: %f\tread:%f write:%f\n", container.ID[:12], container.Name, container.Cpu, container.Memory, container.MemoryPercentage, container.Net["recv"], container.Net["send"], container.Disk["read"], container.Disk["write"])
	}
}
