package MYSQLlogrusHook


import (
    "github.com/Sirupsen/logrus"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "fmt"
)

type hooker struct{
    db *sql.DB
    t string
}

func NewHooker(conn, table string) (*hooker, error){
    db, err := sql.Open("mysql",conn)
    if err != nil {
        return nil, err
    }
    _, err = db.Exec("CREATE TABLE IF NOT EXISTS `"+table+"` (`id` bigint unsigned NOT NULL AUTO_INCREMENT,`level` VARCHAR(10) NOT NULL,`time` DATETIME NOT NULL,`message` LONGTEXT,PRIMARY KEY (`id`),KEY `time` (`time`),KEY `level` (`level`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;")
    return &hooker{db: db, t: table}, nil
}

func (h *hooker) Fire(entry *logrus.Entry) error {
    _, err := h.db.Exec("INSERT INTO `"+h.t+"` (`level`,`time`,`message`) VALUES (?,?,?)",entry.Level.String(), entry.Time, entry.Message)
    if err != nil {
        return fmt.Errorf("Failed to send log entry to mysqldb: %v", err)
    }
	return nil
}

func (h *hooker) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
func (h *hooker) Close(){
    h.db.Close()
}