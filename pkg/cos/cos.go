package cos

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/z2665/git-lfs-cos-agent/pkg/config"
	"github.com/z2665/git-lfs-cos-agent/pkg/types"
	"github.com/z2665/git-lfs-cos-agent/pkg/utils"
)

//support tencent cloud cos
type Cos struct {
	cli        *cos.Client
	tmpdir     string
	remotepath string
	writer     func(r interface{}) error
	errWriter  func(msg string)
}

func NewCos(writer, errWriter *bufio.Writer) types.Client {
	funcerrWriter := utils.BuildWriteToStderr(errWriter)
	tmp := &Cos{
		writer:    utils.BuildSendResponse(writer, funcerrWriter),
		errWriter: funcerrWriter,
	}
	return tmp
}

func (c *Cos) Init(config *config.Config, remotepath string) error {
	//Init cos client
	u, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", config.BucketName, config.Region))
	if err != nil {
		utils.SendInitError(-1, err.Error(), c.writer)
	}
	b := &cos.BaseURL{BucketURL: u}
	cli := cos.NewClient(b, &http.Client{
		//set timeout
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			//ak/sk
			SecretID:  config.SecretID,
			SecretKey: config.SecretKey,
		},
	})
	c.cli = cli
	c.remotepath = remotepath
	c.tmpdir = config.Tmpdir
	// Success!
	resp := &types.InitResponse{}
	return c.writer(resp)
}

func (c *Cos) Download(oid string, size int64, a *types.Action) {
	dlFile, err := ioutil.TempFile(c.tmpdir, "cos-")
	if err != nil {
		utils.SendTransferError(oid, -1, err.Error(), c.writer, c.errWriter)
		return
	}
	defer dlFile.Close()
	dlfilename := dlFile.Name()
	resp, err := c.cli.Object.Get(context.Background(), c.remoteFile(oid), nil)
	if err != nil {
		utils.SendTransferError(oid, -1, err.Error(), c.writer, c.errWriter)
		return
	}
	defer resp.Body.Close()
	io.Copy(dlFile, resp.Body)

	complete := &types.TransferResponse{Event: "complete", Oid: oid, Path: dlfilename, Error: nil}
	err = c.writer(complete)
	if err != nil {
		c.errWriter(fmt.Sprintf("Unable to send completion message: %v\n", err))
	}
}

func (c *Cos) Upload(oid string, size int64, a *types.Action, fromPath string) {
	f, err := os.Open(fromPath)
	if err != nil {
		utils.SendTransferError(oid, -1, err.Error(), c.writer, c.errWriter)
		return
	}
	s, err := f.Stat()
	if err != nil {
		utils.SendTransferError(oid, -1, err.Error(), c.writer, c.errWriter)
		return
	}
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength: s.Size(),
		},
	}
	_, err = c.cli.Object.Put(context.Background(), c.remoteFile(oid), f, opt)
	if err != nil {
		utils.SendTransferError(oid, -1, err.Error(), c.writer, c.errWriter)
		return
	}

	complete := &types.TransferResponse{Event: "complete", Oid: oid, Path: "", Error: nil}
	if err := c.writer(complete); err != nil {
		c.errWriter(fmt.Sprintf("Unable to send completion message: %v\n", err))
	}
}

func (c *Cos) Destory() {

}

func (c *Cos) remoteFile(oid string) string {
	return fmt.Sprintf("%s/%s", c.remotepath, oid)
}
