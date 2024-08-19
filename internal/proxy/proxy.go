package proxy

import "net/http"

func TodoApi() http.Handler() {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse(configs.GetConfig().Service.TODO_ENDPOINT)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		httpProxy := httputil.NewSingleHostReverseProxy(url)
		httpProxy.ServeHTTP(w, r)
	})
}