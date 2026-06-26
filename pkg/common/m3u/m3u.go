package m3u

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/aiqoder/monitor-lite-api/pkg/common/fileutils"
)

// Playlist is a type that represents an m3u playlist containing 0 or more tracks
type Playlist struct {
	Tracks []Track
}

// A Tag is a simple key/value pair
type Tag struct {
	Name  string
	Value string
}

// Track represents an m3u track with a Name, Lengh, URI and a set of tags
type Track struct {
	Name   string
	Length int
	URI    string
	Group  string
	Tags   []Tag
}

// Parse parses an m3u playlist with the given file name and returns a Playlist
func Parse(fileName string) (Playlist, error) {
	var f io.ReadCloser

	if strings.HasPrefix(fileName, "http://") || strings.HasPrefix(fileName, "https://") {
		data, err := http.Get(fileName)
		if err != nil {
			return Playlist{},
				fmt.Errorf("unable to open playlist URL: %v", err)
		}
		f = data.Body
		defer f.Close()
	} else if fileutil.IsExist(fileName) {
		file, err := os.Open(fileName)
		defer f.Close()
		if err != nil {
			return Playlist{},
				fmt.Errorf("unable to open playlist file: %v", err)
		}
		f = file
	} else {
		u := time.Now().Unix()

		filePath := filepath.Join(os.TempDir(), "ygcz", strconv.FormatInt(u, 10))
		_ = fileutils.CreateFile(filePath, []byte(fileName))
		file, err := os.Open(filePath)

		if err != nil {
			return Playlist{},
				fmt.Errorf("unable to open playlist file: %v", err)
		}
		f = file

		defer func() {
			_ = file.Close()
			_ = os.Remove(filePath)
		}()
	}

	onFirstLine := true
	scanner := bufio.NewScanner(f)
	tagsRegExp, _ := regexp.Compile("([a-zA-Z0-9-]+?)=\"([^\"]+)\"")
	playlist := Playlist{}

	for scanner.Scan() {
		line := scanner.Text()
		// 跳过空行
		if len(line) == 0 {
			continue
		}
		if onFirstLine && !strings.HasPrefix(line, "#EXTM3U") {
			return Playlist{},
				errors.New("invalid m3u file format. Expected #EXTM3U file header")
		}

		onFirstLine = false

		if strings.HasPrefix(line, "#EXTINF") {
			// 如果包含多个逗号，只保留最后一个
			split := strings.Split(line, ",")
			if len(split) > 2 {
				line = strings.Join(split[0:len(split)-1], " ") + "," + split[len(split)-1]
			}

			line := strings.Replace(line, "#EXTINF:", "", -1)

			trackInfo := strings.Split(line, ",")

			if len(trackInfo) < 2 {
				return Playlist{},
					errors.New("invalid m3u file format. Expected EXTINF metadata to contain track length and name data")
			}
			length, parseErr := strconv.Atoi(strings.Split(trackInfo[0], " ")[0])
			if parseErr != nil {
				return Playlist{}, errors.New("unable to parse length")
			}
			track := &Track{strings.Trim(trackInfo[1], " "), length, "", "", nil}
			tagList := tagsRegExp.FindAllString(line, -1)
			for i := range tagList {
				tagInfo := strings.Split(tagList[i], "=")
				tag := &Tag{tagInfo[0], strings.Replace(tagInfo[1], "\"", "", -1)}
				// 更新group
				if tag.Name == "group-title" {
					track.Group = tag.Value
				}
				track.Tags = append(track.Tags, *tag)
			}
			playlist.Tracks = append(playlist.Tracks, *track)
		} else if strings.HasPrefix(line, "#") || line == "" {
			continue
		} else if len(playlist.Tracks) == 0 {
			return Playlist{},
				errors.New("URI provided for playlist with no tracks")

		} else {
			playlist.Tracks[len(playlist.Tracks)-1].URI = strings.Trim(line, " ")
		}
	}

	return playlist, nil
}

// Marshall Playlist to an m3u file.
func Marshall(p Playlist) (io.Reader, error) {
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	if err := MarshallInto(p, w); err != nil {
		return nil, err
	}

	return buf, nil
}

// MarshallInto a *bufio.Writer a Playlist.
func MarshallInto(p Playlist, into *bufio.Writer) error {
	into.WriteString("#EXTM3U\n")
	for _, track := range p.Tracks {
		into.WriteString("#EXTINF:")
		into.WriteString(fmt.Sprintf("%d ", track.Length))
		for i := range track.Tags {
			if i == len(track.Tags)-1 {
				into.WriteString(fmt.Sprintf("%s=%q", track.Tags[i].Name, track.Tags[i].Value))
				continue
			}
			into.WriteString(fmt.Sprintf("%s=%q ", track.Tags[i].Name, track.Tags[i].Value))
		}
		into.WriteString(", ")

		into.WriteString(fmt.Sprintf("%s\n%s\n", track.Name, track.URI))
	}

	return into.Flush()
}
