package workers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"movies-auth/users/internal/domain"
	"sync"
	"time"
)

var ErrStopped = errors.New("worker stopped")
var notificationsCount int
var mu sync.Mutex

type UsersStore interface {
	UpdateNotificationSent(ctx context.Context, id int) error
	GetUsersWithExpiredPassword(ctx context.Context) ([]domain.User, error)
}

type PassCheckWorker struct {
	store        UsersStore
	workersCount int
	interval     time.Duration
}

func NewPassCheckWorker(interval time.Duration, store UsersStore, workersCount int) PassCheckWorker {
	return PassCheckWorker{
		interval:     interval,
		store:        store,
		workersCount: workersCount,
	}
}

func (w PassCheckWorker) Run(ctx context.Context) error {
	ch := make(chan int)

	nwg := sync.WaitGroup{}
	for i := 0; i < w.workersCount; i++ {
		workerNum := i + 1
		nwg.Add(1)
		go w.notificationWorker(ctx, &nwg, workerNum, ch)
	}

	ticker := time.NewTicker(w.interval)
	for {
		select {
		case <-ctx.Done():
			close(ch)
			// необходимо дождаться выполнения работы воркеров
			nwg.Wait()
			return ErrStopped
		case <-ticker.C:
			err := w.checkUsersPasswords(ctx, &nwg, ch)
			if errors.Is(err, ErrStopped) {
				return ErrStopped
			}
			if err != nil {
				log.Println("check passwords: %w", err)
			}
		}
	}
}

func (w PassCheckWorker) checkUsersPasswords(ctx context.Context, nwg *sync.WaitGroup, ch chan<- int) error {
	users, err := w.store.GetUsersWithExpiredPassword(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {
		select {
		case <-ctx.Done():
			// необходимо дождаться выполнения работы воркеров
			nwg.Wait()
			fmt.Printf("total notifications sent: %d\n", notificationsCount)
			return ErrStopped
		default:
		}

		ch <- u.ID
	}

	// ждем завершения отправки всех сообщений
	for {
		if notificationsCount == len(users) {
			fmt.Printf("total notifications sent: %d\n", notificationsCount)
			notificationsCount = 0
			return nil
		}
	}
}

func (w PassCheckWorker) notificationWorker(ctx context.Context, wg *sync.WaitGroup, workerNum int, ch <-chan int) {
	fmt.Printf("starting notification worker: %d\n", workerNum)
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case id := <-ch:
			time.Sleep(time.Second * 1)
			dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			err := w.store.UpdateNotificationSent(dbCtx, id)
			if err != nil {
				log.Printf("user %d update failed: %s\n", id, err.Error())
			}
			log.Printf("worker: %d message sent for user: %d\n", workerNum, id)
			mu.Lock()
			notificationsCount++
			mu.Unlock()
		}
	}
}
