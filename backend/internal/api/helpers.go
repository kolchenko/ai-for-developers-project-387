package api

import (
	"encoding/json"
	"net/http"
)

func jsonDecode(r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}
