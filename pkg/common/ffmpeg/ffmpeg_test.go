package ffmpeg

import (
	"fmt"
	"testing"
)

func TestGetOnlineVideoInfo(t *testing.T) {
	url := "http://106.53.99.30/dl/ptbtv.php?id=4"
	info, err := GetOnlineVideoInfo(url)
	fmt.Println(info, err)
}
