package main

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "fmt"
    "math"
    "time"
    "github.com/gizak/termui"
    "github.com/sindbach/gomongo/models"
)


func StatProcess(statCacheOp []float64, currentStatOp int64, previousStatOp int64, steep int, limit int) (newStatCacheOp []float64, err error){
    if len(statCacheOp) > limit {
        statCacheOp  = statCacheOp[1:]
    } else {
        statCacheOp = statCacheOp
    }
    prevCount := float64(previousStatOp)
    if prevCount != 0 {
        prevStat := float64(statCacheOp[len(statCacheOp)-1])
        diff := float64(currentStatOp-int64(prevCount))
        increment := float64((diff-prevStat)/float64(steep))
        smoothers := []float64{diff}
        if math.Abs(diff-prevStat) > 1 && diff != prevStat{
            smoothers = []float64{}
            for i:=1; i<steep; i++ {
                smoothers = append(smoothers, prevStat+(increment*float64(i)))
            }
            smoothers = append(smoothers, diff)
        } 
        newStatCacheOp = append(statCacheOp,  smoothers...)
    } else {
        newStatCacheOp = append(statCacheOp, float64(currentStatOp-previousStatOp))
    }
    return newStatCacheOp, nil
}

func RetrieveStats(statCache *models.StatCache, session *mgo.Session) (err error){
    steep := 6
    limit := 40
    stat := &models.ServerStatus{}
    if statCache.First == true {
        err := session.DB("admin").Run(bson.D{{"serverStatus", 1}, {"recordStats", 0}}, stat)
        if err != nil {
            return err
        }
        statCache.OpCommands = append(statCache.OpCommands, 0)
        statCache.OpInserts = append(statCache.OpInserts, 0)
        statCache.PreviousStat = stat
        statCache.First = false
    }
    err = session.DB("admin").Run(bson.D{{"serverStatus", 1}, {"recordStats", 0}}, stat)
    if err != nil {
        return err
    }

    statCache.OpInserts, err = StatProcess(statCache.OpInserts, stat.Opcounters.Insert, statCache.PreviousStat.Opcounters.Insert, steep, limit)
    if err != nil {
        return err
    }

    statCache.OpCommands, err = StatProcess(statCache.OpCommands, stat.Opcounters.Command, statCache.PreviousStat.Opcounters.Command, steep, limit)
    if err != nil {
        return err
    }

    //fmt.Println("current commands:", stat.Opcounters.Command)
    //fmt.Println("inside commands: ", statCache.OpCommands)
    //fmt.Println("current inserts:", stat.Opcounters.Insert)
    //fmt.Println("primary: ", stat.Repl.Primary)
    statCache.PreviousStat = stat
    return nil
} 

func main() {
    if err := termui.Init(); err != nil {
        fmt.Printf("Error setting up terminal UI: %v", err)
        panic("could not set up termui terminal interface")
    }
    defer termui.Close()

    session, err := mgo.Dial("localhost:27000,localhost:27001,localhost:27002")
    if err != nil {
        fmt.Println(err)
        panic(err)
    }

    session.SetMode(mgo.Monotonic, true)
    session.SetSocketTimeout(0)
    defer session.Close()

    statCache := &models.StatCache{First:true} 

    /*for i:=0;i<1000;i++ {
        err = RetrieveStats(statCache, session)
        if err!=nil{
            panic(err)
        }
        fmt.Println(statCache.OpInserts)

        time.Sleep(300 * time.Millisecond)
    }*/

    lc0 := termui.NewLineChart()
    lc0.BorderLabel = "OpCounter Commands"
    lc0.Mode = "dot"
    lc0.Data = statCache.OpCommands
    lc0.LineColor = termui.ColorGreen

    lc0.Width = 50
    lc0.Height = 12
    lc0.X = 0
    lc0.Y = 0
    lc0.AxesColor = termui.ColorWhite
    lc0.LineColor = termui.ColorGreen | termui.AttrBold

    bc1 := termui.NewBarChart()
    bc1.Data = []int{3, 2, 5}
    bc1.DataLabels = []string{"localhost:27000", "localhost:27001", "localhost:27002"}
    bc1.BorderLabel = "Status"
    bc1.Width = 40
    bc1.Height = 10
    bc1.TextColor = termui.ColorGreen
    bc1.BarColor = termui.ColorGreen
    bc1.NumColor = termui.ColorYellow
    bc1.X = 50
    bc1.Y = 0


    termui.Render(lc0, bc1)
    
    termui.Handle("/sys/kbd/q", func(termui.Event) {
        termui.StopLoop()
    })

    termui.Merge("timer", termui.NewTimerCh(300*time.Millisecond))

    termui.Handle("/timer/300ms", func(e termui.Event) {
        err = RetrieveStats(statCache, session)
        if err!=nil {
            fmt.Println(err)
            session.Refresh()
        } else {
            lc0.Data = statCache.OpCommands
            //bc1.Data = statCache.OpInserts
            //fmt.Println("outside: ", statCache.OpInserts)
            termui.Render(lc0, bc1)
        }
    })
    
    termui.Loop()
}
