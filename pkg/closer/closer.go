package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = New()

func Add(fn ...func() error) {
	globalCloser.Add(fn...)
}

func Wait() {
	globalCloser.Wait()
}

func CloseAll() {
	globalCloser.CloseAll()
}

type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

// New return new closer
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

// Add func to closer
func (c *Closer) Add(fn ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, fn...)
	c.mu.Unlock()
}

// Wait blocks until all closer function are done
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll calls all closer functions
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, fn := range funcs {
			go func(f func() error) {
				errs <- f()
			}(fn)
		}
		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Printf("error returned from Closer: %d", err)
			}
		}
	})
}
