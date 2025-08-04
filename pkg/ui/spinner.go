package ui

import (
	"fmt"
	"sync"
	"time"
)

type Spinner struct {
	message    string
	chars      []string
	delay      time.Duration
	active     bool
	mutex      sync.Mutex
	stopChan   chan bool
	doneChan   chan bool
}

func NewSpinner(message string) *Spinner {
	return &Spinner{
		message:  message,
		chars:    []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		delay:    100 * time.Millisecond,
		stopChan: make(chan bool, 1),
		doneChan: make(chan bool, 1),
	}
}

func (s *Spinner) Start() {
	s.mutex.Lock()
	if s.active {
		s.mutex.Unlock()
		return
	}
	s.active = true
	s.mutex.Unlock()

	go func() {
		i := 0
		for {
			select {
			case <-s.stopChan:
				fmt.Printf("\r\033[K")
				s.doneChan <- true
				return
			default:
				fmt.Printf("\r%s %s", s.chars[i%len(s.chars)], s.message)
				i++
				time.Sleep(s.delay)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.mutex.Lock()
	if !s.active {
		s.mutex.Unlock()
		return
	}
	s.active = false
	s.mutex.Unlock()

	s.stopChan <- true
	<-s.doneChan
}

func (s *Spinner) UpdateMessage(message string) {
	s.mutex.Lock()
	s.message = message
	s.mutex.Unlock()
}

func (s *Spinner) Success(message string) {
	s.Stop()
	fmt.Printf("✅ %s\n", message)
}

func (s *Spinner) Error(message string) {
	s.Stop()
	fmt.Printf("❌ %s\n", message)
}

func (s *Spinner) Warning(message string) {
	s.Stop()
	fmt.Printf("⚠️  %s\n", message)
}

// Utility functions for quick operations
func WithSpinner(message string, fn func() error) error {
	spinner := NewSpinner(message)
	spinner.Start()
	defer spinner.Stop()
	
	return fn()
}

func WithSpinnerResult(message string, fn func() error) error {
	spinner := NewSpinner(message)
	spinner.Start()
	
	err := fn()
	if err != nil {
		spinner.Error(fmt.Sprintf("Failed: %s", message))
		return err
	} else {
		spinner.Success(message)
		return nil
	}
}