package kraken

import (
	"context"
	"errors"
	"github.com/gozelle/logger"
	"net/url"
)

type HandlerFunc func(*Context) error
type HandlersChain []HandlerFunc

type Context struct {
	context.Context
	*Extractor
	*logger.Logger
	URL      url.URL
	Params   Params
	handlers HandlersChain
	index    int8
	cancel   context.CancelFunc

	//abort    func() bool
}

func (c *Context) JustWait() {
	JustWait()
}

func (c *Context) JustThink() {
	JustThink()
}

//func (c *Context) Abort(fn func() bool) {
//	c.abort = fn
//}

func (c *Context) HandleData(data any) {
	if c.Extractor != nil {
		if c.Extractor.crawler != nil {
			c.Extractor.crawler.data <- data
		}
	}
}

func (c *Context) Done(err ...error) {
	var e error
	if len(err) > 0 {
		e = errors.Join(err...)
	}
	c.crawler.sendDone(e)
}

func (c *Context) reset() {
	c.Params = c.Params[:0]
	c.handlers = nil
	c.index = -1
}

//func (c *Context) stop() {
//	if c.Extractor != nil {
//		c.Extractor.stop()
//	}
//}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() (err error) {
	c.index++
	for c.index < int8(len(c.handlers)) {
		if c.handlers[c.index] != nil {
			err = c.handlers[c.index](c)
			if err != nil {
				return
			}
		}
		c.index++
	}
	return nil
}

func (c *Context) Async(func()) {

}

func (c *Context) Wait() {

}
