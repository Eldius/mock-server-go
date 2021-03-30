package request

import (
	"encoding/json"
)

func Serialize(r *Record) ([]byte, error) {
	return json.Marshal(r)
}

func Deserialize(msg []byte) (*Record, error) {
	var r *Record
	err := json.Unmarshal(msg, &r)
	return r, err
}
