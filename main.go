package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func fetchDataDetail(dataChan chan<- string, wg *sync.WaitGroup) {
    defer wg.Done() // Signal bahwa goroutine selesai
    
    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts/1", nil)
    if err != nil {
        dataChan <- "Error creating request: " + err.Error()
        return
    }


    resp, err := client.Do(req)
    if err != nil {
        dataChan <- "Failed to fetch data: " + err.Error()
        return
    }
    
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        dataChan <- "Error reading response: " + err.Error()
        return
    }

    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        dataChan <- "Error parsing JSON: " + err.Error()
        return
    }

    dataChan <- result["title"].(string)
}

func fetchDataArray(dataChan chan<- string, wg *sync.WaitGroup) {
    defer wg.Done() // Signal bahwa goroutine selesai
    
    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts", nil)
    if err != nil {
        dataChan <- "Error creating request: " + err.Error()
        return
    }

    resp, err := client.Do(req)
    if err != nil {
        dataChan <- "Failed to fetch data: " + err.Error()
        return
    }
    
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        dataChan <- "Error reading response: " + err.Error()
        return
    }

    var result []map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        dataChan <- "Error parsing JSON: " + err.Error()
        return
    }

    if len(result) > 0 {
        jsonResult, err := json.Marshal(result)
        if err != nil {
            dataChan <- "Error converting data to JSON: " + err.Error()
        } else {
            dataChan <- string(jsonResult)
        }
    } else {
        dataChan <- "No posts found"
    }
}

func fetchDataWithLoading(dataChan chan<- string, loadingChan chan<- bool) {
    loadingChan <- true
    
    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts", nil)
    if err != nil {
        dataChan <- "Error creating request"
        loadingChan <- false
        return
    }

    resp, err := client.Do(req)
    if err != nil {
        dataChan <- "Failed to fetch data"
        loadingChan <- false
        return
    }
    
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)

    if err != nil {
        dataChan <- "Error reading response"
        loadingChan <- false
        return
    }

    var result []map[string]interface{} // Changed to slice of maps
    err = json.Unmarshal(body, &result)
    if err != nil {
        dataChan <- "Error parsing JSON"
        loadingChan <- false
        return
    }

    // Ambil judul dari post pertama
    if len(result) > 0 {
        title, ok := result[0]["title"].(string)
        if ok {
            dataChan <- title
        } else {
            dataChan <- "Title not found"
        }
    } else {
        dataChan <- "No posts found"
    }
    
    loadingChan <- false
}

type AppState struct {
    Count int
}
func updateData(data *AppState) *AppState{
     data.Count++
     return data
}

func fetchData(dataChan chan<- string, wg *sync.WaitGroup) {
    defer wg.Done()
    
    client := &http.Client{}
    req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts", nil)

    if err != nil {
        dataChan <- "Error creating request"
        return
    }

    resp, err := client.Do(req)

    if err != nil {
        dataChan <- "Failed to fetch data"
        return
    }
    
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)

    if err != nil {
        dataChan <- "Error reading response"
        return
    }

    var result []map[string]interface{} // Changed to slice of maps

    err = json.Unmarshal(body, &result)
    if err != nil {
        dataChan <- "Error parsing JSON"
        return
    }

    // Ambil judul dari post pertama
    if len(result) > 0 {
        title, ok := result[0]["title"].(string)
        if ok {
            dataChan <- title
        } else {
            dataChan <- "Title not found"
        }
    } else {
        dataChan <- "No posts found"
    }
}

type Person struct{
    Age int
}
func modif(p *Person) *Person{
    p.Age = 123
    return p
}

func main() {
    var wg sync.WaitGroup
    dataChannelDetail := make(chan string, 1) // Buffered channel
    dataChannels := make(chan string, 1)      // Buffered channel
    dataChannel := make(chan string, 1)       // Buffered channel;
    fmt.Println("Loading...")

    // Start kedua goroutines
    wg.Add(3) // Kita akan menunggu 3 goroutines
    go fetchDataDetail(dataChannelDetail, &wg)
    go fetchDataArray(dataChannels, &wg)
    go fetchData(dataChannel, &wg)

    // Wait sampai semua selesai di background
    go func() {
        wg.Wait()
        // close(dataChannelDetail)
        // close(dataChannels)
        close(dataChannel)
    }()

    // Read results
    // dataDetail := <-dataChannelDetail
    // fmt.Println("dataDetail:", dataDetail)

    // dataCollection := <-dataChannels
    // fmt.Println("dataCollection:", dataCollection)

    singleCol := <-dataChannel
    fmt.Println("singleCollection:", singleCol)

    fmt.Println("End Loading!")

    countBags := AppState{
        Count: 0,
    }
    
    update := updateData(&countBags)
    fmt.Println(update.Count)
    println(countBags.Count)

}