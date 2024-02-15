package utils

import (
	"fmt"
	"sync"
)

type Background struct {
	logger Logger
	wg     sync.WaitGroup
}

func (b *Background) Run(fn func()) {
	b.wg.Add(1)

	go func() {
		defer b.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				b.logger.LogError(fmt.Errorf("%s", err))
			}
		}()

		fn()
	}()
}

func (b *Background) Wait() {
	b.wg.Wait()
}
