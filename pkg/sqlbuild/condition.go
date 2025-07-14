package sqlbuild

// Condition 条件查询接口
type Condition interface {
	Where(field string) *Field
	WhereByCondition(condition bool, field string) *Field
	addCondition(where string)
}

type PostgresCondition struct {
	builder SQLBuilder
	wheres  []string
}

func (self *PostgresCondition) Where(column string) *Field {
	return self.WhereByCondition(true, column)
}

func (self *PostgresCondition) WhereByCondition(condition bool, column string) *Field {
	whereBUilder := newWhereBuilder(self.builder, self)
	field := defaultPool.GetField()
	field.column = column
	field.condition = condition
	field.builder = whereBUilder
	return field
}

func (self *PostgresCondition) addCondition(where string) {
	if self.wheres == nil {
		self.wheres = make([]string, 0)
	}
	self.wheres = append(self.wheres, where)
}
