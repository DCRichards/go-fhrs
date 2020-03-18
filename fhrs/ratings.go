package fhrs

// RatingsService encapsulates the Ratings methods of the API.
//
// https://api.ratings.food.gov.uk/help#Ratings
type RatingsService service

// Ratings is the list of possible ratings.
type Ratings struct {
	Ratings []Rating `json:"ratings"`
	Meta    Meta     `json:"meta"`
	Links   []Link   `json:"links"`
}

// Rating is a possible hygiene rating value.
type Rating struct {
	RatingID      int    `json:"ratingId"`
	RatingName    string `json:"ratingName"`
	RatingKey     string `json:"ratingKey"`
	RatingKeyName string `json:"ratingKeyName"`
	SchemeTypeID  int    `json:"schemeTypeId"`
	Links         []Link `json:"links"`
}

// Get returns the details of all possible ratings.
//
// https://api.ratings.food.gov.uk/Help/Api/GET-Ratings
func (s *RatingsService) Get() (*Ratings, error) {
	var ratings *Ratings
	if err := s.client.get("Ratings", &ratings); err != nil {
		return nil, err
	}

	return ratings, nil
}
