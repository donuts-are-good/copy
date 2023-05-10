package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var recursive bool

type Progress struct {
	BytesCopied int64
	TotalBytes  int64
	FilesCopied int
	TotalFiles  int
}

func main() {
	flag.BoolVar(&recursive, "r", false, "copy directories recursively")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Println("Usage: ./copy -r /path/to/source /path/to/destination")
		os.Exit(1)
	}

	source := flag.Arg(0)
	destination := flag.Arg(1)

	fmt.Printf("copying: %s -> %s\n", source, destination)

	progressChan := make(chan Progress)
	go displayProgress(progressChan)

	startTime := time.Now()

	var totalFiles int
	var totalBytes int64
	filepath.Walk(source, func(srcPath string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			totalFiles++
			totalBytes += info.Size()
		}
		return nil
	})

	var progress Progress
	progress.TotalFiles = totalFiles
	progress.TotalBytes = totalBytes

	var wg sync.WaitGroup

	if recursive {
		err := filepath.Walk(source, func(srcPath string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(source, srcPath)
			if err != nil {
				return err
			}

			dstPath := filepath.Join(destination, relPath)

			if info.IsDir() {
				return os.Mkdir(dstPath, info.Mode())
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				err := copyFile(srcPath, dstPath, info.Mode(), progressChan, &progress)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					os.Exit(1)
				}
			}()

			return nil
		})
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	} else {
		info, err := os.Stat(source)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err = copyFile(source, destination, info.Mode(), progressChan, &progress)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
		}()
	}

	wg.Wait()
	close(progressChan)

	duration := time.Since(startTime)
	fmt.Printf("\nCopying took %s... done!\n", duration)
}

func copyFile(src, dst string, mode fs.FileMode, progressChan chan<- Progress, progress *Progress) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := dstFile.Write(buf[:n]); err != nil {
			return err
		}

		progress.BytesCopied += int64(n)
		progressChan <- *progress

		if err == io.EOF {
			break
		}
	}

	progress.FilesCopied++
	progressChan <- *progress

	return nil
}

func displayProgress(progressChan <-chan Progress) {
	for progress := range progressChan {
		fmt.Printf("\r%d files - %.1fGB [%s]",
			progress.FilesCopied,
			float64(progress.BytesCopied)/(1<<30),
			progressBar(progress.BytesCopied, progress.TotalBytes, 20),
		)
	}
}

func progressBar(current, total int64, width int) string {
	ratio := float64(current) / float64(total)
	completeChars := int(ratio * float64(width))
	progressBar := strings.Repeat("=", completeChars)
	if completeChars < width {
		progressBar += ">"
		remainingChars := width - completeChars - 1
		progressBar += strings.Repeat(" ", remainingChars)
	}
	return progressBar
}
