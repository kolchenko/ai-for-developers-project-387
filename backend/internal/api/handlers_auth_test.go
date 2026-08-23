package api

import (
	"net/http"
	"testing"
)

func TestAdminLogin(t *testing.T) {
	ts := newTestServer(t)

	cases := []struct {
		name    string
		payload string
		want    int
	}{
		{"success", `{"username":"admin","password":"admin"}`, http.StatusOK},
		{"wrong password", `{"username":"admin","password":"nope"}`, http.StatusUnauthorized},
		{"unknown user", `{"username":"root","password":"admin"}`, http.StatusUnauthorized},
		{"empty username", `{"username":"","password":"admin"}`, http.StatusBadRequest},
		{"empty password", `{"username":"admin","password":""}`, http.StatusBadRequest},
		{"malformed json", `{`, http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := doJSON(t, http.MethodPost, ts.URL+"/admin/login", tc.payload)
			if status != tc.want {
				t.Fatalf("status %d body %s, want %d", status, body, tc.want)
			}
		})
	}

	status, body := doJSON(t, http.MethodPost, ts.URL+"/admin/login", `{"username":"admin","password":"admin"}`)
	if status != http.StatusOK {
		t.Fatalf("login: status %d body %s", status, body)
	}
	res := decode[LoginResponse](t, body)
	if res.Username != "admin" {
		t.Fatalf("unexpected response: %+v", res)
	}
}
