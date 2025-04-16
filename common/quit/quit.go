package quit

import (
	"fmt"
	"os"
	"os/signal"
)

/**********************************************************************************************************************
优雅退出步骤：
1.使用GetQuitEvent()获取私有调用结构quitEvent
2.将退出执行的方法，注册到RegisterQuitFunc中
3.调用WaitQuitSignal()监听系统退出信号，一旦收到退出信号，执行注册方法
***********************************************************************************************************************/

type Event struct{}

var rhs = make([]func(), 0)

func GetQuitEvent() *Event {
	return new(Event)
}

func (event *Event) RegisterFunc(handles ...func()) {
	rhs = handles
}

func WaitSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	s := <-sig
	fmt.Println("接收到系统关闭信号，开始做退出程序操作", s.String())
	for _, rh := range rhs {
		rh()
	}
	fmt.Println("程序优雅退出")
}
