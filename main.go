package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Data string `json:"data"`
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("static/index.html")
	router.StaticFile("/script.js", "./static/js/script.js")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	router.POST("/submit", func(c *gin.Context) {
		body := Body{}
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		w := c.Writer
		header := w.Header()
		header.Set("Content-Type", "audio/mp3")
		w.WriteHeader(http.StatusOK)
		fetchMusic(c, body.Data)
	})

	router.Run(":10001")
}

func fetchMusic(c *gin.Context, data string) {

	r, w := io.Pipe()
	defer r.Close()

	ydl := exec.Command("youtube-dl", data, "-o-")
	ffmpeg := exec.Command("ffmpeg", "-i", "/dev/stdin", "-vn", "-f", "mp3", "-")

	ydl.Stdout = w
	ydl.Stderr = os.Stderr

	ffmpeg.Stdin = r
	ffmpeg.Stdout = c.Writer
	ffmpeg.Stderr = os.Stderr
	fmt.Println("Starting-----------------------")
	go func() {
		defer w.Close()
		if err := ydl.Run(); err != nil {
			panic(err)
		}
	}()

	if err := ffmpeg.Run(); err != nil {
		panic(err)
	}
	r.Close()

	fmt.Println("Done-----------------------")
}
