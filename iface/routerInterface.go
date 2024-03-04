package iface

type RouterInterface interface {
	BeforeHandle(request RequestInterface)

	Handle(request RequestInterface)

	AfterHandle(request RequestInterface)
}
