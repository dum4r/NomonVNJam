package novel

type animation func() bool

func (elm *Element) Opaque(ticks uint8) {
	var aphScale float64 = 1
	elm.animations = append(elm.animations, func() bool {
		aphScale -= (1 / float64(ticks))
		elm.ImageOptions.ColorM.Reset()
		elm.ImageOptions.ColorM.Scale(1, 1, 1, aphScale)
		return aphScale > 0
	})
}

// Animation hide aph color and desactive elm
func (elm *Element) Hide(ticks uint8) {
	if ticks == 0 {
		*elm.state = false
		return
	}
	var aphScale float64 = 1
	elm.animations = append(elm.animations, func() bool {
		aphScale -= (1 / float64(ticks))
		elm.ImageOptions.ColorM.Reset()
		elm.ImageOptions.ColorM.Scale(1, 1, 1, aphScale)
		if aphScale < 0 {
			elm.ImageOptions.ColorM.Reset()
			*elm.state = false
			return false
		}
		return true
	})
}

// Animation hide aph color and desactive elm and exec function
func (elm *Element) HideFunc(ticks uint8, f func()) {
	if ticks == 0 {
		*elm.state = false
		f()
		return
	}
	var aphScale float64 = 1
	elm.animations = append(elm.animations, func() bool {
		aphScale -= (1 / float64(ticks))
		elm.ImageOptions.ColorM.Reset()
		elm.ImageOptions.ColorM.Scale(1, 1, 1, aphScale)
		if aphScale < 1 {
			elm.ImageOptions.ColorM.Reset()
			f()
			*elm.state = false
			return false
		}
		return true
	})
}

func (elm *Element) IsActive() bool { return *elm.state }

// Animation show aph color and active elm
func (elm *Element) Aclare(ticks uint8) {
	elm.ImageOptions.ColorM.Scale(1, 1, 1, 0)
	var aphScale float64 = 0
	elm.animations = append(elm.animations, func() bool {
		aphScale += (1 / float64(ticks))
		elm.ImageOptions.ColorM.Reset()
		elm.ImageOptions.ColorM.Scale(1, 1, 1, aphScale)
		return aphScale < 1
	})
}

// Animation show aph color and active State
func (elm *Element) Show(ticks uint8) {
	*elm.state = true
	if ticks == 0 {
		return
	}
	elm.ImageOptions.ColorM.Scale(1, 1, 1, 0)
	var aphScale float64 = 0
	elm.animations = append(elm.animations, func() bool {
		aphScale += (1 / float64(ticks))
		elm.ImageOptions.ColorM.Reset()
		elm.ImageOptions.ColorM.Scale(1, 1, 1, aphScale)
		if aphScale > 1 {
			elm.ImageOptions.ColorM.Scale(1, 1, 1, 1)
			return false
		}
		return true
	})
}

// Animation show aph color and active elm and exec function
func (elm *Element) ShowFunc(ticks uint8, f func(this *Element)) {
	*elm.state = true
	if ticks == 0 {
		return
	}
	var aphScale float64 = 0
	elm.animations = append(elm.animations, func() bool {
		aphScale += (1 / float64(ticks))
		elm.ImageOptions.ColorM.Reset()
		elm.ImageOptions.ColorM.Scale(1, 1, 1, aphScale)
		if aphScale > 1 {
			elm.ImageOptions.ColorM.Scale(1, 1, 1, 1)
			f(elm)
			return false
		}
		return true
	})
}

func (elm *Element) AniInfinite() {
	switch resour := elm.Resour.(type) {
	case *Oneframe:
		elm.animations = append(elm.animations, func() bool {
			resour.index++
			if resour.index > len(resour.imgs)-1 {
				resour.index = 0
			}
			return true // return true forever because the animation has no end
		})
	case *GroupImgFor:
		elm.animations = append(elm.animations, func() bool {
			resour.index++
			if resour.index > len(resour.imgs[resour.Current])-1 {
				resour.index = 0
			}
			return true // return true forever because the animation has no end
		})
	case *SimpleString:
		elm.animations = append(elm.animations, func() bool {
			elm.Recalc()
			return true // return true forever because the animation has no end
		})
	default:
		println("Alert novel: this resource is not possible infinite loop for = ", resour)
	}
}

func (elm *Element) AniBoomerang() {
	switch resour := elm.Resour.(type) {
	case *Oneframe:
		bln := true
		elm.animations = append(elm.animations, func() bool {
			if bln {
				resour.index++
				if resour.index > len(resour.imgs)-2 {
					bln = false
				}
			} else {
				resour.index--
				if resour.index < 1 {
					bln = true
				}
			}
			return true // return true forever because the animation has no end
		})
	default:
		println("Alert novel: this resource is not possible boomerang loop for = ", resour)
	}
}

func (elm *Element) AniEnding() {
	switch resour := elm.Resour.(type) {
	case *Oneframe:
		elm.animations = append(elm.animations, func() bool {
			if resour.index+1 > len(resour.imgs)-1 {
				return false
			}
			resour.index++
			return true // return true forever because the animation has no end
		})
	case *GroupImgFor:
		elm.animations = append(elm.animations, func() bool {
			if resour.index+1 > len(resour.imgs[resour.Current])-1 {
				return false
			}
			resour.index++
			return true // return true forever because the animation has no end
		})
	default:
		println("Alert novel: this resource is not possible ending for = ", resour)
	}
}

func (elm *Element) AniEndingFunc(f func()) {
	switch resour := elm.Resour.(type) {
	case *Oneframe:
		elm.animations = append(elm.animations, func() bool {
			if resour.index+1 > len(resour.imgs)-1 {
				f()
				return false
			}
			resour.index++
			return true // return true forever because the animation has no end
		})
	case *GroupImgFor:
		elm.animations = append(elm.animations, func() bool {
			if resour.index+1 > len(resour.imgs[resour.Current])-1 {
				f()
				return false
			}
			resour.index++
			return true // return true forever because the animation has no end
		})
	default:
		println("Alert novel: this resource is not possible ending for = ", resour)
	}
}

func (elm *Element) AniEndingFuncRever(f func()) {
	switch resour := elm.Resour.(type) {
	case *GroupImgFor:
		resour.index = len(resour.imgs[resour.Current]) - 1
		elm.animations = append(elm.animations, func() bool {
			if resour.index < 1 {
				f()
				return false
			}
			resour.index--
			return true // return true forever because the animation has no end
		})
	default:
		println("Alert novel: this resource is not possible ending for = ", resour)
	}
}
