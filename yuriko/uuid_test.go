package yuriko

import (
	"sort"
	"sync"
	"testing"
)

func TestUUID(t *testing.T) {
	snowFlake, err := NewSnowFlake(12, 25, 1617922214000)
	idNum := 1000000
	arr := make([]int64, idNum)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < idNum; i++ {
		wg.Add(1)
		go func(i int) {
			id := snowFlake.generateUUID()
			arr[i] = id
			wg.Done()
		}(i)
	}
	wg.Wait()
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})
	t.Logf("sort done, min: %v, max: %v", arr[0], arr[len(arr)-1])
	for i := 1; i < len(arr); i++ {
		if arr[i-1] == arr[i] {
			t.Fatalf("Duplicate uuid found: %d", arr[i])
		}
	}
}
