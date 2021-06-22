package solution

import (
	"fmt"
	"sync"
)

/*
	编写一个程序，开启 3 个线程 A,B,C，这三个线程的输出分别为 A、B、C，
	每个线程将自己的输出在屏幕上打印 10 遍，要求输出的结果必须按顺序显示。
	如：ABCABCABC....
*/

const (
	N = 2  // 需要打印的字母个数
	M = 25 // 打印两次
)

var mainWg sync.WaitGroup

func orderPrint() {

	firstRead := make(chan struct{})
	lastWrite := make(chan struct{})
	read := firstRead

	for i := 0; i < N; i++ {
		mainWg.Add(1)
		s := string('A' + i)
		write := make(chan struct{})
		if i == N-1 {
			write = lastWrite
		}
		go echo(s, read, write, &mainWg)
		read = write // 下一个读取的通道是当前通道的输入
	}

	// 启动打印
	for i := 0; i < M; i++ {
		firstRead <- struct{}{} // 向第一个当中写入
		<-lastWrite             // 从最后一个当中读取
	}
	close(firstRead)

	mainWg.Wait()
}

func echo(s string, read chan struct{}, write chan struct{}, wg *sync.WaitGroup) {
	for _ = range read {
		fmt.Println(s)
		write <- struct{}{}
	}
	close(write)
	wg.Done()
}
