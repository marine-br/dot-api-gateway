package proxy

type Route struct {
	Method  string
	Pattern string
	Backend string
}

var routes = []Route{

	{Method: "GET", Pattern: "/api/v1/users", Backend: "https://jsonplaceholder.typicode.com/users"},
	{Method: "POST", Pattern: "/api/v1/posts", Backend: "https://jsonplaceholder.typicode.com/posts"},
}

func FindBackend(method, path string) (string, bool) {
	for _, r := range routes {
		if r.Method == method && r.Pattern == path {
			return r.Backend, true
		}
	}
	return "", false
}
