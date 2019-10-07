package tuner

import (
  "fmt"
  "strconv"

  "github.com/jinzhu/gorm"

  "bdbu/models"
)

type Variable struct {
  Name      string
  Value     string
}

type TunerController struct {
  db      *gorm.DB
  Vars    []Variable
}

func NewTuner(drv models.Driver) *TunerController {
  return &TunerController{
    db: drv.GetDB(),
  }
}

func (gl *TunerController) Close() {
  gl.db.Close()
}

func (gl *TunerController) AddItem(vars Variable) {
  gl.Vars = append(gl.Vars, vars)
}

func (gl *TunerController) GetGlobals() *TunerController {
  gl.GetGlobal(stStatus)
  gl.GetGlobal(stVariables)

  return gl
}

func (gl *TunerController) GetGlobal(x int) {
  rows, err := gl.db.Raw(fmt.Sprintf("SHOW GLOBAL %s", global[x])).Rows()
  defer rows.Close()

  if err != nil {
    panic(err)
  }

  for rows.Next() {
    var name  string
    var value string

    rows.Scan(&name, &value)

    gl.AddItem(Variable{name, value})
  }
}

func (gl *TunerController) GetVar(x int) string {
  for _, v := range gl.Vars {
    if v.Name == vars[x] {
      return v.Value
    }
  }

  return ""
}

func (gl *TunerController) version() *TunerController {
  ver         := gl.GetVar(vrVersion)
  compileVer  := gl.GetVar(vrCompileVer)
  outInfoF("MySQL Version %s %s\n\n", ver, compileVer)

  return gl
}

func (gl *TunerController) uptime() *TunerController {
  uptime         := gl.GetVar(vrUptime)
  questions      := gl.GetVar(vrQuestions)
  threads        := gl.GetVar(vrThreadsConnected)

  x, _ := strconv.Atoi(uptime)
  y, _ := strconv.Atoi(questions)

  outColInfo(txtUptime, uptime)
  outColInfo(txtAvgQps, x / y)
  outColInfo(txtTotalQuestions, questions)
  outColInfo(txtThreadsConnected, threads)
  outInfo("")

  if x > 172800 {
    outInfo("Server has been running for over 48hrs. It should be safe to follow these recommendations")
  } else {
    outWarnTips("Server has not been running for at least 48hrs. It may not be safe to use these recommendations")
    outInfo("")
  }

  return gl
}

func (gl *TunerController) slow() *TunerController {
  outHeader("Slow queries")

  if gl.GetVar(vrSlowQueryLog) == "ON" {
    outInfo("The slow query log is enabled.")
  } else {
    outError("The slow query log is NOT enabled.")
  }

  outInfoF("Current long_query_time = %v sec.\n", gl.GetVar(vrLongQueryTime))
  outInfo(gl.GetVar(vrLogSlowQueries), " - ", gl.GetVar(vrQuestions))
  return gl
}

func Start(drv models.Driver) {
  outHeader("-- MYSQL PERFORMANCE TUNING --")

  tc := NewTuner(drv)
  defer tc.Close()

  tc.GetGlobals().
    version().
    uptime().
    slow()

}

const (
  stStatus int = iota
  stVariables
)

const (
  vrVersion int = iota
  vrCompileVer
  vrUptime
  vrQuestions
  vrThreadsConnected
  vrSlowQueryLog
  vrLongQueryTime
  vrLogSlowQueries
)

var (
  global = []string{"STATUS", "VARIABLES"}

  vars   = []string{
    "version",
    "version_compile_machine",
    "Uptime",
    "Questions",
    "Threads_connected",
    "slow_query_log",
    "long_query_time",
    "log_slow_queries",
  }
)
