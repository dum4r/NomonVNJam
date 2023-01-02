package novel

type State *bool

func HideState(s State, ticks uint8) {
	var elmts []*Element
	for _, e := range win.elements {
		if e.state == s {
			elmts = append(elmts, e)
		}
	}

	for i, e := range elmts {
		if i != len(elmts)-1 {
			e.Opaque(ticks)
			continue
		}
		e.Hide(ticks)
	}
}

func ShowState(s State, ticks uint8) {
	for _, e := range win.elements {
		if e.state == s {
			e.Aclare(ticks)
		}
	}
	*s = true
}
