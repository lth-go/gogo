%{
package compiler

import (
    "strconv"
)
%}

%union{
    statement            Statement
    statement_list       []Statement

    expression           Expression
    expression_list      []Expression

    parameter            *Parameter
    parameter_list       []*Parameter

    block                *Block
    else_if              []*ElseIf

    type_specifier       *Type

    import_spec          *Import
    import_spec_list     []*Import

    function_definition  *FunctionDefinition

    tok                  Token
}

%token<tok> IF ELSE FOR RETURN BREAK CONTINUE
    LP RP LC RC LB RB
    SEMICOLON COMMA COLON
    ASSIGN
    LOGICAL_AND LOGICAL_OR
    EQ NE GT GE LT LE
    ADD SUB MUL DIV
    INT_LITERAL FLOAT_LITERAL STRING_LITERAL
    TRUE FALSE NIL
    IDENTIFIER
    EXCLAMATION DOT
    IMPORT VAR FUNC
    TYPE STRUCT MAP

%type <import_spec> import_declaration
%type <import_spec_list> import_declaration_list

%type <expression> expression expression_or_nil
    logical_and_expression logical_or_expression
    equality_expression relational_expression
    additive_expression multiplicative_expression
    unary_expression primary_expression
%type <expression_list> expression_list expression_list_or_nil argument_list

%type <statement> statement simple_statement_or_nil
    simple_statement
    if_statement for_statement
    return_statement break_statement continue_statement
    declaration_statement assign_statement
%type <statement_list> statement_list
/* TODO: 临时处理 */
%type <parameter> receiver_or_nil
%type <parameter_list> parameter_list parameters
    result_or_nil result
    type_list_or_nil type_list
%type <block> block block_or_nil
%type <else_if> else_if
%type <type_specifier> type_specifier composite_type array_type_specifier func_type_specifier signature map_type_specifier

%%

translation_unit
        : import_declaration_list_or_nil top_level_decl_list
        ;
import_declaration_list_or_nil
        : /* empty */
        {
            setImportList(nil)
        }
        | import_declaration_list
        {
            setImportList($1)
        }
        ;
import_declaration_list
        : import_declaration
        {
            $$ = CreateImportList($1)
        }
        | import_declaration_list import_declaration
        {
            $$ = append($1, $2)
        }
        ;
import_declaration
        : IMPORT STRING_LITERAL SEMICOLON
        {
            $$ = CreateImport($2.Lit)
        }
        ;
top_level_decl_list
        :
        | top_level_decl_list top_level_decl SEMICOLON
        ;
top_level_decl
        : declaration
        | function_definition
        ;
declaration
        : var_decl
        ;
var_decl
        : VAR IDENTIFIER type_specifier
        {
            AddDeclList(NewDeclaration($1.Position(), $3, $2.Lit, nil))
        }
        | VAR IDENTIFIER type_specifier ASSIGN expression
        {
            AddDeclList(NewDeclaration($1.Position(), $3, $2.Lit, $5))
        }
        ;
array_type_specifier
        : LB RB type_specifier
        {
            $$ = CreateArrayType($3, $1.Position())
            $$.SetPosition($1.Position())
        }
        ;
map_type_specifier
        : MAP LB type_specifier RB type_specifier
        {
            $$ = CreateMapType($3, $5, $1.Position())
        }
        ;
func_type_specifier
        : FUNC signature
        {
            $$ = $2
            $$.SetPosition($1.Position())
        }
        ;
type_specifier
        : IDENTIFIER
        {
            $$ = CreateTypeByName($1.Lit, $1.Position())
        }
        | IDENTIFIER DOT IDENTIFIER
        {
            $$ = CreateTypeByName($1.Lit + "." + $3.Lit, $1.Position())
        }
        | composite_type
        | func_type_specifier
        ;
composite_type
        : array_type_specifier
        | map_type_specifier
        ;
function_definition
        : FUNC receiver_or_nil IDENTIFIER signature block_or_nil
        {
            createFunctionDefine($1.Position(), $2, $3.Lit, $4, $5)
        }
        ;
receiver_or_nil
        :
        {
            $$ = nil
        }
        | LP IDENTIFIER type_specifier RP
        {
            $$ = NewParameter($3, $2.Lit)
        }
        ;
parameter_list
        : IDENTIFIER type_specifier
        {
            $$ = []*Parameter{NewParameter($2, $1.Lit)}
        }
        | parameter_list COMMA IDENTIFIER type_specifier
        {
            $$ = append($1, NewParameter($4, $3.Lit))
        }
        ;
parameters
        : LP RP
        {
            $$ = []*Parameter{}
        }
        | LP parameter_list RP
        {
            $$ = $2
        }
        ;
argument_list
        : expression
        {
            $$ = []Expression{$1}
        }
        | argument_list COMMA expression
        {
            $$ = append($1, $3)
        }
        ;
signature
        : parameters result_or_nil
        {
            $$ = CreateFuncType($1, $2)
        }
        ;
result_or_nil
        :
        {
            $$ = nil
        }
        | result
        ;
result
        : type_specifier
        {
            $$ = []*Parameter{NewParameter($1, "")}
        }
        | LP type_list_or_nil RP
        {
            $$ = $2
        }
        ;
type_list_or_nil
        :
        {
            $$ = nil
        }
        | type_list
        ;
type_list
        : type_specifier
        {
            $$ = []*Parameter{NewParameter($1, "")}
        }
        | type_list COMMA type_specifier
        {
            $$ = append($1, NewParameter($3, ""))
        }
        ;
statement_list
        : statement SEMICOLON
        {
            $$ = []Statement{$1}
        }
        | statement_list statement SEMICOLON
        {
            $$ = append($1, $2)
        }
        ;
expression
        : logical_or_expression
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
        : INT_LITERAL
        {
            value, _ := strconv.Atoi($1.Lit)
            $$ = &IntExpression{Value: value}
            $$.SetPosition($1.Position())
        }
        | FLOAT_LITERAL
        {
            value, _ := strconv.ParseFloat($1.Lit, 64)
            $$ = CreateFloatExpression($1.Position(), value)
        }
        | STRING_LITERAL
        {
            $$ = NewStringExpression($1.Position(), $1.Lit)
        }
        | TRUE
        {
            $$ = &BooleanExpression{Value: true}
            $$.SetPosition($1.Position())
        }
        | FALSE
        {
            $$ = &BooleanExpression{Value: false}
            $$.SetPosition($1.Position())
        }
        | NIL
        {
            $$ = createNilExpression($1.Position())
        }
        | composite_type LC expression_list_or_nil RC
        {
            $$ = NewArrayExpression($1.Position(), $3)
        }
        | composite_type LC expression_list_or_nil COMMA RC
        {
            $$ = NewArrayExpression($1.Position(), $3)
        }
        | IDENTIFIER
        {
            $$ = createIdentifierExpression($1.Lit, $1.Position());
        }
        | primary_expression DOT IDENTIFIER
        {
            $$ = createMemberExpression($1, $3.Lit)
        }
        | primary_expression LB expression RB
        {
            $$ = createIndexExpression($1, $3, $1.Position())
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
        ;
expression_list_or_nil
        :
        {
            $$ = nil
        }
        | expression_list
        ;
expression_list
        : expression
        {
            $$ = []Expression{$1}
        }
        | expression_list_or_nil COMMA expression
        {
            $$ = append($1, $3)
        }
        ;
statement
        : simple_statement
        | if_statement
        | for_statement
        | return_statement
        | break_statement
        | continue_statement
        | declaration_statement
        ;
simple_statement_or_nil
        :
        {
            $$ = nil
        }
        | simple_statement
        ;
simple_statement
        : expression
        {
            $$ = &ExpressionStatement{expression: $1}
            $$.SetPosition($1.Position())
        }
        | assign_statement
        ;
if_statement
        : IF expression block
        {
            $$ = &IfStatement{condition: $2, thenBlock: $3, elifList: []*ElseIf{}, elseBlock: nil}
            $$.SetPosition($1.Position())
        }
        | IF expression block ELSE block
        {
            $$ = &IfStatement{condition: $2, thenBlock: $3, elifList: []*ElseIf{}, elseBlock: $5}
            $$.SetPosition($1.Position())
        }
        | IF expression block else_if
        {
            $$ = &IfStatement{condition: $2, thenBlock: $3, elifList: $4, elseBlock: nil}
            $$.SetPosition($1.Position())
        }
        | IF expression block else_if ELSE block
        {
            $$ = &IfStatement{condition: $2, thenBlock: $3, elifList: $4, elseBlock: $6}
            $$.SetPosition($1.Position())
        }
        ;
else_if
        : ELSE IF expression block
        {
            $$ = []*ElseIf{&ElseIf{condition: $3, block: $4}}
        }
        | else_if ELSE IF expression block
        {
            $$ = append($1, &ElseIf{condition: $4, block: $5})
        }
        ;
for_statement
        : FOR LP simple_statement_or_nil SEMICOLON expression_or_nil SEMICOLON simple_statement_or_nil RP block
        {
            $$ = &ForStatement{init: $3, condition: $5, post: $7, block: $9}
            $$.SetPosition($1.Position())
            $9.parent = &StatementBlockInfo{statement: $$}
        }
        ;
expression_or_nil
        :
        {
            $$ = nil
        }
        | expression
        ;
return_statement
        : RETURN expression_list_or_nil
        {
            $$ = NewReturnStatement($1.Position(), $2)
        }
        ;
break_statement
        : BREAK
        {
            $$ = &BreakStatement{}
            $$.SetPosition($1.Position())
        }
        ;
continue_statement
        : CONTINUE
        {
            $$ = &ContinueStatement{}
            $$.SetPosition($1.Position())
        }
        ;
declaration_statement
        : VAR IDENTIFIER type_specifier
        {
            $$ = NewDeclaration($1.Position(), $3, $2.Lit, nil)
        }
        | VAR IDENTIFIER type_specifier ASSIGN expression
        {
            $$ = NewDeclaration($1.Position(), $3, $2.Lit, $5)
        }
        ;
assign_statement
        : expression_list ASSIGN expression_list
        {
            $$ = &AssignStatement{left: $1, right: $3}
            $$.SetPosition($2.Position())
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
block_or_nil
        :
        {
            $$ = nil
        }
        | block
        ;
%%
