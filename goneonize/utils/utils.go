package utils

import (
	// "go.mau.fi/whatsmeow"
	"github.com/Nubuki-all/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

var MediaType = []whatsmeow.MediaType{
	whatsmeow.MediaImage,
	whatsmeow.MediaVideo,
	whatsmeow.MediaAudio,
	whatsmeow.MediaDocument,
	whatsmeow.MediaHistory,
	whatsmeow.MediaAppState,
	whatsmeow.MediaLinkThumbnail,
	whatsmeow.MediaStickerPack,

}



var ChatPresence = []types.ChatPresence{
	types.ChatPresenceComposing,
	types.ChatPresencePaused,
}

var ChatPresenceMedia = []types.ChatPresenceMedia{
	types.ChatPresenceMediaText,
	types.ChatPresenceMediaAudio,
}
