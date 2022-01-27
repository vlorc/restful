# [Restful](https://github.com/vlorc/restful)
Golang restful minimum project

## Quick Start

```go
func main() {
    engine.Init()
    
    r := engine.NewRouter()
    
    r.Any("/echo", func(req *EchoRequest) (*EchoResponse, error) {
        return &EchoResponse{Data: req.Data}, nil
    })
    
    http.ListenAndServe(":1234", r)
}
```