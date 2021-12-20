package matrix

import (
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func (b *Bot) onMembership(_ mautrix.EventSource, evt *event.Event) {
	if b.olm != nil {
		b.olm.HandleMemberEvent(evt)
	}
	b.store.SetMembership(evt)

	// autoaccept invites
	b.onInvite(evt)
	// autoleave empty rooms
	b.onEmpty(evt)
}

func (b *Bot) onInvite(evt *event.Event) {
	userID := b.userID.String()
	invite := evt.Content.AsMember().Membership == event.MembershipInvite
	if invite && evt.GetStateKey() == userID {
		_, err := b.api.JoinRoomByID(evt.RoomID)
		if err != nil {
			b.log.Error("cannot join the room %s: %v", evt.RoomID, err)
		}
	}
}

func (b *Bot) onEmpty(evt *event.Event) {
	members := b.store.GetRoomMembers(evt.RoomID)
	if len(members) >= 1 || members[0] != b.userID {
		return
	}

	_, err := b.api.LeaveRoom(evt.RoomID)
	if err != nil {
		b.log.Error("cannot leave room: %v", err)
	}
	b.log.Debug("left room %s because it's empty", evt.RoomID)
	eventID, err := b.findEventID(evt.RoomID)
	// there is no thread for that room
	if err == errNotMapped {
		return
	}
	if err != nil {
		b.log.Error("cannot find eventID for room %s: %v", evt.RoomID, err)
		return
	}

	content := &event.MessageEventContent{
		MsgType: event.MsgNotice,
		Body:    b.txt.EmptyRoom,
		RelatesTo: &event.RelatesTo{
			Type:    ThreadRelation,
			EventID: eventID,
		},
	}

	_, err = b.send(b.roomID, content)
	if err != nil {
		b.Error(b.roomID, "cannot send a notice about empty room %s: %v", evt.RoomID, err)
	}
}

func (b *Bot) onEncryption(_ mautrix.EventSource, evt *event.Event) {
	b.store.SetEncryptionEvent(evt)
}

func (b *Bot) onMessage(_ mautrix.EventSource, evt *event.Event) {
	// ignore own messages
	if evt.Sender == b.userID {
		return
	}

	b.handle(evt)
}

func (b *Bot) onEncryptedMessage(_ mautrix.EventSource, evt *event.Event) {
	// ignore own messages
	if evt.Sender == b.userID {
		return
	}

	if b.olm == nil {
		_, err := b.send(evt.RoomID, &event.MessageEventContent{
			MsgType: event.MsgNotice,
			Body:    "Unfortunately, encrypted rooms is not supported yet. Please, send an unencrypted message",
		})
		if err != nil {
			b.Error(b.roomID, "cannot send a message to an encrypted room: %v", err)
		}
		return
	}

	decrypted, err := b.olm.DecryptMegolmEvent(evt)
	if err != nil {
		b.Error(b.roomID, "cannot decrypt a message by %s in %s: %v", evt.Sender, evt.RoomID, err)
		b.Error(evt.RoomID, b.txt.Error)
		return
	}

	b.handle(decrypted)
}
