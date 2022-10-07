package configuration

import (
   "encoding/json"
   "github.com/aws/aws-sdk-go/aws/session"
   "github.com/aws/aws-sdk-go/service/cloudformation"
   "github.com/newrelic-experimental/newrelic-cloudformation-resource-providers-common/logging"
   log "github.com/sirupsen/logrus"
   "os"
   "strings"
)

var usEndpoint = "https://api.newrelic.com/graphql"

// TODO try reading typeName from .rpdk-configuration
var resourceType = cloudformation.RegistryTypeResource
var mockConfig = `{  "APIKey": "mockapikey",  "AccountID": "987654321",  "Endpoint": "https://api.newrelic.com/snafu"}`

type Config struct {
   Endpoint  *string `json:",omitempty"`
   AccountID *string `json:",omitempty"`
   APIKey    *string `json:",omitempty"`
}

func (c *Config) GetEndpoint() string {
   if c.Endpoint == nil || *c.Endpoint == "" {
      return usEndpoint
   }
   return *c.Endpoint
}

func NewConfiguration(s *session.Session, typeName *string) (c *Config) {
   // 1. If we find a TypeConfiguration use it and return
   c = &Config{}
   if c.configurationFromCloudFormation(s, typeName) {
      log.Debugf("SetConfiguration: using CloudFormation")
      return
   }
   // 2. If we find a TypeConfiguration envvar AND the file exists use it and return
   if c.configurationFromFile() {
      log.Debugf("SetConfiguration: using file")
      return
   }
   // 3. If we find nothing use the mock type configuration
   log.Debugf("SetConfiguration: using mock")
   c.setConfiguration(&mockConfig)
   return
}

func (c *Config) setConfiguration(jsonConfig *string) {
   tc := Config{}
   err := json.Unmarshal([]byte(*jsonConfig), &tc)
   if err != nil {
      panic("error unmarshalling typeconfiguration: " + err.Error())
   }

   if tc.APIKey == nil {
      panic("nil APIKey, typeOutput: " + *jsonConfig)
   } else {
      c.APIKey = tc.APIKey
   }

   if tc.Endpoint == nil {
      log.Warnf("no configured Endpoint, using US default")
   } else {
      c.Endpoint = tc.Endpoint
   }

   if tc.AccountID == nil {
      panic("nil AccountID, typeOutput: " + *jsonConfig)
   } else {
      c.AccountID = tc.AccountID
   }

   if strings.Contains(strings.ToLower(*c.APIKey), "mockapikey") {
      err := os.Setenv("Mock", "true")
      if err != nil {
         log.Errorf("error setting Mock envvar: %v", err)
      } else {
         log.Traceln(os.Environ())
      }
   }
}

func (c *Config) configurationFromFile() bool {
   filename := os.Getenv("TypeConfigurationFile")
   if filename == "" {
      return false
   }
   config, err := os.ReadFile(filename)
   if err != nil {
      path, _ := os.Getwd()
      log.Warnf("TypeConfigurationFile %s not found, using mock. PWD: %s Error: %v", filename, path, err)
      return false
   }
   jsonConfig := string(config)
   c.setConfiguration(&jsonConfig)
   return true
}

func (c *Config) configurationFromCloudFormation(s *session.Session, typeName *string) bool {
   // Create a CloudFormation client from just a session.
   svc := cloudformation.New(s)
   bdtci := cloudformation.BatchDescribeTypeConfigurationsInput{}
   tcia := []*cloudformation.TypeConfigurationIdentifier{
      {
         Type:     &resourceType,
         TypeName: typeName,
      },
   }
   bdtci.SetTypeConfigurationIdentifiers(tcia)
   o, err := svc.BatchDescribeTypeConfigurations(&bdtci)
   logging.Dump(log.DebugLevel, o, "SetConfiguration: o")

   if err != nil {
      panic("error getting typeconfiguration: " + err.Error())
   }

   var jsonConfig *string
   for i, tc := range o.TypeConfigurations {
      log.Debugf("SetConfiguration: TypeConfiguration index: %d", i)
      jsonConfig = tc.Configuration
   }
   if jsonConfig == nil {
      return false
   }
   c.setConfiguration(jsonConfig)
   return true
}
