package util

import (
  "golang.org/x/text/language"
  "golang.org/x/text/message"
)

var (
  lang = language.English
  printer = message.NewPrinter(lang)
)

func NumberFormatLong(num uint64) string {
  return printer.Sprintf("%d", num)
}
