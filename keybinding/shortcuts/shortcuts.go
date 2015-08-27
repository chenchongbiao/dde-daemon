/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package shortcuts

const (
	KeyTypeSystem int32 = iota
	KeyTypeMedia
	KeyTypeWM
	KeyTypeCustom
)

type Shortcut struct {
	Id   string //schema key
	Name string

	Type int32

	Accels []string
}
type Shortcuts []*Shortcut

func ListAllShortcuts() Shortcuts {
	list := ListSystemShortcuts()
	list = append(list, ListMediaShortcuts()...)
	list = append(list, ListWMShortcuts()...)
	list = append(list, ListCustomKey().GetShortcuts()...)

	return list
}

func Reset() {
	resetSystemAccels()
	resetWMAccels()
	resetMediaAccels()
}

func (list Shortcuts) GetById(id string, ty int32) *Shortcut {
	for _, s := range list {
		if s.Id == id && s.Type == ty {
			return s
		}
	}

	return nil
}

func (list Shortcuts) GetByAccel(accel string) *Shortcut {
	for _, s := range list {
		if isAccelInList(accel, s.Accels) {
			return s
		}
	}
	return nil
}

func (list Shortcuts) Add(id string, ty int32) Shortcuts {
	item := list.GetById(id, ty)
	if item != nil {
		list = list.Delete(item.Id, item.Type)
	}

	s := newShortcut(id, ty)
	if s != nil {
		list = append(list, s)
	}
	return list
}

func (list Shortcuts) Delete(id string, ty int32) Shortcuts {
	var newList Shortcuts
	for _, s := range list {
		if s.Id == id && s.Type == ty {
			continue
		}
		newList = append(newList, s)
	}

	return newList
}

func (s *Shortcut) Disable() {
	switch s.Type {
	case KeyTypeSystem:
		disableSystemAccels(s.Id)
	case KeyTypeWM:
		disableWMAccels(s.Id)
	case KeyTypeMedia:
		disableMediaAccels(s.Id)
	case KeyTypeCustom:
		disableCustomKey(s.Id)
	}
}

func (s *Shortcut) AddAccel(accel string) {
	switch s.Type {
	case KeyTypeSystem:
		addSystemAccel(s.Id, accel)
	case KeyTypeWM:
		addWMAccel(s.Id, accel)
	case KeyTypeMedia:
		addMediaAccel(s.Id, accel)
	case KeyTypeCustom:
		modifyCustomAccels(s.Id, accel, false)
	}
}

func (s *Shortcut) DeleteAccel(accel string) {
	switch s.Type {
	case KeyTypeSystem:
		delSystemAccel(s.Id, accel)
	case KeyTypeWM:
		delWMAccel(s.Id, accel)
	case KeyTypeMedia:
		delMediaAccel(s.Id, accel)
	case KeyTypeCustom:
		modifyCustomAccels(s.Id, accel, true)
	}
}

func (s *Shortcut) SetName(name string) {
	if s.Type != KeyTypeCustom {
		return
	}
	modifyCustomName(s.Id, name)
}

func (s *Shortcut) SetAction(action string) {
	if s.Type != KeyTypeCustom {
		return
	}
	modifyCustomAction(s.Id, action)
}

func (s *Shortcut) GetAction() string {
	switch s.Type {
	case KeyTypeSystem:
		return getSystemAction(s.Id)
	case KeyTypeCustom:
		info := ListCustomKey().Get(s.Id)
		if info != nil {
			return info.Action
		}
	}
	return ""
}

func newShortcut(id string, ty int32) *Shortcut {
	var list Shortcuts
	switch ty {
	case KeyTypeSystem:
		list = ListSystemShortcuts()
	case KeyTypeMedia:
		list = ListMediaShortcuts()
	case KeyTypeWM:
		list = ListWMShortcuts()
	case KeyTypeCustom:
		list = ListCustomKey().GetShortcuts()
	default:
		return nil
	}

	return list.GetById(id, ty)
}
