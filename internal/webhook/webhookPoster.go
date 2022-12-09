package webhook

import (
)

type Poster interface {
  SendPost(falcoEventStr string, mainEventName string, eventNumber int)      bool
  InitConfig() 
}

func InitConfig() {
  
}
