package proxy

import (
	"bufio"
	"io"
	"log"
	"net/http"
)

type Transformer interface {
	Transform(string) string
}

func Start(transformerConstructor func() Transformer) {
	handleHTTP := func(w http.ResponseWriter, req *http.Request) {
		resp, err := http.Get("http://localhost:9299/metrics")

		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)

		transformer := transformerConstructor()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			io.WriteString(w, transformer.Transform(line))
			w.Write([]byte("\n"))
		}
	}

	http.HandleFunc("/metrics", handleHTTP)
	log.Fatal(http.ListenAndServe("0.0.0.0:19299", nil))
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
