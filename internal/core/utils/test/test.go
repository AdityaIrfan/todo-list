package test

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type PrepareContext struct {
	Method  string
	URL     string
	Payload []byte
	App     *fiber.App
}

func (e *PrepareContext) PrepareHTTPTestContext() *fiber.Ctx {
	var header fasthttp.RequestHeader
	header.SetMethod(e.Method)
	header.SetContentType("application/json")
	req := fasthttp.Request{
		Header: header,
	}
	req.AppendBody(e.Payload)
	req.SetRequestURI(e.URL)

	requestCtx := &fasthttp.RequestCtx{
		Request: req,
	}
	ctx := e.App.AcquireCtx(requestCtx)
	return ctx
}
