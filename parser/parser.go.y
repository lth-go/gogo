%{
package parser
%}

%union{
    identifier          string
    expression          Expression
    statement           Statement
    block               *Block
    elif_list           []*Elif
    parameter_list      []*Parameter
    argument_list       []*Expression
    statement_list      []*Statement
    assignment_operator AssignmentOperator
    type_specifier      VM_BasicType
}

%token <expression>     DOUBLE_LITERAL
%token <expression>     STRING_LITERAL
%token <identifier>     IDENTIFIER
%token IF ELSE ELIF WHILE FOR RETURN_T BREAK CONTINUE
        LP RP LC RC SEMICOLON COLON COMMA ASSIGN_T LOGICAL_AND LOGICAL_OR
        EQ NE GT GE LT LE ADD SUB MUL DIV MOD TRUE_T FALSE_T EXCLAMATION DOT
        ADD_ASSIGN_T SUB_ASSIGN_T MUL_ASSIGN_T DIV_ASSIGN_T MOD_ASSIGN_T
        BOOLEAN_T INT_T DOUBLE_T STRING_T

%type <parameter_list> parameter_list
%type <argument_list> argument_list
%type <expression> expression expression_opt
      assignment_expression logical_and_expression logical_or_expression
      equality_expression relational_expression
      additive_expression multiplicative_expression
      unary_expression postfix_expression primary_expression
%type <statement> statement
      if_statement while_statement for_statement
      return_statement break_statement continue_statement
      declaration_statement
%type <statement_list> statement_list
%type <block> block
%type <elif_list> elif elif_list
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
            $$ = DVM_BOOLEAN_TYPE;
        }
        | INT_T
        {
            $$ = DVM_INT_TYPE;
        }
        | DOUBLE_T
        {
            $$ = DVM_DOUBLE_TYPE;
        }
        | STRING_T
        {
            $$ = DVM_STRING_TYPE;
        }
        ;
function_definition
        : type_specifier IDENTIFIER LP parameter_list RP block
        {
            if l, ok := yylex.(*Lexer); ok {
                l.function_define($1, $2, $4, $6);
            }
        }
        | type_specifier IDENTIFIER LP RP block
        {
            if l, ok := yylex.(*Lexer); ok {
                l.function_define($1, $2, NULL, $5);
            }
        }
        | type_specifier IDENTIFIER LP parameter_list RP SEMICOLON
        {
            if l, ok := yylex.(*Lexer); ok {
                l.function_define($1, $2, $4, NULL);
            }
        }
        | type_specifier IDENTIFIER LP RP SEMICOLON
        {
            if l, ok := yylex.(*Lexer); ok {
                l.function_define($1, $2, NULL, NULL);
            }
        }
        ;
parameter_list
        : type_specifier IDENTIFIER
        {
            $$ = []Parameter{{Type: $1, Name: $2}}
        }
        | parameter_list COMMA type_specifier IDENTIFIER
        {
            $$ = append([]Parameter{{Type: $3, Name: $4}}, $1)
        }
        ;
argument_list
        : assignment_expression
        {
            $$ = []Expression{$1}
        }
        | argument_list COMMA assignment_expression
        {
            $$ = append([]Expression{$3}, $1)
        }
        ;
statement_list
        : statement
        {
            $$ = []Statement{$1}
        }
        | statement_list statement
        {
            $$ = append([]Statement{$2}, $1)
        }
        ;
expression
        : assignment_expression
        | expression COMMA assignment_expression
        {
            $$ = CommaExpression{left: $1, right: $3}
        }
        ;
assignment_expression
        : logical_or_expression
        | postfix_expression assignment_operator assignment_expression
        {
            $$ = AssignExpression{left: $1, operator: $3, operand: $3}
        }
        ;
assignment_operator
        : ASSIGN_T
        {
            $$ = NORMAL_ASSIGN;
        }
        | ADD_ASSIGN_T
        {
            $$ = ADD_ASSIGN;
        }
        | SUB_ASSIGN_T
        {
            $$ = SUB_ASSIGN;
        }
        | MUL_ASSIGN_T
        {
            $$ = MUL_ASSIGN;
        }
        | DIV_ASSIGN_T
        {
            $$ = DIV_ASSIGN;
        }
        | MOD_ASSIGN_T
        {
            $$ = MOD_ASSIGN;
        }
        ;
logical_or_expression
        : logical_and_expression
        | logical_or_expression LOGICAL_OR logical_and_expression
        {
            $$ = BinaryExpression{operator: LOGICAL_OR_EXPRESSION, left: $1, right: $3}
        }
        ;
logical_and_expression
        : equality_expression
        | logical_and_expression LOGICAL_AND equality_expression
        {
            $$ = BinaryExpression{operator: LOGICAL_AND_EXPRESSION, left: $1, right: $3}
        }
        ;
equality_expression
        : relational_expression
        | equality_expression EQ relational_expression
        {
            $$ = BinaryExpression{operator: EQ_EXPRESSION, left: $1, right: $3}
        }
        | equality_expression NE relational_expression
        {
            $$ = BinaryExpression{operator: NE_EXPRESSION, left: $1, right: $3}
        }
        ;
relational_expression
        : additive_expression
        | relational_expression GT additive_expression
        {
            $$ = BinaryExpression{operator: GT_EXPRESSION, left: $1, right: $3}
        }
        | relational_expression GE additive_expression
        {
            $$ = BinaryExpression{operator: GE_EXPRESSION, left: $1, right: $3}
        }
        | relational_expression LT additive_expression
        {
            $$ = BinaryExpression{operator: LT_EXPRESSION, left: $1, right: $3}
        }
        | relational_expression LE additive_expression
        {
            $$ = BinaryExpression{operator: LE_EXPRESSION, left: $1, right: $3}
        }
        ;
additive_expression
        : multiplicative_expression
        | additive_expression ADD multiplicative_expression
        {
            $$ = BinaryExpression{operator: ADD_EXPRESSION, left: $1, right: $3}
        }
        | additive_expression SUB multiplicative_expression
        {
            $$ = BinaryExpression{operator: SUB_EXPRESSION, left: $1, right: $3}
        }
        ;
multiplicative_expression
        : unary_expression
        | multiplicative_expression MUL unary_expression
        {
            $$ = BinaryExpression{operator: MUL_EXPRESSION, left: $1, right: $3}
        }
        | multiplicative_expression DIV unary_expression
        {
            $$ = BinaryExpression{operator: DIV_EXPRESSION, left: $1, right: $3}
        }
        | multiplicative_expression MOD unary_expression
        {
            $$ = BinaryExpression{operator: MOD_EXPRESSION, left: $1, right: $3}
        }
        ;
unary_expression
        : postfix_expression
        | SUB unary_expression
        {
            $$ = MinusExpression{operand: $2}
        }
        | EXCLAMATION unary_expression
        {
            $$ = LogicalNotExpression{operand: $2}
        }
        ;
postfix_expression
        : primary_expression
        | postfix_expression LP argument_list RP
        {
            $$ = FunctionCallExpression{function: $1, argument: $3}
        }
        | postfix_expression LP RP
        {
            $$ = FunctionCallExpression{function: $1, argument: nil}
        }
        ;
primary_expression
        : LP expression RP
        {
            $$ = $2;
        }
        | IDENTIFIER
        {
            $$ = IdentifierExpression{name: $1}
        }
        | INT_LITERAL
        | DOUBLE_LITERAL
        | STRING_LITERAL
        | REGEXP_LITERAL
        | TRUE_T
        {
            $$ = BooleanExpression{boolean_value: DVM_TRUE}
        }
        | FALSE_T
        {
            $$ = BooleanExpression{boolean_value: DVM_FALSE}
        }
        ;
statement
        : expression SEMICOLON
        {
          $$ = ExpressionStatement{expression_s: $1}
        }
        | if_statement
        | while_statement
        | for_statement
        | foreach_statement
        | return_statement
        | break_statement
        | continue_statement
        | try_statement
        | throw_statement
        | declaration_statement
        ;
if_statement
        : IF LP expression RP block
        {
            $$ = IfStatement{condition: $3, then_block: $5, elif_list: nil, else_block: nil}
        }
        | IF LP expression RP block ELSE block
        {
            $$ = IfStatement{condition: $3, then_block: $5, elif_list: nil, else_block: $7}
        }
        | IF LP expression RP block elif_list
        {
            $$ = IfStatement{condition: $3, then_block: $5, elif_list: $6, else_block: nil}
        }
        | IF LP expression RP block elif_list ELSE block
        {
            $$ = IfStatement{condition: $3, then_block: $5, elif_list: $6, else_block: $8}
        }
        ;
elif_list
        : elif
        {
            $$ = []ElifList{$2}
        }
        | elif_list elif
        {
            $$ = append($2, $1)
        }
        ;
elif
        : ELIF LP expression RP block
        {
            $$ = []Elif{{condition: $3, block: $5}}
        }
        ;
while_statement
        : WHILE LP expression RP block
        {
            $$ = WhileStatement{condition: $3, block: $5}
        }
        ;
for_statement
        : FOR LP expression_opt SEMICOLON expression_opt SEMICOLON
          expression_opt RP block
        {
            $$ = ForStatement{init: $3, condition: $5, post: $7, block: $9}
        }
        ;
expression_opt
        : /* empty */
        {
            $$ = NULL;
        }
        | expression
        ;
return_statement
        : RETURN_T expression_opt SEMICOLON
        {
            $$ = ReturnStatement{return_value: $2};
        }
        ;
break_statement 
        : BREAK SEMICOLON
        {
            $$ = BreakStatement{}
        }
        ;
continue_statement
        : CONTINUE SEMICOLON
        {
            $$ = ContinueStatement{}
        }
        ;
declaration_statement
        : type_specifier IDENTIFIER SEMICOLON
        {
            $$ = DeclarationStatement{Type: $1, Name: $2, initializer: nil}
        }
        | type_specifier IDENTIFIER ASSIGN_T expression SEMICOLON
        {
            $$ = DeclarationStatement{type: $1, name: $2, initializer: $4}
        }
        ;
block
        : LC
        {
            if l, ok := yylex.(*Lexer); ok {
                $<block>$ = l.current_block = Block{outer_block: l.current_block}
            }
        }
          statement_list RC
        {

            current_block = $<block>2
            current_block.statement_list = $3
            if l, ok := yylex.(*Lexer); ok {
                $<block>$ = l.current_block = current_block.outer_block
            }
        }
        | LC RC
        {
            if l, ok := yylex.(*Lexer); ok {
                $<block>$ := Block{outer_block: l.current_block, statement_list: nil}
            }
        }
        ;
%%
