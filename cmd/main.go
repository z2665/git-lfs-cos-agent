package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/z2665/git-lfs-cos-agent/pkg/config"
	"github.com/z2665/git-lfs-cos-agent/pkg/cos"
	"github.com/z2665/git-lfs-cos-agent/pkg/types"
	"github.com/z2665/git-lfs-cos-agent/pkg/utils"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	errWriter := bufio.NewWriter(os.Stderr)
	writeToStderr := utils.BuildWriteToStderr(errWriter)

	var confpath *string
	if len(os.Args) >= 2 {
		confpath = &os.Args[1]
	}
	conf, err := config.LoadConfig(confpath)
	if err != nil {
		writeToStderr(fmt.Sprintf("not found conf in %s", *confpath))
		os.Exit(-255)
	}
	cli := cos.NewCos(writer, errWriter)
	for scanner.Scan() {
		line := scanner.Text()
		var req types.Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			writeToStderr(fmt.Sprintf("Unable to parse request: %v\n", line))
			continue
		}
		switch req.Event {
		case "init":
			writeToStderr(fmt.Sprintf("Initialising rsync agent for: %s\n", req.Operation))
			os.MkdirAll(conf.Tmpdir, 0755)
			cli.Init(&conf, req.Remote)
		case "download":
			writeToStderr(fmt.Sprintf("Received download request for: %s\n", req.Oid))
			cli.Download(req.Oid, req.Size, req.Action)
		case "upload":
			writeToStderr(fmt.Sprintf("Received upload request for: %s\n", req.Oid))
			cli.Upload(req.Oid, req.Size, req.Action, req.Path)
		case "terminate":
			writeToStderr("Terminating rsync agent gracefully.\n")
			// clean up and terminate. No response is expected.
			cli.Destory()
			os.Exit(0)
		}
	}
}
