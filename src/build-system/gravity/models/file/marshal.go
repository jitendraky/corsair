package file

import "encoding/json"

// Encode dumps an File to json.
func (file File) Encode() ([]byte, error) {
	return json.Marshal(file)
}

// Decode loads an File from json
func (file *File) Decode(data []byte) error {
	return json.Unmarshal(data, file)
}
