%{
package parser
%}

%union{
    expression          Expression
    statement           Statement
    block               *Block
    elif                *Elif
    elif_list           []*Elif
    parameter_list      []*Parameter
    argument_list       []Expression
    statement_list      []Statement
    assignment_operator AssignmentOperator
    type_specifier      *TypeSpecifier
    tok Token
}

%token<tok> IF ELSE ELIF FOR RETURN_T BREAK CONTINUE
        LP RP LC RC
        SEMICOLON COMMA
        ASSIGN_T
        LOGICAL_AND LOGICAL_OR
        EQ NE GT GE LT LE
        ADD SUB MUL DIV
        DOUBLE_LITERAL STRING_LITERAL
        TRUE_T FALSE_T
        IDENTIFIER
        EXCLAMATION DOT
        BOOLEAN_T NUMBER_T STRING_T

%type <parameter_list> parameter_list
%type <argument_list> argument_list
%type <expression> expression expression_opt
      assignment_expression logical_and_expression logical_or_expression
      equality_expression relational_expression
      additive_expression multiplicative_expression
      unary_expression postfix_expression primary_expression

%type <statement> statement
      if_statement for_statement
      return_statement break_statement continue_statement
      declaration_statement
%type <statement_list> statement_list
%type <block> block
%type <elif> elif
%type <elif_list> elif_list
%type <assignment_operator> assignment_operator
%type <type_specifier> type_specifier

%%

translation_unit
        : definition_or_statement
        | translation_unit definition_or_statement
        ;
definition_or_statement
        : function_definition
        | statement
        {
            if l, ok := yylex.(*Lexer); ok {
                l.stmts = append(l.stmts, $1)
            }
        }
        ;
type_specifier
        : BOOLEAN_T
        {
            $$ = &TypeSpecifier{basicType: BooleanType}
            $$.SetPosition($1.Position())
        }
        | NUMBER_T
        {
            $$ = &TypeSpecifier{basicType: NumberType}
            $$.SetPosition($1.Position())
        }
        | STRING_T
        {
            $$ = &TypeSpecifier{basicType: StringType}
            $$.SetPosition($1.Position())
        }
        ;
function_definition
        : type_specifier IDENTIFIER LP parameter_list RP block
        {
            if l, ok := yylex.(*Lexer); ok {
                l.functionDefine($1, $2.Lit, $4, $6);
            }
        }
        | type_specifier IDENTIFIER LP RP block
        {
            if l, ok := yylex.(*Lexer); ok {
                l.functionDefine($1, $2.Lit, nil, $5);
            }
        }
        | type_specifier IDENTIFIER LP parameter_list RP SEMICOLON
        {
            if l, ok := yylex.(*Lexer); ok {
                l.functionDefine($1, $2.Lit, $4, nil);
            }
        }
        | type_specifier IDENTIFIER LP RP SEMICOLON
        {
            if l, ok := yylex.(*Lexer); ok {
                l.functionDefine($1, $2.Lit, nil, nil);
            }
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
            $$ = append([]*Parameter{&Parameter{typeSpecifier: $3, name: $4.Lit}}, $1...)
        }
        ;
argument_list
        : assignment_expression
        {
            $$ = []Expression{$1}
        }
        | argument_list COMMA assignment_expression
        {
            $$ = append([]Expression{$3}, $1...)
        }
        ;
statement_list
        : statement
        {
            $$ = []Statement{$1}
        }
        | statement_list statement
        {
            $$ = append([]Statement{$2}, $1...)
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
        | postfix_expression assignment_operator assignment_expression
        {
            $$ = &AssignExpression{left: $1, operator: $2, operand: $3}
            $$.SetPosition($1.Position())
        }
        ;
assignment_operator
        : ASSIGN_T
        {
            $$ = NormalAssign;
        }
        ;
logical_or_expression
        : logical_and_expression
        | logical_or_expression LOGICAL_OR logical_and_expression
        {
            $$ = &BinaryExpression{operator: LOGICAL_OR_EXPRESSION, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
logical_and_expression
        : equality_expression
        | logical_and_expression LOGICAL_AND equality_expression
        {
            $$ = &BinaryExpression{operator: LOGICAL_AND_EXPRESSION, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
equality_expression
        : relational_expression
        | equality_expression EQ relational_expression
        {
            $$ = &BinaryExpression{operator: EQ_EXPRESSION, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        | equality_expression NE relational_expression
        {
            $$ = &BinaryExpression{operator: NE_EXPRESSION, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
relational_expression
        : additive_expression
        | relational_expression GT additive_expression
        {
            $$ = &BinaryExpression{operator: GT_EXPRESSION, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        | relational_expression GE additive_expression
        {
            $$ = &BinaryExpression{operator: GE_EXPRESSION, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        | relational_expression LT additive_expression
        {
            $$ = &BinaryExpression{operator: LT_EXPRESSION, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        | relational_expression LE additive_expression
        {
            $$ = &BinaryExpression{operator: LE_EXPRESSION, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
additive_expression
        : multiplicative_expression
        | additive_expression ADD multiplicative_expression
        {
            $$ = &BinaryExpression{operator: ADD_EXPRESSION, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        | additive_expression SUB multiplicative_expression
        {
            $$ = &BinaryExpression{operator: SUB_EXPRESSION, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
multiplicative_expression
        : unary_expression
        | multiplicative_expression MUL unary_expression
        {
            $$ = &BinaryExpression{operator: MUL_EXPRESSION, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        | multiplicative_expression DIV unary_expression
        {
            $$ = &BinaryExpression{operator: DIV_EXPRESSION, left: $1, right: $3}
            $$.SetPosition($1.Position())
        }
        ;
unary_expression
        : postfix_expression
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
postfix_expression
        : primary_expression
        | postfix_expression LP argument_list RP
        {
            $$ = &FunctionCallExpression{function: $1, argument: $3}
            $$.SetPosition($1.Position())
        }
        | postfix_expression LP RP
        {
            $$ = &FunctionCallExpression{function: $1, argument: nil}
            $$.SetPosition($1.Position())
        }
        ;
primary_expression
        : LP expression RP
        {
            $$ = $2;
        }
        | IDENTIFIER
        {
            $$ = &IdentifierExpression{name: $1.Lit}
            $$.SetPosition($1.Position())
        }
        | DOUBLE_LITERAL
        {
            $$ = &BooleanExpression{booleanValue: BooleanTrue}
            $$.SetPosition($1.Position())
        }
        | STRING_LITERAL
        {
            $$ = &BooleanExpression{booleanValue: BooleanTrue}
            $$.SetPosition($1.Position())
        }
        | TRUE_T
        {
            $$ = &BooleanExpression{booleanValue: BooleanTrue}
            $$.SetPosition($1.Position())
        }
        | FALSE_T
        {
            $$ = &BooleanExpression{booleanValue: BooleanFalse}
            $$.SetPosition($1.Position())
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
        : IF LP expression RP block
        {
            $$ = &IfStatement{condition: $3, thenBlock: $5, elifList: nil, elseBlock: nil}
            $$.SetPosition($1.Position())
        }
        | IF LP expression RP block ELSE block
        {
            $$ = &IfStatement{condition: $3, thenBlock: $5, elifList: nil, elseBlock: $7}
            $$.SetPosition($1.Position())
        }
        | IF LP expression RP block elif_list
        {
            $$ = &IfStatement{condition: $3, thenBlock: $5, elifList: $6, elseBlock: nil}
            $$.SetPosition($1.Position())
        }
        | IF LP expression RP block elif_list ELSE block
        {
            $$ = &IfStatement{condition: $3, thenBlock: $5, elifList: $6, elseBlock: $8}
            $$.SetPosition($1.Position())
        }
        ;
elif_list
        : elif
        {
            $$ = []*Elif{$1}
        }
        | elif_list elif
        {
            $$ = append([]*Elif{$2}, $1...)
        }
        ;
elif
        : ELIF LP expression RP block
        {
            $$ = &Elif{condition: $3, block: $5}
        }
        ;
for_statement
        : FOR LP expression_opt SEMICOLON expression_opt SEMICOLON
          expression_opt RP block
        {
            $$ = &ForStatement{init: $3, condition: $5, post: $7, block: $9}
            $$.SetPosition($1.Position())
        }
        ;
expression_opt
        : /* empty */
        {
            $$ = nil;
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
            $$ = &DeclarationStatement{typeSpecifier: $1, name: $2.Lit}
            $$.SetPosition($1.Position())
        }
        | type_specifier IDENTIFIER ASSIGN_T expression SEMICOLON
        {
            $$ = &DeclarationStatement{typeSpecifier: $1, name: $2.Lit, initializer: $4}
            $$.SetPosition($1.Position())
        }
        ;
block
        : LC
        {
            if l, ok := yylex.(*Lexer); ok {
                l.currentBlock = &Block{outerBlock: l.currentBlock}
                $<block>$ = l.currentBlock
            }
        }
          statement_list RC
        {

            currentBlock := $<block>2
            currentBlock.statementList = $3
            if l, ok := yylex.(*Lexer); ok {
                l.currentBlock = currentBlock.outerBlock
                $<block>$ = l.currentBlock
            }
        }
        | LC RC
        {
            if l, ok := yylex.(*Lexer); ok {
                $<block>$ = &Block{outerBlock: l.currentBlock}
            }
        }
        ;
%%
