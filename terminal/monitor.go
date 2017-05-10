package main

import (
    "fmt"
    "math"
    "time"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
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

func RetrieveResults(results *[]string, session *mgo.Session)(err error) {
    users := make([]models.User, 0, 10)
    c := session.DB("gopher").C("users")
    err = c.Find(nil).Limit(10).Sort("-_id").All(&users)
    if err != nil {
        return err
    }
    for i:= range users {
        *results = append(*results, fmt.Sprintf("%d %s", users[i].Assigned, users[i].Name))
    }
    return nil
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

    session.SetMode(mgo.Eventual, true)
    session.SetSocketTimeout(0)
    defer session.Close()

    statCache := &models.StatCache{First:true} 

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

    g1 := termui.NewGauge()
    g1.Percent = 100
    g1.Width = 40
    g1.Height = 5 
    g1.X = 50 
    g1.Y = 0

    g2 := termui.NewGauge()
    g2.Percent = 100
    g2.Width = 40
    g2.Height = 5 
    g2.X = 50 
    g1.Y = 5

    g3 := termui.NewGauge()
    g3.Percent = 100
    g3.Width = 40
    g3.Height = 5
    g3.X = 50 
    g3.Y = 10

    ls := termui.NewList()
    ls.Items = []string{}
    ls.ItemFgColor = termui.ColorYellow
    ls.BorderLabel = "Latest 10"
    ls.Height = 11
    ls.Width = 50
    ls.Y = 15

    termui.Render(lc0, g1, g2, g3)
    
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
            gauges := []*termui.Gauge{g1, g2, g3}
            for i:=0; i<len(statCache.PreviousStat.Repl.Hosts); i++ {
                gauges[i].BorderLabel = statCache.PreviousStat.Repl.Hosts[i]
                gauges[i].BarColor = termui.ColorYellow
                if statCache.PreviousStat.Repl.Primary == statCache.PreviousStat.Repl.Hosts[i] {
                    gauges[i].BarColor = termui.ColorGreen
                }
            }
        }
        results := make([]string, 0, 10)
        err = RetrieveResults(&results, session)
        if err!=nil {
            fmt.Println(err)
            session.Refresh()
        }
        ls.Items = results
        
        termui.Render(lc0, g1, g2, g3, ls)
    })
    
    termui.Loop()
}
