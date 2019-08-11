package nfs

import (
	"reflect"
	"strings"
	"testing"
)

type asserter struct {
	*testing.T
}

func assert(t *testing.T) *asserter {
	return &asserter{t}
}

func TestJsonDecode(t *testing.T) {
	decoded, err := JsonDecode(strings.NewReader(testRequest))
	assert(t).errNil(err, t)

	expected := &ShareRequest{
		SharePath: "/var/nfs/pictures",
		ExportOptions: ExportOptions{
			Clients:        []string{"192.168.1.110", "192.168.1.112"},
			Rw:             true,
			Sync:           true,
			NoSubtreeCheck: true,
		},
		MkDir:       true,
		DirOwnerUid: 33,
		DirOwnerGid: 33,
	}

	assert(t).equals(decoded, expected)
}

func TestJsonDecodeMissingRequiredProperties(t *testing.T) {
	_, err := JsonDecode(strings.NewReader(incompleteTestRequest))
	assert(t).err(err, t)

	assert(t).equals(err.Error(), "required properties of share request are missing")
}

func (a *asserter) equals(actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		a.Errorf("expected:\n%+v\ngot:\n%+v", expected, actual)
	}
}

func (a *asserter) errNil(err error, t *testing.T) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func (a *asserter) err(err error, t *testing.T) {
	if err == nil {
		t.Errorf("expected error: %v", err)
	}
}

const (
	testRequest = `{
  "sharePath": "/var/nfs/pictures",
  "exportOptions": {
    "clients": [
      "192.168.1.110",
      "192.168.1.112"
    ],
    "rw": true,
    "sync": true,
    "noSubtreeCheck": true
  },
  "mkDir": true,
  "dirOwnerUid": 33,
  "dirOwnerGid": 33
}`

	incompleteTestRequest = `{
  "exportOptions": {
    "clients": [
      "192.168.1.110",
      "192.168.1.112"
    ],
    "rw": true,
    "sync": true,
    "noSubtreeCheck": true
  },
  "mkDir": true,
  "dirOwnerUid": 33,
  "dirOwnerGid": 33
}`
)
