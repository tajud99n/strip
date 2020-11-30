package strip

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Code    string `json:"code"`
	DocURL  string `json:"doc_url"`
	Message string `json:"message"`
	Param   string `json:"param"`
	Type    string `json:"type"`
}

func (err Error) Error() string {
	return fmt.Sprintf("%s See %s for more information.", err.Message, err.DocURL)
}

func (err *Error) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Error struct {
			Code    string `json:"code"`
			DocURL  string `json:"doc_url"`
			Message string `json:"message"`
			Param   string `json:"param"`
			Type    string `json:"type"`
		} `json:"error"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*err = tmp.Error
	return nil
}

func (err Error) MarshalJSON() ([]byte, error) {
	var tmp struct {
		Error struct {
			Code    string `json:"code"`
			DocURL  string `json:"doc_url"`
			Message string `json:"message"`
			Param   string `json:"param"`
			Type    string `json:"type"`
		} `json:"error"`
	}

	tmp.Error = err
	return json.Marshal(tmp)
}
