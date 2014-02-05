package library

import "errors"

type MusicEntry struct {
	Id     string
	Name   string
	Artist string
	Source string
	Type   string
}

type MusicManager struct {
	musics []MusicEntry
}

func NewMusicManager() *MusicManager {
	return &MusicManager(make([]MusicEntry, 0))
}

func (m *MusicManager) Len() int {
	return len(m.musics)
}

func (m *MusicManager) Get(index int) (music *MusicEntry, err error) {
	if index >= m.Len() || index < 0 {
		return nil, errors.New("Index out of range.")
	}
	return &m.musics[index], nil
}
func (m *MusicManager) Find(name string) *MusicEntry {
	if m.Len() == 0 {
		return nil
	}

	for _, music := range m.musics {
		if music.Name == name {
			return &music
		}
	}
	return nil
}

func (m *MusicManager) Add(music *MusicEntry) {
	m.musics = append(m.musics, *music)
}

func (m *MusicManager) Remove(index int) *MusicEntry {
	if index < 0 || index >= len(m.musics) {
		return nil
	}
	removedMusic := &m.musics[index]
	// 从数组切片中删除元素
	if index < len(m.musics)-1 { // 中间元素
		m.musics = append(m.musics[:index-1], m.musics[index+1:]...)
	} else if index == 0 { // 删除仅有的一个元素
		m.musics = make([]MusicEntry, 0)
	} else { // 删除的是最后一个元素
		m.musics = m.musics[:index-1]
	}
	return removedMusic
}
