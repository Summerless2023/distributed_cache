package test

import (
	"main/conf"
	"main/src/models"
	"testing"
)

func Test1(t *testing.T) {
	sc := models.NewCacheStorage(conf.DEFAULT_MAX_BYTES)
	t.Log(sc.GetMaxBytes())
	t.Log(sc.GetNbytes())
	sc.AddNBytes(100)
	t.Log(sc.GetNbytes())
}
