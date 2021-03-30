package request

import "testing"

func TestRepository(t *testing.T) {
	r := &Record{}

	Persist(r)
}
