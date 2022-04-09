package yuriko

import (
	"fmt"
	"sync"
	"time"
)

// 比特位划分
const (
	seqBits          = 12 // 序号位共12bits
	workerIdBits     = 5  // workerId 5bits
	dataCenterIdBits = 5  // dataCenterId 5bits
	signBits         = 1  // 符号位 1bit
	// 时间戳 41bits，以毫秒计算，可以容纳（1 << 42 - 1）/ 1000 / 86400 / 365 = 69 年
	timestampBits = 64 - signBits - workerIdBits - dataCenterIdBits - seqBits
)

const (
	maxSeq          = 1<<seqBits - 1
	maxWorkerId     = 1<<workerIdBits - 1
	maxDataCenterId = 1<<dataCenterIdBits - 1
	maxTimestamp    = 1<<timestampBits - 1
)

// 右移位数
const (
	workerIdShift     = seqBits
	dataCenterIdShift = workerIdBits + seqBits
	timestampShift    = dataCenterIdBits + workerIdBits + seqBits
)

// 全局变量
var (
	lastTime int64 = -1 // 上次调用的时间
	seq      int64 = 0  // 序号共12位，即：每台机器最多一秒生成(1 << 12 - 1)个序号
)

type SnowFlake struct {
	mu           sync.Mutex
	startTime    int64
	dataCenterId int64
	workerId     int64
}

func NewSnowFlake(workerId, dataCenterId, startTime int64) (*SnowFlake, error) {
	if workerId > maxWorkerId {
		return nil, fmt.Errorf("workerId > maxWorkerId: %d > %d", workerId, maxWorkerId)
	}
	if dataCenterId > maxDataCenterId {
		return nil, fmt.Errorf("dataCenterId > maxDataCenterId: %d > %d", dataCenterId, maxDataCenterId)
	}
	current := currentTime()
	if startTime > current {
		return nil, fmt.Errorf("startTime > currentTime: %d > %d", startTime, current)
	}
	if current-startTime > maxTimestamp {
		return nil, fmt.Errorf("timstamp exceeded max: %d", maxTimestamp)
	}
	return &SnowFlake{
		startTime:    startTime,
		dataCenterId: dataCenterId,
		workerId:     workerId,
	}, nil
}

// 生成64位的UUID
func (sf *SnowFlake) generateUUID() int64 {
	sf.mu.Lock()
	defer sf.mu.Unlock()
	current := currentTime()
	switch {
	case current < lastTime: // 时钟出现回拨了，这时重新等待，直到时间过了lastTime
		current = waitUntil(lastTime)
		seq = 0
	case current == lastTime: // 在同一毫秒，增加seq，时间戳不变
		seq = (seq + 1) & maxSeq
		if seq == 0 { // seq 溢出了，这时需要等待1ms
			current = waitUntil(current)
		}
	case current > lastTime: // 不在同一毫秒，直接用时间戳即可
		seq = 0
	}
	lastTime = current
	return (current-sf.startTime)<<timestampShift | sf.dataCenterId<<dataCenterIdShift | sf.workerId<<workerIdShift | seq
}

func currentTime() int64 {
	return time.Now().UnixMilli()
}

// 通过sleep的方式，来轮询是否过了目标时间
func waitUntil(target int64) int64 {
	current := currentTime()
	for current <= target {
		time.Sleep(1 * time.Microsecond)
		current = currentTime()
	}
	return current
}
