package docker

import (
	"fmt"
	"testing"
)

func TestGetAllDockerImages(t *testing.T) {
	fmt.Println(GetAllDockerImages())
}
