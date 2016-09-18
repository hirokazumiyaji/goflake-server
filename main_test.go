package main

import (
	"bufio"
	"testing"
	"time"

	"github.com/hirokazumiyaji/goflake"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestGetID(t *testing.T) {
	w, err := goflake.NewIDWorker(uint16(1), uint16(1), time.Now())
	if err != nil {
		t.Fatalf("new worker error. %v", err)
	}

	id, err := getID(w, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36", 5)
	if err != nil {
		t.Errorf("get id error. %v", err)
	}

	if id == 0 {
		t.Errorf("id eq 0")
	}
}

func TestHandlerNotFound(t *testing.T) {
	w, err := goflake.NewIDWorker(uint16(1), uint16(1), time.Now())
	if err != nil {
		t.Fatalf("new worker error. %v", err)
	}

	s := &fasthttp.Server{
		Handler: handler(w, 5),
		Name:    "test-goflake-server",
	}

	il := fasthttputil.NewInmemoryListener()

	sc := make(chan bool)
	go func() {
		if err := s.Serve(il); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		close(sc)
	}()

	cc := make(chan bool)
	go func() {
		c, err := il.Dial()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if _, err = c.Write([]byte("GET / HTTP/1.1\r\nHost: 127.0.0.1\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36\r\n\r\n")); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		br := bufio.NewReader(c)
		var resp fasthttp.Response
		if err = resp.Read(br); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if resp.StatusCode() != fasthttp.StatusNotFound {
			t.Fatalf("unexpected status code: %d. Expecting %d", resp.StatusCode(), fasthttp.StatusNotFound)
		}
		close(cc)
	}()

	select {
	case <-cc:
	case <-time.After(time.Second):
		t.Fatalf("timeout")
	}

	if err := il.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	select {
	case <-sc:
	case <-time.After(time.Second):
		t.Fatalf("timeout")
	}
}

func TestHandlerIDJSON(t *testing.T) {
	w, err := goflake.NewIDWorker(uint16(1), uint16(1), time.Now())
	if err != nil {
		t.Fatalf("new worker error. %v", err)
	}

	s := &fasthttp.Server{
		Handler: handler(w, 5),
		Name:    "test-goflake-server",
	}

	il := fasthttputil.NewInmemoryListener()

	sc := make(chan bool)
	go func() {
		if err := s.Serve(il); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		close(sc)
	}()

	cc := make(chan bool)
	go func() {
		c, err := il.Dial()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if _, err = c.Write([]byte("GET /id HTTP/1.1\r\nHost: 127.0.0.1\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36\r\n\r\n")); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		br := bufio.NewReader(c)
		var resp fasthttp.Response
		if err = resp.Read(br); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if resp.StatusCode() != fasthttp.StatusOK {
			t.Fatalf("unexpected status code: %d. Expecting %d", resp.StatusCode(), fasthttp.StatusOK)
		}
		close(cc)
	}()

	select {
	case <-cc:
	case <-time.After(time.Second):
		t.Fatalf("timeout")
	}

	if err := il.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	select {
	case <-sc:
	case <-time.After(time.Second):
		t.Fatalf("timeout")
	}
}

func TestHandlerIDMsgPack(t *testing.T) {
	w, err := goflake.NewIDWorker(uint16(1), uint16(1), time.Now())
	if err != nil {
		t.Fatalf("new worker error. %v", err)
	}

	s := &fasthttp.Server{
		Handler: handler(w, 5),
		Name:    "test-goflake-server",
	}

	il := fasthttputil.NewInmemoryListener()

	sc := make(chan bool)
	go func() {
		if err := s.Serve(il); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		close(sc)
	}()

	cc := make(chan bool)
	go func() {
		c, err := il.Dial()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if _, err = c.Write([]byte("GET /id.msgpack HTTP/1.1\r\nHost: 127.0.0.1\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36\r\n\r\n")); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		br := bufio.NewReader(c)
		var resp fasthttp.Response
		if err = resp.Read(br); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if resp.StatusCode() != fasthttp.StatusOK {
			t.Fatalf("unexpected status code: %d. Expecting %d", resp.StatusCode(), fasthttp.StatusOK)
		}
		close(cc)
	}()

	select {
	case <-cc:
	case <-time.After(time.Second):
		t.Fatalf("timeout")
	}

	if err := il.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	select {
	case <-sc:
	case <-time.After(time.Second):
		t.Fatalf("timeout")
	}
}
