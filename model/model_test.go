package model

import (
	"FileLake/common/logger"
	"FileLake/common/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestFile_Add(t *testing.T) {
	logger.Init("dbtest")
	Init("root:123456@(172.16.124.84:9000)/filelake?charset=utf8&parseTime=true&loc=Asia%2FShanghai")
	file := File{0, utils.GetRandomShortLink(), "1,2345678", "7eff122b94897ea5b0e2a9abf47b86337fafebdc", 1588210539}
	err := file.SyncTable()
	assert.Equal(t, err, nil)

	err = file.Add()
	assert.Equal(t, err, nil)

	out, err := file.GetByServiceLink()
	assert.Equal(t, err, nil)
	assert.Equal(t, out.ServiceLink, file.ServiceLink)
	assert.Equal(t, out.AccountAddress, file.AccountAddress)
	assert.Equal(t, out.ExpireTime, file.ExpireTime)
	assert.Equal(t, out.SeaweedfsId, file.SeaweedfsId)
	logger.Debug("out= ",zap.Any("rec" ,out))

	outs, err := file.GetbyFid()

	for _,f := range outs {

		f.DeleteByServiceLink()

	}
}
