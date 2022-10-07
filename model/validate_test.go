package model

import (
   json2 "encoding/json"
   "testing"
)

func TestRender(t *testing.T) {
   key := "key1"
   type args struct {
      mutation  string
      variables map[string]interface{}
   }
   tests := []struct {
      name    string
      args    args
      want    string
      wantErr bool
      tags    []TagObject
   }{
      {name: "Render a map of arrays as a tag",
         args: args{
            mutation: `mutation {
            taggingAddTagsToEntity(guid: "", tags: {{{ TAGS }}}) {
            errors {
            message
            type
         }
         }
         }`,
            variables: map[string]interface{}{"TAGS": ""},
         },
         want: `mutation {
            taggingAddTagsToEntity(guid: "", tags: [{"key":"key1","values":["value1","value2"]}]) {
            errors {
            message
            type
         }
         }
         }`,
         wantErr: false,
         tags:    []TagObject{{&key, []string{"value1", "value2"}}},
      },
   }
   for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) {
         json, err := json2.Marshal(tt.tags)
         jsonString := string(json)
         tt.args.variables["TAGS"] = jsonString
         got, err := Render(tt.args.mutation, tt.args.variables)
         if (err != nil) != tt.wantErr {
            t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
            return
         }
         if got != tt.want {
            t.Errorf("Render() got = %v, want %v", got, tt.want)
         }
      })
   }
}
