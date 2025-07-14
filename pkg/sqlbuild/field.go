package sqlbuild

// FieldCondition 数据库字段条件接口
type FieldCondition interface {
	// Eq 等于
	Eq(value any) *WhereBuilder
	// EqRaw 原始值等于
	EqRaw(value string) *WhereBuilder
	Ne(value any) *WhereBuilder
	NeRaw(value string) *WhereBuilder
	Le(value any) *WhereBuilder
	LeRaw(value string) *WhereBuilder
	Ge(value any) *WhereBuilder
	GeRaw(value string) *WhereBuilder
	Lt(value any) *WhereBuilder
	LtRaw(value string) *WhereBuilder
	Gt(value any) *WhereBuilder
	GtRaw(value string) *WhereBuilder
	Original(operate string, value any) *WhereBuilder
	OriginalRaw(operate, value string) *WhereBuilder
	Like(value string) *WhereBuilder
	In(value ...any) *WhereBuilder
	InRaw(value string) *WhereBuilder
	NotIn(value ...any) *WhereBuilder
	NotInRaw(value string) *WhereBuilder
	Between(start, end any) *WhereBuilder
	BetweenRaw(value string) *WhereBuilder
	NotBetween(start, end any) *WhereBuilder
	NotBetWeen(value string) *WhereBuilder
	IsNull() *WhereBuilder
	NotNull() *WhereBuilder
}

type Field struct {
	column    string
	prefix    string
	condition bool
	builder   *WhereBuilder
}

func newField(field, prefix string, condition bool, builder *WhereBuilder) *Field {
	return &Field{column: field, prefix: prefix, condition: condition, builder: builder}
}

func (self *Field) Eq(value any) *WhereBuilder {
	return self.saveCondition("=", value)
}

func (self *Field) Ne(value any) *WhereBuilder {
	return self.saveCondition("!=", value)
}

func (self *Field) Le(value any) *WhereBuilder {
	return self.saveCondition("<=", value)
}
func (self *Field) Ge(value any) *WhereBuilder {
	return self.saveCondition(">=", value)
}
func (self *Field) Lt(value any) *WhereBuilder {
	return self.saveCondition("<", value)
}
func (self *Field) Gt(value any) *WhereBuilder {
	return self.saveCondition(">", value)
}

func (self *Field) Like(value string) *WhereBuilder {
	return self.saveCondition("LIKE", value)
}

func (self *Field) In(value ...any) *WhereBuilder {
	return self.saveCondition("IN", value)
}

func (self *Field) NotIn(value ...any) *WhereBuilder {
	return self.saveCondition("NOT IN", value)
}

func (self *Field) Between(start, end any) *WhereBuilder {
	return self.saveCondition("BETWEEN", []any{start, end})
}

func (self *Field) NotBetween(start, end any) *WhereBuilder {
	return self.saveCondition("NOT BETWEEN", []any{start, end})
}

func (self *Field) IsNull() *WhereBuilder {
	return self.saveCondition("ISNULL", nil)
}

func (self *Field) NotNull() *WhereBuilder {
	return self.saveCondition("NOTNULL", nil)
}

func (self *Field) EqRaw(value string) *WhereBuilder {
	return self.saveConditionRaw("=", value)
}

func (self *Field) NeRaw(value string) *WhereBuilder {
	return self.saveConditionRaw("!=", value)
}

func (self *Field) LeRaw(value string) *WhereBuilder {
	return self.saveConditionRaw(">", value)
}

func (self *Field) GeRaw(value string) *WhereBuilder {
	return self.saveConditionRaw("<", value)
}

func (self *Field) LtRaw(value string) *WhereBuilder {
	return self.saveConditionRaw(">=", value)
}

func (self *Field) GtRaw(value string) *WhereBuilder {
	return self.saveConditionRaw("<=", value)
}

func (self *Field) Original(operate string, value any) *WhereBuilder {
	return self.saveCondition(operate, value)
}

func (self *Field) OriginalRaw(operate, value string) *WhereBuilder {
	return self.saveConditionRaw(operate, value)
}

func (self *Field) InRaw(value string) *WhereBuilder {
	return self.saveConditionRaw("IN", value)
}

func (self *Field) NotInRaw(value string) *WhereBuilder {
	return self.saveConditionRaw("NOT IN", value)
}

func (self *Field) BetweenRaw(value string) *WhereBuilder {
	return self.saveConditionRaw("BETWEEN", value)
}

func (self *Field) NotBetWeen(value string) *WhereBuilder {
	return self.saveConditionRaw("NOT BETWEEN", value)
}

func (self *Field) Recycle() {
	self.prefix = ""
	self.column = ""
	self.condition = false
	self.builder = nil
}

func (self *Field) saveConditionRaw(operator, value string) *WhereBuilder {
	if !self.condition {
		return self.builder
	}
	// 回收字段对象
	defer defaultPool.RecycleField(self)
	self.builder.addRawCondition(self.column, operator, self.prefix, value)
	return self.builder
}
func (self *Field) saveCondition(operator string, value any) *WhereBuilder {
	if !self.condition {
		return self.builder
	}
	defer defaultPool.RecycleField(self)
	self.builder.addCondition(self.column, operator, self.prefix, value)
	return self.builder
}
