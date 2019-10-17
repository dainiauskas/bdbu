package app

import (
  "fmt"
)

const (
  myFoundRows = 2

  myBigPacket = 1 << (2 + iota)
  myNoPrompt
  myDynamicCursor
  myNoSchema
  myNoDefaultCursor
  myNoLocale
  myPadSpace
  myFullColumnNames
  myCompressedProto
  myIgnoreSpace
  myNamedPipe
  myNoBigint
  myNoCatalog
  myUseMyCnf
  mySafe
  myNoTransactions
  myLogQuery
  myNoCache
  myForwardCursor
  myAutoreconnect
  myAutoIsNull
  myZeroDateToMin
  myMinDateToZero
  myMultiStatements
  myColumnSizeS32
  myNoBinaryResult
  myDfltBigintBindStr
  myNoInformationSchema
  myNoDateOverflow = 0
)

var (
  showTemplate = ""

  options = map[string]int{
      "NO_DATE_OVERFLOW":       myNoDateOverflow,
      "FOUND_ROWS":             myFoundRows,
      "BIG_PACKETS":            myBigPacket,
      "NO_PROMPT":              myNoPrompt,
      "DYNAMIC_CURSOR":         myDynamicCursor,
      "NO_SCHEMA":              myNoSchema,
      "NO_DEFAULT_CURSOR":      myNoDefaultCursor,
      "NO_LOCALE":              myNoLocale,
      "PAD_SPACE":              myPadSpace,
      "FULL_COLUMN_NAMES":      myFullColumnNames,
      "COMPRESSED_PROTO":       myCompressedProto,
      "IGNORE_SPACE":           myIgnoreSpace,
      "NAMED_PIPE":             myNamedPipe,
      "NO_BIGINT":              myNoBigint,
      "NO_CATALOG":             myNoCatalog,
      "USE_MYCNF":              myUseMyCnf,
      "SAFE":                   mySafe,
      "NO_TRANSACTIONS":        myNoTransactions,
      "LOG_QUERY":              myLogQuery,
      "NO_CACHE":               myNoCache,
      "FORWARD_CURSOR":         myForwardCursor,
      "AUTO_RECONNECT":         myAutoreconnect,
      "AUTO_IS_NULL":           myAutoIsNull,
      "ZERO_DATE_TO_MIN":       myZeroDateToMin,
      "MIN_DATE_TO_ZERO":       myMinDateToZero,
      "MULTI_STATEMENTS":       myMultiStatements,
      "COLUMN_SIZE_S32":        myColumnSizeS32,
      "NO_BINARY_RESULT":       myNoBinaryResult,
      "DFLT_BIGINT_BIND_STR":   myDfltBigintBindStr,
      "NO_INFORMATION_SCHEMA":  myNoInformationSchema,
  }
)

func init() {
  var maxLength = 0

  for k, _ := range options {
    if maxLength < len(k) {
      maxLength = len(k)
    }
  }

  showTemplate = "%" + fmt.Sprintf("%v", maxLength) + "s = %v\n"
}

func MyOptionExplain(md bool) {
  var option = 5260042

  for k, val := range options {
    status := (option & val) != 0
    if !md || status {
      fmt.Printf(showTemplate, k, status)
    }
  }

}
