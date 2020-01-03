package session

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"sync"
)

const (
	SessionFlagNone = iota
	SessionFlagModify
	SessionFlagLoad
)

type RedisSession struct {
	id string
	pool *redis.Pool
	data map[string]interface{}
	rwlock sync.RWMutex
	flag int
}

func NewRedisSession(id string, pool *redis.Pool) *RedisSession {
	s := &RedisSession{
		id:     id,
		pool:   pool,
		data:   make(map[string]interface{},8),
		rwlock: sync.RWMutex{},
		flag:   0,
	}
	return s
}

func (r *RedisSession) Set (key string, value interface{}) error{
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	r.data[key] = value
	r.flag = SessionFlagModify
	return nil
}

func (r *RedisSession) loadFromRedis() (err error) {

	conn := r.pool.Get()
	reply, err := conn.Do("GET", r.id)
	if err != nil {
		return
	}

	data, err := redis.String(reply, err)
	if err != nil{
		return
	}

	err = json.Unmarshal([]byte(data), &r.data)
	if err != nil {
		return
	}
	return
}

func (r *RedisSession) Get(key string) (result interface{}, err error){

	r.rwlock.RLock()
	defer r.rwlock.RUnlock()
	
	// 实现一个延迟加载功能
	if r.flag == SessionFlagNone {
	//	初始化seesion为0 还未加载
		err = r.loadFromRedis()
		if err != nil {
			return 
		}
	}
	
	result, ok := r.data[key]
	if !ok {
		err = ErrKeyNotExistInSession
		return 
	}
	return 
}

func (m *RedisSession) Id() string {
	return m.id
}

func (r *RedisSession) Del(key string) error {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	
	r.flag = SessionFlagModify
	delete(r.data, key)
	return nil
}

func (r *RedisSession) Save() (err error){
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	
	if r.flag != SessionFlagModify {
		return 
	}
	
	data, err := json.Marshal(r.data)
	if err != nil {
		return 
	}
	
	conn := r.pool.Get()
	_, err = conn.Do("SET", r.id, string(data))
	if err != nil {
		return 
	}
	return 
}

func (m *RedisSession) IsModify() bool {
	if m.flag == SessionFlagModify{
		return true
	}
	return false
}




