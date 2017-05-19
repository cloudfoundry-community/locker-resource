package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Locker struct {
	LockConfig string
}

func (l Locker) GetLocks() (map[string]string, error) {
	locks := map[string]string{}

	data, err := ioutil.ReadFile(l.LockConfig)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &locks)
	if err != nil {
		return nil, err
	}
	return locks, nil
}
func (l Locker) SaveLocks(locks map[string]string) error {
	data, err := json.MarshalIndent(locks, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(l.LockConfig, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
func (l Locker) GetLock(pool string) (string, error) {
	locks, err := l.GetLocks()
	if err != nil {
		return "", err
	}
	return locks[pool], nil
}
func (l Locker) Lock(pool, lock string) error {
	locks, err := l.GetLocks()
	if err != nil {
		return err
	}
	if current, exists := locks[pool]; exists && current != "" {
		return fmt.Errorf("Attempt to steal lock for %s by %s thwarted. Currently held by someone else", pool, lock)
	}
	locks[pool] = lock
	return l.SaveLocks(locks)
}
func (l Locker) Unlock(pool string, lock string) error {
	locks, err := l.GetLocks()
	if err != nil {
		return err
	}
	if current, exists := locks[pool]; exists && current != lock && current != "" {
		return fmt.Errorf("Attempt to unlock %s by %s thwarted. Currently held by someone else", pool, lock)
	}
	locks[pool] = ""
	return l.SaveLocks(locks)
}
