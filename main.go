package main

import (
  "database/sql"
  "fmt"
  "os/exec"
  "time"
  "github.com/gin-gonic/gin"
  _ "github.com/go-sql-driver/mysql"
)

type DeadStatusEntry struct {
  Name      string  `db:"name"`
  Status    string  `db:"status"`
  Timestamp string  `db:"timestamp"`
}

func db_open () *sql.DB {
  db, err := sql.Open("mysql", "dead_core:@/dead")
  if err != nil {
    panic(err.Error())
  }
  return db
}

func get_status (c *gin.Context) {
  db := db_open()
  defer db.Close()
  _, _ = db.Query("UPDATE dead_status SET status = 'alive', timestamp = NOW()")

  rows, _ := db.Query("SELECT * FROM dead_status")
  defer rows.Close()

  entry := DeadStatusEntry{}
  for rows.Next() {
    _ = rows.Scan(&entry.Name, &entry.Status, &entry.Timestamp)
  }

  c.JSON(200, gin.H{
    "name"    : entry.Name,
    "status"  : entry.Status,
    "timestamp" : entry.Timestamp,
  })
}

func only_get_status (c *gin.Context) {
  db := db_open()
  defer db.Close()

  rows, _ := db.Query("SELECT * FROM dead_status")

  entry := DeadStatusEntry{}
  for rows.Next() {
    _ = rows.Scan(&entry.Name, &entry.Status, &entry.Timestamp)
  }

  c.JSON(200, gin.H{
    "name"    : entry.Name,
    "status"  : entry.Status,
    "timestamp" : entry.Timestamp,
  })
}

func update_status (c *gin.Context) {
  db := db_open()
  defer db.Close()
  _, _ = db.Query("UPDATE dead_status SET status = 'alive', timestamp = NOW()")

  c.JSON(200, gin.H{
    "status":  "success",
  })
}

func change_danger (c *gin.Context) {
  db := db_open()
  defer db.Close()
  _, _ = db.Query("UPDATE dead_status SET status = 'danger', timestamp = NOW()")
  // TODO send mail

  go func() {
    time.Sleep(1 * time.Second) // 通常時 24 * time.Hour

    db := db_open()
    defer db.Close()
    rows, _ := db.Query("SELECT * FROM dead_status")
    defer rows.Close()

    entry := DeadStatusEntry{}
    for rows.Next() {
      _ = rows.Scan(&entry.Name, &entry.Status, &entry.Timestamp)
    }
    if (entry.Status == "alive") {
      return
    }

    _, _ = db.Query("UPDATE dead_status SET status = 'dead', timestamp = NOW()")
    out, _ := exec.Command("echo", "you died.").Output()
    fmt.Println(string(out))
  }()

  c.JSON(200, gin.H{
    "status":  "success",
  })
}

func main() {
  router := gin.Default()

  router.GET("/", get_status);
  router.GET("/status", get_status);
  router.GET("/only_get_status", only_get_status);

  router.GET("/update_status", update_status);
  router.POST("/update_status", update_status);
  router.PUT("/update_status", update_status);

  router.GET("/change_danger", change_danger);
  router.POST("/change_danger", change_danger);
  router.PUT("/change_danger", change_danger);

  router.Run(":8080")
}
