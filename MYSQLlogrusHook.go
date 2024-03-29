package MYSQLlogrusHook

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type hooker struct {
	db *sql.DB
	t  string
}

func NewHooker(db *sql.DB, table string, truncate bool) (*hooker, error) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `" + table + "` (`id` bigint unsigned NOT NULL AUTO_INCREMENT,`session_id` bigint unsigned,`level` VARCHAR(10) NOT NULL,`time` DATETIME(6) NOT NULL,`message` LONGTEXT,PRIMARY KEY (`id`),KEY `time` (`time`),KEY `level` (`level`),KEY `session_id` (`session_id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;")

	if err != nil {
		return nil, err
	}
	_, err = db.Exec("SET NAMES utf8mb4;")

	if err != nil {
		return nil, err
	}

	if truncate {
		_, err = db.Exec("TRUNCATE `" + table + "`;")
		if err != nil {
			return nil, err
		}
	}

	return &hooker{db: db, t: table}, nil
}

func (h *hooker) Fire(entry *logrus.Entry) error {
	_, err := h.db.Exec("INSERT INTO `"+h.t+"` (`level`,`time`,`session_id`,`message`) VALUES (?,?,?,?);", entry.Level.String(), entry.Time, entry.Data["session_id"], entry.Message)
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
func (h *hooker) Close() {
	h.db.Close()
}
