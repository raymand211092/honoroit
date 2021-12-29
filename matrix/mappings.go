package matrix

import (
	"errors"
	"strings"

	"maunium.net/go/mautrix/id"
)

type accountDataMappings struct {
	Rooms  map[id.RoomID]id.EventID `json:"rooms"`
	Events map[id.EventID]id.RoomID `json:"events"`
}

// errNotMapped returned if roomID or eventID doesn't exist in room<->event map (yet)
var errNotMapped = errors.New("cannot find appropriate mapping")

func (b *Bot) getMappings() (*accountDataMappings, error) {
	mappings := &accountDataMappings{
		Rooms:  make(map[id.RoomID]id.EventID),
		Events: make(map[id.EventID]id.RoomID),
	}
	b.log.Debug("trying to get mappings")

	cached := b.cache.Get(accountDataRooms)
	if cached != nil {
		var ok bool
		mappings, ok = cached.(*accountDataMappings)
		if ok {
			b.log.Debug("got mappings from cache")
			return mappings, nil
		}
	}

	b.log.Debug("mappings not cached yet, trying to get them from account data")
	err := b.api.GetAccountData(accountDataRooms, &mappings)
	if err != nil {
		if strings.Contains(err.Error(), "M_NOT_FOUND") {
			return nil, nil
		}
		return nil, err
	}
	b.cache.Set(accountDataRooms, mappings)
	return mappings, err
}

func (b *Bot) addMapping(roomID id.RoomID, eventID id.EventID) error {
	b.log.Debug("adding new mapping: %s<->%s", roomID, eventID)
	data, err := b.getMappings()
	if err != nil {
		return err
	}

	if data == nil {
		data = &accountDataMappings{
			Rooms:  make(map[id.RoomID]id.EventID),
			Events: make(map[id.EventID]id.RoomID),
		}
	}

	data.Rooms[roomID] = eventID
	data.Events[eventID] = roomID

	b.cache.Set(accountDataRooms, data)
	return b.api.SetAccountData(accountDataRooms, data)
}

func (b *Bot) removeMapping(roomID id.RoomID, eventID id.EventID) error {
	b.log.Debug("removing mapping %s<->%s...", roomID, eventID)
	data, err := b.getMappings()
	if err != nil {
		return err
	}

	if data == nil {
		b.log.Debug("no mappings, so nothing to remove")
		return nil
	}

	delete(data.Rooms, roomID)
	delete(data.Events, eventID)

	b.log.Debug("mapping has been removed, uploading data...")
	b.cache.Set(accountDataRooms, data)
	return b.api.SetAccountData(accountDataRooms, data)
}

// findRoomID by eventID
func (b *Bot) findRoomID(eventID id.EventID) (id.RoomID, error) {
	b.log.Debug("trying to find room ID by eventID = %s", eventID)
	mappings, err := b.getMappings()
	if err != nil {
		return "", err
	}

	if mappings == nil {
		b.log.Debug("cannot get mappings: no error returned, but mappings = nil")
		return "", errNotMapped
	}

	roomID, ok := mappings.Events[eventID]
	if !ok || roomID == "" {
		b.log.Debug("room not found in existing mappings")
		return "", errNotMapped
	}

	return roomID, nil
}

// findEventID by roomID
func (b *Bot) findEventID(roomID id.RoomID) (id.EventID, error) {
	b.log.Debug("trying to find event ID by roomID = %s", roomID)
	mappings, err := b.getMappings()
	if err != nil {
		return "", err
	}

	if mappings == nil {
		b.log.Debug("mappings not created yet in account data, seems like first run")
		return "", errNotMapped
	}

	eventID, ok := mappings.Rooms[roomID]
	if !ok || eventID == "" {
		b.log.Debug("event not found in existing mappings")
		return "", errNotMapped
	}

	return eventID, nil
}
