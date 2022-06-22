package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/z2665/git-lfs-cos-agent/pkg/types"
)

func BuildWriteToStderr(errWriter *bufio.Writer) func(msg string) {
	return func(msg string) {
		if !strings.HasSuffix(msg, "\n") {
			msg = msg + "\n"
		}
		errWriter.WriteString(msg)
		errWriter.Flush()
	}
}

func BuildSendResponse(writer *bufio.Writer, errWriter func(msg string)) func(r interface{}) error {
	return func(r interface{}) error {
		b, err := json.Marshal(r)
		if err != nil {
			return err
		}
		// Line oriented JSON
		b = append(b, '\n')
		_, err = writer.Write(b)
		if err != nil {
			return err
		}
		writer.Flush()
		errWriter(fmt.Sprintf("Sent message %v", string(b)))
		return nil
	}

}

func SendTransferError(oid string, code int, message string, writer func(r interface{}) error, errWriter func(msg string)) {
	resp := &types.TransferResponse{
		Event: "complete",
		Oid:   oid,
		Path:  "",
		Error: &types.OperationError{Code: code, Message: message},
	}
	err := writer(resp)
	if err != nil {
		errWriter(fmt.Sprintf("Unable to send transfer error: %v\n", err))
	}
}

func SendProgress(oid string, bytesSoFar int64, bytesSinceLast int, writer func(r interface{}) error, errWriter func(msg string)) {
	resp := &types.ProgressResponse{
		Event:          "progress",
		Oid:            oid,
		BytesSoFar:     bytesSoFar,
		BytesSinceLast: bytesSinceLast,
	}
	err := writer(resp)
	if err != nil {
		errWriter(fmt.Sprintf("Unable to send progress update: %v\n", err))
	}
}
func SendInitError(code int, message string, writer func(r interface{}) error) {
	resp := &types.InitResponse{Error: &types.OperationError{Code: code, Message: message}}
	writer(resp)
}
