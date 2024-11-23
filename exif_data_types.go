package media_image

import (
	"fmt"
	"strconv"
	"strings"
)

// Rational represents a rational number with a numerator and a denominator.
type Rational struct {
	Numerator   int
	Denominator int
}

func NewRational(s string) (Rational, error) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return Rational{}, fmt.Errorf("invalid format for Rational: %s", s)
	}
	numerator, err := strconv.Atoi(parts[0])
	if err != nil {
		return Rational{}, err
	}
	denominator, err := strconv.Atoi(parts[1])
	if err != nil {
		return Rational{}, err
	}
	return Rational{Numerator: numerator, Denominator: denominator}, nil
}

// String converts a Rational to a string formatted as "numerator/denominator".
func (r Rational) String() string {
	return fmt.Sprintf("%d/%d", r.Numerator, r.Denominator)
}
