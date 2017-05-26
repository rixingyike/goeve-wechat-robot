package wxweb

import (
	"fmt"
	"sync"
	"../logs"
)

// Handler: message function wrapper
type Handler func(*Session, *ReceivedMessage)

// HandlerWrapper: message handler wrapper
type HandlerWrapper struct {
	handle  Handler
	enabled bool
	name    string
}

// Run: message handler callback
func (s *HandlerWrapper) Run(session *Session, msg *ReceivedMessage) {
	if s.enabled {
		s.handle(session, msg)
	}
}

func (s *HandlerWrapper) getName() string {
	return s.name
}

func (s *HandlerWrapper) IsEnabled() bool {
	return s.enabled
}

func (s *HandlerWrapper) enableHandle() {
	s.enabled = true
	return
}

func (s *HandlerWrapper) disableHandle() {
	s.enabled = false
	return
}

// HandlerRegister: message handler manager
type HandlerRegister struct {
	mu   sync.RWMutex
	hmap map[int][]*HandlerWrapper
}

// CreateHandlerRegister: create handler register
func CreateHandlerRegister() *HandlerRegister {
	return &HandlerRegister{
		hmap: make(map[int][]*HandlerWrapper),
	}
}

// Add: add message callback handle to handler register
func (hr *HandlerRegister) Add(key int, h Handler, name string) error {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	for _, v := range hr.hmap {
		for _, handle := range v {
			if handle.getName() == name {
				return fmt.Errorf("handler name %s has been registered", name)
			}
		}
	}
	hr.hmap[key] = append(hr.hmap[key], &HandlerWrapper{handle: h, enabled: false, name: name})
	return nil
}

// 批量处理某个消息
func (hr *HandlerRegister) Runs(session *Session, msg *ReceivedMessage) {
	err, handles := hr.Get(msg.MsgType)
	if err != nil {
		logs.Error(err)
		return
	}
	for _, v := range handles {
		go v.Run(session, msg)
	}
}

// Get: get message handler
func (hr *HandlerRegister) Get(key int) (error, []*HandlerWrapper) {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	if v, ok := hr.hmap[key]; ok {
		return nil, v
	}
	return fmt.Errorf("handlers for key [%d] not registered", key), nil
}

// EnableByType: enable handler by message type
func (hr *HandlerRegister) EnableByType(key int) error {
	err, handles := hr.Get(key)
	if err != nil {
		return err
	}
	hr.mu.Lock()
	defer hr.mu.Unlock()
	// all
	for _, v := range handles {
		v.enableHandle()
	}
	return nil
}

// DisableByType: disable handler by message type
func (hr *HandlerRegister) DisableByType(key int) error {
	err, handles := hr.Get(key)
	if err != nil {
		return err
	}
	hr.mu.Lock()
	defer hr.mu.Unlock()
	// all
	for _, v := range handles {
		v.disableHandle()
	}
	return nil
}

// 使用名称查询插件
func (hr *HandlerRegister) GetByName(name string) *HandlerWrapper {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	for _, handles := range hr.hmap {
		for _, v := range handles {
			if v.getName() == name {
				return v
			}
		}
	}
	return nil
}

func (hr *HandlerRegister) EnableByName(name string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	for _, handles := range hr.hmap {
		for _, v := range handles {
			if v.getName() == name {
				v.enableHandle()
			}
		}
	}
}

// DisableByName: disable message handler by name
func (hr *HandlerRegister) DisableByName(name string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	for _, handles := range hr.hmap {
		for _, v := range handles {
			if v.getName() == name {
				v.disableHandle()
			}
		}
	}
}

// Dump: output all message handlers
func (hr *HandlerRegister) Dump() string {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	str := "[plugins dump]\n"
	for k, handles := range hr.hmap {
		for _, v := range handles {
			str += fmt.Sprintf("%d %s [%v]\n", k, v.getName(), v.enabled)
		}
	}
	return str
}
