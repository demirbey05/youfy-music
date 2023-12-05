package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Data string `json:"data"`
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("static/index.html")
	router.StaticFile("/script.js", "./static/js/script.js")
	router.StaticFile("/styles.css", "./static/css/styles.css")
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
		err := checkUrl(body.Data)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		fetchMusic(c, body.Data)
	})

	router.Run(":80")
}

func fetchMusic(c *gin.Context, data string) {

	ydl := exec.Command("yt-dlp", data, "-o-")
	ffmpeg := exec.Command("ffmpeg", "-i", "/dev/stdin", "-vn", "-f", "mp3", "-")

	ydlOut, err := ydl.StdoutPipe()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ydl.Stderr = os.Stderr

	ffmpeg.Stdin = ydlOut
	ffmpeg.Stdout = c.Writer
	ffmpeg.Stderr = os.Stderr
	fmt.Println("Starting-----------------------")
	go func() {
		defer ydlOut.Close()
		if err := ydl.Run(); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
	}()

	if err := ffmpeg.Run(); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	fmt.Println("Done-----------------------")
}

func checkUrl(url string) error {

	// Define the regular expression pattern
	pattern := `^https://www\.youtube\.com/watch\?v=[a-zA-Z0-9_\-]+$`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Check if the string matches the pattern
	if re.MatchString(url) {
		return nil
	} else {
		return fmt.Errorf("invalid url")
	}
}
