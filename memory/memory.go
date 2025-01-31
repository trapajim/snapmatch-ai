package memory

import (
	"container/list"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"sync"
	"time"
)

var m = &Provider{list: list.New()}

func init() {
	m.sessions = make(map[string]*list.Element, 0)
	snapmatchai.Register("memory", m)
}

type SessionStore struct {
	sid          string              // unique session id
	timeAccessed time.Time           // last access time
	value        map[any]interface{} // session value stored inside
}

func (st *SessionStore) Set(key, value any) error {
	st.value[key] = value
	return m.SessionUpdate(st.sid)
}

func (st *SessionStore) Get(key any) any {
	err := m.SessionUpdate(st.sid)
	if err != nil {
		return nil
	}
	if v, ok := st.value[key]; ok {
		return v
	}
	return nil
}

func (st *SessionStore) Delete(key any) error {
	delete(st.value, key)
	return m.SessionUpdate(st.sid)
}

func (st *SessionStore) SessionID() string {
	return st.sid
}

type Provider struct {
	lock     sync.Mutex               // lock
	sessions map[string]*list.Element // save in memory
	list     *list.List               // gc
}

func (m *Provider) SessionInit(sid string) (snapmatchai.Session, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	v := make(map[any]any)
	newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := m.list.PushBack(newsess)
	m.sessions[sid] = element
	return newsess, nil
}

func (m *Provider) SessionRead(sid string) (snapmatchai.Session, error) {
	if element, ok := m.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	}
	sess, err := m.SessionInit(sid)
	return sess, err
}

func (m *Provider) SessionDestroy(sid string) error {
	if element, ok := m.sessions[sid]; ok {
		delete(m.sessions, sid)
		m.list.Remove(element)
		return nil
	}
	return nil
}

func (m *Provider) SessionGC(maxLifetime int64) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for {
		element := m.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxLifetime) < time.Now().Unix() {
			m.list.Remove(element)
			delete(m.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (m *Provider) SessionUpdate(sid string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if element, ok := m.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		m.list.MoveToFront(element)
		return nil
	}
	return nil
}
