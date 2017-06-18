package version

import (
  "fmt"
)

func (self Version) ToString() (string) {
  return fmt.Sprintf("v%d.%d.%d", self.Major, self.Minor, self.Patch)
}

type Version struct {
  Major int
  Minor int
  Patch int
}
