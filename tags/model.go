package tags

import (
   log "github.com/sirupsen/logrus"
   "strings"
)

type Model struct {
   Guid *string `json:",omitempty"`
   // EntityGuid      *string           `json:",omitempty"`
   // ListQueryFilter *string           `json:",omitempty"`
   // TagString       *string           `json:",omitempty"`
   Variables map[string]string `json:",omitempty"`
   Tags      []TagObject       `json:",omitempty"`
}

type TagObject struct {
   Key    string   `json:",omitempty"`
   Values []string `json:",omitempty"`
}

func NewTagModel(guid *string, tags map[string]string, vars map[string]string) *Model {
   log.Debugf("NewTagModel: guid: %s tags: %+v vars: %+v", *guid, tags, vars)
   m := &Model{
      Guid:      guid,
      Tags:      make([]TagObject, 0, len(tags)),
      Variables: vars,
   }
   for k, v := range tags {
      to := TagObject{
         Key:    k,
         Values: strings.Split(v, ","),
      }
      m.Tags = append(m.Tags, to)
   }
   log.Debugf("NewTagModel: m.Tags: %+v", m.Tags)
   return m
}
