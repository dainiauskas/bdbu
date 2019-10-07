package tuner

import (
  "strings"
  "fmt"

  "github.com/gookit/color"
)

func init() {
  x := 0

  for _, s := range columns {
    if x < len(s) {
      x = len(s)
    }
  }

  txtTemplate = "%" + fmt.Sprintf("%d", x) + "s: %v\n"
}

func outHeader(msg string) {
  color.Notice.Println(strings.ToUpper(msg))
}

func outInfo(args ...interface{}) {
  color.Info.Println(args...)
}

func outInfoF(msg string, args ...interface{}) {
  color.Info.Printf(msg, args...)
}

func outError(msg ...interface{}) {
  color.Info.Println(msg...)
}

func outInfoTips(msg string, args ...interface {}) {
  color.Info.Tips(msg, args...)
}

func outWarnTips(msg string, args ...interface {}) {
  color.Warn.Tips(msg, args...)
}

func outColInfo(hdr int, arg interface{}) {
  color.Comment.Printf(txtTemplate, columns[hdr], arg)
}

const (
  txtUptime int = iota
  txtAvgQps
  txtTotalQuestions
  txtThreadsConnected
)

var (
  txtTemplate string

  columns = []string{
    "Uptime",
    "Avg. qps",
    "Total Questions",
    "Threads Connected",
  }
)
