package util

import "sync"

// SafeMap for map
type SafeMap struct {
	Data map[string]interface{}
	Lock sync.RWMutex
}

// Get map data
func (s *SafeMap) Get(k string) interface{} {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	if v, exit := s.Data[k]; exit {
		return v
	}
	return nil
}

// Set map data
func (s *SafeMap) Set(k string, v interface{}) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if s.Data == nil {
		s.Data = make(map[string]interface{})
	}
	s.Data[k] = v
}
