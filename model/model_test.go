package model

import (
	"FileLake/common/utils"
	"testing"
)

func TestFile_Add(t *testing.T) {
	file := File{utils.GetRandomShortLink(), "1,2345678", "7eff122b94897ea5b0e2a9abf47b86337fafebdc", 1588210539}
	file.Add()
}
