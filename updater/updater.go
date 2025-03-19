package updater

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/System233/ll-killer-go/config"
	"github.com/System233/ll-killer-go/utils"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/sys/unix"
)

//go:embed publickey.asc
var publicKeyData string

const (
	latestURL       = "/releases/latest"
	githubAPIURL    = "https://api.github.com/repos/" + config.Repo
	killerMirrorURL = "https://ll-killer.win"
	SHA256SUMS      = "SHA256SUMS"
)

type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int    `json:"size"`
	} `json:"assets"`
}
type UpdateAsset struct {
	URL       string
	OriginURL string
	Name      string
	SHA256    string
	Size      int
}
type UpdateInfo struct {
	Tag    string
	Body   string
	Note   string
	Assets []UpdateAsset
}
type FetchError struct {
	StatusCode int
	Status     string
}

func (e *FetchError) Error() string {
	return e.Status
}

func fetch(ctx context.Context, url string) ([]byte, error) {
	utils.Debug("fetch", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, &FetchError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}
	return io.ReadAll(resp.Body)
}
func verifyData(dataReader io.Reader, signatureReader io.Reader) error {
	keyring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(publicKeyData))
	if err != nil {
		return fmt.Errorf("读取公钥失败: %v", err)
	}

	_, err = openpgp.CheckArmoredDetachedSignature(keyring, dataReader, signatureReader)
	if err != nil {
		return fmt.Errorf("签名验证失败: %v", err)
	}

	return nil
}
func fetchLatestSHA256(ctx context.Context, base string, tag string) (map[string]string, error) {
	hashURL := fmt.Sprintf("%s/releases/download/%s/%s", base, tag, SHA256SUMS)
	hashAscURL := fmt.Sprint(hashURL, ".asc")
	hashData, err := fetch(ctx, hashURL)
	if err != nil {
		return nil, err
	}
	ascData, err := fetch(ctx, hashAscURL)
	if err != nil {
		return nil, err
	}
	if err := verifyData(strings.NewReader(string(hashData)), strings.NewReader(string(ascData))); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(hashData))
	var result map[string]string = make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		chunks := strings.SplitN(line, " ", 2)
		if len(chunks) != 2 {
			continue
		}
		key := strings.TrimSpace(chunks[0])
		value := strings.TrimSpace(chunks[1])
		result[value] = key
	}
	return result, nil
}

func fetchRelease(ctx context.Context, base string, note string, ch chan<- *UpdateInfo, errCh chan<- error) {
	targetURL := fmt.Sprint(base, "/releases/latest")
	utils.Debug("fetchRelease", targetURL)
	resp, err := fetch(ctx, targetURL)
	if err != nil {
		errCh <- err
		return
	}

	var release ReleaseInfo
	if err := json.Unmarshal(resp, &release); err != nil {
		errCh <- err
		return
	}

	hashmap, err := fetchLatestSHA256(ctx, base, release.TagName)
	if err != nil {
		errCh <- err
		return
	}
	var assets []UpdateAsset
	for _, asset := range release.Assets {
		value, ok := hashmap[asset.Name]
		if !ok {
			continue
		}
		assets = append(assets, UpdateAsset{
			URL:       fmt.Sprintf("%s/releases/download/%s/%s", base, release.TagName, asset.Name),
			SHA256:    value,
			Name:      asset.Name,
			Size:      asset.Size,
			OriginURL: asset.BrowserDownloadURL,
		})
	}
	ch <- &UpdateInfo{
		Tag:    release.TagName,
		Body:   release.Body,
		Note:   note,
		Assets: assets,
	}
}

func CheckForUpdate(duration time.Duration) (*UpdateInfo, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan *UpdateInfo, 2)
	errCh := make(chan error, 2)

	go fetchRelease(ctx, githubAPIURL, "Github地址", ch, errCh)
	go fetchRelease(ctx, killerMirrorURL, "镜像地址", ch, errCh)

	timeout := time.After(duration)
	for i := 0; i < 2; i++ {
		select {
		case rel := <-ch:
			return rel, nil
		case err := <-errCh:
			if i == 1 {
				return nil, err
			}
		case <-timeout:
			return nil, errors.New("请求超时")
		}
	}
	return nil, errors.New("未能获取版本信息")
}
func FindUpdateAsset(info UpdateInfo) *UpdateAsset {
	name := fmt.Sprintf("%s-%s", config.KillerExec, config.Variant)
	for _, item := range info.Assets {
		if item.Name == name {
			return &item
		}
	}
	return nil
}

type UpdateOption struct {
	Yes     bool
	Retry   int
	Timeout time.Duration
}

func Update(opt UpdateOption) error {
	log.Println("正在检查更新...")
	var info *UpdateInfo
	var err error
	for i := range opt.Retry {
		info, err = CheckForUpdate(opt.Timeout)
		if err == nil {
			break
		}
		if e, ok := err.(*FetchError); ok {
			if e.StatusCode == 404 {
				return fmt.Errorf("更新失败, 请浏览器打开 %s 手动更新: %v", config.GithubURL, e)
			}
		}
		log.Printf("第%d次检查更新失败:%v, 正在重试...", i+1, err)
	}
	if err != nil {
		return fmt.Errorf("已达最大重试次数，检查更新失败:%v", err)
	}
	if info.Tag == config.Tag {
		log.Println("当前已是最新版本: ", info.Tag)
		return nil
	}
	asset := FindUpdateAsset(*info)
	if asset == nil {
		return fmt.Errorf("不支持此架构类型的更新，你可能需要手动编译本项目:%s", config.Variant)
	}
	target, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取可执行文件位置:%s", config.Variant)
	}
	log.Println("检测到新版本：", info.Tag)
	log.Println("当前版本:", config.Tag)
	log.Println("更新位置:", target)
	log.Println("文件名:  ", asset.Name)
	log.Println("原始地址:", asset.OriginURL)
	log.Printf("下载地址: %s (%s)", asset.URL, info.Note)
	log.Println("SHA-256: ", asset.SHA256)
	log.Println("大小: ", fmt.Sprintf("%d (%.2fMiB)", asset.Size, float64(asset.Size)/(1024*1024)))
	if !opt.Yes {
		var input string
		fmt.Printf("是否更新？(Y/N): ")
		if _, err := fmt.Scanf("%s", &input); err != nil {
			log.Println("输入错误:", err)
		}
		input = strings.ToUpper(strings.TrimSpace(input)) // 处理大小写和空格
		if input == "Y" {
		} else {
			log.Println("更新已取消")
			return nil
		}
	}
	newTarget := fmt.Sprint(target, ".new")
	ctx := context.Background()
	if err := DownloadFile(ctx, asset.URL, newTarget, asset.SHA256); err != nil {
		return err
	}
	if err := unix.Chmod(newTarget, 0755); err != nil {
		return fmt.Errorf("无法设置执行权限:%v", err)
	}
	if err := unix.Rename(newTarget, target); err != nil {
		return fmt.Errorf("重命名原文件失败:%v", err)
	}
	log.Println("更新成功")
	return nil
}
