package sqlbuild

import (
	"slices"
	"strconv"
)

// SelectBuilder Select Sql构造器
type SelectBuilder interface {
	Select(fields ...string) SelectBuilder
	CountField(field string) SelectBuilder
	Join(table string) *JoinBuilder
	LeftJoin(table string) *JoinBuilder
	RightJoin(table string) *JoinBuilder
	InnerJoin(table string) *JoinBuilder
	OrderBy(orderFields ...string) SelectBuilder
	OrderByAsc(fields ...string) SelectBuilder
	OrderByDesc(fields ...string) SelectBuilder
	GroupBy(fields ...string) SelectBuilder
	Limit(limit int64) SelectBuilder
	Offset(offset int64) SelectBuilder
	CountSQL() string
	addJoin(join string)
	Condition
	SQLBuilder
}

type PostgresSelectBuilder struct {
	tableName    string
	countField   string
	columns      []string
	orderColumns []string
	groupColumns []string
	joins        []string
	offset       int64
	limit        int64
	PostgresStatementParameterBuffer
	PostgresCondition
}

func (self *PostgresSelectBuilder) table(tableName string) SQLBuilder {
	self.tableName = tableName
	return self
}

func (self *PostgresSelectBuilder) SQL() string {
	builder := defaultPool.GetStringBuilder()
	builder.WriteString("SELECT ")
	if self.columns == nil || len(self.columns) == 0 {
		builder.WriteString("*")
	} else {
		handleStringsSplice(self.columns, ", ", builder)
	}
	builder.WriteString(" FROM " + self.tableName)
	if self.joins != nil && len(self.joins) > 0 {
		builder.WriteByte(' ')
		handleStringsSplice(self.joins, " ", builder)
	}
	if self.wheres != nil && len(self.wheres) > 0 {
		builder.WriteString(" WHERE ")
		handleStringsSplice(self.wheres, " AND ", builder)
	}
	if self.groupColumns != nil && len(self.groupColumns) > 0 {
		builder.WriteString(" GROUP BY ")
		handleStringsSplice(self.groupColumns, ", ", builder)
	}
	if self.orderColumns != nil && len(self.orderColumns) > 0 {
		builder.WriteString(" ORDER BY ")
		handleStringsSplice(self.orderColumns, ", ", builder)
	}
	if self.limit > -1 {
		builder.WriteString(" LIMIT " + strconv.FormatInt(self.limit, 10))
	}
	if self.offset > -1 {
		builder.WriteString(" OFFSET " + strconv.FormatInt(self.offset, 10))
	}
	defer defaultPool.RecycleStringBuilder(builder)
	return builder.String()
}

func (self *PostgresSelectBuilder) CountSQL() string {
	builder := defaultPool.GetStringBuilder()
	if self.countField == "" {
		builder.WriteString("SELECT COUNT(*) as total FROM " + self.tableName)
	} else {
		builder.WriteString("SELECT COUNT(DISTINCT " + self.countField + ") as total FROM " + self.tableName)
	}
	if self.joins != nil && len(self.joins) > 0 {
		builder.WriteByte(' ')
		handleStringsSplice(self.joins, " ", builder)
	}
	if self.wheres != nil && len(self.wheres) > 0 {
		builder.WriteString(" WHERE ")
		handleStringsSplice(self.wheres, " AND ", builder)
	}
	defer defaultPool.RecycleStringBuilder(builder)
	return builder.String()
}

func (self *PostgresSelectBuilder) Select(fields ...string) SelectBuilder {
	if self.columns == nil {
		self.columns = make([]string, 0)
	}
	self.columns = slices.Concat(self.columns, fields)
	return self
}

func (self *PostgresSelectBuilder) CountField(field string) SelectBuilder {
	self.countField = field
	return self
}

func (self *PostgresSelectBuilder) Join(table string) *JoinBuilder {
	return newJoinBuilder(table, "JOIN", self)
}
func (self *PostgresSelectBuilder) LeftJoin(table string) *JoinBuilder {
	return newJoinBuilder(table, "LEFT JOIN", self)
}
func (self *PostgresSelectBuilder) RightJoin(table string) *JoinBuilder {
	return newJoinBuilder(table, "RIGHT JOIN", self)
}
func (self *PostgresSelectBuilder) InnerJoin(table string) *JoinBuilder {
	return newJoinBuilder(table, "INNER JOIN", self)
}

func (self *PostgresSelectBuilder) OrderBy(orderFields ...string) SelectBuilder {
	return self.addOrder("", orderFields)
}

func (self *PostgresSelectBuilder) OrderByAsc(orderFields ...string) SelectBuilder {
	return self.addOrder("ASC", orderFields)
}

func (self *PostgresSelectBuilder) OrderByDesc(orderFields ...string) SelectBuilder {
	return self.addOrder("DESC", orderFields)
}

func (self *PostgresSelectBuilder) addOrder(suffix string, fields []string) SelectBuilder {
	if self.orderColumns == nil {
		self.orderColumns = make([]string, 0)
	}
	if suffix == "" {
		self.orderColumns = slices.Concat(self.orderColumns, fields)
	} else {
		for _, field := range fields {
			self.orderColumns = append(self.orderColumns, field+" "+suffix)
		}
	}
	return self
}

func (self *PostgresSelectBuilder) GroupBy(fields ...string) SelectBuilder {
	if self.groupColumns == nil {
		self.groupColumns = make([]string, 0)
	}
	self.groupColumns = slices.Concat(self.groupColumns, fields)
	return self
}

func (self *PostgresSelectBuilder) Limit(limit int64) SelectBuilder {
	self.limit = limit
	return self
}

func (self *PostgresSelectBuilder) Offset(offset int64) SelectBuilder {
	self.offset = offset
	return self
}

func (self *PostgresSelectBuilder) addJoin(join string) {
	if self.joins == nil {
		self.joins = make([]string, 0)
	}
	self.joins = append(self.joins, join)
}
