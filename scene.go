package goNet

var scenes []Scene

type Scene interface {
	Handler(msg *Msg)
}

func AddCommonScene(sceneID uint8, scene Scene) {
	if scenes == nil {
		scenes = make([]Scene, int(sceneID)+1)
	}
	more := sceneID + 1 - uint8(len(scenes))
	for i := uint8(0); i < more; i++ {
		scenes = append(scenes, nil)
	}
	scenes[sceneID] = scene
}

func GetCommonScene(sceneID uint8) Scene {
	if uint8(len(scenes)) < sceneID {
		return nil
	}
	return scenes[sceneID]
}
