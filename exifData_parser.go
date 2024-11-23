package media_image

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	"github.com/go-mods/tags"
	"github.com/ringsaturn/tzf"
)

// Common time formats used in EXIF data
var timeFormats = []string{
	"2006:01:02 15:04:05",
	"2006:01:02 15:04:05-0700",
	"2006:01:02 15:04:05.000Z",
	"2006:01:02 15:04:05Z07:00",
}

// Common error messages
const (
	errNoGPSInfo    = "no GPS info found: %v"
	errParseGPSInfo = "failed to parse GPS info: %v"
	errParseTime    = "unable to parse time: %s"
	errParseTag     = "failed to parse tags for field %s: %v"
	errLoadTimezone = "failed to load timezone location: %v"
)

// ExifDataParser is responsible for extracting and parsing EXIF metadata from images.
// It maintains a cache of parsed tags to improve performance when processing multiple images.
type ExifDataParser struct {
	tagCache map[string][]*tags.Tag
}

// NewExifDataParser creates and initializes a new instance of ExifDataParser.
// It initializes the tag cache used to store parsed EXIF tags for better performance.
//
// Returns:
//   - *ExifDataParser: A pointer to the newly created parser instance
func NewExifDataParser() *ExifDataParser {
	return &ExifDataParser{
		tagCache: make(map[string][]*tags.Tag),
	}
}

// Global variables for EXIF parsing and timezone lookup
var (
	// exifIfdMapping stores the Image File Directory mapping information
	exifIfdMapping *exifcommon.IfdMapping

	// exifTagIndex maintains an index of all EXIF tags for quick lookup
	exifTagIndex = exif.NewTagIndex()

	// tzFinder is used to determine timezone information from GPS coordinates
	tzFinder tzf.F
)

// init initializes the global variables required for EXIF parsing.
// It loads standard IFD mappings and initializes the timezone finder.
// If initialization of critical components fails, the program will terminate.
func init() {
	// Initialize IFD mapping
	exifIfdMapping = exifcommon.NewIfdMapping()
	if err := exifcommon.LoadStandardIfds(exifIfdMapping); err != nil {
		log.Fatalf("Failed to load standard IFDs: %s", err)
	}

	// Initialize timezone finder
	var err error
	tzFinder, err = tzf.NewDefaultFinder()
	if err != nil {
		log.Printf("Warning: Failed to initialize timezone finder: %v", err)
	}
}

// getExifTags retrieves the EXIF tags for a given struct field.
// It uses a cache to avoid repeated parsing of the same tags.
//
// Parameters:
//   - field: The struct field to get EXIF tags for
//
// Returns:
//   - []*tags.Tag: Slice of parsed EXIF tags
//   - error: Any error encountered while parsing tags
func (p *ExifDataParser) getExifTags(field reflect.StructField) ([]*tags.Tag, error) {
	// Check cache first
	if cachedTags, ok := p.tagCache[field.Name]; ok {
		return cachedTags, nil
	}

	// Parse tags if not in cache
	parsedTags, err := tags.Parse(string(field.Tag))
	if err != nil {
		return nil, fmt.Errorf(errParseTag, field.Name, err)
	}

	// Filter and store only EXIF tags
	var exifTags []*tags.Tag
	for _, tag := range parsedTags {
		if tag.Key == "exif" {
			exifTags = append(exifTags, tag)
		}
	}

	// Cache the results
	p.tagCache[field.Name] = exifTags
	return exifTags, nil
}

// getValueFromMetadata searches for the first non-empty value among the given tags
// in the metadata map.
//
// Parameters:
//   - metadata: Map of EXIF tag names to their values
//   - fieldTags: Slice of tags to search for
//
// Returns:
//   - string: The found value
//   - bool: Whether a value was found
func (p *ExifDataParser) getValueFromMetadata(metadata map[string]string, fieldTags []*tags.Tag) (string, bool) {
	for _, tag := range fieldTags {
		names := strings.Split(tag.Value, ",")
		for _, name := range names {
			if value, ok := metadata[name]; ok && value != "" {
				return value, true
			}
		}
	}
	return "", false
}

// Parse extracts and processes EXIF metadata from raw image data.
// It performs a comprehensive extraction of all available EXIF information
// and organizes it into a structured ImageData object.
//
// Parameters:
//   - exifData: Raw EXIF data bytes from the image
//
// Returns:
//   - ImageData: Structured representation of the extracted metadata
//   - error: Any error encountered during parsing
func (p *ExifDataParser) Parse(exifData []byte) (ImageData, error) {
	// Extract all EXIF entries and build metadata map
	metadata, err := p.buildMetadataMap(exifData)
	if err != nil {
		return ImageData{}, fmt.Errorf("failed to extract EXIF data: %v", err)
	}

	// Build IFD index for structured access to EXIF data
	var ifdIndex exif.IfdIndex
	_, ifdIndex, err = exif.Collect(exifIfdMapping, exifTagIndex, exifData)
	if err != nil {
		return ImageData{}, fmt.Errorf("failed to build IFD index: %v", err)
	}

	return p.parseWithReflection(metadata, ifdIndex)
}

// buildMetadataMap creates a map of EXIF tag names to their values from raw EXIF data.
// It handles null-terminated strings and filters out empty or invalid entries.
//
// Parameters:
//   - exifData: Raw EXIF data bytes
//
// Returns:
//   - map[string]string: Processed metadata map
//   - error: Any error encountered during extraction
func (p *ExifDataParser) buildMetadataMap(exifData []byte) (map[string]string, error) {
	metadata := make(map[string]string)

	entries, _, err := exif.GetFlatExifDataUniversalSearch(exifData, nil, true)
	if err != nil {
		return nil, err
	}

	//// affiche dans la console les données exif
	//for _, entry := range entries {
	//	fmt.Printf("entry: %v\n", entry)
	//}

	for _, entry := range entries {
		// get formatted tag entry
		s := strings.Split(entry.FormattedFirst, "\x00")

		// Skip empty or invalid entries
		if entry.TagName == "" || len(s) == 0 {
			continue
		}

		// Ignore IFD1 data.exif as it is usually a thumbnail
		if entry.IfdPath == exif.ThumbnailFqIfdPath {
			continue
		}

		// Handle null-terminated strings
		if len(s) > 0 && s[0] != "" {
			metadata[entry.TagName] = s[0]
		}
	}

	return metadata, nil
}

// parseWithReflection processes the metadata map and IFD index using reflection
// to populate an ImageData struct with the extracted information.
//
// Parameters:
//   - metadata: Map of EXIF tag names to their values
//   - ifdIndex: Index of Image File Directory information
//
// Returns:
//   - ImageData: Populated structure containing the image metadata
//   - error: Any error encountered during processing
func (p *ExifDataParser) parseWithReflection(metadata map[string]string, ifdIndex exif.IfdIndex) (ImageData, error) {
	imageData := ImageData{}
	v := reflect.ValueOf(&imageData).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		// Skip GPS fields as they are handled separately
		if p.isSpecialField(field.Name) {
			continue
		}

		// Get and validate EXIF tags for the field
		fieldTags, err := p.getExifTags(field)
		if err != nil {
			log.Printf("Warning: %v", err)
			continue
		}

		if len(fieldTags) == 0 {
			continue
		}

		// Extract and set field value
		if value, ok := p.getValueFromMetadata(metadata, fieldTags); ok {
			if err := p.setFieldValue(fieldValue, value); err != nil {
				log.Printf("Warning: failed to set field %s: %v", field.Name, err)
			}
		}
	}

	// Process GPS information separately due to its complex nature
	if err := p.extractGPSInfo(&imageData, metadata, ifdIndex); err != nil {
		log.Printf("Warning: GPS extraction failed: %v", err)
	}

	return imageData, nil
}

// isSpecialField determines if a field requires special handling
// and should not be processed using the standard reflection approach.
// All GPS-related fields are considered special and handled separately.
//
// Parameters:
//   - fieldName: Name of the field to check
//
// Returns:
//   - bool: True if the field requires special handling
func (p *ExifDataParser) isSpecialField(fieldName string) bool {
	// All GPS fields are considered special
	return strings.HasPrefix(fieldName, "GPS")
}

// setFieldValue sets a field's value based on its type and the provided string value.
// It handles various data types including strings, integers, floats, and time.Time.
//
// Parameters:
//   - field: Reflect.Value of the field to set
//   - value: String value to parse and set
//
// Returns:
//   - error: Any error encountered while setting the value
func (p *ExifDataParser) setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse int: %v", err)
		}
		field.SetInt(i)
	case reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float: %v", err)
		}
		field.SetFloat(f)
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			t, err := p.parseTime(value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(t))
		} else if field.Type() == reflect.TypeOf(Rational{}) {
			r, err := NewRational(value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(r))
		}
	}
	return nil
}

// parseTime attempts to parse a time string using multiple common EXIF time formats.
// It iterates through known formats until it finds one that successfully parses the input.
//
// Parameters:
//   - value: The time string to parse
//
// Returns:
//   - time.Time: The parsed time value
//   - error: Any error encountered during parsing
func (p *ExifDataParser) parseTime(value string) (time.Time, error) {
	var lastErr error
	for _, format := range timeFormats {
		if t, err := time.Parse(format, value); err == nil {
			return t, nil
		} else {
			lastErr = err
		}
	}
	return time.Time{}, fmt.Errorf("%s: %v", errParseTime, lastErr)
}

// extractGPSInfo processes and extracts GPS-related information from EXIF data.
// This includes coordinates, altitude, timestamp, timezone information, and all other GPS fields.
//
// Parameters:
//   - imageData: Pointer to the ImageData struct to populate
//   - metadata: Map of EXIF tag names to their values
//   - ifdIndex: Index of Image File Directory information
//
// Returns:
//   - error: Any error encountered during GPS data extraction
func (p *ExifDataParser) extractGPSInfo(imageData *ImageData, metadata map[string]string, ifdIndex exif.IfdIndex) error {
	// Get GPS IFD (Image File Directory)
	ifd, err := ifdIndex.RootIfd.ChildWithIfdPath(exifcommon.IfdGpsInfoStandardIfdIdentity)
	if err != nil {
		return fmt.Errorf(errNoGPSInfo, err)
	}

	// Extract GPS info using the dedicated GPS parser
	gpsInfo, err := ifd.GpsInfo()
	if err != nil {
		return fmt.Errorf(errParseGPSInfo, err)
	}

	// Process coordinates and timezone
	if err := p.processGPSCoordinates(imageData, gpsInfo); err != nil {
		return err
	}

	// Set altitude if available
	if gpsInfo.Altitude != 0 {
		imageData.GPSAltitude = float64(gpsInfo.Altitude)
	}

	// Set GPS timestamp and process local time
	if !gpsInfo.Timestamp.IsZero() {
		imageData.GPSTimestamp = gpsInfo.Timestamp
		p.processLocalTime(imageData)
	}

	// Process additional GPS metadata
	p.processAdditionalGPSMetadata(imageData, metadata)

	return nil
}

// processAdditionalGPSMetadata handles the extraction of additional GPS-related metadata
// that is not covered by the standard GPS parser.
//
// Parameters:
//   - imageData: Pointer to the ImageData struct to populate
//   - metadata: Map of EXIF tag names to their values
func (p *ExifDataParser) processAdditionalGPSMetadata(imageData *ImageData, metadata map[string]string) {
	// Processing method
	if method, ok := metadata["GPSProcessingMethod"]; ok {
		imageData.GPSProcessingMethod = method
	}

	// Status
	if status, ok := metadata["GPSStatus"]; ok {
		imageData.GPSStatus = status
	}

	// Satellites
	if satellites, ok := metadata["GPSSatellites"]; ok {
		imageData.GPSSatellites = satellites
	}

	// Positioning error
	if hError, ok := metadata["GPSHPositioningError"]; ok {
		imageData.GPSHPositioningError, _ = strconv.ParseFloat(hError, 64)
	}

	// Movement information
	if speed, ok := metadata["GPSSpeed"]; ok {
		imageData.GPSSpeed, _ = strconv.ParseFloat(speed, 64)
	}
	if track, ok := metadata["GPSTrack"]; ok {
		imageData.GPSTrack, _ = strconv.ParseFloat(track, 64)
	}
	if imgDir, ok := metadata["GPSImgDirection"]; ok {
		imageData.GPSImgDirection, _ = strconv.ParseFloat(imgDir, 64)
	}

	// Destination information
	if destLat, ok := metadata["GPSDestLatitude"]; ok {
		imageData.GPSDestLatitude, _ = strconv.ParseFloat(destLat, 64)
	}
	if destLong, ok := metadata["GPSDestLongitude"]; ok {
		imageData.GPSDestLongitude, _ = strconv.ParseFloat(destLong, 64)
	}
	if bearing, ok := metadata["GPSDestBearing"]; ok {
		imageData.GPSDestBearing, _ = strconv.ParseFloat(bearing, 64)
	}
	if distance, ok := metadata["GPSDestDistance"]; ok {
		imageData.GPSDestDistance, _ = strconv.ParseFloat(distance, 64)
	}
}

// processGPSCoordinates handles the extraction and validation of GPS coordinates
// and associated timezone information.
//
// Parameters:
//   - imageData: Pointer to the ImageData struct to populate
//   - gpsInfo: GPS information from EXIF data
//
// Returns:
//   - error: Any error encountered during processing
func (p *ExifDataParser) processGPSCoordinates(imageData *ImageData, gpsInfo *exif.GpsInfo) error {
	// Validate and set coordinates
	if math.IsNaN(gpsInfo.Latitude.Decimal()) || math.IsNaN(gpsInfo.Longitude.Decimal()) {
		return fmt.Errorf("invalid GPS coordinates")
	}

	imageData.GPSLatitude = gpsInfo.Latitude.Decimal()
	imageData.GPSLongitude = gpsInfo.Longitude.Decimal()

	// Get timezone from coordinates if possible
	if tzFinder != nil {
		timezoneName := tzFinder.GetTimezoneName(
			imageData.GPSLongitude,
			imageData.GPSLatitude,
		)
		if timezoneName != "" {
			imageData.GPSTimeZone = timezoneName
			// Adjust all time fields with the found timezone
			p.adjustTimeWithTimezone(imageData)
		}
	}

	return nil
}

// processLocalTime attempts to create a local timestamp using the GPS timezone
// if both GPS timestamp and timezone information are available.
//
// Parameters:
//   - imageData: Pointer to the ImageData struct containing GPS information
func (p *ExifDataParser) processLocalTime(imageData *ImageData) {
	if imageData.GPSTimeZone == "" || imageData.GPSTimestamp.IsZero() {
		return
	}

	loc, err := time.LoadLocation(imageData.GPSTimeZone)
	if err != nil {
		log.Printf("Warning: %s", fmt.Sprintf(errLoadTimezone, err))
		return
	}

	imageData.GPSTimestampLocal = imageData.GPSTimestamp.In(loc)
}

// adjustTimeWithTimezone updates all time fields with the timezone information
// when available. This includes DateTimeOriginal and DateTimeDigitized.
//
// Parameters:
//   - imageData: Pointer to the ImageData struct to update
func (p *ExifDataParser) adjustTimeWithTimezone(imageData *ImageData) {
	// Skip if no timezone was found
	if imageData.GPSTimeZone == "" {
		return
	}

	// Load the location for the timezone
	loc, err := time.LoadLocation(imageData.GPSTimeZone)
	if err != nil {
		log.Printf("Warning: %s", fmt.Sprintf(errLoadTimezone, err))
		return
	}

	// Calculate the timezone offset for the current time
	now := time.Now().UTC().In(loc)
	_, offset := now.Zone()

	// Format the offset as "+HHMM" or "-HHMM"
	sign := "+"
	if offset < 0 {
		sign = "-"
		offset = -offset
	}
	hours := offset / 3600
	minutes := (offset % 3600) / 60
	imageData.TimeOffset = fmt.Sprintf("%s%02d%02d", sign, hours, minutes)
	imageData.HasTimeOffset = true

	// Adjust DateTimeOriginal if it exists
	if !imageData.DateTimeOriginal.IsZero() {
		// If the time already has a timezone, convert it
		if imageData.DateTimeOriginal.Location() != time.UTC {
			imageData.DateTimeOriginal = imageData.DateTimeOriginal.In(loc)
		} else {
			// If the time is in UTC, treat it as local time in the new timezone
			imageData.DateTimeOriginal = time.Date(
				imageData.DateTimeOriginal.Year(),
				imageData.DateTimeOriginal.Month(),
				imageData.DateTimeOriginal.Day(),
				imageData.DateTimeOriginal.Hour(),
				imageData.DateTimeOriginal.Minute(),
				imageData.DateTimeOriginal.Second(),
				imageData.DateTimeOriginal.Nanosecond(),
				loc,
			)
		}
	}

	// Adjust DateTimeDigitized if it exists
	if !imageData.DateTimeDigitized.IsZero() {
		// If the time already has a timezone, convert it
		if imageData.DateTimeDigitized.Location() != time.UTC {
			imageData.DateTimeDigitized = imageData.DateTimeDigitized.In(loc)
		} else {
			// If the time is in UTC, treat it as local time in the new timezone
			imageData.DateTimeDigitized = time.Date(
				imageData.DateTimeDigitized.Year(),
				imageData.DateTimeDigitized.Month(),
				imageData.DateTimeDigitized.Day(),
				imageData.DateTimeDigitized.Hour(),
				imageData.DateTimeDigitized.Minute(),
				imageData.DateTimeDigitized.Second(),
				imageData.DateTimeDigitized.Nanosecond(),
				loc,
			)
		}
	}

	// Also adjust DateTime if it exists
	if !imageData.DateTime.IsZero() {
		// If the time already has a timezone, convert it
		if imageData.DateTime.Location() != time.UTC {
			imageData.DateTime = imageData.DateTime.In(loc)
		} else {
			// If the time is in UTC, treat it as local time in the new timezone
			imageData.DateTime = time.Date(
				imageData.DateTime.Year(),
				imageData.DateTime.Month(),
				imageData.DateTime.Day(),
				imageData.DateTime.Hour(),
				imageData.DateTime.Minute(),
				imageData.DateTime.Second(),
				imageData.DateTime.Nanosecond(),
				loc,
			)
		}
	}
}
