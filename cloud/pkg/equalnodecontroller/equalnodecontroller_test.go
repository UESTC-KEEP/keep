package equalnodecontroller

import (
	"keep/cloud/pkg/equalnodecontroller/controller"
	"testing"
)

func TestWorker(t *testing.T) {
	controller.StartEqualNodecontroller()
}
