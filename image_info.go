package media_image

import (
	"errors"
	"image"
	"os"
	"path/filepath"

	"github.com/dsoprea/go-exif/v3"
	"github.com/smartmediafiles/media.fs/fs"
	"github.com/smartmediafiles/media/media/types"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// ImageInfo is a structure that contains information about an image file.
// This information are extracted from the image file exif data.
type ImageInfo struct {
	// File information
	FileInfo fs.FileInfo
	FileType types.FileType
	FileExt  types.FileExtension

	// Image information
	ImageData ImageData
}

// NewImageInfo creates a new ImageInfo struct.
func NewImageInfo(path string) (*ImageInfo, error) {
	// Retrieve file information
	fileInfo, err := fs.NewFileInfo(path)
	if err != nil {
		return nil, err
	}

	// Retrieve file type and extension
	fileType, fileExt := ImageFileTypesExtensions.GetFileTypeAndExtension(fileInfo.Name())

	// Assign values
	i := new(ImageInfo)
	i.FileInfo = fileInfo
	i.FileType = fileType
	i.FileExt = fileExt

	// extract minimal information from the image file
	_ = i.extractData()

	return i, nil
}

// Exif extracts the image information from the exif data.
func (i *ImageInfo) Exif() (*ImageInfo, error) {

	// Parse the file to extract exif data
	exifParser := NewExifParser()
	rawExif, err := exifParser.Parse(filepath.Join(i.FileInfo.Abs(), i.FileInfo.Name()), i.FileType)
	if errors.Is(err, exif.ErrNoExif) {
		return i, nil
	}
	if err != nil {
		return i, err
	}

	// Parse the exif data to extract image data
	exifDataParser := NewExifDataParser()
	imageData, err := exifDataParser.Parse(rawExif)
	if err != nil {
		return i, err
	}

	// Assign values
	i.ImageData = imageData

	return i, nil
}

// IsPhoto checks if the image is a photo.
func (i *ImageInfo) IsPhoto() bool {
	return IsPhoto(i.FileType)
}

// IsImage checks if the image is an image.
func (i *ImageInfo) IsImage() bool {
	return IsImage(i.FileType)
}

// extractData extracts minimal information from the image file.
func (i *ImageInfo) extractData() error {
	file, err := os.Open(filepath.Join(i.FileInfo.Abs(), i.FileInfo.Name()))
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image to get its dimensions
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return err
	}

	// Assign dimensions to ImageData
	i.ImageData.ImageWidth = img.Width
	i.ImageData.ImageHeight = img.Height

	// Use file date as image date
	if i.FileInfo.CreationTime().IsZero() {
		i.ImageData.DateTime = i.FileInfo.LastWriteTime()
	} else {
		i.ImageData.DateTime = i.FileInfo.CreationTime()
	}

	return nil
}
