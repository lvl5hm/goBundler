package jsLoader

import "testing"

var ps parser

func setParser(text string) {
	toks := lex([]byte(text))
	ps = parser{
		tokens: toks,
		i:      0,
	}
}

func TestExpressions(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"a + foo * 32;",
			"a+foo*32;",
		},
		{
			"fee = a**(b+'ds');",
			"fee=a**(b+'ds');",
		},
		{
			"a + b / /[A*^?-Z]/g;",
			"a+b//[A*^?-Z]/g;",
		},
		{
			"a + {foo:bar} * 3;",
			"a+{foo:bar}*3;",
		},
		{
			"!!(a+b);",
			"!!(a+b);",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestObjectLiteral(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"a={a,...foo, ...{bar} = 3};",
			"a={a,...foo,...{bar}=3};",
		},
		{
			"a = {};",
			"a={};",
		},
		{
			"a = {a:b,c,};",
			"a={a:b,c};",
		},
		{
			"a = {a:()=>{},c,};",
			"a={a:()=>{},c};",
		},
		{
			"a = {32: foo, 'bar': bar};",
			"a={32:foo,'bar':bar};",
		},
		{
			"a = {0xff: foo};",
			"a={0xff:foo};",
		},
		{
			"a = {[foo+32]:a};",
			"a={[foo+32]:a};",
		},
		{
			"a = {foo(){}};",
			"a={foo(){}};",
		},
		{
			"a = {get foo(){}, set bar(){}};",
			"a={get foo(){},set bar(){}};",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestLambda(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"foo=>bar;",
			"(foo)=>bar;",
		},
		{
			"()=>bar;",
			"()=>bar;",
		},
		{
			"(a,b,c)=>{bar;};",
			"(a,b,c)=>{bar;};",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestMemberExpression(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"foo[a].b().c;",
			"foo[a].b().c;",
		},
		{
			"new a.b().c;",
			"new a.b().c;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestFunctionDeclaration(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"function foo() {}",
			"function foo(){}",
		},
		{
			"function foo(foo = ee = 321, bar) {}",
			"function foo(foo=ee=321,bar){}",
		},
		{
			"function foo(...{}) {}",
			"function foo(...{}){}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"[foo, bar, 213*(21+3), () => foo,];",
			"[foo,bar,213*(21+3),()=>foo];",
		},
		{
			"[foo, , , bar,];",
			"[foo,,,bar];",
		},
		{
			"[foo, ...bar];",
			"[foo,...bar];",
		},
		{
			"[foo(a, b, c, d), bar()];",
			"[foo(a,b,c,d),bar()];",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestBlockStatement(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"{foo; bar = 321;;}",
			"{foo;bar=321;;}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestForStatement(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"for (;;);",
			"for(;;);",
		},
		{
			"for(var i=0;i<10;i++);",
			"for(var i=0;i<10;i++);",
		},
		{
			"for(;i<10;i++);",
			"for(;i<10;i++);",
		},
		{
			"for(i;;i++);",
			"for(i;;i++);",
		},
		{
			"for(i;i<23;);",
			"for(i;i<23;);",
		},
		{
			"for(i;i<23;) {}",
			"for(i;i<23;){}",
		},
		{
			"for(i;i<23;) foo = 3;",
			"for(i;i<23;)foo=3;",
		},
		{
			"for(foo of bar()) foo = 3;",
			"for(foo of bar())foo=3;",
		},
		{
			"for(const foo in bar) {foo = 3;}",
			"for(const foo in bar){foo=3;}",
		},
		{
			"for(a in b; i < 21; i++) {foo = 3;}",
			"for(a in b;i<21;i++){foo=3;}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestInAndInstanceof(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"for(a in b; i < 21; i++) {foo = 3;}",
			"for(a in b;i<21;i++){foo=3;}",
		},
		{
			"a instanceof b;",
			"a instanceof b;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestStringLiteral(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"'foo \\' + fsbds';",
			"'foo \\' + fsbds';",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestWhileStatement(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"while (foo);",
			"while(foo);",
		},
		{
			"while (foo*bar < 3) {}",
			"while(foo*bar<3){}",
		},
		{
			"while (foo, bar += 3) bar();",
			"while(foo,bar+=3)bar();",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestDoWhileStatement(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"do {} while(foo);",
			"do {}while(foo);",
		},
		{
			"do ; while (foo*bar < 3);",
			"do ;while(foo*bar<3);",
		},
		{
			"do bar();while(foo,bar+=3);",
			"do bar();while(foo,bar+=3);",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestIfStatement(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"if (foo) bar;",
			"if(foo)bar;",
		},
		{
			"if(foo){} else bar;",
			"if(foo){}else bar;",
		},
		{
			"if(foo, bar = 3){foo();}",
			"if(foo,bar=3){foo();}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestExpressionStatement(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"var foo = 3, bar;",
			"var foo=3,bar;",
		},
		{
			"break;",
			"break;",
		},
		{
			"continue;",
			"continue;",
		},
		{
			"debugger;",
			"debugger;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestImportStatement(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"import 'foo';",
			"import'foo';",
		},
		{
			"import * as foo from 'foo';",
			"import*as foo from'foo';",
		},
		{
			"import bar, * as foo from 'foo';",
			"import bar,*as foo from'foo';",
		},
		{
			"import bar, {foo as bar} from 'foo';",
			"import bar,{foo as bar} from'foo';",
		},
		{
			"import foo, {default as foo, bar, baz} from 'foo';",
			"import foo,{default as foo,bar as bar,baz as baz} from'foo';",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestExportStatement(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"export default foo;",
			"export default foo;",
		},
		{
			"export default foo + 231;",
			"export default foo+231;",
		},
		{
			"export default function() {};",
			"export default function(){};",
		},
		{
			"export default function foo() {};",
			"function foo(){}export default foo;",
		},
		{
			"export var foo = 4, bar;",
			"var foo=4,bar;export{foo as foo,bar as bar};",
		},
		{
			"export {};",
			"export{};",
		},
		{
			"export {foo as fee, bar as default, wee, };",
			"export{foo as fee,bar as default,wee as wee};",
		},
		{
			"export {} from 'foo';",
			"export{} from'foo';",
		},
		{
			"export function foo() {};",
			"function foo(){}export{foo as foo};;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestObjectPattern(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"var {} = foo;",
			"var {}=foo;",
		},
		{
			"var {a:b} = foo;",
			"var {a:b}=foo;",
		},
		{
			"var {a:b=32} = foo;",
			"var {a:b=32}=foo;",
		},
		{
			"var {a:b={c:f=5,d:e}=3} = foo;",
			"var {a:b={c:f=5,d:e}=3}=foo;",
		},
		{
			"({} = foo);",
			"({}=foo);",
		},
		{
			"var {foo,...bar}=doo;",
			"var {foo,...bar}=doo;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestReturnStatement(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			`return {
				result: mapResult,
				keyPrefix: keyPrefix,
				func: mapFunction,
				context: mapContext,
				count: 0
			};`,
			"return {result:mapResult,keyPrefix:keyPrefix," +
				"func:mapFunction,context:mapContext,count:0};",
		},
	}
	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestConditionalExpression(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"foo?bar:baz;",
			"foo?bar:baz;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestArrayPattern(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"[,,] = foo;",
			"[,]=foo;",
		},
		{
			"[foo,,bar] = a;",
			"[foo,,bar]=a;",
		},
		{
			"[a = 23, foo]=a;",
			"[a=23,foo]=a;",
		},
		{
			"[{foo:bar = 23} = 23, foo]=a;",
			"[{foo:bar=23}=23,foo]=a;",
		},
		{
			"[a, ...b] = 32;",
			"[a,...b]=32;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestClass(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"class foo{}",
			"class foo{}",
		},
		{
			"class foo extends bar{}",
			"class foo extends bar{}",
		},
		{
			"class foo{bar:3;}",
			"class foo{bar=3;}",
		},
		{
			"class foo{23=12;['ffp']:321;}",
			"class foo{23=12;['ffp']=321;}",
		},
		{
			"class foo{[foo](){}}",
			"class foo{[foo](){}}",
		},
		{
			"class foo{get [foo](){}}",
			"class foo{get [foo](){}}",
		},
		{
			"class foo{static bar;}",
			"class foo{static bar;}",
		},
		{
			"a = class foo{};",
			"a=class foo{};",
		},
		{
			"a = class{};",
			"a=class{};",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestNewlineAndSemi(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"var\nfoo\n",
			"var foo;",
		},
		{
			"{foo}",
			"{foo;}",
		},
		{
			"var a = {foo}",
			"var a={foo};",
		},
		{
			`const {
				addEntityType,

				addEntity,
			} = engine`,
			"const {addEntityType,addEntity}=engine;",
		},
		{
			`for(
				i
				;
				i<321;
				i++
				);`,
			"for(i;i<321;i++);",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := program(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}
