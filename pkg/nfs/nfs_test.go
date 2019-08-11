package nfs

import (
	"testing"
)

func TestExportOptionsStringerNoRootSquash(t *testing.T) {
	opts := ExportOptions{
		Clients:        []string{"181.128.88.50", "184.121.78.14"},
		Rw:             true,
		Sync:           true,
		NoRootSquash:   true,
		NoSubtreeCheck: true,
	}
	optsStr := opts.String()

	if optsStr != "181.128.88.50(rw,sync,no_root_squash,no_subtree_check) 184.121.78.14(rw,sync,no_root_squash,no_subtree_check)" {
		t.Errorf("%s", optsStr)
	}
}

func TestExportOptionsStringerRootSquash(t *testing.T) {
	opts := ExportOptions{
		Clients:        []string{"181.128.88.50", "184.121.78.14"},
		Rw:             true,
		Sync:           true,
		NoRootSquash:   false,
		NoSubtreeCheck: true,
	}
	optsStr := opts.String()

	if optsStr != "181.128.88.50(rw,sync,no_subtree_check) 184.121.78.14(rw,sync,no_subtree_check)" {
		t.Errorf("%s", optsStr)
	}
}
