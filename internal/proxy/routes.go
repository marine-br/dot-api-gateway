package proxy

type Route struct {
	Method  string
	Pattern string
	Backend string
}

var routes = []Route{

	{Method: "GET", Pattern: "/api/v1/service1", Backend: "http://service-1:8080"},
	{Method: "POST", Pattern: "/api/v1/service2", Backend: "http://service-2:8080"},
}

func FindBackend(method, path string) (string, bool) {
	for _, r := range routes {
		if r.Method == method && r.Pattern == path {
			return r.Backend, true
		}
	}
	return "", false
}
