package utils

import (
    "dw/poker/lib"
    "dw/poker/context"
    "database/sql"
    "fmt"
)


var NoticeLog *lib.Logger
var AccessLog *lib.Logger
var DebugLog *lib.Logger
var WarningLog *lib.Logger
var FatalLog *lib.Logger
var Mailer *lib.Mailer
var sysMail lib.Mail


var MainDB *sql.DB

func Init(conf *context.Config) error {
    dir := conf.Log.LogDir
    rotate := conf.Log.LogRotate
    DebugLog   = lib.NewLogger(dir, "debug", rotate, true)
    WarningLog = lib.NewLogger(dir, "warning", rotate, true)
    FatalLog   = lib.NewLogger(dir, "fatal", rotate, true)
    AccessLog  = lib.NewLogger(dir, "access", rotate, false)
    NoticeLog  = lib.NewLogger(dir, "notice", rotate, false)
    if !conf.Server.Debug {
        DebugLog.Mute()
    }

    mailConf := lib.MailConfig{}
    mailConf.Username = conf.AlertMail.Username
    mailConf.Password = conf.AlertMail.Password
    mailConf.Host = conf.AlertMail.Host
    mailConf.Server = conf.AlertMail.Server
    mailConf.HostName = conf.Server.Hostname

    Mailer = lib.NewMailer(mailConf)
    sysMail = lib.Mail{}
    sysMail.Sender = conf.AlertMail.Sender
    sysMail.Receivers = conf.AlertMail.Receiver
    sysMail.Subject = conf.AlertMail.Subject

    c := conf.Sqldb.Main
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", c.Username, c.Password,
        c.Host, c.Port, c.Dbname, c.Charset)
    if len(c.Local) > 0 {
        dsn = fmt.Sprintf("%s&parseTime=true&loc=%s", dsn, c.Local)
    }
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    MainDB = db
    return nil
}

//content, subject, receivers...
func SendSysMail(args ...string) {
    m := sysMail
    if len(args) > 0 {
        m.Content = args[0]
    }
    if len(args) > 1 {
        m.Subject = args[1]
    }
    if len(args) > 2 {
        m.Receivers = append(m.Receivers, args[2:]...)
    }
    Mailer.Send(m)
}


