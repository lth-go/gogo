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

package_name         []string
    require_list         []*Require

extends_list         []*Extend
    member_declaration   []MemberDeclaration
    function_definition  *FunctionDefinition

tok                  Token
}

%token<tok> IF ELSE ELIF FOR RETURN_T BREAK CONTINUE
        LP RP LC RC LB RB
        SEMICOLON COMMA COLON
        ASSIGN_T
        LOGICAL_AND LOGICAL_OR
        EQ NE GT GE LT LE
        ADD SUB MUL DIV
        INT_LITERAL DOUBLE_LITERAL STRING_LITERAL TRUE_T FALSE_T
        NULL_T
        IDENTIFIER
        EXCLAMATION DOT
        VOID_T BOOLEAN_T INT_T DOUBLE_T STRING_T
        NEW
        REQUIRE
        CLASS_T THIS_T

%type   <package_name> package_name
%type   <require_list> require_list require_declaration

%type <expression> expression expression_opt
      assignment_expression
      logical_and_expression logical_or_expression
      equality_expression relational_expression
      additive_expression multiplicative_expression
      unary_expression postfix_expression primary_expression primary_no_new_array
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

%type <type_specifier> basic_type_specifier type_specifier class_type_specifier array_type_specifier

%type <array_dimension> dimension_expression
%type <array_dimension_list> dimension_expression_list dimension_list

%type   <extends_list> extends_list extends
%type   <member_declaration> member_declaration member_declaration_list method_member field_member
%type   <function_definition> method_function_definition

%%

translation_unit
        : initial_declaration definition_or_statement
        | translation_unit definition_or_statement
        ;
initial_declaration
        : /* empty */
        {
            setRequireList(nil)
        }
        | require_list
        {
            setRequireList($1)
        }
        ;
require_list
        : require_declaration
        | require_list require_declaration
        {
            $$ = chainRequireList($1, $2)
        }
        ;
require_declaration
        : REQUIRE package_name SEMICOLON
        {
            $$ = createRequireList($2)
        }
        ;
package_name
        : IDENTIFIER
        {
            $$ = createPackageName($1.Lit)
        }
        | package_name DOT IDENTIFIER
        {
            $$ = chainPackageName($1, $3.Lit)
        }
        ;
definition_or_statement
        : function_definition
        | class_definition
        | statement
        {
            l := yylex.(*Lexer)
            l.compiler.statementList = append(l.compiler.statementList, $1)
        }
        ;
basic_type_specifier
        : VOID_T
        {
            $$ = createTypeSpecifier(vm.VoidType, $1.Position())
        }
        | BOOLEAN_T
        {
            $$ = createTypeSpecifier(vm.BooleanType, $1.Position())
        }
        | INT_T
        {
            $$ = createTypeSpecifier(vm.IntType, $1.Position())
        }
        | DOUBLE_T
        {
            $$ = createTypeSpecifier(vm.DoubleType, $1.Position())
        }
        | STRING_T
        {
            $$ = createTypeSpecifier(vm.StringType, $1.Position())
        }
        ;
class_type_specifier
        : IDENTIFIER
        {
            $$ = create_class_type_specifier($1.Lit, $1.Position())
        }
        ;
array_type_specifier
        : basic_type_specifier LB RB
        {
            $$ = create_array_type_specifier($1)
            $$.SetPosition($1.Position())
        }
        | IDENTIFIER LB RB
        {
            class_type := create_class_type_specifier($1.Lit, $1.Position())
            $$ = create_array_type_specifier(class_type)
        }
        | array_type_specifier LB RB
        {
            $$ = create_array_type_specifier($1)
        }
        ;
type_specifier
        : basic_type_specifier
        {
            $$ = $1
        }
        | array_type_specifier
        | class_type_specifier
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
        ;
primary_expression
        : primary_no_new_array
        | array_creation
        | IDENTIFIER
        {
            $$ = createIdentifierExpression($1.Lit, $1.Position());
        }
        ;
primary_no_new_array
        : primary_no_new_array LB expression RB
        {
            $$ = createIndexExpression($1, $3, $1.Position())
        }
        | IDENTIFIER LB expression RB
        {
            identifier := createIdentifierExpression($1.Lit, $1.Position());
            $$ = createIndexExpression(identifier, $3, $1.Position())
        }
        | primary_expression DOT IDENTIFIER
        {
            $$ = createMemberExpression($1, $3.Lit, $1.Position())
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
        | THIS_T
        {
            $$ = createThisExpression($1.Position())
        }
        | NEW IDENTIFIER LP RP
        {
            $$ = createNewExpression($2.Lit, "", nil, $1.Position())
        }
        | NEW IDENTIFIER LP argument_list RP
        {
            $$ = createNewExpression($2.Lit, "", $4, $1.Position())
        }
        | NEW IDENTIFIER DOT IDENTIFIER LP RP
        {
            $$ = createNewExpression($2.Lit, $4.Lit, nil, $1.Position())
        }
        | NEW IDENTIFIER DOT IDENTIFIER LP argument_list RP
        {
            $$ = createNewExpression($2.Lit, $4.Lit, $6, $1.Position())
        }
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
            $$ = createBasicArrayCreation($2, $3, nil, $1.Position())
        }
        | NEW basic_type_specifier dimension_expression_list dimension_list
        {
            $$ = createBasicArrayCreation($2, $3, $4, $1.Position())
        }
        | NEW class_type_specifier dimension_expression_list
        {
            $$ = createClassArrayCreation($2, $3, nil, $1.Position())
        }
        | NEW class_type_specifier dimension_expression_list dimension_list
        {
            $$ = createClassArrayCreation($2, $3, $4, $1.Position())
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
class_definition
        : CLASS_T IDENTIFIER extends LC
        {
            startClassDefine($2.Lit, $3, $1.Position())
        }
          member_declaration_list RC
        {
            endClassDefine($6)
        }
        | CLASS_T IDENTIFIER extends LC
        {
            startClassDefine($2.Lit, $3, $1.Position())
        }
          RC
        {
            endClassDefine(nil)
        }
        ;
extends
        : /* empty */
        {
            $$ = nil;
        }
        | COLON extends_list
        {
            $$ = $2;
        }
        ;
extends_list
        : IDENTIFIER
        {
            $$ = createExtendList($1.Lit)
        }
        | extends_list COMMA IDENTIFIER
        {
            $$ = chainExtendList($1, $3.Lit)
        }
        ;
member_declaration_list
        : member_declaration
        | member_declaration_list member_declaration
        {
            $$ = chainMemberDeclaration($1, $2);
        }
        ;
member_declaration
        : method_member
        | field_member
        ;
method_member
        : method_function_definition
        {
            $$ = createMethodMember($1, $1.typeSpecifier.Position())
        }
        ;
method_function_definition
        : type_specifier IDENTIFIER LP parameter_list RP block
        {
            $$ = methodFunctionDefine($1, $2.Lit, $4, $6);
        }
        | type_specifier IDENTIFIER LP RP block
        {
            $$ = methodFunctionDefine($1, $2.Lit, nil, $5);
        }
        | type_specifier IDENTIFIER LP parameter_list RP SEMICOLON
        {
            $$ = methodFunctionDefine($1, $2.Lit, $4, nil);
        }
        | type_specifier IDENTIFIER LP RP SEMICOLON
        {
            $$ = methodFunctionDefine($1, $2.Lit, nil, nil);
        }
        ;
field_member
        : type_specifier IDENTIFIER SEMICOLON
        {
            $$ = createFieldMember($1, $2.Lit, $1.Position())
        }
        ;
%%
