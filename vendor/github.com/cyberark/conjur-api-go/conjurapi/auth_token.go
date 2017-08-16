package conjurapi

import (
	"time"
	"encoding/json"
	"encoding/base64"
)

type AuthnToken struct {
	Data          string `json:"data"`
	Timestamp     time.Time
	Base64encoded string
	Signature 	  string `json:"signature"`
	Key       	  string `json:"key"`
}

func (t *AuthnToken) UnmarshalJSON(data []byte) error {
	type Alias AuthnToken
	var (
		err error
		timestamp time.Time
	)
	aux := struct {
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	timestamp, err = time.Parse("2006-01-02 15:04:05 MST", aux.Timestamp)
	if err != nil {
		return err
	}

	t.Timestamp = timestamp
	t.Base64encoded = base64.StdEncoding.EncodeToString(data)
	return nil
}

func (t *AuthnToken) ValidAtTime(refTime time.Time) bool {
	return refTime.Before(t.Timestamp.Add(5 * time.Minute))
}
