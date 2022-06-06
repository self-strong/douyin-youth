package core

// core中工具函数

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

// GetVideoCover 获取视频的缩略图
func GetVideoCover(videoPath string) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg_go.Input(videoPath).
		Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", 5)}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		return "", err
	}

	base := filepath.Base(videoPath)
	coverPath := filepath.Join("./public/cover/", strings.Replace(base, ".mp4", ".jpeg", -1))
	err = imaging.Save(img, coverPath)
	if err != nil {
		return "", err
	}
	return filepath.Base(coverPath), nil
}
