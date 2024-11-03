package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/AniComix/mpeg"
	"github.com/AniComix/query"
	"github.com/AniComix/server/storage"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	blockSize       = 8 * 1024 * 1024
	maxUploadingAge = 60 * 60 * 24
	maxUploadedAge  = 60 * 60 * 24 * 1
)

var (
	uploading = make(map[string]*Uploading)
	uploaded  = make(map[string]*UploadedFile)
)

type UploadedFile struct {
	ID       string
	Username string
	FileName string
	Time     int64
}

type uploadFileRequest struct {
	FileName   string `json:"file_name"`
	BlockCount int    `json:"block_count"`
	TotalSize  int    `json:"total_size"`
	Md5        string `json:"md5"`
}

type Uploading struct {
	ID          string
	Username    string
	FileName    string
	BlockCount  int
	TotalSize   int
	Md5         string
	BlockState  []bool
	CreateAt    int64
	SeriesID    int32
	Season      int32
	Title       string
	Description string
}

func UploadFile(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		unauthorized(c)
		return
	}
	var req uploadFileRequest
	if err := c.BindJSON(&req); err != nil {
		badRequest(c, "invalid request")
		return
	}
	if req.BlockCount < 1 || req.TotalSize < 1 {
		badRequest(c, "invalid request")
		return
	}
	if (req.BlockCount-1)*blockSize > req.TotalSize || req.BlockCount*blockSize < req.TotalSize {
		badRequest(c, "invalid request")
		return
	}
	var id string
	t := time.Now().Unix()
	hash := md5.Sum([]byte(user.Username + req.FileName + strconv.FormatInt(t, 10)))
	id = hex.EncodeToString(hash[:])
	seriesID, err := strconv.ParseInt(c.PostForm("series_id"), 10, 32)
	if err != nil {
		badRequest(c, "invalid series id")
		return
	}
	season, err := strconv.ParseInt(c.PostForm("season"), 10, 32)
	if err != nil {
		badRequest(c, "invalid season")
		return
	}
	uploading[id] = &Uploading{
		ID:          id,
		Username:    user.Username,
		FileName:    req.FileName,
		BlockCount:  req.BlockCount,
		TotalSize:   req.TotalSize,
		Md5:         req.Md5,
		BlockState:  make([]bool, req.BlockCount),
		CreateAt:    t,
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		SeriesID:    int32(seriesID),
		Season:      int32(season),
	}
	err = os.MkdirAll(filepath.Join(storage.CacheDir(), "upload", id+"uploading"), 0755)
	if err != nil {
		internalServerError(c)
		return
	}
	c.JSON(200, gin.H{
		"message":   "ok",
		"upload_id": id,
	})
}

func UploadBlock(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		unauthorized(c)
		return
	}
	id := c.DefaultQuery("id", "")
	block := c.DefaultQuery("block", "")
	task := uploading[id]
	if task == nil {
		badRequest(c, "invalid id")
		return
	}
	if task.Username != user.Username {
		unauthorized(c)
		return
	}
	i, err := strconv.Atoi(block)
	if err != nil {
		badRequest(c, "invalid block")
		return
	}
	if i < 0 || i >= task.BlockCount {
		badRequest(c, "invalid block")
		return
	}
	if task.BlockState[i] {
		badRequest(c, "block already uploaded")
		return
	}
	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		internalServerError(c)
		return
	}
	if len(bytes) > blockSize || (len(bytes) != blockSize && i != task.BlockCount-1) {
		badRequest(c, "invalid block size")
		return
	}
	err = os.WriteFile(filepath.Join(storage.CacheDir(), "upload", id+"uploading", block), bytes, 0644)
	if err != nil {
		internalServerError(c)
		return
	}
	task.BlockState[i] = true
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func FinishUpload(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		unauthorized(c)
		return
	}
	id := c.DefaultQuery("id", "")
	task := uploading[id]
	if task == nil {
		badRequest(c, "invalid id")
		return
	}
	if task.Username != user.Username {
		unauthorized(c)
		return
	}
	for _, v := range task.BlockState {
		if !v {
			badRequest(c, "not all blocks uploaded")
			return
		}
	}
	file, err := os.OpenFile(filepath.Join(storage.CacheDir(), "upload", id), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		internalServerError(c)
		return
	}
	defer func(file *os.File) {
		if file == nil {
			return
		}
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	}(file)
	for i := 0; i < task.BlockCount; i++ {
		block, err := os.Open(filepath.Join(storage.CacheDir(), "upload", id+"uploading", strconv.Itoa(i)))
		if err != nil {
			internalServerError(c)
			return
		}
		bytes, err := io.ReadAll(block)
		err = block.Close()
		if err != nil {
			internalServerError(c)
			return
		}
		_, err = file.Write(bytes)
		if err != nil {
			internalServerError(c)
			return
		}
	}
	err = os.RemoveAll(filepath.Join(storage.CacheDir(), "upload", id+"uploading"))
	if err != nil {
		internalServerError(c)
		return
	}
	uploaded[id] = &UploadedFile{
		ID:       id,
		Username: user.Username,
		FileName: task.FileName,
		Time:     time.Now().Unix(),
	}
	err = file.Close()
	if err != nil {
		internalServerError(c)
		return
	}
	filename := file.Name()
	ok := mpeg.CheckVideoFileIntegrity(filename)
	file = nil
	if !ok {
		c.JSON(404, "invalid uploaded file")
		return
	}
	series, err := query.Series.Where(query.Series.ID.Eq(uint(task.SeriesID))).First()
	if err != nil {
		internalServerError(c)
		return
	}
	outputPath := fmt.Sprintf("./data/%s/Season%d/%s.mpd", series.Title, task.Season, task.Title)
	go func(inputPath, outputPath string) {
		ok := mpeg.TransformVideoToDASHMultipleResolution(inputPath, outputPath)
		if ok {
			err := os.Remove(filename)
			if err != nil {
				log.Print(err)
			}
		}
	}(filename, outputPath)
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func GetUploadingStatus(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		unauthorized(c)
		return
	}
	id := c.DefaultQuery("id", "")
	task := uploading[id]
	if task == nil {
		badRequest(c, "invalid id")
		return
	}
	if task.Username != user.Username {
		unauthorized(c)
		return
	}
	statusBuilder := strings.Builder{}
	for _, v := range task.BlockState {
		if v {
			statusBuilder.WriteString("1")
		} else {
			statusBuilder.WriteString("0")
		}
	}
	c.JSON(200, gin.H{
		"message":     "ok",
		"status":      statusBuilder.String(),
		"create_at":   task.CreateAt,
		"block_count": task.BlockCount,
		"total_size":  task.TotalSize,
	})
}

func CancelUploadTask(c *gin.Context) {
	user, err := getCurrentUser(c)
	if err != nil {
		unauthorized(c)
		return
	}
	id := c.DefaultQuery("id", "")
	task := uploading[id]
	if task == nil {
		badRequest(c, "invalid id")
		return
	}
	if task.Username != user.Username {
		unauthorized(c)
		return
	}
	err = os.RemoveAll(filepath.Join(storage.CacheDir(), "upload", id+"uploading"))
	if err != nil {
		internalServerError(c)
		return
	}
	delete(uploading, id)
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func checkExpiredTasks() {
	for k, v := range uploading {
		if time.Now().Unix()-v.CreateAt > maxUploadingAge {
			err := os.RemoveAll(filepath.Join(storage.CacheDir(), "upload", k+"uploading"))
			if err != nil {
				log.Print(err)
			}
			delete(uploading, k)
		}
	}
	for k, v := range uploaded {
		if time.Now().Unix()-v.Time > maxUploadedAge {
			err := os.Remove(filepath.Join(storage.CacheDir(), "upload", k))
			if err != nil {
				log.Print(err)
			}
			delete(uploaded, k)
		}
	}
}

func StartUploadTaskCleaner() {
	go func() {
		for {
			checkExpiredTasks()
			time.Sleep(2 * time.Hour)
		}
	}()
}
