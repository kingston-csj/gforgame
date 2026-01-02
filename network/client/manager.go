package client

import (
	"errors"
	"sync"
	"time"
)

// CallBackService 定义了回调服务
type CallBackService struct {
	mapper     map[int32]*RequestResponseFuture
	mapperLock sync.Mutex
	ticker     *time.Ticker
	stopChan   chan bool
}

var (
	CallBackManager *CallBackService
	ErrTimeOut      = errors.New("request timeout, no reply")
)

func init() {
	CallBackManager = NewCallBackService()
	CallBackManager.Start()
}

// NewCallBackService 初始化回调服务
func NewCallBackService() *CallBackService {
	return &CallBackService{
		mapper:   make(map[int32]*RequestResponseFuture),
		stopChan: make(chan bool),
	}
}

// Register 注册请求
func (s *CallBackService) Register(correlationId int32, future *RequestResponseFuture) {
	s.mapperLock.Lock()
	defer s.mapperLock.Unlock()
	s.mapper[correlationId] = future
}

// Remove 移除请求
func (s *CallBackService) Remove(correlationId int32) *RequestResponseFuture {
	s.mapperLock.Lock()
	defer s.mapperLock.Unlock()
	if future, ok := s.mapper[correlationId]; ok {
		delete(s.mapper, correlationId)
		return future
	}
	return nil
}

// Start 启动定时任务
func (s *CallBackService) Start() {
	s.ticker = time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.scanExpiredRequest()
			case <-s.stopChan:
				s.ticker.Stop()
				return
			}
		}
	}()
}

// Stop 停止定时任务
func (s *CallBackService) Stop() {
	s.stopChan <- true
}

// ScanExpiredRequest 扫描超时请求
func (s *CallBackService) scanExpiredRequest() {
	s.mapperLock.Lock()
	defer s.mapperLock.Unlock()
	var rfList []*RequestResponseFuture
	for correlationId, rf := range s.mapper {
		if rf.isTimeout() {
			delete(s.mapper, correlationId)
			rfList = append(rfList, rf)
		}
	}

	for _, rf := range rfList {
		cause := ErrTimeOut
		rf.Cause = cause
		// 异步回调
		if rf.RequestCallback != nil {
			rf.RequestCallback.OnError(rf.Cause)
		} else {
			// 同步调用
			rf.waitCause <- rf.Cause
		}
	}
}

// FillCallBack 填充回调
func (s *CallBackService) FillCallBack(index int32, message any) {
	future := s.Remove(index)
	if future == nil {
		return
	}
	callback := future.RequestCallback
	// 异步回调
	if callback != nil {
		callback.OnSuccess(message)
	} else {
		// 同步调用
		future.waitResponse <- message
	}
}
