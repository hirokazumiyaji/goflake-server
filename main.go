package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hirokazumiyaji/goflake"
	"github.com/ugorji/go/codec"
	"github.com/valyala/fasthttp"
)

const version = "0.0.1"

var mh = &codec.MsgpackHandle{RawToString: true}

func main() {

	var (
		ip           string
		port         int
		startTime    string
		datacenterID int
		workerID     int
		retry        int
		err          error
	)

	defaultIP := os.Getenv("IP_ADDR")
	if defaultIP == "" {
		defaultIP = "127.0.0.1"
	}
	flag.StringVar(&ip, "addr", defaultIP, "ip address (default 127.0.0.1)")
	flag.StringVar(&ip, "a", defaultIP, "ip address (short)")

	defaultPort := 8000
	if v := os.Getenv("PORT"); v != "" {
		defaultPort, err = strconv.Atoi(v)
		if err != nil {
			log.Fatalf("port value type not integer. %v\n", err)
		}
	}
	flag.IntVar(&port, "port", defaultPort, "port to use (default 8000)")
	flag.IntVar(&port, "p", defaultPort, "port to use (short)")

	defaultStartTime := os.Getenv("START_TIME")
	if defaultStartTime == "" {
		defaultStartTime = "2016-01-01 00:00:00 +0000"
	}
	flag.StringVar(&startTime, "start-time", defaultStartTime, "id generate start time (default '2016-01-01 00:00:00 +0000')")
	flag.StringVar(&startTime, "s", defaultStartTime, "id generate start time (short)")

	defaultDatacenterID := 1
	if v := os.Getenv("DATACENTER_ID"); v != "" {
		defaultDatacenterID, err = strconv.Atoi(v)
		if err != nil {
			log.Fatalf("datacenter id value type not integer. %v\n", err)
		}
	}
	flag.IntVar(&datacenterID, "datacenter-id", defaultDatacenterID, "datacenter id (default 1)")
	flag.IntVar(&datacenterID, "d", defaultDatacenterID, "datacenter id (short)")

	defaultWorkerID := 1
	if v := os.Getenv("WORKER_ID"); v != "" {
		defaultWorkerID, err = strconv.Atoi(v)
		if err != nil {
			log.Fatalf("worker id value type not integer. %v\n", err)
		}
	}
	flag.IntVar(&workerID, "worker-id", defaultWorkerID, "worker id (default 1)")
	flag.IntVar(&workerID, "w", defaultWorkerID, "worker id (short)")

	defaultRetry := 5
	if v := os.Getenv("RETRY"); v != "" {
		defaultRetry, err = strconv.Atoi(v)
		if err != nil {
			log.Fatalf("retry value type not integer. %v\n", err)
		}
	}
	flag.IntVar(&retry, "retry", defaultRetry, "generate id retry count (default 5)")
	flag.IntVar(&retry, "r", defaultRetry, "generate id retry count (short)")

	flag.Parse()

	t, err := time.Parse("2006-01-02 15:04:05 -0700", startTime)
	if err != nil {
		log.Fatalf("start time parse error. %v\n", err)
	}

	idWorker, err := goflake.NewIdWorker(uint16(datacenterID), uint16(workerID), t)
	if err != nil {
		log.Fatalf("could not create id worker. %v\n", err)
	}

	m := func(ctx *fasthttp.RequestCtx) {
		if strings.HasPrefix(string(ctx.Path()), "/id") {
			idHandlerFunc(ctx, idWorker, retry)
		} else {
			ctx.Error("", fasthttp.StatusNotFound)
		}
	}

	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), m))
}

func idHandlerFunc(ctx *fasthttp.RequestCtx, idWorker *goflake.IdWorker, retry int) {
	ua := string(ctx.UserAgent())

	var (
		id  uint64
		err error
	)

	for i := 0; i < retry; i++ {
		id, err = idWorker.GetId(ua)
		if err == nil {
			break
		}
	}

	r := map[string]string{
		"id": strconv.FormatUint(id, 10),
	}

	if strings.HasSuffix(string(ctx.Path()), ".msgpack") {
		ctx.SetContentType("application/x-msgpack; charset=UTF-8")
		if err := codec.NewEncoder(ctx, mh).Encode(r); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		}
	} else {
		ctx.SetContentType("application/json; charset=UTF-8")
		if err := json.NewEncoder(ctx).Encode(r); err != nil {
			ctx.Error(fmt.Sprintf(`{"error":"%v"}`, err.Error()), fasthttp.StatusInternalServerError)
		}
	}
}
