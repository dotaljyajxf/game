package gateway

type TContextHandler struct {
	context *TContext
}

func (this *TContextHandler) SetContext(aContext *TContext) {
	this.context = aContext
}
func (this *TContextHandler) GetContext() *TContext {
	return this.context
}
