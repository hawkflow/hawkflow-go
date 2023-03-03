![HawkFLow.ai](https://hawkflow.ai/static/images/emails/bars.png)

# HawkFlow.ai

## Monitoring for anyone that writes code

Usage

```go
package main

import (
	"fmt"
	"github.com/hawkflow/hawkflow-go"
	"time"
)

func main() {
	// Authenticate with your API key
	hf := hawkflow.New("YOUR_API_KEY")

	// Start timing your code - pass through process (required) and meta (optional) parameters
	fmt.Println("Sending timing start data to hawkflow")
	err := hf.Start("hawkflow_examples", "your_meta_data", "")
	if err != nil {
		fmt.Println(err)
	}

	// Do some work
	fmt.Println("Sleeping for 5 seconds...")
	time.Sleep(5 * time.Second)

	// End timing this piece of code - process (required) and meta (optional) parameters should match the start
	fmt.Println("Sending timing end data to hawkflow")
	err = hf.End("hawkflow_examples", "your_meta_data", "")
	if err != nil {
		fmt.Println(err)
	}
}
```

More examples: [HawkFlow.ai Go examples](https://github.com/hawkflow/hawkflow-examples/tree/master/go)

Read the docs: [HawkFlow.ai documentation](http://docs.hawkflow.ai/)

## What is HawkFlow.ai?

HawkFlow.ai is a new monitoring platform that makes it easier than ever to make monitoring part of your development
process. Whether you are an Engineer, a Data Scientist, an Analyst, or anyone else that writes code, HawkFlow.ai helps
you and your team take ownership of monitoring.