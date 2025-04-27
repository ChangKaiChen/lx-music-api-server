package crypto

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

var (
	part1Indexes   = []int{23, 14, 6, 36, 16, 40, 7, 19}
	part2Indexes   = []int{16, 1, 32, 12, 19, 27, 8, 5}
	scrambleValues = []int{
		89, 39, 179, 150, 218, 82, 58, 252, 177, 52,
		186, 123, 120, 64, 242, 133, 143, 161, 121, 179,
	}
)

func Sign(payload string, flag bool) (string, error) {
	var filteredPart1Indexes []int
	for _, idx := range part1Indexes {
		if idx < 40 {
			filteredPart1Indexes = append(filteredPart1Indexes, idx)
		}
	}
	h := sha1.New()
	h.Write([]byte(payload))
	hash := strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
	var part1 strings.Builder
	for _, idx := range filteredPart1Indexes {
		part1.WriteByte(hash[idx])
	}
	if flag {
		part1.Reset()
	}
	var part2 strings.Builder
	for _, idx := range part2Indexes {
		part2.WriteByte(hash[idx])
	}
	part3 := make([]byte, 20)
	for i, v := range scrambleValues {
		hexValue := hash[i*2 : i*2+2]
		intValue := 0
		_, err := fmt.Sscanf(hexValue, "%x", &intValue)
		if err != nil {
			return "", err
		}
		part3[i] = byte(v ^ intValue)
	}
	b64Part := base64.StdEncoding.EncodeToString(part3)
	re := regexp.MustCompile(`[/\\+=]`)
	b64Part = re.ReplaceAllString(b64Part, "")
	return fmt.Sprintf("zzc%s%s%s", strings.ToLower(part1.String()),
		strings.ToLower(b64Part), strings.ToLower(part2.String())), nil
}
