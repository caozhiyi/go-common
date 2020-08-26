package timer

import (
	"time"
)

var defaultTimer *Timer

func init() {
	defaultTimer = NewTimer()
}

// AddTask 添加一个定时器任务
// duration： 执行时间间隔，单位为纳秒
// runTimes： 执行次数，<=0时没有次数限制
// function： 执行函数
// params：   函数参数
func AddTask(duration time.Duration, runTimes int, function interface{}, params ...interface{}) error {
	return defaultTimer.AddTask(duration, runTimes, function, params...)
}

// RemoveTask 移除一个定时器任务
func RemoveTask(function interface{}) error {
	return defaultTimer.RemoveTask(function)
}

// RemoveAll 移除所有定时器
func RemoveAll() {
	defaultTimer.RemoveAll()
}

// Wait 等待所有定时器协程结束
func Wait() {
	defaultTimer.Wait()
}
