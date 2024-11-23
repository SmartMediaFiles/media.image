package media_image

import (
	"fmt"

	"github.com/dsoprea/go-exif/v3"
	heicexif "github.com/dsoprea/go-heic-exif-extractor/v2"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
	pngstructure "github.com/dsoprea/go-png-image-structure/v2"
	tiffstructure "github.com/dsoprea/go-tiff-image-structure/v2"
	"github.com/smartmediafiles/media/media/types"
)

// ExifParser is a struct that contains the EXIF parser.
type ExifParser struct{}

// NewExifParser creates a new ExifParser struct.
func NewExifParser() *ExifParser {
	return new(ExifParser)
}

// Parse parses the EXIF data from the file.
func (p *ExifParser) Parse(path string, fileType types.FileType) ([]byte, error) {

	// Switch on the file type
	switch fileType {
	case ImageBmp:
		return p.parseRaw(path)

	case ImageGif:
		return p.parseRaw(path)

	case ImageHeic, ImageHeif:
		return p.parseHeic(path)

	case ImageJpeg:
		return p.parseJpeg(path)

	case ImagePng:
		return p.parsePng(path)

	case ImageTiff:
		return p.parseTiff(path)

	case ImageWebp:
		return p.parseRaw(path)
	}

	return nil, fmt.Errorf("unsupported file type: %s", fileType)
}

// parseRaw parses the EXIF data from the file using exif.SearchFileAndExtractExif.
func (p *ExifParser) parseRaw(path string) ([]byte, error) {
	// Search the file for the EXIF data
	rawExif, err := exif.SearchFileAndExtractExif(path)
	if err != nil {
		return nil, err
	}
	return rawExif, nil
}

// parseHeic parses the EXIF data from the HEIC, HEIF file.
func (p *ExifParser) parseHeic(path string) ([]byte, error) {
	// Create a new HEIC media parser
	heicMediaParser := heicexif.NewHeicExifMediaParser()

	// Parse the HEIC file
	mediaContext, err := heicMediaParser.ParseFile(path)
	if err != nil {
		return nil, err
	}

	// Get the EXIF data
	_, rawExif, err := mediaContext.Exif()
	if err != nil {
		return nil, err
	}
	return rawExif, nil
}

// parseJpeg parses the EXIF data from the JPEG file.
func (p *ExifParser) parseJpeg(path string) ([]byte, error) {
	// Create a new JPEG media parser
	jpegMediaParser := jpegstructure.NewJpegMediaParser()

	// Parse the JPEG file
	mediaContext, err := jpegMediaParser.ParseFile(path)
	if err != nil {
		return nil, err
	}

	// Get the EXIF data
	_, rawExif, err := mediaContext.Exif()
	if err != nil {
		return nil, err
	}
	return rawExif, nil
}

// parsePng parses the EXIF data from the PNG file.
func (p *ExifParser) parsePng(path string) ([]byte, error) {
	// Create a new PNG media parser
	pngMediaParser := pngstructure.NewPngMediaParser()

	// Parse the PNG file
	mediaContext, err := pngMediaParser.ParseFile(path)
	if err != nil {
		return nil, err
	}

	// Get the EXIF data
	_, rawExif, err := mediaContext.Exif()
	if err != nil {
		return nil, err
	}
	return rawExif, nil
}

// parseTiff parses the EXIF data from the TIFF file.
func (p *ExifParser) parseTiff(path string) ([]byte, error) {
	// Create a new TIFF media parser
	tiffMediaParser := tiffstructure.NewTiffMediaParser()

	// Parse the TIFF file
	mediaContext, err := tiffMediaParser.ParseFile(path)
	if err != nil {
		return nil, err
	}

	// Get the EXIF data
	_, rawExif, err := mediaContext.Exif()
	if err != nil {
		return nil, err
	}
	return rawExif, nil
}
