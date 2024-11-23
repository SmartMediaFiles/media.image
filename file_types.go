package media_image

import (
	"github.com/smartmediafiles/media/media/types"
)

// List of supported media.Image file types.
const (
	ImageBmp  types.FileType = "bmp"  // Bitmap Image
	ImageGif  types.FileType = "gif"  // Graphics Interchange Format (GIF)
	ImageHeic types.FileType = "heic" // High Efficiency Image Container (HEIC)
	ImageHeif types.FileType = "heif" // High Efficiency Image File Format (HEIF)
	ImageJpeg types.FileType = "jpg"  // Joint Photographic Experts Group (JPEG)
	ImagePng  types.FileType = "png"  // Portable Network Graphics (PNG)
	ImageTiff types.FileType = "tiff" // Tagged Image File Format (TIFF)
	ImageWebp types.FileType = "webp" // Google WebP Image
)

// ImageFileTypes is a list of supported media.Image file types.
var ImageFileTypes = []types.FileType{
	ImageBmp,
	ImageGif,
	ImageHeic,
	ImageHeif,
	ImageJpeg,
	ImagePng,
	ImageTiff,
	ImageWebp,
}

// IsPhoto checks if the given file type is considered a photo.
func IsPhoto(fileType types.FileType) bool {
	switch fileType {
	case ImageJpeg, ImageHeic, ImageHeif:
		return true
	default:
		return false
	}
}

// IsImage checks if the given file type is considered an image.
func IsImage(fileType types.FileType) bool {
	switch fileType {
	case ImageBmp, ImageGif, ImagePng, ImageTiff, ImageWebp:
		return true
	default:
		return false
	}
}
