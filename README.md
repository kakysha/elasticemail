Partially implemented ElasticEmail API in Go. Currently implemented all Email methods and methods to work with Subaccounts. Methods are named identical to the API: https://api.elasticemail.com/public/help#Email_header.

### Installing

##### using `go dep`to put into project/vendor folder:
1) in your .go source file import this lib
```go
import ee "github.com/kakysha/elasticemail"
```
2) then run in your project folder
```
go dep ensure
```

##### using default `go get` way to put into $GOPATH folder
```
go get -u github.com/kakysha/elasticemail
```

### Documentation

You can get quick overview of the code by looking at godoc page: https://godoc.org/github.com/kakysha/elasticemail

All methods are provided in two variants: `Method` and `MethodContext`, latter has the first argument of an instance of `context.Context`.

You can set `apikey` per every request by providing different values in context's key `'apikey'` (or in query params directly invoking HTTPGet/Post methods):

```go
c := context.WithValue(c, "apikey", s.ProviderAPIKey)
res := EE.SendContext(c, eeMail)
```

### Example

1. Create ElasticEmail API Client instance:
```go
var (
	EE  ee.Client
	res *ee.Response
)
EE.Init(&ee.Config{APIKey: "dca6c423-4251-0225-cc22-ade2d28168db"}) // your API key
```
2. Email API:
```go
var eeMail = &ee.Email{
	From:          "sender@example.com",
	FromName:      "EE Test",
	To:            "my@me.com",
	BodyHTML:      "<h1>It works!</h1>",
	BodyText:      "It works!",
	Subject:       "Assalamu Alaykum",
	CustomHeaders: map[string]string{"myheader1": "heavervalue1"},
}

res = EE.Send(eeMail)
log.Printf("Send\tSuccess:%t\tError:%v\tData:%v", res.Success, res.Error, res.Data)
message := res.Data.(map[string]interface{})

res = EE.Status(message["messageid"].(string))
log.Printf("Status\tSuccess:%t\tError:%v\tData:%v", res.Success, res.Error, res.Data)

res = EE.View(message["messageid"].(string))
log.Printf("View\tSuccess:%t\tError:%v\tData:%v", res.Success, res.Error, res.Data)
```

3. Subaccounts API:
```go
res = EE.AddSubAccount(&ee.Subaccount{Email: "test@example.com", Password: "123"})
log.Printf("AddSubAccount\tSuccess:%t\tError:%v\tData:%v", res.Success, res.Error, res.Data)

res = EE.UpdateSubAccountSettings(&ee.Subaccount{DailySendLimit: 10}, map[string]string{"subAccountEmail": "test@example.com"})
log.Printf("UpdateSubAccountSettings\tSuccess:%t\tError:%v\tData:%v", res.Success, res.Error, res.Data)

res = EE.GetSubAccountAPIKey(map[string]string{"subAccountEmail": "test@example.com"})
log.Printf("GetSubAccountApiKey\tSuccess:%t\tError:%v\tData:%v", res.Success, res.Error, res.Data)

res = EE.DeleteSubAccount(map[string]string{"subAccountEmail": "test@example.com", "notify": "false"})
log.Printf("DeleteSubAccount\tSuccess:%t\tError:%v\tData:%v", res.Success, res.Error, res.Data)
```

## License

This project is licensed under the MIT License.
