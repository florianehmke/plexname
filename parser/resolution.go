package parser

// Resolution is the internal representation of a resolution.
type Resolution int

// All known resolutions.
const (
	ResNA Resolution = iota
	R360
	R480
	R540
	R576
	R720
	R1080
	R2160
)

// Resolutions mapped to their string representations.
var (
	resNames = map[Resolution]string{
		ResNA: "--",
		R360:  "360p",
		R480:  "480p",
		R540:  "540p",
		R576:  "576p",
		R720:  "720p",
		R1080: "1080p",
		R2160: "2160p",
	}

	resMap = map[string]Resolution{
		"360p":    R360,
		"640x480": R480,
		"848x480": R480,
		"480p":    R480,
		"540p":    R540,
		"576p":    R576,
		"720p":    R720,
		"1080p":   R1080,
		"2160p":   R2160,
	}
)

// String returns the string representation of r.
func (r Resolution) String() string {
	return resNames[r]
}

