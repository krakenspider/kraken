package kraken

import (
	"fmt"
	"math"
	"net/http"
)

const abortIndex int8 = math.MaxInt8 >> 1

func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		trees:              make(methodTrees, 0, 9),
		UseRawPath:         false,
		UnescapePathValues: true,
		RemoveExtraSlash:   false,
	}
	engine.RouterGroup.engine = engine
	return engine
}

type Engine struct {
	RouterGroup
	trees              methodTrees
	maxParams          uint16
	maxSections        uint16
	UseRawPath         bool
	UnescapePathValues bool
	RemoveExtraSlash   bool
}

func (engine *Engine) addRoute(path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(len(handlers) > 0, "there must be at least one handler")

	root := engine.trees.get(http.MethodGet)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		engine.trees = append(engine.trees, methodTree{method: http.MethodGet, root: root})
	}
	root.addRoute(path, handlers)

	if paramsCount := countParams(path); paramsCount > engine.maxParams {
		engine.maxParams = paramsCount
	}

	if sectionsCount := countSections(path); sectionsCount > engine.maxSections {
		engine.maxSections = sectionsCount
	}
}

func (engine *Engine) handleHTTPRequest(c *Context) (err error) {
	httpMethod := http.MethodGet
	rPath := c.URL.Path
	unescape := false
	if engine.UseRawPath && len(c.URL.RawPath) > 0 {
		rPath = c.URL.RawPath
		unescape = engine.UnescapePathValues
	}
	if engine.RemoveExtraSlash {
		rPath = cleanPath(rPath)
	}

	// Find root of the tree for the given HTTP method
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// Find route in tree
		var params = new(Params)
		var skip = new([]skippedNode)
		value := root.getValue(rPath, params, skip, unescape)
		if value.params != nil {
			c.Params = *value.params
		}
		if value.handlers != nil {
			c.handlers = value.handlers
			c.fullPath = value.fullPath
			for _, handler := range value.handlers {
				err = handler(c)
				if err != nil {
					return
				}
			}
			return
		}
		break
	}
	err = fmt.Errorf("url: %s handler not found", c.URL.String())
	return
}
