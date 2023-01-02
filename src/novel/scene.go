package novel

var currentScene *sceneWin

type Scene1 interface {
	InitElmts()
	Update() error
	Destroy()
}

type sceneWin interface {
	init()
	IsDone() bool
	Update() error
	Destroy()
}

func AddScene(name string, object Scene1) (cursor sceneWin) {
	if object == nil {
		panic("Error novel: can't create scene without an object")
	}
	win.scenes[name] = &cursor
	return &scene{&cursor, false, object}
}

type scene struct {
	pointer *sceneWin
	done    bool
	Scene   Scene1
}

func (s *scene) init() {
	currentScene = s.pointer
	s.Scene.InitElmts()
	currentScene = nil
	s.done = true
}

func (s *scene) IsDone() bool  { return s.done }
func (s *scene) Destroy()      { s.Scene.Destroy() }
func (s *scene) Update() error { return s.Scene.Update() }
