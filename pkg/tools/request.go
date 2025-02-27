package tools

import (
	"BTM-backend/pkg/error_code"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
)

// JsonRequest 向外部發送請求，限時一分鐘。
// 記住錯誤有分，此函式回傳的錯誤只是內部 500 錯誤，第三方服務回傳錯誤要自行判斷
// response 請再給個指標(pointer)
// 這隻不管業務邏輯、http status，專注在請求上
func JsonRequest(method, url string, header map[string]string, body any) (int, []byte, map[string]string,
	error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		err = errors.InternalServer(error_code.ErrToolsHttpMarshal, "json.Marshal").WithCause(err).
			WithMetadata(map[string]string{"body": fmt.Sprintf("%v", body)})
		return 0, nil, map[string]string{}, err
	}
	if body == nil {
		jsonBody = nil
	}

	if header == nil {
		header = make(map[string]string)
	}

	header["Content-Type"] = "application/json"
	return Request(method, url, header, jsonBody)
}

func Request(method, url string, header map[string]string, body []byte) (status int,
	respBody []byte, headers map[string]string, err error) {
	// 製作請求內容
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		err = errors.InternalServer(error_code.ErrToolsHttpRequest, "http.NewRequest").WithCause(err).
			WithMetadata(map[string]string{"method": method, "url": url})
		return 0, nil, map[string]string{}, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	// 製作請求載體
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// 發送請求
	resp, err := client.Do(req)
	if err != nil {
		err = errors.InternalServer(error_code.ErrToolsHttpRequestDo, "client.Do").WithCause(err).
			WithMetadata(map[string]string{"method": method, "url": url})
		return 0, nil, map[string]string{}, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	headers = map[string]string{}
	for key, values := range resp.Header {
		for _, value := range values {
			headers[key] = value
		}
	}

	// 讀取請求內容
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.InternalServer(error_code.ErrToolsHttpRequestIoRead, "io.ReadAll").WithCause(err)
		return 0, nil, map[string]string{}, err
	}
	return resp.StatusCode, result, headers, nil
}

func DownloadFile(url string) (body io.ReadCloser, err error) {
	// #nosec G107 - 注意 url 來源，不要被 300 攻擊了
	resp, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("http.Get() failed: %s", err)
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Get() is not 200: %s", resp.Status)
		return
	}

	return resp.Body, nil
}

func DownloadFileWithHeader(url string, header map[string]string) (body io.ReadCloser, err error) {
	// #nosec G107 - 注意 url 來源，不要被 300 攻擊了
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = errors.InternalServer(error_code.ErrToolsHttpRequest, "http.NewRequest").WithCause(err).
			WithMetadata(map[string]string{"method": "GET", "url": url})
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	// 製作請求載體
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// 發送請求
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("http.Get() failed: %s", err)
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Get() is not 200: %s", resp.Status)
		return
	}

	return resp.Body, nil
}
