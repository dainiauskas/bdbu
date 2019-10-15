package tuner

import (
  "fmt"
  "strconv"

  "github.com/jinzhu/gorm"
  "github.com/dustin/go-humanize"
  "github.com/pbnjay/memory"

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

func (gl *TunerController) ParseBytes(x int) (value uint64) {
  value, _ = humanize.ParseBytes(gl.GetVar(x))
  return
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

  qTime := gl.GetVar(vrLongQueryTime)
  outInfoF("Current long_query_time = %v sec.\n", qTime)
  outInfoF("You have %v  out of %v that take longer than %v sec. to complete\n",
    gl.GetVar(vrSlowQueries), gl.GetVar(vrQuestions), qTime)

  x, _ := strconv.Atoi(qTime)
  if x > 5 {
    outWarnTips("Your long_query_time may be too high, I typically set this under %d sec.", 5)
  } else {
    outInfoTips("Your long_query_time seems to be fine.")
  }

  return gl
}

func (gl *TunerController) binlog() *TunerController {
  outHeader("Binary update log")

  if gl.GetVar(vrBinLog) != "ON" {
    outWarnTips("The binary update log is NOT enabled.")
    outWarnTips("You will not be able to do point in time recovery.")
    return gl
  }

  outInfo("The binary update log is enabled")

  x, _ := strconv.Atoi(gl.GetVar(vrBinMaxSize))
  if x == 0 {
    outWarnTips("The max_binlog_size is not set. The binary log will rotate when it reaches 1GB.")
  }

  x, _ = strconv.Atoi(gl.GetVar(vrBinExpire))
  if x == 0 {
    outWarnTips("The expire_logs_days is not set.")
    outWarnTips("The mysqld will retain the entire binary log until RESET MASTER or PURGE MASTER LOGS commands are run manually")
    outWarnTips("Setting expire_logs_days will allow you to remove old binary logs automatically")
  }

  x, _ = strconv.Atoi(gl.GetVar(vrBinSync))
  if x == 0 {
    outWarnTips("Binlog sync is not enabled, you could loose binlog records during a server crash")
  }

  return gl
}

func (gl *TunerController) threads() *TunerController {
  outHeader("WORKER THREADS")

  thCreated, _ := strconv.Atoi(gl.GetVar(vrThreadsCreated))
  thCached     := gl.GetVar(vrThreadsCached)
  thCacheSz    := gl.GetVar(vrThreadCacheSize)
  uptime, _    := strconv.Atoi(gl.GetVar(vrUptime))
  thHist       := thCreated/uptime

  outInfoF("Current thread_cache_size = %v\n", thCacheSz)
  outInfoF("Current threads_cached = %v\n", thCached)
  outInfoF("Historic threads_per_sec = %v\n", thHist)

  if thHist > 2 {
    outWarnTips("%v\n%v",
      "Threads created per/sec are overrunning threads cached",
      "You should raise thread_cache_size",
    )
  } else {
    outInfoTips("Your thread_cache_size is fine")
  }

  return gl
}

func (gl *TunerController) keybuffer() *TunerController {
  // For MyIsam maybe letter
  outHeader("KEY BUFFER")

  row := gl.db.Raw("/*!50000 SELECT IFNULL(SUM(INDEX_LENGTH),0) as length from information_schema.TABLES where ENGINE='MyISAM' */").Row()

  var size int
  row.Scan(&size)

  outInfoF("Size: %v\n", size)
  outInfoF("Dir: %v\n", gl.GetVar(vrDataDir))

  return gl
}

func (gl *TunerController) queryCache() *TunerController {
  outHeader("Query cache")

  // version := gl.GetVar(vrVersion)

  qSize, _   := strconv.Atoi(gl.GetVar(vrQueryCacheSize))
  qLimit, _  := strconv.Atoi(gl.GetVar(vrQueryCacheLimit))
  qMinRes, _ := strconv.Atoi(gl.GetVar(vrQueryCacheMinRes))
  //
  qcFreeM, _ := strconv.Atoi(gl.GetVar(vrQcacheFreeMemory))
  qcTotal, _ := strconv.Atoi(gl.GetVar(vrQcacheTotalBlocks))
  qcFree, _  := strconv.Atoi(gl.GetVar(vrQcacheFreeBlocks))
  qLowMem, _ := strconv.Atoi(gl.GetVar(vrQcacheLowMemPrunes))

  if qSize == 0 {
    outWarnTips("%v %v\n",
      "Query cache is supported but not enabled",
  		"Perhaps you should set the query_cache_size but not recommend.",
    )

    return gl
  }

  usedMem  := qSize - qcFreeM
  memRatio := usedMem * 100 / qSize

  outInfoTips("Query cache is enabled")
  outInfoF("Current query_cache_size = %d K\n", qSize / 1024)
  outInfoF("Current query_cache_used = %d K\n", usedMem / 1024)
  outInfoF("Current query_cache_limit = %d K\n", qLimit / 1024)
  outInfoF("Current Query cache Memory fill ratio = %v%s\n", memRatio, "%")

  if qMinRes == 0 {
    outWarnTips("No query_cache_min_res_unit is defined.  Using MySQL < 4.1 cache fragmentation can be inpredictable")
  } else {
    outInfoF("Current query_cache_min_res_unit = %d K\n", qMinRes / 1024)
  }

  if qcFree > 2 && qcTotal > 0 {
    x := (qcFree * 100 / qcTotal)
    if x > 20 {
      outWarnTips("Query Cache is %v %v fragmented\n%v\n%v", x, "%",
        "Run \"FLUSH QUERY CACHE\" periodically to defragment the query cache memory",
        "If you have many small queries lower 'query_cache_min_res_unit' to reduce fragmentation.",
      )
    }
  }

  if memRatio < 25 {
    outWarnTips("Your query_cache_size seems to be too high. Perhaps you can use these resources elsewhere")
  }

  if qLowMem > 50 && memRatio > 80 {
    outWarnTips("However, %v %v %v", qLowMem,
      "queries have been removed from the query cache due to lack of memory",
      "Perhaps you should raise query_cache_size",
    )
  }

  return gl
}

func (gl *TunerController) innodb() *TunerController {
  outHeader("InnoDB")

  if gl.GetVar(vrInnoDBIgnoreBuiltIn) == "ON" {
    outWarnTips("No InnoDB Support Enabled!")
    return gl
  }

  var innodb_indexes, innodb_data uint64

  gl.db.Raw("/*!50000 SELECT IFNULL(SUM(INDEX_LENGTH),0) as size from information_schema.TABLES where ENGINE='InnoDB' */").
    Row().Scan(&innodb_indexes)

  gl.db.Raw("/*!50000 SELECT IFNULL(SUM(DATA_LENGTH),0) as size from information_schema.TABLES where ENGINE='InnoDB' */").
    Row().Scan(&innodb_data)

  if innodb_indexes == 0 {
    outWarnTips("InnoDB indexes empty")
    return gl
  }

  InnoDBBufferPoolSize, _        := humanize.ParseBytes(gl.GetVar(vrInnoDBBufferPoolSize))
  InnoDBLogFileSize, _           := humanize.ParseBytes(gl.GetVar(vrInnoDBLogFileSize))
  InnoDBLogBufferSize, _         := humanize.ParseBytes(gl.GetVar(vrInnoDBLogBufferSize))
  InnoDBBufferPoolPagesFree, _   := humanize.ParseBytes(gl.GetVar(vrInnoDBBufferPoolPagesFree))
  InnoDBBufferPoolPagesTotal, _  := humanize.ParseBytes(gl.GetVar(vrInnoDBBufferPoolPagesTotal))
  InnoDBFlushLogAtTrxCommit, _   := humanize.ParseBytes(gl.GetVar(vrInnoDBFlushLogAtTrxCommit))

  outInfoF("Current InnoDB index space = %v\n", humanize.Bytes(uint64(innodb_indexes)))
  outInfoF("Current InnoDB data space = %v\n", humanize.Bytes(innodb_data))
  outInfoF("Current InnoDB buffer pool free = %d%s\n", InnoDBBufferPoolPagesFree * 100 / InnoDBBufferPoolPagesTotal, "%")
  outInfoF("Current innodb_buffer_pool_size = %v\n", humanize.Bytes(InnoDBBufferPoolSize))

  logRate := InnoDBLogFileSize * 100 / InnoDBBufferPoolSize
  if logRate < 25 {
    outWarnTips("Current innodb_log_file_size %v is LOW", humanize.Bytes(InnoDBLogFileSize))
  } else if logRate > 50 {
    outWarnTips("Current innodb_log_file_size %v is HIGH", humanize.Bytes(InnoDBLogFileSize))
  }

  outInfoF("Current innodb_log_buffer_size = %v\n", humanize.Bytes(InnoDBLogBufferSize))
  if (InnoDBLogBufferSize / 1024 / 1024) != 64 {
    outWarnTips("It is suggested to set innodb_log_buffer_size to 64M on all servers.")
  }

  outInfoF("Current innodb_flush_log_at_trx_commit = %v\n", InnoDBFlushLogAtTrxCommit)
  if InnoDBFlushLogAtTrxCommit != 2 {
    outWarnTips("%v %v",
      "Writes the log buffer out to file on each commit but flushes to disk every second.",
      "If the disk cache has a battery backup (for instance a battery backed cache raid controller) this is generally the best balance of performance and safety.",
    )
  }
  outInfo("Depending on how much space your innodb indexes take up it may be safe")
  outInfo("to increase this value to up to 2 / 3 of total system memory")

  return gl
}

func (gl *TunerController) memory() *TunerController {
  outHeader("Memory usage")

  var binlogCacheSize, effectiveTmpTableSize uint64

  if gl.GetVar(vrBinLog) == "ON" {
    binlogCacheSize = gl.ParseBytes(vrBinlogCacheSize)
  }

  maxHeapTableSize    := gl.ParseBytes(vrMaxHeapTableSize)
  tmpTableSize        := gl.ParseBytes(vrTmpTableSize)
  readBufferSize      := gl.ParseBytes(vrReadBufferSize)
  readRndBufferSize   := gl.ParseBytes(vrReadRndBufferSize)
  sortBufferSize      := gl.ParseBytes(vrSortBufferSize)
  threadStack         := gl.ParseBytes(vrThreadStack)
  joinBufferSize      := gl.ParseBytes(vrJoinBufferSize)
  maxConnections      := gl.ParseBytes(vrMaxConnections)
  maxUsedConnections  := gl.ParseBytes(vrMaxUsedConnections)
  keyBufferSize       := gl.ParseBytes(vrKeyBufferSize)
  queryCacheSize      := gl.ParseBytes(vrQueryCacheSize)

  innoDBBufferPoolSize        := gl.ParseBytes(vrInnoDBBufferPoolSize)
  innoDBAdditionalMemPoolSize := gl.ParseBytes(vrInnoDBAdditionalMemPoolSize)
  innoDBLogBufferSize         := gl.ParseBytes(vrInnoDBLogBufferSize)

  if maxHeapTableSize <= tmpTableSize {
    effectiveTmpTableSize = maxHeapTableSize
  } else {
    effectiveTmpTableSize = tmpTableSize
  }

  bufferSize := (readBufferSize + readRndBufferSize + sortBufferSize +
    threadStack + joinBufferSize + binlogCacheSize)

  perThreadBuffers    := bufferSize * maxConnections
  perThreadMaxBuffers := bufferSize * maxUsedConnections

  globalBuffers := innoDBBufferPoolSize + innoDBAdditionalMemPoolSize +
    innoDBLogBufferSize + keyBufferSize + queryCacheSize


  maxMemory       := globalBuffers + perThreadMaxBuffers
  totalMemory     := globalBuffers + perThreadBuffers
  physicalMemory  := memory.TotalMemory()
  pctOfSysMem     := totalMemory * 100 / physicalMemory

  outInfoF("Current physical memory           : %v\n", humanize.Bytes(physicalMemory))
  outInfoF("Max Memory Ever Allocated         : %v\n", humanize.Bytes(maxMemory))
  outInfoF("Configured Max Per-thread Buffers : %v\n", humanize.Bytes(perThreadBuffers))
  outInfoF("Configured Max Global Buffers     : %v\n", humanize.Bytes(globalBuffers))
  outInfoF("Configured Max Memory Limit       : %v\n", humanize.Bytes(totalMemory))
  outInfoF("Plus %v per temporary table created\n", humanize.Bytes(effectiveTmpTableSize))

  if pctOfSysMem > 90 {
    outWarnTips("Max memory limit exceeds 90% of physical memory")
  } else {
    outInfoTips("Max memory limit seem to be within acceptable norms")
  }

  return gl
}

func (gl *TunerController) joins() *TunerController {
  outHeader("Joins")

  selectFullJoin    := gl.ParseBytes(vrSelectFullJoin)
  selectRangeCheck  := gl.ParseBytes(vrSelectRangeCheck)
  joinBufferSize    := gl.ParseBytes(vrJoinBufferSize) + 4096

  outInfoF("Current join_buffer_size   = %v\n", humanize.Bytes(joinBufferSize))
  outInfoF("Current Select range check = %v\n", humanize.Bytes(joinBufferSize))
  outInfoF("Current Select full join   = %v\n", humanize.Bytes(selectFullJoin))

  outInfoF("You have had %v queries where a join could not use an index properly\n",
    humanize.Bytes(selectFullJoin))

  printError, raiseBuffer := false, false
  if selectFullJoin > 0 {
    printError, raiseBuffer = true, true
  }
  if selectRangeCheck == 0 {
    outInfoTips("Your joins seem to be using indexes properly")
  } else if selectRangeCheck > 0 {
    outWarnTips("You have had %v joins without keys that check for key usage after each row",
      humanize.Bytes(selectRangeCheck))
    printError, raiseBuffer = true, true
  }

  if joinBufferSize >= 4194304 {
    outWarnTips("join_buffer_size >= 4 M. This is not advised.")
    raiseBuffer = false
  }

  if printError {
    outWarnTips("You should enable \"log-queries-not-using-indexes\" Then look for non indexed joins in the slow query log.")

    if raiseBuffer {
      outWarnTips("If you are unable to optimize your queries you may want to increase your join_buffer_size to accommodate larger joins in one pass.")
      outNoticeTips("This script will still suggest raising the join_buffer_size when ANY joins not using indexes are found.")
    }
  }

  return gl
}

func Start(drv models.Driver) {
  outHeader("-- MYSQL PERFORMANCE TUNING --")

  tc := NewTuner(drv)
  defer tc.Close()

  tc.GetGlobals().
    version().
    uptime().
    slow().
    binlog().
    connection().
    threads().
    queryCache().
    innodb().
    memory().
    joins()

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
  vrDataDir
  vrThreadsConnected
  vrThreadsCreated
  vrThreadsCached
  vrThreadCacheSize
  vrThreadStack
  vrSlowQueryLog
  vrLongQueryTime
  vrSlowQueries
  vrBinLog
  vrBinMaxSize
  vrBinExpire
  vrBinSync
  vrBinlogCacheSize
  vrMaxConnections
  vrMaxUsedConnections
  vrKeyReadReq
  vrKeyReads
  vrKeyBlocksUsed
  vrKeyBlocksUnused
  vrKeyCacheBlockSize
  vrKeyBufferSize
  vrQueryCacheSize
  vrQueryCacheLimit
  vrQueryCacheMinRes
  vrQcacheFreeMemory
  vrQcacheTotalBlocks
  vrQcacheFreeBlocks
  vrQcacheLowMemPrunes
  vrInnoDBIgnoreBuiltIn
  vrInnoDBBufferPoolSize
  vrInnoDBBufferPoolPagesFree
  vrInnoDBBufferPoolPagesTotal
  vrInnoDBLogFileSize
  vrInnoDBLogBufferSize
  vrInnoDBFlushLogAtTrxCommit
  vrInnoDBAdditionalMemPoolSize
  vrMaxHeapTableSize
  vrTmpTableSize
  vrReadBufferSize
  vrReadRndBufferSize
  vrSortBufferSize
  vrJoinBufferSize
  vrSelectFullJoin
  vrSelectRangeCheck
)

var (
  global = []string{"STATUS", "VARIABLES"}

  vars   = []string{
    "version",
    "version_compile_machine",
    "Uptime",
    "Questions",
    "datadir",

    "Threads_connected",
    "Threads_created",
    "Threads_cached",
    "thread_cache_size",
    "thread_stack",

    "slow_query_log",
    "long_query_time",
    "Slow_queries",

    "log_bin",
    "max_binlog_size",
    "expire_logs_days",
    "sync_binlog",
    "binlog_cache_size",

    "max_connections",
    "Max_used_connections",

    "Key_read_requests",
    "Key_reads",
    "Key_blocks_used",
    "Key_blocks_unused",
    "key_cache_block_size",
    "key_buffer_size",

    "query_cache_size",
    "query_cache_limit",
    "query_cache_min_res_unit",
    "Qcache_free_memory",
    "Qcache_total_blocks",
    "Qcache_free_blocks",
    "Qcache_lowmem_prunes",

    "ignore_builtin_innodb",
    "innodb_buffer_pool_size",
    "Innodb_buffer_pool_pages_free",
    "Innodb_buffer_pool_pages_total",
    "innodb_log_file_size",
    "innodb_log_buffer_size",
    "innodb_flush_log_at_trx_commit",
    "innodb_additional_mem_pool_size",

    "max_heap_table_size",
    "tmp_table_size",
    "read_buffer_size",
    "read_rnd_buffer_size",
    "sort_buffer_size",

    "join_buffer_size",
    "Select_full_join",
    "Select_range_check",
  }
)
