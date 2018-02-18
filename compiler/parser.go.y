%{
package compiler

import (
    "../vm"
    "strconv"
)
%}

%union{
    parameter_list       []*Parameter
    argument_list        []Expression

    statement            Statement
    statement_list       []Statement

    expression           Expression
    expression_list      []Expression

    block                *Block
    elif_list            []*Elif

    basic_type_specifier *TypeSpecifier
    type_specifier       *TypeSpecifier

    array_dimension      *ArrayDimension
    array_dimension_list []*ArrayDimension

    tok                  Token
}

%token<tok> IF ELSE ELIF FOR RETURN_T BREAK CONTINUE
        LP RP LC RC LB RB
        SEMICOLON COMMA
        ASSIGN_T
        LOGICAL_AND LOGICAL_OR
        EQ NE GT GE LT LE
        ADD SUB MUL DIV
        INT_LITERAL DOUBLE_LITERAL STRING_LITERAL TRUE_T FALSE_T
        NULL_T
        IDENTIFIER
        EXCLAMATION DOT
        BOOLEAN_T INT_T DOUBLE_T STRING_T NEW

%type <expression> expression expression_opt
      assignment_expression
      logical_and_expression logical_or_expression
      equality_expression relational_expression
      additive_expression multiplicative_expression
      unary_expression primary_expression primary_no_new_array
      array_literal array_creation
%type   <expression_list> expression_list

%type <statement> statement
      if_statement for_statement
      return_statement break_statement continue_statement
      declaration_statement
%type <statement_list> statement_list
%type <parameter_list> parameter_list
%type <argument_list> argument_list
%type <block> block
%type <elif_list> elif_list
%type <type_specifier> type_specifier basic_type_specifier
%type <array_dimension> dimension_expression
%type <array_dimension_list> dimension_expression_list dimension_list

%%

translation_unit
        : definition_or_statement
        | translation_unit definition_or_statement
        ;
definition_or_statement
        : function_definition
        | statement
        {
            l := yylex.(*Lexer)
            l.compiler.statementList = append(l.compiler.statementList, $1)
        }
        ;
basic_type_specifier
        : BOOLEAN_T
        {
            $$ = &TypeSpecifier{basicType: vm.BooleanType}
            $$.SetPosition($1.Position())
        }
        | INT_T
        {
            $$ = &TypeSpecifier{basicType: vm.IntType}
            $$.SetPosition($1.Position())
        }
        | DOUBLE_T
        {
            $$ = &TypeSpecifier{basicType: vm.DoubleType}
            $$.SetPosition($1.Position())
        }
        | STRING_T
        {
            $$ = &TypeSpecifier{basicType: vm.StringType}
            $$.SetPosition($1.Position())
        }
        ;
type_specifier
        : basic_type_specifier
        {
            $$ = $1
        }
        | type_specifier LB RB
        {
            $1.appendDerive(&ArrayDerive{})
            $$ = $1
        }
        ;
function_definition
        : type_specifier IDENTIFIER LP parameter_list RP block
        {
            l := yylex.(*Lexer)
            l.compiler.functionDefine($1, $2.Lit, $4, $6)
        }
        | type_specifier IDENTIFIER LP RP block
        {
            l := yylex.(*Lexer)
            l.compiler.functionDefine($1, $2.Lit, []*Parameter{}, $5)
        }
        | type_specifier IDENTIFIER LP parameter_list RP SEMICOLON
        {
            l := yylex.(*Lexer)
            l.compiler.functionDefine($1, $2.Lit, $4, nil)
        }
        | type_specifier IDENTIFIER LP RP SEMICOLON
        {
            l := yylex.(*Lexer)
            l.compiler.functionDefine($1, $2.Lit, []*Parameter{}, nil)
        }
        ;
parameter_list
        : type_specifier IDENTIFIER
        {
            parameter := &Parameter{typeSpecifier: $1, name: $2.Lit}
            parameter.SetPosition($1.Position())
            $$ = []*Parameter{parameter}
        }
        | parameter_list COMMA type_specifier IDENTIFIER
        {
            $$ = append($1, &Parameter{typeSpecifier: $3, name: $4.Lit})
        }
        ;
argument_list
        : assignment_expression
        {
            $$ = []Expression{$1}
        }
        | argument_list COMMA assignment_expression
        {
            $$ = append($1, $3)
        }
        ;
statement_list
        : statement
        {
            $$ = []Statement{$1}
        }
        | statement_list statement
        {
            $$ = append($1, $2)
        }
        ;
expression
        : assignment_expression
        | expression COMMA assignment_expression
        {
            $$ = &CommaExpression{left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
assignment_expression
        : logical_or_expression
        | primary_expression ASSIGN_T assignment_expression
        {
            $$ = &AssignExpression{left: $1, operand: $3}
            $$.SetPosition($1.Position())
        }
        ;
logical_or_expression
        : logical_and_expression
        | logical_or_expression LOGICAL_OR logical_and_expression
        {
            $$ = &BinaryExpression{operator: LogicalOrOperator, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
logical_and_expression
        : equality_expression
        | logical_and_expression LOGICAL_AND equality_expression
        {
            $$ = &BinaryExpression{operator: LogicalAndOperator, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
equality_expression
        : relational_expression
        | equality_expression EQ relational_expression
        {
            $$ = &BinaryExpression{operator: EqOperator, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        | equality_expression NE relational_expression
        {
            $$ = &BinaryExpression{operator: NeOperator, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
relational_expression
        : additive_expression
        | relational_expression GT additive_expression
        {
            $$ = &BinaryExpression{operator: GtOperator, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        | relational_expression GE additive_expression
        {
            $$ = &BinaryExpression{operator: GeOperator, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        | relational_expression LT additive_expression
        {
            $$ = &BinaryExpression{operator: LtOperator, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        | relational_expression LE additive_expression
        {
            $$ = &BinaryExpression{operator: LeOperator, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
additive_expression
        : multiplicative_expression
        | additive_expression ADD multiplicative_expression
        {
            $$ = &BinaryExpression{operator: AddOperator, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        | additive_expression SUB multiplicative_expression
        {
            $$ = &BinaryExpression{operator: SubOperator, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
multiplicative_expression
        : unary_expression
        | multiplicative_expression MUL unary_expression
        {
            $$ = &BinaryExpression{operator: MulOperator, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        | multiplicative_expression DIV unary_expression
        {
            $$ = &BinaryExpression{operator: DivOperator, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
unary_expression
        : primary_expression
        | SUB unary_expression
        {
            $$ = &MinusExpression{operand: $2}
            $$.SetPosition($1.Position())
        }
        | EXCLAMATION unary_expression
        {
            $$ = &LogicalNotExpression{operand: $2}
            $$.SetPosition($1.Position())
        }
        ;
primary_expression
        : primary_no_new_array
        | array_creation
        ;
primary_no_new_array
        : primary_no_new_array LB expression RB
        {
            $$ = &IndexExpression{array: $1, index: $3}
            $$.SetPosition($1.Position())
        }
        | primary_expression DOT IDENTIFIER
        {
            $$ = &MemberExpression{expression: $1, memberName: $3.Lit}
            $$.SetPosition($1.Position())
        }
        | primary_expression LP argument_list RP
        {
            $$ = &FunctionCallExpression{function: $1, argumentList: $3}
            $$.SetPosition($1.Position())
        }
        | primary_expression LP RP
        {
            $$ = &FunctionCallExpression{function: $1, argumentList: []Expression{}}
            $$.SetPosition($1.Position())
        }
        | LP expression RP
        {
            $$ = $2
        }
        | IDENTIFIER
        {
            $$ = &IdentifierExpression{name: $1.Lit}
            $$.SetPosition($1.Position())
        }
        | INT_LITERAL
        {
            value, _ := strconv.Atoi($1.Lit)
            $$ = &IntExpression{intValue: value}
            $$.SetPosition($1.Position())
        }
        | DOUBLE_LITERAL
        {
            value, _ := strconv.ParseFloat($1.Lit, 64)
            $$ = &DoubleExpression{doubleValue: value}
            $$.SetPosition($1.Position())
        }
        | STRING_LITERAL
        {
            $$ = &StringExpression{stringValue: $1.Lit}
            $$.SetPosition($1.Position())
        }
        | TRUE_T
        {
            $$ = &BooleanExpression{booleanValue: true}
            $$.SetPosition($1.Position())
        }
        | FALSE_T
        {
            $$ = &BooleanExpression{booleanValue: false}
            $$.SetPosition($1.Position())
        }
        | NULL_T
        {
            $$ = &NullExpression{}
            $$.SetPosition($1.Position())
        }
        | array_literal
        ;
array_literal
        : LC expression_list RC
        {
            $$ = &ArrayLiteralExpression{arrayLiteral: $2}
            $$.SetPosition($1.Position())
        }
        | LC expression_list COMMA RC
        {
            $$ = &ArrayLiteralExpression{arrayLiteral: $2}
            $$.SetPosition($1.Position())
        }
        ;
array_creation
        : NEW basic_type_specifier dimension_expression_list
        {
            $$ = &ArrayCreation{dimensionList: $3}
            $$.setType($2)
            $$.SetPosition($1.Position())
        }
        | NEW basic_type_specifier dimension_expression_list dimension_list
        {
            $$ = &ArrayCreation{dimensionList: append($3, $4...)}
            $$.setType($2)
            $$.SetPosition($1.Position())
        }
        ;
dimension_expression_list
        : dimension_expression
        {
            $$ = []*ArrayDimension{$1}
        }
        | dimension_expression_list dimension_expression
        {
            $$ = append($1, $2)
        }
        ;
dimension_expression
        : LB expression RB
        {
            $$ = &ArrayDimension{expression: $2}
        }
        ;
dimension_list
        : LB RB
        {
            $$ = []*ArrayDimension{&ArrayDimension{}}
        }
        | dimension_list LB RB
        {
            $$ = append($1, &ArrayDimension{})
        }
        ;
expression_list
        :
        {
            $$ = nil
        }
        | assignment_expression
        {
            $$ = []Expression{$1}
        }
        | expression_list COMMA assignment_expression
        {
            $$ = append($1, $3)
        }
        ;
statement
        : expression SEMICOLON
        {
            $$ = &ExpressionStatement{expression: $1}
            $$.SetPosition($1.Position())
        }
        | if_statement
        | for_statement
        | return_statement
        | break_statement
        | continue_statement
        | declaration_statement
        ;
if_statement
        : IF expression block
        {
            $$ = &IfStatement{condition: $2, thenBlock: $3, elifList: []*Elif{}, elseBlock: nil}
            $$.SetPosition($1.Position())
        }
        | IF expression block ELSE block
        {
            $$ = &IfStatement{condition: $2, thenBlock: $3, elifList: []*Elif{}, elseBlock: $5}
            $$.SetPosition($1.Position())
        }
        | IF expression block elif_list
        {
            $$ = &IfStatement{condition: $2, thenBlock: $3, elifList: $4, elseBlock: nil}
            $$.SetPosition($1.Position())
        }
        | IF expression block elif_list ELSE block
        {
            $$ = &IfStatement{condition: $2, thenBlock: $3, elifList: $4, elseBlock: $6}
            $$.SetPosition($1.Position())
        }
        ;
elif_list
        : ELIF expression block
        {
            $$ = []*Elif{&Elif{condition: $2, block: $3}}
        }
        | elif_list ELIF expression block
        {
            $$ = append($1, &Elif{condition: $3, block: $4})
        }
        ;
for_statement
        : FOR LP expression_opt SEMICOLON expression_opt SEMICOLON expression_opt RP block
        {
            $$ = &ForStatement{init: $3, condition: $5, post: $7, block: $9}
            $$.SetPosition($1.Position())
            $9.parent = &StatementBlockInfo{statement: $$}
        }
        ;
expression_opt
        :
        {
            $$ = nil
        }
        | expression
        ;
return_statement
        : RETURN_T expression_opt SEMICOLON
        {
            $$ = &ReturnStatement{returnValue: $2};
            $$.SetPosition($1.Position())
        }
        ;
break_statement
        : BREAK SEMICOLON
        {
            $$ = &BreakStatement{}
            $$.SetPosition($1.Position())
        }
        ;
continue_statement
        : CONTINUE SEMICOLON
        {
            $$ = &ContinueStatement{}
            $$.SetPosition($1.Position())
        }
        ;
declaration_statement
        : type_specifier IDENTIFIER SEMICOLON
        {
            $$ = &Declaration{typeSpecifier: $1, name: $2.Lit, variableIndex: -1}
            $$.SetPosition($1.Position())
        }
        | type_specifier IDENTIFIER ASSIGN_T expression SEMICOLON
        {
            $$ = &Declaration{typeSpecifier: $1, name: $2.Lit, initializer: $4, variableIndex: -1}
            $$.SetPosition($1.Position())
        }
        ;
block
        : LC
        {
            l := yylex.(*Lexer)
            l.compiler.currentBlock = &Block{outerBlock: l.compiler.currentBlock}
            $<block>$ = l.compiler.currentBlock
        }
          statement_list RC
        {
            currentBlock := $<block>2
            currentBlock.statementList = $3

            l := yylex.(*Lexer)

            $<block>$ = l.compiler.currentBlock
            l.compiler.currentBlock = currentBlock.outerBlock
        }
        | LC RC
        {
            l := yylex.(*Lexer)
            $<block>$ = &Block{outerBlock: l.compiler.currentBlock}
        }
        ;
%%
