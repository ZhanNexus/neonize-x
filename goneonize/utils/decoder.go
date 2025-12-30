package utils

import (
	"C"
	"time"

	defproto "github.com/krypton-byte/neonize/defproto"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	waBinary "go.mau.fi/whatsmeow/binary"
	waVname "go.mau.fi/whatsmeow/proto/waVnameCert"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
)

import (
	"go.mau.fi/whatsmeow/types/events"
)

func DecodeJidProto(data *defproto.JID) types.JID {
	return types.JID{
		User:       *data.User,
		RawAgent:   uint8(*data.RawAgent),
		Device:     uint16(*data.Device),
		Integrator: uint16(*data.Integrator),
		Server:     *data.Server,
	}
}

func DecodeGroupParent(groupParent *defproto.GroupParent) types.GroupParent {
	return types.GroupParent{
		IsParent:                      *groupParent.IsParent,
		DefaultMembershipApprovalMode: *groupParent.DefaultMembershipApprovalMode,
	}
}

func DecodeGroupLinkedParent(groupLinkedParent *defproto.GroupLinkedParent) types.GroupLinkedParent {
	return types.GroupLinkedParent{
		LinkedParentJID: DecodeJidProto(groupLinkedParent.LinkedParentJID),
	}
}

func DecodeReqCreateGroup(reqCreateGroup *defproto.ReqCreateGroup) whatsmeow.ReqCreateGroup {
	participants := []types.JID{}
	for _, participant := range reqCreateGroup.Participants {
		participants = append(participants, DecodeJidProto(participant))
	}
	new_type := whatsmeow.ReqCreateGroup{
		Name:         *reqCreateGroup.Name,
		Participants: participants,
		CreateKey:    *reqCreateGroup.CreateKey,
	}
	if reqCreateGroup.GroupParent != nil {
		new_type.GroupParent = DecodeGroupParent(reqCreateGroup.GroupParent)
	}
	if reqCreateGroup.GroupLinkedParent != nil {
		new_type.GroupLinkedParent = DecodeGroupLinkedParent(reqCreateGroup.GroupLinkedParent)
	}
	return new_type
}

func DecodeAddressingMode(mode_types *defproto.AddressingMode) types.AddressingMode {
	var AddressingMode types.AddressingMode
	switch mode_types {
	case defproto.AddressingMode_PN.Enum():
		AddressingMode = types.AddressingModePN
	case defproto.AddressingMode_LID.Enum():
		AddressingMode = types.AddressingModeLID
	}
	return AddressingMode
}

func DecodeMessageSource(messageSource *defproto.MessageSource) types.MessageSource {
	model := types.MessageSource{
		Chat:     DecodeJidProto(messageSource.Chat),
		Sender:   DecodeJidProto(messageSource.Sender),
		IsFromMe: *messageSource.IsFromMe,
		IsGroup:  *messageSource.IsGroup,

		SenderAlt:    DecodeJidProto(messageSource.SenderAlt),
		RecipientAlt: DecodeJidProto(messageSource.RecipientAlt),

		BroadcastListOwner: DecodeJidProto(messageSource.BroadcastListOwner),
	}
	if messageSource.AddressingMode != nil {
		model.AddressingMode = DecodeAddressingMode(messageSource.AddressingMode)
	}
	return model
}

func DecodeVerifiedNameCertificate(verifiedNameCertificate *waVname.VerifiedNameCertificate) *waVname.VerifiedNameCertificate {
	// passing types through protobuf
	return verifiedNameCertificate
}

func DecodeVerifiedNameDetails(verifiedNameDetails *waVname.VerifiedNameCertificate_Details) *waVname.VerifiedNameCertificate_Details {
	return verifiedNameDetails
}

func DecodeVerifiedName(verifiedName *defproto.VerifiedName) *types.VerifiedName {
	verifiednametypes := types.VerifiedName{}
	if verifiedName.Certificate != nil {
		verifiednametypes.Certificate = verifiedName.Certificate
	}
	if verifiedName.Details != nil {
		verifiednametypes.Details = verifiedName.Details
	}
	return &verifiednametypes
}

func DecodeDeviceSentMeta(deviceSentMeta *defproto.DeviceSentMeta) *types.DeviceSentMeta {
	return &types.DeviceSentMeta{
		DestinationJID: *deviceSentMeta.DestinationJID,
		Phash:          *deviceSentMeta.Phash,
	}
}

func DecodeMessageInfo(messageInfo *defproto.MessageInfo) *types.MessageInfo {
	ts := *messageInfo.Timestamp
	model := &types.MessageInfo{
		MessageSource: DecodeMessageSource(messageInfo.MessageSource),
		ID:            *messageInfo.ID,
		ServerID:      int(*messageInfo.ServerID),
		Type:          *messageInfo.Type,
		PushName:      *messageInfo.Pushname,
		Timestamp:     time.Unix(0, ts),
		Category:      *messageInfo.Category,
		Multicast:     *messageInfo.Multicast,
		MediaType:     *messageInfo.MediaType,
		Edit:          types.EditAttribute(*messageInfo.Edit),
	}
	if messageInfo.VerifiedName != nil {
		model.VerifiedName = DecodeVerifiedName(messageInfo.VerifiedName)
	}
	if messageInfo.DeviceSentMeta != nil {
		model.DeviceSentMeta = DecodeDeviceSentMeta(messageInfo.DeviceSentMeta)
	}
	return model
}

func DecodeCreateNewsletterParams(createletterNewsParams *defproto.CreateNewsletterParams) whatsmeow.CreateNewsletterParams {
	return whatsmeow.CreateNewsletterParams{
		Name:        *createletterNewsParams.Name,
		Description: *createletterNewsParams.Description,
		Picture:     createletterNewsParams.Picture,
	}
}

func DecodeGetProfilePictureParams(params *defproto.GetProfilePictureParams) *whatsmeow.GetProfilePictureParams {
	if params.Preview == nil || params.ExistingID == nil || params.IsCommunity == nil {
		return nil
	}
	return &whatsmeow.GetProfilePictureParams{
		Preview:     *params.Preview,
		ExistingID:  *params.ExistingID,
		IsCommunity: *params.IsCommunity,
	}
}

func DecodeMutationInfo(mutationInfo *defproto.MutationInfo) appstate.MutationInfo {
	return appstate.MutationInfo{
		Index:   mutationInfo.Index,
		Version: *mutationInfo.Version,
		Value:   mutationInfo.Value,
	}
}

func DecodePatchInfo(patchInfo *defproto.PatchInfo) *appstate.PatchInfo {
	var Type appstate.WAPatchName
	switch patchInfo.Type {
	case defproto.PatchInfo_CRITICAL_BLOCK.Enum():
		Type = appstate.WAPatchCriticalBlock
	case defproto.PatchInfo_CRITICAL_UNBLOCK_LOW.Enum():
		Type = appstate.WAPatchCriticalUnblockLow
	case defproto.PatchInfo_REGULAR.Enum():
		Type = appstate.WAPatchRegular
	}
	mutationInfo := []appstate.MutationInfo{}
	for _, mutation := range patchInfo.Mutations {
		mutationInfo = append(mutationInfo, DecodeMutationInfo(mutation))
	}
	return &appstate.PatchInfo{
		Timestamp: time.Unix(0, *patchInfo.Timestamp),
		Type:      Type,
		Mutations: mutationInfo,
	}
}

func DecodeContactEntry(entry *defproto.ContactEntry) *store.ContactEntry {
	return &store.ContactEntry{
		JID:       DecodeJidProto(entry.JID),
		FirstName: *entry.FirstName,
		FullName:  *entry.FullName,
	}
}

func DecodeSendRequestExtra(extra *defproto.SendRequestExtra) whatsmeow.SendRequestExtra {
	var additionalNodes *[]waBinary.Node

	if len(extra.AdditionalNodes) > 0 {
		nodes := make([]waBinary.Node, 0, len(extra.AdditionalNodes))

		for _, pn := range extra.AdditionalNodes {
			n := DecodeNodeProto(pn)
			if n != nil {
				nodes = append(nodes, *n)
			}
		}

		if len(nodes) > 0 {
			additionalNodes = &nodes
		}
	}

	return whatsmeow.SendRequestExtra{
		ID:              types.MessageID(extra.GetID()),
		InlineBotJID:    DecodeJidProto(extra.InlineBotJID),
		Peer:            extra.GetPeer(),
		Timeout:         time.Duration(extra.GetTimeout()),
		MediaHandle:     extra.GetMediaHandle(),
		AdditionalNodes: additionalNodes,
	}
}

func DecodeNewsLetterMessageMeta(defproto.NewsLetterMessageMeta) {
}

func DecodeNodeProto(n *defproto.Node) *waBinary.Node {
	if n == nil {
		return nil
	}

	// explicit nil node
	if n.Nil != nil && *n.Nil {
		return nil
	}

	node := &waBinary.Node{
		Tag: n.Tag,
	}

	// attrs
	if len(n.Attrs) > 0 {
		node.Attrs = make(map[string]interface{}, len(n.Attrs))
		for _, a := range n.Attrs {
			switch v := a.Value.(type) {
			case *defproto.NodeAttrs_Boolean:
				node.Attrs[a.Name] = v.Boolean
			case *defproto.NodeAttrs_Integer:
				node.Attrs[a.Name] = v.Integer
			case *defproto.NodeAttrs_Text:
				node.Attrs[a.Name] = v.Text
			case *defproto.NodeAttrs_Jid:
				node.Attrs[a.Name] = DecodeJidProto(v.Jid)
			}
		}
	}

	// children
	for _, c := range n.Nodes {
		if child := DecodeNodeProto(c); child != nil {
			node.Children = append(node.Children, *child)
		}
	}

	// raw bytes node
	if len(n.Bytes) > 0 {
		node.Data = n.Bytes
	}

	return node
}

func DecodeEventTypesMessage(message *defproto.Message) *events.Message {
	model := &events.Message{
		Info:                  *DecodeMessageInfo(message.Info),
		IsEphemeral:           *message.IsEphemeral,
		IsViewOnce:            *message.IsViewOnce,
		IsViewOnceV2:          *message.IsViewOnceV2,
		IsEdit:                *message.IsEdit,
		IsViewOnceV2Extension: *message.IsViewOnceV2Extension,
		IsDocumentWithCaption: *message.IsDocumentWithCaption,
		IsLottieSticker:       *message.IsLottieSticker,
		UnavailableRequestID:  *message.UnavailableRequestID,
		RetryCount:            int(*message.RetryCount),
		RawMessage:            message.Message,
	}
	if message.NewsLetterMeta != nil {
		model.NewsletterMeta = &events.NewsletterMessageMeta{
			EditTS:     time.Unix(0, *message.NewsLetterMeta.EditTS),
			OriginalTS: time.Unix(0, *message.NewsLetterMeta.OriginalTS),
		}
	}
	if message.SourceWebMsg != nil {
		model.SourceWebMsg = message.SourceWebMsg
	}
	if message.Message != nil {
		model.Message = message.Message
	}
	return model
}
