package SkyLine

func (i *Integer) Type_Object() Type_Object {
	return IntegerType
}

func (f *Float) Type_Object() Type_Object {
	return FloatType
}

func (b *Boolean_Object) Type_Object() Type_Object {
	return BooleanType
}

func (n *Nil) Type_Object() Type_Object {
	return NilType
}

func (rv *ReturnValue) Type_Object() Type_Object {
	return ReturnValueType
}

func (e *Error) Type_Object() Type_Object {
	return ErrorType
}

func (f *Function) Type_Object() Type_Object {
	return FunctionType
}

func (s *String) Type_Object() Type_Object {
	return StringType
}

func (b *Builtin) Type_Object() Type_Object {
	return BuiltinType
}

func (*Array) Type_Object() Type_Object {
	return ArrayType
}

func (*Hash) Type_Object() Type_Object {
	return HashType
}

func (q *Quote) Type_Object() Type_Object {
	return QuoteType
}

func (m *Macro) Type_Object() Type_Object {
	return MacroType
}
