package gobetafaceapi

type DetectionFlags struct {
	Bestface    bool
	Basicpoints bool
	Propoints   bool
	Classifiers bool
	Extended    bool
}

// ************************************************************************************************
// ** Error
// ************************************************************************************************
// ** Error response (500 status code)
// ************************************************************************************************

type ErrorResponse struct {
	Code        int    `json:"error_code"`
	Description string `json:"error_description"`
}

// ************************************************************************************************
// ** DetectionFlags
// ************************************************************************************************
// ** To be used as a SendMedia parameter
// ************************************************************************************************

func (st DetectionFlags) String() string {
	ret := getValueOrEmpty(st.Bestface, "bestface,") +
		getValueOrEmpty(st.Basicpoints, "basicpoints,") +
		getValueOrEmpty(st.Propoints, "propoints,") +
		getValueOrEmpty(st.Classifiers, "classifiers") +
		getValueOrEmpty(st.Extended, "extended,") +
		"content"

	return ret
}

func getValueOrEmpty(condition bool, value string) string {
	if condition {
		return value
	}

	return ""
}

// ************************************************************************************************
// ** SendMedia
// ************************************************************************************************
// ** SendMedia request
// ************************************************************************************************

type SendMediaRequest struct {
	ApiKey           string   `json:"api_key"`
	FileURI          string   `json:"file_uri"`
	DetectionFlags   string   `json:"detection_flags"`
	RecognizeTargets []string `json:"recognize_targets"`
	OriginalFilename string   `json:"original_filename"`
}

// ************************************************************************************************
// ** SendMedia
// ************************************************************************************************
// ** SendMedia response
// ************************************************************************************************

type ExtendedMedia struct {
	Media     Media `json:"media"`
	Recognize struct {
		RecognizeUUID string `json:"recognize_uuid"`
		Results       []struct {
			FaceUUID string `json:"face_uuid"`
			Matches  []struct {
				FaceUUID   string  `json:"face_uuid"`
				Confidence float64 `json:"confidence"`
				IsMatch    bool    `json:"is_match"`
				PersonID   string  `json:"person_id"`
			} `json:"matches"`
		} `json:"results"`
	} `json:"recognize"`
}

type Media struct {
	MediaUUID string `json:"media_uuid"`
	Checksum  string `json:"checksum"`
	Faces     []struct {
		FaceUUID       string  `json:"face_uuid"`
		MediaUUID      string  `json:"media_uuid"`
		X              float64 `json:"x"`
		Y              float64 `json:"y"`
		Width          float64 `json:"width"`
		Height         float64 `json:"height"`
		Angle          float64 `json:"angle"`
		DetectionScore float64 `json:"detection_score"`
		Points         []struct {
			X    float64 `json:"x"`
			Y    float64 `json:"y"`
			Type int     `json:"type"`
			Name string  `json:"name"`
		} `json:"points"`
		UserPoints interface{} `json:"user_points"`
		Tags       []struct {
			Name       string  `json:"name"`
			Value      string  `json:"value"`
			Confidence float64 `json:"confidence"`
		} `json:"tags"`
		PersonID     string `json:"person_id"`
		AppearanceID int    `json:"appearance_id"`
		Start        string `json:"start"`
		Duration     string `json:"duration"`
	} `json:"faces"`
	Tags []struct {
		Name       string  `json:"name"`
		Value      string  `json:"value"`
		Confidence float64 `json:"confidence"`
		X          float64 `json:"x"`
		Y          float64 `json:"y"`
		Width      float64 `json:"width"`
		Height     float64 `json:"height"`
		Angle      float64 `json:"angle"`
		InstanceID int     `json:"instance_id"`
		Start      string  `json:"start"`
		Duration   string  `json:"duration"`
	} `json:"tags"`
	OriginalFilename string `json:"original_filename"`
	Duration         string `json:"duration"`
}
