package test

import (
	"main/conf"
	"main/src/models"
	"testing"
)

func Test1(t *testing.T) {
	sc := models.NewCacheStorage(conf.Default_Max_Bytes)
	t.Log(sc.GetMaxBytes())
	t.Log(sc.GetNbytes())
	sc.AddNBytes(100)
	t.Log(sc.GetNbytes())
}
