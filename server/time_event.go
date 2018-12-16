package server

import (
	"redis_go/client"
	"sync"
	"sync/atomic"
	"time"
)

type EventLoop struct {
	eventIDSequence uint64                // 用于生成time event id
	events          map[uint64]*TimeEvent // 所有的timeEvents
	lastTime        time.Time             // 上次执行时间
	stop            bool                  // 是否停止
	lock            sync.Mutex            // 锁
}

type TimeEvent struct {
	ID           uint64               // event id
	ArriveSecond int64                // 到达时间 second
	Interval     int64                // 时间间隔 milliseconds
	mask         int64                // 事件类型掩码，可以是 AE_READABLE 或 AE_WRITABLE
	ClientDate   *client.Client       // client data
	Proc         func(*client.Client) // 处理函数
}

func NewEventLoop() *EventLoop {
	return &EventLoop{
		eventIDSequence: 0,
		lastTime:        time.Now(),
		stop:            false,
		events:          make(map[uint64]*TimeEvent),
	}
}

func (el *EventLoop) NewTimeEvent(interval int64, mask int64, clientDate *client.Client, proc func(*client.Client)) error {
	el.lock.Lock()
	defer el.lock.Unlock()

	// TODO 参数校验

	timeEvent := &TimeEvent{
		ID:         atomic.AddUint64(&el.eventIDSequence, 1),
		Interval:   interval,
		mask:       mask,
		ClientDate: clientDate,
		Proc:       proc,
	}
	el.events[timeEvent.ID] = timeEvent
	return nil
}
