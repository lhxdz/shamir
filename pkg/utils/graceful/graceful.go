package graceful

import (
	"os"
	"os/signal"
)

type Handler func()

type Option func(g *Graceful)

func WithSignalHandlers(handlers map[os.Signal]Handler) Option {
	return func(g *Graceful) {
		g.handlers = handlers
	}
}

type Graceful struct {
	handlers map[os.Signal]Handler
	finish   chan struct{}
}

func NewGraceFul(opts ...Option) *Graceful {
	g := &Graceful{
		handlers: map[os.Signal]Handler{},
		finish:   make(chan struct{}, 1),
	}

	for _, opt := range opts {
		opt(g)
	}

	go g.listen()
	return g
}

func (g *Graceful) listen() {
	signals := make([]os.Signal, 0, len(g.handlers))
	for s := range g.handlers {
		signals = append(signals, s)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, signals...)
	s := <-quit
	handler, ok := g.handlers[s]
	if ok {
		handler()
	}

	close(g.finish)
}

// Wait 创建完成graceful模块后，需要调用 Wait 方法阻塞式等待
func (g *Graceful) Wait() {
	<-g.finish
}
