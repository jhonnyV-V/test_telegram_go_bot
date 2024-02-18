package utils

import (
  "database/sql"
  "log"
  "errors"

  "github.com/jmoiron/sqlx"
  "github.com/mattn/go-sqlite3"
)

type Chat struct {
  ID          int64  `db:"id"`
  IsSuscribed bool   `db:"isSubscribed"`
  CreatedAt   string `db:"createdAt"`
  UpdatedAt   string `db:"updatedAt"`
}

type CustomListeer struct {
  ID        int64  `db:"id"`
  Module    string `db:"module"`
  Event     string `db:"event"`
  CreatedAt string `db:"createdAt"`
  UpdatedAt string `db:"updatedAt"`
}
type ChatToListener struct {
  ChatID     string `db:"user_id"`
  ListenerID int64  `db:"listenerId"`
  CreatedAt  string `db:"createdAt"`
}

var (
  ErrDuplicate    = errors.New("record already exists")
  ErrNotExists    = errors.New("row not exists")
  ErrUpdateFailed = errors.New("update failed")
  ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
  db *sqlx.DB
}

var DB *SQLiteRepository;

func OpenDB(driver string, datasource string) (*SQLiteRepository, error) {
  rawdb, err := sqlx.Connect(driver, datasource)
  if err != nil {
    return nil, err
  }
  DB = &SQLiteRepository{db: rawdb};
  return DB, nil
}

func (db *SQLiteRepository) Migrate() {
  setPragma := "PRAGMA foreign_keys;"
  createChatTable := `
    CREATE TABLE IF NOT EXISTS chat(
      id           INTEGER PRIMARY KEY NOT NULL, 
      isSubscribed BOOLEAN DEFAULT FALSE,
      createdAt    DATETIME DEFAULT CURRENT_DATE,
      updatedAt    DATETIME DEFAULT CURRENT_DATE
    );
  `
  createCustomListenerTable := `
    CREATE TABLE IF NOT EXISTS custom_listener(
      id        INTEGER PRIMARY KEY NOT NULL,
      module    TEXT NOT NULL,
      event     TEXT NOT NULL,
      createdAt DATETIME DEFAULT CURRENT_DATE,
      updatedAt DATETIME DEFAULT CURRENT_DATE,
      UNIQUE(module, event)
    );
  `
  createUserChatListenerTable := `
    CREATE TABLE IF NOT EXISTS chat_to_listener(
      chatId     INTEGER NOT NULL,
      listenerId INTEGER NOT NULL,
      createdAt  DATETIME DEFAULT CURRENT_DATE,
      UNIQUE(chatId, listenerId),
      FOREIGN KEY(chatId) REFERENCES chat(id) ON DELETE CASCADE,
      FOREIGN KEY(listenerId) REFERENCES custom_listener(id) ON DELETE CASCADE
    );
  `

  tx := db.db.MustBegin()

  tx.MustExec(setPragma)
  log.Println("setting pragma");
  
  tx.MustExec(createChatTable)
  log.Println("setting chat table");

  tx.MustExec(createCustomListenerTable)
  log.Println("setting custom_listener table");

  tx.MustExec(createUserChatListenerTable)
  log.Println("setting chat_to_listener table");

  err := tx.Commit();
  if err != nil {
    log.Println(err.Error());
    panic(err)
  }
}

func (db *SQLiteRepository) SaveChat(id int64) (error) {
  res, err := db.db.Exec("INSERT INTO chat(id) values(?)", id)
  if err != nil {
    var sqliteErr sqlite3.Error
    if errors.As(err, &sqliteErr) {
      if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
        return ErrDuplicate
      }
    }
    return err
  }
  _, err = res.LastInsertId()
  if err != nil {
    return err
  }

  return nil
}

func (db *SQLiteRepository) GetChatByID(id int64) (Chat, error) {
  var chat Chat
  err := db.db.Get(&chat, "SELECT id,isSubscribed FROM chat WHERE id=?", id)
  if errors.Is(err, sql.ErrNoRows) {
    return Chat{}, nil
  }
  if err != nil {
    return Chat{}, err
  }
  return chat, nil
}

//if the user was not subscribed and it updates the user data correctly
//it will return true else it will return false
func (db *SQLiteRepository) AddSubscription(id int64) (bool, error) {
  chat, err := db.GetChatByID(id)
  if err != nil {
    return false, nil
  }

  if chat == (Chat{}) {
    err = db.SaveChat(id)
    if err != nil {
      return false, err
    }
    chat.IsSuscribed = false
  }

  if chat.IsSuscribed {
    return false, nil 
  }

  _, err = db.db.Exec("UPDATE chat SET isSubscribed=TRUE WHERE id=?", id)

  if err != nil {
    return false, err
  }
  
  return true, nil
}

//if the user was subscribed and it updates the user data correctly
//it will return true else it will return false
func (db *SQLiteRepository) RemoveSubscription(id int64) (bool, error) {
  chat, err := db.GetChatByID(id)
  log.Println("id: ", id)
  if err != nil {
    log.Println("can't get chat")
    log.Println(err.Error())
    return false, nil
  }

  if chat == (Chat{}) {
    return false, nil
  }

  if !chat.IsSuscribed {
    return false, nil 
  }

  _, err = db.db.Exec("UPDATE chat SET isSubscribed=FALSE WHERE id=?", id)

  if err != nil {
    return false, err
  }
  
  return true, nil
}
