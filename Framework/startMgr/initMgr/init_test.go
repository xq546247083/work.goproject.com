package initMgr

import (
	"fmt"
	"testing"
)

func TestRegister(t *testing.T) {
	Register("first", first, true)
	Register("second", second, true)
	Register("third", third, true)
	Register("fourth", fourth, true)
}

func TestCallOne(t *testing.T) {
	name := "first"
	if err := CallOne(name); err != nil {
		t.Errorf("there should be no error, but now it has:%s", err)
	}
}

func TestCallAny(t *testing.T) {
	errList := CallAny("second", "third")
	if len(errList) != 1 {
		t.Errorf("there should be 1 error, but now:%d", len(errList))
	}
}

func TestCallAll(t *testing.T) {
	errList := CallAll()
	if len(errList) != 2 {
		t.Errorf("there should be 1 error, but now:%d", len(errList))
	}
}

func first() error {
	fmt.Println("first")
	return nil
}

func second() error {
	fmt.Println("second")
	return fmt.Errorf("the second error")
}

func third() error {
	fmt.Println("third")
	return nil
}

func fourth() error {
	fmt.Println("fourth")
	return fmt.Errorf("the fourth error")
}
