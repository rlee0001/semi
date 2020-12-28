grammar Semi;

MUL: '*';
DIV: '/';
ADD: '+';
SUB: '-';
EQUALS: '=';
IS_EQUAL_TO: '==';

PAREN_OPEN: '(';
PAREN_CLOSE: ')';
BRACE_OPEN: '{';
BRACE_CLOSE: '}';
SEMI: ';';
COMMA: ',';
COLON: ':';

OUTPUT: 'output';
VAR: 'var';
FUN: 'fun';

TRUE: 'true';
FALSE: 'false';

IDENT: [a-zA-Z_$][a-zA-Z_$0-9?]*;
INTEGER: [-]?[0-9]+;
FLOAT: [-]?[0-9]+[.][0-9]+;
WS: [ \r\n\t]+;

module
    : WS? moduleScope WS? EOF
    ;

moduleScope
    : moduleScopeItems?
    ;

moduleScopeItems
    : moduleScopeItem
    | moduleScopeItem WS? moduleScopeItems
    ;

moduleScopeItem
    : functionDeclaration
    ;

functionDeclaration
    : FUN WS ident=IDENT WS? parameterDeclarationList WS? COLON WS? typeExpression WS? block SEMI
    ;

parameterDeclarationList
    : PAREN_OPEN WS? parameterDeclarations? WS? PAREN_CLOSE
    ;

parameterDeclarations
    : parameterDeclaration
    | parameterDeclaration WS? COMMA WS? parameterDeclarations
    ;

parameterDeclaration
    : ident=IDENT WS? COLON WS? typeExpression
    ;

typeExpression
    : ident=IDENT
    ;

block
    : BRACE_OPEN WS? functionScope WS? BRACE_CLOSE
    ;

functionScope
    : functionScopeItems?
    ;

functionScopeItems
    : functionScopeItem
    | functionScopeItem WS? functionScopeItems
    ;

functionScopeItem
    : statement
    | localDeclaration
    | block
    ;

statement
    : OUTPUT WS expression SEMI
    ;

localDeclaration
    : VAR WS ident=IDENT WS? EQUALS WS? expression SEMI
    ;

expression
    : (SUB WS)? ident=IDENT                                                         # IdentifierFactor
    | (SUB WS)? integer=INTEGER                                                     # IntegerFactor
    | (SUB WS)? float=FLOAT                                                         # FloatFactor
    | (SUB WS)? TRUE                                                                # TrueFactor
    | (SUB WS)? FALSE                                                               # FalseFactor
    | (SUB WS)? ident=IDENT WS? PAREN_OPEN WS? argumentListItems? WS? PAREN_CLOSE   # CallExpression
    | (SUB WS)? PAREN_OPEN WS? expression WS? PAREN_CLOSE                           # SubExpression
    | expression WS? op=(MUL|DIV) WS? expression                                    # MulDivTerm
    | expression WS? op=(ADD|SUB) WS? expression                                    # AddSubTerm
    | expression WS? IS_EQUAL_TO WS? expression                                     # IsEqualToTerm
    ;

argumentListItems
    : argumentListItem
    | argumentListItem WS? COMMA WS? argumentListItems
    ;

argumentListItem
    : expression
    ;
