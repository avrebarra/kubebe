package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/avrebarra/minivalidator"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	stdlog "log"
)

var (
	appname       = "kubebe"
	port          = 2323
	starttime     = time.Now()
	prettylogging = true
)

func main() {
	// setup structured log
	logsink := log.Logger
	if prettylogging {
		logsink = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	logsink = logsink.With().Str("app", appname).Int("port", port).Timestamp().Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(logsink.With().Str("level", "debug").Logger())
	log.Logger = logsink

	// prepare basic server
	stdlog.Println("setting up server...")
	r := mux.NewRouter()
	r.HandleFunc("/", HandleIndex()).Methods("GET")
	r.HandleFunc("/diceroll", HandleDiceRoll()).Methods("POST")
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			log.Info().
				Str("url", r.URL.String()).
				Str("verb", r.Method).
				Msg("request received")
		})
	})

	// start serving, bon appetité
	stdlog.Printf("using port http://localhost:%d to start server... ", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		err = fmt.Errorf("cannot start server: %w", err)
		panic(err)
	}
}

// ***

func HandleIndex() http.HandlerFunc {
	type ResponseData struct {
		ID             string    `json:"id"`
		State          string    `json:"state"`
		StartTimestamp time.Time `json:"start_time"`
		Uptime         string    `json:"uptime"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		respdata := Response{
			Code: "ok",
			Data: ResponseData{
				ID:             appname,
				State:          "healthy",
				StartTimestamp: starttime,
				Uptime:         time.Since(starttime).Round(time.Second).String(),
			},
		}
		respond(w, http.StatusOK, respdata)
	}
}

func HandleDiceRoll() http.HandlerFunc {
	type RequestData struct {
		DiceNumber int `json:"dice_number" validate:"required"`
	}
	type ResponseData struct {
		Total    int   `json:"total"`
		Sequence []int `json:"sequence"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// handle input
		reqdata := RequestData{}
		if err := json.NewDecoder(r.Body).Decode(&reqdata); err != nil {
			respdata := Response{Code: "bad", Message: err.Error()}
			respond(w, http.StatusInternalServerError, respdata)
			return
		}
		if err := minivalidator.Validate(reqdata); err != nil {
			respdata := Response{Code: "bad", Message: err.Error()}
			respond(w, http.StatusInternalServerError, respdata)
			return
		}

		// process request
		sequence := []int{}
		total := 0
		for i := 0; i < reqdata.DiceNumber; i++ {
			num := rand.Intn(6-1) + 1
			total += num
			sequence = append(sequence, num)
		}

		// prepare output
		respdata := Response{
			Code: "ok",
			Data: ResponseData{
				Sequence: sequence,
				Total:    total,
			},
		}
		respond(w, http.StatusOK, respdata)
	}
}

// ***

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func respond(writer http.ResponseWriter, status int, data any) {
	respout, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("cannot marshall output: %w", err)
		log.Error().Err(err).Msg("marshall failure")
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(respout)
}
