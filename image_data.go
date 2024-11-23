package media_image

import "time"

type ImageData struct {
	// GPS information extracted from the EXIF data
	GPSLatitude          float64   `exif:"GPSLatitude"`
	GPSLongitude         float64   `exif:"GPSLongitude"`
	GPSAltitude          float64   `exif:"GPSAltitude"`
	GPSTimeZone          string    // Determined from coordinates
	GPSTimestamp         time.Time `exif:"GPSDateStamp,GPSTimeStamp"`
	GPSTimestampLocal    time.Time // Computed from GPSTimestamp and GPSTimeZone
	GPSProcessingMethod  string    `exif:"GPSProcessingMethod"`
	GPSStatus            string    `exif:"GPSStatus"`
	GPSSatellites        string    `exif:"GPSSatellites"`
	GPSHPositioningError float64   `exif:"GPSHPositioningError"`
	GPSSpeed             float64   `exif:"GPSSpeed"`
	GPSTrack             float64   `exif:"GPSTrack"`
	GPSImgDirection      float64   `exif:"GPSImgDirection"`
	GPSDestLatitude      float64   `exif:"GPSDestLatitude"`
	GPSDestLongitude     float64   `exif:"GPSDestLongitude"`
	GPSDestBearing       float64   `exif:"GPSDestBearing"`
	GPSDestDistance      float64   `exif:"GPSDestDistance"`

	// Camera information extracted from the EXIF data
	CameraMake        string    `exif:"Make,CameraMake"`
	CameraModel       string    `exif:"Model,CameraModel"`
	CameraExposure    string    `exif:"ExposureTime,Exposure"`
	ISOSpeed          int       `exif:"ISOSpeedRatings,ISO"`
	ShutterSpeed      string    `exif:"ShutterSpeedValue"`
	Software          string    `exif:"Software"`
	DateTime          time.Time `exif:"DateTime,CreateDate"`
	DateTimeOriginal  time.Time `exif:"DateTimeOriginal,OriginalDateTime"`
	DateTimeDigitized time.Time `exif:"DateTimeDigitized,DigitizedDateTime"`
	TimeOffset        string    `exif:"OffsetTime,OffsetTimeOriginal,OffsetTimeDigitized"` // Format: "+0200" or "-0700"
	SubSecOriginal    string    `exif:"SubSecTimeOriginal,SubSecTime"`                     // Subsecond precision
	HasTimeOffset     bool      // Indicates if time offset was found

	// Lens information extracted from the EXIF data
	LensMake            string `exif:"LensMake"`
	LensModel           string `exif:"LensModel,Lens"`
	LensFocalLength     string `exif:"FocalLength"`
	LensAperture        string `exif:"FNumber,ApertureValue"`
	LensFocalLength35mm string `exif:"FocalLengthIn35mmFilm"`
	LensMaxAperture     string `exif:"MaxApertureValue"`
	LensMinAperture     string `exif:"MinApertureValue"`
	LensMaxFocalLength  string `exif:"MaxFocalLength"`

	// Image information
	ImageWidth       int      `exif:"ImageWidth,PixelXDimension,ExifImageWidth,SourceImageWidth"`
	ImageHeight      int      `exif:"ImageHeight,PixelYDimension,ExifImageHeight,SourceImageHeight"`
	ImageOrientation int      `exif:"Orientation"`
	ColorSpace       string   `exif:"ColorSpace"`
	Compression      string   `exif:"Compression"`
	XResolution      Rational `exif:"XResolution"`
	YResolution      Rational `exif:"YResolution"`
	ResolutionUnit   string   `exif:"ResolutionUnit"`

	// Additional EXIF information
	Artist           string  `exif:"Artist,Creator"`
	Copyright        string  `exif:"Copyright,CopyrightNotice"`
	Description      string  `exif:"ImageDescription,Description"`
	WhiteBalance     string  `exif:"WhiteBalance"`
	Flash            string  `exif:"Flash,FlashFired"`
	MeteringMode     string  `exif:"MeteringMode"`
	ExposureProgram  string  `exif:"ExposureProgram"`
	SceneCaptureType string  `exif:"SceneCaptureType"`
	SubjectDistance  float64 `exif:"SubjectDistance"`
	DigitalZoomRatio float64 `exif:"DigitalZoomRatio"`
}
