package models

import (
  "strings"
  "fmt"
)

type Propert struct {
  Table     string
  Field     string
  Type      string
  Size      int
  Decimals  int
  Unique    string
  Required  string
  ColNull   int
}

func (Propert) TableName() string {
  return "propert"
}

func (p *Propert) MyCreateField(primaryKey, relKey string) string {
  return p.GetField() + " " + p.GetType(primaryKey, relKey)
}

func (p *Propert) GetField() string {
  return "`" + strings.ToLower(strings.TrimSpace(p.Field)) + "`"
}

func (p *Propert) GetType(primaryKey, relKey string) string {
  var null string

  if p.Field == primaryKey {
    null = "NOT NULL PRIMARY KEY"
  } else if p.Field == relKey {
    null = "NOT NULL UNIQUE"
  } else if p.ColNull == 0 {
    null = "NOT NULL"
  } else {
    null = "DEFAULT NULL"
  }

  switch p.Type {
  case "I":
    return fmt.Sprintf("INT(%d) %s", p.Size, null)
  case "C":
    return fmt.Sprintf("CHAR(%d) %s", p.Size, null)
  case "N":
    return fmt.Sprintf("DECIMAL(%d,%d) %s", p.Size, p.Decimals, null)
  case "M":
    return fmt.Sprintf("TEXT")
  case "D":
    return fmt.Sprintf("DATE %s", null)
  case "T":
    return fmt.Sprintf("DATETIME %s", null)
  }

  return ""
}
