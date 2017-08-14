package tspec

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"sync"

	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

// Parser ...
type Parser struct {
	fset       *token.FileSet
	fileMap    map[string]*ast.File
	dirPkgMap  map[string]*ast.Package
	pkgObjsMap map[*ast.Package]map[string]*ast.Object
	typeMap    map[string]*spec.Schema
}

// NewParser ...
func NewParser() (parser *Parser) {
	parser = new(Parser)
	parser.fset = token.NewFileSet()
	parser.fileMap = make(map[string]*ast.File)
	parser.dirPkgMap = make(map[string]*ast.Package)
	parser.pkgObjsMap = make(map[*ast.Package]map[string]*ast.Object)
	parser.typeMap = make(map[string]*spec.Schema)
	return
}

// ParseFile ...
func (t *Parser) ParseFile(filePath string) (f *ast.File, err error) {
	if tmpF, ok := t.fileMap[filePath]; ok {
		f = tmpF
		return
	}
	f, err = parser.ParseFile(t.fset, filePath, nil, parser.ParseComments)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	t.fileMap[filePath] = f
	return
}

// ParseDir ...
func (t *Parser) ParseDir(dirPath string, pkgName string) (pkg *ast.Package, err error) {
	if tmpPkg, ok := t.dirPkgMap[dirPath]; ok {
		pkg = tmpPkg
		return
	}

	pkgs, err := parser.ParseDir(t.fset, dirPath, nil, parser.ParseComments)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	for k := range pkgs {
		if k == pkgName {
			pkg = pkgs[k]
			break
		}
	}
	if pkg == nil {
		err = errors.Errorf("%s not found in %s", pkgName, dirPath)
		return
	}

	t.dirPkgMap[dirPath] = pkg
	return
}

// Import ...
func (t *Parser) Import(ispec *ast.ImportSpec) (pkg *ast.Package, err error) {
	pkgPath := strings.Trim(ispec.Path.Value, "\"")

	wd, err := os.Getwd()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	importPkg, err := build.Import(pkgPath, wd, build.ImportComment)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	pkg, err = t.ParseDir(importPkg.Dir, importPkg.Name)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// ParsePkg ...
func (t *Parser) ParsePkg(pkg *ast.Package) (objs map[string]*ast.Object, err error) {
	if tmpObjs, ok := t.pkgObjsMap[pkg]; ok {
		objs = tmpObjs
		return
	}

	objs = make(map[string]*ast.Object)
	for _, f := range pkg.Files {
		for key, obj := range f.Scope.Objects {
			if obj.Kind == ast.Typ {
				objs[key] = obj
			}
		}
	}

	t.pkgObjsMap[pkg] = objs
	return
}

func (t *Parser) parseTypeStr(oPkg *ast.Package, typeStr string) (pkg *ast.Package,
	obj *ast.Object, err error) {
	var objs map[string]*ast.Object
	var pkgName, typeTitle string
	var ok bool

	strs := strings.Split(typeStr, ".")
	l := len(strs)
	if l == 0 || l > 2 {
		err = errors.Errorf("invalid type str %s", typeStr)
		return
	}
	if l == 1 {
		pkgName = oPkg.Name
		typeTitle = strs[0]
	} else {
		pkgName = strs[0]
		typeTitle = strs[1]
	}
	if pkgName == oPkg.Name {
		pkg = oPkg
		objs, err = t.ParsePkg(pkg)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		if obj, ok = objs[typeTitle]; !ok {
			err = errors.Errorf("%s not found in package %s", typeTitle, pkg.Name)
			return
		}
		return
	}
	var p *ast.Package
	for _, file := range oPkg.Files {
		for _, ispec := range file.Imports {
			p, err = t.Import(ispec)
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			if !(pkgName == p.Name || (ispec.Name != nil && pkgName == ispec.Name.Name)) {
				continue
			}
			objs, err = t.ParsePkg(p)
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			if _, ok = objs[typeTitle]; !ok {
				continue
			} else {
				pkg = p
				obj = objs[typeTitle]
				break
			}
		}
	}
	if pkg == nil || obj == nil {
		err = errors.Errorf("%s.%s not found", pkgName, typeTitle)
		return
	}

	return
}

func (t *Parser) parseIdentExpr(oExpr ast.Expr, pkg *ast.Package) (expr ast.Expr, err error) {
	expr = starExprX(oExpr)
	if ident, ok := expr.(*ast.Ident); ok {
		if ident.Obj != nil {
			ts, e := objDeclTypeSpec(ident.Obj)
			if e != nil {
				err = errors.WithStack(e)
				return
			}
			expr = starExprX(ts.Type)
		} else {
			// try to find obj in pkg
			objs, e := t.ParsePkg(pkg)
			if e != nil {
				err = errors.WithStack(e)
				return
			}
			if obj, ok := objs[ident.Name]; ok {
				ts, e := objDeclTypeSpec(obj)
				if e != nil {
					err = errors.WithStack(e)
					return
				}
				expr = starExprX(ts.Type)
			}
		}
	}
	return
}

func (t *Parser) parseTypeRef(pkg *ast.Package, expr ast.Expr, typeTitle, typeID string) (
	schema *spec.Schema, err error) {
	ident, isIdent := starExprX(expr).(*ast.Ident)
	typeExpr, err := t.parseIdentExpr(expr, pkg)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	switch typ := typeExpr.(type) {
	case *ast.StructType:
		if isIdent {
			typeID = pkg.Name + "." + ident.Name
			typeTitle = ident.Name
		}
		schema = spec.RefProperty("#/definitions/" + typeID)
		_, err = t.parseType(pkg, typ, typeTitle, typeID)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		return
	case *ast.SelectorExpr:
		typeStr, e := selectorExprTypeStr(typ)
		if e != nil {
			err = errors.WithStack(err)
			return
		}
		if typeStr != "time.Time" {
			typeID := typeStr
			schema = spec.RefProperty("#/definitions/" + typeID)
			_, err = t.Parse(pkg, typeStr)
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			return
		}
	}
	return t.parseType(pkg, typeExpr, typeTitle, typeID)
}

var parseTypeLock sync.Mutex

func (t *Parser) parseType(pkg *ast.Package, expr ast.Expr, typeTitle, typeID string) (schema *spec.Schema,
	err error) {
	parseTypeLock.Lock()
	if tmpSchema, ok := t.typeMap[typeID]; ok {
		schema = tmpSchema
		parseTypeLock.Unlock()
		return
	}
	if typeID != "" {
		t.typeMap[typeID] = nil
	}
	parseTypeLock.Unlock()

	// parse ident expr
	expr, err = t.parseIdentExpr(expr, pkg)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	schema = new(spec.Schema)
	schema.WithID(typeID)
	schema.WithTitle(typeTitle)
	switch typ := expr.(type) {
	case *ast.StructType:
		schema.Typed("object", "")
		if typ.Fields.List == nil {
			break
		}
		for _, field := range typ.Fields.List {
			if len(field.Names) != 0 {
				fieldName := field.Names[0].Name
				var fTypeID, fTypeTitle string
				if _, isAnonymousStruct := starExprX(field.Type).(*ast.StructType); isAnonymousStruct {
					fTypeTitle = typeTitle + "_" + fieldName
					fTypeID = pkg.Name + "." + fTypeTitle
				}
				prop, e := t.parseTypeRef(pkg, field.Type, fTypeTitle, fTypeID)
				if e != nil {
					err = errors.WithStack(e)
					return
				}
				schema.SetProperty(fieldName, *prop)
			} else {
				// inherited struct
				var fieldTypeTitle, fieldTypeID string
				ident, isIdent := starExprX(field.Type).(*ast.Ident)
				fieldExpr, e := t.parseIdentExpr(field.Type, pkg)
				if e != nil {
					err = errors.WithStack(e)
					return
				}
				switch fieldTyp := fieldExpr.(type) {
				case *ast.StructType:
					if isIdent {
						fieldTypeID = pkg.Name + "." + ident.Name
						fieldTypeTitle = ident.Name
					}
				case *ast.SelectorExpr:
					typeStr, e := selectorExprTypeStr(fieldTyp)
					if e != nil {
						err = errors.WithStack(err)
						return
					}
					if typeStr != "time.Time" {
						fieldTypeID = typeStr
						fieldTypeTitle = typeStr
					}
				}
				inheritedSchema, e := t.parseType(pkg, field.Type, fieldTypeTitle,
					fieldTypeID)
				if e != nil {
					err = errors.WithStack(e)
					return
				}
				for propKey, propSchema := range inheritedSchema.Properties {
					if _, ok := schema.Properties[propKey]; !ok {
						schema.SetProperty(propKey, propSchema)
					}
				}
			}
		}
	case *ast.ArrayType:
		itemsSchema, e := t.parseTypeRef(pkg, typ.Elt, "", "")
		if e != nil {
			err = errors.WithStack(e)
			return
		}
		schema = spec.ArrayProperty(itemsSchema)
	case *ast.MapType:
		valueSchema, e := t.parseTypeRef(pkg, typ.Value, "", "")
		if e != nil {
			err = errors.WithStack(e)
			return
		}
		schema = spec.MapProperty(valueSchema)
	case *ast.SelectorExpr:
		typeStr, e := selectorExprTypeStr(typ)
		if e != nil {
			err = errors.WithStack(err)
			return
		}
		if typeStr == "time.Time" {
			typeType, typeFormat, e := parseBasicType("time")
			if e != nil {
				err = errors.WithStack(e)
				return
			}
			schema.Typed(typeType, typeFormat)
		} else {
			schema, err = t.Parse(pkg, typeStr)
			if err != nil {
				err = errors.WithStack(err)
				return
			}
		}
	case *ast.Ident: // basic type only
		typeType, typeFormat, e := parseBasicType(typ.Name)
		if e != nil {
			err = errors.WithStack(e)
			return
		}
		schema.Typed(typeType, typeFormat)
	case *ast.InterfaceType:
	default:
		err = errors.Errorf("invalid expr type %T", typ)
		return
	}

	if typeID != "" {
		t.typeMap[typeID] = schema
	}
	return
}

// Parse ...
func (t *Parser) Parse(oPkg *ast.Package, typeStr string) (
	schema *spec.Schema, err error) {
	pkg, obj, err := t.parseTypeStr(oPkg, typeStr)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	ts, err := objDeclTypeSpec(obj)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	schema, err = t.parseType(pkg, ts.Type, obj.Name, pkg.Name+"."+obj.Name)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// Definitions ...
func (t *Parser) Definitions() (defs spec.Definitions) {
	defs = make(spec.Definitions)
	for k, v := range t.typeMap {
		defs[k] = *v
	}
	return
}

func starExprX(expr ast.Expr) ast.Expr {
	if star, ok := expr.(*ast.StarExpr); ok {
		return star.X
	}
	return expr
}

func objDeclTypeSpec(obj *ast.Object) (ts *ast.TypeSpec, err error) {
	ts, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		err = errors.Errorf("invalid object decl, want *ast.TypeSpec, got %T", ts)
		return
	}
	return
}

func selectorExprTypeStr(expr *ast.SelectorExpr) (typeStr string, err error) {
	xIdent, ok := expr.X.(*ast.Ident)
	if !ok {
		err = errors.Errorf("invalid selector expr %#v", expr)
		return
	}
	typeStr = xIdent.Name + "." + expr.Sel.Name
	return
}

var basicTypes = map[string]string{
	"bool": "boolean:",
	"uint": "integer:uint", "uint8": "integer:uint8", "uint16": "integer:uint16",
	"uint32": "integer:uint32", "uint64": "integer:uint64",
	"int": "integer:int", "int8": "integer:int8", "int16": "integer:int16",
	"int32": "integer:int32", "int64": "integer:int64",
	"uintptr": "integer:int64",
	"float32": "number:float32", "float64": "number:float64",
	"string":    "string",
	"complex64": "number:float", "complex128": "number:double",
	"byte": "string:byte", "rune": "string:byte", "time": "string:date-time",
}

func parseBasicType(typeTitle string) (typ, format string, err error) {
	typeStr, ok := basicTypes[typeTitle]
	if !ok {
		err = errors.Errorf("invalid ident %s", typeTitle)
		return
	}
	exprs := strings.Split(typeStr, ":")
	typ = exprs[0]
	if len(exprs) > 1 {
		format = exprs[1]
	}
	return
}
