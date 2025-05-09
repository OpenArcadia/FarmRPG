package scenes

type Scene interface {
	Create()
	Render()
	Dispose()
}

var current Scene
var needInit = true

func ChangeScreen(s Scene) {
	if current != nil {
		current.Dispose()
	}
	current = s
	needInit = true
}

func Update() {
	if current == nil {
		return
	}

	if needInit {
		current.Create()
		needInit = false
	}

	current.Render()
}
