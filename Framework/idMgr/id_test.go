package idMgr

import (
	"fmt"
	"sync"
	"testing"
)

var (
	wg       sync.WaitGroup
	mapMutex sync.Mutex
)

func TestGenerateNewId(t *testing.T) {
	dataMap := make(map[int64]bool)
	serverGroupId := int64(32768)

	if _, err := GenerateNewId(serverGroupId); err == nil {
		t.Errorf("there should be err, but now not.")
	}

	serverGroupId = 32767
	count := 1048576

	for num := 0; num < 10; num++ {
		go func() {
			wg.Add(1)
			defer wg.Done()

			for i := 0; i < count; i++ {
				if id, err := GenerateNewId(serverGroupId); err != nil {
					t.Errorf("there should be no error, but now it has")
				} else {
					mapMutex.Lock()
					if _, exists := dataMap[id]; exists {
						t.Errorf("there should be not duplicate, but now it does.%d", id)
					} else {
						dataMap[id] = true
					}
					mapMutex.Unlock()
					fmt.Println(id)
				}
			}
		}()
	}

	wg.Wait()
}
