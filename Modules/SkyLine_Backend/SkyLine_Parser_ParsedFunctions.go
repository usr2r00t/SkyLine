package SkyLine

import (
	"fmt"
	"strconv"
)

func (parser *Parser) parseStatement() Statement {
	switch parser.CurrentToken.Token_Type {
	case LET:
		return parser.parseLetStatement()
	case RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: parser.CurrentToken}
	if !parser.ExpectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Ident{
		Token: parser.CurrentToken,
		Value: parser.CurrentToken.Literal,
	}

	if !parser.ExpectPeek(ASSIGN) {
		return nil
	}

	parser.NextToken()

	stmt.Value = parser.parseExpression(LOWEST)

	for parser.PeekTokenIs(SEMICOLON) {
		parser.NextToken()
	}

	return stmt
}

func (parser *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{
		Token: parser.CurrentToken,
	}

	parser.NextToken()

	stmt.ReturnValue = parser.parseExpression(LOWEST)

	for parser.PeekTokenIs(SEMICOLON) {
		parser.NextToken()
	}

	return stmt
}

func (parser *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{
		Token:      parser.CurrentToken,
		Expression: parser.parseExpression(LOWEST),
	}

	if parser.PeekTokenIs(SEMICOLON) {
		parser.NextToken()
	}

	return stmt
}

func (parser *Parser) parseExpression(precedence int) Expression {
	prefix := parser.PrefixParseFns[parser.CurrentToken.Token_Type]
	if prefix == nil {
		msg := StatementMap_Parser_SingleArgument["PREFIX_PARSE_ERR"](
			"No prefix parse function found for token",
			"Could not locate parse function for the parsed token",
			parser.CurrentToken.Literal,
		)
		parser.Errors = append(parser.Errors, msg)
		return nil
	}

	leftExp := prefix()

	for !parser.CurrentTokenIs(SEMICOLON) && precedence < parser.peekPrecedence() {
		infix := parser.InfixParseFns[parser.PeekToken.Token_Type]
		if infix == nil {
			return leftExp
		}

		parser.NextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (parser *Parser) parseIdent() Expression {
	return &Ident{
		Token: parser.CurrentToken,
		Value: parser.CurrentToken.Literal,
	}
}

func (parser *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: parser.CurrentToken}

	val, err := strconv.ParseInt(parser.CurrentToken.Literal, 0, 64)
	if err != nil {
		msg := StatementMap_Parser_SingleArgument["INT_PARSE_ERR"](
			"Could not parse value as INTEGER",
			fmt.Sprint(err),
			parser.CurrentToken.Literal,
		)
		parser.Errors = append(parser.Errors, msg)
		return nil
	}

	lit.Value = val
	return lit
}

func (parser *Parser) parseFloatLiteral() Expression {
	val, err := strconv.ParseFloat(parser.CurrentToken.Literal, 64)
	if err != nil {
		msg := StatementMap_Parser_SingleArgument["FLOAT_PARSE_ERR"](
			"Could not parse value as FLOAT",
			fmt.Sprint(err),
			parser.CurrentToken.Literal,
		)
		parser.Errors = append(parser.Errors, msg)
		return nil
	}

	return &FloatLiteral{
		Token: parser.CurrentToken,
		Value: val,
	}
}

func (parser *Parser) parsePrefixExpression() Expression {
	expr := &PrefixExpression{
		Token:    parser.CurrentToken,
		Operator: parser.CurrentToken.Literal,
	}

	parser.NextToken()

	expr.Right = parser.parseExpression(PREFIX)
	return expr
}

func (parser *Parser) peekPrecedence() int {
	if p, ok := Precedences[parser.PeekToken.Token_Type]; ok {
		return p
	}
	return LOWEST
}

func (parser *Parser) curPrecedence() int {
	if p, ok := Precedences[parser.CurrentToken.Token_Type]; ok {
		return p
	}
	return LOWEST
}

func (parser *Parser) parseInfixExpression(left Expression) Expression {
	expr := &InfixExpression{
		Token:    parser.CurrentToken,
		Operator: parser.CurrentToken.Literal,
		Left:     left,
	}

	prec := parser.curPrecedence()

	parser.NextToken()

	expr.Right = parser.parseExpression(prec)
	return expr
}

func (parser *Parser) parseBoolean() Expression {
	return &Boolean_AST{
		Token: parser.CurrentToken,
		Value: parser.CurrentTokenIs(TRUE),
	}
}

func (parser *Parser) parseGroupedExpression() Expression {
	parser.NextToken()

	expr := parser.parseExpression(LOWEST)

	if !parser.ExpectPeek(RPAREN) {
		return nil
	}

	return expr
}

func (parser *Parser) parseIfExpression() Expression {
	expr := &ConditionalExpression{Token: parser.CurrentToken}

	parser.NextToken()
	expr.Condition = parser.parseExpression(LOWEST)

	if !parser.ExpectPeek(LBRACE) {
		return nil
	}

	expr.Consequence = parser.parseBlockStatement()

	if parser.PeekTokenIs(ELSE) {
		parser.NextToken()

		if !parser.ExpectPeek(LBRACE) {
			return nil
		}

		expr.Alternative = parser.parseBlockStatement()
	}

	return expr
}

func (parser *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{
		Token:      parser.CurrentToken,
		Statements: []Statement{},
	}

	parser.NextToken()
	for !parser.CurrentTokenIs(RBRACE) && !parser.CurrentTokenIs(EOF) {
		stmt := parser.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		parser.NextToken()
	}

	return block
}

func (parser *Parser) parseFunctionLiteral() Expression {
	lit := &FunctionLiteral{Token: parser.CurrentToken}

	if !parser.ExpectPeek(LPAREN) {
		return nil
	}

	lit.Parameters = parser.parseFunctionParameters()

	if !parser.ExpectPeek(LBRACE) {
		return nil
	}

	lit.Body = parser.parseBlockStatement()

	return lit
}

func (parser *Parser) parseFunctionParameters() []*Ident {
	idents := []*Ident{}

	if parser.PeekTokenIs(RPAREN) {
		parser.NextToken()
		return idents
	}

	parser.NextToken()

	ident := &Ident{
		Token: parser.CurrentToken,
		Value: parser.CurrentToken.Literal,
	}
	idents = append(idents, ident)

	for parser.PeekTokenIs(COMMA) || parser.PeekTokenIs(COLON) {
		parser.NextToken()
		parser.NextToken()
		ident := &Ident{
			Token: parser.CurrentToken,
			Value: parser.CurrentToken.Literal,
		}
		idents = append(idents, ident)
	}

	if !parser.ExpectPeek(RPAREN) {
		return nil
	}

	return idents
}

func (parser *Parser) parseExpressionList(end Token_Type) []Expression {
	list := make([]Expression, 0)

	if parser.PeekTokenIs(end) {
		parser.NextToken()
		return list
	}

	parser.NextToken()
	list = append(list, parser.parseExpression(LOWEST))

	for parser.PeekTokenIs(COMMA) {
		parser.NextToken()
		parser.NextToken()
		list = append(list, parser.parseExpression(LOWEST))
	}

	if !parser.ExpectPeek(end) {
		return nil
	}

	return list
}

func (parser *Parser) parseCallExpression(function Expression) Expression {
	return &CallExpression{
		Token:     parser.CurrentToken,
		Function:  function,
		Arguments: parser.parseExpressionList(RPAREN),
	}
}

func (parser *Parser) parseStringLiteral() Expression {
	return &StringLiteral{
		Token: parser.CurrentToken,
		Value: parser.CurrentToken.Literal,
	}
}

func (parser *Parser) parseArrayLiteral() Expression {
	return &ArrayLiteral{
		Token:    parser.CurrentToken,
		Elements: parser.parseExpressionList(RBRACKET),
	}
}

func (parser *Parser) parseIndexExpression(left Expression) Expression {
	expr := &IndexExpression{
		Token: parser.CurrentToken,
		Left:  left,
	}

	parser.NextToken()
	expr.Index = parser.parseExpression(LOWEST)

	if !parser.ExpectPeek(RBRACKET) {
		return nil
	}

	return expr
}

func (parser *Parser) parseHashLiteral() Expression {
	hash := &HashLiteral{
		Token: parser.CurrentToken,
		Pairs: make(map[Expression]Expression),
	}

	for !parser.PeekTokenIs(RBRACE) {
		parser.NextToken()
		key := parser.parseExpression(LOWEST)

		if !parser.ExpectPeek(COLON) {
			return nil
		}

		parser.NextToken()
		value := parser.parseExpression(LOWEST)
		hash.Pairs[key] = value

		if !parser.PeekTokenIs(RBRACE) && !parser.ExpectPeek(COMMA) {
			return nil
		}
	}

	if !parser.ExpectPeek(RBRACE) {
		return nil
	}

	return hash
}

func (parser *Parser) parseMacroLiteral() Expression {
	tok := parser.CurrentToken

	if !parser.ExpectPeek(LPAREN) {
		return nil
	}

	params := parser.parseFunctionParameters()

	if !parser.ExpectPeek(LBRACE) {
		return nil
	}

	body := parser.parseBlockStatement()

	return &MacroLiteral{
		Token:      tok,
		Parameters: params,
		Body:       body,
	}
}
