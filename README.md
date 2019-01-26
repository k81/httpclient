# httpclient library support customization and didi TraceInfo

# Examples
## Simple GET
```go
    client := httpclient.New(context.TODO(), Timeout(time.Second), DisableRedirect)
    query := url.Values{}
    query.Add("q1", "v1")
    query.Add("q2", "v2")
    resultBody, err := client.Get("http://api.example.com/hello", "", SetQuery(query))
```

## Simple POST
```go
    client := httpclient.New(context.TODO(), Timeout(time.Second), DisableRedirect)
    form := url.Values{}
    form.Add("q1", "v1")
    form.Add("q2", "v2")
    resultBody, err := client.Post("http://api.example.com/hello", form.Encode(), SetTypeForm())
```

## JSON POST
```go
    type Request struct {
        Hello string `json:"hello"`
    }

    type Response struct {
        ErrNo   int     `json:"errno"`
        ErrMsg  string  `json:"errmsg"`
        Data    string  `json:"data"`
    }

    client := httpclient.NewJSON(context.TODO(), Timeout(time.Second), DisableRedirect)
    req := &Request{
        Hello: "world",
    }
    result := &Response{}
    err := client.Post("http://api.example.com/hello", req, result)
```

# Customization Support
```go
    // customize client
    transport := http.Transport{
        MaxIdleConns: 16,
        //...
    }
    client := httpclient.NewJSON(context.TODO(), Timeout, DisableRedirect, SetTransport(transport))

    // customize request for signature
    func SignRequest(req *http.Request) error {
        // add signature to request
    }

    req := &Request{
        // build request...
    }
    resp := &Response{}

    err := client.Post("http://api.example.com/hello", req, resp, SignRequest)
```
