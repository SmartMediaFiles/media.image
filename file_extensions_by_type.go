package media_image

import "github.com/smartmediafiles/media/media/maps"

// ImageFileTypesExtensions is a map of media.Image file types to their file extensions.
var ImageFileTypesExtensions = maps.MapFileTypeExtensions{
	ImageBmp:  {ExtensionBmp, ExtensionDib},
	ImageGif:  {ExtensionGif},
	ImageHeic: {ExtensionHeic},
	ImageHeif: {ExtensionHeif},
	ImageJpeg: {ExtensionJpg, ExtensionJpeg, ExtensionJpe, ExtensionJif, ExtensionJfif, ExtensionJfi},
	ImagePng:  {ExtensionPng},
	ImageTiff: {ExtensionTiff, ExtensionTif},
	ImageWebp: {ExtensionWebp},
}
