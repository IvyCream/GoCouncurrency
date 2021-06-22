package examples

// description
/*
	高并发的web服务器中，要限制IP的频繁访问
	模拟100个IP同时并发访问服务器，
	每个IP重复访问1000次，三分钟之内只能访问一次
	修改以下代码，要求能成功输出success: 100
*/

import (
	"fmt"
	"time"
)

func ErrExample() {
	var flag int32 = 1
	go func() {
		count := 1
		for flag == 1 {
			fmt.Println(flag)
			count++
		}
		println("thread1 end") //这个循环永远也不会结束，为什么？
	}()

	go func() {
		for {
			flag = 0
			fmt.Println(flag)
		}
	}()
	time.Sleep(time.Hour)
	return
}
