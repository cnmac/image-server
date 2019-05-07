package httphandler

import (
	"errors"
	"fmt"
	"github.com/cnmac/image-server/common"
	"github.com/cnmac/image-server/mysql"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
	auth := r.FormValue("auth")
	if auth == "" {
		//do authentication
	}
	uploadFile, header, e := r.FormFile("img")
	common.ErrorHandle(e, w)
	// 检查图片后缀
	ext := strings.ToLower(path.Ext(header.Filename))
	//auth := strings.ToLower(header.)
	if ext != ".jpg" && ext != ".png" {
		common.ErrorHandle(errors.New("只支持jpg/png图片上传"), w)
		return
	}
	if header.Size > 1*1024*1024 {
		common.ErrorHandle(errors.New("上传文件大小超过1mb限制"), w)
		return
	}
	dirName, fileName := generateRamdomFileName()
	trueName := header.Filename
	imginfo := mysql.Img{1, trueName, dirName, fileName, "", time.Now()}
	err := mysql.InsertImg(&imginfo)
	//if err != nil{
	//	common.ErrorHandle(err, w)
	//	return
	//}

	// 保存图片
	os.MkdirAll("./uploaded/"+strconv.Itoa(dirName)+"/", 0777)
	saveFile, err := os.OpenFile("./uploaded/"+strconv.Itoa(dirName)+"/"+trueName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		common.ErrorHandle(err, w)
		return
	}
	_, err = io.Copy(saveFile, uploadFile)
	if err != nil {
		w.Write([]byte("图片未保存成功！请重新上传"))
	}

	w.Write([]byte("alias:" + fileName))

	defer uploadFile.Close()
	defer saveFile.Close()
}

func generateRamdomFileName() (int, string) {
	timenow := time.Now().UnixNano()
	dirname := (time.Now().Nanosecond() / 100) % 500
	return int(dirname), strconv.FormatInt(timenow, 10)
}

func Download(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	img, err := mysql.QueryByAlias(ps.ByName("alias"))
	if err != nil {
		common.ErrorHandle(errors.New("查询mysql出差或该文件不存在"), w)
	}
	fileUrl := "./uploaded/" + strconv.Itoa(img.Dir) + "/" + img.Name
	exist, err := pathExists(fileUrl)
	if err != nil {
		common.ErrorHandle(errors.New("判断文件夹是否存在时出错"), w)
		return
	}
	if !exist {
		common.ErrorHandle(errors.New("mysql记录有该文件，而文件夹不存在该文件"), w)
	}
	http.ServeFile(w, r, fileUrl)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
