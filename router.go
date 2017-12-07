package main

func (r *Router) Handle(msgName string, handler Handler) {
	// r.rules[msgName] = handler
}

// func (r *Router) FindHandler(msgName string) (Handler, bool) {
// 	// handler, found := r.rules[msgName]
// 	return handler, found
// }
