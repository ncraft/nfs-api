// Package nfs provides functions to manage linux NFS shares.
package nfs

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ncraft/nfs-api/pkg/types"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
)

const defaultExportsFile = "/etc/exports"

var mu sync.Mutex

// Add creates an NFS share as specified by the share request.
func Add(request *types.ShareRequest) error {
	mu.Lock()
	defer mu.Unlock()

	currentExports, err := ioutil.ReadFile(defaultExportsFile)

	if err != nil {
		return err
	}

	if request.MkDir {
		if err := os.Mkdir(request.SharePath, 0640); err != nil {
			return err
		}
		if err := os.Chown(request.SharePath, request.DirOwnerUid, request.DirOwnerGid); err != nil {
			return err
		}
	}

	if containsExport(currentExports, request.SharePath) {
		return errors.New("export share path already exists")
	}

	newExports := currentExports
	if len(newExports) > 0 && !bytes.HasSuffix(currentExports, []byte("\n")) {
		newExports = append(newExports, '\n')
	}

	newExports = append(newExports, []byte(fmt.Sprintf("\"%s\"    %s", request.SharePath, request.ExportOptions))...)

	if err := ioutil.WriteFile(defaultExportsFile, newExports, 0644); err != nil {
		return err
	}

	if err := reloadNfsExports(); err != nil {
		if err := ioutil.WriteFile(defaultExportsFile, currentExports, 0644); err != nil {
			return err // TODO: merge both errors?
		}

		return err
	}

	return nil
}

func reloadNfsExports() error {
	cmd := exec.Command("sudo", "exportfs", "-ra")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	errorOutput, _ := ioutil.ReadAll(stderr)

	if err := cmd.Wait(); err != nil {
		return err
	}

	if len(errorOutput) > 0 {
		return fmt.Errorf("%s", errorOutput)
	}

	return nil
}

func containsExport(exports []byte, sharePath string) bool {
	return bytes.Contains(exports, []byte(sharePath))
}
