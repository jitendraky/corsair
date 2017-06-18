package models

import (
  "time"

  "framework/version"
)

type Build struct {
  Date        time.Time
  Development bool
  Version     version.Version
}

type ServerType struct {
  Name       string
  Directives []Directive
}

type Directive struct {
  Name      string
  // 1-100 to be reserved for order critical items
  Weight    string
}

type TLS struct {
  CAUrl     string
  Email     string
}
