package sqlbuild

import (
	"slices"
	"strconv"
)

// UpdateBuilder update语句构造器
type UpdateBuilder interface {
	Set(column string, value any) UpdateBuilder
	SetRaw(column, value string) UpdateBuilder
	SetByCondition(condition bool, column string, value any) UpdateBuilder
	SetBySlice(columns []string, values []any) UpdateBuilder
	SetByMap(kv map[string]any) UpdateBuilder
	Returning(fields ...string) UpdateBuilder
	SQLBuilder
	Condition
}

type PostgresUpdateBuilder struct {
	tableName string
	sets      []string
	returns   []string
	PostgresStatementParameterBuffer
	PostgresCondition
}

func (self *PostgresUpdateBuilder) table(tableName string) SQLBuilder {
	self.tableName = tableName
	return self
}

func (self *PostgresUpdateBuilder) SQL() string {
	builder := defaultPool.GetStringBuilder()
	builder.WriteString("UPDATE " + self.tableName)
	if self.sets != nil && len(self.sets) > 0 {
		builder.WriteString(" SET ")
		handleStringsSplice(self.sets, ", ", builder)
	}
	if self.wheres != nil && len(self.wheres) > 0 {
		builder.WriteString(" WHERE ")
		handleStringsSplice(self.wheres, " AND ", builder)
	}
	if self.returns != nil && len(self.returns) > 0 {
		builder.WriteString(" RETURNING ")
		handleStringsSplice(self.returns, ", ", builder)
	}
	defer defaultPool.RecycleStringBuilder(builder)
	return builder.String()
}

func (self *PostgresUpdateBuilder) Set(column string, value any) UpdateBuilder {
	self.addSet(column, value)
	return self
}

func (self *PostgresUpdateBuilder) SetRaw(column, value string) UpdateBuilder {
	if self.sets == nil {
		self.sets = make([]string, 0)
	}
	self.sets = append(self.sets, column+" = "+value)
	return self
}

func (self *PostgresUpdateBuilder) SetByCondition(condition bool, column string, value any) UpdateBuilder {
	if !condition {
		return self
	}
	self.addSet(column, value)
	return self
}

func (self *PostgresUpdateBuilder) SetBySlice(columns []string, values []any) UpdateBuilder {
	if columns == nil || len(columns) == 0 {
		return self
	}
	for i := 0; i < len(columns); i++ {
		self.addSet(columns[i], values[i])
	}
	return self
}

func (self *PostgresUpdateBuilder) SetByMap(kv map[string]any) UpdateBuilder {
	if kv == nil || len(kv) == 0 {
		return self
	}
	for k, v := range kv {
		self.addSet(k, v)
	}
	return self
}

func (self *PostgresUpdateBuilder) Returning(fields ...string) UpdateBuilder {
	if self.returns == nil {
		self.returns = make([]string, 0)
	}
	self.returns = slices.Concat(self.returns, fields)
	return self
}

func (self *PostgresUpdateBuilder) addSet(column string, value any) {
	if self.sets == nil {
		self.sets = make([]string, 0)
	}
	index := self.addParameter(value)
	self.sets = append(self.sets, column+" = $"+strconv.Itoa(index))
}
