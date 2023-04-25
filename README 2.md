### 需要更改源码协议

```go
// ConnectToWS connects using webshocket protocol
func (c *Client) ConnectToWS(addr string, path string, rawQuery string, tlsConfig ...*tls.Config) error {
u := url.URL{Scheme: "ws", Host: addr, Path: path, RawQuery: rawQuery}
}
```