# clarifai-api
Client which provides basic functionality of clarifai.com API.

## Installation
```bash
go get github.com/Kisulken/clarifai-api
```

## Example
```Go
import (
	clarifai "clarifai-api"
	"fmt"
)

const (
	client_id = "<your client id>"
	client_secret = "<your client secret>"

	test_url = "any image url here"
)

func main()  {
	// create a new client with your id and secret
	client := clarifai.NewClient(client_id, client_secret)
	// retrive tags from provided image using general model
	tags, err := client.GetTags(test_url, clarifai.GeneralModelID)
	if err != nil {
		panic(err)
	}
	fmt.Println(tags)
	// give a feedback about recgonized image
	err = client.Feedback(clarifai.FeedbackForm{
		URLs: []string{test_url},
		AddTags: []string{"tag to add"},
		RemoveTags: []string{"tag to remove"},
	})
	if err != nil {
		panic(err)
	}
}
```
