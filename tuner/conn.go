package tuner

import "strconv"

func (gl *TunerController) connection() *TunerController {
  outHeader("Max connections")

  connDefault, _ := strconv.Atoi(gl.GetVar(vrMaxConnections))
  connUsed, _    := strconv.Atoi(gl.GetVar(vrMaxUsedConnections))
  connRatio      := connUsed * 100 / connDefault

  outInfoF("Current max_connections = %v\n", connDefault)
  outInfoF("Current threads_connected = %v\n", gl.GetVar(vrThreadsConnected))
  outInfoF("Historic max_used_connections = %v\n", connUsed)

  errCode := 0

  outInfoF("The number of used connections is ")
  if connRatio > 85 {
    outErrorF("%v%s", connRatio, "%")
    errCode = 1
  } else if connRatio < 10 {
    outErrorF("%v%s", connRatio, "%")
    errCode = 1
  } else {
    outInfoF("%v%s", connRatio, "%")
  }
  outInfo(" of the configured maximum.")

  if errCode == 1 {
    outWarnTips("You should raise max_connections")
  } else if errCode == 2 {
    outWarnTips("You are using less than 10% of your configured max_connections.")
    outWarnTips("Lowering max_connections could help to avoid an over-allocation of memory")
    outWarnTips("See \"MEMORY USAGE\" section to make sure you are not over-allocating")
  } else {
    outInfoTips("Your max_connections variable seems to be fine.")
  }

  return gl
}
