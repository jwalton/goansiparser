package ansiparser

import "github.com/rivo/uniseg"

// This is the "slow" parser which is used when parseAscii finds unicode characters.
func parseUnicode(str string, startFG string, startBG string) (tokens []AnsiToken, consumed int) {
	tokens = make([]AnsiToken, 0, 1)

	currentFG := startFG
	currentBG := startBG
	i := 0
	currentStringStart := -1

	// Called when we get to the end of a string of ascii characters.
	// This passes the string of ascii charaters through the existing
	// `parseASCII()` function and appends all the resulting tokens.
	endString := func() {
		if currentStringStart != -1 {
			substr := str[currentStringStart:i]
			substrTokens, _ := parseASCII(substr, true, currentFG, currentBG)
			tokens = append(tokens, substrTokens...)

			if len(tokens) > 0 {
				lastToken := tokens[len(tokens)-1]
				// Copy the FG and BG colors from the last token generated by parseASCII.
				currentFG = lastToken.FG
				currentBG = lastToken.BG
			}
		}
		currentStringStart = -1
	}

	gr := uniseg.NewGraphemes(str)
	for gr.Next() {
		charBytes := gr.Bytes()
		charLen := len(charBytes)
		if charLen == 1 {
			// Add the character to the current string.
			if currentStringStart == -1 {
				currentStringStart = i
			}
		} else {
			endString()
			// Add a token for the multi-byte character.
			tokens = append(tokens, AnsiToken{
				Type:    ComplexChar,
				Content: string(charBytes),
				FG:      currentFG,
				BG:      currentBG,
			})
		}
		i += charLen
	}

	endString()

	return tokens, i
}
