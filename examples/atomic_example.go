package examples

import (
	"sync"
	"sync/atomic"
)

type UserInfo struct {
	Name string
}

var lock sync.Mutex

// 操作的原子性
func getInstanceErr(instance interface{}) interface{} {
	/*
		? 为什么instance的修改在lock 临界区中,依然会发生DataRace？
		* 因为锁内对instance值的修改步骤并不是原子操作，这个赋值步骤可能会有好几步命令，例如：
			1. 先 new 一个UserInfo
			2. 设置 Name = "test"
			3. UserInfo 赋值给instance
		如果这其中发生了乱序，例如，1,3,2,这会导致instance != nil
		此时，如果其他线程进入 getInstanceErr，会导致instance的状态还没有修改完 (处于中间状态)
		而 instance != nil,跳过first-check，直接返回的情况
	*/
	if instance == nil { // first-check
		lock.Lock()
		defer lock.Unlock()
		if instance == nil { // second-check
			instance = &UserInfo{
				Name: "test",
			}
		}
	}
	return instance
}

var flag uint32

func getInstanceFix(instance interface{}) interface{} {
	if atomic.LoadUint32(&flag) != 1 {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &UserInfo{
				Name: "test",
			}
			atomic.StoreUint32(&flag, 1)
		}
	}
	return instance
}
