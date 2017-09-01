package tspec_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/wy-z/tspec/samples"
	"github.com/wy-z/tspec/tspec"
)

var samplesPackage *ast.Package

func TestTSpec(t *testing.T) {
	suite.Run(t, new(TSpecTestSuite))
}

type TSpecTestSuite struct {
	suite.Suite
	parser *tspec.Parser
	pkg    *ast.Package
}

func (s *TSpecTestSuite) SetupTest() {
	s.parser = tspec.NewParser()
	s.pkg = samplesPackage
}

func (s *TSpecTestSuite) testParse(typeStr, assert string) {
	require := s.Require()

	schema, err := s.parser.Parse(s.pkg, typeStr)
	require.NoError(err)
	require.NotNil(schema)

	defs := s.parser.Definitions()
	bts, err := json.MarshalIndent(defs, "", "\t")
	require.NoError(err)
	require.Equal(string(bytes.TrimSpace(samples.MustAsset(assert))),
		string(bytes.TrimSpace(bts)))
	s.parser.Reset()
}

func (s *TSpecTestSuite) TestParse() {
	s.testParse("BasicTypes", "samples/source/basic_types.json")
	s.testParse("NormalStruct", "samples/source/normal_struct.json")
	s.testParse("StructWithNoExportField", "samples/source/struct_with_no_export_field.json")
	s.testParse("StructWithAnonymousField", "samples/source/struct_with_anonymous_field.json")
	s.testParse("StructWithCircularReference", "samples/source/struct_with_circular_reference.json")
	s.testParse("StructWithInheritance", "samples/source/struct_with_inheritance.json")
	s.testParse("MapType", "samples/source/map_type.json")
	s.testParse("ArrayType", "samples/source/array_type.json")
}

func (s *TSpecTestSuite) TestParseInvalidMap() {
	require := s.Require()

	schema, err := s.parser.Parse(s.pkg, "InvalidMap")
	require.Error(err)
	require.Nil(schema)
}

func init() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "tspec_test.go", `
package tspec_test
import (
    "github.com/wy-z/tspec/samples"
)
`, parser.ImportsOnly)
	pkg, err := tspec.NewParser().Import(f.Imports[0])
	if err != nil {
		msg := fmt.Sprintf("failed to import 'github.com/wy-z/tspec/samples': %s", err)
		panic(msg)
	}
	samplesPackage = pkg
}
