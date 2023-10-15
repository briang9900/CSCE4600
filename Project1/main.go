package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "log"
    "os"
    "strconv"
    "strings"
    "math"
    "github.com/olekukonko/tablewriter"
)

func main() {
    // CLI args
    f, closeFile, err := openProcessingFile(os.Args...)
    if err != nil {
        log.Fatal(err)
    }
    defer closeFile()

    // Load and parse processes
    processes, err := loadProcesses(f)
    if err != nil {
        log.Fatal(err)
    }

    // First-come, first-serve scheduling
    FCFSSchedule(os.Stdout, "First-come, first-serve", processes)
    // Shortest Job First (SJF)
    SJFSchedule(os.Stdout, "Shortest-job-first", processes)
   // SJF Priority
    SJFPrioritySchedule(os.Stdout, "Priority", processes)
   // Round-robin (RR)
    RRSchedule(os.Stdout, "Round-robin", processes)
}

func openProcessingFile(args ...string) (*os.File, func(), error) {
    if len(args) != 2 {
        return nil, nil, fmt.Errorf("%w: must give a scheduling file to process", ErrInvalidArgs)
    }
    // Read in CSV process CSV file
    f, err := os.Open(args[1])
    if err != nil {
        return nil, nil, fmt.Errorf("%v: error opening scheduling file", err)
    }
    closeFn := func() {
        if err := f.Close(); err != nil {
            log.Fatalf("%v: error closing scheduling file", err)
        }
    }
    return f, closeFn, nil
}

type (
    Process struct {
        ProcessID     int64
        ArrivalTime   int64
        BurstDuration int64
        Priority      int64
    }
    RunTime struct {
        ProcessID int64
        waitTime  int64
        remainTime int64
    }
       TimeSlice struct {
        PID   int64
        Start int64
        Stop  int64
    }
)
//region Schedulers
// FCFSSchedule outputs a schedule of processes in a GANTT chart and a table of timing given:
// • an output writer
// • a title for the chart
// • a slice of processes
func FCFSSchedule(w io.Writer, title string, processes []Process) {
    var (
        serviceTime     int64
        totalWait       float64
        turnAroundTotal float64
        lasttimeComplete  float64
        tWait     int64
        schedule        = make([][]string, len(processes))
        gantt           = make([]TimeSlice, 0)
    )
    for i := range processes {
        if processes[i].ArrivalTime > 0 {
            tWait = serviceTime - processes[i].ArrivalTime
        }
        totalWait += float64(tWait)

        start := tWait + processes[i].ArrivalTime

        turnaround := processes[i].BurstDuration + tWait
        turnAroundTotal += float64(turnaround)

        timeComplete := processes[i].BurstDuration + processes[i].ArrivalTime + tWait
        lasttimeComplete = float64(timeComplete)

		schedule[i] = []string{
            fmt.Sprint(processes[i].ProcessID),
            fmt.Sprint(processes[i].Priority),
            fmt.Sprint(processes[i].BurstDuration),
            fmt.Sprint(processes[i].ArrivalTime),
            fmt.Sprint(tWait),
            fmt.Sprint(turnaround),
            fmt.Sprint(timeComplete),
        }
        serviceTime += processes[i].BurstDuration

		gantt = append(gantt, TimeSlice{
            PID:   processes[i].ProcessID,
            Start: start,
            Stop:  serviceTime,
        })
    }

    count := float64(len(processes))
    aveWait := totalWait / count
    aveTurnaround := turnAroundTotal / count
    aveThroughput := count / lasttimeComplete

	outputTitle(w, title)
    outputGantt(w, gantt)
    outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

//Implement SJF priority scheduling (preemptive) and report average turnaround time, average waiting time, and average throughput.
func SJFPrioritySchedule(w io.Writer, title string, processes []Process) {
    var (
        totalWait       float64
        turnAroundTotal float64
        lasttimeComplete  float64
        tWait     int64
        totalTime       int64
        activeProc      int64
        highestAvail    int64
        processStart    int64
        trackSJF = make([]RunTime, len(processes))
        schedule = make([][]string, len(processes))
        gantt = make([]TimeSlice, 0)
    )
      
    //This calculate total needed time and poluate the time tracker
    for i := range processes {
        trackSJF[i].ProcessID = processes[i].ProcessID
        trackSJF[i].waitTime = 0
        trackSJF[i].remainTime = processes[i].BurstDuration
        totalTime += processes[i].BurstDuration
    }

    activeProc = 0
    processStart = 0
    highestAvail = getShortest(trackSJF, processes, 0)
    activeProc = highestAvail
   
	for t := 0; t <= int(totalTime); t++ {
        for i := range processes {
            //Checks if higher priority process arrives
            if (processes[i].Priority < processes[activeProc].Priority) && (t == int(processes[i].ArrivalTime)) {
                highestAvail = int64(i)
            }
            //Checks an equal priorty process arrives and it is shorter
            if (processes[i].Priority == processes[activeProc].Priority) && (t == int(processes[i].ArrivalTime)) && (trackSJF[i].remainTime < trackSJF[activeProc].remainTime) {
                highestAvail = int64(i)
            }
            //Will increment wait time if process has arrived & not executing
            if (i != int(activeProc) && i != int(highestAvail)) && (trackSJF[i].remainTime > 0) && (t > int(processes[i].ArrivalTime)) {
                trackSJF[i].waitTime += 1
            }
            //Check if running process has completed
            if (i == int(activeProc)) && (trackSJF[i].remainTime == 0) {
                highestAvail = getHighest(trackSJF, processes, int64(t))
            }
            //Decrement the running process remainTime
            if i == int(activeProc) {
                trackSJF[i].remainTime -= 1
            }
        }
        //Update Gannt chart
        if activeProc != highestAvail || t == int(totalTime) {
            gantt = append(gantt, TimeSlice{
                PID:   processes[activeProc].ProcessID,
                Start: processStart,
                Stop:  int64(t),
            })
            processStart = int64(t)
            activeProc = highestAvail
        }
    }
    //Record final results
    for i := range processes {
		tWait = trackSJF[i].waitTime
        totalWait += float64(tWait)
        turnaround := processes[i].BurstDuration + trackSJF[i].waitTime
        turnAroundTotal += float64(turnaround)
        timeComplete := processes[i].BurstDuration + processes[i].ArrivalTime + tWait
        lasttimeComplete = float64(timeComplete)
        schedule[i] = []string{
            fmt.Sprint(processes[i].ProcessID),
            fmt.Sprint(processes[i].Priority),
            fmt.Sprint(processes[i].BurstDuration),
            fmt.Sprint(processes[i].ArrivalTime),
            fmt.Sprint(tWait),
            fmt.Sprint(turnaround),
            fmt.Sprint(timeComplete),
        }
    }
    //Calculate average and output results
    count := float64(len(processes))
    aveWait := totalWait / count
    aveTurnaround := turnAroundTotal / count
    aveThroughput := count / lasttimeComplete
    outputTitle(w, title)
    outputGantt(w, gantt)
    outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}
//Implement SJF (preemptive) and report average turnaround time, average waiting time, and average throughput.
func SJFSchedule(w io.Writer, title string, processes []Process) {
    var (
        serviceTime     int64
        lastStart       int64
		timeLess         int64
        totalWait       float64
        turnAroundTotal float64
        lasttimeComplete  float64  
        schedule        = make([][]string, len(processes))
        gantt           = make([]TimeSlice, 0)
        tWait           = make([]int64, len(processes))
        turnAroundTime  = make([]int64, len(processes))
        tRemain         = make([]int64, len(processes))
        timeComplete      = make([]int64, len(processes))
    )
    completed := 0
    shortest := 0
    lastShortest := 0
    timeLess = math.MaxInt64
    check := false
    count := len(processes)
   
    for i := range processes {
                tRemain[i] = processes[i].BurstDuration
    }
    for completed != count {
        for j := 0; j < count; j++ {
            if processes[j].ArrivalTime <= serviceTime && tRemain[j] < timeLess && tRemain[j] > 0 {
                timeLess = tRemain[j]
                shortest = j
                check = true
            }
        }
        // Update Gantt schedule
        if shortest != lastShortest {
            gantt = append(gantt, TimeSlice{
                PID:   processes[lastShortest].ProcessID,
                Start: lastStart,
                Stop:  serviceTime,
            })
            lastStart = serviceTime
            lastShortest = shortest
        }
        if check == false {
            serviceTime++
            continue
        }
        tRemain[shortest]--
        timeLess = tRemain[shortest]
        if (timeLess == 0) {
            timeLess = math.MaxInt64
        }
        if tRemain[shortest] == 0 {
            completed++
            check = false
            timeComplete[shortest] = serviceTime + 1
            lasttimeComplete = float64(timeComplete[shortest])
            tWait[shortest] = timeComplete[shortest] - processes[shortest].BurstDuration - processes[shortest].ArrivalTime
            if tWait[shortest] < 0 {
                tWait[shortest] = 0
            }
        }
        serviceTime++
    }
    //Calculate total wait and turn around time
    for i := range tWait {
        totalWait += float64(tWait[i])
        turnAroundTime[i] = processes[i].BurstDuration + tWait[i]
        turnAroundTotal += float64(turnAroundTime[i])
    }
    //Calculate average
    aveWait := totalWait / float64(count)
    aveTurnaround := turnAroundTotal / float64(count)
    aveThroughput := float64(count) / lasttimeComplete
   
    for i := range processes {
        schedule[i] = []string{
            fmt.Sprint(processes[i].ProcessID),
            fmt.Sprint(processes[i].Priority),
            fmt.Sprint(processes[i].BurstDuration),
            fmt.Sprint(processes[i].ArrivalTime),
            fmt.Sprint(tWait[i]),
            fmt.Sprint(turnAroundTime[i]),
            fmt.Sprint(timeComplete[i]),
        }
    }
    //Add last entry and generate results
    gantt = append(gantt, TimeSlice{
        PID:   processes[lastShortest].ProcessID,
        Start: lastStart,
        Stop:  serviceTime,
    })
    outputTitle(w, title)
    outputGantt(w, gantt)
    outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

//Round-round (preemptive) and report average turnaround time, average waiting time, and average throughput.
func RRSchedule(w io.Writer, title string, processes []Process) {
    var (
        serviceTime     int64
        lastStart       int64
        timeQuantum     int64
        totalWait       float64
        turnAroundTotal float64
        lasttimeComplete  float64  
        schedule = make([][]string, len(processes))
        gantt = make([]TimeSlice, 0)
        turnAroundTime = make([]int64, len(processes))
        tWait = make([]int64, len(processes))
        tRemain = make([]int64, len(processes))
        timeComplete = make([]int64, len(processes))
    )
    timeQuantum = 1
    stuck := 0
    completed := 0
    turn := 0
    count := len(processes)
    check := false
   
    for i := range processes {
                tRemain[i] = processes[i].BurstDuration
    }

    for completed != count {
        if processes[turn].ArrivalTime > serviceTime || tRemain[turn] == 0 {
            turn = (turn + 1) % count
            if check == false {
                check = true
                stuck = turn
            } else if stuck == turn {
                serviceTime++
                lastStart = serviceTime
                check = false
                turn = 0
            }
            continue
        }
        check = false
        if tRemain[turn] > timeQuantum {
            serviceTime += timeQuantum
            tRemain[turn] -= timeQuantum
        } else {
            serviceTime += tRemain[turn]
            tRemain[turn] = 0
            completed++
            timeComplete[turn] = serviceTime
            lasttimeComplete = float64(timeComplete[turn])
            tWait[turn] = timeComplete[turn] - processes[turn].BurstDuration - processes[turn].ArrivalTime
            if tWait[turn] < 0 {
                tWait[turn] = 0
            }  
        }
        gantt = append(gantt, TimeSlice{
            PID:   processes[turn].ProcessID,
            Start: lastStart,
            Stop:  serviceTime,
        })
        lastStart = serviceTime
        turn = (turn + 1) % count
    }
    //Calculate total wait and turn around time
    for i := range tWait {
        totalWait += float64(tWait[i])
        turnAroundTime[i] = processes[i].BurstDuration + tWait[i]
        turnAroundTotal += float64(turnAroundTime[i])
    }
    aveWait := totalWait / float64(count)
    aveTurnaround := turnAroundTotal / float64(count)
    aveThroughput := float64(count) / lasttimeComplete

    for i := range processes {
        schedule[i] = []string{
            fmt.Sprint(processes[i].ProcessID),
            fmt.Sprint(processes[i].Priority),
            fmt.Sprint(processes[i].BurstDuration),
            fmt.Sprint(processes[i].ArrivalTime),
            fmt.Sprint(tWait[i]),
            fmt.Sprint(turnAroundTime[i]),
            fmt.Sprint(timeComplete[i]),
        }
    }
    //Generate results
    outputTitle(w, title)
    outputGantt(w, gantt)
    outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)  
}

// Returns the index of the shortest job that has an arrival time at or before the specified current time
func getShortest(tracker []RunTime, processes []Process, current int64) (shortest int64) {
    shortest = 0

    for i := range processes {
        if tracker[shortest].remainTime <= 0 {
            shortest += 1
            continue
        }
        if (tracker[i].remainTime < tracker[shortest].remainTime) && (processes[i].ArrivalTime <= current) && (tracker[i].remainTime > 0) {
            shortest = int64(i)
        }
    }
    return
}
// Returns index of the highest priortity and shortest job avaiable at the current clock time
func getHighest(tracker []RunTime, processes []Process, current int64) (highest int64) {
    highest = 0

    for i := range processes {
        if tracker[highest].remainTime <= 0 {
            highest += 1
            continue
        }
        if (processes[i].Priority < processes[highest].Priority) && (processes[i].ArrivalTime <= current) && (tracker[i].remainTime > 0) {
            highest = int64(i)
        }
    }
    return
}
//region Output helpers
func outputTitle(w io.Writer, title string) {
    _, _ = fmt.Fprintln(w, strings.Repeat("-", len(title)*2))
    _, _ = fmt.Fprintln(w, strings.Repeat(" ", len(title)/2), title)
    _, _ = fmt.Fprintln(w, strings.Repeat("-", len(title)*2))
}

func outputGantt(w io.Writer, gantt []TimeSlice) {
    _, _ = fmt.Fprintln(w, "Gantt schedule")
    _, _ = fmt.Fprint(w, "|")
    for i := range gantt {
        pid := fmt.Sprint(gantt[i].PID)
        padding := strings.Repeat(" ", (8-len(pid))/2)
        _, _ = fmt.Fprint(w, padding, pid, padding, "|")
    }
    _, _ = fmt.Fprintln(w)
    for i := range gantt {
        _, _ = fmt.Fprint(w, fmt.Sprint(gantt[i].Start), "\t")
        if len(gantt)-1 == i {
            _, _ = fmt.Fprint(w, fmt.Sprint(gantt[i].Stop))
        }
    }
    _, _ = fmt.Fprintf(w, "\n\n")
}

func outputSchedule(w io.Writer, rows [][]string, wait, turnaround, throughput float64) {
    _, _ = fmt.Fprintln(w, "Schedule table")
    table := tablewriter.NewWriter(w)
    table.SetHeader([]string{"ID", "Priority", "Burst", "Arrival", "Wait", "Turnaround", "Exit"})
    table.AppendBulk(rows)
    table.SetFooter([]string{"", "", "", "",
        fmt.Sprintf("Average\n%.2f", wait),
        fmt.Sprintf("Average\n%.2f", turnaround),
        fmt.Sprintf("Throughput\n%.2f/t", throughput)})
    table.Render()
}
//endregion
//region Loading processes.
var ErrInvalidArgs = errors.New("invalid args")

func loadProcesses(r io.Reader) ([]Process, error) {
    rows, err := csv.NewReader(r).ReadAll()
    if err != nil {
        return nil, fmt.Errorf("%w: reading CSV", err)
    }
	processes := make([]Process, len(rows))
    for i := range rows {
        processes[i].ProcessID = mustStrToInt(rows[i][0])
        processes[i].BurstDuration = mustStrToInt(rows[i][1])
        processes[i].ArrivalTime = mustStrToInt(rows[i][2])
        if len(rows[i]) == 4 {
            processes[i].Priority = mustStrToInt(rows[i][3])
        }
    }
    return processes, nil
}

func mustStrToInt(s string) int64 {
    i, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        _, _ = fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    return i
}
//endregion
