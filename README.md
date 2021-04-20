## supervisor-client

supervisor client by golang

## examples

### stop process and get process info

```go
package main

import (
	"fmt"
	sc "github.com/lixianyang/supervisor-client"
)

func main() {
	client, err := sc.New("http://user:123@127.0.0.1:9001", nil)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	err = client.StopProcess("web", true)
	if err != nil {
		panic(err)
	}
	info, err := client.GetProcessInfo("web")
	if err != nil {
		panic(err)
	}
	fmt.Print(info)
	return
}
```
