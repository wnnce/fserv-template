package sqlbuild

// JoinBuilder join语句构造器
type JoinBuilder struct {
	builder   SelectBuilder
	joinType  string
	tableName string
}

func newJoinBuilder(table, joinType string, builder SelectBuilder) *JoinBuilder {
	return &JoinBuilder{
		builder:   builder,
		joinType:  joinType,
		tableName: table,
	}
}

func (self *JoinBuilder) On(column string) *Field {
	whereBuilder := newWhereBuilder(self.builder, self)
	field := defaultPool.GetField()
	field.column = column
	field.condition = true
	field.builder = whereBuilder
	return field
}

func (self *JoinBuilder) addCondition(where string) {
	self.builder.addJoin(self.joinType + " " + self.tableName + " ON " + where)
}

func (self *JoinBuilder) Where(column string) *Field {
	panic("join not where")
}

func (self *JoinBuilder) WhereByCondition(condition bool, column string) *Field {
	panic("join not where")
}
