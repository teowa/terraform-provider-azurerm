package azuresdkhacks

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/encodings"
)

// Hacking the SDK model

var _ encodings.Codec = PngImage{}

type PngImage struct {
	KeyFrameInterval *string                  `json:"keyFrameInterval,omitempty"`
	Layers           *[]PngLayer              `json:"layers,omitempty"`
	Range            *string                  `json:"range,omitempty"`
	Start            string                   `json:"start"`
	Step             *string                  `json:"step,omitempty"`
	StretchMode      *encodings.StretchMode   `json:"stretchMode,omitempty"`
	SyncMode         *encodings.VideoSyncMode `json:"syncMode,omitempty"`

	// Fields inherited from Codec
	Label *string `json:"label,omitempty"`
}

var _ json.Marshaler = PngImage{}

func (s PngImage) MarshalJSON() ([]byte, error) {
	type wrapper PngImage
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PngImage: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PngImage: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.PngImage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PngImage: %+v", err)
	}

	return encoded, nil
}