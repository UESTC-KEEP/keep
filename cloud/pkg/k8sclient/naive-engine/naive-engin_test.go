package naive_engine

import (
	"fmt"
	"testing"
)

func TestListPods(t *testing.T) {
	fmt.Println(NewNaiveEngine().ListPods("default"))
}
