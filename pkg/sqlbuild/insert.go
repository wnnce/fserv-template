package sqlbuild

import (
	"slices"
	"strconv"
)

// InsertBuilder insert sql语句构造器
type InsertBuilder interface {
	Insert(column string, value any) InsertBuilder
	InsertRaw(column, value string) InsertBuilder
	InsertByCondition(condition bool, column string, value any) InsertBuilder
	InsertBySlice(columns []string, values []any) InsertBuilder
	InsertByMap(kv map[string]any) InsertBuilder
	Fields(columns ...string) InsertBuilder
	Values(values ...any) InsertBuilder
	Returning(columns ...string) InsertBuilder
	SQLBuilder
}

type PostgresInsertBuilder struct {
	tableName    string
	rawColumnMap map[string]string
	columns      []string
	values       []any
	returns      []string
}

func (self *PostgresInsertBuilder) table(tableName string) SQLBuilder {
	self.tableName = tableName
	return self
}

func (self *PostgresInsertBuilder) SQL() string {
	builder, columnBuilder, valueBuilder := defaultPool.GetStringBuilder(), defaultPool.GetStringBuilder(), defaultPool.GetStringBuilder()
	valueIndex := 1
	builder.WriteString("INSERT INTO " + self.tableName + " ")
	columnBuilder.WriteByte('(')
	valueBuilder.WriteByte('(')
	if self.columns != nil && len(self.columns) > 0 {
		for i, column := range self.columns {
			if i > 0 {
				columnBuilder.WriteByte(',')
				valueBuilder.WriteByte(',')
			}
			columnBuilder.WriteString(column)
			raw, ok := self.rawColumnMap[column]
			if ok {
				valueBuilder.WriteString(raw)
			} else {
				valueBuilder.WriteString("$" + strconv.Itoa(valueIndex))
				valueIndex++
			}
		}
	}
	columnBuilder.WriteByte(')')
	valueBuilder.WriteByte(')')
	builder.WriteString(columnBuilder.String() + " VALUES " + valueBuilder.String())
	defaultPool.RecycleStringBuilder(columnBuilder)
	defaultPool.RecycleStringBuilder(valueBuilder)
	if self.returns != nil && len(self.returns) > 0 {
		builder.WriteString(" RETURNING ")
		handleStringsSplice(self.returns, ", ", builder)
	}
	defer defaultPool.RecycleStringBuilder(builder)
	return builder.String()
}

func (self *PostgresInsertBuilder) Insert(column string, value any) InsertBuilder {
	return self.InsertByCondition(true, column, value)
}

func (self *PostgresInsertBuilder) InsertRaw(column, value string) InsertBuilder {
	if self.rawColumnMap == nil {
		self.rawColumnMap = make(map[string]string)
	}
	self.rawColumnMap[column] = value
	self.columns = append(self.columns, column)
	return self
}

func (self *PostgresInsertBuilder) InsertByCondition(condition bool, column string, value any) InsertBuilder {
	if !condition {
		return self
	}
	if self.columns == nil {
		self.columns = make([]string, 0)
	}
	self.addParameter(value)
	self.columns = append(self.columns, column)
	return self
}

func (self *PostgresInsertBuilder) InsertBySlice(columns []string, values []any) InsertBuilder {
	if columns == nil || len(columns) == 0 {
		return self
	}
	for i := 0; i < len(columns); i++ {
		self.InsertByCondition(true, columns[i], values[i])
	}
	return self
}

func (self *PostgresInsertBuilder) InsertByMap(kv map[string]any) InsertBuilder {
	if kv == nil || len(kv) == 0 {
		return self
	}
	for column, value := range kv {
		self.InsertByCondition(true, column, value)
	}
	return self
}

func (self *PostgresInsertBuilder) Fields(columns ...string) InsertBuilder {
	self.columns = slices.Concat(self.columns, columns)
	return self
}

func (self *PostgresInsertBuilder) Values(values ...any) InsertBuilder {
	self.values = slices.Concat(self.values, values)
	return self
}

func (self *PostgresInsertBuilder) Returning(columns ...string) InsertBuilder {
	if self.returns == nil {
		self.returns = make([]string, 0)
	}
	self.returns = slices.Concat(self.returns, columns)
	return self
}

func (self *PostgresInsertBuilder) Args() []any {
	return self.values
}

func (self *PostgresInsertBuilder) addParameter(param any) int {
	if self.values == nil {
		self.values = make([]any, 0)
	}
	self.values = append(self.values, param)
	return len(self.values)
}

func BatchInsertBuilder[T any](table string, list []T, rowHandler func(T) []any, columns ...string) (string, []any) {
	if list == nil || len(list) == 0 {
		return "", nil
	}
	if len(list) == 1 {
		builder := NewInsertBuilder(table).
			Fields(columns...).
			Values(rowHandler(list[0])...)
		return builder.SQL(), builder.Args()
	}
	args := make([]any, 0, len(list)*len(rowHandler(list[0])))
	builder := defaultPool.GetStringBuilder()
	defer defaultPool.RecycleStringBuilder(builder)
	builder.WriteString("INSERT INTO " + table)
	builder.WriteString(" (")
	handleStringsSplice(columns, ", ", builder)
	builder.WriteString(") VALUES ")
	for i, row := range list {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteByte('(')
		for y, value := range rowHandler(row) {
			if y > 0 {
				builder.WriteByte(',')
			}
			args = append(args, value)
			builder.WriteString("$" + strconv.Itoa(len(args)))
		}
		builder.WriteByte(')')
	}
	return builder.String(), args
}
