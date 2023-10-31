package main

/*
log in
cd .. && ssh -i "itra-ec2-key-pair.pem" ec2-user@ec2-15-206-90-198.ap-south-1.compute.amazonaws.com

rsync -avz -e "ssh -i C:\Users\Rutik\Desktop\code\irta\itra-ec2-key-pair.pem" C:\Users\Rutik\Desktop\code\irta\itra_robo_aws_lambda ec2-user@3.110.204.86:
scp -i "C:\Users\Rutik\Desktop\code\irta\itra-ec2-key-pair.pem" -r C:\Users\Rutik\Desktop\code\irta\itra_robo_aws_lambda ec2-user@15.206.90.198:

sudo systemctl daemon-reload
sudo systemctl enable itra-robo-video-merger-app
sudo systemctl start itra-robo-video-merger-app
sudo systemctl status itra-robo-video-merger-app
sudo systemctl stop itra-robo-video-merger-app
sudo systemctl restart itra-robo-video-merger-app
sudo systemctl disable itra-robo-video-merger-app

*/

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}
	
	
	app := fiber.New(fiber.Config{
		BodyLimit:    100 * 1024 * 1024,
		IdleTimeout:  time.Minute * 20,
		ReadTimeout:  time.Minute * 20,
		WriteTimeout: time.Minute * 20,
	})

	app.Post("/merge-video", handleMergeVideo)




	log.Fatal(app.Listen("0.0.0.0" + port))
}
func generateUniqueFilename(fileType string) string {
	uuidString := uuid.New().String()
	uniqueFilename := fmt.Sprintf("./files/temp_%s_%s", uuidString, fileType)
	return uniqueFilename
}
func generateUniqueFilenameForOutput(fileType string) string {
	uuidString := uuid.New().String()
	uniqueFilename := fmt.Sprintf("./files/output/temp_%s_%s", uuidString, fileType)
	return uniqueFilename
}

func handleMergeVideo(c *fiber.Ctx) error {
	imageType := "image.jpg"
	videoType := "video.mp4"
	// Create temporary file paths
	tempImageFile := generateUniqueFilename(imageType)
	tempVideoFile := generateUniqueFilename(videoType)
	outputVideoFile := generateUniqueFilenameForOutput(videoType)

	for i := 0; i < 2; i++ {
		var file multipart.FileHeader
		if i == 0 {
			ifile, erro := c.FormFile("image")
			if erro != nil {
				return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
					"messege": "Image Required",
					"error":   erro.Error(),
				})
			}
			file = *ifile

		} else {
			vfile, erro := c.FormFile("video")
			if erro != nil {
				return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
					"messege": "Image Required",
					"error":   erro.Error(),
				})
			}
			file = *vfile

		}
		src, err := file.Open()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"messege": "Failed to open file",
				"error":   err.Error(),
			})

		}
		defer src.Close()

		dstPath := tempImageFile
		if i == 1 {
			dstPath = tempVideoFile
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"messege": "Failed to create destination file",
				"error":   err.Error(),
			})

		}
		defer os.Remove(dstPath)
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"messege": "Failed to write file content",
				"error":   err.Error(),
			})

		}
	}

	log.Print(outputVideoFile)
	// Place tempImageFile as overlay on tempVideoFile and resize outputVideoFile to 1500x1500 using FFmpeg
	ffmpegPath := "./ffmpeg" // for devlopment
	// ffmpegPath := "ffmpeg" // for production
	cmd := exec.Command(ffmpegPath, "-i", tempVideoFile, "-i", tempImageFile, "-filter_complex", "[0:v]scale=1500:1500[base];[1:v]scale=1500:1500[overlay];[base][overlay]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/2", outputVideoFile)

	err := cmd.Run()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"messege": "Failed to process video",
			"error":   err.Error(),
		})
	}

	// Open the processed video file
	file, err := os.Open(outputVideoFile)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"messege": "Failed to open processed video file",
			"error":   err.Error(),
		})

	}

	// Get the file size
	fileStat, err := file.Stat()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"messege": "Failed to get file information",
			"error":   err.Error(),
		})
	}
	fileSize := fileStat.Size()

	// Send the processed video data back as the response
	c.Set("Content-Disposition", "attachment; filename=output.mp4")
	c.Set("Content-Length", strconv.FormatInt(fileSize, 10))
	c.Set("Content-Type", "video/mp4")

	// Stream the file to the client
	_, _ = io.Copy(c, file)
	file.Close()

	return nil
}
