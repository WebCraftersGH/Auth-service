package otpsrepo

import (
	"sync"
	"time"

	"github.com/WebCraftersGH/Auth-service/internal/domain"
)

type Repo struct {
	mu      sync.RWMutex
	storage map[string]otpRecord
}

func New() *Repo {
	return &Repo{storage: make(map[string]otpRecord)}
}

func (r *Repo) SaveCode(otp domain.OTP, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[key] = otpRecord{Code: cloneOTP(otp)}
	return nil
}

func (r *Repo) GetCode(key string) (domain.OTP, error) {
	r.mu.RLock()
	record, ok := r.storage[key]
	r.mu.RUnlock()
	if !ok {
		return domain.OTP{}, domain.ErrOTPNotFound
	}

	if !record.Code.ExpiresAt.IsZero() && time.Now().After(record.Code.ExpiresAt) {
		_ = r.DeleteCode(key)
		return domain.OTP{}, domain.ErrOTPExpired
	}

	return cloneOTP(record.Code), nil
}

func (r *Repo) IncAttempts(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.storage[key]
	if !ok {
		return domain.ErrOTPNotFound
	}
	record.Code.Attempts++
	r.storage[key] = record
	return nil
}

func (r *Repo) DeleteCode(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.storage, key)
	return nil
}
