package entity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectToJSON(t *testing.T) {
	order := Order{UUID: "ci-test", Destination: "1"}
	want := []byte(`{"order":"ci-test","destination":"1"}`)
	got, _ := json.Marshal(order)
	assert.Equal(t, string(got), string(want))
}
