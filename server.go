package main

import (
	"fmt"
	"os"
)

type LockOperation int

const (
	LockOp LockOperation = iota
	UnlockOp
	ListOp
)

type LockStatus int

const (
	Locked LockStatus = iota
	Unlocked
	Error
	None
)

type LockRequest struct {
	Command  LockOperation
	Pool     string
	Lock     string
	Response chan LockResponse
}
type LockResponse struct {
	Status  LockStatus
	Message interface{}
	Error   error
}

type LockInput struct {
	Lock string `json: "lock"`
}

func lockServer(lockRequests chan LockRequest, lockConfig string) {
	locker := Locker{LockConfig: lockConfig}
	for req := range lockRequests {
		if req.Command == ListOp {
			locks, err := locker.GetLocks()

			res := LockResponse{
				Status:  None,
				Message: locks,
				Error:   err,
			}
			req.Response <- res
		} else if req.Command == LockOp {
			res := LockResponse{}

			current, err := locker.GetLock(req.Pool)
			if err != nil {
				res.Status = Error
				res.Error = err
				req.Response <- res
				continue
			}

			err = locker.Lock(req.Pool, req.Lock)
			if err != nil {
				res.Status = Error
				res.Error = err
				req.Response <- res
				continue
			}
			current, err = locker.GetLock(req.Pool)
			if err != nil {
				res.Status = Error
				res.Error = err
				req.Response <- res
				continue
			}
			if current != req.Lock {
				res.Status = Error
				res.Error = fmt.Errorf("Locking failed. Should be locked by %s, but found %s", req.Lock, current)
				req.Response <- res
				continue
			}

			res.Status = Locked
			res.Message = map[string]string{
				"response": fmt.Sprintf("lock for %s acquired by %s", req.Pool, req.Lock),
			}
			req.Response <- res
		} else if req.Command == UnlockOp {
			res := LockResponse{}

			err := locker.Unlock(req.Pool, req.Lock)
			if err != nil {
				res.Status = Error
				res.Error = err
				req.Response <- res
				continue
			}

			res.Status = Unlocked
			res.Message = map[string]string{"response": fmt.Sprintf("Lock released on %s", req.Pool)}
			req.Response <- res
		} else {
			fmt.Fprintf(os.Stderr, "Invalid lock request '%s'", req.Command)
			req.Response <- LockResponse{Error: fmt.Errorf("Invalid lock request '%s'", req.Command)}
		}
	}
}
