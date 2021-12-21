package serv

import (
	"lzero/internal/data"
	"lzero/internal/utils"
	"net/http"
)

const (
	HTTP_ADDR = "127.0.0.1"
	HTTP_PORT = ":3000"
)

type Server struct {
	http.Handler
}

func NewServer(orderSet *data.Cache) Server {
	l := utils.NewLogger()

	l.InfoLog.Println("New server")

	s := new(Server)

	router := *http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Origin", "*")

		case http.MethodGet:
			l.InfoLog.Printf("Requested order UID: %s", r.Header.Get("order_uid"))

			resp, ok := orderSet.GetOrder(r.Header.Get("order_uid"))
			if !ok {
				l.ErrorLog.Println("No order")
				w.WriteHeader(http.StatusNoContent)
			} else {
				l.InfoLog.Println("Sending order")

				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Access-Control-Allow-Headers", "*")
				w.Header().Set("Access-Control-Allow-Methods", "*")
				w.Header().Set("Access-Control-Allow-Origin", "*")

				parsedResp, err := resp.Marshal()
				if err != nil {
					l.ErrorLog.Printf("Cant marshal order: %s", err)
				}

				w.Write(parsedResp)
			}

		default:
			http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		}
	})

	s.Handler = &router

	return *s
}