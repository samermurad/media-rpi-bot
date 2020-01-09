package config

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

var APPROVAL_REG = regexp.MustCompile("(?i)^y$")

func envOrPanic(key string) string {
	if _env := os.Getenv(key); _env != "" {
		return _env
	}
	panic("Must Set " + key)
}

func envOrDefault(key string, defaultV string) string {
	if _env := os.Getenv(key); _env != "" {
		return _env
	}
	return defaultV
}
func BOT_TOKEN() string {
	return envOrPanic("BOT_TOKEN")
}

func MEDIA_SRC_FOLDER() string {
	return envOrPanic("MEDIA_SRC_FOLDER")
}

func MEDIA_DEST_FOLDER() string {
	return envOrPanic("MEDIA_DEST_FOLDER")
}

var chtsId []int64

func ALLOWED_CHATS_IDS() []int64 {
	if chtsId != nil {
		return chtsId
	}
	str := envOrPanic("ALLOWED_CHATS_IDS")
	strs := strings.Split(str, ",")
	chtsId = make([]int64, 0)
	for _, v := range strs {
		_v, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		chtsId = append(chtsId, _v)
	}
	return chtsId
}

var chatOffset int64 = 0

func CHAT_OFFSET() int64 {
	return chatOffset
}

func SET_CHAT_OFFSET(i int64) {
	chatOffset = i
}
