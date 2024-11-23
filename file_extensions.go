package media_image

import "github.com/smartmediafiles/media/media/types"

// List of supported media.Image file extensions.
const (
	ExtensionBmp  types.FileExtension = ".bmp"  // Bitmap Image
	ExtensionDib  types.FileExtension = ".dib"  // Device Independent Bitmap
	ExtensionGif  types.FileExtension = ".gif"  // Graphics Interchange Format (GIF)
	ExtensionHeic types.FileExtension = ".heic" // High Efficiency Image Container (HEIC)
	ExtensionHeif types.FileExtension = ".heif" // High Efficiency Image File Format (HEIF)
	ExtensionJpg  types.FileExtension = ".jpg"  // Joint Photographic Experts Group (JPEG)
	ExtensionJpeg types.FileExtension = ".jpeg" // Joint Photographic Experts Group (JPEG)
	ExtensionJpe  types.FileExtension = ".jpe"  // Joint Photographic Experts Group (JPEG)
	ExtensionJif  types.FileExtension = ".jif"  // Joint Photographic Experts Group (JPEG)
	ExtensionJfif types.FileExtension = ".jfif" // Joint Photographic Experts Group (JPEG)
	ExtensionJfi  types.FileExtension = ".jfi"  // Joint Photographic Experts Group (JPEG)
	ExtensionPng  types.FileExtension = ".png"  // Portable Network Graphics (PNG)
	ExtensionTiff types.FileExtension = ".tiff" // Tagged Image File Format (TIFF)
	ExtensionTif  types.FileExtension = ".tif"  // Tagged Image File Format (TIFF)
	ExtensionWebp types.FileExtension = ".webp" // Google WebP Image
)

// ImageFileExtensions is a list of supported media.Image file extensions.
var ImageFileExtensions = []types.FileExtension{
	ExtensionBmp,
	ExtensionDib,
	ExtensionGif,
	ExtensionHeic,
	ExtensionHeif,
	ExtensionJpg,
	ExtensionJpeg,
	ExtensionJpe,
	ExtensionJif,
	ExtensionJfif,
	ExtensionJfi,
	ExtensionPng,
	ExtensionTiff,
	ExtensionTif,
	ExtensionWebp,
}
