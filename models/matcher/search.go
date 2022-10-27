package matcher

import (
	"unicode/utf8"
)

func optimizeQuery(q string) (formattedQ string, apearingLetters uint32) {
	// byte layout of apearingLetters:
	// 0-25 = represends a-z
	// 26   = query contains a space
	// 27   = query contains unicode characters

	resp := []rune{}

	lastAppearingSpaceIdx := 0
	for _, c := range q {
		if c >= utf8.RuneSelf {
			replacementC, garbage := checkAndCorredUnicodeChar(c)
			if garbage {
				continue
			}
			if replacementC == c {
				// Just commit all the unicode letters
				// These are characters we don't care about
				resp = append(resp)
				apearingLetters |= 1 << 27
				continue
			}
			c = replacementC
		}

		// If the letter is uppercase we make it lowercase
		if 'a' <= c && c <= 'z' {
			apearingLetters |= 1 << (c - 'a')
			resp = append(resp, c)
			continue
		} else if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
			apearingLetters |= 1 << (c - 'a')
			resp = append(resp, c)
			continue
		}

		switch c {
		case ' ', '\t', '\n':
			if len(resp) == 0 {
				// Strip leading spaces
				continue
			}

			lastAppearingSpaceIdx = len(resp)
			resp = append(resp, ' ')
			apearingLetters |= 1 << 26
			continue
		}
	}

	if lastAppearingSpaceIdx != 0 && len(q) != 0 && resp[len(resp)-1] == ' ' {
		// Strip trailing spaces
		resp = resp[lastAppearingSpaceIdx:]
	}

	return string(resp), apearingLetters
}

func checkAndCorredUnicodeChar(c rune) (replacementCharacter rune, garbage bool) {
	switch c {
	case 'à', 'À', 'á', 'Á', 'â', 'Â', 'ã', 'Ã', 'ä', 'Ä', 'å', 'Å', 'æ', 'Æ':
		return 'a', false
	case 'è', 'È', 'é', 'É', 'ê', 'Ê', 'ë', 'Ë':
		return 'e', false
	case 'ì', 'Ì', 'í', 'Í', 'î', 'Î', 'ï', 'Ï':
		return 'i', false
	case 'ò', 'Ò', 'ó', 'Ó', 'ô', 'Ô', 'õ', 'Õ', 'ö', 'Ö', 'ð', 'Ð', 'ø', 'Ø':
		return 'o', false
	case 'ù', 'Ù', 'ú', 'Ú', 'û', 'Û', 'ü', 'Ü':
		return 'u', false
	case 'ß':
		return 's', false
	case 'ñ', 'Ñ':
		return 'n', false
	case 'ý', 'Ý', 'ÿ', 'Ÿ':
		return 'y', false
	case 'ç', 'Ç', '©':
		return 'c', false
	case '®':
		return 'r', false
	case 768, // accent of: à
		769, // accent of: á
		770, // accent of: â
		771, // accent of: ã
		776, // accent of: ä
		778, // accent of: å
		'¿',
		'¡',
		0x2002, // En space
		0x2003, // Em space
		0x2004, // Three-per-em space
		0x2005, // Four-per-em space
		0x2006, // Six-per-em space
		0x2007, // Figure space
		0x2008, // Punctuation space
		0x2009, // Thin space
		0x200A, // Hair space
		0x200B, // Zero width space
		0x202F, // Narrow no-break space
		0x205F, // Medium mathematical space
		0x3000, // Ideographic space
		'“',
		'”',
		'’',
		'‵',
		'‹',
		'›',
		'»',
		'«',
		utf8.RuneError:
		return utf8.RuneError, true
	default:
		return c, false
	}
}
