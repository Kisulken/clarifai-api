package clarifai_api

type ResponseStatus struct {
	Code int `json:"code"`
	Description string `json:"descripton"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	Status      ResponseStatus `json:"status"`
}

type TaggingImage struct {
	URL string `json:"url"`
}

type TaggingRequestData struct {
	Image TaggingImage `json:"image"`
}

type TaggingRequestInputs struct {
	Data TaggingRequestData `json:"data"`
}

type TaggingRequest struct {
	Inputs []TaggingRequestInputs `json:"inputs"`
}

type FeedbackForm struct {
	DocIDs           []string `json:"docids,omitempty"`
	URLs             []string `json:"url,omitempty"`
	AddTags          []string `json:"add_tags,omitempty"`
	RemoveTags       []string `json:"remove_tags,omitempty"`
	DissimilarDocIDs []string `json:"dissimilar_docids,omitempty"`
	SimilarDocIDs    []string `json:"similar_docids,omitempty"`
	SearchClick      []string `json:"search_click,omitempty"`
}

type FeedbackResponse struct {
	Code    string `json:"status_code"`
	Message string `json:"status_msg"`
}

// This structure was automatically generated
type TaggingResponse struct {
	Outputs []struct {
		CreatedAt string `json:"created_at"`
		Data      struct {
				  Concepts []struct {
					  AppID interface{} `json:"app_id"`
					  ID    string      `json:"id"`
					  Name  string      `json:"name"`
					  Value float64     `json:"value"`
				  } `json:"concepts"`
				  Regions []struct {
					  Data struct {
						       Faces []struct {
							       Identity []struct {
								       AppID interface{} `json:"app_id"`
								       ID    string      `json:"id"`
								       Name  string      `json:"name"`
								       Value float64     `json:"value"`
							       } `json:"identity"`
						       } `json:"faces"`
					       } `json:"data"`
					  RegionInfo struct {
						       BoundingBox struct {
									   BottomRow float64 `json:"bottom_row"`
									   LeftCol   float64 `json:"left_col"`
									   RightCol  float64 `json:"right_col"`
									   TopRow    float64 `json:"top_row"`
								   } `json:"bounding_box"`
					       } `json:"region_info"`
				  } `json:"regions"`
			  } `json:"data"`
		ID    string `json:"id"`
		Input struct {
				  Data struct {
					       Image struct {
							     URL string `json:"url"`
						     } `json:"image"`
				       } `json:"data"`
				  ID string `json:"id"`
			  } `json:"input"`
		Model struct {
				  AppID        interface{} `json:"app_id"`
				  CreatedAt    string      `json:"created_at"`
				  ID           string      `json:"id"`
				  ModelVersion struct {
						       CreatedAt string `json:"created_at"`
						       ID        string `json:"id"`
						       Status    struct {
									 Code        int    `json:"code"`
									 Description string `json:"description"`
								 } `json:"status"`
					       } `json:"model_version"`
				  Name       string `json:"name"`
				  OutputInfo struct {
						       Message string `json:"message"`
						       Type    string `json:"type"`
					       } `json:"output_info"`
			  } `json:"model"`
		Status struct {
				  Code        int    `json:"code"`
				  Description string `json:"description"`
			  } `json:"status"`
	} `json:"outputs"`
	Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
}