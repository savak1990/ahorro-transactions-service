package handler

import "net/http"

func EnsureAwsRegionHeader(region string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-AWS-Region") == "" {
				r.Header.Set("X-AWS-Region", region)
			}
			next.ServeHTTP(w, r)
		})
	}
}
