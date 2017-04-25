# kochava reporting client

This package provides basic (see TODOs) access to the Kochava reporting system. Full documentation of API capabilities found here: http://support.kochava.com

Sign up/sign in here: https://kochava.com

## Setup/Usage

```go
    go get github.com/ericdmann/kochava-reporting
```
Examples are located inside of `/examples`

### Create a client

```go
  kClient, err := kochava.NewClient(apiKey, appGUID)
```

### Retrieving available columns
A new client pre-fetches all available columns, which are made available in `.Templates`

```go
  for _, template := range kClient.Templates {
    fmt.Println(template.ReportType, " - ", template.ReportComponent)
  }
```

### Requesting a detail report
The below requests a click report for the last 24 hours (UTC), to be emailed in JSON format.

```go
  reportRequest, err := kClient.NewDetailRequest("click",
    "json",
    []string{"eric@ericmann.net"},
    time.Now().Add(-time.Hour*24),
    time.Now(),
    "UTC")
```

### Checking a request
Each report request (schedule or one-time) will return a report token which may be queried.

```go
  repResponse, err := kClient.CheckRequest(token)
  fmt.Println("["+time.Now().Local().String()+"] Status update: ", repResponse.Status, "Progress: ", repResponse.Progress)
```

# Retrieving a report
Once your report has been generated a link to retrieve the report is returned.
```go
	fmt.Println("Report link: ", repResponse.ReportLink)
```

### Example
Publisher and subscription examples are located inside of `/examples`

```go
	go run retrieve.go
```

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D


## TODO
  1. Testing
  2. Add recurring report submissions
  3. Add additional documentation/implementation for groupings
	4. For JSON requests, unmarshal into structs
