package pool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type PoolTask interface {
	// Execute запускает выполнение задачи и возвращает nil,
	// либо возникшую ошибку.
	Execute() error
	// OnFailure будет обрабатывать ошибки, возникшие в Execute(), то есть
	// пул должен вызывать OnFailure в случае, если Execute возвращает ошибку.
	OnFailure(error)
}

type WorkerPool interface {
	// Start подготавливает пул для обработки задач. Должен вызываться один раз
	// перед использованием пула. Очередные вызовы должны игнорироваться.
	Start()
	// Stop останавливает обработку в пуле. Должен вызываться один раз.
	// Очередные вызовы должны игнорироваться.
	Stop()
	// AddWork добавляет задачу для обработки пулом. Добавлять задачи
	// можно после вызова Start() и до вызова Stop().
	// Если на момент добавления в пуле нет
	// свободных ресурсов (очередь заполнена) -
	// эту функция ожидает их освобождения (либо вызова Stop).
	AddWork(PoolTask)
}

type MyPool struct {
	// из этого канала будем брать задачи для обработки
	tasks chan PoolTask
	// для синхронизации работы
	wg         sync.WaitGroup
	onceStart  sync.Once
	onceStop   sync.Once
	inStarting int32

	NumWorkers  int
	ChannelSize int
}

func (p *MyPool) Start() {
	p.onceStart.Do(func() {
		p.tasks = make(chan PoolTask, p.ChannelSize)
		atomic.StoreInt32(&p.inStarting, 1)
		// для ожидания завершения
		p.wg.Add(p.NumWorkers)
		for i := 0; i < p.NumWorkers; i++ {
			go func() {
				// забираем задачи из канала
				for w := range p.tasks {
					// и выполняем
					err := w.Execute()
					if err != nil {
						w.OnFailure(err)
					}
				}
				p.wg.Done()
			}()
		}
	})
}

func (p *MyPool) Stop() {
	p.onceStop.Do(func() {
		atomic.StoreInt32(&p.inStarting, 0)
		// закроем канал с задачами
		close(p.tasks)
		// дождемся завершения работы уже запущенных задач
		p.wg.Wait()
	})
}

func (p *MyPool) AddWork(pt PoolTask) {
	if atomic.LoadInt32(&p.inStarting) == 1 {
		// добавляем задачи в канал, из которого забирает работу пул
		p.tasks <- pt
	}
}

// NewWorkerPool возвращает новый пул
// numWorkers - количество воркеров
// channelSize - размер очереди ожидания
// В случае ошибок верните nil и описание ошибки
func NewWorkerPool(numWorkers int, channelSize int) (*MyPool, error) {
	if (numWorkers <= 0) || (channelSize < 0) {
		return nil, fmt.Errorf("bad params")
	}

	p := MyPool{
		ChannelSize: channelSize,
		NumWorkers:  numWorkers,
	}
	return &p, nil
}
