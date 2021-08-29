package routeBuilder

type pathInformation struct {
	AsOpenAPIPath string
	Params        []string
}

func parseFiberPath(p string) pathInformation {
	return (&pathParser{path: []byte(p)}).start()
}

type pathParser struct {
	path []byte
	x    int
}

// c returns the character at x and a boolean to indicate the end of the path
// x is also incremented by one
func (p *pathParser) c() (char byte, eof bool) {
	if p.x >= len(p.path) {
		return 0, true
	}
	char = p.path[p.x]
	p.x++
	return char, false
}

func (p *pathParser) start() pathInformation {
	openAPIPath := []byte{}
	params := []string{}

	for {
		c, eof := p.c()
		if eof {
			break
		}
		switch c {
		case ':':
			paramName := p.param()
			if len(paramName) == 0 {
				openAPIPath = append(openAPIPath, ':')
			} else {
				params = append(params, string(paramName))
				openAPIPath = append(openAPIPath, '{')
				openAPIPath = append(openAPIPath, paramName...)
				openAPIPath = append(openAPIPath, '}')
			}
		case '*':
			panic("FIXME: converting fiber unnamed parameters to open api parameters unsupported")
		case '+':
			panic("FIXME: converting fiber unnamed parameters to open api parameters unsupported")
		default:
			openAPIPath = append(openAPIPath, c)
		}
	}

	return pathInformation{AsOpenAPIPath: string(openAPIPath), Params: params}
}

func (p *pathParser) param() []byte {
	paramName := []byte{}

	for {
		char, eof := p.c()
		if eof {
			break
		}
		if char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z' || char >= '0' && char <= '9' || char == '_' {
			paramName = append(paramName, char)
		} else {
			p.x--
			break
		}
	}

	return paramName
}
