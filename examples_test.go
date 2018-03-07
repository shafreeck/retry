package retry

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"
)

func ExampleHTTPGet() {
	var resp *http.Response
	var err error
	get := func() error {
		resp, err = http.Get("http://www.example.com")
		if err != nil {
			return Retriable(err)
		}
		return nil
	}

	r := New()
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	if err := r.Ensure(ctx, get); err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)
}

func ExampleHTTPGetEasier() {
	r := New()
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err := r.Ensure(ctx, func() error {
		resp, err := http.Get("http://www.example.com")
		if err != nil {
			return Retriable(err)
		}

		log.Println(resp)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleHTTPGetAndPost() {
	r := New()
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	err := r.Ensure(ctx, func() error {
		resp, err := http.Get("http://www.example.com")
		// get error can be retried
		if err != nil {
			log.Println(err)
			return Retriable(err)
		}
		log.Println(resp)

		buf := bytes.NewBuffer(nil)
		resp, err = http.Post("http://example.com/upload", "image/jpeg", buf)
		// post error should not be retried
		if err != nil {
			return err
		}
		log.Println(resp)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
