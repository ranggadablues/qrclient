package qrclient

import (
	"sync"
	"time"
)

const TokenLength = 32

type tokenMap struct {
	lock sync.Mutex
	list map[string]tokenData
	run  bool
}

type tokenData struct {
	LoginName string
	ExpireAt  time.Time
}

func newTokenMap() *tokenMap {
	t := &tokenMap{
		list: make(map[string]tokenData),
		run:  false,
	}

	return t
}

func (t *tokenMap) get(token string) string {
	t.lock.Lock()
	defer t.lock.Unlock()

	data, found := t.list[token]
	if found && time.Now().Before(data.ExpireAt) {
		return data.LoginName
	}

	delete(t.list, token)

	return ""
}

func (t *tokenMap) add(login string) string {
	token := randStr(TokenLength)

	t.lock.Lock()
	defer t.lock.Unlock()
	for {
		_, found := t.list[token]
		if !found {
			break
		}

		token = randStr(TokenLength)
	}

	data := tokenData{
		LoginName: login,
		ExpireAt:  time.Now().Add(1 * time.Minute),
	}

	t.list[token] = data

	return token
}

func (t *tokenMap) start() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.run {
		return
	}

	t.run = true
	go t.cleanerLoop()
}

func (t *tokenMap) stop() {
	t.run = false
}

func (t *tokenMap) clean() {
	t.lock.Lock()
	defer t.lock.Unlock()

	now := time.Now()
	dellist := []string{}
	for k, v := range t.list {
		if v.ExpireAt.After(now) {
			dellist = append(dellist, k)
		}
	}

	for _, i := range dellist {
		delete(t.list, i)
	}
}

func (t *tokenMap) cleanerLoop() {
	for {
		time.Sleep(10 * time.Second)
	}
}
