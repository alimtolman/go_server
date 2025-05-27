package Tables

import "Server/DataBase"

var AppData = DataBase.NewJsonTable("AppData", DataBase.AppDataDefault())
