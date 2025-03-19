package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// Configuration
const (
	streamURL    = "https://example.com/upload" // Replace with your website URL
	captureWidth = 640
	captureHeight = 480
	frameRate    = 30
	uploadInterval = 2 * time.Second
)

func main() {
	log.Println("Starting Raspberry Pi camera streaming service")
	
	// Check if raspivid is available
	if _, err := exec.LookPath("raspivid"); err != nil {
		log.Fatal("raspivid not found. Please ensure the Raspberry Pi camera is enabled and raspivid is installed")
	}
	
	// Create a temporary file to store video segments
	tempFile, err := os.CreateTemp("", "video-*.h264")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFileName := tempFile.Name()
	defer os.Remove(tempFileName) // Clean up temp file when done
	
	log.Printf("Capturing video and uploading to %s", streamURL)
	
	// Start a ticker for regular uploads
	ticker := time.NewTicker(uploadInterval)
	defer ticker.Stop()
	
	// Start the camera capture in a separate goroutine
	go captureVideo(tempFileName)
	
	// Begin upload loop
	for range ticker.C {
		// We'll read from the file each interval and upload it
		if err := uploadVideoSegment(tempFileName); err != nil {
			log.Printf("Error uploading video segment: %v", err)
		}
	}
}

func captureVideo(outputFile string) {
	for {
		cmd := exec.Command("raspivid",
			"-o", outputFile,
			"-w", fmt.Sprint(captureWidth),
			"-h", fmt.Sprint(captureHeight),
			"-fps", fmt.Sprint(frameRate),
			"-t", "0", // Run continuously
			"-n",      // No preview window
			"-ih",     // Insert inline headers for better streaming
		)
		
		if err := cmd.Start(); err != nil {
			log.Printf("Failed to start raspivid: %v", err)
			time.Sleep(5 * time.Second) // Wait before retry
			continue
		}
		
		if err := cmd.Wait(); err != nil {
			log.Printf("raspivid process ended with error: %v", err)
			time.Sleep(5 * time.Second) // Wait before retry
		}
	}
}

func uploadVideoSegment(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open video file: %v", err)
	}
	defer file.Close()
	
	// Create buffer for the multipart form
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	
	// Create form file
	formFile, err := multipartWriter.CreateFormFile("video", "segment.h264")
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}
	
	// Copy file content to form file
	if _, err = io.Copy(formFile, file); err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}
	
	// Add timestamp
	if err = multipartWriter.WriteField("timestamp", time.Now().Format(time.RFC3339)); err != nil {
		return fmt.Errorf("failed to add timestamp: %v", err)
	}
	
	// Close multipart writer
	if err = multipartWriter.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %v", err)
	}
	
	// Create request
	req, err := http.NewRequest("POST", streamURL, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	
	// Set content type
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	
	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code: %d", resp.StatusCode)
	}
	
	log.Println("Video segment uploaded successfully")
	return nil
}
module facetracker

go 1.24.1

require (
	github.com/dhowden/raspicam v0.0.0-20190323051945-60ef25a6629f // indirect
	github.com/vladimirvivien/go4vl v0.0.5 // indirect
	gocv.io/x/gocv v0.40.0 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
)
github.com/dhowden/raspicam v0.0.0-20190323051945-60ef25a6629f h1:oT+wXjVnCofGqUTbPqRJKs3ne9vhxNumXk3PxMRkRTc=
github.com/dhowden/raspicam v0.0.0-20190323051945-60ef25a6629f/go.mod h1:OMzt06oC5E5YW+9MyiiuX5oNfBGzpTeMRmGOK5OB8SE=
github.com/vladimirvivien/go4vl v0.0.5 h1:jHuo/CZOAzYGzrSMOc7anOMNDr03uWH5c1B5kQ+Chnc=
github.com/vladimirvivien/go4vl v0.0.5/go.mod h1:FP+/fG/X1DUdbZl9uN+l33vId1QneVn+W80JMc17OL8=
gocv.io/x/gocv v0.40.0 h1:kGBu/UVj+dO6A9dhQmGOnCICSL7ke7b5YtX3R3azdXI=
gocv.io/x/gocv v0.40.0/go.mod h1:zYdWMj29WAEznM3Y8NsU3A0TRq/wR/cy75jeUypThqU=
golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c h1:F1jZWGFhYfh0Ci55sIpILtKKK8p3i2/krTr0H1rg74I=
golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
package main

import(
  "fmt"
  "log"
  "context"
  "net/http"
  "mime/multipart"
  "net/textproto"

  "github.com/vladimirvivien/go4vl/device"
  "github.com/vladimirvivien/go4vl/v4l2"
)

var (
  frames <-chan []byte
)

func imageServ(w http.ResponseWriter, req *http.Request) {
  mimeWriter := multipart.NewWriter(w)
  w.Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", mimeWriter.Boundary()))
  partHeader := make(textproto.MIMEHeader)
  partHeader.Add("Content-Type", "image/jpeg")

  var frame []byte
  for frame = range frames {
    partWriter, err := mimeWriter.CreatePart(partHeader)
    if err != nil {
      log.Printf("failed to create multi-part writer: %s", err)
      return
    }

    if _, err := partWriter.Write(frame); err != nil {
      log.Printf("failed to write image: %s", err)
    }
  }
}

func main() {
  port := ":9090"
  devName := "/dev/video0"

  camera, err := device.Open(
    devName,
	device.WithPixFormat(v4l2.PixFormat{
	    PixelFormat: v4l2.PixelFmtMJPEG, // Try MJPEG instead of YUYV
	    Width:       640,
	    Height:      480,
	    Field:       v4l2.FieldNone,
	    BytesPerLine: 0,
	    SizeImage:    0,
	}),
	device.WithFPS(15),
//    device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtYUYV, Width: 640, Height: 480}),
//    device.WithBufferSize(1),
 )
  if err != nil {
    log.Fatalf("failed to open device: %s", err)
  }
  defer camera.Close()

	return

  if err := camera.Start(context.TODO()); err != nil {
    log.Fatalf("camera start: %s", err)
  }

  frames = camera.GetOutput()
  http.HandleFunc("/stream", imageServ)
  log.Fatal(http.ListenAndServe(port, nil))
}
package main

import (
	"log"
	"fmt"
	"net"
	"os"

	"github.com/dhowden/raspicam"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6666")
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v", err)
		return
	}
	log.Println("Listening on 0.0.0.0:6666")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "accept: %v", err)
			return
		}
		log.Printf("Accepted connection from: %v\n", conn.RemoteAddr())
		go func() {
			s := raspicam.NewStill()
			errCh := make(chan error)
			go func() {
				for x := range errCh {
					fmt.Fprintf(os.Stderr, "%v\n", x)
			}
			}()
			log.Println("Capturing image...")
			raspicam.Capture(s, conn, errCh)
			log.Println("Done")
			conn.Close()
		}()
	}
}
package main 

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"gocv.io/x/gocv"
)

func imageServ(w http.ResponseWriter, req *http.Request, webcam *gocv.VideoCapture) {
	mimeWriter := multipart.NewWriter(w)
	w.Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", mimeWriter.Boundary()))
	partHeader := make(textproto.MIMEHeader)
	partHeader.Add("Content-Type", "image/jpeg")

	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			log.Printf("Cannot read device")
			break
		}
		if img.Empty() {
			continue
		}

		buf, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			log.Printf("Failed to encode jpeg: %v", err)
			continue
		}

		partWriter, err := mimeWriter.CreatePart(partHeader)
		if err != nil {
			log.Printf("failed to create multi-part writer: %s", err)
			return
		}

		partWriter.Write(buf.GetBytes())
		buf.Close()
		
		// Rate limiting to avoid overwhelming the network or CPU
		time.Sleep(33 * time.Millisecond) // ~30fps
	}
}

func main() {
	port := ":9090"
	devName := "/dev/video0"
	
	webcam, err := gocv.OpenVideoCapture(devName)
	if err != nil {
		log.Fatalf("Error opening video capture device: %v", err)
	}
	defer webcam.Close()

	return

	// Set webcam properties
	webcam.Set(gocv.VideoCaptureFrameWidth, 640)
	webcam.Set(gocv.VideoCaptureFrameHeight, 480)

	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		imageServ(w, r, webcam)
	})
	
	fmt.Printf("Starting video server on http://localhost%s/stream\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

