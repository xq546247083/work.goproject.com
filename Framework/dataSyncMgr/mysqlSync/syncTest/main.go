package main

import (
	"fmt"
	"work.goproject.com/goutil/mathUtil"
	"work.goproject.com/goutil/stringUtil"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func init() {
	wg.Add(1)
}

func main() {
	playerMgr := newPlayerMgr()

	// insert
	go func() {
		for {
			id := stringUtil.GetNewGUID()
			name := fmt.Sprintf("Hero_%s", id)
			obj := newPlayer(id, name)
			playerMgr.insert(obj)

			insert(obj)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// update
	go func() {
		for {
			obj := playerMgr.randomSelect()
			if obj == nil {
				continue
			}
			suffix := mathUtil.GetRandInt(1000)
			newName := fmt.Sprintf("Hero_%d", suffix)
			obj.resetName(newName)

			update(obj)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// delete
	go func() {
		for {
			obj := playerMgr.randomSelect()
			if obj == nil {
				continue
			}
			playerMgr.delete(obj)

			clear(obj)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	//	errorFile
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			id := stringUtil.GetNewGUID()
			name := fmt.Sprintf("Hero_%s%s", id, id)
			obj := newPlayer(id, name)
			playerMgr.insert(obj)
			print("errorFile")

			insert(obj)
		}

	}()

	wg.Wait()
}
