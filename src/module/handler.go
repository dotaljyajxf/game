package module

import (
	"public"
)

type TContextHandler struct {
	context *public.TContext
}

func (this *TContextHandler) SetContext(aContext *public.TContext) {
	this.context = aContext
}
func (this *TContextHandler) GetContext() *public.TContext {
	return this.context
}
