Retry is a pretty simple library to ensure your work to be done

[godoc](https://godoc.org/github.com/shafreeck/retry)

## Features
* Retry to run a workflow(Ex. rpc or db access)
* Customize backoff strategy
* Retry accoding to your type of error

## Examples

```
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
```
