/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package restclient

import (
	"testing"
)

func TestValidatesHostParameter(t *testing.T) {
	testCases := []struct {
		Host    string
		APIPath string

		URL string
		Err bool
	}{
		{"127.0.0.1", "", "http://127.0.0.1/", false},
		{"127.0.0.1:8080", "", "http://127.0.0.1:8080/", false},
		{"foo.bar.com", "", "http://foo.bar.com/", false},
		{"http://host/prefix", "", "http://host/prefix/", false},
		{"http://host", "", "http://host/", false},
		{"http://host", "/", "http://host/", false},
		{"http://host", "/other", "http://host/other/", false},
		{"host/server", "", "", true},
	}
	for i, testCase := range testCases {
		u, err := DefaultServerURL(testCase.Host, testCase.APIPath, false)
		switch {
		case err == nil && testCase.Err:
			t.Errorf("expected error but was nil")
			continue
		case err != nil && !testCase.Err:
			t.Errorf("unexpected error %v", err)
			continue
		case err != nil:
			continue
		}
		if e, a := testCase.URL, u.String(); e != a {
			t.Errorf("%d: expected host %s, got %s", i, e, a)
			continue
		}
	}
}
