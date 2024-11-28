package workerpool

import (
	"sync"
)

// WorkerPool - структура для пула рабочих
type WorkerPool struct {
	// Канал для управления количеством одновременно работающих горутин
	workerChan chan struct{}
	// Группа ожидания для ожидания завершения всех задач
	wg sync.WaitGroup
}

// NewWorkerPool - конструктор для создания нового пула рабочих
func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		workerChan: make(chan struct{}, workerCount),
	}
}

// Run - метод для выполнения задач с использованием пула рабочих
func (wp *WorkerPool) Run(task func()) {
	// Захватываем место в пуле
	wp.workerChan <- struct{}{}
	// Увеличиваем счетчик ожидания
	wp.wg.Add(1)

	go func() {
		defer wp.wg.Done() // Уменьшаем счетчик после выполнения задачи
		defer func() {
			// Освобождаем место в пуле
			<-wp.workerChan
		}()
		task()
	}()
}

// Wait - метод для ожидания завершения всех задач
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
