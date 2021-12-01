package gonet

//共同场景
var commonScenes []Scene

type Scene interface {
	Handler(msg *Message) //场景消息处理入口
}

func AddCommonScene(sceneID uint8, scene Scene) {
	if commonScenes == nil {
		commonScenes = make([]Scene, int(sceneID)+1)
	}
	more := sceneID + 1 - uint8(len(commonScenes))
	for i := uint8(0); i < more; i++ {
		commonScenes = append(commonScenes, nil)
	}
	commonScenes[sceneID] = scene
}

func getCommonScene(sceneID uint8) Scene {
	if uint8(len(commonScenes)) <= sceneID {
		return nil
	}
	return commonScenes[sceneID]
}
