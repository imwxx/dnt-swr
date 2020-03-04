package lib
import (
    "fmt"
    "database/sql"
    _"github.com/go-sql-driver/mysql"
)

var MDBC *sql.DB

type RESPONSEACTION struct {
    Status bool
    Msg string
    SyncStatus string
}

func check(err error) {
    if err !=nil {
        fmt.Println(err)
    }
}
func InitDb(dbInfo string) (*sql.DB, error){
    db, err := sql.Open("mysql", dbInfo)
    if err != nil {
        return db, err
    }
//    defer db.Close()
    return db, nil
}

func MQueryRow(harborRepo string, db *sql.DB) (bool){
    var dbRsp bool
    rows, err := db.Query("SELECT * from syncpool where harborRepo = ?", harborRepo)
    check(err)
    count := 0
    for rows.Next() {
        count += 1
    }
    if count == 0 {
        dbRsp = false
    } else {
        dbRsp = true
    }
    return dbRsp
}

func Insert(harborRepo, swrRepo, syncStatus string, db *sql.DB) (RESPONSEACTION){
    var dbRsp RESPONSEACTION
    count := MQueryRow(harborRepo, db)
    if !count {
        _, err := db.Exec(
            "INSERT INTO syncpool (harborRepo, swrRepo, syncStatus) VALUES (?, ?, ?)",
            harborRepo,
            swrRepo,
            syncStatus,
        )
        if err != nil {
            fmt.Println(err)
            dbRsp.Status = false
            dbRsp.Msg = "turn on repo sync failed"
        } else {
            dbRsp.Status = true
            dbRsp.Msg = "turn on repo sync successed"
        }
    } else {
        data := map[string]string{"harborRepo":harborRepo, "swrRepo":swrRepo}
        dbRsp = CheckRepoSync(data, db)
        status := dbRsp.SyncStatus
        switch status {
            case "off":
                _, err := db.Exec(
                    "UPDATE syncpool set syncStatus = 'on' where harborRepo = ? and swrRepo = ? and syncStatus = ?",
                    data["harborRepo"],
                    data["swrRepo"],
                    "off",
                )
                if err != nil {
                    fmt.Println(err)
                    dbRsp.Status = false
                    dbRsp.Msg = fmt.Sprintf("internal error, %s", err)
                } else {
                    dbRsp.Status = true
                    dbRsp.Msg = "turn on repo sync successed"
                }
            case "on":
                dbRsp.Status = false
                dbRsp.Msg = "repo sync has been added, status was 'on'"
        }
    }
    return dbRsp
}

func DelRepoSync(data map[string]string, db *sql.DB) (RESPONSEACTION) {
    var resp RESPONSEACTION
    result, err := db.Exec(
        "DELETE from syncpool where syncStatus = 'off' and harborRepo = ? and swrRepo = ?",
        data["harborRepo"],
        data["swrRepo"],
    )
    rowsaffected, _:= result.RowsAffected()
    if err != nil || rowsaffected == 0 {
        fmt.Println(err)
        fmt.Println(rowsaffected)
        resp.Status = false
        resp.Msg = "del repo sync record failed, did it has been turn on sync, turn it \"off\" first"
    } else {
        resp.Status = true
        resp.Msg = "del repo sync record successed"
    }
    return resp
}

func OffRepoSync(data map[string]string, db *sql.DB) (RESPONSEACTION) {
    isCheck := MQueryRow(data["harborRepo"], db)
    var resp RESPONSEACTION
    //若数据存在,为true
    if isCheck {
        _, err := db.Exec(
            "UPDATE syncpool set syncStatus = 'off' where harborRepo = ? and swrRepo = ? and syncStatus = ?",
            data["harborRepo"],
            data["swrRepo"],
            "on",
        )
        if err != nil {
            fmt.Println(err)
            resp.Status = false
            resp.Msg = "repo sync turn off failed"
        } else {
            resp.Status = true
            resp.Msg = "repo sync turn off successed"
        }
    } else {
        resp = CheckRepoSync(data, db)
    }
    return resp
}

func CheckRepoSync(data map[string]string, db *sql.DB) (RESPONSEACTION) {
    var s string
    var resp RESPONSEACTION
    err := db.QueryRow("SELECT syncStatus from syncpool where harborRepo = ?", data["harborRepo"]).Scan(&s)
    if err == nil {
        switch s {
            case "on":
                 resp.Status = true
                 resp.Msg = "repo sync has been added, status was 'on'"
                 resp.SyncStatus = "on"
            case "off":
                resp.Status = true
                resp.Msg = "repo sync hase been added, status was 'off'"
                 resp.SyncStatus = "off"
            default:
                resp.Status = false
                resp.Msg = "unknown mistake"
                resp.SyncStatus = s
            }
    } else {
        resp.Status = false
        resp.Msg = fmt.Sprintf("internal error, %s", err)
        fmt.Println(err)
    }
    return resp
}

type CREPO struct {
    Status bool
    Res map[string]string
}

func AllowSync(repo string, db *sql.DB) (CREPO) {
    repoMap := make(map[string]string)
    var crepo CREPO
    crepo.Status = false
    crepo.Res = repoMap

    var h, s string
    rows, err := db.Query("SELECT harborRepo, swrRepo from syncpool where syncStatus = ?", "on")
    if err == nil {
        for rows.Next() {
            err = rows.Scan(&h, &s)
            if err == nil {
                if repo == h {
                    crepo.Status = true
                    repoMap[h] = s
                    crepo.Res = repoMap
                }
            } else {
                fmt.Println(err)
            }
        }
    } else {
        fmt.Println(err)
    }
    return crepo
}
