package services

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Определение ошибок на русском языке
var (
	ErrJobQueueIsFull = errors.New("очередь задач заполнена")
)

// Тип Job представляет функцию, выполняющуюся в контексте
type Job func(ctx context.Context)

// JobQueueService управляет очередью задач и их выполнением
type JobQueueService struct {
	jobs   chan Job      // Канал для очереди задач
	resume chan struct{} // Канал для возобновления выполнения при паузе
	paused int32         // Флаг состояния паузы (0 - работа, 1 - пауза)
	wg     sync.WaitGroup // Группа ожидания для управления горутинами
}

// NewJobQueueService создает новый сервис очереди задач с заданной емкостью и числом рабочих
func NewJobQueueService(ctx context.Context, capacity, workers int) *JobQueueService {
	service := &JobQueueService{
		jobs:   make(chan Job, capacity),
		resume: make(chan struct{}),
		wg:     sync.WaitGroup{},
	}
	service.start(ctx, workers)

	return service
}

// start запускает рабочих (воркеров), обрабатывающих задачи из очереди
func (jqs *JobQueueService) start(ctx context.Context, workers int) {
	for i := 0; i < workers; i++ {
		jqs.wg.Add(1)

		go func() {
			defer jqs.wg.Done()

			for {
				select {
				case job, ok := <-jqs.jobs:
					if !ok {
						return // Выход, если канал закрыт
					}

					// Проверка состояния паузы
					if atomic.LoadInt32(&jqs.paused) == 1 {
						<-jqs.resume // Ожидание сигнала возобновления
					}

					job(ctx) // Выполнение задачи
				case <-ctx.Done():
					return // Выход при завершении контекста
				}
			}
		}()
	}
}

// Enqueue добавляет задачу в очередь. Возвращает ошибку, если очередь заполнена.
func (jqs *JobQueueService) Enqueue(job Job) {
	jqs.jobs <- job
}


// ScheduleJob планирует выполнение задачи через заданную задержку
func (jqs *JobQueueService) ScheduleJob(job Job, delay time.Duration) {
	time.AfterFunc(delay, func() {
		jqs.jobs <- job
	})
}
// Pause приостанавливает выполнение задач
func (jqs *JobQueueService) Pause() {
	atomic.StoreInt32(&jqs.paused, 1)
}

// Resume возобновляет выполнение задач, если они были приостановлены
func (jqs *JobQueueService) Resume() {
	if atomic.CompareAndSwapInt32(&jqs.paused, 1, 0) {
		close(jqs.resume)            // Сигнал для возобновления
		jqs.resume = make(chan struct{}) // Создание нового канала для следующей паузы
	}
}

// PauseAndResume приостанавливает выполнение задач на заданную длительность, затем возобновляет
func (jqs *JobQueueService) PauseAndResume(delay time.Duration) {
	jqs.Pause()
	time.AfterFunc(delay, func() {
		jqs.Resume()
	})
}

// Shutdown корректно завершает работу сервиса, ожидая завершения всех рабочих
func (jqs *JobQueueService) Shutdown() {
	close(jqs.jobs) // Закрытие канала задач
	jqs.wg.Wait()   // Ожидание завершения всех рабочих
}
