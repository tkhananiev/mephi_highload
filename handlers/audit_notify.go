package handlers

import (
    "fmt"
    "time"

    "go-microservice/utils"
)

type Audit struct {
    ch chan string
}

func NewAudit(buffer int) *Audit {
    a := &Audit{ch: make(chan string, buffer)}
    go func() {
        for msg := range a.ch {
            utils.Logger.Printf("[AUDIT] %s", msg)
        }
    }()
    return a
}

func (a *Audit) Log(action string, userID int) {
    msg := fmt.Sprintf("%s user_id=%d ts=%s", action, userID, time.Now().UTC().Format(time.RFC3339Nano))
    select {
    case a.ch <- msg:
    default:
        go func() { a.ch <- msg }()
    }
}

type Notifier struct {
    ch chan string
}

func NewNotifier(buffer int) *Notifier {
    n := &Notifier{ch: make(chan string, buffer)}
    go func() {
        for msg := range n.ch {
            utils.Logger.Printf("[NOTIFY] %s", msg)
        }
    }()
    return n
}

func (n *Notifier) Send(event string, userID int) {
    msg := fmt.Sprintf("%s user_id=%d ts=%s", event, userID, time.Now().UTC().Format(time.RFC3339Nano))
    select {
    case n.ch <- msg:
    default:
        go func() { n.ch <- msg }()
    }
}
