package sqlbuild

type DeleteBuilder interface {
	SQLBuilder
	Condition
}

type PostgresDeleteBuilder struct {
	tableName string
	PostgresStatementParameterBuffer
	PostgresCondition
}

func (self *PostgresDeleteBuilder) table(tableName string) SQLBuilder {
	self.tableName = tableName
	return self
}

func (self *PostgresDeleteBuilder) SQL() string {
	builder := defaultPool.GetStringBuilder()
	builder.WriteString("DELETE FROM " + self.tableName)
	if self.wheres != nil && len(self.wheres) > 0 {
		builder.WriteString(" WHERE ")
		handleStringsSplice(self.wheres, " AND ", builder)
	}
	defer defaultPool.RecycleStringBuilder(builder)
	return builder.String()
}
