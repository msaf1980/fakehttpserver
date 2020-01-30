package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

var (
	minDelay int64
	maxDelay int64
)

func main() {
	var err error
	var t time.Duration

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	addr := os.Getenv("LISTEN")
	if len(addr) == 0 {
		addr = ":8080"
	}
	var s string
	s = os.Getenv("MINDELAY")
	if len(s) > 0 {
		t, err = time.ParseDuration(s)
		if err != nil {
			log.Fatal().Str("error", err.Error()).Msg("error in MINDELAY")
			os.Exit(1)
		}
		minDelay = t.Nanoseconds()
	}
	s = os.Getenv("MAXDELAY")
	if len(s) > 0 {
		t, err = time.ParseDuration(s)
		if err != nil {
			log.Fatal().Str("error", err.Error()).Msg("error in MAXDELAY")
			os.Exit(1)
		}
		maxDelay = t.Nanoseconds()
	}
	if minDelay > 0 && maxDelay == 0 {
		maxDelay = minDelay
	}
	if minDelay < 0 || maxDelay < 0 || minDelay > maxDelay {
		log.Fatal().Str("error", "incorrect values").Msg("error in MINDELAY/MAXDELAY")
		os.Exit(1)
	}

	h := requestHandler

	if err := fasthttp.ListenAndServe(addr, h); err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error in ListenAndServe")
		os.Exit(1)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	bodyLen := len(ctx.Request.Body())
	reqDuration := ctx.Time().Sub(ctx.ConnTime())
	ctx.SetContentType("text/plain; charset=utf8")
	resp := []byte("OK\n")
	r := time.Duration(rand.Int63n(maxDelay-minDelay) + minDelay)
	time.Sleep(r)
	n, _ := ctx.Write(resp)
	//fmt.Fprintf(ctx, "OK\n")
	log.Info().
		Str("method", string(ctx.Method())).
		Str("uri", string(ctx.RequestURI())).
		Str("server", string(ctx.Host())).
		Str("remote", ctx.RemoteIP().String()).
		Int("headersize", ctx.Request.Header.ContentLength()).
		Int("reqsize", bodyLen).
		Int("respsize", n).
		Str("reqtime", fmt.Sprintf("%.2f", float64(reqDuration.Nanoseconds())/1000000)).
		Str("sleeptime", fmt.Sprintf("%.2f", float64(r.Nanoseconds())/1000000)).
		Msg("query")
	/*
		fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
		fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
		fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
		fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
		fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
		fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
		fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
		fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
		fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
		fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

		fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)
	*/
}
