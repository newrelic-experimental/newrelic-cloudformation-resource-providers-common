package model

import (
   "fmt"
   "github.com/cbroglie/mustache"
   "github.com/graphql-go/graphql/language/ast"
   "github.com/graphql-go/graphql/language/parser"
   "github.com/graphql-go/graphql/language/source"
   log "github.com/sirupsen/logrus"
)

// TODO Fix logging

func Validate(mutation *string) (err error) {
   // Parse the GetGraphQL so we can test for multiple
   doc, err := parser.Parse(parser.ParseParams{
      Source: &source.Source{
         Body: []byte(*mutation),
         Name: "GetGraphQL",
      },
   })
   if err != nil {
      return
   }

   if len(doc.Definitions) == 0 {
      err = fmt.Errorf("result document contains no definitions")
      return
   }
   if len(doc.Definitions) > 1 {
      err = fmt.Errorf("result document contains multiple definitions")
      return
   }

   // At this point we know we only have one Definition, ensure it's an OperationDefinition
   def := doc.Definitions[0]
   kind := def.GetKind()
   if kind != "OperationDefinition" {
      err = fmt.Errorf("expected OperationDefintion, received %s", kind)
      return
   }

   // Cast the ast.Node to a *ast.OperationDefinition
   opDef := doc.Definitions[0].(*ast.OperationDefinition)
   // Ensure only one operation is present
   if len(opDef.SelectionSet.Selections) > 1 {
      err = fmt.Errorf("%d operations found, only 1 is allowed ", len(opDef.SelectionSet.Selections))
      return
   }

   // Ensure each Operation is valid
   for _, s := range opDef.SelectionSet.Selections {
      op := s.(*ast.Field)
      log.Debugf("Validate: field name: %s field: %+v", op.Name.Value, op.SelectionSet)

      // All Operations must be aliased so we can match them to their delete
      // if op.Alias == nil {
      //    log.Warnf("unaliased operation, will result in orahan resource if rollback required. %s", printer.Print(op))
      // }

      // valid := false
      // for _, r := range s.GetSelectionSet().Selections {
      //    selection := r.(*ast.Field)
      //    log.Debugf("Validate: field name: %s field: %+v", selection.Name.Value, selection.SelectionSet)
      //    if selection.Name.Value == "id" || selection.Name.Value == "guid" {
      //       valid = true
      //    }
      // }

      // FIXME this is not valid for Tags
      //      if FindFieldInSelectionSet(m.GuidField, s.GetSelectionSet()) {
      //         err = fmt.Errorf("every mutation must return either an id or guid")
      //      }
   }
   return
}

func FindFieldInSelectionSet(field string, set *ast.SelectionSet) (found bool) {
   if set == nil {
      return
   }
   for _, r := range set.Selections {
      selection := r.(*ast.Field)
      log.Debugf("Validate: field name: %s field: %+v", selection.Name.Value, selection.SelectionSet)
      if selection.Name.Value == field {
         found = true
         break
      } else {
         found = FindFieldInSelectionSet(field, selection.SelectionSet)
      }
   }
   return
}

func Render(mutation string, variables interface{}) (string, error) {
   return mustache.Render(mutation, variables)
}
