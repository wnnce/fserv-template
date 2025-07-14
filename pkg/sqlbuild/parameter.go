package sqlbuild

// StatementParameterBuffer SQL命令参数缓存接口
type StatementParameterBuffer interface {
	Args() []any
	addParameter(param any) int
}

type PostgresStatementParameterBuffer struct {
	args []any
}

func (self *PostgresStatementParameterBuffer) Args() []any {
	return self.args
}

func (self *PostgresStatementParameterBuffer) addParameter(param any) int {
	if self.args == nil {
		self.args = make([]any, 0)
	}
	// map不能使用切片做为key 其次同一条sql中重复参数的比例不会太高 故而先不使用map存储参数索引
	/*if self.argsIndexMap == nil {
		self.argsIndexMap = make(map[any]int)
	}
	index, ok := self.argsIndexMap[param]
	if ok {
		return index
	}*/
	self.args = append(self.args, param)
	/*index = len(self.args)
	self.argsIndexMap[param] = index*/
	return len(self.args)
}
