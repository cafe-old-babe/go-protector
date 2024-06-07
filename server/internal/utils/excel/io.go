package excel

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go-protector/server/internal/config"
	"go-protector/server/internal/custom/c_result"
	"go-protector/server/internal/utils"
	"os"
	"path/filepath"
)

// Write 写入模板路径
// dst 选填:
// 无效: 使用配置路径,使用随机数补全文件名称
// 有效 无文件路径:使用随机数补全文件名称,有文件路径 os.Create
func Write(file *excelize.File, dst ...string) (fp string, err error) {
	if len(dst) <= 0 {
		fp = config.GetConfig().Server.TempPath

	} else {
		fp = filepath.Join(dst...)
	}
	if len(fp) <= 0 {
		err = errors.New("写入文件失败,无目录信息")
		return
	}
	if fp, err = filepath.Abs(fp); err != nil {
		return
	}
	ext := filepath.Ext(fp)
	if len(ext) <= 0 {
		fp = filepath.Join(fp, utils.GetNextIdStr()+".xlsx")
	}

	systemFile, err := os.Create(fp)
	if err != nil {
		return
	}
	defer systemFile.Close()
	if err = file.Write(systemFile); err != nil {
		return
	}
	return
}

// Export 导出至页面
func Export(c *gin.Context, file *excelize.File, fileName string) (err error) {
	var filePath string
	if filePath, err = Write(file); err != nil {
		return
	}
	c_result.ExportFile(c, filePath, fileName)

	return
}
