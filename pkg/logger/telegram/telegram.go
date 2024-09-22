package telegram

import (
	"bytes"
	"fmt"
	"go-clean-template/pkg/logger/common"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

const dtMask = time.RFC3339Nano

type Logs struct {
	mu sync.Mutex
	m  map[string][]string
}

func (l *Logs) Load(reqID string) ([]string, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	msgs, ok := l.m[reqID]
	return msgs, ok
}

func (l *Logs) Store(reqID string, logMsg string) {
	l.mu.Lock()
	if _, ok := l.m[reqID]; !ok {
		l.m[reqID] = make([]string, 0)
	}
	l.m[reqID] = append(l.m[reqID], logMsg)
	l.mu.Unlock()
}

func (l *Logs) Delete(reqID string) {
	l.mu.Lock()
	delete(l.m, reqID)
	l.mu.Unlock()
}

type Message struct {
	Header      string
	AppName     string
	Version     string
	Environment string
	InstanceID  string
	RequestID   string
	Stack       []string
}

func NewMessage(header, appName, version, env, instanceID, reqID string, logs []string) *Message {
	return &Message{header, appName, version, env, instanceID, reqID, logs}
}

func (m *Message) ToString() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("<b>%s</b>\n", m.Header))
	buf.WriteString(fmt.Sprintf("<b>AppName:</b> %s\n", m.AppName))
	buf.WriteString(fmt.Sprintf("<b>Version:</b> %s\n", m.Version))
	buf.WriteString(fmt.Sprintf("<b>Environment:</b> %s\n", m.Environment))
	buf.WriteString(fmt.Sprintf("<b>InstanceID:</b> %s\n", m.InstanceID))
	buf.WriteString(fmt.Sprintf("<b>Timestamp:</b> %s\n", time.Now().Format(dtMask)))
	if m.Header == "ERROR" || m.Header == "FATAL" {
		reqID := "-"
		if m.RequestID != "" {
			reqID = m.RequestID
		}
		buf.WriteString(fmt.Sprintf("<b>RequestID:</b> %s\n", reqID))
		buf.WriteString("<b>Log stack:</b>\n")
		for _, log := range m.Stack {
			buf.WriteString(fmt.Sprintf("<code>%s</code>\n", log))
		}
	}

	return buf.String()
}

type Logger struct {
	appName    string
	version    string
	env        string
	instanceID string
	logs       *Logs
	ch         chan string
	bot        *tgbotapi.BotAPI
	chatID     int64
}

type TelegramLoggerOpts struct {
	Enabled      bool
	Level        string
	TargetChatID int64
	BotAPIToken  string
}

func NewLogger(selfOpts *TelegramLoggerOpts, opts *common.GeneralOpts) *Logger {
	if selfOpts.TargetChatID == 0 {
		return nil
	}
	bot, err := tgbotapi.NewBotAPI(selfOpts.BotAPIToken)
	if err != nil {
		return nil
	}
	bot.Debug = false

	chBuf := 100
	l := &Logger{
		appName:    opts.AppName,
		version:    opts.AppVersion,
		env:        opts.Env,
		instanceID: opts.InstanceID.String(),
		logs:       &Logs{m: make(map[string][]string)},
		ch:         make(chan string, chBuf),
		bot:        bot,
		chatID:     selfOpts.TargetChatID,
	}

	go l.senderToChat()

	msg := NewMessage("STARTED", l.appName, l.version, l.env, l.instanceID, "", nil)
	l.ch <- msg.ToString()

	return l
}

func (l *Logger) Close() {
	msg := NewMessage("STOPPED", l.appName, l.version, l.env, l.instanceID, "", nil)
	l.ch <- msg.ToString()
	close(l.ch)
}

func (l *Logger) Debug(v ...interface{}) {
	ss := strings.Split(v[0].(string), " ")
	if _, err := uuid.Parse(ss[0]); err == nil {
		logMsg := fmt.Sprintf("%s DEBUG %s %s", time.Now().Format(dtMask), common.GetFuncName(),
			"["+strings.ReplaceAll(v[0].(string), ss[0]+" ", "")+"]")
		l.logs.Store(ss[0], logMsg)

		go func() {
			time.Sleep(30 * time.Second) //nolint:mnd //timeout 30s
			l.logs.Delete(ss[0])
		}()
	}
}

func (l *Logger) Info(v ...interface{}) {
	ss := strings.Split(v[0].(string), " ")
	if _, err := uuid.Parse(ss[0]); err == nil {
		logMsg := fmt.Sprintf("%s INFO %s %s", time.Now().Format(dtMask), common.GetFuncName(),
			"["+strings.ReplaceAll(v[0].(string), ss[0]+" ", "")+"]")
		l.logs.Store(ss[0], logMsg)

		go func() {
			time.Sleep(30 * time.Second) //nolint:mnd //timeout 30s
			l.logs.Delete(ss[0])
		}()
	}
}

func (l *Logger) Warning(v ...interface{}) {
	ss := strings.Split(v[0].(string), " ")
	if _, err := uuid.Parse(ss[0]); err == nil {
		logMsg := fmt.Sprintf("%s WARNING %s %s", time.Now().Format(dtMask), common.GetFuncName(),
			"["+strings.ReplaceAll(v[0].(string), ss[0]+" ", "")+"]")
		l.logs.Store(ss[0], logMsg)

		go func() {
			time.Sleep(30 * time.Second) //nolint:mnd //timeout 30s
			l.logs.Delete(ss[0])
		}()
	}
}

func (l *Logger) Error(v ...interface{}) {
	ss := strings.Split(v[0].(string), " ")
	if _, err := uuid.Parse(ss[0]); err == nil {
		logMsg := fmt.Sprintf("%s ERROR %s %s", time.Now().Format(dtMask), common.GetFuncName(),
			"["+strings.ReplaceAll(v[0].(string), ss[0]+" ", "")+"]")
		l.logs.Store(ss[0], logMsg)

		go func() {
			time.Sleep(time.Second)
			if logs, ok := l.logs.Load(ss[0]); ok {
				msg := NewMessage("ERROR", l.appName, l.version, l.env, l.instanceID, ss[0], logs)
				l.ch <- msg.ToString()
			}
		}()

		go func() {
			time.Sleep(30 * time.Second) //nolint:mnd //timeout 30s
			l.logs.Delete(ss[0])
		}()
	} else {
		logMsg := fmt.Sprintf("%s ERROR %s %s", time.Now().Format(dtMask), common.GetFuncName(), "["+v[0].(string)+"]")
		msg := NewMessage("ERROR", l.appName, l.version, l.env, l.instanceID, "", []string{logMsg})
		l.ch <- msg.ToString()
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	ss := strings.Split(v[0].(string), " ")
	if _, err := uuid.Parse(ss[0]); err == nil {
		logMsg := fmt.Sprintf("%s FATAL %s %s", time.Now().Format(dtMask), common.GetFuncName(),
			"["+strings.ReplaceAll(v[0].(string), ss[0]+" ", "")+"]")
		l.logs.Store(ss[0], logMsg)

		go func() {
			time.Sleep(time.Second)
			if logs, ok := l.logs.Load(ss[0]); ok {
				msg := NewMessage("FATAL", l.appName, l.version, l.env, l.instanceID, ss[0], logs)
				l.ch <- msg.ToString()
			}
		}()
	} else {
		logMsg := fmt.Sprintf("%s FATAL %s %s", time.Now().Format(dtMask), common.GetFuncName(), "["+v[0].(string)+"]")
		msg := NewMessage("FATAL", l.appName, l.version, l.env, l.instanceID, "", []string{logMsg})
		l.ch <- msg.ToString()
	}
}

func (l *Logger) senderToChat() {
	for msg := range l.ch {
		outMsg := tgbotapi.NewMessage(l.chatID, msg)
		outMsg.ParseMode = tgbotapi.ModeHTML
		_, _ = l.bot.Send(outMsg)
	}
}
