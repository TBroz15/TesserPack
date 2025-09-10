package helpers

import (
	"context"
	"runtime"

	"github.com/marusama/semaphore"
)

type Semaphore struct {
	sem semaphore.Semaphore
}

/*
Linux handles alot of concurrent disk operation just fine.

But Windows will shat itself and lag the entire computer
since it had to go through CGO calls and stuff, 
leading to high memory usage and panics.

This only solves for Windows systems and when the pack is too large to handle.

For the longest time I have not yet tested TesserPack for Windows. Since I code on Fedora btw, dont tell me to use arch BTW -TuxeBro, 2025
*/
func NewSemaphore(n int) *Semaphore {
	if runtime.GOOS != "windows" {
		return nil
	}

	sem := semaphore.New(n)

	return &Semaphore{
		sem: sem,
	}
}

func (s *Semaphore) Acquire() {
	if s == nil {
		return
	}
	s.sem.Acquire(context.Background(), 1)
}

func (s *Semaphore) Release() {
	if s == nil {
		return
	}
	s.sem.Release(1)
}
