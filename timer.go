package timer

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// task 定时执行的任务
type task struct {
	function   interface{}
	params     []interface{}
	cancelChan chan interface{}
	duration   time.Duration
	timer      *Timer
	runTimes   int
}

func newTask(timer *Timer, duration time.Duration, runTimes int, function interface{}, params ...interface{}) (*task, error) {
	typ := reflect.TypeOf(function)
	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("timer: wrong type %q set in timer", strconv.Itoa(int(typ.Kind())))
	}
	if len(params) != typ.NumIn() {
		return nil, fmt.Errorf("timer: %d params set but function need %d", len(params), typ.NumIn())
	}

	return &task{
		function:   function,
		params:     params,
		cancelChan: make(chan interface{}),
		duration:   duration,
		timer:      timer,
		runTimes:   runTimes,
	}, nil
}

type Timer struct {
	mutex   sync.Mutex
	wg      sync.WaitGroup
	taskMap map[string]*task
}

// NewTimer 创建定时器
func NewTimer() *Timer {
	return &Timer{
		taskMap: make(map[string]*task),
	}
}

// AddTask 添加一个定时器任务
// duration： 执行时间间隔，单位为毫秒
// runTimes： 执行次数，<=0时没有次数限制
// function： 执行函数
// params：   函数参数
func (t *Timer) AddTask(duration time.Duration, runTimes int, function interface{}, params ...interface{}) error {
	tk, err := newTask(t, duration, runTimes, function, params...)
	if err != nil {
		return err
	}

	funcName := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()

	t.mutex.Lock()
	defer t.mutex.Unlock()

	oldTk, ok := t.taskMap[funcName]
	if ok {
		close(oldTk.cancelChan)
	}
	t.taskMap[funcName] = tk

	// 启动定时器协程
	t.wg.Add(1)
	go runTask(tk)

	return nil
}

// RemoveTask 移除一个定时器任务
func (t *Timer) RemoveTask(function interface{}) error {
	typ := reflect.TypeOf(function)
	if typ.Kind() != reflect.Func {
		return fmt.Errorf("timer: wrong function type %q", strconv.Itoa(int(typ.Kind())))
	}
	funcName := runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()

	t.mutex.Lock()
	defer t.mutex.Unlock()

	oldTk, ok := t.taskMap[funcName]
	if ok {
		close(oldTk.cancelChan)
		delete(t.taskMap, funcName)
	}
	return nil
}

// RemoveAll 移除所有定时器
func (t *Timer) RemoveAll() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, value := range t.taskMap {
		close(value.cancelChan)
	}
	t.taskMap = make(map[string]*task)
}

// Wait 等待所有定时器协程结束
func (t *Timer) Wait() {
	t.wg.Wait()
}

func runTask(t *task) {
	defer func() {
		// 增加recover，防止process执行panic，导致整个服务退出.
		if r := recover(); r != nil {
			const size = 2048
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Fprintf(os.Stderr, "[ERROR] timer task process panic. %v\n%s", r, string(buf))
		}
	}()

	ticker := time.NewTicker(t.duration)
	defer ticker.Stop()
	defer t.timer.wg.Done()

	in := make([]reflect.Value, len(t.params))
	f := reflect.ValueOf(t.function)

	check := false
	if t.runTimes > 0 {
		check = true
	}

	for k, param := range t.params {
		in[k] = reflect.ValueOf(param)
	}
	for {
		select {
		case <-ticker.C:
			f.Call(in)
			if check {
				t.runTimes--
				if t.runTimes <= 0 {
					return
				}
			}
		case <-t.cancelChan:
			return
		}
	}
}
