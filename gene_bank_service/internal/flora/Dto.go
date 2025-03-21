package flora

type Flora struct {
	ID             string `json:"id,omitempty"`              // Unique identifier for the plant
	CommonName     string `json:"common_name,omitempty"`     // Common name of the plant
	ScientificName string `json:"scientific_name,omitempty"` // Scientific name of the plant
	Image          []byte `json:"image,omitempty"`           // Image data (bytes)
	Description    string `json:"description,omitempty"`     // Description of the plant
	Origin         string `json:"origin,omitempty"`          // Origin of the plant
	OtherDetails   string `json:"other_details,omitempty"`   // Additional details about the plant
	Type           string `json:"type,omitempty"`            // Type of post
}

type FloraResponse struct {
	Flora []Flora `json:"flora,omitempty"`
}

type FloraRequest struct {
	CommonName     string `json:"common_name,omitempty"`     // Common name of the plant
	ScientificName string `json:"scientific_name,omitempty"` // Scientific name of the plant
	ImageURL       string `json:"image_url,omitempty"`       // Image URL or file reference
	ImagePath      string `json:"image_path,omitempty"`      // Image URL or file reference
	Description    string `json:"description,omitempty"`     // Description of the plant
	Origin         string `json:"origin,omitempty"`          // Origin of the plant
	OtherDetails   string `json:"other_details,omitempty"`   // Additional details about the plant
	Type           string `json:"type,omitempty"`            // Type of post
}
