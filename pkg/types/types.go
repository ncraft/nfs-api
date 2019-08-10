package types

import (
	"bytes"
	"encoding/json"
	"errors"
	util "github.com/ncraft/go-util/pkg/base"
	"io"
	"log"
	"strings"
)

type ShareResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ShareRequest struct {
	SharePath     string        `json:"sharePath"`
	ExportOptions ExportOptions `json:"exportOptions"`
	MkDir         bool          `json:"mkDir"`
	DirOwnerUid   int           `json:"dirOwnerUid"`
	DirOwnerGid   int           `json:"dirOwnerGid"`
}

type ExportOptions struct {
	Clients        []string `json:"clients"`
	Rw             bool     `json:"rw"`
	Sync           bool     `json:"sync"`
	NoRootSquash   bool     `json:"noRootSquash"`
	NoSubtreeCheck bool     `json:"noSubtreeCheck"`
}

func (o ExportOptions) String() string {
	var buffer bytes.Buffer

	for _, client := range o.Clients {
		buffer.WriteString(client)
		buffer.WriteString("(")
		buffer.WriteString(optsStr(o))
		buffer.WriteString(")")
		buffer.WriteString(" ")
	}

	return strings.TrimSpace(buffer.String())
}

func optsStr(o ExportOptions) string {
	var buffer bytes.Buffer

	condPrependComma := func(b *bytes.Buffer) {
		if b.Len() > 0 {
			b.WriteString(",")
		}
	}

	if o.Rw {
		condPrependComma(&buffer)
		buffer.WriteString("rw")
	} else {
		condPrependComma(&buffer)
		buffer.WriteString("ro")
	}
	if o.Sync {
		condPrependComma(&buffer)
		buffer.WriteString("sync")
	} else {
		condPrependComma(&buffer)
		buffer.WriteString("async")
	}
	if o.NoRootSquash {
		condPrependComma(&buffer)
		buffer.WriteString("no_root_squash")
	}
	if o.NoSubtreeCheck {
		condPrependComma(&buffer)
		buffer.WriteString("no_subtree_check")
	}

	return strings.Trim(buffer.String(), ",")
}

func DecodeShareRequest(body io.ReadCloser) (*ShareRequest, error) {
	decoder := json.NewDecoder(body)
	var shareRequest ShareRequest
	err := decoder.Decode(&shareRequest)
	if err != nil {
		return nil, err
	}

	if missingRequiredProps(&shareRequest) {
		log.Printf("invalid share request (required properties missing): %+v", shareRequest)
		return nil, errors.New("required properties of share request are missing")
	}

	return &shareRequest, nil
}

func missingRequiredProps(r *ShareRequest) bool {
	if len(r.SharePath) == 0 || len(r.ExportOptions.Clients) == 0 || util.AnyStringEmpty(r.ExportOptions.Clients) {
		return true
	}
	return false
}
