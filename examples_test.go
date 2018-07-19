package retry

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"
)

func ExampleEnsure() {
	r := New()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

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
