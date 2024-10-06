package media_image

import (
	"github.com/smartmediafiles/media/media/maps"
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

// ImageFileTypesExtensions is a map of media.Image file types to their file extensions.
var ImageFileTypesExtensions = maps.MapFileTypeExtensions{
	ImageBmp:  {".bmp", ".dib"},
	ImageGif:  {".gif"},
	ImageHeic: {".heic"},
	ImageHeif: {".heif"},
	ImageJpeg: {".jpg", ".jpeg", ".jpe", ".jif", ".jfif", ".jfi"},
	ImagePng:  {".png"},
	ImageTiff: {".tiff", ".tif"},
	ImageWebp: {".webp"},
}
