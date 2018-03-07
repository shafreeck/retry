Retry 是一个通用重试框架，可以适配所有需要重试操作的场合

## 功能
* 支持重试取消
* 支持自定义重试策略
* 支持非幂等函数在提交数据前出错重试

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
