package media_image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BmpImages(t *testing.T) {
	t.Log("Testing BMP images")
}

func Test_GifImages(t *testing.T) {
	t.Log("Testing GIF images")

	t.Run("sunflower-plants", func(t *testing.T) {
		imgInfo, err := NewImageInfo("samples/gif/sunflower-plants.gif")
		if err != nil {
			t.Fatal(err)
		}
		i, err := imgInfo.Exif()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, i.FileInfo.Name(), "sunflower-plants.gif")
		assert.Equal(t, i.FileType, ImageGif)
		assert.Equal(t, i.FileExt, ExtensionGif)
		assert.True(t, i.IsImage())
		assert.True(t, i.ImageData.ImageWidth > 0)
		assert.True(t, i.ImageData.ImageHeight > 0)
	})
}

func Test_HeicImages(t *testing.T) {
	t.Log("Testing HEIC images")

	t.Run("netherlands", func(t *testing.T) {
		imgInfo, err := NewImageInfo("samples/heic/netherlands.heic")
		if err != nil {
			t.Fatal(err)
		}
		i, err := imgInfo.Exif()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, i.FileInfo.Name(), "netherlands.heic")
		assert.Equal(t, i.FileType, ImageHeic)
		assert.Equal(t, i.FileExt, ExtensionHeic)
		assert.True(t, i.IsPhoto())
		assert.Equal(t, i.ImageData.GPSTimeZone, "Europe/Amsterdam")
	})

	t.Run("holyroodhouse", func(t *testing.T) {
		imgInfo, err := NewImageInfo("samples/heic/holyroodhouse.heic")
		if err != nil {
			t.Fatal(err)
		}
		i, err := imgInfo.Exif()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, i.FileInfo.Name(), "holyroodhouse.heic")
		assert.Equal(t, i.FileType, ImageHeic)
		assert.Equal(t, i.FileExt, ExtensionHeic)
		assert.True(t, i.IsPhoto())
		assert.Equal(t, i.ImageData.GPSTimeZone, "Europe/London")
	})
}

func Test_HeifImages(t *testing.T) {
	t.Log("Testing HEIF images")

	t.Run("madrid", func(t *testing.T) {
		imgInfo, err := NewImageInfo("samples/heif/madrid.heif")
		if err != nil {
			t.Fatal(err)
		}
		i, err := imgInfo.Exif()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, i.FileInfo.Name(), "madrid.heif")
		assert.Equal(t, i.FileType, ImageHeif)
		assert.Equal(t, i.FileExt, ExtensionHeif)
		assert.True(t, i.IsPhoto())
		assert.Equal(t, i.ImageData.GPSTimeZone, "Europe/Madrid")
	})
}

func Test_JpegImages(t *testing.T) {
	t.Log("Testing JPEG images")

	t.Run("exif-org-1", func(t *testing.T) {
		imgInfo, err := NewImageInfo("samples/jpg/exif-org/exif-org-1.jpg")
		if err != nil {
			t.Fatal(err)
		}
		i, err := imgInfo.Exif()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, i.FileInfo.Name(), "exif-org-1.jpg")
		assert.Equal(t, i.FileType, ImageJpeg)
		assert.Equal(t, i.FileExt, ExtensionJpg)
		assert.True(t, i.IsPhoto())
		assert.Equal(t, i.ImageData.GPSTimeZone, "")
	})

	t.Run("gps-1", func(t *testing.T) {
		imgInfo, err := NewImageInfo("samples/jpg/gps/gps-1.jpg")
		if err != nil {
			t.Fatal(err)
		}
		i, err := imgInfo.Exif()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, i.FileInfo.Name(), "gps-1.jpg")
		assert.Equal(t, i.FileType, ImageJpeg)
		assert.Equal(t, i.FileExt, ExtensionJpg)
		assert.Equal(t, i.ImageData.GPSTimeZone, "Europe/Rome")
	})
}

func Test_PngImages(t *testing.T) {
	t.Log("Testing PNG images")
}

func Test_TiffImages(t *testing.T) {
	t.Log("Testing TIFF images")
}

func Test_WebpImages(t *testing.T) {
	t.Log("Testing WEBP images")

	t.Run("giphy", func(t *testing.T) {
		imgInfo, err := NewImageInfo("samples/webp/giphy.webp")
		if err != nil {
			t.Fatal(err)
		}
		i, err := imgInfo.Exif()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, i.FileInfo.Name(), "giphy.webp")
		assert.Equal(t, i.FileType, ImageWebp)
		assert.Equal(t, i.FileExt, ExtensionWebp)
		assert.True(t, i.IsImage())
		assert.True(t, i.ImageData.ImageWidth > 0)
		assert.True(t, i.ImageData.ImageHeight > 0)
	})

	t.Run("Nærøyfjorden", func(t *testing.T) {
		imgInfo, err := NewImageInfo("samples/webp/Nærøyfjorden.webp")
		if err != nil {
			t.Fatal(err)
		}
		i, err := imgInfo.Exif()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, i.FileInfo.Name(), "Nærøyfjorden.webp")
		assert.Equal(t, i.FileType, ImageWebp)
		assert.Equal(t, i.FileExt, ExtensionWebp)
		assert.True(t, i.IsImage())
		assert.True(t, i.ImageData.ImageWidth > 0)
		assert.True(t, i.ImageData.ImageHeight > 0)
	})
}
