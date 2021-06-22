package solution

import (
	"fmt"
	"sync"
	"time"
)

// description
/*
	高并发的web服务器中，要限制IP的频繁访问
	模拟100个IP同时并发访问服务器，
	每个IP重复访问1000次，三分钟之内只能访问一次
	修改以下代码，要求能成功输出success: 100
*/

type Ban struct {
	visitIPs map[string]struct{}
}

var lock sync.Mutex

func NewBan() *Ban {
	return &Ban{visitIPs: make(map[string]struct{})}
}

func (o *Ban) visit(ip string) bool {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := o.visitIPs[ip]; ok {
		return true
	}
	o.visitIPs[ip] = struct{}{}
	go o.expire(ip)
	return false
}

func (o *Ban) expire(ip string) {
	// 不建议采取这种方法，会导致 n个空余的goroutine在等待
	time.Sleep(time.Minute * 3)
	lock.Lock()
	visitIPs := o.visitIPs
	delete(visitIPs, ip)
	o.visitIPs = visitIPs
	lock.Unlock()
}

func Example() {
	success := 0
	ban := NewBan()
	for i := 0; i < 1000; i++ { // 重复访问1000 次
		for j := 0; j < 100; j++ { // 共有 100个IP
			go func(ij int) {
				ip := fmt.Sprintf("192.168.1.%d", ij)
				if !ban.visit(ip) {
					// 如果当前IP 还没有visit过，success+1
					// 注意对ban的访问是并发的，所以需要对ban的读写加锁
					success++
				}
			}(j)
		}
	}
	fmt.Println("success:", success)
}
