package test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type User struct {
	itemMap sync.Map
}

func f1(rangeMap sync.Map) {
	rangeMap.Range(func(key, value interface{}) bool {
		v := value.(int)
		if v < 0 {
			fmt.Errorf("error")
		}
		fmt.Printf("key->%v,value->%v", key, value)
		return true
	})
}

func TestSyncMap(t *testing.T) {
	user := &User{}
	map1 := user.itemMap
	go func() {
		for i := 0; i < 900000; i++ {
			user.itemMap.Store(i, i)
			time.Sleep(1 * time.Millisecond)
		}
	}()
	go func() {
		for i := 0; i < 900000; i++ {
			//map1.Range(func(key, value interface{}) bool {
			//	fmt.Printf("key->%v,value->%v", key, value)
			//	return true
			//})
			//函数传参为值传递（sync.Map的副本），等同如下（值类型深拷贝）
			f1(map1)
			time.Sleep(1 * time.Millisecond)
		}
	}()

	time.Sleep(5 * time.Second)
}
