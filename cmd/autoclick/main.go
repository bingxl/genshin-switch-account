package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var (
	active bool // 是否处于激活状态， 活跃窗口在games 中，按下 alt+f 时对 active 取反
	window string
	mu     sync.Mutex
)

// 窗口进程名：激活后重复的动作
var games = map[string]func(){
	"YuanShen.exe": func() {
		// 2600 x 1980 下
		// 原神对话最后一个选项
		xPos := 1817
		yPos := 1235

		robotgo.MoveClick(xPos, yPos)
	},

	// 鸣潮，激活后循环按 F
	"wave": func() {
		robotgo.KeyTap("f")
	},
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go listenHotKey(&wg)

	wg.Add(1)
	go clickLoop(&wg)

	wg.Wait()
}

// 检查是否已有目标程序启动， 若有则将 window变量设为为启动的程序名
func checkWindow() bool {
	for k := range games {
		if _, err := robotgo.FindIds(k); err == nil {
			window = k
			return true
		}
	}
	return false
}

func listenHotKey(wg *sync.WaitGroup) {
	defer wg.Done()
	defer hook.End()

	hook.Register(hook.KeyDown, []string{"alt", "f"}, func(e hook.Event) {
		// Alt+F
		mu.Lock()
		defer mu.Unlock()

		active = !active
		if !checkWindow() {
			active = false
			fmt.Println("没有启动的目标程序")
		}
		fmt.Printf("Alt+F pressed, active: %t\n", active)

		// hook.End()
	})

	s := hook.Start()
	<-hook.Process(s)
}

func clickLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// mu.Lock()
		if active {
			ru, ok := games[window]
			if ok {
				ru()
				// fmt.Println("handle")
			}
			time.Sleep(time.Duration(500+randOffset()) * time.Millisecond)
		} else {
			time.Sleep(1000 * time.Millisecond)
		}
		// mu.Unlock()
	}
}

func randOffset() int64 {
	return int64((rand.Intn(40) - 20)) // 生成-20到20之间的随机数
}
