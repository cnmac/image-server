package httphandler

import (
	"errors"
	"fmt"
	"github.com/cnmac/image-server/common"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
	auth := r.FormValue("auth")
	if auth == "" {

	}
	uploadFile, header, e := r.FormFile("uploadFile")
	common.ErrorHandle(e, w)
	// 检查图片后缀
	ext := strings.ToLower(path.Ext(header.Filename))
	//auth := strings.ToLower(header.)
	if ext != ".jpg" && ext != ".png" {
		common.ErrorHandle(errors.New("只支持jpg/png图片上传"), w)
		return
		//defer os.Exit(2)
	}
	dirName, fileName := generateRamdomFileName()

	// 保存图片
	os.Mkdir("./uploaded/", 0777)
	saveFile, err := os.OpenFile("./uploaded/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	common.ErrorHandle(err, w)
	io.Copy(saveFile, uploadFile)

	defer uploadFile.Close()
	defer saveFile.Close()
}

func generateRamdomFileName() (string, string) {
	dirname := time.Now().UnixNano() % 500
	return string(dirname), string(time.Now().UnixNano())
}

func Download(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
