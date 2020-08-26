package timer

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer1(t *testing.T) {
	fn1 := func() {
		fmt.Println("no params task 1")
	}

	fn2 := func() {
		fmt.Println("no params task 2")
	}

	if err := AddTask(2*time.Second, 3, fn1); err != nil {
		t.Errorf("add timer task failed.")
	}
	if err := AddTask(3*time.Second, 4, fn2); err != nil {
		t.Errorf("add timer task failed.")
	}
	Wait()
}

func TestTimer2(t *testing.T) {
	fn1 := func(param int) {
		fmt.Println("params task 3", param)
	}

	fn2 := func(param1, param2 int) {
		fmt.Println("no params task 4", param1, param2)
	}

	if err := AddTask(1*time.Second, 3, fn1, 100); err != nil {
		t.Errorf("add timer task failed.")
	}
	if err := AddTask(2*time.Second, 4, fn2, 100, 200); err != nil {
		t.Errorf("add timer task failed.")
	}
	Wait()
}

func TestTimer3(t *testing.T) {
	fn1 := func(param int) {
		fmt.Println("params task 5", param)
	}

	fn2 := func() {
		RemoveTask(fn1)
	}

	if err := AddTask(2*time.Second, 0, fn1, 100); err != nil {
		t.Errorf("add timer task failed.")
	}
	if err := AddTask(1*time.Second, 1, fn2); err != nil {
		t.Errorf("add timer task failed.")
	}
	Wait()
}

func TestTimer4(t *testing.T) {
	fn1 := func(param int) {
		fmt.Println("params task 6", param)
	}

	fn2 := func() {
		RemoveAll()
	}

	if err := AddTask(2*time.Second, 0, fn1, 100); err != nil {
		t.Errorf("add timer task failed.")
	}
	if err := AddTask(1*time.Second, 1, fn2); err != nil {
		t.Errorf("add timer task failed.")
	}
	Wait()
}
