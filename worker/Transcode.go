package main

import (
	"../common"
	"github.com/google/uuid"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Transcode(job_id string) error {
	var video common.Video

	db := common.Init()
	db.First(&video, common.Video{JobId: uuid.MustParse(job_id)})
	db.Model(&video).UpdateColumn("status", "PENDING")
	out_dir := filepath.Join("/Users/bseenu/Downloads", job_id)
	m3u8 := filepath.Join(out_dir, "playlist.m3u8")
	dest := filepath.Join(out_dir, "%d.ts")
	if _, err := os.Stat(out_dir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(out_dir, 0777)
		}
	}
	time_start := time.Now()
	out, err := exec.Command("/usr/local/bin/ffmpeg", "-y", "-i", video.Input, "-profile:v", "baseline", "-level", "3", "-s", "1280x720", "-vbsf", "h264_mp4toannexb", "-c:v", "libx264", "-c:a", "aac", "-f", "segment", "-segment_start_number", "0", "-output_ts_offset", "0", "-segment_time", "5", "-segment_time_delta", "0.05", "-segment_format", "mpegts", "-segment_list", m3u8, dest).CombinedOutput()
	time_end := time.Now()
	if err != nil {
		db.Model(&video).UpdateColumn("status", "FAILED")
		db.Model(&video).UpdateColumn("completed", time_end)
		log.Printf("Video transcoding failed - %s", out)
		return err
	} else {
		log.Printf("Video successfully transcoded, it too %g seconds to transcode it", time_end.Sub(time_start).Seconds())
		db.Model(&video).UpdateColumn("status", "SUCCESSFUL")
		db.Model(&video).UpdateColumn("completed", time_end)
		db.Model(&video).UpdateColumn("output", m3u8)
		return nil
	}
}
