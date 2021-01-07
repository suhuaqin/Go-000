package Week06

import (
	"container/list"
	"log"
	"sync"
	"time"
)

const (
	typeSuccess int8 = iota
	typeFail
)

type metrics struct {
	success int
	fail    int
}

type SlidingWindow struct {
	bucket int                //桶数
	curKey int64              //当前key
	m      map[int64]*metrics //统计
	data   *list.List
	sync.RWMutex
}

func NewSlidingWindow(bucket int) *SlidingWindow {
	return &SlidingWindow{
		bucket: bucket,
		//m:      make(map[int64]*metrics),
		data: list.New(),
	}
}

func (sw *SlidingWindow) AddSuccess() {
	sw.incr(typeSuccess)
}

func (sw *SlidingWindow) AddFail() {
	sw.incr(typeFail)
}

func (sw *SlidingWindow) incr(t int8) {
	sw.Lock()
	defer sw.Unlock()

	nowTime := time.Now().Unix()
	if _, ok := sw.m[nowTime]; !ok {
		sw.m = make(map[int64]*metrics)
		sw.m[nowTime] = &metrics{}
	}
	if sw.curKey == 0 {
		sw.curKey = nowTime
	}

	//一秒一个bucket
	if sw.curKey != nowTime {
		sw.data.PushBack(sw.m[nowTime])
		delete(sw.m, sw.curKey)
		sw.curKey = nowTime
		if sw.data.Len() > sw.bucket {
			for i := 0; i <= sw.data.Len()-sw.bucket; i++ {
				sw.data.Remove(sw.data.Front())
			}
		}
	}

	switch t {
	case typeSuccess:
		sw.m[nowTime].success++
	case typeFail:
		sw.m[nowTime].fail++
	default:
		log.Fatal("err type")
	}
}

func (sw *SlidingWindow) Len() int {
	return sw.data.Len()
}

func (sw *SlidingWindow) Data(space int) []*metrics {
	sw.RLock()
	defer sw.RUnlock()

	var (
		data []*metrics
		num  = 0
		m    = &metrics{}
	)
	for i := sw.data.Front(); i != nil; i = i.Next() {
		one := i.Value.(*metrics)
		m.success += one.success
		m.fail += one.fail
		if num%space == 0 {
			data = append(data, m)
			m = &metrics{}
		}
		num++
	}
	return data
}
