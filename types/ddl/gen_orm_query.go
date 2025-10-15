package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var path2Table2FieldMapping = map[string]map[string]map[string]any{
	"apps/user/domain/internal/dal/query": {
		"user": {},
	},
	"apps/task/domain/internal/dal/query": {
		"task": {},
	},
}

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("gorm.Open failed, err=%v", err)
	}

	rootPath, err := findProjectRoot()
	if err != nil {
		log.Fatalf("failed to find project room: %v", err)
	}

	for path, mapping := range path2Table2FieldMapping {

		g := gen.NewGenerator(gen.Config{
			OutPath: filepath.Join(rootPath, path),
			Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		})

		parts := strings.Split(path, "/")
		modelPath := strings.Join(append(parts[:len(parts)-1], g.Config.ModelPkgPath), "/")

		g.UseDB(db)
		g.WithOpts(gen.FieldType("deleted_at", "gorm.DeletedAt"))

		var resolveType func(typ reflect.Type, required bool) string
		resolveType = func(typ reflect.Type, required bool) string {
			switch typ.Kind() {
			case reflect.Ptr:
				return resolveType(typ.Elem(), false)
			case reflect.Slice:
				return "[]" + resolveType(typ.Elem(), required)
			default:
				prefix := "*"
				if required {
					prefix = ""
				}

				if strings.HasSuffix(typ.PkgPath(), modelPath) {
					return prefix + typ.Name()
				}

				return prefix + typ.String()
			}
		}

		genModify := func(col string, model any) func(f gen.Field) gen.Field {
			return func(f gen.Field) gen.Field {
				if f.ColumnName != col {
					return f
				}

				st := reflect.TypeOf(model)
				// f.Name = st.Name()
				f.Type = resolveType(st, true)
				f.GORMTag.Set("serializer", "json")
				return f
			}
		}

		timeModify := func(f gen.Field) gen.Field {
			if f.ColumnName == "updated_at" {
				// https://gorm.io/zh_CN/docs/models.html#%E5%88%9B%E5%BB%BA-x2F-%E6%9B%B4%E6%96%B0%E6%97%B6%E9%97%B4%E8%BF%BD%E8%B8%AA%EF%BC%88%E7%BA%B3%E7%A7%92%E3%80%81%E6%AF%AB%E7%A7%92%E3%80%81%E7%A7%92%E3%80%81Time%EF%BC%89
				f.GORMTag.Set("autoUpdateTime", "milli")
			}
			if f.ColumnName == "created_at" {
				f.GORMTag.Set("autoCreateTime", "milli")
			}
			return f
		}

		var models []any
		for table, col2Model := range mapping {
			opts := make([]gen.ModelOpt, 0, len(col2Model))
			for column, m := range col2Model {
				cp := m
				opts = append(opts, gen.FieldModify(genModify(column, cp)))
			}
			opts = append(opts, gen.FieldModify(timeModify))
			models = append(models, g.GenerateModel(table, opts...))
		}

		g.ApplyBasic(models...)

		g.Execute()
	}
}

func findProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	projectDir := filepath.Dir(filepath.Dir(filepath.Dir(filename))) // notice: the relative path of the script file is assumed here

	if _, err := os.Stat(filepath.Join(projectDir, "apps")); os.IsNotExist(err) {
		return "", fmt.Errorf("could not find 'domain' directory in backend path: %s", projectDir)
	}

	return projectDir, nil
}
