package notify

type Handler func([]byte) error

type Router struct {
	handlers map[string]Handler
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]Handler),
	}
}

func (r *Router) Register(eventType string, handler Handler) {
	r.handlers[eventType] = handler
}

func (r *Router) Handle(eventType string, body []byte) error {
	h, ok := r.handlers[eventType]
	if !ok {
		return nil
	}
	return h(body)
}