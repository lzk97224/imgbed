package utils

import (
	"testing"
	"time"
)

func Test_createName(t *testing.T) {
	t.Log(CreateName(time.Now(), "123.jpg", 1))
}

func Test_createPath(t *testing.T) {
	t.Log(CreatePath("", time.Now()))
	t.Log(CreatePath("lizhikui", time.Now()))
}

func Test_createFullName(t *testing.T) {
	now := time.Now()
	t.Log(CreateFullName("", now, "sdf.jpg", 1))
	t.Log(CreateFullName("lihzikui", now, "sdf.jpg", 1))
}
