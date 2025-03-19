package updater

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/System233/ll-killer-go/utils"
)

// DownloadFile 以 "已下载 XMB/YMB (Z%)" 格式显示进度
func DownloadFile(ctx context.Context, url, filepath string, hash string) error {
	utils.Debug("DownloadFile", url, filepath)
	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败:%v", err)
	}

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("执行请求失败:%v", err)
	}
	defer resp.Body.Close()

	// 获取文件大小
	size := resp.ContentLength
	if size <= 0 {
		fmt.Println("无法获取文件大小，进度显示可能不准确")
	}

	// 创建目标文件
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("无法打开文件:%v", err)
	}
	defer out.Close()

	// 开始下载
	hasher := sha256.New()
	startTime := time.Now()
	progressReader := &ProgressReader{
		Reader: io.TeeReader(resp.Body, hasher),
		Total:  size,
		Start:  startTime,
	}

	_, err = io.Copy(out, progressReader)
	if err != nil {
		return fmt.Errorf("复制失败:%v", err)
	}
	log.Println("\n下载完成:", filepath)

	calculatedSHA256 := hex.EncodeToString(hasher.Sum(nil))
	fmt.Printf("计算的 SHA-256: %s\n", calculatedSHA256)

	if hash != "" {
		if calculatedSHA256 == hash {
			fmt.Println("SHA-256 校验通过")
		} else {
			fmt.Println("SHA-256 校验失败！")
			return fmt.Errorf("文件哈希不匹配，可能已损坏")
		}
	}
	return nil
}

// ProgressReader 计算下载进度
type ProgressReader struct {
	Reader io.Reader
	Total  int64
	Loaded int64
	Start  time.Time
}

// Read 实现 io.Reader 接口
func (p *ProgressReader) Read(buf []byte) (int, error) {
	n, err := p.Reader.Read(buf)
	if n > 0 {
		p.Loaded += int64(n)

		// 计算进度
		percentage := float64(p.Loaded) / float64(p.Total) * 100
		elapsed := time.Since(p.Start).Seconds()
		speed := float64(p.Loaded) / elapsed / (1024 * 1024) // MB/s

		// 计算单位
		fmt.Printf("\r已下载 %.1fMB/%.1fMB (%.1f%%)  ⚡ %.2f MB/s",
			float64(p.Loaded)/(1024*1024),
			float64(p.Total)/(1024*1024),
			percentage,
			speed,
		)
	}
	return n, err
}
