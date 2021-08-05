package strtr

type Tr struct {
	tr map[rune]rune
}

func NewTr(spec *Spec) *Tr {
	tr := make(map[rune]rune)
	ra1 := []rune(spec.Inputs)
	ra2 := []rune(spec.Replacements)
	for i := 0; i < len(ra1); i++ {
		tr[ra1[i]] = ra2[i]
	}

	return &Tr{
		tr: tr,
	}
}

func (tr *Tr) Apply(input string) string {
	inRunes := []rune(input)
	outRunes := make([]rune, len(inRunes))

	for i, inRune := range inRunes {
		outRune, ok := tr.tr[inRune]
		if ok {
			outRunes[i] = outRune
		} else {
			outRunes[i] = inRune
		}
	}

	return string(outRunes)
}