package module

import "netserver"

type TContextHandler struct {
	context *netserver.TContext
}

func (this *TContextHandler) SetContext(aContext *netserver.TContext) {
	this.context = aContext
}
func (this *TContextHandler) GetContext() *netserver.TContext {
	return this.context
}
