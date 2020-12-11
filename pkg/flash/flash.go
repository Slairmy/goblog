package flash

import (
	"encoding/gob"
	"goblog/pkg/session"
)

// Flashes 闪存消息
type Flashes map[string]interface{}

var flashKey = "_flashes"

func init() {
	gob.Register(Flashes{})
}

func addFlash(key string, message string) {
	flashes := Flashes{}
	flashes[key] = message
	session.Put(flashKey, flashes)
	session.Save()
}

// All 返回所有的消息
func All() Flashes {
	val := session.Get(flashKey)

	flashMessages, ok := val.(Flashes)

	if !ok {
		return nil
	}

	session.Forget(flashKey)
	return flashMessages
}

// Info info类型消息
func Info(message string) {
	addFlash("info", message)
}

// Warning warning类型消息
func Warning(message string) {
	addFlash("warning", message)
}

// Success success类型消息
func Success(message string) {
	addFlash("success", message)
}

// Danger danger类型消息
func Danger(message string) {
	addFlash("danger", message)
}
