package sqlbuild

import (
	"strconv"
	"strings"
)

// WhereBuilder 通用where语句构造器
type WhereBuilder struct {
	buffer      *strings.Builder
	isInit      bool
	builder     SQLBuilder
	whereBuffer Condition
}

func newWhereBuilder(sqlBuilder SQLBuilder, condition Condition) *WhereBuilder {
	return &WhereBuilder{
		buffer:      defaultPool.GetStringBuilder(),
		builder:     sqlBuilder,
		whereBuffer: condition,
	}
}

func (self *WhereBuilder) And(field string) *Field {
	return newField(field, "AND", true, self)
}
func (self *WhereBuilder) AndByCondition(condition bool, field string) *Field {
	return newField(field, "AND", condition, self)
}
func (self *WhereBuilder) Or(field string) *Field {
	return newField(field, "OR", true, self)
}
func (self *WhereBuilder) OrByCondition(condition bool, field string) *Field {
	return newField(field, "OR", condition, self)
}
func (self *WhereBuilder) addCondition(field, operator, prefix string, value any) {
	if self.isInit {
		self.buffer.WriteString(" " + prefix + " ")
	} else {
		self.isInit = true
	}
	self.buffer.WriteString(field + " " + operator + " ")
	if value == nil {
		return
	}
	sliceValue, ok := value.([]any)
	if !ok {
		paramIndex := self.builder.addParameter(value)
		self.buffer.WriteString("$" + strconv.Itoa(paramIndex))
		return
	}
	switch operator {
	case "BETWEEN":
		startIndex, endIndex := self.builder.addParameter(sliceValue[0]), self.builder.addParameter(sliceValue[1])
		self.buffer.WriteString("$" + strconv.Itoa(startIndex) + " AND " + "$" + strconv.Itoa(endIndex))
	case "NOT BETWEEN":
		startIndex, endIndex := self.builder.addParameter(sliceValue[0]), self.builder.addParameter(sliceValue[1])
		self.buffer.WriteString("$" + strconv.Itoa(startIndex) + " AND " + "$" + strconv.Itoa(endIndex))
	default:
		self.buffer.WriteByte('(')
		for index, param := range sliceValue {
			if index > 0 {
				self.buffer.WriteByte(',')
			}
			paramIndex := self.builder.addParameter(param)
			self.buffer.WriteString("$" + strconv.Itoa(paramIndex))
		}
		self.buffer.WriteByte(')')
	}
}

func (self *WhereBuilder) addRawCondition(field, operator, prefix, value string) {
	if self.isInit {
		self.buffer.WriteString(" " + prefix + " ")
	} else {
		self.isInit = true
	}
	self.buffer.WriteString(field + " " + operator + " " + value)
}

func (self *WhereBuilder) build() SQLBuilder {
	if self.buffer.Len() > 0 {
		self.whereBuffer.addCondition(self.buffer.String())
	}
	defaultPool.RecycleStringBuilder(self.buffer)
	return self.builder
}
func (self *WhereBuilder) BuildAsSelect() SelectBuilder {
	return self.build().(SelectBuilder)
}
func (self *WhereBuilder) BuildAsUpdate() UpdateBuilder {
	return self.build().(UpdateBuilder)
}
func (self *WhereBuilder) BuildAsDelete() DeleteBuilder {
	return self.build().(DeleteBuilder)
}
