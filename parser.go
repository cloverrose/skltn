package skltn

import (
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
)

type Method struct {
	ReceiverTyp string
	MethodName string
	Args []Arg
	Returns []string
}

type Arg struct {
	Name string
	Typ string
}

type parser struct {
	s   *scanner.Scanner
	src []byte
}

func newParser(src []byte) *parser {
	var (
		fs = token.NewFileSet()
		s  = &scanner.Scanner{}
	)
	// currently the only scanner.Mode option is to include
	// comments (and a private option for not including semicolons).
	// Selecting 0 to *not* scan comments.
	bareEssentials := scanner.Mode(0)
	s.Init(fs.AddFile("", fs.Base(), len(src)), src, nil, bareEssentials)
	return &parser{s: s, src: src}
}

// clone returns new parser with same src
func (p *parser) clone() *parser {
	return newParser(p.src)
}

// it try 2 patterns
// 1. single return
// 2. multiple return
func (p *parser) parse() (*Method, error) {
	method, err := p.parseSingleReturn()
	if err == nil {
		return method, nil
	}

	p = p.clone()
	method, err = p.parseMultiReturn()
	if err != nil {
		return nil, err
	}
	return method, nil
}

func (p *parser) parseBeforeReturn() (*Method, error) {
	receiver, err := p.parseReceiverName()
	if err != nil {
		return nil, err
	}
	methodName, err := p.lookForMethodName()
	if err != nil {
		return nil, err
	}
	args, err := p.parseArgs()
	if err != nil {
		return nil, err
	}
	return &Method{
		ReceiverTyp: receiver,
		MethodName:  methodName,
		Args:        args,
	}, nil
}

func (p *parser) parseSingleReturn() (*Method, error) {
	meethod, err := p.parseBeforeReturn()
	if err != nil {
		return nil, err
	}
	returns, err := p.parseReturnsWithoutParen()
	if err != nil {
		return nil, err
	}
	meethod.Returns = returns
	return meethod, nil
}

func (p *parser) parseMultiReturn() (*Method, error) {
	method, err := p.parseBeforeReturn()
	if err != nil {
		return nil, err
	}
	returns, err := p.parseReturnsWithParen()
	if err != nil {
		return nil, err
	}
	method.Returns = returns
	return method, nil
}

// func (d *db)
func (p *parser) parseReceiverName() (string, error) {
	tempErr := errors.New("unable to find receiver name")
	tok, _ := p.scan(1)
	if tok != token.FUNC {
		return "", tempErr
	}
	tok, _ = p.scan(2)
	if tok != token.LPAREN {
		return "", tempErr
	}
	tok, lit := p.scan(3)
	if tok != token.IDENT {
		return "", tempErr
	}
	tok, lit = p.scan(4)
	if tok == token.MUL {
		tok, lit = p.scan(5)
		if tok != token.IDENT {
			return "", tempErr
		}
	} else if tok != token.IDENT {
		return "", tempErr
	}
	tok, _ = p.scan(5)
	if tok != token.RPAREN {
		return "", errors.New("does not close")
	}
	return lit, nil
}

// Update
func (p *parser) lookForMethodName() (string, error) {
	tempErr := errors.New("unable to find method name")
	tok, lit := p.scan(1)
	if tok != token.IDENT {
		return "", tempErr
	}
	return lit, nil
}

// (ctx context.Context, userID data.UserID, value data.UserDTO)
func (p *parser) parseArgs() ([]Arg, error) {
	tok, _ := p.scan(10)
	if tok != token.LPAREN {
		return nil, errors.New("does not begin")
	}

	args := []Arg{}
	for {
		tok, argName := p.scan(11)
		switch tok {
		case token.RPAREN:
			return args, nil
		case token.IDENT:
			argTyp, cont, err := p.parseArgType()
			if err != nil {
				return nil, err
			}
			args = append(args, Arg{argName, argTyp})
			if !cont {
				return args, nil
			}
		default:
			return nil, errors.New(fmt.Sprintf("unexpected tokens: tok=[%v], lit=[%v]", tok, argName))
		}
	}

	return args, nil
}

// *[]a.B
// []*a.B
// a.B
// a
func (p *parser) parseArgType() (string, bool, error) {
	ret := ""
	for {
		tok, lit := p.scan(21)
		switch tok {
		case token.MUL:
			ret += "*"
		case token.LBRACK:
			ret += "["
		case token.RBRACK:
			ret += "]"
		case token.PERIOD:
			ret += "."
		case token.IDENT:
			ret += lit
		case token.ELLIPSIS:
			ret += "..."
		case token.COMMA:
			return ret, true, nil
		case token.RPAREN:
			return ret, false, nil
		default:
			return "", false, errors.New(fmt.Sprintf("unexpected token: tok=[%v], lit=[%v]", tok, lit))
		}
	}
}

func (p *parser) parseReturnType() (string, bool, error) {
	ret := ""
	for {
		tok, lit := p.scan(21)
		switch tok {
		case token.MUL:
			ret += "*"
		case token.LBRACK:
			ret += "["
		case token.RBRACK:
			ret += "]"
		case token.PERIOD:
			ret += "."
		case token.IDENT:
			ret += lit
		case token.COMMA:
			return ret, true, nil
		case token.RPAREN:
			return ret, false, nil
		default:
			return "", false, errors.New(fmt.Sprintf("unexpected token: tok=[%v], lit=[%v]", tok, lit))
		}
	}
}

func (p *parser) parseReturnTypeWithoutParen() (string, bool, error) {
	ret := ""
	for {
		tok, lit := p.scan(21)
		switch tok {
		case token.MUL:
			ret += "*"
		case token.LBRACK:
			ret += "["
		case token.RBRACK:
			ret += "]"
		case token.PERIOD:
			ret += "."
		case token.IDENT:
			ret += lit
		case token.LBRACE, token.EOF:
			return ret, false, nil
		default:
			return "", false, errors.New(fmt.Sprintf("unexpected token: tok=[%v], lit=[%v]", tok, lit))
		}
	}
}

// (string, error)
func (p *parser) parseReturnsWithParen() ([]string, error) {
	returns := []string{}
	tok, _ := p.scan(30)
	if tok != token.LPAREN {
		return nil, errors.New("does not begin")
	}

	for {
		retTyp, cont, err := p.parseReturnType()
		if err != nil {
			return nil, err
		}
		returns = append(returns, retTyp)
		if !cont {
			return returns, nil
		}
	}
	return returns, nil
}

// error
func (p *parser) parseReturnsWithoutParen() ([]string, error) {
	returns := []string{}
	retTyp, _, err := p.parseReturnTypeWithoutParen()
	if err != nil {
		return nil, err
	}
	returns = append(returns, retTyp)
	return returns, nil
}

// scan reads the next token to scan.
// The input argument only serves as a code-base location for debugging, if
// necessary.
func (p *parser) scan(_ int) (token.Token, string) {
	_, tok, lit := p.s.Scan()
	// NOTE: this is a good spot to do fmt.Println debugging if needed
	return tok, lit
}
