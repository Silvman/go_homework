package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

const numThs = 6 // Количество хешей на каждое входящее значение (для MultiHash)

func ExecutePipeline(jobs ...job) {
	in, out := make(chan interface{}), make(chan interface{})
	wg := &sync.WaitGroup{} // чтобы дождаться всех горутин с задачами

	for _, currentJob := range jobs {
		wg.Add(1)
		go func(group *sync.WaitGroup, currentJob job, in, out chan interface{}) { // запуск горутины с задачей
			defer group.Done()
			currentJob(in, out)
			close(out) // закрываем поток для остановки циклов
		}(wg, currentJob, in, out)

		in, out = out, make(chan interface{}) // выход становится новым входом
	}

	wg.Wait()
}

func asyncCrc32(data string, returnCh chan string) { // ассинхронная обертка над DataSignerCrc32
	returnCh <- DataSignerCrc32(data)
}

func asyncMd5(data string, returnCh chan string, mutex *sync.Mutex) { // ассинхронная обертка над DataSignerMd5
	mutex.Lock()
	returnCh <- DataSignerMd5(data)
	mutex.Unlock()
}

func asyncMd5Crc32(data string, returnCh chan string, mutex *sync.Mutex) {
	md5Ch := make(chan string) // канал для получения md5 хэша
	go asyncMd5(data, md5Ch, mutex)
	returnCh <- DataSignerCrc32(<-md5Ch)
}

func asyncTwoHashes(wg *sync.WaitGroup, data string, out chan interface{}, mutex *sync.Mutex) {
	defer wg.Done()

	crc32Ch, crc32md5Ch := make(chan string), make(chan string) // каналы для получения первого и второго хэшей
	go asyncCrc32(data, crc32Ch)
	go asyncMd5Crc32(data, crc32md5Ch, mutex)

	crc32Hash := <-crc32Ch
	crc32md5Hash := <-crc32md5Ch

	out <- crc32Hash + "~" + crc32md5Hash
}

func SingleHash(in, out chan interface{}) {
	mutex := &sync.Mutex{} // не более одной ф-ии md5 одновременно
	wg := &sync.WaitGroup{}

	for data := range in {
		strData := strconv.Itoa(data.(int)) // получаем строчку для хэширования
		wg.Add(1)
		go asyncTwoHashes(wg, strData, out, mutex)
	}

	wg.Wait()
}

func asyncSixHashes(wg *sync.WaitGroup, data string, out chan interface{}) {
	defer wg.Done()

	var hashChArray [numThs]chan string
	for i := range hashChArray {
		hashChArray[i] = make(chan string)                  // канал для приёма хэшей
		go asyncCrc32(strconv.Itoa(i)+data, hashChArray[i]) // параллельно запускаем по хэш-функции на каждое значение i
	}

	result := ""
	for _, ch := range hashChArray {
		result += <-ch
	}

	out <- result
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for data := range in {
		strData := data.(string)
		wg.Add(1)
		go asyncSixHashes(wg, strData, out)
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var hashes []string

	for data := range in {
		hashes = append(hashes, data.(string)) // собираем все полученные хэши
	}
	sort.Strings(hashes)

	out <- strings.Join(hashes, "_")
}
