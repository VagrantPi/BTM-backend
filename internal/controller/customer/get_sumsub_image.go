package customer

import (
	"BTM-backend/pkg/api"
	"BTM-backend/pkg/error_code"
	"BTM-backend/pkg/logger"
	"BTM-backend/third_party/sumsub"
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type GetSumsubImageReq struct {
	InspectionId string `form:"inspection_id" binding:"required"`
	ImageId      string `form:"image_id" binding:"required"`
}

func GetSumsubImage(c *gin.Context) {
	log := logger.Zap().WithClassFunction("api", "GetSumsubImage")
	defer func() {
		_ = log.Sync()
	}()
	c.Set("log", log)

	req := GetSumsubImageReq{}
	err := c.BindQuery(&req)
	if err != nil {
		log.Error("c.BindQuery(req)", zap.Any("err", err))
		api.ErrResponse(c, "c.BindQuery(&req)", errors.BadRequest(error_code.ErrInvalidRequest, "c.BindQuery(&req)").WithCause(err))
		return
	}

	imageReader, err := sumsub.GettingDocumentImages(req.InspectionId, req.ImageId)
	if err != nil {
		log.Error("sumsub.GettingDocumentImages", zap.Any("err", err))
		api.ErrResponse(c, "sumsub.GettingDocumentImages", errors.InternalServer(error_code.ErrBTMSumsubGetItem, "sumsub.GettingDocumentImages").WithCause(err))
		return
	}
	if imageReader != nil {
		defer imageReader.Close()
	}

	// 設定 Content-Type 為圖片
	c.Header("Content-Type", "image/jpeg") // 根據實際圖片類型調整

	// 直接將圖片串流回傳

	// 使用 DataFromReader 回傳
	c.DataFromReader(
		http.StatusOK,
		-1, // 未知長度
		"image/jpeg",
		imageReader,
		map[string]string{
			"Cache-Control": "no-cache",
		},
	)
}

// 支持自訂縮放比例的圖片壓縮函式
func resizeImage(originalReader io.Reader, scaleFactor float64) (io.Reader, error) {
	// 解碼原始圖片
	originalImage, _, err := image.Decode(originalReader)
	if err != nil {
		return nil, err
	}

	bounds := originalImage.Bounds()

	// 計算新的寬高
	newWidth := int(float64(bounds.Dx()) * scaleFactor)
	newHeight := int(float64(bounds.Dy()) * scaleFactor)

	// 建立新圖片
	resizedImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 使用協程數量
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	// 建立等待組
	var wg sync.WaitGroup

	// 分區塊縮放
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for y := start; y < end; y++ {
				for x := 0; x < newWidth; x++ {
					// 對應原圖的像素位置
					srcX := int(float64(x) / scaleFactor)
					srcY := int(float64(y) / scaleFactor)
					color := originalImage.At(srcX, srcY)
					resizedImage.Set(x, y, color)
				}
			}
		}(i*newHeight/numCPU, (i+1)*newHeight/numCPU)
	}

	wg.Wait()

	// 編碼為 JPEG
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resizedImage, &jpeg.Options{Quality: 80})
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
