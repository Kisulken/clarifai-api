package clarifai_api

import (
	"encoding/json"
	"errors"
)

const (
	min_accuracy = 0.7 // entities with accuracy lower than 70% will be ignored
	max_results = "10" // per request
)

type Tag struct {
	Value 	 string
	Accuracy float64
}

type Tags struct {
	Success bool
	Faces bool
	Concepts []Tag
	Persons []Tag
}

// converts a default response into the simpler format
func responseToTags(res TaggingResponse) *Tags {
	tags := new(Tags)

	faceDetected := false
	identified := false

	if len(res.Outputs) > 0 {
		output := res.Outputs[0]
		if output.Data.Concepts != nil {
			concepts := output.Data.Concepts
			tags.Concepts = make([]Tag, len(concepts))
			for i, c := range concepts {
				tags.Concepts[i] = Tag{c.Name, c.Value}
			}
		}
		if output.Data.Regions != nil {
			regions := output.Data.Regions
			tags.Persons = make([]Tag, len(regions))
			for i, r := range regions {
				identified = false
				if r.Data.Faces != nil && len(r.Data.Faces) > 0 {
					faces := r.Data.Faces[0]
					for _, identity := range faces.Identity {
						if identity.Value >= min_accuracy {
							tags.Persons[i] = Tag{identity.Name, identity.Value}
							faceDetected = true
							identified = true
							break
						}
					}
				}
				if !identified {
					tags.Persons[i] = Tag{"unknown", 0.0}
				}
			}
		}
	}

	tags.Success = faceDetected || len(tags.Concepts) > 0
	tags.Faces = faceDetected

	return tags
}

// recognizing the image and returns tags
func (client *Client) GetTags(url, model string) (*Tags, error) {
	img := TaggingImage{url}
	inputs := make([]TaggingRequestInputs, 1)
	inputs[0].Data = TaggingRequestData{img}
	tagReq := TaggingRequest{inputs}

	body, err := json.Marshal(tagReq)
	if err != nil {
		return nil, errors.New("Error during marhaling the response due " + err.Error())
	}

	result, err := client.CustomRequest(rootURLv2, "models/" + model + "/outputs?per_page=" + max_results, "POST", body)
	if err != nil {
		if err.Error() == "TOKEN_INVALID" {
			err = client.requestAccessToken()
			if err != nil {
				return nil, err
			}
			result, err = client.CustomRequest(rootURLv2, "models/" + model + "/outputs?per_page=" + max_results, "POST", body)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var response TaggingResponse

	err = json.Unmarshal(result, &response)
	if err != nil {
		return nil, errors.New("Error during unmarhaling the response due " + err.Error())
	}

	return responseToTags(response), nil
}

// sends feedback regarding recognized images in order to improve recognition
func (client *Client) Feedback(form FeedbackForm) (error) {
	if form.DocIDs == nil && form.URLs == nil {
		return errors.New("Requires at least one DocID or url")
	}

	if form.DocIDs != nil && form.URLs != nil {
		return errors.New("Request must provide exactly one of the following fields: {'DocIDs', 'URLs'}")
	}

	encodedForm, err := json.Marshal(form)
	if err != nil {
		return err
	}

	result, err := client.CustomRequest(rootURLv1, "feedback", "POST", encodedForm)

	response := new(FeedbackResponse)
	err = json.Unmarshal(result, response)
	if err != nil {
		return err
	}

	if response.Code != "OK" {
		return errors.New(response.Message)
	}

	return nil
}